package routes

import (
	"github.com/gin-gonic/gin"

	"bico-admin/internal/api/handler"
	sharedMiddleware "bico-admin/internal/shared/middleware"
)

// RegisterRoutes 注册API端路由
func RegisterRoutes(r *gin.Engine, handlers *Handlers) {
	// API端路由组
	apiGroup := r.Group("/api")
	apiGroup.Use(sharedMiddleware.CORS()) // 跨域中间件

	// 简单的Hello接口
	{
		apiGroup.GET("/hello", handlers.HelloHandler.Hello)
	}
}

// Handlers 处理器集合
type Handlers struct {
	HelloHandler *handler.HelloHandler
}
