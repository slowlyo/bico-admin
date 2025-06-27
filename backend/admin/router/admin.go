package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"bico-admin/admin/handler"
	"bico-admin/core/middleware"
)

// SetupRoutes 设置后台管理路由
func SetupRoutes(app fiber.Router, db *gorm.DB) {
	// TODO: 从配置中获取JWT密钥
	jwtSecret := "your-secret-key"
	
	// 应用认证中间件
	app.Use(middleware.AuthMiddleware(jwtSecret))

	// 后台管理处理器
	dashboardHandler := handler.NewDashboardHandler(db)
	
	// API路由组
	api := app.Group("/api")
	{
		// 仪表板
		api.Get("/dashboard", dashboardHandler.GetDashboard)
		api.Get("/dashboard/stats", dashboardHandler.GetStats)
		
		// 用户管理
		users := api.Group("/users")
		{
			users.Get("/", dashboardHandler.GetUsers)
			users.Post("/", dashboardHandler.CreateUser)
			users.Get("/:id", dashboardHandler.GetUser)
			users.Put("/:id", dashboardHandler.UpdateUser)
			users.Delete("/:id", dashboardHandler.DeleteUser)
		}
		
		// 角色管理
		roles := api.Group("/roles")
		{
			roles.Get("/", dashboardHandler.GetRoles)
			roles.Post("/", dashboardHandler.CreateRole)
			roles.Get("/:id", dashboardHandler.GetRole)
			roles.Put("/:id", dashboardHandler.UpdateRole)
			roles.Delete("/:id", dashboardHandler.DeleteRole)
		}
		
		// 权限管理
		permissions := api.Group("/permissions")
		{
			permissions.Get("/", dashboardHandler.GetPermissions)
			permissions.Post("/", dashboardHandler.CreatePermission)
			permissions.Get("/:id", dashboardHandler.GetPermission)
			permissions.Put("/:id", dashboardHandler.UpdatePermission)
			permissions.Delete("/:id", dashboardHandler.DeletePermission)
		}
	}
}
