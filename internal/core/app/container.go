package app

import (
	"fmt"
	"math"

	"bico-admin/internal/admin"
	adminHandler "bico-admin/internal/admin/handler"
	adminMiddleware "bico-admin/internal/admin/middleware"
	adminService "bico-admin/internal/admin/service"
	"bico-admin/internal/api"
	"bico-admin/internal/core/cache"
	"bico-admin/internal/core/config"
	"bico-admin/internal/core/db"
	"bico-admin/internal/core/logger"
	"bico-admin/internal/core/middleware"
	"bico-admin/internal/core/server"
	"bico-admin/internal/core/upload"
	"bico-admin/internal/job"
	"bico-admin/internal/pkg/captcha"
	"bico-admin/internal/pkg/crud"
	"bico-admin/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// BuildContainer 构建 DI 容器
func BuildContainer(configPath string) (*dig.Container, error) {
	container := dig.New()

	// 按模块注册依赖（模块化 + 批量错误处理）
	providers := []struct {
		provider interface{}
		name     string
	}{
		// 基础设施层
		{provideConfig(configPath), "Config"},
		{provideLogger, "Logger"},
		{provideConfigManager(configPath), "ConfigManager"},
		{provideDatabase, "Database"},
		{provideCache, "Cache"},
		{provideJWT, "JWT"},
		{provideRateLimiter, "RateLimiter"},
		{provideGinEngine, "GinEngine"},
		{provideUploader, "Uploader"},
		{provideScheduler, "Scheduler"},
		{provideCaptcha, "Captcha"},

		// 服务层（核心服务）
		{provideAuthService, "AuthService"},
		{provideConfigService, "ConfigService"},

		// 非 CRUD 处理器（需手动注册）
		{adminHandler.NewAuthHandler, "AuthHandler"},
		{adminHandler.NewCommonHandler, "CommonHandler"},

		// 路由层
		{provideAdminRouter, "AdminRouter"},
		{api.NewRouter, "ApiRouter"},

		// 模块路由注册器
		{provideModuleRouter, "ModuleRouter"},

		// 应用实例
		{NewApp, "App"},
	}

	for _, p := range providers {
		if err := container.Provide(p.provider); err != nil {
			return nil, fmt.Errorf("provide %s failed: %w", p.name, err)
		}
	}

	// 自动注册 CRUD 模块（admin 分组）
	if err := crud.ProvideModules(container); err != nil {
		return nil, fmt.Errorf("provide crud modules failed: %w", err)
	}

	return container, nil
}

// provideConfig 提供配置实例
func provideConfig(configPath string) func() (*config.Config, error) {
	return func() (*config.Config, error) {
		return config.LoadConfig(configPath)
	}
}

// provideConfigManager 提供配置管理器（支持热更新）
func provideConfigManager(configPath string) func(*zap.Logger) (*config.ConfigManager, error) {
	return func(zapLogger *zap.Logger) (*config.ConfigManager, error) {
		return config.NewConfigManager(configPath, zapLogger)
	}
}

// provideLogger 提供日志实例
func provideLogger(cfg *config.Config) (*zap.Logger, error) {
	return logger.InitLogger(&cfg.Log)
}

// provideDatabase 提供数据库连接
func provideDatabase(cfg *config.Config, zapLogger *zap.Logger) (*gorm.DB, error) {
	isDebug := cfg.Server.Mode == "debug"
	return db.InitDB(&cfg.Database, zapLogger, isDebug)
}

// provideGinEngine 提供 Gin 引擎
func provideGinEngine(cfg *config.Config, rateLimiter *middleware.RateLimiter) *gin.Engine {
	return server.NewServer(&cfg.Server, rateLimiter)
}

// provideCache 提供缓存实例
func provideCache(cfg *config.Config) (cache.Cache, error) {
	return cache.NewCache(&cfg.Cache)
}

// provideJWT 提供 JWT 管理器
func provideJWT(cfg *config.Config) *jwt.JWTManager {
	return jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)
}

// provideAuthService 提供认证服务接口（dig 要求返回接口类型）
func provideAuthService(database *gorm.DB, jwtManager *jwt.JWTManager, cacheInstance cache.Cache) adminService.IAuthService {
	return adminService.NewAuthService(database, jwtManager, cacheInstance)
}

// provideConfigService 提供配置服务接口（dig 要求返回接口类型，支持热更新）
func provideConfigService(cm *config.ConfigManager) adminService.IConfigService {
	return adminService.NewConfigService(cm)
}

// provideUploader 提供文件上传器
func provideUploader(cfg *config.Config) (upload.Uploader, error) {
	return upload.NewUploader(upload.ConfigFromAppConfig(cfg))
}

// provideScheduler 提供定时任务调度器
func provideScheduler(zapLogger *zap.Logger) *job.Scheduler {
	return job.NewScheduler(zapLogger)
}

// provideCaptcha 提供验证码实例
func provideCaptcha(cacheInstance cache.Cache) *captcha.Captcha {
	return captcha.NewCaptcha(cacheInstance)
}

const (
	disabledRateLimitValue = math.MaxInt32
)

// provideRateLimiter 提供限流器（支持热更新）
func provideRateLimiter(cm *config.ConfigManager) *middleware.RateLimiter {
	cfg := cm.GetConfig()
	if !cfg.RateLimit.Enabled {
		return middleware.NewRateLimiter(disabledRateLimitValue, disabledRateLimitValue)
	}
	return middleware.NewRateLimiter(cfg.RateLimit.RPS, cfg.RateLimit.Burst)
}

// AdminRouterParams 使用 dig.In 简化依赖注入
type AdminRouterParams struct {
	dig.In
	AuthHandler   *adminHandler.AuthHandler
	CommonHandler *adminHandler.CommonHandler
	JWTManager    *jwt.JWTManager
	AuthService   adminService.IAuthService
	DB            *gorm.DB
}

// provideAdminRouter 提供 Admin 路由（使用 dig.In 简化参数）
func provideAdminRouter(params AdminRouterParams) *admin.Router {
	// 创建权限中间件
	permMiddleware := adminMiddleware.NewPermissionMiddleware(params.AuthService)
	// 创建用户状态中间件
	userStatusMiddleware := adminMiddleware.NewUserStatusMiddleware(params.DB)

	return admin.NewRouter(
		params.AuthHandler,
		params.CommonHandler,
		middleware.JWTAuth(params.JWTManager, params.AuthService),
		permMiddleware,
		userStatusMiddleware,
		params.DB,
	)
}

// ModuleRouterParams 模块路由参数
type ModuleRouterParams struct {
	dig.In
	JWTManager  *jwt.JWTManager
	AuthService adminService.IAuthService
	DB          *gorm.DB
}

// provideModuleRouter 提供模块路由注册器
func provideModuleRouter(params ModuleRouterParams) *crud.ModuleRouter {
	permMiddleware := adminMiddleware.NewPermissionMiddleware(params.AuthService)
	userStatusMiddleware := adminMiddleware.NewUserStatusMiddleware(params.DB)

	return crud.NewModuleRouter(
		middleware.JWTAuth(params.JWTManager, params.AuthService),
		permMiddleware,
		userStatusMiddleware.Check(),
	)
}
