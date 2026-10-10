package crud

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type testCRUDModel struct {
	ID      uint   `gorm:"primarykey" json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type testListReq struct {
	Enabled *bool `form:"enabled"`
}

type testCreateReq struct {
	Name string `json:"name" binding:"required"`
}

type testUpdateReq struct {
	Name    string `json:"name"`
	Enabled *bool  `json:"enabled"`
}

// newTestCRUDHandler 创建覆盖 CRUD 通用逻辑的测试 handler。
func newTestCRUDHandler(t *testing.T) (*CRUDHandler[testCRUDModel, testListReq, testCreateReq, testUpdateReq], *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("创建测试数据库失败: %v", err)
	}
	if err := db.AutoMigrate(&testCRUDModel{}); err != nil {
		t.Fatalf("迁移测试表失败: %v", err)
	}

	h := &CRUDHandler[testCRUDModel, testListReq, testCreateReq, testUpdateReq]{
		DB:          db,
		NotFoundMsg: "测试记录不存在",
	}
	h.BuildListQuery = func(db *gorm.DB, req *testListReq) *gorm.DB {
		query := db.Model(&testCRUDModel{})
		if req.Enabled != nil {
			query = query.Where("enabled = ?", *req.Enabled)
		}
		return query
	}
	h.NewModelFromCreate = func(req *testCreateReq) (*testCRUDModel, error) {
		return &testCRUDModel{Name: req.Name, Enabled: true}, nil
	}
	h.BuildUpdates = func(tx *gorm.DB, req *testUpdateReq, existing *testCRUDModel) (map[string]interface{}, error) {
		updates := map[string]interface{}{}
		if req.Name != "" {
			updates["name"] = req.Name
		}
		if req.Enabled != nil {
			updates["enabled"] = *req.Enabled
		}
		return updates, nil
	}

	return h, db
}

// performRequest 执行 handler 并返回响应记录器。
func performRequest(method string, path string, body string, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler(c)
	return w
}

// decodeResponse 解析统一响应体。
func decodeResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("响应不是合法 JSON: %v", err)
	}
	return resp
}

// TestCRUDListReturnsBadRequestOnInvalidQuery 验证 query 绑定失败不再静默忽略。
func TestCRUDListReturnsBadRequestOnInvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, _ := newTestCRUDHandler(t)

	w := performRequest(http.MethodGet, "/test?enabled=bad", "", h.List)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("非法 query 应返回 400，实际: %d", w.Code)
	}
}

// TestCRUDUpdateEnabledUsesUpdateHooks 验证启用状态更新复用标准更新生命周期。
func TestCRUDUpdateEnabledUsesUpdateHooks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, db := newTestCRUDHandler(t)

	item := testCRUDModel{Name: "test", Enabled: true}
	if err := db.Create(&item).Error; err != nil {
		t.Fatalf("创建测试数据失败: %v", err)
	}

	h.EnabledField = "enabled"
	h.BeforeUpdate = func(tx *gorm.DB, id uint, existing *testCRUDModel, req *testUpdateReq, updates map[string]interface{}) error {
		updates["name"] = "hooked"
		return nil
	}

	w := performRequest(http.MethodPatch, "/test/1", `{"enabled":false}`, func(c *gin.Context) {
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		h.UpdateEnabled(c)
	})
	if w.Code != http.StatusOK {
		t.Fatalf("更新启用状态应成功，实际: %d, body: %s", w.Code, w.Body.String())
	}

	var updated testCRUDModel
	if err := db.First(&updated, item.ID).Error; err != nil {
		t.Fatalf("读取更新后记录失败: %v", err)
	}
	if updated.Enabled || updated.Name != "hooked" {
		t.Fatalf("UpdateEnabled 未复用更新 hook，记录: %+v", updated)
	}
}

// TestCRUDDeleteBatchRequiresAllIDs 验证批量删除必须全部命中。
func TestCRUDDeleteBatchRequiresAllIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, db := newTestCRUDHandler(t)

	item := testCRUDModel{Name: "test", Enabled: true}
	if err := db.Create(&item).Error; err != nil {
		t.Fatalf("创建测试数据失败: %v", err)
	}

	w := performRequest(http.MethodDelete, "/test/batch", `{"ids":[1,999]}`, h.DeleteBatch)
	resp := decodeResponse(t, w)

	if w.Code != http.StatusOK || resp["code"].(float64) != 404 {
		t.Fatalf("部分 ID 不存在应返回业务 404，status=%d body=%s", w.Code, w.Body.String())
	}

	var count int64
	if err := db.Model(&testCRUDModel{}).Where("id = ?", item.ID).Count(&count).Error; err != nil {
		t.Fatalf("统计测试记录失败: %v", err)
	}
	if count != 1 {
		t.Fatalf("批量删除部分失败时应回滚，剩余记录数: %d", count)
	}
}

// seedSortItems 写入用于排序断言的测试数据。
func seedSortItems(t *testing.T, db *gorm.DB) {
	t.Helper()

	items := []testCRUDModel{{Name: "b"}, {Name: "a"}, {Name: "c"}}
	if err := db.Create(&items).Error; err != nil {
		t.Fatalf("创建测试数据失败: %v", err)
	}
}

// listNames 请求列表并按返回顺序提取 name。
func listNames(t *testing.T, h *CRUDHandler[testCRUDModel, testListReq, testCreateReq, testUpdateReq], query string) []string {
	t.Helper()

	w := performRequest(http.MethodGet, "/test?"+query, "", h.List)
	resp := decodeResponse(t, w)
	data := resp["data"].(map[string]interface{})
	list := data["list"].([]interface{})
	names := make([]string, 0, len(list))
	for _, item := range list {
		names = append(names, item.(map[string]interface{})["name"].(string))
	}
	return names
}

// TestCRUDListSortWhitelist 验证只有白名单字段允许客户端排序，其余回退默认排序。
func TestCRUDListSortWhitelist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, db := newTestCRUDHandler(t)
	seedSortItems(t, db)
	h.DefaultOrder = "id ASC"
	h.SortableFields = []string{"name"}

	// 白名单内字段按客户端要求排序。
	if got := strings.Join(listNames(t, h, "sortField=name&sortOrder=ascend"), ""); got != "abc" {
		t.Fatalf("白名单字段应生效，实际: %s", got)
	}
	// 不在白名单内的字段被忽略，回退默认排序。
	if got := strings.Join(listNames(t, h, "sortField=enabled&sortOrder=descend"), ""); got != "bac" {
		t.Fatalf("非白名单字段应回退默认排序，实际: %s", got)
	}

	// 未配置白名单时一律忽略客户端排序。
	h.SortableFields = nil
	if got := strings.Join(listNames(t, h, "sortField=name&sortOrder=ascend"), ""); got != "bac" {
		t.Fatalf("未配置白名单应使用默认排序，实际: %s", got)
	}
}

// TestCRUDCommitHooksRunAfterSuccess 验证提交后钩子只在事务成功后执行。
func TestCRUDCommitHooksRunAfterSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, db := newTestCRUDHandler(t)

	var created []string
	var deleted [][]uint
	h.AfterCreateCommit = func(item *testCRUDModel, req *testCreateReq) {
		created = append(created, item.Name)
	}
	h.AfterDeleteCommit = func(ids []uint) {
		deleted = append(deleted, ids)
	}

	// 创建成功后触发一次创建钩子。
	performRequest(http.MethodPost, "/test", `{"name":"x"}`, h.Create)
	if len(created) != 1 || created[0] != "x" {
		t.Fatalf("创建成功应触发 AfterCreateCommit，实际: %v", created)
	}

	// 创建事务失败时不触发提交后钩子。
	h.BeforeCreate = func(tx *gorm.DB, item *testCRUDModel, req *testCreateReq) error {
		return gorm.ErrInvalidData
	}
	performRequest(http.MethodPost, "/test", `{"name":"y"}`, h.Create)
	if len(created) != 1 {
		t.Fatalf("创建失败不应触发 AfterCreateCommit，实际: %v", created)
	}

	// 单个删除以单元素 ids 触发钩子。
	performRequest(http.MethodDelete, "/test/1", "", func(c *gin.Context) {
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		h.Delete(c)
	})
	if len(deleted) != 1 || len(deleted[0]) != 1 || deleted[0][0] != 1 {
		t.Fatalf("单个删除应触发 AfterDeleteCommit，实际: %v", deleted)
	}

	// 批量删除部分不存在时整体回滚，不触发钩子。
	performRequest(http.MethodDelete, "/test/batch", `{"ids":[1,999]}`, h.DeleteBatch)
	if len(deleted) != 1 {
		t.Fatalf("批量删除失败不应触发 AfterDeleteCommit，实际: %v", deleted)
	}

	// 批量删除成功时携带去重后的全部 id。
	items := []testCRUDModel{{Name: "p"}, {Name: "q"}}
	if err := db.Create(&items).Error; err != nil {
		t.Fatalf("创建测试数据失败: %v", err)
	}
	body := fmt.Sprintf(`{"ids":[%d,%d,%d]}`, items[0].ID, items[1].ID, items[1].ID)
	performRequest(http.MethodDelete, "/test/batch", body, h.DeleteBatch)
	if len(deleted) != 2 || len(deleted[1]) != 2 {
		t.Fatalf("批量删除成功应携带去重后的 ids，实际: %v", deleted)
	}
}

// requestWithID 以指定路由 ID 执行 handler。
func requestWithID(method string, id string, body string, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	return performRequest(method, "/test/"+id, body, func(c *gin.Context) {
		c.Params = gin.Params{{Key: "id", Value: id}}
		handler(c)
	})
}

// TestCRUDUniqueFields 验证唯一列在创建和更新时自动查重，更新时排除自身。
func TestCRUDUniqueFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, db := newTestCRUDHandler(t)
	h.UniqueFields = []Unique{{Column: "name", Err: NewBizError("名称已存在")}}

	if err := db.Create(&[]testCRUDModel{{Name: "a"}, {Name: "b"}}).Error; err != nil {
		t.Fatalf("创建测试数据失败: %v", err)
	}

	// 创建重名记录被拒绝，且返回业务错误文案。
	resp := decodeResponse(t, performRequest(http.MethodPost, "/test", `{"name":"a"}`, h.Create))
	if resp["code"].(float64) != 400 || resp["msg"] != "名称已存在" {
		t.Fatalf("创建重名应被拒绝，实际: %v", resp)
	}

	// 更新为其他记录占用的名称被拒绝。
	resp = decodeResponse(t, requestWithID(http.MethodPut, "1", `{"name":"b"}`, h.Update))
	if resp["code"].(float64) != 400 {
		t.Fatalf("更新为已占用名称应被拒绝，实际: %v", resp)
	}

	// 更新为自己当前的名称不算冲突。
	resp = decodeResponse(t, requestWithID(http.MethodPut, "1", `{"name":"a"}`, h.Update))
	if resp["code"].(float64) != 0 {
		t.Fatalf("保持自身名称应允许更新，实际: %v", resp)
	}

	// 不冲突的新名称可以创建。
	resp = decodeResponse(t, performRequest(http.MethodPost, "/test", `{"name":"c"}`, h.Create))
	if resp["code"].(float64) != 0 {
		t.Fatalf("新名称应允许创建，实际: %v", resp)
	}
}

// TestCRUDFailHidesSystemErrors 验证只有业务错误会展示给客户端，系统错误只返回通用提示。
func TestCRUDFailHidesSystemErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, _ := newTestCRUDHandler(t)

	// 业务错误原样返回。
	h.BeforeCreate = func(tx *gorm.DB, item *testCRUDModel, req *testCreateReq) error {
		return NewBizError("可以展示的提示")
	}
	resp := decodeResponse(t, performRequest(http.MethodPost, "/test", `{"name":"x"}`, h.Create))
	if resp["code"].(float64) != 400 || resp["msg"] != "可以展示的提示" {
		t.Fatalf("业务错误应原样返回，实际: %v", resp)
	}

	// 系统错误不暴露内部细节。
	h.BeforeCreate = func(tx *gorm.DB, item *testCRUDModel, req *testCreateReq) error {
		return errors.New("SELECT secret FROM internal_table failed")
	}
	w := performRequest(http.MethodPost, "/test", `{"name":"x"}`, h.Create)
	resp = decodeResponse(t, w)
	if resp["code"].(float64) != 500 || strings.Contains(w.Body.String(), "internal_table") {
		t.Fatalf("系统错误不应暴露细节，实际: %s", w.Body.String())
	}
}

// TestCRUDWriteReloadsDetail 验证配置详情加载后，创建和更新返回与详情一致的数据。
func TestCRUDWriteReloadsDetail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, db := newTestCRUDHandler(t)
	h.AfterGet = func(item *testCRUDModel) error {
		item.Name += "!"
		return nil
	}

	resp := decodeResponse(t, performRequest(http.MethodPost, "/test", `{"name":"x"}`, h.Create))
	if resp["data"].(map[string]interface{})["name"] != "x!" {
		t.Fatalf("创建应返回详情口径数据，实际: %v", resp)
	}

	item := testCRUDModel{Name: "y"}
	if err := db.Create(&item).Error; err != nil {
		t.Fatalf("创建测试数据失败: %v", err)
	}
	resp = decodeResponse(t, requestWithID(http.MethodPut, fmt.Sprint(item.ID), `{"name":"z"}`, h.Update))
	if resp["data"].(map[string]interface{})["name"] != "z!" {
		t.Fatalf("更新应返回详情口径数据，实际: %v", resp)
	}
}

// TestCRUDAfterUpdateCommitReceivesUpdates 验证 UpdateEnabled 路径下提交后钩子可依据 updates 判断变更。
func TestCRUDAfterUpdateCommitReceivesUpdates(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, db := newTestCRUDHandler(t)

	item := testCRUDModel{Name: "test", Enabled: true}
	if err := db.Create(&item).Error; err != nil {
		t.Fatalf("创建测试数据失败: %v", err)
	}

	var got map[string]interface{}
	h.AfterUpdateCommit = func(id uint, existing *testCRUDModel, req *testUpdateReq, updates map[string]interface{}) {
		got = updates
	}
	requestWithID(http.MethodPatch, "1", `{"enabled":false}`, h.UpdateEnabled)

	if enabled, ok := got["enabled"].(bool); !ok || enabled {
		t.Fatalf("提交后钩子应收到实际写入的字段，实际: %v", got)
	}
}

// TestCRUDListWithoutBuildQuery 验证未配置列表查询时查询整张表。
func TestCRUDListWithoutBuildQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, db := newTestCRUDHandler(t)
	seedSortItems(t, db)
	h.BuildListQuery = nil
	h.DefaultOrder = "id ASC"

	if got := strings.Join(listNames(t, h, ""), ""); got != "bac" {
		t.Fatalf("未配置列表查询应返回整张表，实际: %s", got)
	}
}
