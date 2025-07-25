package handler

import (
	"github.com/gin-gonic/gin"

	"bico-admin/internal/admin/models"
	"bico-admin/internal/admin/service"
	"bico-admin/internal/admin/types"
	sharedTypes "bico-admin/internal/shared/types"
	"bico-admin/pkg/response"
	"bico-admin/pkg/utils"
)

// AdminUserHandler 管理员用户处理器
type AdminUserHandler struct {
	adminUserService service.AdminUserService
	adminRoleService service.AdminRoleService
}

// NewAdminUserHandler 创建管理员用户处理器
func NewAdminUserHandler(adminUserService service.AdminUserService, adminRoleService service.AdminRoleService) *AdminUserHandler {
	return &AdminUserHandler{
		adminUserService: adminUserService,
		adminRoleService: adminRoleService,
	}
}

// GetList 获取管理员用户列表
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
		userResponses = append(userResponses, h.convertToAdminUserResponse(c, user))
	}

	response.Page(c, userResponses, total, req.GetPage(), req.GetPageSize())
}

// GetByID 根据ID获取管理员用户
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

	userResponse := h.convertToAdminUserResponse(c, user)
	response.Success(c, userResponse)
}

// convertToAdminUserResponse 转换为管理员用户响应
func (h *AdminUserHandler) convertToAdminUserResponse(ctx *gin.Context, user *models.AdminUser) types.AdminUserResponse {
	// 检查权限
	canDelete, _ := h.adminUserService.CanUserBeDeleted(ctx.Request.Context(), user.ID)
	canDisable, _ := h.adminUserService.CanUserBeDisabled(ctx.Request.Context(), user.ID)

	// 获取用户角色
	var roles []types.AdminUserRoleResponse
	for _, userRole := range user.Roles {
		roles = append(roles, types.AdminUserRoleResponse{
			ID:          userRole.Role.ID,
			Name:        userRole.Role.Name,
			Code:        userRole.Role.Code,
			Description: userRole.Role.Description,
		})
	}

	// 转换LastLoginAt
	var lastLoginAt *utils.FormattedTime
	if user.LastLoginAt != nil {
		ft := utils.NewFormattedTime(*user.LastLoginAt)
		lastLoginAt = &ft
	}

	status := 0
	if user.Status != nil {
		status = *user.Status
	}

	return types.AdminUserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Avatar:      user.Avatar,
		Email:       user.Email,
		Phone:       user.Phone,
		Status:      status,
		StatusText:  user.GetStatusText(),
		LastLoginAt: lastLoginAt,
		Remark:      user.Remark,
		CanDelete:   canDelete,
		CanDisable:  canDisable,
		Roles:       roles,
		CreatedAt:   utils.NewFormattedTime(user.CreatedAt),
		UpdatedAt:   utils.NewFormattedTime(user.UpdatedAt),
	}
}

// Create 创建管理员用户
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

	userResponse := h.convertToAdminUserResponse(c, user)
	response.Success(c, userResponse)
}

// Update 更新管理员用户
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

	userResponse := h.convertToAdminUserResponse(c, user)
	response.Success(c, userResponse)
}

// Delete 删除管理员用户
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
