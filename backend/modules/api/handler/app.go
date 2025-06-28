package handler

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"bico-admin/core/middleware"
	"bico-admin/pkg/response"
)

// AppHandler 应用处理器 - 业务逻辑直接在handler中实现
type AppHandler struct {
	db *gorm.DB
}

// NewAppHandler 创建应用处理器实例
func NewAppHandler(db *gorm.DB) *AppHandler {
	return &AppHandler{
		db: db,
	}
}

// GetAppInfo 获取应用信息
func (h *AppHandler) GetAppInfo(c *fiber.Ctx) error {
	// 简化实现 - 返回基本应用信息
	info := fiber.Map{
		"name":        "Bico Admin",
		"version":     "1.0.0",
		"description": "AI友好的管理后台框架",
	}

	return response.Success(c, info)
}

// GetAppConfig 获取应用配置
func (h *AppHandler) GetAppConfig(c *fiber.Ctx) error {
	// 简化实现 - 返回基本配置信息
	config := fiber.Map{
		"theme":    "default",
		"language": "zh-CN",
		"timezone": "Asia/Shanghai",
	}

	return response.Success(c, config)
}

// GetUserProfile 获取用户资料
func (h *AppHandler) GetUserProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	// 简化实现 - 返回基本用户信息
	profile := fiber.Map{
		"id":       userID,
		"username": "user",
		"email":    "user@example.com",
		"nickname": "用户",
	}

	return response.Success(c, profile)
}

// UpdateUserProfile 更新用户资料
func (h *AppHandler) UpdateUserProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 简化实现 - 返回更新成功消息
	return response.SuccessWithMessage(c, "Profile updated successfully", fiber.Map{
		"id":       userID,
		"username": "user",
		"email":    "user@example.com",
		"nickname": "用户",
	})
}

// GetContentList 获取内容列表
func (h *AppHandler) GetContentList(c *fiber.Ctx) error {
	// 获取查询参数
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)
	category := c.Query("category", "")

	// 简化实现 - 返回示例内容列表
	content := fiber.Map{
		"data": []fiber.Map{
			{
				"id":       1,
				"title":    "示例内容1",
				"category": category,
				"content":  "这是示例内容1",
			},
			{
				"id":       2,
				"title":    "示例内容2",
				"category": category,
				"content":  "这是示例内容2",
			},
		},
		"total":       2,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": 1,
	}

	return response.Success(c, content)
}

// GetContent 获取单个内容
func (h *AppHandler) GetContent(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid content ID")
	}

	// 简化实现 - 返回示例内容
	content := fiber.Map{
		"id":      id,
		"title":   "示例内容",
		"content": "这是示例内容的详细信息",
		"author":  "系统管理员",
	}

	return response.Success(c, content)
}

// HealthCheck 健康检查
func (h *AppHandler) HealthCheck(c *fiber.Ctx) error {
	return response.Success(c, fiber.Map{
		"status":  "ok",
		"message": "API is running",
		"time":    fiber.Map{},
	})
}

// GetVersion 获取版本信息
func (h *AppHandler) GetVersion(c *fiber.Ctx) error {
	return response.Success(c, fiber.Map{
		"version": "1.0.0",
		"name":    "Bico Admin API",
		"build":   "20240101",
	})
}
