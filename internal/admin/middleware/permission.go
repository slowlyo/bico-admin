package middleware

import (
	"bico-admin/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// PermissionMiddleware 权限检查中间件
type PermissionMiddleware struct {
	userService interface {
		GetUserPermissions(userID uint) ([]string, error)
	}
}

// NewPermissionMiddleware 创建权限中间件
func NewPermissionMiddleware(userService interface {
	GetUserPermissions(userID uint) ([]string, error)
}) *PermissionMiddleware {
	return &PermissionMiddleware{
		userService: userService,
	}
}

// RequirePermission 要求指定权限
func (pm *PermissionMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithCode(c, 401, "未授权")
			c.Abort()
			return
		}

		permissions, err := pm.userService.GetUserPermissions(userID.(uint))
		if err != nil {
			response.ErrorWithCode(c, 500, "获取权限失败")
			c.Abort()
			return
		}

		hasPermission := false
		for _, perm := range permissions {
			if perm == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			response.ErrorWithCode(c, 403, "无权访问")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission 要求任意一个权限
func (pm *PermissionMiddleware) RequireAnyPermission(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithCode(c, 401, "未授权")
			c.Abort()
			return
		}

		permissions, err := pm.userService.GetUserPermissions(userID.(uint))
		if err != nil {
			response.ErrorWithCode(c, 500, "获取权限失败")
			c.Abort()
			return
		}

		hasPermission := false
		for _, userPerm := range permissions {
			for _, reqPerm := range requiredPermissions {
				if userPerm == reqPerm {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			response.ErrorWithCode(c, 403, "无权访问")
			c.Abort()
			return
		}

		c.Next()
	}
}
