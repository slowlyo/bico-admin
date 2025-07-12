package definitions

import "bico-admin/internal/admin/types"

// MenuDefinition 菜单定义（扩展版本，包含权限要求）
type MenuDefinition struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`        // 菜单名称
	Path        string           `json:"path"`        // 路由路径
	Icon        string           `json:"icon"`        // 图标
	Sort        int              `json:"sort"`        // 排序
	ParentID    uint             `json:"parent_id"`   // 父菜单ID，0表示顶级菜单
	Permissions []string         `json:"permissions"` // 访问该菜单需要的权限
	Children    []MenuDefinition `json:"children"`    // 子菜单
}

// GetAllMenus 获取所有菜单定义
func GetAllMenus() []MenuDefinition {
	return []MenuDefinition{
		{
			ID:          1,
			Name:        "用户管理",
			Path:        "/admin/users",
			Icon:        "user",
			Sort:        1,
			ParentID:    0,
			Permissions: []string{"user:read"},
			Children: []MenuDefinition{
				{
					ID:          11,
					Name:        "用户列表",
					Path:        "/admin/users/list",
					Icon:        "list",
					Sort:        1,
					ParentID:    1,
					Permissions: []string{"user:read"},
				},
				{
					ID:          12,
					Name:        "用户详情",
					Path:        "/admin/users/detail",
					Icon:        "detail",
					Sort:        2,
					ParentID:    1,
					Permissions: []string{"user:read"},
				},
			},
		},
		{
			ID:          2,
			Name:        "管理员管理",
			Path:        "/admin/admin-users",
			Icon:        "admin",
			Sort:        2,
			ParentID:    0,
			Permissions: []string{"admin_user:read"},
			Children: []MenuDefinition{
				{
					ID:          21,
					Name:        "管理员列表",
					Path:        "/admin/admin-users/list",
					Icon:        "list",
					Sort:        1,
					ParentID:    2,
					Permissions: []string{"admin_user:read"},
				},
				{
					ID:          22,
					Name:        "添加管理员",
					Path:        "/admin/admin-users/create",
					Icon:        "plus",
					Sort:        2,
					ParentID:    2,
					Permissions: []string{"admin_user:write"},
				},
			},
		},
		{
			ID:          3,
			Name:        "系统管理",
			Path:        "/admin/system",
			Icon:        "setting",
			Sort:        3,
			ParentID:    0,
			Permissions: []string{"system:read"},
			Children: []MenuDefinition{
				{
					ID:          31,
					Name:        "系统信息",
					Path:        "/admin/system/info",
					Icon:        "info",
					Sort:        1,
					ParentID:    3,
					Permissions: []string{"system:read"},
				},
				{
					ID:          32,
					Name:        "系统配置",
					Path:        "/admin/system/config",
					Icon:        "config",
					Sort:        2,
					ParentID:    3,
					Permissions: []string{"system:write"},
				},
			},
		},
		{
			ID:          4,
			Name:        "配置管理",
			Path:        "/admin/config",
			Icon:        "tool",
			Sort:        4,
			ParentID:    0,
			Permissions: []string{"config:read"},
			Children: []MenuDefinition{
				{
					ID:          41,
					Name:        "应用配置",
					Path:        "/admin/config/app",
					Icon:        "app",
					Sort:        1,
					ParentID:    4,
					Permissions: []string{"config:read"},
				},
			},
		},
	}
}

// FilterMenusByPermissions 根据用户权限过滤菜单
func FilterMenusByPermissions(userPermissions []string) []types.Menu {
	allMenus := GetAllMenus()
	var filteredMenus []types.Menu

	// 创建权限映射表，提高查找效率
	permissionMap := make(map[string]bool)
	for _, perm := range userPermissions {
		permissionMap[perm] = true
	}

	// 递归过滤菜单
	for _, menu := range allMenus {
		if hasPermission(menu.Permissions, permissionMap) {
			filteredMenu := convertToMenu(menu, permissionMap)
			filteredMenus = append(filteredMenus, filteredMenu)
		}
	}

	return filteredMenus
}

// hasPermission 检查是否有权限访问菜单
func hasPermission(requiredPermissions []string, userPermissions map[string]bool) bool {
	if len(requiredPermissions) == 0 {
		return true // 无权限要求，所有人都可访问
	}

	// 只要有一个权限匹配就可以访问
	for _, perm := range requiredPermissions {
		if userPermissions[perm] {
			return true
		}
	}
	return false
}

// convertToMenu 将MenuDefinition转换为types.Menu，并递归处理子菜单
func convertToMenu(menuDef MenuDefinition, userPermissions map[string]bool) types.Menu {
	menu := types.Menu{
		ID:   menuDef.ID,
		Name: menuDef.Name,
		Path: menuDef.Path,
		Icon: menuDef.Icon,
		Sort: menuDef.Sort,
	}

	// 递归处理子菜单
	for _, child := range menuDef.Children {
		if hasPermission(child.Permissions, userPermissions) {
			childMenu := convertToMenu(child, userPermissions)
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
		if menu.ParentID == 0 {
			// 复制菜单但不包含子菜单
			topMenu := menu
			topMenu.Children = nil
			topMenus = append(topMenus, topMenu)
		}
	}

	return topMenus
}

// GetMenuByID 根据ID获取菜单
func GetMenuByID(id uint) *MenuDefinition {
	allMenus := GetAllMenus()
	return findMenuByID(allMenus, id)
}

// findMenuByID 递归查找菜单
func findMenuByID(menus []MenuDefinition, id uint) *MenuDefinition {
	for _, menu := range menus {
		if menu.ID == id {
			return &menu
		}
		if found := findMenuByID(menu.Children, id); found != nil {
			return found
		}
	}
	return nil
}
