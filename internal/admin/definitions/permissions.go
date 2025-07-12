package definitions

// Permission 权限定义
type Permission struct {
	Code   string `json:"code"`   // 权限代码，如 "user:read"
	Name   string `json:"name"`   // 权限名称，如 "查看用户"
	Module string `json:"module"` // 所属模块，如 "user", "system"
}

// PermissionGroup 权限组
type PermissionGroup struct {
	Module      string       `json:"module"`      // 模块名称
	Name        string       `json:"name"`        // 模块显示名称
	Permissions []Permission `json:"permissions"` // 该模块下的权限列表
}

// ModuleConfig 模块配置
type ModuleConfig struct {
	Module          string            // 模块名称
	Name            string            // 模块显示名称
	StandardActions []string          // 标准操作：list, create, update, delete
	CustomActions   map[string]string // 自定义操作：code -> name
}

// getModuleConfigs 获取模块配置
func getModuleConfigs() []ModuleConfig {
	return []ModuleConfig{
		{
			Module:          "user",
			Name:            "用户",
			StandardActions: []string{"list", "create", "update", "delete"},
			CustomActions: map[string]string{
				"export": "导出用户",
				"import": "导入用户",
			},
		},
		{
			Module:          "admin_user",
			Name:            "管理员",
			StandardActions: []string{"list", "create", "update", "delete"},
			CustomActions: map[string]string{
				"reset_password": "重置密码",
			},
		},
		{
			Module:          "system",
			Name:            "系统",
			StandardActions: []string{"list", "update"},
			CustomActions: map[string]string{
				"backup":  "系统备份",
				"restore": "系统恢复",
			},
		},
		{
			Module:          "config",
			Name:            "配置",
			StandardActions: []string{"list", "update"},
			CustomActions:   nil,
		},
	}
}

// getActionName 获取操作的中文名称
func getActionName(action string) string {
	standardActionNames := map[string]string{
		"list":   "列表",
		"create": "创建",
		"update": "修改",
		"delete": "删除",
	}
	if name, ok := standardActionNames[action]; ok {
		return name
	}
	return action
}

// buildPermissions 根据模块配置构建权限
func buildPermissions(moduleConfig ModuleConfig) []Permission {
	var permissions []Permission

	// 构建标准操作权限
	for _, action := range moduleConfig.StandardActions {
		permissions = append(permissions, Permission{
			Code:   moduleConfig.Module + ":" + action,
			Name:   getActionName(action) + moduleConfig.Name,
			Module: moduleConfig.Module,
		})
	}

	// 构建自定义操作权限
	for action, name := range moduleConfig.CustomActions {
		permissions = append(permissions, Permission{
			Code:   moduleConfig.Module + ":" + action,
			Name:   name,
			Module: moduleConfig.Module,
		})
	}

	return permissions
}

// GetAllPermissions 获取所有权限定义
func GetAllPermissions() []PermissionGroup {
	var groups []PermissionGroup

	for _, moduleConfig := range getModuleConfigs() {
		group := PermissionGroup{
			Module:      moduleConfig.Module,
			Name:        moduleConfig.Name,
			Permissions: buildPermissions(moduleConfig),
		}
		groups = append(groups, group)
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

// GetPermissionsByModule 根据模块获取权限
func GetPermissionsByModule(module string) []Permission {
	for _, group := range GetAllPermissions() {
		if group.Module == module {
			return group.Permissions
		}
	}
	return nil
}
