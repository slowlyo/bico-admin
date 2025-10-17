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

	// 加载配置
	if err := container.Provide(func() (*config.Config, error) {
		return config.LoadConfig(configPath)
	}); err != nil {
		return nil, err
	}

	// 提供数据库连接
	if err := container.Provide(func(cfg *config.Config) (*gorm.DB, error) {
		return db.InitDB(&cfg.Database)
	}); err != nil {
		return nil, err
	}

	// 提供 Gin 引擎
	if err := container.Provide(func(cfg *config.Config) *gin.Engine {
		return server.NewServer(&cfg.Server)
	}); err != nil {
		return nil, err
	}

	// 提供缓存
	if err := container.Provide(func(cfg *config.Config) (cache.Cache, error) {
		return cache.NewCache(&cfg.Cache)
	}); err != nil {
		return nil, err
	}

	// 提供 JWT 管理器
	if err := container.Provide(func(cfg *config.Config) *jwt.JWTManager {
		return jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)
	}); err != nil {
		return nil, err
	}

	// 提供 AuthService
	if err := container.Provide(func(db *gorm.DB, jwtManager *jwt.JWTManager, cacheInstance cache.Cache) *adminService.AuthService {
		return adminService.NewAuthService(db, jwtManager, cacheInstance)
	}); err != nil {
		return nil, err
	}

	// 提供 AuthHandler
	if err := container.Provide(func(authService *adminService.AuthService) *adminHandler.AuthHandler {
		return adminHandler.NewAuthHandler(authService)
	}); err != nil {
		return nil, err
	}

	// 提供 ConfigService
	if err := container.Provide(func(cfg *config.Config) *adminService.ConfigService {
		return adminService.NewConfigService(cfg)
	}); err != nil {
		return nil, err
	}

	// 提供 CommonHandler
	if err := container.Provide(func(configService *adminService.ConfigService) *adminHandler.CommonHandler {
		return adminHandler.NewCommonHandler(configService)
	}); err != nil {
		return nil, err
	}

	// 提供路由
	if err := container.Provide(func(authHandler *adminHandler.AuthHandler, commonHandler *adminHandler.CommonHandler, jwtManager *jwt.JWTManager, authService *adminService.AuthService) *admin.Router {
		return admin.NewRouter(authHandler, commonHandler, middleware.JWTAuth(jwtManager, authService))
	}); err != nil {
		return nil, err
	}

	if err := container.Provide(func() *api.Router {
		return api.NewRouter()
	}); err != nil {
		return nil, err
	}

	// 提供应用实例
	if err := container.Provide(NewApp); err != nil {
		return nil, err
	}

	return container, nil
}

