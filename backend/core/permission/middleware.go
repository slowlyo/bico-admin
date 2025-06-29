package permission

import (
	"github.com/gofiber/fiber/v2"

	"bico-admin/core/middleware"
	"bico-admin/pkg/response"
)

// RequirePermission 权限验证中间件
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 获取用户角色（从认证中间件设置的上下文中获取）
		userRole := getUserRole(c)
		if userRole == "" {
			return response.Unauthorized(c, "User not authenticated")
		}

		// 检查权限
		if !HasPermission(userRole, permission) {
			return response.Forbidden(c, "Insufficient permissions")
		}

		return c.Next()
	}
}

// RequireAnyPermission 需要任意一个权限的中间件
func RequireAnyPermission(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := getUserRole(c)
		if userRole == "" {
			return response.Unauthorized(c, "User not authenticated")
		}

		if !HasAnyPermission(userRole, permissions) {
			return response.Forbidden(c, "Insufficient permissions")
		}

		return c.Next()
	}
}

// RequireAllPermissions 需要所有权限的中间件
func RequireAllPermissions(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := getUserRole(c)
		if userRole == "" {
			return response.Unauthorized(c, "User not authenticated")
		}

		if !HasAllPermissions(userRole, permissions) {
			return response.Forbidden(c, "Insufficient permissions")
		}

		return c.Next()
	}
}

// getUserRole 从上下文中获取用户角色
func getUserRole(c *fiber.Ctx) string {
	// 这里需要根据实际的用户模型来获取角色
	// 暂时返回一个默认值，实际使用时需要从数据库查询用户信息
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return ""
	}

	// TODO: 从数据库查询用户角色
	// 这里需要注入数据库连接或用户服务来查询用户角色
	// 暂时返回默认角色，实际使用时需要实现
	return "user"
}

// CheckPermissionInHandler 在处理器中检查权限的辅助函数
func CheckPermissionInHandler(c *fiber.Ctx, permission string) bool {
	userRole := getUserRole(c)
	return HasPermission(userRole, permission)
}
