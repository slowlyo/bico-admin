package frontend

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"bico-admin/pkg/config"
	"bico-admin/pkg/logger"
	"bico-admin/web"
)

// 初始化 MIME 类型
func init() {
	// 确保 JavaScript 模块文件有正确的 MIME 类型
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".mjs", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".json", "application/json")
	mime.AddExtensionType(".wasm", "application/wasm")
}

// getEmbeddedFileSystem 获取嵌入的文件系统
func getEmbeddedFileSystem() (http.FileSystem, error) {
	return web.GetFileSystem()
}

// isEmbedded 检查是否使用嵌入模式
func isEmbedded() bool {
	return web.IsEmbedded()
}

// Service 前端服务
type Service struct {
	config *config.FrontendConfig
}

// NewService 创建前端服务
func NewService(cfg *config.FrontendConfig) *Service {
	return &Service{
		config: cfg,
	}
}

// SetupRoutes 设置前端路由
func (s *Service) SetupRoutes(r *gin.Engine) error {
	if s.config.IsEmbedMode() {
		return s.setupEmbeddedRoutes(r)
	}
	return s.setupExternalRoutes(r)
}

// setupEmbeddedRoutes 设置嵌入式文件路由
func (s *Service) setupEmbeddedRoutes(r *gin.Engine) error {
	// 检查是否支持嵌入模式
	if !isEmbedded() {
		logger.Warn("嵌入模式未启用，回退到外部文件模式")
		return s.setupExternalRoutes(r)
	}

	// 获取嵌入的文件系统
	fileSystem, err := getEmbeddedFileSystem()
	if err != nil {
		logger.Error("获取嵌入文件系统失败: " + err.Error())
		return err
	}

	// 设置静态文件路由，使用自定义处理器确保正确的 MIME 类型
	r.GET("/assets/*filepath", s.serveEmbeddedAssets(fileSystem))

	// 设置主页和图标
	r.GET("/", s.serveEmbeddedIndex(fileSystem))
	r.GET("/favicon.ico", s.serveEmbeddedFile(fileSystem, "favicon.ico"))

	// 设置 NoRoute 处理器
	r.NoRoute(s.createNoRouteHandler(fileSystem, true))

	logger.Info("前端服务已启用 (嵌入模式)")
	return nil
}

// setupExternalRoutes 设置外部文件路由
func (s *Service) setupExternalRoutes(r *gin.Engine) error {
	// 检查静态文件目录是否存在
	staticDir := s.config.GetStaticDir()
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		logger.Warn("静态文件目录不存在: " + staticDir)
		return err
	}

	// 设置静态文件路由
	assetsDir := s.config.GetAssetsDir()
	r.Static("/assets", assetsDir)

	// 设置主页和图标
	indexFile := s.config.GetIndexFile()
	r.StaticFile("/", indexFile)
	r.StaticFile("/favicon.ico", filepath.Join(staticDir, "favicon.ico"))

	// 设置 NoRoute 处理器
	r.NoRoute(s.createNoRouteHandler(nil, false))

	logger.Info("前端服务已启用 (外部文件模式)")
	return nil
}

// serveEmbeddedIndex 服务嵌入的主页文件
func (s *Service) serveEmbeddedIndex(fileSystem http.FileSystem) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := fileSystem.Open("index.html")
		if err != nil {
			c.String(http.StatusNotFound, "Index file not found")
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to get file info")
			return
		}

		c.DataFromReader(http.StatusOK, stat.Size(), "text/html", file, nil)
	}
}

// serveEmbeddedAssets 服务嵌入的静态资源文件
func (s *Service) serveEmbeddedAssets(fileSystem http.FileSystem) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的文件路径，去掉前缀 "/assets/"
		filepath := strings.TrimPrefix(c.Param("filepath"), "/")
		if filepath == "" {
			c.String(http.StatusNotFound, "File not found")
			return
		}

		// 构建完整的文件路径
		fullPath := "assets/" + filepath

		file, err := fileSystem.Open(fullPath)
		if err != nil {
			c.String(http.StatusNotFound, "File not found")
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to get file info")
			return
		}

		// 根据文件扩展名设置正确的 Content-Type
		contentType := s.getContentType(filepath)

		c.DataFromReader(http.StatusOK, stat.Size(), contentType, file, nil)
	}
}

// getContentType 根据文件扩展名获取 Content-Type
func (s *Service) getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".js", ".mjs":
		return "application/javascript"
	case ".css":
		return "text/css"
	case ".json":
		return "application/json"
	case ".wasm":
		return "application/wasm"
	case ".html", ".htm":
		return "text/html"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	case ".eot":
		return "application/vnd.ms-fontobject"
	default:
		return "application/octet-stream"
	}
}

// serveEmbeddedFile 服务嵌入的文件
func (s *Service) serveEmbeddedFile(fileSystem http.FileSystem, filename string) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := fileSystem.Open(filename)
		if err != nil {
			c.String(http.StatusNotFound, "File not found")
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to get file info")
			return
		}

		// 根据文件扩展名设置 Content-Type
		contentType := s.getContentType(filename)

		c.DataFromReader(http.StatusOK, stat.Size(), contentType, file, nil)
	}
}

// createNoRouteHandler 创建 NoRoute 处理器
func (s *Service) createNoRouteHandler(fileSystem http.FileSystem, isEmbedded bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 如果是API请求，返回404
		if strings.HasPrefix(path, "/admin-api") ||
			strings.HasPrefix(path, "/master") ||
			strings.HasPrefix(path, "/api") ||
			strings.HasPrefix(path, "/uploads") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}

		// 其他路径返回前端index.html
		if isEmbedded {
			s.serveEmbeddedIndex(fileSystem)(c)
		} else {
			c.File(s.config.GetIndexFile())
		}
	}
}
