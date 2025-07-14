package handler

import (
	"github.com/gin-gonic/gin"

	"bico-admin/internal/admin/service"
	"bico-admin/internal/admin/types"
	sharedTypes "bico-admin/internal/shared/types"
	"bico-admin/pkg/response"
)

// AdminUserHandler 管理员用户处理器
type AdminUserHandler struct {
	adminUserService service.AdminUserService
}

// NewAdminUserHandler 创建管理员用户处理器
func NewAdminUserHandler(adminUserService service.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{
		adminUserService: adminUserService,
	}
}

// GetList 获取管理员用户列表
// @Summary 获取管理员用户列表
// @Description 分页获取管理员用户列表
// @Tags 管理员用户管理
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param username query string false "用户名"
// @Param name query string false "姓名"
// @Param status query int false "状态"
// @Success 200 {object} response.PageResponse{list=[]types.AdminUserResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /admin/admin-users [get]
func (h *AdminUserHandler) GetList(c *gin.Context) {
	var req types.AdminUserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	users, total, err := h.adminUserService.ListWithFilter(c.Request.Context(), &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	// 转换为响应格式
	var userResponses []types.AdminUserResponse
	for _, user := range users {
		userResponses = append(userResponses, types.AdminUserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Name:        user.Name,
			Avatar:      user.Avatar,
			Email:       user.Email,
			Phone:       user.Phone,
			Status:      user.Status,
			StatusText:  user.GetStatusText(),
			LastLoginAt: user.LastLoginAt,
			Remark:      user.Remark,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		})
	}

	response.Page(c, userResponses, total, req.Page, req.PageSize)
}

// GetByID 根据ID获取管理员用户
// @Summary 根据ID获取管理员用户
// @Description 根据用户ID获取管理员用户详细信息
// @Tags 管理员用户管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.ApiResponse{data=types.AdminUserResponse}
// @Failure 404 {object} response.ApiResponse
// @Router /admin/admin-users/{id} [get]
func (h *AdminUserHandler) GetByID(c *gin.Context) {
	var req sharedTypes.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	user, err := h.adminUserService.GetByID(c.Request.Context(), req.ID)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeNotFound, err.Error())
		return
	}

	userResponse := types.AdminUserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Avatar:      user.Avatar,
		Email:       user.Email,
		Phone:       user.Phone,
		Status:      user.Status,
		StatusText:  user.GetStatusText(),
		LastLoginAt: user.LastLoginAt,
		Remark:      user.Remark,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	response.Success(c, userResponse)
}

// Create 创建管理员用户
// @Summary 创建管理员用户
// @Description 创建新的管理员用户
// @Tags 管理员用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body types.AdminUserCreateRequest true "创建管理员用户请求"
// @Success 200 {object} response.ApiResponse{data=types.AdminUserResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /admin/admin-users [post]
func (h *AdminUserHandler) Create(c *gin.Context) {
	var req types.AdminUserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	user, err := h.adminUserService.Create(c.Request.Context(), &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	userResponse := types.AdminUserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Avatar:      user.Avatar,
		Email:       user.Email,
		Phone:       user.Phone,
		Status:      user.Status,
		StatusText:  user.GetStatusText(),
		LastLoginAt: user.LastLoginAt,
		Remark:      user.Remark,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	response.Success(c, userResponse)
}

// Update 更新管理员用户
// @Summary 更新管理员用户
// @Description 更新管理员用户信息
// @Tags 管理员用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param request body types.AdminUserUpdateRequest true "更新管理员用户请求"
// @Success 200 {object} response.ApiResponse{data=types.AdminUserResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /admin/admin-users/{id} [put]
func (h *AdminUserHandler) Update(c *gin.Context) {
	var uriReq sharedTypes.IDRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	var req types.AdminUserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	user, err := h.adminUserService.Update(c.Request.Context(), uriReq.ID, &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	userResponse := types.AdminUserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Avatar:      user.Avatar,
		Email:       user.Email,
		Phone:       user.Phone,
		Status:      user.Status,
		StatusText:  user.GetStatusText(),
		LastLoginAt: user.LastLoginAt,
		Remark:      user.Remark,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	response.Success(c, userResponse)
}

// Delete 删除管理员用户
// @Summary 删除管理员用户
// @Description 软删除管理员用户
// @Tags 管理员用户管理
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.ApiResponse
// @Failure 404 {object} response.ApiResponse
// @Router /admin/admin-users/{id} [delete]
func (h *AdminUserHandler) Delete(c *gin.Context) {
	var req sharedTypes.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.adminUserService.Delete(c.Request.Context(), req.ID); err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateStatus 更新管理员用户状态
// @Summary 更新管理员用户状态
// @Description 启用或禁用管理员用户
// @Tags 管理员用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param request body types.StatusRequest true "状态更新请求"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /admin/admin-users/{id}/status [patch]
func (h *AdminUserHandler) UpdateStatus(c *gin.Context) {
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

	enabled := req.Status == sharedTypes.StatusActive
	if err := h.adminUserService.UpdateStatus(c.Request.Context(), uriReq.ID, enabled); err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
