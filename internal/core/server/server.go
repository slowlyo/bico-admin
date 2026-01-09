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
			// 有业务错误时，明确标记为失败，便于控制台快速区分
			logger.Warn("http_fail", append(fields, zap.String("result", "fail"), zap.String("errors", c.Errors.String()))...)
			return
		}

		if status >= 500 {
			// 5xx 认为是服务端失败
			logger.Error("http_fail", append(fields, zap.String("result", "fail"))...)
			return
		}
		if status >= 400 {
			// 4xx 认为是请求失败
			logger.Warn("http_fail", append(fields, zap.String("result", "fail"))...)
			return
		}

		// 非 debug 环境默认不输出成功请求日志，避免刷屏
		if !isDebug {
			return
		}
		logger.Info("http_success", append(fields, zap.String("result", "success"))...)
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
		serveEmbedStatic(engine, cfg.Server.AdminPath, embedFS)
	}
}

// serveEmbedStatic 服务嵌入的前端静态文件
func serveEmbedStatic(engine *gin.Engine, adminPath string, embedFS embed.FS) {
	if adminPath == "" || adminPath == "/" {
		adminPath = "/"
	}

	subFS, err := fs.Sub(embedFS, "dist")
	if err != nil {
		panic("failed to create sub filesystem: " + err.Error())
	}

	fileServer := http.FileServer(http.FS(subFS))

	// 如果配置了特定路径，则挂载到该路径下
	if adminPath != "/" {
		engine.StaticFS(adminPath, http.FS(subFS))
	}

	engine.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// API 路由返回 JSON 404
		if len(path) >= 10 && path[:10] == "/admin-api" {
			response.NotFound(c, "路由不存在")
			return
		}
		if len(path) >= 4 && path[:4] == "/api" {
			response.NotFound(c, "路由不存在")
			return
		}

		// 如果访问的是根路径且配置了 admin_path，则跳转到 admin_path
		if path == "/" && adminPath != "/" {
			c.Redirect(http.StatusMovedPermanently, adminPath)
			return
		}

		// 处理 SPA 路由：如果请求的路径不是文件，则返回 index.html
		if adminPath != "/" {
			// 如果是以 adminPath 开头的请求，且不是静态资源请求（简单判断，通常静态资源有后缀）
			// 这里交给 fileServer 处理，它会自动处理静态文件
			// 如果文件不存在，则重定向到 index.html 以支持 SPA
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		// 其他请求交给文件服务器处理
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}
