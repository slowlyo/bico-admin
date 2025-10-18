package consts

// Permission 权限结构
type Permission struct {
	Key      string       `json:"key"`
	Label    string       `json:"label"`
	Children []Permission `json:"children,omitempty"`
}

// 权限常量定义
const (
	// 系统管理
	PermSystemManage = "system:manage"
	
	// 用户管理
	PermAdminUserMenu   = "system:admin_user:menu"
	PermAdminUserList   = "system:admin_user:list"
	PermAdminUserCreate = "system:admin_user:create"
	PermAdminUserEdit   = "system:admin_user:edit"
	PermAdminUserDelete = "system:admin_user:delete"
	
	// 角色管理
	PermAdminRoleMenu   = "system:admin_role:menu"
	PermAdminRoleList   = "system:admin_role:list"
	PermAdminRoleCreate = "system:admin_role:create"
	PermAdminRoleEdit   = "system:admin_role:edit"
	PermAdminRoleDelete = "system:admin_role:delete"
	PermAdminRolePermission = "system:admin_role:permission"
	
	// Dashboard
	PermDashboardMenu = "dashboard:menu"
)

// AllPermissions 所有权限树
var AllPermissions = []Permission{
	{
		Key:   PermDashboardMenu,
		Label: "Dashboard",
	},
	{
		Key:   PermSystemManage,
		Label: "系统管理",
		Children: []Permission{
			{
				Key:   PermAdminUserMenu,
				Label: "用户管理",
				Children: []Permission{
					{Key: PermAdminUserList, Label: "查看列表"},
					{Key: PermAdminUserCreate, Label: "创建用户"},
					{Key: PermAdminUserEdit, Label: "编辑用户"},
					{Key: PermAdminUserDelete, Label: "删除用户"},
				},
			},
			{
				Key:   PermAdminRoleMenu,
				Label: "角色管理",
				Children: []Permission{
					{Key: PermAdminRoleList, Label: "查看列表"},
					{Key: PermAdminRoleCreate, Label: "创建角色"},
					{Key: PermAdminRoleEdit, Label: "编辑角色"},
					{Key: PermAdminRoleDelete, Label: "删除角色"},
					{Key: PermAdminRolePermission, Label: "配置权限"},
				},
			},
		},
	},
}

// GetAllPermissionKeys 获取所有权限的 key 列表
func GetAllPermissionKeys() []string {
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
	
	collectKeys(AllPermissions)
	return keys
}
