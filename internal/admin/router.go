package admin

import (
	"bico-admin/internal/admin/consts"
	"bico-admin/internal/admin/handler"
	"bico-admin/internal/admin/middleware"

	"github.com/gin-gonic/gin"
)

// Router 实现路由注册
type Router struct {
	authHandler        *handler.AuthHandler
	commonHandler      *handler.CommonHandler
	adminUserHandler   *handler.AdminUserHandler
	adminRoleHandler   *handler.AdminRoleHandler
	jwtAuth            gin.HandlerFunc
	permMiddleware     *middleware.PermissionMiddleware
	userStatusMiddleware *middleware.UserStatusMiddleware
}

// NewRouter 创建路由实例
func NewRouter(
	authHandler *handler.AuthHandler,
	commonHandler *handler.CommonHandler,
	adminUserHandler *handler.AdminUserHandler,
	adminRoleHandler *handler.AdminRoleHandler,
	jwtAuth gin.HandlerFunc,
	permMiddleware *middleware.PermissionMiddleware,
	userStatusMiddleware *middleware.UserStatusMiddleware,
) *Router {
	return &Router{
		authHandler:        authHandler,
		commonHandler:      commonHandler,
		adminUserHandler:   adminUserHandler,
		adminRoleHandler:   adminRoleHandler,
		jwtAuth:            jwtAuth,
		permMiddleware:     permMiddleware,
		userStatusMiddleware: userStatusMiddleware,
	}
}

// Register 注册路由
func (r *Router) Register(engine *gin.Engine) {
	admin := engine.Group("/admin-api")

	// 公开路由
	{
		// 登录
		admin.POST("/auth/login", r.authHandler.Login)
		// 应用配置
		admin.GET("/app-config", r.commonHandler.GetAppConfig)
	}

	// 需要认证的路由
	authorized := admin.Group("", r.jwtAuth, r.userStatusMiddleware.Check())
	{
		// 认证相关
		auth := authorized.Group("/auth")
		{
			auth.POST("/logout", r.authHandler.Logout)
			auth.GET("/current-user", r.authHandler.CurrentUser)
			auth.PUT("/profile", r.authHandler.UpdateProfile)
			auth.PUT("/password", r.authHandler.ChangePassword)
			auth.POST("/avatar", r.authHandler.UploadAvatar)
		}

		// 系统管理 - 用户管理
		adminUsers := authorized.Group("/admin-users")
		{
			adminUsers.GET("", r.permMiddleware.RequirePermission(consts.PermAdminUserList), r.adminUserHandler.List)
			adminUsers.GET("/:id", r.permMiddleware.RequirePermission(consts.PermAdminUserList), r.adminUserHandler.Get)
			adminUsers.POST("", r.permMiddleware.RequirePermission(consts.PermAdminUserCreate), r.adminUserHandler.Create)
			adminUsers.PUT("/:id", r.permMiddleware.RequirePermission(consts.PermAdminUserEdit), r.adminUserHandler.Update)
			adminUsers.DELETE("/:id", r.permMiddleware.RequirePermission(consts.PermAdminUserDelete), r.adminUserHandler.Delete)
		}

		// 系统管理 - 角色管理
		adminRoles := authorized.Group("/admin-roles")
		{
			adminRoles.GET("", r.permMiddleware.RequirePermission(consts.PermAdminRoleList), r.adminRoleHandler.List)
			adminRoles.GET("/all", r.permMiddleware.RequirePermission(consts.PermAdminRoleList), r.adminRoleHandler.GetAll)
			adminRoles.GET("/permissions", r.permMiddleware.RequirePermission(consts.PermAdminRoleList), r.adminRoleHandler.GetAllPermissions)
			adminRoles.GET("/:id", r.permMiddleware.RequirePermission(consts.PermAdminRoleList), r.adminRoleHandler.Get)
			adminRoles.POST("", r.permMiddleware.RequirePermission(consts.PermAdminRoleCreate), r.adminRoleHandler.Create)
			adminRoles.PUT("/:id", r.permMiddleware.RequirePermission(consts.PermAdminRoleEdit), r.adminRoleHandler.Update)
			adminRoles.DELETE("/:id", r.permMiddleware.RequirePermission(consts.PermAdminRoleDelete), r.adminRoleHandler.Delete)
			adminRoles.PUT("/:id/permissions", r.permMiddleware.RequirePermission(consts.PermAdminRolePermission), r.adminRoleHandler.UpdatePermissions)
		}
	}
}
