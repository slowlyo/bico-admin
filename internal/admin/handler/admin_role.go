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
// @Summary 获取角色列表
// @Description 获取角色列表，支持分页和条件筛选
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param name query string false "角色名称"
// @Param code query string false "角色代码"
// @Param status query int false "状态"
// @Success 200 {object} response.Response{data=types.PageResponse[types.RoleResponse]}
// @Router /api/admin/roles [get]
func (h *AdminRoleHandler) GetRoleList(ctx *gin.Context) {
	var req types.RoleListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMessage(ctx, response.CodeBadRequest, err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	result, err := h.roleService.GetRoleList(ctx.Request.Context(), &req)
	if err != nil {
		response.ErrorWithMessage(ctx, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(ctx, result)
}

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建新角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param role body types.RoleCreateRequest true "角色信息"
// @Success 200 {object} response.Response{data=types.RoleResponse}
// @Router /api/admin/roles [post]
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
// @Summary 更新角色
// @Description 更新角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param role body types.RoleUpdateRequest true "角色信息"
// @Success 200 {object} response.Response{data=types.RoleResponse}
// @Router /api/admin/roles/{id} [put]
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

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response
// @Router /api/admin/roles/{id} [delete]
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
// @Summary 获取角色详情
// @Description 根据ID获取角色详情
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response{data=types.RoleResponse}
// @Router /api/admin/roles/{id} [get]
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
// @Summary 获取权限树
// @Description 获取权限树，用于角色权限分配
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param role_id query int false "角色ID，用于标记已选中的权限"
// @Success 200 {object} response.Response{data=[]types.PermissionTreeNode}
// @Router /api/admin/roles/permissions [get]
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

// AssignRolesToUser 为用户分配角色
// @Summary 为用户分配角色
// @Description 为管理员用户分配角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param assignment body types.RoleAssignRequest true "角色分配信息"
// @Success 200 {object} response.Response
// @Router /api/admin/roles/assign [post]
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

// GetUserRoles 获取用户角色
// @Summary 获取用户角色
// @Description 获取管理员用户的角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param user_id path int true "用户ID"
// @Success 200 {object} response.Response{data=types.UserRoleResponse}
// @Router /api/admin/users/{user_id}/roles [get]
func (h *AdminRoleHandler) GetUserRoles(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMessage(ctx, response.CodeBadRequest, "无效的用户ID")
		return
	}

	result, err := h.roleService.GetUserRoles(ctx.Request.Context(), uint(userID))
	if err != nil {
		response.ErrorWithMessage(ctx, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(ctx, result)
}

// GetRoleStats 获取角色统计
// @Summary 获取角色统计
// @Description 获取角色统计信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=types.RoleStatsResponse}
// @Router /api/admin/roles/stats [get]
func (h *AdminRoleHandler) GetRoleStats(ctx *gin.Context) {
	result, err := h.roleService.GetRoleStats(ctx.Request.Context())
	if err != nil {
		response.ErrorWithMessage(ctx, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(ctx, result)
}
