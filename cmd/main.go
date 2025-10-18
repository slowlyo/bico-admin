package main

import (
	"fmt"
	"os"

	"bico-admin/internal/admin"
	"bico-admin/internal/api"
	"bico-admin/internal/core/app"
	"bico-admin/internal/core/config"
	"bico-admin/internal/core/server"
	"bico-admin/internal/migrate"
	"bico-admin/web"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var (
	configPath string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "bico-admin",
	Short: "Bico Admin 管理系统",
	Long:  "基于 Gin + GORM + Viper + Dig 构建的管理系统",
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 HTTP 服务",
	Long:  "启动 Web 服务器，监听 HTTP 请求",
	Run: func(cmd *cobra.Command, args []string) {
		container, err := app.BuildContainer(configPath)
		if err != nil {
			fmt.Printf("❌ 构建容器失败: %v\n", err)
			os.Exit(1)
		}

		if err := container.Invoke(func(
			engine *gin.Engine,
			adminRouter *admin.Router,
			apiRouter *api.Router,
			cfg *config.Config,
			application *app.App,
		) error {
			server.RegisterRoutes(engine, adminRouter, apiRouter, cfg, web.DistFS)
			return application.Run()
		}); err != nil {
			fmt.Printf("❌ 启动失败: %v\n", err)
			os.Exit(1)
		}
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "执行数据库迁移",
	Long:  "自动创建或更新数据库表结构",
	Run: func(cmd *cobra.Command, args []string) {
		container, err := app.BuildContainer(configPath)
		if err != nil {
			fmt.Printf("❌ 构建容器失败: %v\n", err)
			os.Exit(1)
		}

		if err := container.Invoke(func(db *gorm.DB) error {
			fmt.Println("📦 开始数据库迁移...")
			if err := migrate.AutoMigrate(db); err != nil {
				return fmt.Errorf("迁移失败: %w", err)
			}
			fmt.Println("✅ 数据库迁移完成")
			return nil
		}); err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config/config.yaml", "配置文件路径")
	
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(migrateCmd)
}
