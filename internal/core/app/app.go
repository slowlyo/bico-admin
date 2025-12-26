package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bico-admin/internal/core/cache"
	"bico-admin/internal/core/config"
	"bico-admin/internal/core/scheduler"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// App 应用结构体
type App struct {
	cfg       *config.Config
	engine    *gin.Engine
	server    *http.Server
	scheduler *scheduler.Scheduler
	db        *gorm.DB
	cache     cache.Cache
	logger    *zap.Logger
}

// NewApp 创建应用实例
func NewApp(
	cfg *config.Config,
	engine *gin.Engine,
	scheduler *scheduler.Scheduler,
	db *gorm.DB,
	cache cache.Cache,
	logger *zap.Logger,
) *App {
	return &App{
		cfg:       cfg,
		engine:    engine,
		scheduler: scheduler,
		db:        db,
		cache:     cache,
		logger:    logger,
	}
}

// Run 运行应用
func (a *App) Run() error {
	addr := fmt.Sprintf(":%d", a.cfg.Server.Port)

	a.server = &http.Server{
		Addr:              addr,
		Handler:           a.engine,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 启动定时任务调度器（任务由各模块自行注册）
	a.scheduler.Start()

	// 启动服务器
	go func() {
		a.logger.Info("服务启动成功", zap.String("addr", addr))
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("服务启动失败", zap.Error(err))
			os.Exit(1)
		}
	}()

	// 优雅关闭
	a.gracefulShutdown()

	return nil
}

// Run 使用 AppContext 运行应用
func Run(ctx *AppContext) error {
	application := NewApp(ctx.Cfg, ctx.Engine, ctx.Scheduler, ctx.DB, ctx.Cache, ctx.Logger)
	return application.Run()
}

// gracefulShutdown 优雅关闭
func (a *App) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("正在关闭服务")

	// 停止定时任务调度器
	a.scheduler.Stop()

	// 关闭缓存连接
	if err := a.cache.Close(); err != nil {
		a.logger.Error("关闭缓存失败", zap.Error(err))
	}

	// 关闭数据库连接池
	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err != nil {
			a.logger.Error("获取数据库连接池失败", zap.Error(err))
		} else {
			if err := sqlDB.Close(); err != nil {
				a.logger.Error("关闭数据库连接池失败", zap.Error(err))
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("服务关闭异常", zap.Error(err))
	}

	// 同步日志
	if err := a.logger.Sync(); err != nil {
		// 忽略 sync 错误（stdout/stderr 在某些系统上会报错）
	}

	a.logger.Info("服务已关闭")
}
