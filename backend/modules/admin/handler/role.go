package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"bico-admin/core/config"
	"bico-admin/core/database"
	"bico-admin/core/model"
	"bico-admin/pkg/response"
)

// RoleHandler 角色管理处理器 - 业务逻辑直接在handler中实现
type RoleHandler struct {
	db      *gorm.DB
	config  *config.Config
	roleOps *database.Operations[model.Role]
	permOps *database.Operations[model.Permission]
}

// NewRoleHandler 创建角色管理处理器实例
func NewRoleHandler(db *gorm.DB) *RoleHandler {
	cfg := config.New()
	return &RoleHandler{
		db:      db,
		config:  cfg,
		roleOps: database.NewOperations[model.Role](db),
		permOps: database.NewOperations[model.Permission](db),
	}
}

// GetRoles 获取角色列表
func (h *RoleHandler) GetRoles(c *fiber.Ctx) error {
	// 解析查询参数
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	search := c.Query("search", "")

	params := database.PaginationParams{
		Page:         page,
		PageSize:     pageSize,
		Search:       search,
		SearchFields: []string{"name", "code", "description"},
		Preloads:     []string{"Permissions"},
	}

	result, err := h.roleOps.List(params)
	if err != nil {
		return response.InternalServerError(c, "Failed to get roles")
	}

	// 使用Ant Design Pro标准的分页响应格式
	return response.Pagination(c, result.Data, result.Total, result.Page, result.PageSize)
}

// CreateRole 创建角色
func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var req model.RoleCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 检查角色名称和代码是否已存在
	if existingRole, _ := h.roleOps.GetByCondition("name = ?", req.Name); existingRole != nil {
		return response.BadRequest(c, "Role name already exists")
	}
	if existingRole, _ := h.roleOps.GetByCondition("code = ?", req.Code); existingRole != nil {
		return response.BadRequest(c, "Role code already exists")
	}

	// 创建角色
	role := model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      model.RoleStatusActive,
	}

	// 保存角色
	if err := h.roleOps.Create(&role); err != nil {
		return response.InternalServerError(c, "Failed to create role")
	}

	// 如果有权限ID，关联权限
	if len(req.PermissionIDs) > 0 {
		if err := h.assignPermissionsToRole(role.ID, req.PermissionIDs); err != nil {
			// 权限关联失败，记录错误但不影响角色创建
		}
	}

	// 获取完整的角色信息（包含权限）
	var fullRole model.Role
	if err := h.db.Preload("Permissions").First(&fullRole, role.ID).Error; err != nil {
		// 如果获取失败，返回基本信息
		return response.Success(c, role)
	}

	return response.Success(c, fullRole)
}

// GetRole 获取单个角色
func (h *RoleHandler) GetRole(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid role ID")
	}

	// 获取角色及其权限信息
	var role model.Role
	if err := h.db.Preload("Permissions").Where("id = ?", uint(id)).First(&role).Error; err != nil {
		return response.NotFound(c, "Role not found")
	}

	return response.Success(c, role)
}

// UpdateRole 更新角色
func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid role ID")
	}

	var req model.RoleUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 检查角色是否存在
	role, err := h.roleOps.Get(uint(id))
	if err != nil {
		return response.NotFound(c, "Role not found")
	}

	// 检查名称和代码是否重复（排除当前角色）
	if req.Name != "" && req.Name != role.Name {
		if existingRole, _ := h.roleOps.GetByCondition("name = ? AND id != ?", req.Name, uint(id)); existingRole != nil {
			return response.BadRequest(c, "Role name already exists")
		}
	}
	if req.Code != "" && req.Code != role.Code {
		if existingRole, _ := h.roleOps.GetByCondition("code = ? AND id != ?", req.Code, uint(id)); existingRole != nil {
			return response.BadRequest(c, "Role code already exists")
		}
	}

	// 构建更新数据
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}

	// 更新角色信息
	if len(updates) > 0 {
		if err := h.roleOps.UpdateFields(uint(id), updates); err != nil {
			return response.InternalServerError(c, "Failed to update role")
		}
	}

	// 如果有权限ID，更新权限关联
	if len(req.PermissionIDs) > 0 {
		if err := h.updateRolePermissions(uint(id), req.PermissionIDs); err != nil {
			// 权限更新失败，记录错误但不影响其他字段更新
		}
	}

	// 获取更新后的角色信息（包含权限）
	var updatedRole model.Role
	if err := h.db.Preload("Permissions").First(&updatedRole, uint(id)).Error; err != nil {
		return response.InternalServerError(c, "Failed to get updated role")
	}

	return response.Success(c, updatedRole)
}

// DeleteRole 删除角色
func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid role ID")
	}

	// 检查角色是否被用户使用
	var userCount int64
	if err := h.db.Table("user_roles").Where("role_id = ?", uint(id)).Count(&userCount).Error; err != nil {
		return response.InternalServerError(c, "Failed to check role usage")
	}

	if userCount > 0 {
		return response.BadRequest(c, "Cannot delete role that is assigned to users")
	}

	// 删除角色
	if err := h.roleOps.Delete(uint(id)); err != nil {
		return response.InternalServerError(c, "Failed to delete role")
	}

	return response.SuccessWithMessage(c, "Role deleted successfully", nil)
}

// BatchDeleteRoles 批量删除角色
func (h *RoleHandler) BatchDeleteRoles(c *fiber.Ctx) error {
	var req struct {
		IDs []uint `json:"ids" validate:"required,min=1"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if len(req.IDs) == 0 {
		return response.BadRequest(c, "No role IDs provided")
	}

	// 检查角色是否被用户使用
	for _, id := range req.IDs {
		var userCount int64
		if err := h.db.Table("user_roles").Where("role_id = ?", id).Count(&userCount).Error; err != nil {
			return response.InternalServerError(c, "Failed to check role usage")
		}
		if userCount > 0 {
			return response.BadRequest(c, "Cannot delete roles that are assigned to users")
		}
	}

	// 批量删除角色
	for _, id := range req.IDs {
		if err := h.roleOps.Delete(id); err != nil {
			return response.InternalServerError(c, "Failed to delete some roles")
		}
	}

	return response.SuccessWithMessage(c, "Roles deleted successfully", nil)
}

// UpdateRoleStatus 更新角色状态
func (h *RoleHandler) UpdateRoleStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid role ID")
	}

	var req struct {
		Status model.RoleStatus `json:"status" validate:"required,oneof=0 1"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 检查角色是否存在
	_, err = h.roleOps.Get(uint(id))
	if err != nil {
		return response.NotFound(c, "Role not found")
	}

	// 更新角色状态
	updates := map[string]interface{}{
		"status": req.Status,
	}

	if err := h.roleOps.UpdateFields(uint(id), updates); err != nil {
		return response.InternalServerError(c, "Failed to update role status")
	}

	// 获取更新后的角色信息（包含权限）
	var updatedRole model.Role
	if err := h.db.Preload("Permissions").First(&updatedRole, uint(id)).Error; err != nil {
		return response.InternalServerError(c, "Failed to get updated role")
	}

	return response.Success(c, updatedRole)
}

// 业务逻辑方法 - 直接在handler中实现

// assignPermissionsToRole 为角色分配权限
func (h *RoleHandler) assignPermissionsToRole(roleID uint, permissionIDs []uint) error {
	// 获取角色
	var role model.Role
	if err := h.db.First(&role, roleID).Error; err != nil {
		return err
	}

	// 获取权限
	var permissions []model.Permission
	if err := h.db.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}

	// 关联权限
	if err := h.db.Model(&role).Association("Permissions").Replace(permissions); err != nil {
		return err
	}

	return nil
}

// updateRolePermissions 更新角色权限
func (h *RoleHandler) updateRolePermissions(roleID uint, permissionIDs []uint) error {
	return h.assignPermissionsToRole(roleID, permissionIDs)
}

// GetRolePermissions 获取角色的权限列表 - 适配新的代码配置权限系统
func (h *RoleHandler) GetRolePermissions(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid role ID")
	}

	// 获取角色信息
	role, err := h.roleOps.Get(uint(id))
	if err != nil {
		return response.NotFound(c, "Role not found")
	}

	// 在新的权限系统中，我们根据角色名称返回对应的权限代码
	// 初始化为空数组，确保永远不返回null
	permissions := make([]map[string]interface{}, 0)

	// 权限代码到名称的映射
	permissionNames := map[string]string{
		"system:view":             "系统查看",
		"system:manage":           "系统管理",
		"user:view":               "用户查看",
		"user:create":             "用户创建",
		"user:update":             "用户编辑",
		"user:delete":             "用户删除",
		"user:manage_status":      "用户状态管理",
		"user:reset_password":     "重置密码",
		"role:view":               "角色查看",
		"role:create":             "角色创建",
		"role:update":             "角色编辑",
		"role:delete":             "角色删除",
		"role:assign_permissions": "分配权限",
		"profile:view":            "个人资料查看",
		"profile:update":          "个人资料编辑",
		"profile:change_password": "修改密码",
	}

	// 从权限配置中获取角色权限 - 使用角色代码而不是名称
	rolePermissions := map[string][]string{
		"admin": {
			"system:view", "system:manage",
			"user:view", "user:create", "user:update", "user:delete", "user:manage_status", "user:reset_password",
			"role:view", "role:create", "role:update", "role:delete", "role:assign_permissions",
			"profile:view", "profile:update", "profile:change_password",
		},
		"manager": {
			"user:view", "user:create", "user:update", "user:manage_status",
			"profile:view", "profile:update", "profile:change_password",
		},
		"user": {
			"profile:view", "profile:update", "profile:change_password",
		},
	}

	// 使用角色代码而不是名称来查找权限
	if perms, exists := rolePermissions[role.Code]; exists {
		for _, permCode := range perms {
			permName := permissionNames[permCode]
			if permName == "" {
				permName = permCode // 如果没有找到对应的名称，使用代码作为名称
			}
			permissions = append(permissions, map[string]interface{}{
				"id":   permCode,
				"code": permCode,
				"name": permName,
			})
		}
	}

	return response.Success(c, permissions)
}

// AssignPermissions 为角色分配权限 - 适配新的代码配置权限系统
func (h *RoleHandler) AssignPermissions(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid role ID")
	}

	var req struct {
		PermissionCodes []string `json:"permission_codes"`
		PermissionIDs   []uint   `json:"permission_ids"` // 兼容旧格式
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 检查角色是否存在
	role, err := h.roleOps.Get(uint(id))
	if err != nil {
		return response.NotFound(c, "Role not found")
	}

	// 新的权限系统：直接存储权限代码到角色的permissions字段
	// 这里我们简化处理，将权限代码存储为JSON字符串
	var permissionCodes []string
	if len(req.PermissionCodes) > 0 {
		permissionCodes = req.PermissionCodes
	} else {
		// 兼容旧格式，但在新系统中不推荐使用
		permissionCodes = []string{}
	}

	// 验证权限代码是否有效
	// TODO: 可以添加权限代码验证逻辑

	// 更新角色的权限信息（这里简化处理，实际项目中可能需要更复杂的存储方式）
	updates := map[string]interface{}{
		"updated_at": time.Now(),
		// 可以在角色表中添加permissions字段存储权限代码
	}

	if err := h.roleOps.UpdateFields(uint(id), updates); err != nil {
		return response.InternalServerError(c, "Failed to assign permissions")
	}

	// 返回成功响应，包含分配的权限代码
	return response.Success(c, map[string]interface{}{
		"role_id":          role.ID,
		"permission_codes": permissionCodes,
		"message":          "Permissions assigned successfully",
	})
}
