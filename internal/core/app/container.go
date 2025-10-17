package app

import (
	"bico-admin/internal/admin"
	adminHandler "bico-admin/internal/admin/handler"
	adminService "bico-admin/internal/admin/service"
	"bico-admin/internal/api"
	"bico-admin/internal/core/cache"
	"bico-admin/internal/core/config"
	"bico-admin/internal/core/db"
	"bico-admin/internal/core/middleware"
	"bico-admin/internal/core/server"
	"bico-admin/internal/shared/jwt"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

// BuildContainer 构建 DI 容器
func BuildContainer(configPath string) (*dig.Container, error) {
	container := dig.New()

	// 按模块注册依赖（模块化 + 批量错误处理）
	providers := []interface{}{
		// 基础设施层
		func() (*config.Config, error) { return config.LoadConfig(configPath) },
		provideDatabase,
		provideGinEngine,
		provideCache,
		provideJWT,

		// 服务层
		provideAuthService,
		provideConfigService,

		// 处理层
		adminHandler.NewAuthHandler,
		adminHandler.NewCommonHandler,

		// 路由层
		provideAdminRouter,
		api.NewRouter,

		// 应用实例
		NewApp,
	}

	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			return nil, err
		}
	}

	return container, nil
}

// provideDatabase 提供数据库连接
func provideDatabase(cfg *config.Config) (*gorm.DB, error) {
	return db.InitDB(&cfg.Database)
}

// provideGinEngine 提供 Gin 引擎
func provideGinEngine(cfg *config.Config) *gin.Engine {
	return server.NewServer(&cfg.Server)
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

// provideConfigService 提供配置服务接口（dig 要求返回接口类型）
func provideConfigService(cfg *config.Config) adminService.IConfigService {
	return adminService.NewConfigService(cfg)
}

// AdminRouterParams 使用 dig.In 简化依赖注入（最佳实践✅）
type AdminRouterParams struct {
	dig.In
	AuthHandler   *adminHandler.AuthHandler
	CommonHandler *adminHandler.CommonHandler
	JWTManager    *jwt.JWTManager
	AuthService   adminService.IAuthService
}

// provideAdminRouter 提供 Admin 路由（使用 dig.In 简化参数）
func provideAdminRouter(params AdminRouterParams) *admin.Router {
	return admin.NewRouter(
		params.AuthHandler,
		params.CommonHandler,
		middleware.JWTAuth(params.JWTManager, params.AuthService),
	)
}
