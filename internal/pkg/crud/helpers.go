package crud

import "fmt"

// CRUDPerms 标准 CRUD 权限集合
type CRUDPerms struct {
	Menu   string
	List   string
	Create string
	Edit   string
	Delete string
	Tree   []Permission
}

// NewCRUDPerms 生成标准 CRUD 权限
func NewCRUDPerms(module, label string) CRUDPerms {
	prefix := fmt.Sprintf("system:%s", module)
	p := CRUDPerms{
		Menu:   prefix + ":menu",
		List:   prefix + ":list",
		Create: prefix + ":create",
		Edit:   prefix + ":edit",
		Delete: prefix + ":delete",
	}
	p.Tree = []Permission{{
		Key:   p.Menu,
		Label: label,
		Children: []Permission{
			{Key: p.List, Label: "查看列表"},
			{Key: p.Create, Label: "创建"},
			{Key: p.Edit, Label: "编辑"},
			{Key: p.Delete, Label: "删除"},
		},
	}}
	return p
}

// WithExtra 添加额外权限
func (p CRUDPerms) WithExtra(perms ...Permission) CRUDPerms {
	if len(p.Tree) > 0 {
		p.Tree[0].Children = append(p.Tree[0].Children, perms...)
	}
	return p
}

// Routes 生成标准 CRUD 路由
func (p CRUDPerms) Routes() []Route {
	return []Route{
		{Method: "GET", Path: "", Handler: "List", Permission: p.List},
		{Method: "GET", Path: "/:id", Handler: "Get", Permission: p.List},
		{Method: "POST", Path: "", Handler: "Create", Permission: p.Create},
		{Method: "PUT", Path: "/:id", Handler: "Update", Permission: p.Edit},
		{Method: "DELETE", Path: "/:id", Handler: "Delete", Permission: p.Delete},
	}
}

// RoutesWithExtra 生成标准 CRUD 路由 + 额外路由
func (p CRUDPerms) RoutesWithExtra(extra ...Route) []Route {
	return append(p.Routes(), extra...)
}

// UniqueUints 数组去重
func UniqueUints(ids []uint) []uint {
	seen := make(map[uint]bool)
	result := make([]uint, 0, len(ids))
	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}
	return result
}

// CRUDPermissions 快速生成 CRUD 权限配置
// Deprecated: 使用 NewCRUDPerms 替代
// name: 模块名，如 "admin_user"
// label: 显示名称，如 "用户管理"
// parentKey: 父级权限 key，如 "system:manage"
func CRUDPermissions(name, label, parentKey string) (menuPerm, listPerm, createPerm, editPerm, deletePerm string, permissions []Permission) {
	prefix := name
	if parentKey != "" {
		// 从 parentKey 提取前缀，如 "system:manage" -> "system"
		for i := len(parentKey) - 1; i >= 0; i-- {
			if parentKey[i] == ':' {
				prefix = parentKey[:i] + ":" + name
				break
			}
		}
	}

	menuPerm = prefix + ":menu"
	listPerm = prefix + ":list"
	createPerm = prefix + ":create"
	editPerm = prefix + ":edit"
	deletePerm = prefix + ":delete"

	permissions = []Permission{
		{
			Key:   menuPerm,
			Label: label,
			Children: []Permission{
				{Key: listPerm, Label: "查看列表"},
				{Key: createPerm, Label: "创建"},
				{Key: editPerm, Label: "编辑"},
				{Key: deletePerm, Label: "删除"},
			},
		},
	}

	return
}

// CRUDRoutes 快速生成标准 CRUD 路由配置
// Deprecated: 使用 CRUDPerms.Routes() 替代
func CRUDRoutes(listPerm, createPerm, editPerm, deletePerm string) []Route {
	return []Route{
		{Method: "GET", Path: "", Handler: "List", Permission: listPerm},
		{Method: "GET", Path: "/:id", Handler: "Get", Permission: listPerm},
		{Method: "POST", Path: "", Handler: "Create", Permission: createPerm},
		{Method: "PUT", Path: "/:id", Handler: "Update", Permission: editPerm},
		{Method: "DELETE", Path: "/:id", Handler: "Delete", Permission: deletePerm},
	}
}

// CRUDRoutesWithExtra 生成标准 CRUD 路由 + 额外路由
// Deprecated: 使用 CRUDPerms.RoutesWithExtra() 替代
func CRUDRoutesWithExtra(listPerm, createPerm, editPerm, deletePerm string, extra ...Route) []Route {
	routes := CRUDRoutes(listPerm, createPerm, editPerm, deletePerm)
	return append(routes, extra...)
}

// PublicRoute 快速定义公开路由（不需要登录）
func PublicRoute(method, path, handler string) Route {
	return Route{Method: method, Path: path, Handler: handler, Public: true}
}

// AuthRoute 快速定义需要登录但不需要权限的路由
func AuthRoute(method, path, handler string) Route {
	return Route{Method: method, Path: path, Handler: handler}
}

// PermRoute 快速定义需要权限的路由
func PermRoute(method, path, handler, permission string) Route {
	return Route{Method: method, Path: path, Handler: handler, Permission: permission}
}
