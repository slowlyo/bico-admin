package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"bico-admin/pkg/response"
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 执行下一个处理器
		err := c.Next()
		
		if err != nil {
			// 记录错误日志
			log.Printf("Error: %v, Path: %s, Method: %s, IP: %s", 
				err, c.Path(), c.Method(), c.IP())

			// 检查是否是Fiber错误
			if e, ok := err.(*fiber.Error); ok {
				return response.Error(c, e.Code, e.Message)
			}

			// 默认内部服务器错误
			return response.InternalServerError(c, "Internal server error")
		}

		return nil
	}
}

// RecoverMiddleware 恢复中间件，防止panic导致程序崩溃
func RecoverMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v, Path: %s, Method: %s, IP: %s", 
					r, c.Path(), c.Method(), c.IP())
				
				// 返回内部服务器错误
				response.InternalServerError(c, "Internal server error")
			}
		}()

		return c.Next()
	}
}
