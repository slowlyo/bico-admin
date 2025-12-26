package app

import (
	"fmt"
	"math"

	"bico-admin/internal/core/cache"
	"bico-admin/internal/core/config"
	"bico-admin/internal/core/db"
	"bico-admin/internal/core/logger"
	"bico-admin/internal/core/middleware"
	"bico-admin/internal/core/scheduler"
	"bico-admin/internal/core/server"
	"bico-admin/internal/core/upload"
	"bico-admin/internal/pkg/captcha"
	"bico-admin/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AppContext 应用运行时上下文
//
// 说明：只包含框架与基础设施层能力，业务模块只能依赖这里暴露的对象。
// 业务模块需要的 service/handler/router 必须在模块内部自行装配。
type AppContext struct {
	Cfg           *config.Config
	ConfigManager *config.ConfigManager
	Logger        *zap.Logger
	DB            *gorm.DB
	Cache         cache.Cache
	JWT           *jwt.JWTManager
	Engine        *gin.Engine
	Uploader      upload.Uploader
	Scheduler     *scheduler.Scheduler
	Captcha       *captcha.Captcha
}

// Module 业务模块接口
//
// 说明：模块在 Register 中自行完成依赖装配与路由/任务注册。
type Module interface {
	Name() string
	Register(ctx *AppContext) error
}

const disabledRateLimitMaxIntValue = math.MaxInt32

// BuildContext 构建应用运行时上下文
//
// 说明：此处仅负责基础设施的创建与装配，业务依赖必须由模块自行处理。
func BuildContext(configPath string) (*AppContext, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	zapLogger, err := logger.InitLogger(&cfg.Log)
	if err != nil {
		return nil, err
	}

	cm, err := config.NewConfigManager(configPath, zapLogger)
	if err != nil {
		return nil, err
	}

	isDebug := cfg.Server.Mode == "debug"
	database, err := db.InitDB(&cfg.Database, zapLogger, isDebug)
	if err != nil {
		return nil, err
	}

	cacheInstance, err := cache.NewCache(&cfg.Cache)
	if err != nil {
		return nil, err
	}

	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	rateLimiter := buildRateLimiter(cm)
	engine := server.NewServer(&cfg.Server, rateLimiter, zapLogger)

	uploader, err := upload.NewUploader(upload.ConfigFromAppConfig(cfg))
	if err != nil {
		return nil, err
	}

	schedulerInstance := scheduler.NewScheduler(zapLogger)

	cap := captcha.NewCaptcha(cacheInstance)

	return &AppContext{
		Cfg:           cfg,
		ConfigManager: cm,
		Logger:        zapLogger,
		DB:            database,
		Cache:         cacheInstance,
		JWT:           jwtManager,
		Engine:        engine,
		Uploader:      uploader,
		Scheduler:     schedulerInstance,
		Captcha:       cap,
	}, nil
}

func buildRateLimiter(cm *config.ConfigManager) *middleware.RateLimiter {
	cfg := cm.GetConfig()
	if !cfg.RateLimit.Enabled {
		return middleware.NewRateLimiter(disabledRateLimitMaxIntValue, disabledRateLimitMaxIntValue)
	}
	return middleware.NewRateLimiter(cfg.RateLimit.RPS, cfg.RateLimit.Burst)
}

// RegisterModules 批量注册模块
func RegisterModules(ctx *AppContext, modules ...Module) error {
	for _, m := range modules {
		if err := m.Register(ctx); err != nil {
			return fmt.Errorf("register module %s failed: %w", m.Name(), err)
		}
	}
	return nil
}
