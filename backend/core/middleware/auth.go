package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"bico-admin/core/cache"
	"bico-admin/pkg/response"
)

// JWTClaims JWT声明
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AuthMiddleware JWT认证中间件
func AuthMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 获取Authorization头
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "Missing authorization header")
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return response.Unauthorized(c, "Invalid authorization header format")
		}

		// 提取token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return response.Unauthorized(c, "Missing token")
		}

		// 检查token是否在黑名单中
		cacheInstance := cache.GetCache()
		isBlacklisted, err := cacheInstance.IsTokenBlacklisted(tokenString)
		if err == nil && isBlacklisted {
			return response.Unauthorized(c, "Token has been revoked")
		}

		// 解析token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			return response.Unauthorized(c, "Invalid token")
		}

		// 验证token
		if !token.Valid {
			return response.Unauthorized(c, "Invalid token")
		}

		// 获取claims
		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			return response.Unauthorized(c, "Invalid token claims")
		}

		// 检查token是否过期
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return response.Unauthorized(c, "Token expired")
		}

		// 将用户信息存储到上下文
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)

		return c.Next()
	}
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint, username, secret string, expireDuration time.Duration) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *fiber.Ctx) uint {
	if userID, ok := c.Locals("user_id").(uint); ok {
		return userID
	}
	return 0
}

// GetUsername 从上下文获取用户名
func GetUsername(c *fiber.Ctx) string {
	if username, ok := c.Locals("username").(string); ok {
		return username
	}
	return ""
}
