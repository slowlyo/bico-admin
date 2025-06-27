package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"bico-admin/core/config"
)

// SetupRoutes 设置核心路由
func SetupRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	// 设置认证相关路由
	SetupAuthRoutes(app, db, cfg)

	// 设置系统路由
	SetupSystemRoutes(app, db)
}

// SetupSystemRoutes 设置系统路由
func SetupSystemRoutes(app *fiber.App, db *gorm.DB) {
	// 系统信息路由
	system := app.Group("/system")
	{
		system.Get("/info", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"name":    "Bico Admin",
				"version": "1.0.0",
				"status":  "running",
			})
		})

		system.Get("/health", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"status":   "ok",
				"database": "connected",
			})
		})
	}
}
