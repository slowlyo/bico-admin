package handler

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"bico-admin/internal/admin/model"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// recordingInvalidator 记录缓存失效调用，用于断言各个 hook 触发了哪些失效。
type recordingInvalidator struct {
	userPermission []uint
	rolePermission []uint
	rolesUsers     [][]uint
	userStatus     []uint
	userToken      []uint
	usersAuth      [][]uint
}

func (r *recordingInvalidator) InvalidateUserPermissionCache(id uint) {
	r.userPermission = append(r.userPermission, id)
}

func (r *recordingInvalidator) InvalidateRoleUsersPermissionCache(id uint) {
	r.rolePermission = append(r.rolePermission, id)
}

func (r *recordingInvalidator) InvalidateRolesUsersPermissionCache(ids []uint) {
	r.rolesUsers = append(r.rolesUsers, ids)
}

func (r *recordingInvalidator) InvalidateUserStatusCache(id uint) {
	r.userStatus = append(r.userStatus, id)
}

func (r *recordingInvalidator) InvalidateUserTokenVersionCache(id uint) {
	r.userToken = append(r.userToken, id)
}

func (r *recordingInvalidator) InvalidateUsersAuthCache(ids []uint) {
	r.usersAuth = append(r.usersAuth, ids)
}

// newAdminTestDB 创建独立的内存数据库并迁移后台表。
//
// 说明：使用按测试名区分的共享缓存内存库，避免连接池里的多个连接各自拿到一个空库，
// 同时不同测试之间互不影响。
func newAdminTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := "file:" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("创建测试数据库失败: %v", err)
	}
	if err := db.AutoMigrate(
		&model.AdminUser{},
		&model.AdminRole{},
		&model.AdminRolePermission{},
		&model.AdminUserRole{},
	); err != nil {
		t.Fatalf("迁移测试数据库失败: %v", err)
	}
	return db
}

// callHandler 通过 HTTP 入口调用处理器方法并返回响应体；id 非空时作为路径参数。
func callHandler(t *testing.T, method, id, body string, handler gin.HandlerFunc) map[string]interface{} {
	return callHandlerAs(t, method, id, body, 0, handler)
}

// callHandlerAs 在 Gin 上下文中注入 operatorID 后调用处理器（0 表示不注入登录用户）。
func callHandlerAs(t *testing.T, method, id, body string, operatorID uint, handler gin.HandlerFunc) map[string]interface{} {
	t.Helper()

	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(method, "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	if operatorID != 0 {
		c.Set("user_id", operatorID)
	}
	if id != "" {
		c.Params = gin.Params{{Key: "id", Value: id}}
	}
	handler(c)

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("响应不是合法 JSON: %v, body=%s", err, w.Body.String())
	}
	return resp
}

// expectOK 断言响应成功。
func expectOK(t *testing.T, resp map[string]interface{}) {
	t.Helper()
	if resp["code"].(float64) != 0 {
		t.Fatalf("期望成功，实际: %v", resp)
	}
}

// expectBizError 断言响应是带指定文案的业务错误。
func expectBizError(t *testing.T, resp map[string]interface{}, msg string) {
	t.Helper()
	if resp["code"].(float64) != 400 || resp["msg"] != msg {
		t.Fatalf("期望业务错误 %q，实际: %v", msg, resp)
	}
}

// seedRole 创建角色并返回，code 需要在测试内唯一。
func seedRole(t *testing.T, db *gorm.DB, name, code string) model.AdminRole {
	t.Helper()
	role := model.AdminRole{Name: name, Code: code, Enabled: true}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("创建测试角色失败: %v", err)
	}
	return role
}

// seedUser 创建启用的用户并可选绑定角色。
func seedUser(t *testing.T, db *gorm.DB, username string, roles ...model.AdminRole) model.AdminUser {
	t.Helper()
	user := model.AdminUser{Username: username, Password: "old-hash", Enabled: true}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}
	for _, role := range roles {
		if err := db.Create(&model.AdminUserRole{UserID: user.ID, RoleID: role.ID}).Error; err != nil {
			t.Fatalf("绑定测试角色失败: %v", err)
		}
	}
	return user
}

// countRows 统计表中满足条件的行数。
func countRows(t *testing.T, db *gorm.DB, value interface{}, query string, args ...interface{}) int64 {
	t.Helper()
	var n int64
	if err := db.Model(value).Where(query, args...).Count(&n).Error; err != nil {
		t.Fatalf("统计失败: %v", err)
	}
	return n
}
