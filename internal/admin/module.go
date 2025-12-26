package admin

import (
	"bico-admin/internal/admin/handler"
	"bico-admin/internal/admin/middleware"
	"bico-admin/internal/admin/service"
	"bico-admin/internal/core/app"
	coreMiddleware "bico-admin/internal/core/middleware"
	"bico-admin/internal/pkg/crud"
)

// Module admin 模块
type Module struct{}

// NewModule 创建 admin 模块
func NewModule() *Module {
	return &Module{}
}

// Name 模块名称
func (m *Module) Name() string {
	return "admin"
}

// Register 注册 admin 模块
func (m *Module) Register(ctx *app.AppContext) error {
	authSvc := service.NewAuthService(ctx.DB, ctx.JWT, ctx.Cache)
	cfgSvc := service.NewConfigService(ctx.ConfigManager)

	jwtAuth := coreMiddleware.JWTAuth(ctx.JWT, authSvc)
	permMiddleware := middleware.NewPermissionMiddleware(authSvc)
	userStatusMiddleware := middleware.NewUserStatusMiddleware(ctx.DB)

	authHandler := handler.NewAuthHandler(authSvc, ctx.Uploader, ctx.Captcha)
	uploadHandler := handler.NewUploadHandler(ctx.Uploader)
	commonHandler := handler.NewCommonHandler(cfgSvc)

	modules := []crud.Module{
		handler.NewAdminUserHandler(ctx.DB),
		handler.NewAdminRoleHandler(ctx.DB),
	}

	r := NewRouter(authHandler, uploadHandler, commonHandler, jwtAuth, permMiddleware, userStatusMiddleware, ctx.DB, modules)
	r.Register(ctx.Engine)

	return nil
}
