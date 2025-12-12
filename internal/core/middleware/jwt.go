package middleware

import (
	"bico-admin/internal/pkg/jwt"
	"bico-admin/internal/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth(jwtManager *jwt.JWTManager, authService interface{ IsTokenBlacklisted(token string) bool }) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorWithCode(c, 401, "请先登录")
			c.Abort()
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorWithCode(c, 401, "token 格式错误")
			c.Abort()
			return
		}

		token := parts[1]

		// 检查 token 是否在黑名单中
		if authService.IsTokenBlacklisted(token) {
			response.ErrorWithCode(c, 401, "token 已失效")
			c.Abort()
			return
		}

		// 验证 token
		claims, err := jwtManager.ParseToken(token)

		if err != nil {
			response.ErrorWithCode(c, 401, "token 无效或已过期")
			c.Abort()
			return
		}
		if claims == nil || claims.UserID == 0 || claims.Username == "" {
			response.ErrorWithCode(c, 401, "token 无效")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
