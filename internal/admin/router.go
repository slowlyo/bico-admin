package admin

import (
	"reflect"

	"bico-admin/internal/admin/handler"
	"bico-admin/internal/admin/middleware"
	"bico-admin/internal/pkg/crud"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Router 实现路由注册
type Router struct {
	authHandler          *handler.AuthHandler
	commonHandler        *handler.CommonHandler
	jwtAuth              gin.HandlerFunc
	permMiddleware       *middleware.PermissionMiddleware
	userStatusMiddleware *middleware.UserStatusMiddleware
	db                   *gorm.DB
}

// NewRouter 创建路由实例
func NewRouter(
	authHandler *handler.AuthHandler,
	commonHandler *handler.CommonHandler,
	jwtAuth gin.HandlerFunc,
	permMiddleware *middleware.PermissionMiddleware,
	userStatusMiddleware *middleware.UserStatusMiddleware,
	db *gorm.DB,
) *Router {
	// 初始化基础权限树
	initBasePermissions()

	return &Router{
		authHandler:          authHandler,
		commonHandler:        commonHandler,
		jwtAuth:              jwtAuth,
		permMiddleware:       permMiddleware,
		userStatusMiddleware: userStatusMiddleware,
		db:                   db,
	}
}

// initBasePermissions 初始化基础权限树
func initBasePermissions() {
	crud.SetBasePermissions([]crud.Permission{
		{Key: handler.PermDashboardMenu, Label: "工作台"},
		{
			Key:      handler.PermSystemManage,
			Label:    "系统管理",
			Children: []crud.Permission{},
		},
	})
}

// Register 注册路由
func (r *Router) Register(engine *gin.Engine) {
	admin := engine.Group("/admin-api")

	// 公开路由
	{
		admin.POST("/auth/login", r.authHandler.Login)
		admin.GET("/captcha", r.authHandler.GetCaptcha)
		admin.GET("/app-config", r.commonHandler.GetAppConfig)
	}

	// 需要认证的路由
	authorized := admin.Group("", r.jwtAuth, r.userStatusMiddleware.Check())
	{
		// 认证相关（特殊路由，不走 CRUD 模块）
		auth := authorized.Group("/auth")
		{
			auth.POST("/logout", r.authHandler.Logout)
			auth.GET("/current-user", r.authHandler.CurrentUser)
			auth.PUT("/profile", r.authHandler.UpdateProfile)
			auth.PUT("/password", r.authHandler.ChangePassword)
			auth.POST("/avatar", r.authHandler.UploadAvatar)
		}

		// 自动注册所有 CRUD 模块路由
		r.registerModules(authorized)
	}
}

// registerModules 注册所有通过 crud.RegisterModule 注册的模块
func (r *Router) registerModules(group *gin.RouterGroup) {
	moduleRouter := crud.NewModuleRouter(
		r.jwtAuth,
		r.permMiddleware,
		r.userStatusMiddleware.Check(),
	)

	for _, reg := range crud.GetRegisteredModules() {
		module := r.instantiateModule(reg.Constructor)
		if module != nil {
			moduleRouter.RegisterModule(group, module)
		}
	}
}

// instantiateModule 实例化模块
func (r *Router) instantiateModule(constructor interface{}) crud.Module {
	constructorVal := reflect.ValueOf(constructor)
	constructorType := constructorVal.Type()

	numIn := constructorType.NumIn()
	args := make([]reflect.Value, numIn)

	for i := 0; i < numIn; i++ {
		argType := constructorType.In(i)
		switch argType.String() {
		case "*gorm.DB":
			args[i] = reflect.ValueOf(r.db)
		default:
			return nil
		}
	}

	results := constructorVal.Call(args)
	if len(results) == 0 {
		return nil
	}

	if module, ok := results[0].Interface().(crud.Module); ok {
		return module
	}

	return nil
}
