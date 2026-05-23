package crud

import (
	"encoding/json"
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
	h.BuildUpdates = func(req *testUpdateReq, existing *testCRUDModel) (map[string]interface{}, error) {
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
