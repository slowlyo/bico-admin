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

	// 公开路由（无需认证）
	admin.POST("/login", r.authHandler.Login)
	admin.GET("/app-config", r.commonHandler.GetAppConfig)

	// 需要认证的路由
	admin.POST("/logout", r.jwtAuth, r.authHandler.Logout)
	admin.GET("/current-user", r.jwtAuth, r.authHandler.CurrentUser)

	// 临时路由
	admin.GET("/menus", func(c *gin.Context) {
		c.JSON(200, response.Success([]string{}))
	})
}
