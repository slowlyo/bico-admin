package handler

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"bico-admin/core/middleware"
	"bico-admin/modules/api/service"
	"bico-admin/pkg/response"
)

// AppHandler 应用处理器
type AppHandler struct {
	appService service.AppService
}

// NewAppHandler 创建应用处理器实例
func NewAppHandler(db *gorm.DB) *AppHandler {
	return &AppHandler{
		appService: service.NewAppService(db),
	}
}

// GetAppInfo 获取应用信息
func (h *AppHandler) GetAppInfo(c *fiber.Ctx) error {
	info, err := h.appService.GetAppInfo()
	if err != nil {
		return response.InternalServerError(c, "Failed to get app info")
	}

	return response.Success(c, info)
}

// GetAppConfig 获取应用配置
func (h *AppHandler) GetAppConfig(c *fiber.Ctx) error {
	config, err := h.appService.GetAppConfig()
	if err != nil {
		return response.InternalServerError(c, "Failed to get app config")
	}

	return response.Success(c, config)
}

// GetUserProfile 获取用户资料
func (h *AppHandler) GetUserProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	profile, err := h.appService.GetUserProfile(userID)
	if err != nil {
		return response.NotFound(c, "User not found")
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

	profile, err := h.appService.UpdateUserProfile(userID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "Profile updated successfully", profile)
}

// GetContentList 获取内容列表
func (h *AppHandler) GetContentList(c *fiber.Ctx) error {
	// 获取查询参数
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)
	category := c.Query("category", "")

	content, err := h.appService.GetContentList(page, pageSize, category)
	if err != nil {
		return response.InternalServerError(c, "Failed to get content list")
	}

	return response.Success(c, content)
}

// GetContent 获取单个内容
func (h *AppHandler) GetContent(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid content ID")
	}

	content, err := h.appService.GetContent(uint(id))
	if err != nil {
		return response.NotFound(c, "Content not found")
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
