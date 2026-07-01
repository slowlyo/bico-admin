package admin

import (
	"bico-admin/internal/admin/handler"
	"bico-admin/internal/admin/middleware"
	"bico-admin/internal/admin/service"
	"bico-admin/internal/core/app"
	coreMiddleware "bico-admin/internal/core/middleware"
	"bico-admin/internal/pkg/crud"

	"gorm.io/gorm"
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
	userStatusMiddleware := middleware.NewUserStatusMiddleware(authSvc)

	authHandler := handler.NewAuthHandler(authSvc, ctx.Uploader, ctx.Captcha)
	uploadHandler := handler.NewUploadHandler(ctx.Uploader)
	commonHandler := handler.NewCommonHandler(cfgSvc)
	dashboardHandler := handler.NewDashboardHandler(ctx.Cfg, ctx.DB)

	modules := NewCRUDModules(ctx.DB, authSvc)
	r := NewRouter(authHandler, uploadHandler, commonHandler, dashboardHandler, jwtAuth, permMiddleware, userStatusMiddleware, ctx.DB, modules)
	r.Register(ctx.Engine)

	return nil
}

// NewCRUDModules 创建后台声明式 CRUD 模块列表。
//
// 说明：运行时路由注册和 Swagger 文档增强共用同一份模块配置，避免接口路径与文档漂移。
func NewCRUDModules(db *gorm.DB, cacheInvalidator service.AuthCacheInvalidator) []crud.Module {
	return []crud.Module{
		handler.NewAdminUserHandler(db, cacheInvalidator),
		handler.NewAdminRoleHandler(db, cacheInvalidator),
	}
}
