package main

import (
	"os"

	_ "bico-admin/docs"
	"bico-admin/internal/admin"
	"bico-admin/internal/api"
	"bico-admin/internal/core/app"
	"bico-admin/internal/core/logger"
	"bico-admin/internal/core/server"
	"bico-admin/internal/job"
	"bico-admin/internal/migrate"
	"bico-admin/web"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// @title Bico Admin API
// @version 1.0
// @description 基于 Gin + GORM 构建的管理系统 API 文档
// @termsOfService https://github.com/slowlyo/bico-admin

// @contact.name API Support
// @contact.url https://github.com/slowlyo/bico-admin

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT 认证，格式: Bearer {token}

var (
	configPath string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("命令执行失败", zap.Error(err))
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "bico-admin",
	Short: "Bico Admin 管理系统",
	Long:  "基于 Gin + GORM + Viper 构建的管理系统",
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 HTTP 服务",
	Long:  "启动 Web 服务器，监听 HTTP 请求",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, err := app.BuildContext(configPath)
		if err != nil {
			logger.Error("构建上下文失败", zap.Error(err))
			os.Exit(1)
		}

		server.RegisterCoreRoutes(ctx.Engine, ctx.Cfg, web.DistFS)

		if err := app.RegisterModules(
			ctx,
			admin.NewModule(),
			api.NewModule(),
			job.NewModule(),
		); err != nil {
			ctx.Logger.Error("注册模块失败", zap.Error(err))
			os.Exit(1)
		}

		if err := app.Run(ctx); err != nil {
			ctx.Logger.Error("启动失败", zap.Error(err))
			os.Exit(1)
		}
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "执行数据库迁移",
	Long:  "自动创建或更新数据库表结构",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, err := app.BuildContext(configPath)
		if err != nil {
			logger.Error("构建上下文失败", zap.Error(err))
			os.Exit(1)
		}

		ctx.Logger.Info("开始数据库迁移")
		if err := migrate.AutoMigrate(ctx.DB); err != nil {
			ctx.Logger.Error("数据库迁移失败", zap.Error(err))
			os.Exit(1)
		}
		ctx.Logger.Info("数据库迁移完成")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "配置文件路径（默认自动查找 config.yaml 或 config/config.yaml）")

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(migrateCmd)
}
