package server

import (
	"embed"
	"io/fs"
	"net/http"
	"runtime/debug"
	"strings"
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
	engine.RedirectTrailingSlash = false
	engine.RedirectFixedPath = false
	engine.Use(zapRecovery(zapLogger))
	engine.Use(zapAccessLogger(zapLogger, cfg.Mode == "debug"))

	// 添加 CORS 中间件
	engine.Use(middleware.CORS(cfg.AllowedOrigins))

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

// RegisterCoreRoutes 注册框架级路由。
//
// 说明：静态后台入口通过配置管理器按请求读取，修改 admin_path 后无需重启服务。
func RegisterCoreRoutes(engine *gin.Engine, configManager *config.ConfigManager, embedFS embed.FS) {
	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		response.SuccessWithData(c, gin.H{"status": "ok"})
	})

	// Swagger 文档，分别暴露对外 API 与后台管理 API。
	engine.GET("/swagger/api/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("api")))
	engine.GET("/swagger/admin/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("admin")))

	// 静态文件服务（用于访问上传的文件）
	cfg := configManager.GetConfig()
	if cfg.Upload.Driver == "local" {
		engine.Static(cfg.Upload.Local.ServePath, cfg.Upload.Local.BasePath)
	}

	// 前端静态文件服务
	if cfg.Server.EmbedStatic {
		serveEmbedStatic(engine, configManager, embedFS)
	}
}

// serveEmbedStatic 服务嵌入的前端静态文件。
//
// 说明：Gin 的路由表无法在运行时增删，因此通过 NoRoute 分发后台入口并在每次请求时读取最新路径。
func serveEmbedStatic(engine *gin.Engine, configManager *config.ConfigManager, embedFS embed.FS) {
	subFS, err := fs.Sub(embedFS, "dist")
	if err != nil {
		panic("failed to create sub filesystem: " + err.Error())
	}

	// 未命中业务路由时才进入静态资源分发，避免覆盖 API、上传和 Swagger 路由。
	handler := func(c *gin.Context) {
		// 强制不缓存，解决开发/调试期间的 301 缓存问题
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Header("X-Bico-Debug", "v2-fixed-redirects")

		path := c.Request.URL.Path
		prefix := normalizeAdminPath(configManager.GetConfig().Server.AdminPath)

		// API 未命中时保持 404，避免被 SPA 回退页掩盖接口错误。
		if strings.HasPrefix(path, "/admin-api") || strings.HasPrefix(path, "/api") {
			response.NotFound(c, "资源不存在")
			return
		}

		// 仅当前配置的入口及其子路径可以访问后台，旧入口在热更新后立即失效。
		if prefix != "" && path != prefix && !strings.HasPrefix(path, prefix+"/") {
			response.NotFound(c, "资源不存在")
			return
		}

		// 访问入口根路径时补齐斜杠，使相对资源始终以后台目录为基准解析。
		if prefix != "" && path == prefix {
			c.Redirect(http.StatusFound, prefix+"/")
			return
		}

		// 移除后台入口前缀后，才能从嵌入的 dist 目录读取对应资源。
		filePath := path
		if prefix != "" && strings.HasPrefix(path, prefix+"/") {
			filePath = strings.TrimPrefix(path, prefix)
		}
		filePath = strings.TrimPrefix(filePath, "/")

		// 目录或空路径指向 index.html，由前端 Hash 路由继续处理页面切换。
		if filePath == "" || strings.HasSuffix(filePath, "/") {
			filePath = "index.html"
		}

		// 优先读取静态文件；没有扩展名的路径按 SPA 入口回退。
		f, err := subFS.Open(filePath)
		if err != nil {
			// 如果是带后缀的静态资源（如 .js, .css），文件不存在则直接 404
			if strings.Contains(filePath, ".") && !strings.HasSuffix(filePath, "index.html") {
				response.NotFound(c, "资源不存在")
				return
			}
			// 否则作为 SPA 路由，返回 index.html
			filePath = "index.html"
			f, err = subFS.Open(filePath)
			if err != nil {
				response.NotFound(c, "入口文件不存在")
				return
			}
		}
		defer f.Close()

		fi, _ := f.Stat()
		if fi.IsDir() {
			filePath = "index.html"
		}

		// index.html 手动输出，避免 FileServer 对入口文件触发 301。
		// 对于 index.html，必须手动读取并返回，绝对不能使用 http.FileServer 或 c.FileFromFS
		// 因为它们检测到 index.html 时会尝试做 301 重定向，这是导致死循环的根源
		if strings.HasSuffix(filePath, "index.html") {
			content, err := fs.ReadFile(subFS, "index.html")
			if err != nil {
				response.NotFound(c, "读取入口文件失败")
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", content)
			return
		}

		// 其他普通资源使用 FileFromFS
		c.FileFromFS(filePath, http.FS(subFS))
	}

	engine.NoRoute(handler)
}

// normalizeAdminPath 统一后台入口格式。
//
// 说明：空值与根路径使用空前缀，其他值确保仅保留一个前导斜杠且没有尾随斜杠。
func normalizeAdminPath(adminPath string) string {
	path := strings.Trim(adminPath, "/")
	if path == "" {
		return ""
	}

	return "/" + path
}
