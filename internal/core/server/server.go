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
func RegisterRoutes(engine *gin.Engine, adminRouter, apiRouter Router, cfg *config.Config) {
	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, response.Success(gin.H{"status": "ok"}))
	})
	
	// 静态文件服务（用于访问上传的文件）
	if cfg.Upload.Driver == "local" {
		engine.Static(cfg.Upload.Local.ServePath, cfg.Upload.Local.BasePath)
	}
	
	// 注册模块路由
	adminRouter.Register(engine)
	apiRouter.Register(engine)
}

// Router 路由接口
type Router interface {
	Register(engine *gin.Engine)
}
