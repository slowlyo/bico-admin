//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"gorm.io/gorm"

	"bico-admin/internal/admin"
	adminRoutes "bico-admin/internal/admin/routes"
	"bico-admin/internal/api"
	apiRoutes "bico-admin/internal/api/routes"
	"bico-admin/internal/master"
	masterRoutes "bico-admin/internal/master/routes"
	"bico-admin/internal/shared"
	sharedMiddleware "bico-admin/internal/shared/middleware"
	"bico-admin/pkg/cache"
	"bico-admin/pkg/config"
	"bico-admin/pkg/logger"
)

// initializeApp 初始化应用
func initializeApp(cfg *config.Config) (*gin.Engine, error) {
	wire.Build(
		// 共享组件
		shared.ProviderSet,

		// Admin端组件
		admin.ProviderSet,

		// Master端组件
		master.ProviderSet,

		// API端组件
		api.ProviderSet,

		// 提供Gin引擎
		ProvideGinEngine,
	)
	return &gin.Engine{}, nil
}

// ProvideGinEngine 提供Gin引擎
func ProvideGinEngine(
	cfg *config.Config,
	db *gorm.DB,
	cache cache.Cache,
	adminHandlers *adminRoutes.Handlers,
	adminPermissionMiddleware gin.HandlerFunc,
	masterHandlers *masterRoutes.Handlers,
	apiHandlers *apiRoutes.Handlers,
) *gin.Engine {
	// 执行Admin模块数据库迁移
	if err := admin.AutoMigrateAdminModels(db); err != nil {
		logger.Error("Admin模块数据库迁移失败")
		panic(err)
	}

	// 设置Gin模式
	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建Gin引擎
	r := gin.New()

	// 添加全局中间件
	r.Use(gin.Recovery())
	r.Use(sharedMiddleware.CORS()) // 全局CORS中间件
	r.Use(sharedMiddleware.Logging())

	// 注册路由
	adminRoutes.RegisterRoutes(r, adminHandlers, cache, adminPermissionMiddleware)
	masterRoutes.RegisterRoutes(r, masterHandlers)
	apiRoutes.RegisterRoutes(r, apiHandlers)

	return r
}
