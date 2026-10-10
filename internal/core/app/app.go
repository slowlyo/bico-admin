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
	serverErrCh := make(chan error, 1)
	go func() {
		a.logger.Info("服务启动成功", zap.String("addr", addr))
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case err := <-serverErrCh:
		a.logger.Error("服务启动失败", zap.Error(err))
		_ = a.shutdown()
		return err
	case <-quit:
		a.logger.Info("收到退出信号，开始优雅关闭")
	}

	if err := a.shutdown(); err != nil {
		return err
	}

	return nil
}

// Run 使用 AppContext 运行应用
func Run(ctx *AppContext) error {
	application := NewApp(ctx.Cfg, ctx.Engine, ctx.Scheduler, ctx.DB, ctx.Cache, ctx.Logger)
	return application.Run()
}

// shutdown 先排空 HTTP，再停任务，最后释放缓存和数据库。
//
// 说明：进行中的请求和定时任务都还要访问缓存与数据库，先关连接会让排空阶段直接失败。
func (a *App) shutdown() error {
	a.logger.Info("正在关闭服务")

	var shutdownErr error
	if a.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := a.server.Shutdown(ctx); err != nil {
			a.logger.Error("服务关闭异常", zap.Error(err))
			shutdownErr = err
		}
	}

	// HTTP 已停止接新请求后再停任务，避免任务与请求同时被掐断。
	if a.scheduler != nil {
		a.scheduler.Stop()
	}

	if a.cache != nil {
		if err := a.cache.Close(); err != nil {
			a.logger.Error("关闭缓存失败", zap.Error(err))
		}
	}

	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err != nil {
			a.logger.Error("获取数据库连接池失败", zap.Error(err))
		} else if err := sqlDB.Close(); err != nil {
			a.logger.Error("关闭数据库连接池失败", zap.Error(err))
		}
	}

	a.logger.Info("服务已关闭")
	// stdout 上的 Sync 在部分系统会返回错误，不影响进程退出。
	_ = a.logger.Sync()
	return shutdownErr
}
