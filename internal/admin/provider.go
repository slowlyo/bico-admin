package admin

import (
	"github.com/google/wire"

	"bico-admin/internal/admin/handler"
	"bico-admin/internal/admin/repository"
	"bico-admin/internal/admin/routes"
	"bico-admin/internal/admin/service"
)

// ProviderSet Admin端Provider集合
var ProviderSet = wire.NewSet(
	// Repository层
	repository.NewUserRepository,
	repository.NewAdminUserRepository,

	// Service层
	service.NewUserService,
	service.NewAdminUserService,
	service.NewAuthService,

	// Handler层
	handler.NewAuthHandler,
	handler.NewUserHandler,
	handler.NewSystemHandler,

	// 路由处理器集合
	ProvideHandlers,
)

// ProvideHandlers 提供处理器集合
func ProvideHandlers(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	systemHandler *handler.SystemHandler,
) *routes.Handlers {
	return &routes.Handlers{
		AuthHandler:   authHandler,
		UserHandler:   userHandler,
		SystemHandler: systemHandler,
	}
}
