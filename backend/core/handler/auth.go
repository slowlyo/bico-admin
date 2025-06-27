package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"bico-admin/core/middleware"
	"bico-admin/core/model"
	"bico-admin/core/service"
	"bico-admin/pkg/response"
	"bico-admin/pkg/validator"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	// TODO: 从配置中获取JWT密钥和过期时间
	jwtSecret := "your-secret-key"
	jwtExpire := 24 * time.Hour

	return &AuthHandler{
		authService: service.NewAuthService(db, jwtSecret, jwtExpire),
	}
}

// Login 用户登录
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 验证请求参数
	if errors := validator.Validate(req); len(errors) > 0 {
		return response.ValidationError(c, errors)
	}

	// 执行登录
	user, token, err := h.authService.Login(req)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	return response.Success(c, fiber.Map{
		"user":  user,
		"token": token,
	})
}

// Register 用户注册
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req model.UserCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 验证请求参数
	if errors := validator.Validate(req); len(errors) > 0 {
		return response.ValidationError(c, errors)
	}

	// 执行注册
	user, err := h.authService.Register(req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "User registered successfully", user)
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// TODO: 实现token黑名单机制
	return response.SuccessWithMessage(c, "Logout successful", nil)
}

// GetProfile 获取用户资料
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	user, err := h.authService.GetProfile(userID)
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	return response.Success(c, user)
}

// UpdateProfile 更新用户资料
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	var req model.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 验证请求参数
	if errors := validator.Validate(req); len(errors) > 0 {
		return response.ValidationError(c, errors)
	}

	// 更新用户资料
	user, err := h.authService.UpdateProfile(userID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "Profile updated successfully", user)
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	var req struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=6"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 验证请求参数
	if errors := validator.Validate(req); len(errors) > 0 {
		return response.ValidationError(c, errors)
	}

	// 修改密码
	if err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.SuccessWithMessage(c, "Password changed successfully", nil)
}
