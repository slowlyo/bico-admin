package server

import (
	"embed"
	"io/fs"
	"net/http"
	"runtime/debug"
	"time"

	"bico-admin/internal/core/config"
	"bico-admin/internal/core/middleware"
	"bico-admin/internal/pkg/response"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// NewServer 创建 Gin 服务器
func NewServer(cfg *config.ServerConfig, rateLimiter *middleware.RateLimiter, zapLogger *zap.Logger) *gin.Engine {
	gin.SetMode(cfg.Mode)

	engine := gin.New()
	engine.Use(zapRecovery(zapLogger))
	engine.Use(zapAccessLogger(zapLogger, cfg.Mode == "debug"))

	// 添加 CORS 中间件
	engine.Use(middleware.CORS())

	// 添加全局限流中间件（如果启用）
	if rateLimiter != nil {
		engine.Use(rateLimiter.RateLimit())
	}

	return engine
}

func zapAccessLogger(logger *zap.Logger, isDebug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := c.ClientIP()
		ua := c.Request.UserAgent()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("client_ip", clientIP),
			zap.Duration("latency", latency),
			zap.String("user_agent", ua),
		}

		if len(c.Errors) > 0 {
			logger.Warn("http", append(fields, zap.String("errors", c.Errors.String()))...)
			return
		}

		if status >= 500 {
			logger.Error("http", fields...)
			return
		}
		if status >= 400 {
			logger.Warn("http", fields...)
			return
		}
		if !isDebug {
			return
		}
		logger.Info("http", fields...)
	}
}

func zapRecovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error(
					"panic",
					zap.Any("error", rec),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("client_ip", c.ClientIP()),
					zap.ByteString("stack", debug.Stack()),
				)
				response.ErrorWithCode(c, 500, "服务器内部错误")
				c.Abort()
			}
		}()
		c.Next()
	}
}

// RegisterCoreRoutes 注册框架级路由
//
// 说明：业务路由由各模块自行注册，core 只处理健康检查/Swagger/静态资源等框架能力。
func RegisterCoreRoutes(engine *gin.Engine, cfg *config.Config, embedFS embed.FS) {
	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		response.SuccessWithData(c, gin.H{"status": "ok"})
	})

	// Swagger 文档
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 静态文件服务（用于访问上传的文件）
	if cfg.Upload.Driver == "local" {
		engine.Static(cfg.Upload.Local.ServePath, cfg.Upload.Local.BasePath)
	}

	// 前端静态文件服务
	if cfg.Server.EmbedStatic {
		serveEmbedStatic(engine, embedFS)
	}
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
			response.NotFound(c, "路由不存在")
			return
		}
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			response.NotFound(c, "路由不存在")
			return
		}

		// 其他请求交给文件服务器处理
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}
