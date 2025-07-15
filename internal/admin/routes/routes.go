package routes

import (
	"github.com/gin-gonic/gin"

	"bico-admin/internal/admin/handler"
	"bico-admin/internal/admin/middleware"
	"bico-admin/pkg/cache"
)

// RegisterRoutes 注册admin端路由
func RegisterRoutes(r *gin.Engine, handlers *Handlers, cache cache.Cache, permissionMiddleware gin.HandlerFunc) {
	// admin端路由组
	adminGroup := r.Group("/admin")

	// 认证相关路由（无需认证）
	authGroup := adminGroup.Group("/auth")
	{
		authGroup.POST("/login", handlers.AuthHandler.Login)
		authGroup.POST("/refresh", handlers.AuthHandler.RefreshToken)
	}

	// 需要认证的路由
	protectedGroup := adminGroup.Group("")
	protectedGroup.Use(middleware.AuthWithCache(cache)) // 带缓存的认证中间件

	// 如果提供了权限中间件，则使用它进行权限检查
	if permissionMiddleware != nil {
		protectedGroup.Use(permissionMiddleware)
	}

	{
		// 认证相关
		protectedGroup.POST("/auth/logout", handlers.AuthHandler.Logout)
		protectedGroup.GET("/auth/profile", handlers.AuthHandler.GetProfile)
		protectedGroup.PUT("/auth/profile", handlers.AuthHandler.UpdateProfile)

		// 用户管理
		userGroup := protectedGroup.Group("/users")
		{
			userGroup.GET("", handlers.UserHandler.GetList)
			userGroup.GET("/stats", handlers.UserHandler.GetStats)
			userGroup.GET("/:id", handlers.UserHandler.GetByID)
			userGroup.POST("", handlers.UserHandler.Create)
			userGroup.PUT("/:id", handlers.UserHandler.Update)
			userGroup.DELETE("/:id", handlers.UserHandler.Delete)
			userGroup.PATCH("/:id/status", handlers.UserHandler.UpdateStatus)
			userGroup.PATCH("/:id/password", handlers.UserHandler.ResetPassword)
		}

		// 管理员用户管理
		adminUserGroup := protectedGroup.Group("/admin-users")
		{
			adminUserGroup.GET("", handlers.AdminUserHandler.GetList)
			adminUserGroup.GET("/:id", handlers.AdminUserHandler.GetByID)
			adminUserGroup.POST("", handlers.AdminUserHandler.Create)
			adminUserGroup.PUT("/:id", handlers.AdminUserHandler.Update)
			adminUserGroup.DELETE("/:id", handlers.AdminUserHandler.Delete)
			adminUserGroup.PATCH("/:id/status", handlers.AdminUserHandler.UpdateStatus)
		}

		// 角色管理
		roleGroup := protectedGroup.Group("/roles")
		{
			roleGroup.GET("", handlers.AdminRoleHandler.GetRoleList)
			roleGroup.GET("/permissions", handlers.AdminRoleHandler.GetPermissionTree)
			roleGroup.GET("/:id", handlers.AdminRoleHandler.GetRoleByID)
			roleGroup.POST("", handlers.AdminRoleHandler.CreateRole)
			roleGroup.PUT("/:id", handlers.AdminRoleHandler.UpdateRole)
			roleGroup.PUT("/:id/permissions", handlers.AdminRoleHandler.UpdateRolePermissions)
			roleGroup.DELETE("/:id", handlers.AdminRoleHandler.DeleteRole)
			roleGroup.POST("/assign", handlers.AdminRoleHandler.AssignRolesToUser)
		}

		// 系统管理
		systemGroup := protectedGroup.Group("/system")
		{
			systemGroup.GET("/info", handlers.SystemHandler.GetInfo)
			systemGroup.GET("/stats", handlers.SystemHandler.GetStats)
			systemGroup.GET("/cache/stats", handlers.SystemHandler.GetCacheStats)
			systemGroup.DELETE("/cache", handlers.SystemHandler.ClearCache)
		}

		// 配置管理
		configGroup := protectedGroup.Group("/configs")
		{
			configGroup.GET("", handlers.SystemHandler.GetConfigList)
			configGroup.GET("/:id", handlers.SystemHandler.GetConfig)
			configGroup.POST("", handlers.SystemHandler.CreateConfig)
			configGroup.PUT("/:id", handlers.SystemHandler.UpdateConfig)
			configGroup.DELETE("/:id", handlers.SystemHandler.DeleteConfig)
		}

		// 日志管理
		logGroup := protectedGroup.Group("/logs")
		{
			logGroup.GET("", handlers.SystemHandler.GetLogList)
			logGroup.DELETE("", handlers.SystemHandler.ClearLogs)
		}
	}
}

// Handlers 处理器集合
type Handlers struct {
	AuthHandler      *handler.AuthHandler
	UserHandler      *handler.UserHandler
	AdminUserHandler *handler.AdminUserHandler
	AdminRoleHandler *handler.AdminRoleHandler
	SystemHandler    *handler.SystemHandler
}
