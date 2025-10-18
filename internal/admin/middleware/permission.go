package middleware

import (
	"net/http"

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
			c.JSON(http.StatusOK, gin.H{
				"code": 401,
				"msg":  "未授权",
			})
			c.Abort()
			return
		}
		
		permissions, err := pm.userService.GetUserPermissions(userID.(uint))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "获取权限失败",
			})
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
			c.JSON(http.StatusOK, gin.H{
				"code": 403,
				"msg":  "无权访问",
			})
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
			c.JSON(http.StatusOK, gin.H{
				"code": 401,
				"msg":  "未授权",
			})
			c.Abort()
			return
		}
		
		permissions, err := pm.userService.GetUserPermissions(userID.(uint))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "获取权限失败",
			})
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
			c.JSON(http.StatusOK, gin.H{
				"code": 403,
				"msg":  "无权访问",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}
