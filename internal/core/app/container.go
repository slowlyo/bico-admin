package app

import (
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
	providers := []interface{}{
		// 基础设施层
		func() (*config.Config, error) { return config.LoadConfig(configPath) },
		provideLogger,
		provideDatabase,
		provideGinEngine,
		provideCache,
		provideJWT,
		provideUploader,
		provideScheduler,
		provideCaptcha,

		// 服务层
		provideAuthService,
		provideConfigService,
		adminService.NewAdminUserService,
		adminService.NewAdminRoleService,

		// 处理层
		adminHandler.NewAuthHandler,
		adminHandler.NewCommonHandler,
		adminHandler.NewAdminUserHandler,
		adminHandler.NewAdminRoleHandler,

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

// provideUploader 提供文件上传器
func provideUploader(cfg *config.Config) (upload.Uploader, error) {
	uploaderConfig := &upload.UploaderConfig{
		Driver:       cfg.Upload.Driver,
		MaxSize:      cfg.Upload.MaxSize,
		AllowedTypes: cfg.Upload.AllowedTypes,
		LocalConfig: upload.LocalConfig{
			BasePath:  cfg.Upload.Local.BasePath,
			URLPrefix: cfg.Upload.Local.URLPrefix,
		},
		QiniuConfig: upload.QiniuConfig{
			AccessKey:    cfg.Upload.Qiniu.AccessKey,
			SecretKey:    cfg.Upload.Qiniu.SecretKey,
			Bucket:       cfg.Upload.Qiniu.Bucket,
			Domain:       cfg.Upload.Qiniu.Domain,
			Zone:         cfg.Upload.Qiniu.Zone,
			UseHTTPS:     cfg.Upload.Qiniu.UseHTTPS,
			UseCDNDomain: cfg.Upload.Qiniu.UseCDNDomain,
		},
		AliyunConfig: upload.AliyunConfig{
			AccessKeyId:     cfg.Upload.Aliyun.AccessKeyId,
			AccessKeySecret: cfg.Upload.Aliyun.AccessKeySecret,
			Bucket:          cfg.Upload.Aliyun.Bucket,
			Endpoint:        cfg.Upload.Aliyun.Endpoint,
			Domain:          cfg.Upload.Aliyun.Domain,
			UseHTTPS:        cfg.Upload.Aliyun.UseHTTPS,
		},
	}
	return upload.NewUploader(uploaderConfig)
}

// provideScheduler 提供定时任务调度器
func provideScheduler(zapLogger *zap.Logger) *job.Scheduler {
	return job.NewScheduler(zapLogger)
}

// provideCaptcha 提供验证码实例
func provideCaptcha(cacheInstance cache.Cache) *captcha.Captcha {
	return captcha.NewCaptcha(cacheInstance)
}

// AdminRouterParams 使用 dig.In 简化依赖注入
type AdminRouterParams struct {
	dig.In
	AuthHandler      *adminHandler.AuthHandler
	CommonHandler    *adminHandler.CommonHandler
	AdminUserHandler *adminHandler.AdminUserHandler
	AdminRoleHandler *adminHandler.AdminRoleHandler
	JWTManager       *jwt.JWTManager
	AuthService      adminService.IAuthService
	DB               *gorm.DB
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
		params.AdminUserHandler,
		params.AdminRoleHandler,
		middleware.JWTAuth(params.JWTManager, params.AuthService),
		permMiddleware,
		userStatusMiddleware,
	)
}
