package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"

	"bico-admin/internal/admin/handler"
	"bico-admin/internal/admin/initializer"
	"bico-admin/internal/admin/middleware"
	"bico-admin/internal/admin/repository"
	"bico-admin/internal/admin/routes"
	"bico-admin/internal/admin/service"
)

// ProviderSet Admin端Provider集合
var ProviderSet = wire.NewSet(
	// Repository层
	repository.NewAdminUserRepository,
	repository.NewAdminRoleRepository,

	// Service层
	service.NewAdminUserService,
	service.NewAdminRoleService,
	service.NewAuthService,

	// Handler层
	handler.NewAuthHandler,
	handler.NewAdminUserHandler,
	handler.NewAdminRoleHandler,

	// 路由处理器集合
	ProvideHandlers,

	// 权限中间件
	ProvidePermissionMiddleware,
)

// ProvideHandlers 提供处理器集合
func ProvideHandlers(
	authHandler *handler.AuthHandler,
	adminUserHandler *handler.AdminUserHandler,
	adminRoleHandler *handler.AdminRoleHandler,
) *routes.Handlers {
	return &routes.Handlers{
		AuthHandler:      authHandler,
		AdminUserHandler: adminUserHandler,
		AdminRoleHandler: adminRoleHandler,
	}
}

// ProvidePermissionMiddleware 提供权限中间件
func ProvidePermissionMiddleware(adminUserService service.AdminUserService) gin.HandlerFunc {
	return middleware.PermissionMiddlewareFactory(adminUserService)
}

// AutoMigrateAdminModels 自动迁移Admin模块的数据库表
func AutoMigrateAdminModels(db *gorm.DB) error {
	dbInitializer := initializer.NewDatabaseInitializer(db)
	return dbInitializer.MigrateAndSeed()
}
