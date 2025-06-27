package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"bico-admin/core/config"
	"bico-admin/core/middleware"
	coreRouter "bico-admin/core/router"
	"bico-admin/modules/admin/handler"
)

// SetupRoutes 设置后台管理路由
func SetupRoutes(app fiber.Router, db *gorm.DB) {
	// 获取配置
	cfg := config.New()

	// API路由组
	api := app.Group("/api")

	// 设置认证路由（不需要认证的路由）
	coreRouter.SetupAuthRoutes(api, db, cfg)

	// 需要认证的路由
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))

	// 后台管理处理器
	dashboardHandler := handler.NewDashboardHandler(db)

	// 后台管理路由（需要认证）
	{
		// 仪表板
		protected.Get("/dashboard", dashboardHandler.GetDashboard)
		protected.Get("/dashboard/stats", dashboardHandler.GetStats)

		// 用户管理
		users := protected.Group("/users")
		{
			users.Get("/", dashboardHandler.GetUsers)
			users.Post("/", dashboardHandler.CreateUser)
			users.Get("/:id", dashboardHandler.GetUser)
			users.Put("/:id", dashboardHandler.UpdateUser)
			users.Delete("/:id", dashboardHandler.DeleteUser)
		}

		// 角色管理
		roles := protected.Group("/roles")
		{
			roles.Get("/", dashboardHandler.GetRoles)
			roles.Post("/", dashboardHandler.CreateRole)
			roles.Get("/:id", dashboardHandler.GetRole)
			roles.Put("/:id", dashboardHandler.UpdateRole)
			roles.Delete("/:id", dashboardHandler.DeleteRole)
		}

		// 权限管理
		permissions := protected.Group("/permissions")
		{
			permissions.Get("/", dashboardHandler.GetPermissions)
			permissions.Post("/", dashboardHandler.CreatePermission)
			permissions.Get("/:id", dashboardHandler.GetPermission)
			permissions.Put("/:id", dashboardHandler.UpdatePermission)
			permissions.Delete("/:id", dashboardHandler.DeletePermission)
		}
	}
}
