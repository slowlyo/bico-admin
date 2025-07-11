package routes

import (
	"github.com/gin-gonic/gin"

	"bico-admin/internal/master/handler"
	"bico-admin/internal/shared/middleware"
)

// RegisterRoutes 注册master端路由
func RegisterRoutes(r *gin.Engine, handlers *Handlers) {
	// master端路由组
	masterGroup := r.Group("/master")
	masterGroup.Use(middleware.CORS()) // 跨域中间件

	// 简单的Hello接口
	{
		masterGroup.GET("/hello", handlers.HelloHandler.Hello)
	}
}

// Handlers 处理器集合
type Handlers struct {
	HelloHandler *handler.HelloHandler
}
