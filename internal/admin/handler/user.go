package handler

import (
	"github.com/gin-gonic/gin"

	"bico-admin/internal/admin/service"
	"bico-admin/internal/admin/types"
	sharedTypes "bico-admin/internal/shared/types"
	"bico-admin/pkg/response"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetList 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Param user_type query string false "用户类型"
// @Param status query int false "状态"
// @Success 200 {object} response.ApiResponse{data=response.PageResponse}
// @Router /admin/users [get]
func (h *UserHandler) GetList(c *gin.Context) {
	var req types.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.userService.GetList(c.Request.Context(), &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	response.Page(c, result.List, result.Total, result.Page, result.PageSize)
}

// GetByID 根据ID获取用户
// @Summary 根据ID获取用户
// @Description 根据用户ID获取用户详细信息
// @Tags 用户管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.ApiResponse{data=types.UserResponse}
// @Failure 404 {object} response.ApiResponse
// @Router /admin/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	var req sharedTypes.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), req.ID)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeNotFound, err.Error())
		return
	}

	response.Success(c, user)
}

// Create 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body types.UserCreateRequest true "创建用户请求"
// @Success 200 {object} response.ApiResponse{data=types.UserResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /admin/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req types.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	user, err := h.userService.Create(c.Request.Context(), &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(c, user)
}

// Update 更新用户
// @Summary 更新用户
// @Description 更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param request body types.UserUpdateRequest true "更新用户请求"
// @Success 200 {object} response.ApiResponse{data=types.UserResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /admin/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	var uriReq sharedTypes.IDRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	var req types.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	user, err := h.userService.Update(c.Request.Context(), uriReq.ID, &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(c, user)
}

// Delete 删除用户
// @Summary 删除用户
// @Description 软删除用户
// @Tags 用户管理
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Router /admin/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	var req sharedTypes.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.userService.Delete(c.Request.Context(), req.ID); err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateStatus 更新用户状态
// @Summary 更新用户状态
// @Description 更新用户状态（激活/禁用）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param request body object{status=int} true "状态更新请求" example({"status": 1})
// @Success 200 {object} response.ApiResponse
// @Router /admin/users/{id}/status [patch]
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	var uriReq sharedTypes.IDRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	var req types.StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.userService.UpdateStatus(c.Request.Context(), uriReq.ID, req.Status); err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// ResetPassword 重置用户密码
// @Summary 重置用户密码
// @Description 管理员重置用户密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param request body types.UserPasswordRequest true "密码重置请求"
// @Success 200 {object} response.ApiResponse
// @Router /admin/users/{id}/password [patch]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var uriReq sharedTypes.IDRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	var req types.UserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.userService.ResetPassword(c.Request.Context(), uriReq.ID, req.Password); err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetStats 获取用户统计
// @Summary 获取用户统计
// @Description 获取用户统计信息
// @Tags 用户管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.ApiResponse{data=types.UserStatsResponse}
// @Router /admin/users/stats [get]
func (h *UserHandler) GetStats(c *gin.Context) {
	stats, err := h.userService.GetStats(c.Request.Context())
	if err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(c, stats)
}
