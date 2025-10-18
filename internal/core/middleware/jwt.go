package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth(jwtManager interface{}, authService interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 401,
				"msg":  "请先登录",
			})
			c.Abort()
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusOK, gin.H{
				"code": 401,
				"msg":  "token 格式错误",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// 检查 token 是否在黑名单中
		if authService.(interface {
			IsTokenBlacklisted(token string) bool
		}).IsTokenBlacklisted(token) {
			c.JSON(http.StatusOK, gin.H{
				"code": 401,
				"msg":  "token 已失效",
			})
			c.Abort()
			return
		}

		// 验证 token
		claims, err := jwtManager.(interface {
			ValidateToken(token string) (map[string]interface{}, error)
		}).ValidateToken(token)

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 401,
				"msg":  "token 无效或已过期",
			})
			c.Abort()
			return
		}

		userID := uint(claims["user_id"].(float64))

		// 将用户信息存入上下文
		c.Set("user_id", userID)
		c.Set("username", claims["username"].(string))

		c.Next()
	}
}
