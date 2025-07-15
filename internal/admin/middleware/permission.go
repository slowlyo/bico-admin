package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"bico-admin/internal/admin/definitions"
	"bico-admin/internal/admin/service"
	"bico-admin/pkg/response"
)

// 权限白名单 - 这些路径不需要权限检查
var permissionWhitelist = []string{
	"/admin/auth/logout",  // 登出
	"/admin/auth/profile", // 获取个人信息
	"/admin/auth/refresh", // 刷新token
}

// PermissionMiddleware 权限检查中间件
func PermissionMiddleware(adminUserService service.AdminUserService) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "用户未登录")
			c.Abort()
			return
		}

		userType, exists := c.Get("user_type")
		if !exists {
			response.Unauthorized(c, "用户未登录")
			c.Abort()
			return
		}

		// 只对admin用户进行权限检查
		if userType != "admin" {
			c.Next()
			return
		}

		// 获取当前请求的API路径
		apiPath := c.Request.URL.Path

		// 查找匹配的权限
		requiredPermission := findPermissionByAPI(apiPath)
		if requiredPermission == "" {
			// 如果没有找到匹配的权限配置，检查是否在白名单中
			if !isWhitelistedPath(apiPath) {
				response.Forbidden(c, "该接口需要权限配置")
				c.Abort()
				return
			}
			// 对于白名单中的路径，允许访问
			c.Next()
			return
		}

		// 获取用户信息并检查权限
		adminUser, err := adminUserService.GetByID(c.Request.Context(), userID.(uint))
		if err != nil {
			response.Forbidden(c, "权限验证失败")
			c.Abort()
			return
		}

		// 检查用户是否有指定权限（超级管理员会自动通过）
		if !adminUser.HasPermission(requiredPermission) {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	})
}

// findPermissionByAPI 根据API路径查找对应的权限代码
func findPermissionByAPI(apiPath string) string {
	allPermissions := definitions.GetAllPermissionsFlat()

	for _, permission := range allPermissions {
		// 只检查操作类型的权限
		if permission.Type != definitions.PermissionTypeAction {
			continue
		}

		// 检查API路径是否匹配
		for _, permissionAPI := range permission.APIs {
			// 处理多个API路径（用逗号分隔）
			apiPaths := strings.Split(permissionAPI, ",")
			for _, singleAPI := range apiPaths {
				singleAPI = strings.TrimSpace(singleAPI)
				if matchAPIPath(singleAPI, apiPath) {
					return permission.Code
				}
			}
		}
	}

	return ""
}

// matchAPIPath 匹配API路径（支持参数路径如 /admin/admin-users/:id）
func matchAPIPath(pattern, path string) bool {
	// 简单的完全匹配
	if pattern == path {
		return true
	}

	// 处理参数路径匹配，如 /admin/admin-users/:id 匹配 /admin/admin-users/123
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	// 路径段数量必须相同
	if len(patternParts) != len(pathParts) {
		return false
	}

	// 逐段比较
	for i, patternPart := range patternParts {
		pathPart := pathParts[i]

		// 如果是参数段（以:开头），则跳过比较
		if strings.HasPrefix(patternPart, ":") {
			continue
		}

		// 普通段必须完全匹配
		if patternPart != pathPart {
			return false
		}
	}

	return true
}

// isWhitelistedPath 检查路径是否在权限白名单中
func isWhitelistedPath(path string) bool {
	for _, whitelistPath := range permissionWhitelist {
		if matchAPIPath(whitelistPath, path) {
			return true
		}
	}
	return false
}

// PermissionMiddlewareFactory 权限中间件工厂
func PermissionMiddlewareFactory(adminUserService service.AdminUserService) gin.HandlerFunc {
	return PermissionMiddleware(adminUserService)
}
