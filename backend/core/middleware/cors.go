package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORSMiddleware CORS中间件配置
func CORSMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type",
		MaxAge:           86400, // 24小时
	})
}
