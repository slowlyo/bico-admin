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
			roleGroup.GET("/options", handlers.AdminRoleHandler.GetRoleOptions)
			roleGroup.GET("/permissions", handlers.AdminRoleHandler.GetPermissionTree)
			roleGroup.GET("/:id", handlers.AdminRoleHandler.GetRoleByID)
			roleGroup.POST("", handlers.AdminRoleHandler.CreateRole)
			roleGroup.PUT("/:id", handlers.AdminRoleHandler.UpdateRole)
			roleGroup.PATCH("/:id/status", handlers.AdminRoleHandler.UpdateRoleStatus)
			roleGroup.PUT("/:id/permissions", handlers.AdminRoleHandler.UpdateRolePermissions)
			roleGroup.DELETE("/:id", handlers.AdminRoleHandler.DeleteRole)
			roleGroup.POST("/assign", handlers.AdminRoleHandler.AssignRolesToUser)
		}

	}
}

// Handlers 处理器集合
type Handlers struct {
	AuthHandler      *handler.AuthHandler
	AdminUserHandler *handler.AdminUserHandler
	AdminRoleHandler *handler.AdminRoleHandler
}
