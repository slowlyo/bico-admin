package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"bico-admin/core/handler"
)

// SetupAuthRoutes 设置认证路由
func SetupAuthRoutes(app fiber.Router, db *gorm.DB) {
	authHandler := handler.NewAuthHandler(db)

	auth := app.Group("/auth")
	{
		auth.Post("/login", authHandler.Login)
		auth.Post("/register", authHandler.Register)
		auth.Post("/logout", authHandler.Logout)
		auth.Get("/profile", authHandler.GetProfile)
		auth.Put("/profile", authHandler.UpdateProfile)
		auth.Post("/change-password", authHandler.ChangePassword)
	}
}
