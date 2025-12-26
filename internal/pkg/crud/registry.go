package crud

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gin-gonic/gin"
)

// Permission 权限定义
type Permission struct {
	Key      string       `json:"key"`
	Label    string       `json:"label"`
	Children []Permission `json:"children,omitempty"`
}

// Route 路由定义
type Route struct {
	Method     string // GET, POST, PUT, DELETE
	Path       string // 路由路径，如 "", "/:id", "/all"
	Handler    string // Handler 方法名
	Permission string // 权限 key，为空则不校验
	Public     bool   // 是否公开（不需要登录）
}

// ModuleConfig 模块配置（在 Handler 中定义）
type ModuleConfig struct {
	// 基础信息
	Name        string // 模块名称，如 "admin_user"
	Group       string // 路由分组，如 "/admin-users"
	Description string // 描述

	// 权限配置
	ParentPermission string       // 父级权限 key
	Permissions      []Permission // 权限树

	// 路由配置
	Routes []Route
}

// Module Handler 模块接口，实现此接口的 Handler 会被自动注册
type Module interface {
	// ModuleConfig 返回模块配置
	ModuleConfig() ModuleConfig
}

var (
	allPermissions []Permission
	permissionsMu  sync.RWMutex
)

// AddPermissions 添加权限到全局权限树
func AddPermissions(parentKey string, perms []Permission) {
	permissionsMu.Lock()
	defer permissionsMu.Unlock()

	if parentKey == "" {
		allPermissions = append(allPermissions, perms...)
		return
	}

	// 在现有权限树中找到父节点并添加
	addToParent(&allPermissions, parentKey, perms)
}

func addToParent(tree *[]Permission, parentKey string, children []Permission) bool {
	for i := range *tree {
		if (*tree)[i].Key == parentKey {
			(*tree)[i].Children = append((*tree)[i].Children, children...)
			return true
		}
		if len((*tree)[i].Children) > 0 {
			if addToParent(&(*tree)[i].Children, parentKey, children) {
				return true
			}
		}
	}
	return false
}

// SetBasePermissions 设置基础权限树
func SetBasePermissions(perms []Permission) {
	permissionsMu.Lock()
	defer permissionsMu.Unlock()
	allPermissions = perms
}

// GetAllPermissions 获取完整的权限树
func GetAllPermissions() []Permission {
	permissionsMu.RLock()
	defer permissionsMu.RUnlock()
	return allPermissions
}

// GetAllPermissionKeys 获取所有权限的 key 列表
func GetAllPermissionKeys() []string {
	permissionsMu.RLock()
	defer permissionsMu.RUnlock()

	var keys []string
	var collectKeys func(perms []Permission)

	collectKeys = func(perms []Permission) {
		for _, perm := range perms {
			keys = append(keys, perm.Key)
			if len(perm.Children) > 0 {
				collectKeys(perm.Children)
			}
		}
	}

	collectKeys(allPermissions)
	return keys
}

// PermissionChecker 权限检查接口
type PermissionChecker interface {
	RequirePermission(permission string) gin.HandlerFunc
}

// RouterConfig 路由配置
type RouterConfig struct {
	// 认证中间件（可选，为 nil 时所有路由都是公开的）
	AuthMiddleware gin.HandlerFunc
	// 用户状态检查（可选）
	UserStatusMiddleware gin.HandlerFunc
	// 权限中间件（可选）
	PermMiddleware PermissionChecker
}

// ModuleRouter 模块路由注册器
type ModuleRouter struct {
	config RouterConfig
}

// NewModuleRouter 创建模块路由注册器
func NewModuleRouter(
	jwtAuth gin.HandlerFunc,
	permMiddleware PermissionChecker,
	userStatusMiddleware gin.HandlerFunc,
) *ModuleRouter {
	return &ModuleRouter{
		config: RouterConfig{
			AuthMiddleware:       jwtAuth,
			PermMiddleware:       permMiddleware,
			UserStatusMiddleware: userStatusMiddleware,
		},
	}
}

// NewModuleRouterWithConfig 使用配置创建模块路由注册器
func NewModuleRouterWithConfig(config RouterConfig) *ModuleRouter {
	return &ModuleRouter{config: config}
}

// RegisterModule 注册单个模块的路由
func (r *ModuleRouter) RegisterModule(engine *gin.RouterGroup, module Module) {
	config := module.ModuleConfig()

	// 注册权限
	if len(config.Permissions) > 0 {
		AddPermissions(config.ParentPermission, config.Permissions)
	}

	// 获取 handler 的反射值
	handlerVal := reflect.ValueOf(module)

	// 创建路由分组
	group := engine.Group(config.Group)

	// 是否需要认证
	hasPublic := false
	hasPrivate := false
	for _, route := range config.Routes {
		if route.Public {
			hasPublic = true
		} else {
			hasPrivate = true
		}
	}

	// 注册公开路由
	if hasPublic {
		for _, route := range config.Routes {
			if !route.Public {
				continue
			}
			r.registerRoute(group, handlerVal, route)
		}
	}

	// 注册需要认证的路由
	if hasPrivate {
		// 构建认证中间件链
		var authMiddlewares []gin.HandlerFunc
		if r.config.AuthMiddleware != nil {
			authMiddlewares = append(authMiddlewares, r.config.AuthMiddleware)
		}
		if r.config.UserStatusMiddleware != nil {
			authMiddlewares = append(authMiddlewares, r.config.UserStatusMiddleware)
		}

		authGroup := group.Group("", authMiddlewares...)
		for _, route := range config.Routes {
			if route.Public {
				continue
			}
			r.registerRoute(authGroup, handlerVal, route)
		}
	}
}

func (r *ModuleRouter) registerRoute(group *gin.RouterGroup, handlerVal reflect.Value, route Route) {
	// 获取 handler 方法
	method := handlerVal.MethodByName(route.Handler)
	if !method.IsValid() {
		panic(fmt.Sprintf("handler method %s not found", route.Handler))
	}

	// 将反射方法转换为 gin.HandlerFunc
	handlerFunc := func(c *gin.Context) {
		method.Call([]reflect.Value{reflect.ValueOf(c)})
	}

	// 构建中间件链
	handlers := make([]gin.HandlerFunc, 0)
	if route.Permission != "" && r.config.PermMiddleware != nil {
		handlers = append(handlers, r.config.PermMiddleware.RequirePermission(route.Permission))
	}
	handlers = append(handlers, handlerFunc)

	// 注册路由
	switch route.Method {
	case "GET":
		group.GET(route.Path, handlers...)
	case "POST":
		group.POST(route.Path, handlers...)
	case "PUT":
		group.PUT(route.Path, handlers...)
	case "DELETE":
		group.DELETE(route.Path, handlers...)
	case "PATCH":
		group.PATCH(route.Path, handlers...)
	default:
		panic(fmt.Sprintf("unsupported HTTP method: %s", route.Method))
	}
}
