package admin

import (
	"bico-admin/internal/shared/response"

	"github.com/gin-gonic/gin"
)

// Router 实现路由注册
type Router struct{
	authHandler interface{}
}

// NewRouter 创建路由实例
func NewRouter(authHandler interface{}) *Router {
	return &Router{authHandler: authHandler}
}

// Register 注册路由
func (r *Router) Register(engine *gin.Engine) {
	admin := engine.Group("/admin-api")
	{
		admin.GET("/menus", func(c *gin.Context) {
			c.JSON(200, response.Success([]string{}))
		})
		
		handler := r.authHandler.(interface {
			Login(c interface {
				ShouldBindJSON(obj interface{}) error
				JSON(code int, obj interface{})
			})
			Logout(c interface {
				GetHeader(key string) string
				JSON(code int, obj interface{})
			})
		})
		
		admin.POST("/login", func(c *gin.Context) {
			handler.Login(c)
		})
		admin.POST("/logout", func(c *gin.Context) {
			handler.Logout(c)
		})
	}
}
