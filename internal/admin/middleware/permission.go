package middleware

import (
	"bico-admin/internal/core/middleware"
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
		set, ok := pm.permissionSet(c)
		if !ok {
			return
		}
		if _, allowed := set[permission]; !allowed {
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
		set, ok := pm.permissionSet(c)
		if !ok {
			return
		}
		for _, permission := range requiredPermissions {
			if _, allowed := set[permission]; allowed {
				c.Next()
				return
			}
		}
		response.ErrorWithCode(c, 403, "无权访问")
		c.Abort()
	}
}

// permissionSet 取本次请求的权限集合。
//
// 说明：同一请求里多次校验只查一次缓存，结果放进上下文。返回 false 时响应已经写完。
func (pm *PermissionMiddleware) permissionSet(c *gin.Context) (map[string]struct{}, bool) {
	if value, ok := c.Get(middleware.CtxUserPermissions); ok {
		if set, ok := value.(map[string]struct{}); ok {
			return set, true
		}
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorWithCode(c, 401, "未授权")
		c.Abort()
		return nil, false
	}
	uid, ok := userID.(uint)
	if !ok {
		response.ErrorWithCode(c, 401, "未授权")
		c.Abort()
		return nil, false
	}

	permissions, err := pm.userService.GetUserPermissions(uid)
	if err != nil {
		response.ErrorWithCode(c, 500, "获取权限失败")
		c.Abort()
		return nil, false
	}

	set := make(map[string]struct{}, len(permissions))
	for _, permission := range permissions {
		set[permission] = struct{}{}
	}
	c.Set(middleware.CtxUserPermissions, set)
	return set, true
}
