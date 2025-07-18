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

	// 调用所有注册的路由注册器
	for _, registrar := range routeRegistrars {
		registrar.RegisterRoutes(r, handlers)
	}
}

// Handlers 处理器集合
type Handlers struct {
	HelloHandler *handler.HelloHandler
}

// RouteRegistrar 路由注册器接口
// 用于支持动态路由注册，生成的路由代码可以实现此接口
type RouteRegistrar interface {
	RegisterRoutes(router *gin.Engine, handlers *Handlers)
}

// routeRegistrars 存储所有注册的路由注册器
var routeRegistrars []RouteRegistrar

// RegisterRouteRegistrar 注册路由注册器
// 生成的路由代码可以调用此函数来注册自己
func RegisterRouteRegistrar(registrar RouteRegistrar) {
	routeRegistrars = append(routeRegistrars, registrar)
}
