package admin

import (
	"bico-admin/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// Router 实现路由注册
type Router struct{}

// NewRouter 创建路由实例
func NewRouter() *Router {
	return &Router{}
}

// Register 注册路由
func (r *Router) Register(engine *gin.Engine) {
	admin := engine.Group("/admin")
	{
		admin.GET("/menus", func(c *gin.Context) {
			c.JSON(200, response.Success([]string{}))
		})
	}
}
