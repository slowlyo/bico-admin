package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"bico-admin/core/config"
	"bico-admin/core/handler"
	"bico-admin/core/middleware"
)

// SetupAuthRoutes 设置认证路由
func SetupAuthRoutes(app fiber.Router, db *gorm.DB, cfg *config.Config) {
	authHandler := handler.NewAuthHandler(db, cfg)

	auth := app.Group("/auth")
	{
		// 不需要认证的路由
		auth.Post("/login", authHandler.Login)

		// 需要认证的路由
		protected := auth.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			protected.Post("/logout", authHandler.Logout)
			protected.Get("/profile", authHandler.GetProfile)
			protected.Put("/profile", authHandler.UpdateProfile)
			protected.Post("/change-password", authHandler.ChangePassword)
			protected.Get("/permissions", authHandler.GetUserPermissions)
		}
	}
}
