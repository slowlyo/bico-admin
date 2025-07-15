package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"bico-admin/internal/admin/service"
	"bico-admin/pkg/cache"
	"bico-admin/pkg/config"
	"bico-admin/pkg/jwt"
	"bico-admin/pkg/response"
)

// Auth 认证中间件（向后兼容，不支持黑名单检查）
func Auth() gin.HandlerFunc {
	return AuthWithCache(nil)
}

// AuthWithCache 带缓存支持的认证中间件
func AuthWithCache(cache cache.Cache) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证令牌")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			response.Unauthorized(c, "认证令牌格式错误")
			c.Abort()
			return
		}

		// 提取令牌
		token := authHeader[len(bearerPrefix):]
		if token == "" {
			response.Unauthorized(c, "认证令牌为空")
			c.Abort()
			return
		}

		// 验证JWT令牌
		userID, userType, err := validateTokenWithCache(c.Request.Context(), token, cache)
		if err != nil {
			response.Unauthorized(c, "认证令牌无效")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", userID)
		c.Set("user_type", userType)
		c.Set("token", token)

		c.Next()
	})
}

// validateTokenWithCache 验证令牌并检查黑名单
func validateTokenWithCache(ctx context.Context, token string, cache cache.Cache) (uint, string, error) {
	cfg := config.Get()
	var jwtManager *jwt.JWTManager

	if cache != nil {
		jwtManager = jwt.NewJWTManagerWithCache(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.ExpireTime, cache)

		// 检查令牌是否在黑名单中
		isBlacklisted, err := jwtManager.IsBlacklisted(ctx, token)
		if err != nil {
			return 0, "", err
		}
		if isBlacklisted {
			return 0, "", errors.New("令牌已失效")
		}
	} else {
		jwtManager = jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.ExpireTime)
	}

	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		return 0, "", err
	}

	return claims.UserID, claims.UserType, nil
}

// RequirePermission 权限检查中间件
func RequirePermission(permission string) gin.HandlerFunc {
	return RequirePermissionWithService(permission, nil)
}

// RequirePermissionWithService 带服务的权限检查中间件
func RequirePermissionWithService(permission string, adminUserService service.AdminUserService) gin.HandlerFunc {
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

		// 如果有AdminUserService，使用真正的权限检查
		if adminUserService != nil && userType == "admin" {
			adminUser, err := adminUserService.GetByID(c.Request.Context(), userID.(uint))
			if err != nil {
				response.Forbidden(c, "权限验证失败")
				c.Abort()
				return
			}

			// 检查用户是否有指定权限（超级管理员会自动通过）
			if !adminUser.HasPermission(permission) {
				response.Forbidden(c, "权限不足")
				c.Abort()
				return
			}
		} else {
			// 回退到简单的权限检查
			if !hasPermission(userType.(string), permission) {
				response.Forbidden(c, "权限不足")
				c.Abort()
				return
			}
		}

		c.Next()
	})
}

// hasPermission 检查用户是否有指定权限（临时实现）
func hasPermission(userType, permission string) bool {
	// TODO: 实现真正的权限检查逻辑
	switch userType {
	case "admin":
		return true // 管理员拥有所有权限
	case "master":
		// 主控用户的权限列表
		masterPermissions := []string{
			"user:read", "user:write",
			"system:read",
		}
		for _, p := range masterPermissions {
			if p == permission {
				return true
			}
		}
	}
	return false
}
