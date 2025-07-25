package definitions

import (
	"strings"
)

// 注意：生成的权限代码应该直接添加到 getAllPermissionDefs 函数中

// Permission 权限定义（静态配置，程序中定义）
type Permission struct {
	Code     string       `json:"code"`      // 权限代码，如 "system.admin_user:list"
	Name     string       `json:"name"`      // 权限名称，如 "查看管理员列表"
	ParentID string       `json:"parent_id"` // 父权限ID，空表示根权限
	Type     string       `json:"type"`      // 权限类型：module(模块), action(操作)
	Buttons  []string     `json:"buttons"`   // 关联的按钮标识列表
	APIs     []string     `json:"apis"`      // 关联的API路径列表
	Level    int          `json:"level"`     // 权限级别：1-查看，2-操作，3-管理，4-超级
	Children []Permission `json:"children"`  // 子权限列表
}

// PermissionTree 权限树结构
type PermissionTree struct {
	Permissions []Permission `json:"permissions"` // 权限树
}

// 权限级别常量
const (
	PermissionLevelView   = 1 // 查看权限
	PermissionLevelAction = 2 // 操作权限
	PermissionLevelManage = 3 // 管理权限
	PermissionLevelSuper  = 4 // 超级权限
)

// 权限类型常量
const (
	PermissionTypeModule = "module" // 模块权限
	PermissionTypeAction = "action" // 操作权限
)

// PermissionDef 权限定义的简化结构
type PermissionDef struct {
	Code    string
	Name    string
	Parent  string
	Type    string
	Level   int
	Buttons string // 逗号分隔
	APIs    string // 逗号分隔
}

// getAllPermissionDefs 获取所有权限定义（包括生成的）
func getAllPermissionDefs() []PermissionDef {
	// 基础权限定义
	baseDefs := []PermissionDef{
		// 系统管理
		{"system", "系统管理", "", "module", 1, "", ""},
		{"system.admin_user", "管理员", "system", "module", 1, "", ""},
		{"system.admin_user:list", "查看管理员列表", "system.admin_user", "action", 1, "search,filter", "/admin-api/admin-users,/admin-api/admin-users/:id"},
		{"system.admin_user:create", "创建管理员", "system.admin_user", "action", 3, "create", "/admin-api/admin-users,/admin-api/roles/options"},
		{"system.admin_user:update", "编辑管理员", "system.admin_user", "action", 3, "edit,save", "/admin-api/admin-users/:id,/admin-api/admin-users/:id/status,/admin-api/roles/options"},
		{"system.admin_user:delete", "删除管理员", "system.admin_user", "action", 4, "delete", "/admin-api/admin-users/:id"},

		{"system.role", "角色", "system", "module", 1, "", ""},
		{"system.role:list", "查看角色列表", "system.role", "action", 1, "search,filter", "/admin-api/roles,/admin-api/roles/:id"},
		{"system.role:create", "创建角色", "system.role", "action", 3, "create", "/admin-api/roles,/admin-api/roles/permissions"},
		{"system.role:update", "编辑角色", "system.role", "action", 3, "edit,save,assign_permissions", "/admin-api/roles/:id,/admin-api/roles/:id/status,/admin-api/roles/:id/permissions,/admin-api/roles/permissions,/admin-api/roles/assign"},
		{"system.role:delete", "删除角色", "system.role", "action", 4, "delete", "/admin-api/roles/:id"},
	}

	// 注意：生成的权限定义应该直接添加到上面的 baseDefs 数组中
	// 不再使用动态注册模式

	return baseDefs
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

// buildPermissionTree 从扁平化定义构建权限树
func buildPermissionTree() []Permission {
	// 获取所有权限定义
	allDefs := getAllPermissionDefs()

	// 先构建所有权限的映射
	permMap := make(map[string]*Permission)
	for _, def := range allDefs {
		perm := &Permission{
			Code:     def.Code,
			Name:     def.Name,
			ParentID: def.Parent,
			Type:     def.Type,
			Level:    def.Level,
			Buttons:  splitString(def.Buttons),
			APIs:     splitString(def.APIs),
			Children: []Permission{},
		}
		permMap[def.Code] = perm
	}

	// 构建树形结构 - 递归构建完整的权限树
	var buildChildren func(parentCode string) []Permission
	buildChildren = func(parentCode string) []Permission {
		var children []Permission
		for _, def := range allDefs {
			if def.Parent == parentCode {
				perm := permMap[def.Code]
				// 递归构建子权限
				perm.Children = buildChildren(def.Code)
				children = append(children, *perm)
			}
		}
		return children
	}

	// 收集根权限并构建完整的树
	var roots []Permission
	for _, def := range allDefs {
		if def.Parent == "" {
			perm := permMap[def.Code]
			perm.Children = buildChildren(def.Code)
			roots = append(roots, *perm)
		}
	}

	return roots
}

// 权限数据（懒加载）
var permissionData []Permission

// 初始化权限数据
func init() {
	permissionData = buildPermissionTree()
}

// GetAllPermissions 获取所有权限定义（静态配置）
func GetAllPermissions() []Permission {
	return permissionData
}

// GetPermissionTree 获取权限树
func GetPermissionTree() PermissionTree {
	return PermissionTree{
		Permissions: permissionData,
	}
}

// flattenPermissions 递归展平权限树，获取所有权限的扁平列表
func flattenPermissions(permissions []Permission) []Permission {
	var result []Permission
	for _, perm := range permissions {
		result = append(result, perm)
		if len(perm.Children) > 0 {
			result = append(result, flattenPermissions(perm.Children)...)
		}
	}
	return result
}

// GetAllPermissionsFlat 获取所有权限的扁平列表
func GetAllPermissionsFlat() []Permission {
	return flattenPermissions(permissionData)
}

// GetPermissionCodes 获取所有权限代码列表
func GetPermissionCodes() []string {
	var codes []string
	allPermissions := GetAllPermissionsFlat()
	for _, permission := range allPermissions {
		codes = append(codes, permission.Code)
	}
	return codes
}

// GetPermissionByCode 根据权限代码获取权限详情
func GetPermissionByCode(code string) *Permission {
	return findPermissionByCode(permissionData, code)
}

// findPermissionByCode 递归查找权限
func findPermissionByCode(permissions []Permission, code string) *Permission {
	for _, perm := range permissions {
		if perm.Code == code {
			return &perm
		}
		if len(perm.Children) > 0 {
			if found := findPermissionByCode(perm.Children, code); found != nil {
				return found
			}
		}
	}
	return nil
}

// GetPermissionsByParent 根据父权限ID获取子权限
func GetPermissionsByParent(parentID string) []Permission {
	if parentID == "" {
		return permissionData // 返回根权限
	}

	parent := GetPermissionByCode(parentID)
	if parent != nil {
		return parent.Children
	}
	return nil
}

// GetPermissionsByType 根据权限类型获取权限
func GetPermissionsByType(permType string) []Permission {
	var result []Permission
	allPermissions := GetAllPermissionsFlat()
	for _, perm := range allPermissions {
		if perm.Type == permType {
			result = append(result, perm)
		}
	}
	return result
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
	// 简单的完全匹配
	if pattern == path {
		return true
	}

	// 处理参数路径匹配，如 /admin-api/admin-users/:id 匹配 /admin-api/admin-users/123
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	// 路径段数量必须相同
	if len(patternParts) != len(pathParts) {
		return false
	}

	// 逐段比较
	for i, patternPart := range patternParts {
		pathPart := pathParts[i]

		// 如果是参数段（以:开头），则跳过比较
		if strings.HasPrefix(patternPart, ":") {
			continue
		}

		// 普通段必须完全匹配
		if patternPart != pathPart {
			return false
		}
	}

	return true
}
