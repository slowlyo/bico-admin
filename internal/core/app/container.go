package app

import (
	"bico-admin/internal/core/config"
	"bico-admin/internal/core/db"
	"bico-admin/internal/core/server"
	"bico-admin/internal/shared/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

// BuildContainer 构建 DI 容器
func BuildContainer(configPath string) (*dig.Container, error) {
	container := dig.New()

	// 加载配置
	if err := container.Provide(func() (*config.Config, error) {
		return config.LoadConfig(configPath)
	}); err != nil {
		return nil, err
	}

	// 提供数据库连接
	if err := container.Provide(func(cfg *config.Config) (*gorm.DB, error) {
		return db.InitDB(&cfg.Database)
	}); err != nil {
		return nil, err
	}

	// 提供 Gin 引擎
	if err := container.Provide(func(cfg *config.Config) *gin.Engine {
		return server.NewServer(&cfg.Server)
	}); err != nil {
		return nil, err
	}

	// 提供路由
	if err := container.Provide(NewAdminRouter); err != nil {
		return nil, err
	}

	if err := container.Provide(NewAPIRouter); err != nil {
		return nil, err
	}

	// 提供应用实例
	if err := container.Provide(NewApp); err != nil {
		return nil, err
	}

	return container, nil
}

func NewAdminRouter() server.Router {
	return &adminRouter{}
}

func NewAPIRouter() server.Router {
	return &apiRouter{}
}

type adminRouter struct{}

func (r *adminRouter) Register(engine *gin.Engine) {
	admin := engine.Group("/admin")
	{
		admin.GET("/ping", func(c *gin.Context) {
			c.JSON(200, response.Success("admin pong"))
		})
	}
}

type apiRouter struct{}

func (r *apiRouter) Register(engine *gin.Engine) {
	api := engine.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, response.Success("api pong"))
		})
	}
}
