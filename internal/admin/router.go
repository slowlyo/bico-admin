package admin

import (
	"bico-admin/internal/admin/handler"
	"bico-admin/internal/shared/response"

	"github.com/gin-gonic/gin"
)

// Router 实现路由注册
type Router struct {
	authHandler   *handler.AuthHandler
	commonHandler *handler.CommonHandler
	jwtAuth       gin.HandlerFunc
}

// NewRouter 创建路由实例
func NewRouter(authHandler *handler.AuthHandler, commonHandler *handler.CommonHandler, jwtAuth gin.HandlerFunc) *Router {
	return &Router{
		authHandler:   authHandler,
		commonHandler: commonHandler,
		jwtAuth:       jwtAuth,
	}
}

// Register 注册路由
func (r *Router) Register(engine *gin.Engine) {
	admin := engine.Group("/admin-api")

	// 认证相关路由
	auth := admin.Group("/auth")
	{
		// 公开路由（无需认证）
		auth.POST("/login", r.authHandler.Login)
		
		// 需要认证的路由
		auth.POST("/logout", r.jwtAuth, r.authHandler.Logout)
		auth.GET("/current-user", r.jwtAuth, r.authHandler.CurrentUser)
		auth.PUT("/profile", r.jwtAuth, r.authHandler.UpdateProfile)
		auth.PUT("/password", r.jwtAuth, r.authHandler.ChangePassword)
		auth.POST("/avatar", r.jwtAuth, r.authHandler.UploadAvatar)
	}

	// 应用配置
	admin.GET("/app-config", r.commonHandler.GetAppConfig)

	// 临时路由
	admin.GET("/menus", func(c *gin.Context) {
		c.JSON(200, response.Success([]string{}))
	})
}
