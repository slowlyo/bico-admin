package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bico-admin/internal/core/config"
	"github.com/gin-gonic/gin"
)

// App 应用结构体
type App struct {
	cfg    *config.Config
	engine *gin.Engine
	server *http.Server
}

// NewApp 创建应用实例
func NewApp(cfg *config.Config, engine *gin.Engine) *App {
	return &App{
		cfg:    cfg,
		engine: engine,
	}
}

// Run 运行应用
func (a *App) Run() error {
	addr := fmt.Sprintf(":%d", a.cfg.Server.Port)
	
	a.server = &http.Server{
		Addr:    addr,
		Handler: a.engine,
	}

	// 启动服务器
	go func() {
		fmt.Printf("🚀 服务启动成功，监听端口: %s\n", addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("❌ 服务启动失败: %v\n", err)
			os.Exit(1)
		}
	}()

	// 优雅关闭
	a.gracefulShutdown()
	
	return nil
}

// gracefulShutdown 优雅关闭
func (a *App) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("🛑 正在关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		fmt.Printf("❌ 服务关闭异常: %v\n", err)
	}

	fmt.Println("👋 服务已关闭")
}
