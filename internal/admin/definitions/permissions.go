package definitions

import (
	"strings"
)

// Permission 权限定义（静态配置，程序中定义）
type Permission struct {
	Code      string   `json:"code"`       // 权限代码，如 "user:list"
	Name      string   `json:"name"`       // 权限名称，如 "查看用户列表"
	Module    string   `json:"module"`     // 所属模块，如 "user", "system"
	MenuSigns []string `json:"menu_signs"` // 关联的菜单标识列表
	Buttons   []string `json:"buttons"`    // 关联的按钮标识列表，如 ["create", "edit", "delete"]
	APIs      []string `json:"apis"`       // 关联的API路径列表，如 ["/api/admin/users", "/api/admin/users/:id"]
	Level     int      `json:"level"`      // 权限级别：1-查看，2-操作，3-管理，4-超级
}

// PermissionGroup 权限组（按模块分组）
type PermissionGroup struct {
	Module      string       `json:"module"`      // 模块名称
	Name        string       `json:"name"`        // 模块显示名称
	Permissions []Permission `json:"permissions"` // 该模块下的权限列表
}

// 权限级别常量
const (
	PermissionLevelView   = 1 // 查看权限
	PermissionLevelAction = 2 // 操作权限
	PermissionLevelManage = 3 // 管理权限
	PermissionLevelSuper  = 4 // 超级权限
)

// PermissionDef 权限定义的简化结构（用于生成完整权限）
type PermissionDef struct {
	Code      string
	Name      string
	Level     int
	MenuSigns string // 逗号分隔
	Buttons   string // 逗号分隔
	APIs      string // 逗号分隔
}

// ModuleDef 模块定义
type ModuleDef struct {
	Module      string
	Name        string
	Permissions []PermissionDef
}

// 权限定义数据（大幅简化的格式）
var permissionData = []ModuleDef{
	{"user", "用户管理", []PermissionDef{
		{"user:list", "查看用户列表", 1, "user,user.list", "search,filter,export", "/api/admin/users,/api/admin/users/stats"},
		{"user:create", "创建用户", 2, "user,user.list", "create", "/api/admin/users"},
		{"user:update", "编辑用户", 2, "user,user.list,user.detail", "edit,save", "/api/admin/users/:id"},
		{"user:delete", "删除用户", 3, "user,user.list", "delete,batch_delete", "/api/admin/users/:id"},
		{"user:export", "导出用户", 2, "user,user.list", "export", "/api/admin/users/export"},
	}},
	{"admin_user", "管理员管理", []PermissionDef{
		{"admin_user:list", "查看管理员列表", 1, "admin_user,admin_user.list", "search,filter", "/api/admin/admin-users"},
		{"admin_user:create", "创建管理员", 3, "admin_user,admin_user.create", "create", "/api/admin/admin-users"},
		{"admin_user:update", "编辑管理员", 3, "admin_user,admin_user.list", "edit,save", "/api/admin/admin-users/:id"},
		{"admin_user:delete", "删除管理员", 4, "admin_user,admin_user.list", "delete", "/api/admin/admin-users/:id"},
		{"admin_user:reset_password", "重置管理员密码", 3, "admin_user,admin_user.list", "reset_password", "/api/admin/admin-users/:id/password"},
	}},
	{"system", "系统管理", []PermissionDef{
		{"system:info", "查看系统信息", 1, "system,system.info", "refresh", "/api/admin/system/info,/api/admin/system/stats"},
		{"system:config", "系统配置", 4, "system,system.config", "save,reset", "/api/admin/system/config"},
	}},
	{"config", "配置管理", []PermissionDef{
		{"config:list", "查看配置列表", 1, "config,config.app", "search,filter", "/api/admin/config"},
		{"config:update", "修改配置", 3, "config,config.app", "edit,save", "/api/admin/config/:id"},
	}},
	{"role", "角色管理", []PermissionDef{
		{"role:list", "查看角色列表", 1, "role,role.list", "search,filter", "/api/admin/roles"},
		{"role:create", "创建角色", 3, "role,role.list", "create", "/api/admin/roles"},
		{"role:update", "编辑角色", 3, "role,role.list", "edit,save,assign_permissions", "/api/admin/roles/:id,/api/admin/roles/:id/permissions"},
		{"role:delete", "删除角色", 4, "role,role.list", "delete", "/api/admin/roles/:id"},
	}},
}

// splitString 分割字符串并去除空白
func splitString(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// buildPermission 从简化定义构建完整权限
func buildPermission(module string, def PermissionDef) Permission {
	return Permission{
		Code:      def.Code,
		Name:      def.Name,
		Module:    module,
		MenuSigns: splitString(def.MenuSigns),
		Buttons:   splitString(def.Buttons),
		APIs:      splitString(def.APIs),
		Level:     def.Level,
	}
}

// GetAllPermissions 获取所有权限定义（静态配置）
func GetAllPermissions() []PermissionGroup {
	groups := make([]PermissionGroup, 0, len(permissionData))

	for _, moduleDef := range permissionData {
		permissions := make([]Permission, 0, len(moduleDef.Permissions))
		for _, permDef := range moduleDef.Permissions {
			permissions = append(permissions, buildPermission(moduleDef.Module, permDef))
		}

		groups = append(groups, PermissionGroup{
			Module:      moduleDef.Module,
			Name:        moduleDef.Name,
			Permissions: permissions,
		})
	}

	return groups
}

// GetPermissionCodes 获取所有权限代码列表
func GetPermissionCodes() []string {
	var codes []string
	for _, group := range GetAllPermissions() {
		for _, permission := range group.Permissions {
			codes = append(codes, permission.Code)
		}
	}
	return codes
}

// GetPermissionByCode 根据权限代码获取权限详情
func GetPermissionByCode(code string) *Permission {
	for _, group := range GetAllPermissions() {
		for _, permission := range group.Permissions {
			if permission.Code == code {
				return &permission
			}
		}
	}
	return nil
}

// GetPermissionsByModule 根据模块获取权限
func GetPermissionsByModule(module string) []Permission {
	for _, group := range GetAllPermissions() {
		if group.Module == module {
			return group.Permissions
		}
	}
	return nil
}

// GetMenuSignsByPermissions 根据权限列表获取所有关联的菜单标识
func GetMenuSignsByPermissions(permissionCodes []string) []string {
	menuSignSet := make(map[string]bool)

	for _, code := range permissionCodes {
		if permission := GetPermissionByCode(code); permission != nil {
			for _, menuSign := range permission.MenuSigns {
				menuSignSet[menuSign] = true
			}
		}
	}

	var menuSigns []string
	for menuSign := range menuSignSet {
		menuSigns = append(menuSigns, menuSign)
	}

	return menuSigns
}

// GetAPIPathsByPermissions 根据权限列表获取所有关联的API路径
func GetAPIPathsByPermissions(permissionCodes []string) []string {
	apiPathSet := make(map[string]bool)

	for _, code := range permissionCodes {
		if permission := GetPermissionByCode(code); permission != nil {
			for _, apiPath := range permission.APIs {
				apiPathSet[apiPath] = true
			}
		}
	}

	var apiPaths []string
	for apiPath := range apiPathSet {
		apiPaths = append(apiPaths, apiPath)
	}

	return apiPaths
}

// HasPermissionForMenu 检查权限列表是否包含访问指定菜单的权限
func HasPermissionForMenu(userPermissions []string, menuSign string) bool {
	for _, code := range userPermissions {
		if permission := GetPermissionByCode(code); permission != nil {
			for _, sign := range permission.MenuSigns {
				if sign == menuSign {
					return true
				}
			}
		}
	}
	return false
}

// HasPermissionForAPI 检查权限列表是否包含访问指定API的权限
func HasPermissionForAPI(userPermissions []string, apiPath string) bool {
	for _, code := range userPermissions {
		if permission := GetPermissionByCode(code); permission != nil {
			for _, path := range permission.APIs {
				if matchAPIPath(path, apiPath) {
					return true
				}
			}
		}
	}
	return false
}

// HasPermissionForButton 检查权限列表是否包含指定按钮的权限
func HasPermissionForButton(userPermissions []string, buttonKey string) bool {
	for _, code := range userPermissions {
		if permission := GetPermissionByCode(code); permission != nil {
			for _, button := range permission.Buttons {
				if button == buttonKey {
					return true
				}
			}
		}
	}
	return false
}

// matchAPIPath 匹配API路径（支持参数路径如 /api/users/:id）
func matchAPIPath(pattern, path string) bool {
	// 简单实现，实际项目中可以使用更复杂的路径匹配
	if pattern == path {
		return true
	}
	// TODO: 实现参数路径匹配，如 /api/users/:id 匹配 /api/users/123
	return false
}
