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
// @Summary 管理员登录
// @Description 管理员登录接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body types.AdminLoginRequest true "登录请求"
// @Success 200 {object} response.ApiResponse{data=types.AdminLoginResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /admin/auth/login [post]
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
// @Summary 登出
// @Description 管理员登出接口
// @Tags 认证
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.ApiResponse
// @Failure 401 {object} response.ApiResponse
// @Router /admin/auth/logout [post]
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
// @Summary 刷新令牌
// @Description 刷新访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body types.RefreshTokenRequest true "刷新令牌请求"
// @Success 200 {object} response.ApiResponse{data=types.AdminLoginResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /admin/auth/refresh [post]
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

// GetProfile 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.ApiResponse{data=types.AdminUserResponse}
// @Failure 401 {object} response.ApiResponse
// @Router /admin/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "用户未登录")
		return
	}

	profile, err := h.authService.GetProfile(c.Request.Context(), userID.(uint))
	if err != nil {
		response.ErrorWithMessage(c, response.CodeNotFound, err.Error())
		return
	}

	response.Success(c, profile)
}

// UpdateProfile 更新当前用户信息
// @Summary 更新当前用户信息
// @Description 更新当前登录用户的信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body types.AdminUserUpdateRequest true "更新用户请求"
// @Success 200 {object} response.ApiResponse{data=types.AdminUserResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /admin/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "用户未登录")
		return
	}

	var req types.AdminUserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.authService.UpdateProfile(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}
