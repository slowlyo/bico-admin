package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,              // 最大请求数
		Expiration: 1 * time.Minute,  // 时间窗口
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // 基于IP限流
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
				"message": "Rate limit exceeded",
			})
		},
	})
}

// AuthRateLimitMiddleware 认证接口限流中间件（更严格）
func AuthRateLimitMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,               // 最大请求数
		Expiration: 1 * time.Minute,  // 时间窗口
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // 基于IP限流
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
				"message": "Authentication rate limit exceeded",
			})
		},
	})
}
