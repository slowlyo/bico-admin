package permission

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"bico-admin/core/middleware"
	"bico-admin/core/model"
	"bico-admin/pkg/response"
)

// PermissionMiddleware 权限中间件结构体
type PermissionMiddleware struct {
	db *gorm.DB
}

// NewPermissionMiddleware 创建权限中间件实例
func NewPermissionMiddleware(db *gorm.DB) *PermissionMiddleware {
	return &PermissionMiddleware{db: db}
}

// 注意：旧版本的静态权限验证方法已移除，请使用 PermissionMiddleware 实例方法

// RequirePermission 权限验证中间件
func (pm *PermissionMiddleware) RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 获取用户ID
		userID := middleware.GetUserID(c)
		if userID == 0 {
			return response.Unauthorized(c, "User not authenticated")
		}

		// 检查是否为个人资料相关权限（无需验证）
		if pm.isProfileRoute(permission) {
			return c.Next()
		}

		// 检查权限
		hasPermission, err := pm.hasPermission(userID, permission)
		if err != nil {
			return response.InternalServerError(c, "Permission check failed")
		}

		if !hasPermission {
			return response.Forbidden(c, "Insufficient permissions")
		}

		return c.Next()
	}
}

// RequireAnyPermission 需要任意一个权限的中间件
func (pm *PermissionMiddleware) RequireAnyPermission(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := middleware.GetUserID(c)
		if userID == 0 {
			return response.Unauthorized(c, "User not authenticated")
		}

		// 检查是否有任意一个权限
		for _, permission := range permissions {
			if pm.isProfileRoute(permission) {
				return c.Next()
			}
			hasPermission, err := pm.hasPermission(userID, permission)
			if err != nil {
				continue // 继续检查下一个权限
			}
			if hasPermission {
				return c.Next()
			}
		}

		return response.Forbidden(c, "Insufficient permissions")
	}
}

// RequireAllPermissions 需要所有权限的中间件
func (pm *PermissionMiddleware) RequireAllPermissions(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := middleware.GetUserID(c)
		if userID == 0 {
			return response.Unauthorized(c, "User not authenticated")
		}

		// 检查是否有所有权限
		for _, permission := range permissions {
			if pm.isProfileRoute(permission) {
				continue // 个人资料权限跳过检查
			}
			hasPermission, err := pm.hasPermission(userID, permission)
			if err != nil || !hasPermission {
				return response.Forbidden(c, "Insufficient permissions")
			}
		}

		return c.Next()
	}
}

// CheckPermissionInHandler 在处理器中检查权限的辅助函数
func (pm *PermissionMiddleware) CheckPermissionInHandler(c *fiber.Ctx, permission string) bool {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return false
	}

	if pm.isProfileRoute(permission) {
		return true
	}

	hasPermission, err := pm.hasPermission(userID, permission)
	if err != nil {
		return false
	}

	return hasPermission
}

// hasPermission 检查用户是否有指定权限
func (pm *PermissionMiddleware) hasPermission(userID uint, permission string) (bool, error) {
	// 检查用户是否为超级管理员
	isSuperAdmin, err := pm.isSuperAdmin(userID)
	if err != nil {
		return false, err
	}

	if isSuperAdmin {
		return true, nil
	}

	// 查询用户是否有该权限
	var count int64
	err = pm.db.Table("permissions p").
		Joins("JOIN role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN user_roles ur ON rp.role_id = ur.role_id").
		Where("ur.user_id = ? AND p.code = ? AND p.status = ?",
			userID, permission, model.PermissionStatusActive).
		Count(&count).Error

	return count > 0, err
}

// isSuperAdmin 检查用户是否为超级管理员
func (pm *PermissionMiddleware) isSuperAdmin(userID uint) (bool, error) {
	var count int64
	err := pm.db.Table("users u").
		Joins("JOIN user_roles ur ON u.id = ur.user_id").
		Joins("JOIN roles r ON ur.role_id = r.id").
		Where("u.id = ? AND r.code = ? AND r.status = ?",
			userID, "super_admin", model.RoleStatusActive).
		Count(&count).Error

	if err != nil {
		// 如果查询失败，尝试通过用户表的role字段检查（兼容旧版本）
		var user model.User
		if err := pm.db.First(&user, userID).Error; err != nil {
			return false, err
		}
		return user.Role == "super_admin", nil
	}

	return count > 0, nil
}

// isProfileRoute 检查是否为个人资料相关权限
func (pm *PermissionMiddleware) isProfileRoute(permission string) bool {
	profileRoutes := []string{
		"/auth/profile",
		"/auth/change-password",
		"profile:view",
		"profile:update",
		"profile:change_password",
	}

	for _, route := range profileRoutes {
		if permission == route {
			return true
		}
	}
	return false
}

// 兼容旧版本的静态方法（已废弃，建议使用权限中间件实例）
// RequirePermission 权限验证中间件（兼容旧版本）
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 这是兼容旧版本的方法，建议使用 PermissionMiddleware.RequirePermission
		return response.InternalServerError(c, "Please use PermissionMiddleware.RequirePermission instead")
	}
}
