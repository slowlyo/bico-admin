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

// 注意：生成的 Provider 代码应该直接添加到 ProviderSet 中
// 不再使用动态注册模式，代码生成器会提供具体的插入位置指导

// ProviderSet Admin端Provider集合
// 注意：生成的 Provider 需要手动添加到这里，或者使用 wire.Build 时包含生成的 ProviderSet
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
	handler.NewCommonHandler,

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
	commonHandler *handler.CommonHandler,
) *routes.Handlers {
	return &routes.Handlers{
		AuthHandler:      authHandler,
		AdminUserHandler: adminUserHandler,
		AdminRoleHandler: adminRoleHandler,
		CommonHandler:    commonHandler,
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
