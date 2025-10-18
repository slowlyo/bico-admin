package server

import (
	"bico-admin/internal/core/config"
	"bico-admin/internal/core/middleware"
	"bico-admin/internal/pkg/response"
	"embed"
	"io/fs"
	"net/http"
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
func RegisterRoutes(engine *gin.Engine, adminRouter, apiRouter Router, cfg *config.Config, embedFS embed.FS) {
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
	
	// 前端静态文件服务
	if cfg.Server.EmbedStatic {
		serveEmbedStatic(engine, embedFS)
	}
}

// Router 路由接口
type Router interface {
	Register(engine *gin.Engine)
}

// serveEmbedStatic 服务嵌入的前端静态文件
func serveEmbedStatic(engine *gin.Engine, embedFS embed.FS) {
	subFS, err := fs.Sub(embedFS, "dist")
	if err != nil {
		panic("failed to create sub filesystem: " + err.Error())
	}
	
	fileServer := http.FileServer(http.FS(subFS))
	
	engine.NoRoute(func(c *gin.Context) {
		// API 路由返回 JSON 404
		if len(c.Request.URL.Path) >= 10 && c.Request.URL.Path[:10] == "/admin-api" {
			c.JSON(404, response.Error(404, "路由不存在"))
			return
		}
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(404, response.Error(404, "路由不存在"))
			return
		}
		
		// 其他请求交给文件服务器处理
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}
