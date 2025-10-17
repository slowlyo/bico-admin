package server

import (
	"bico-admin/internal/core/config"
	"bico-admin/internal/core/middleware"
	"bico-admin/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// NewServer 创建 Gin 服务器
func NewServer(cfg *config.ServerConfig) *gin.Engine {
	gin.SetMode(cfg.Mode)
	
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	
	// 添加 CORS 中间件
	engine.Use(middleware.CORS())
	
	return engine
}

// RegisterRoutes 注册所有路由
func RegisterRoutes(engine *gin.Engine, adminRouter, apiRouter Router) {
	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, response.Success(gin.H{"status": "ok"}))
	})
	
	// 注册模块路由
	adminRouter.Register(engine)
	apiRouter.Register(engine)
}

// Router 路由接口
type Router interface {
	Register(engine *gin.Engine)
}
