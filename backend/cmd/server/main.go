package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"bico-admin/core/config"
	"bico-admin/core/middleware"
	coreRouter "bico-admin/core/router"
	adminRouter "bico-admin/modules/admin/router"
	apiRouter "bico-admin/modules/api/router"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 初始化配置
	cfg := config.New()

	// 初始化数据库
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 初始化Redis
	rdb, err := config.InitRedis(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Println("Cache functionality will be disabled")
	} else {
		config.SetRedis(rdb)
	}

	// 创建Fiber应用
	app := fiber.New(fiber.Config{
		AppName: "Bico Admin v1.0.0",
	})

	// 全局中间件
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// 自定义中间件
	app.Use(middleware.ErrorHandler())

	// 静态文件服务
	app.Static("/uploads", "./storage/uploads")

	// 路由组
	api := app.Group("/api")
	admin := app.Group("/admin")

	// 注册路由
	coreRouter.SetupSystemRoutes(app, db) // 核心系统路由（健康检查等）
	apiRouter.SetupRoutes(api, db, cfg)   // 对外API路由（包含认证）
	adminRouter.SetupRoutes(admin, db)    // 后台管理路由（包含认证）

	// 健康检查
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Bico Admin is running",
		})
	})

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
