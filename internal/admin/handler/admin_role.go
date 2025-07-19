package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"bico-admin/internal/admin/service"
	"bico-admin/internal/admin/types"
	"bico-admin/pkg/response"
)

// AdminRoleHandler 管理员角色处理器
type AdminRoleHandler struct {
	roleService service.AdminRoleService
}

// NewAdminRoleHandler 创建管理员角色处理器
func NewAdminRoleHandler(roleService service.AdminRoleService) *AdminRoleHandler {
	return &AdminRoleHandler{
		roleService: roleService,
	}
}

// GetRoleList 获取角色列表
func (h *AdminRoleHandler) GetRoleList(ctx *gin.Context) {
	var req types.RoleListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMessage(ctx, response.CodeBadRequest, err.Error())
		return
	}

	result, err := h.roleService.GetRoleList(ctx.Request.Context(), &req)
	if err != nil {
		response.ErrorWithMessage(ctx, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(ctx, result)
}

// CreateRole 创建角色
func (h *AdminRoleHandler) CreateRole(ctx *gin.Context) {
	var req types.RoleCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(ctx, response.CodeBadRequest, err.Error())
		return
	}

	result, err := h.roleService.CreateRole(ctx.Request.Context(), &req)
	if err != nil {
		response.ErrorWithMessage(ctx, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(ctx, result)
}

// UpdateRole 更新角色
func (h *AdminRoleHandler) UpdateRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.CodeBadRequest, "无效的角色ID")
		return
	}

	var req types.RoleUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(ctx, response.CodeBadRequest, err.Error())
		return
	}

	result, err := h.roleService.UpdateRole(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		response.ErrorWithMessage(ctx, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(ctx, result)
}

// UpdateRoleStatus 更新角色状态
func (h *AdminRoleHandler) UpdateRoleStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.CodeBadRequest, "无效的角色ID")
		return
	}

	var req types.StatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(ctx, response.CodeBadRequest, err.Error())
		return
	}

	// 获取角色信息
	role, err := h.roleService.GetRoleByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.ErrorWithMessage(ctx, response.CodeNotFound, err.Error())
		return
	}

	// 检查是否可以编辑
	if !role.CanEdit {
		response.ErrorWithMessage(ctx, response.CodeForbidden, "该角色不可编辑")
		return
	}

	// 构造更新请求
	updateReq := types.RoleUpdateRequest{
		Name:        role.Name,
		Description: role.Description,
		Status:      req.Status,
		Permissions: make([]string, len(role.Permissions)),
	}

	// 提取权限代码
	for i, perm := range role.Permissions {
		updateReq.Permissions[i] = perm.PermissionCode
	}

	// 更新角色
	if _, err := h.roleService.UpdateRole(ctx.Request.Context(), uint(id), &updateReq); err != nil {
		response.ErrorWithMessage(ctx, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// DeleteRole 删除角色
func (h *AdminRoleHandler) DeleteRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.CodeBadRequest, "无效的角色ID")
		return
	}

	if err := h.roleService.DeleteRole(ctx.Request.Context(), uint(id)); err != nil {
		response.ErrorWithMessage(ctx, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// GetRoleByID 根据ID获取角色
func (h *AdminRoleHandler) GetRoleByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.CodeBadRequest, "无效的角色ID")
		return
	}

	result, err := h.roleService.GetRoleByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.ErrorWithMessage(ctx, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(ctx, result)
}

// GetPermissionTree 获取权限树
func (h *AdminRoleHandler) GetPermissionTree(ctx *gin.Context) {
	var roleID *uint
	if roleIDStr := ctx.Query("role_id"); roleIDStr != "" {
		if id, err := strconv.ParseUint(roleIDStr, 10, 32); err == nil {
			roleIDUint := uint(id)
			roleID = &roleIDUint
		}
	}

	result, err := h.roleService.GetPermissionTree(ctx.Request.Context(), roleID)
	if err != nil {
		response.ErrorWithMessage(ctx, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(ctx, result)
}

// UpdateRolePermissions 更新角色权限
func (h *AdminRoleHandler) UpdateRolePermissions(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.CodeBadRequest, "无效的角色ID")
		return
	}

	var req types.RolePermissionUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(ctx, response.CodeBadRequest, err.Error())
		return
	}

	if err := h.roleService.UpdateRolePermissions(ctx.Request.Context(), uint(id), &req); err != nil {
		response.ErrorWithMessage(ctx, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// AssignRolesToUser 为用户分配角色
func (h *AdminRoleHandler) AssignRolesToUser(ctx *gin.Context) {
	var req types.RoleAssignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMessage(ctx, response.CodeBadRequest, err.Error())
		return
	}

	if err := h.roleService.AssignRolesToUser(ctx.Request.Context(), &req); err != nil {
		response.ErrorWithMessage(ctx, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// GetRoleOptions 获取角色选项（用于下拉选择）
func (h *AdminRoleHandler) GetRoleOptions(c *gin.Context) {
	roles, err := h.roleService.GetActiveRoles(c.Request.Context())
	if err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	// 转换为简化的选项格式
	var options []types.RoleOptionResponse
	for _, role := range roles {
		options = append(options, types.RoleOptionResponse{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
		})
	}

	response.Success(c, options)
}
