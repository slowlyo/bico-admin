package definitions

import (
	"bico-admin/internal/admin/types"
	"strings"
)

// MenuDefinition 菜单定义（扩展版本，包含权限要求）
type MenuDefinition struct {
	Sign        string           `json:"sign"`        // 菜单唯一标识，如 "user.list"
	Name        string           `json:"name"`        // 菜单名称
	Path        string           `json:"path"`        // 路由路径
	Icon        string           `json:"icon"`        // 图标
	Sort        int              `json:"sort"`        // 排序
	ParentSign  string           `json:"parent_sign"` // 父菜单标识，空表示顶级菜单
	Permissions []string         `json:"permissions"` // 访问该菜单需要的权限
	Children    []MenuDefinition `json:"children"`    // 子菜单
}

// MenuDef 菜单定义的简化结构（用于生成完整菜单）
type MenuDef struct {
	Sign        string
	Name        string
	Path        string
	Icon        string
	Sort        int
	ParentSign  string
	Permissions string // 逗号分隔的权限列表
}

// 菜单定义数据（按层级关系组织）
var menuData = []MenuDef{
	// 用户管理模块
	{"user", "用户管理", "/admin/users", "user", 1, "", "user:list"},
	{"user.list", "用户列表", "/admin/users/list", "list", 1, "user", "user:list"},
	{"user.detail", "用户详情", "/admin/users/detail", "detail", 2, "user", "user:list"},

	// 管理员管理模块
	{"admin_user", "管理员管理", "/admin/admin-users", "admin", 2, "", "admin_user:list"},
	{"admin_user.list", "管理员列表", "/admin/admin-users/list", "list", 1, "admin_user", "admin_user:list"},
	{"admin_user.create", "添加管理员", "/admin/admin-users/create", "plus", 2, "admin_user", "admin_user:create"},

	// 系统管理模块
	{"system", "系统管理", "/admin/system", "setting", 3, "", "system:info"},
	{"system.info", "系统信息", "/admin/system/info", "info", 1, "system", "system:info"},
	{"system.config", "系统配置", "/admin/system/config", "config", 2, "system", "system:config"},

	// 配置管理模块
	{"config", "配置管理", "/admin/config", "tool", 4, "", "config:list"},
	{"config.app", "应用配置", "/admin/config/app", "app", 1, "config", "config:list"},

	// 角色管理模块
	{"role", "角色管理", "/admin/roles", "team", 5, "", "role:list"},
	{"role.list", "角色列表", "/admin/roles/list", "list", 1, "role", "role:list"},
}

// splitPermissions 分割权限字符串
func splitPermissions(s string) []string {
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

// buildMenuDefinition 从简化定义构建完整菜单
func buildMenuDefinition(def MenuDef) MenuDefinition {
	return MenuDefinition{
		Sign:        def.Sign,
		Name:        def.Name,
		Path:        def.Path,
		Icon:        def.Icon,
		Sort:        def.Sort,
		ParentSign:  def.ParentSign,
		Permissions: splitPermissions(def.Permissions),
		Children:    []MenuDefinition{}, // 初始化为空，后续填充
	}
}

// GetAllMenus 获取所有菜单定义
func GetAllMenus() []MenuDefinition {
	// 构建所有菜单的映射
	menuMap := make(map[string]*MenuDefinition)
	var topMenus []MenuDefinition

	// 第一遍：创建所有菜单
	for _, def := range menuData {
		menu := buildMenuDefinition(def)
		menuMap[def.Sign] = &menu

		if def.ParentSign == "" {
			topMenus = append(topMenus, menu)
		}
	}

	// 第二遍：建立父子关系
	for _, def := range menuData {
		if def.ParentSign != "" {
			if parent, exists := menuMap[def.ParentSign]; exists {
				if child, exists := menuMap[def.Sign]; exists {
					parent.Children = append(parent.Children, *child)
				}
			}
		}
	}

	// 更新顶级菜单的子菜单
	for i, topMenu := range topMenus {
		if updated, exists := menuMap[topMenu.Sign]; exists {
			topMenus[i] = *updated
		}
	}

	return topMenus
}

// FilterMenusByPermissions 根据用户权限过滤菜单（优化版本）
func FilterMenusByPermissions(userPermissions []string) []types.Menu {
	allMenus := GetAllMenus()
	var filteredMenus []types.Menu

	// 获取用户权限对应的所有菜单标识
	allowedMenuSigns := GetMenuSignsByPermissions(userPermissions)
	menuSignSet := make(map[string]bool)
	for _, sign := range allowedMenuSigns {
		menuSignSet[sign] = true
	}

	// 递归过滤菜单
	for _, menu := range allMenus {
		if menuSignSet[menu.Sign] || hasPermissionLegacy(menu.Permissions, userPermissions) {
			filteredMenu := convertToMenuWithPermissionCheck(menu, menuSignSet, userPermissions)
			filteredMenus = append(filteredMenus, filteredMenu)
		}
	}

	return filteredMenus
}

// hasPermissionLegacy 兼容旧版权限检查（使用字符串切片）
func hasPermissionLegacy(requiredPermissions []string, userPermissions []string) bool {
	if len(requiredPermissions) == 0 {
		return true // 无权限要求，所有人都可访问
	}

	// 创建用户权限映射表
	userPermMap := make(map[string]bool)
	for _, perm := range userPermissions {
		userPermMap[perm] = true
	}

	// 只要有一个权限匹配就可以访问
	for _, perm := range requiredPermissions {
		if userPermMap[perm] {
			return true
		}
	}
	return false
}

// convertToMenuWithPermissionCheck 将MenuDefinition转换为types.Menu，并递归处理子菜单（带权限检查）
func convertToMenuWithPermissionCheck(menuDef MenuDefinition, menuSignSet map[string]bool, userPermissions []string) types.Menu {
	menu := types.Menu{
		Sign: menuDef.Sign,
		Name: menuDef.Name,
		Path: menuDef.Path,
		Icon: menuDef.Icon,
		Sort: menuDef.Sort,
	}

	// 递归处理子菜单
	for _, child := range menuDef.Children {
		if menuSignSet[child.Sign] || hasPermissionLegacy(child.Permissions, userPermissions) {
			childMenu := convertToMenuWithPermissionCheck(child, menuSignSet, userPermissions)
			menu.Children = append(menu.Children, childMenu)
		}
	}

	return menu
}

// GetTopLevelMenus 获取顶级菜单（不包含子菜单）
func GetTopLevelMenus() []MenuDefinition {
	allMenus := GetAllMenus()
	var topMenus []MenuDefinition

	for _, menu := range allMenus {
		if menu.ParentSign == "" {
			// 复制菜单但不包含子菜单
			topMenu := menu
			topMenu.Children = nil
			topMenus = append(topMenus, topMenu)
		}
	}

	return topMenus
}

// GetMenuBySign 根据标识获取菜单
func GetMenuBySign(sign string) *MenuDefinition {
	allMenus := GetAllMenus()
	return findMenuBySign(allMenus, sign)
}

// findMenuBySign 递归查找菜单
func findMenuBySign(menus []MenuDefinition, sign string) *MenuDefinition {
	for _, menu := range menus {
		if menu.Sign == sign {
			return &menu
		}
		if found := findMenuBySign(menu.Children, sign); found != nil {
			return found
		}
	}
	return nil
}
