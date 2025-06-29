package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"bico-admin/core/database"
	"bico-admin/core/model"
	"bico-admin/core/permission"
	"bico-admin/pkg/response"
	"bico-admin/pkg/validator"
)

// RolePermissionHandler 角色权限管理处理器 - 遵循简化架构，业务逻辑直接在handler中实现
type RolePermissionHandler struct {
	db      *gorm.DB
	roleOps *database.Operations[model.Role]
}

// NewRolePermissionHandler 创建角色权限管理处理器
func NewRolePermissionHandler(db *gorm.DB) *RolePermissionHandler {
	return &RolePermissionHandler{
		db:      db,
		roleOps: database.NewOperations[model.Role](db),
	}
}

// GetAllPermissions 获取所有权限定义（从代码中获取）
func (h *RolePermissionHandler) GetAllPermissions(c *fiber.Ctx) error {
	// 从代码中获取权限定义
	permissions := permission.AllPermissions

	// 按分类组织权限
	categories := permission.GetPermissionsByCategory()

	result := fiber.Map{
		"permissions": permissions,
		"categories":  categories,
	}

	return response.Success(c, result)
}

// GetRolePermissions 获取角色的权限列表
func (h *RolePermissionHandler) GetRolePermissions(c *fiber.Ctx) error {
	roleIDStr := c.Params("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid role ID")
	}

	// 检查角色是否存在
	role, err := h.roleOps.GetById(uint(roleID))
	if err != nil {
		return response.NotFound(c, "Role not found")
	}

	// 检查是否为受保护角色
	if h.isProtectedRole(role.Code) {
		// 超级管理员拥有所有权限
		allPermissions := permission.GetAllPermissionCodes()
		return response.Success(c, allPermissions)
	}

	// 从 role_permissions 表查询角色权限代码
	var rolePermissions []model.RolePermission
	err = h.db.Where("role_id = ?", roleID).Find(&rolePermissions).Error
	if err != nil {
		return response.InternalServerError(c, "Failed to get role permissions")
	}

	// 提取权限代码
	var permissionCodes []string
	for _, rp := range rolePermissions {
		permissionCodes = append(permissionCodes, rp.PermissionCode)
	}

	return response.Success(c, permissionCodes)
}

// AssignRolePermissions 为角色分配权限
func (h *RolePermissionHandler) AssignRolePermissions(c *fiber.Ctx) error {
	roleIDStr := c.Params("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid role ID")
	}

	var req struct {
		PermissionCodes []string `json:"permission_codes" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 验证请求参数
	if errors := validator.Validate(req); len(errors) > 0 {
		return response.ValidationError(c, errors)
	}

	// 检查角色是否存在
	role, err := h.roleOps.GetById(uint(roleID))
	if err != nil {
		return response.NotFound(c, "Role not found")
	}

	// 检查是否为受保护角色
	if h.isProtectedRole(role.Code) {
		return response.BadRequest(c, "Cannot modify permissions for protected role")
	}

	// 验证权限代码是否有效
	allPermissionCodes := permission.GetAllPermissionCodes()
	for _, code := range req.PermissionCodes {
		if !h.isValidPermissionCode(code, allPermissionCodes) {
			return response.BadRequest(c, "Invalid permission code: "+code)
		}
	}

	// 开始事务
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 清除角色现有权限
	if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
		tx.Rollback()
		return response.InternalServerError(c, "Failed to clear existing permissions")
	}

	// 为每个权限代码创建角色权限关联
	for _, code := range req.PermissionCodes {
		// 直接创建角色权限关联，存储权限代码
		rolePermission := model.RolePermission{
			RoleID:         uint(roleID),
			PermissionCode: code,
		}
		if err := tx.Create(&rolePermission).Error; err != nil {
			tx.Rollback()
			return response.InternalServerError(c, "Failed to assign permission: "+code)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return response.InternalServerError(c, "Failed to commit transaction")
	}

	return response.SuccessWithMessage(c, "Permissions assigned successfully", nil)
}

// RemoveRolePermission 移除角色的特定权限
func (h *RolePermissionHandler) RemoveRolePermission(c *fiber.Ctx) error {
	roleIDStr := c.Params("roleId")
	permissionCode := c.Params("permissionCode")

	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid role ID")
	}

	// 检查角色是否存在
	role, err := h.roleOps.GetById(uint(roleID))
	if err != nil {
		return response.NotFound(c, "Role not found")
	}

	// 检查是否为受保护角色
	if h.isProtectedRole(role.Code) {
		return response.BadRequest(c, "Cannot modify permissions for protected role")
	}

	// 验证权限代码是否有效
	if permission.GetPermissionByCode(permissionCode) == nil {
		return response.BadRequest(c, "Invalid permission code")
	}

	// 删除角色权限关联
	result := h.db.Where("role_id = ? AND permission_code = ?", roleID, permissionCode).
		Delete(&model.RolePermission{})

	if result.Error != nil {
		return response.InternalServerError(c, "Failed to remove permission")
	}

	if result.RowsAffected == 0 {
		return response.NotFound(c, "Permission not assigned to this role")
	}

	return response.SuccessWithMessage(c, "Permission removed successfully", nil)
}

// 注意：不再需要 getOrCreatePermission 方法
// 权限代码直接存储在 role_permissions 表中，不依赖权限表

// isProtectedRole 检查是否为受保护角色
func (h *RolePermissionHandler) isProtectedRole(roleCode string) bool {
	return roleCode == "super_admin"
}

// isValidPermissionCode 检查权限代码是否有效
func (h *RolePermissionHandler) isValidPermissionCode(code string, validCodes []string) bool {
	for _, validCode := range validCodes {
		if code == validCode {
			return true
		}
	}
	return false
}
