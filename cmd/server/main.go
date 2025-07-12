package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "bico-admin/docs" // 导入swagger文档
	"bico-admin/pkg/config"
	"bico-admin/pkg/logger"

	"go.uber.org/zap"
)

// @title Bico Admin API
// @version 1.0
// @description Bico Admin 管理系统API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Bearer token

func main() {
	// 解析命令行参数
	var configPath = flag.String("config", "", "配置文件路径")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 初始化日志
	logConfig := logger.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		Output:     cfg.Log.Output,
		Filename:   cfg.Log.Filename,
		MaxSize:    cfg.Log.MaxSize,
		MaxAge:     cfg.Log.MaxAge,
		MaxBackups: cfg.Log.MaxBackups,
		Compress:   cfg.Log.Compress,
	}
	if err := logger.Init(logConfig); err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}
	defer logger.Sync()

	logger.Info("启动应用",
		zap.String("app_name", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.String("environment", cfg.App.Environment),
	)

	// 初始化应用
	app, err := initializeApp(cfg)
	if err != nil {
		logger.Fatal("初始化应用失败", zap.Error(err))
	}

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      app,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 启动服务器
	go func() {
		logger.Info("启动HTTP服务器",
			zap.String("addr", server.Addr),
		)

		// 输出访问地址信息
		host := cfg.Server.Host
		if host == "0.0.0.0" {
			host = "localhost"
		}
		baseURL := fmt.Sprintf("http://%s:%d", host, cfg.Server.Port)

		fmt.Printf("\n🚀 服务启动成功！\n")
		fmt.Printf("📍 访问地址:\n")
		fmt.Printf("   • 主页: %s\n", baseURL)
		fmt.Printf("   • Admin端: %s/admin\n", baseURL)
		fmt.Printf("   • Master端: %s/master\n", baseURL)
		fmt.Printf("   • API端: %s/api\n", baseURL)
		fmt.Printf("📚 API文档:\n")
		fmt.Printf("   • Swagger UI: %s/swagger/index.html\n", baseURL)
		fmt.Printf("   • API JSON: %s/swagger/doc.json\n", baseURL)
		fmt.Printf("🔧 环境: %s\n", cfg.App.Environment)
		fmt.Printf("📝 日志级别: %s\n\n", cfg.Log.Level)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("启动HTTP服务器失败", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭失败", zap.Error(err))
	} else {
		logger.Info("服务器已关闭")
	}
}
