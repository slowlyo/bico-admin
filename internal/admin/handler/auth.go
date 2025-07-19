package handler

import (
	"github.com/gin-gonic/gin"

	"bico-admin/internal/admin/service"
	"bico-admin/internal/admin/types"
	"bico-admin/pkg/response"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 管理员登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req types.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeUnauthorized, err.Error())
		return
	}

	response.Success(c, result)
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从上下文获取令牌
	token, exists := c.Get("token")
	if !exists {
		response.Unauthorized(c, "未找到认证令牌")
		return
	}

	if err := h.authService.Logout(c.Request.Context(), token.(string)); err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// RefreshToken 刷新令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req types.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeUnauthorized, err.Error())
		return
	}

	response.Success(c, result)
}

// GetProfile 获取当前用户信息和权限
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "用户未登录")
		return
	}

	profile, err := h.authService.GetProfileWithPermissions(c.Request.Context(), userID.(uint))
	if err != nil {
		response.ErrorWithMessage(c, response.CodeNotFound, err.Error())
		return
	}

	response.Success(c, profile)
}

// UpdateProfile 更新当前用户信息
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "用户未登录")
		return
	}

	var req types.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.authService.UpdateProfileInfo(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "用户未登录")
		return
	}

	var req types.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.authService.ChangePassword(c.Request.Context(), userID.(uint), &req); err != nil {
		response.ErrorWithMessage(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}
