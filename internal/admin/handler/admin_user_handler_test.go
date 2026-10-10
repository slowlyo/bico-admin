package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bico-admin/internal/admin/model"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// TestAdminUserUsernameUpdate 验证用户管理可修改用户名（自动去除首尾空格）且不允许重复。
func TestAdminUserUsernameUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		// 数据库不可用时无法验证用户名查重。
		t.Fatalf("创建测试数据库失败: %v", err)
	}
	if err := database.AutoMigrate(&model.AdminUser{}, &model.AdminRole{}, &model.AdminUserRole{}); err != nil {
		// 用户表未建立时无法测试更新逻辑。
		t.Fatalf("迁移测试数据库失败: %v", err)
	}

	users := []model.AdminUser{
		{Username: "first-user", Password: "unused", Enabled: true},
		{Username: "existing-user", Password: "unused", Enabled: true},
	}
	if err := database.Create(&users).Error; err != nil {
		// 测试数据缺失时无法区分自身与重复用户名。
		t.Fatalf("创建测试用户失败: %v", err)
	}

	handler := NewAdminUserHandler(database, nil)

	// 未被占用的用户名应允许更新，并统一去除首尾空格。
	resp := updateUser(t, handler, "1", `{"username":"  renamed-user  "}`)
	if resp["code"].(float64) != 0 {
		t.Fatalf("未占用的用户名应允许更新，实际: %v", resp)
	}
	var renamed model.AdminUser
	if err := database.First(&renamed, users[0].ID).Error; err != nil || renamed.Username != "renamed-user" {
		t.Fatalf("用户名未正确标准化: %v, err=%v", renamed.Username, err)
	}

	// 被其他用户占用的用户名必须被拒绝。
	resp = updateUser(t, handler, "1", `{"username":"existing-user"}`)
	if resp["code"].(float64) != 400 || resp["msg"] != "用户名已存在" {
		t.Fatalf("期望用户名已存在错误，实际: %v", resp)
	}
}

// updateUser 通过 HTTP 入口执行用户更新并返回响应体。
func updateUser(t *testing.T, handler *AdminUserHandler, id string, body string) map[string]interface{} {
	t.Helper()

	req := httptest.NewRequest(http.MethodPut, "/admin-users/"+id, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: id}}
	handler.Update(c)

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("响应不是合法 JSON: %v", err)
	}
	return resp
}
