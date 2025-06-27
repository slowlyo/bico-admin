package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"bico-admin/core/config"
	"bico-admin/core/middleware"
	coreRouter "bico-admin/core/router"
	"bico-admin/modules/api/handler"
)

// SetupRoutes 设置对外API路由
func SetupRoutes(app fiber.Router, db *gorm.DB, cfg *config.Config) {
	// 设置核心认证路由（不需要认证的路由）
	coreRouter.SetupAuthRoutes(app, db, cfg)

	// 从配置中获取JWT密钥
	jwtSecret := cfg.JWT.Secret

	// 需要认证的路由
	protected := app.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSecret))

	// API处理器
	appHandler := handler.NewAppHandler(db)

	// 应用相关接口
	appGroup := protected.Group("/app")
	{
		appGroup.Get("/info", appHandler.GetAppInfo)
		appGroup.Get("/config", appHandler.GetAppConfig)
	}

	// 用户相关接口
	userGroup := protected.Group("/user")
	{
		userGroup.Get("/profile", appHandler.GetUserProfile)
		userGroup.Put("/profile", appHandler.UpdateUserProfile)
	}

	// 内容相关接口
	contentGroup := protected.Group("/content")
	{
		contentGroup.Get("/", appHandler.GetContentList)
		contentGroup.Get("/:id", appHandler.GetContent)
	}

	// 公开接口（不需要认证）
	public := app.Group("/public")
	{
		public.Get("/health", appHandler.HealthCheck)
		public.Get("/version", appHandler.GetVersion)
	}
}
