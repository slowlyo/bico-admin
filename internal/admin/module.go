package admin

import (
	"fmt"

	"bico-admin/internal/admin/handler"
	"bico-admin/internal/admin/middleware"
	"bico-admin/internal/admin/service"
	"bico-admin/internal/core/app"
	coreMiddleware "bico-admin/internal/core/middleware"
	"bico-admin/internal/pkg/crud"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
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
	container := dig.New()

	if err := container.Provide(func() *gorm.DB { return ctx.DB }); err != nil {
		return fmt.Errorf("provide db failed: %w", err)
	}
	if err := container.Provide(func() service.IAuthService {
		return service.NewAuthService(ctx.DB, ctx.JWT, ctx.Cache)
	}); err != nil {
		return fmt.Errorf("provide auth service failed: %w", err)
	}
	if err := container.Provide(func() service.IConfigService {
		return service.NewConfigService(ctx.ConfigManager)
	}); err != nil {
		return fmt.Errorf("provide config service failed: %w", err)
	}
	if err := container.Provide(func(authSvc service.IAuthService) gin.HandlerFunc {
		return coreMiddleware.JWTAuth(ctx.JWT, authSvc)
	}); err != nil {
		return fmt.Errorf("provide jwt auth middleware failed: %w", err)
	}
	if err := container.Provide(func(authSvc service.IAuthService) *middleware.PermissionMiddleware {
		return middleware.NewPermissionMiddleware(authSvc)
	}); err != nil {
		return fmt.Errorf("provide permission middleware failed: %w", err)
	}
	if err := container.Provide(func() *middleware.UserStatusMiddleware {
		return middleware.NewUserStatusMiddleware(ctx.DB)
	}); err != nil {
		return fmt.Errorf("provide user status middleware failed: %w", err)
	}
	if err := container.Provide(func(authSvc service.IAuthService) *handler.AuthHandler {
		return handler.NewAuthHandler(authSvc, ctx.Uploader, ctx.Captcha)
	}); err != nil {
		return fmt.Errorf("provide auth handler failed: %w", err)
	}
	if err := container.Provide(func(cfgSvc service.IConfigService) *handler.CommonHandler {
		return handler.NewCommonHandler(cfgSvc)
	}); err != nil {
		return fmt.Errorf("provide common handler failed: %w", err)
	}

	if err := container.Provide(func(db *gorm.DB) []crud.Module {
		return []crud.Module{
			handler.NewAdminUserHandler(db),
			handler.NewAdminRoleHandler(db),
		}
	}); err != nil {
		return fmt.Errorf("provide crud modules failed: %w", err)
	}

	if err := container.Provide(func(
		authHandler *handler.AuthHandler,
		commonHandler *handler.CommonHandler,
		jwtAuth gin.HandlerFunc,
		permMiddleware *middleware.PermissionMiddleware,
		userStatusMiddleware *middleware.UserStatusMiddleware,
		db *gorm.DB,
		modules []crud.Module,
	) *Router {
		return NewRouter(authHandler, commonHandler, jwtAuth, permMiddleware, userStatusMiddleware, db, modules)
	}); err != nil {
		return fmt.Errorf("provide admin router failed: %w", err)
	}

	if err := container.Invoke(func(r *Router) {
		r.Register(ctx.Engine)
	}); err != nil {
		return fmt.Errorf("invoke admin router register failed: %w", err)
	}

	return nil
}
