package permission

// Permission 权限定义
type Permission struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// 权限常量定义
const (
	// 系统管理权限
	SystemView   = "system:view"
	SystemManage = "system:manage"

	// 用户管理权限
	UserView          = "user:view"
	UserCreate        = "user:create"
	UserUpdate        = "user:update"
	UserDelete        = "user:delete"
	UserManageStatus  = "user:manage_status"
	UserResetPassword = "user:reset_password"

	// 角色管理权限
	RoleView            = "role:view"
	RoleCreate          = "role:create"
	RoleUpdate          = "role:update"
	RoleDelete          = "role:delete"
	RoleAssignPermissions = "role:assign_permissions"

	// 个人资料权限
	ProfileView         = "profile:view"
	ProfileUpdate       = "profile:update"
	ProfileChangePassword = "profile:change_password"
)

// AllPermissions 所有权限定义
var AllPermissions = []Permission{
	// 系统管理
	{Code: SystemView, Name: "查看系统", Description: "查看系统管理页面", Category: "系统管理"},
	{Code: SystemManage, Name: "管理系统", Description: "管理系统设置", Category: "系统管理"},

	// 用户管理
	{Code: UserView, Name: "查看用户", Description: "查看用户列表和详情", Category: "用户管理"},
	{Code: UserCreate, Name: "创建用户", Description: "创建新用户", Category: "用户管理"},
	{Code: UserUpdate, Name: "编辑用户", Description: "编辑用户信息", Category: "用户管理"},
	{Code: UserDelete, Name: "删除用户", Description: "删除用户", Category: "用户管理"},
	{Code: UserManageStatus, Name: "管理用户状态", Description: "启用/禁用用户", Category: "用户管理"},
	{Code: UserResetPassword, Name: "重置密码", Description: "重置用户密码", Category: "用户管理"},

	// 角色管理
	{Code: RoleView, Name: "查看角色", Description: "查看角色列表和详情", Category: "角色管理"},
	{Code: RoleCreate, Name: "创建角色", Description: "创建新角色", Category: "角色管理"},
	{Code: RoleUpdate, Name: "编辑角色", Description: "编辑角色信息", Category: "角色管理"},
	{Code: RoleDelete, Name: "删除角色", Description: "删除角色", Category: "角色管理"},
	{Code: RoleAssignPermissions, Name: "分配权限", Description: "为角色分配权限", Category: "角色管理"},

	// 个人资料
	{Code: ProfileView, Name: "查看个人资料", Description: "查看个人资料", Category: "个人资料"},
	{Code: ProfileUpdate, Name: "更新个人资料", Description: "更新个人资料", Category: "个人资料"},
	{Code: ProfileChangePassword, Name: "修改密码", Description: "修改个人密码", Category: "个人资料"},
}

// RolePermissions 角色权限映射
var RolePermissions = map[string][]string{
	"admin": {
		// 系统管理
		SystemView, SystemManage,
		// 用户管理
		UserView, UserCreate, UserUpdate, UserDelete, UserManageStatus, UserResetPassword,
		// 角色管理
		RoleView, RoleCreate, RoleUpdate, RoleDelete, RoleAssignPermissions,
		// 个人资料
		ProfileView, ProfileUpdate, ProfileChangePassword,
	},
	"manager": {
		// 用户管理（部分权限）
		UserView, UserCreate, UserUpdate, UserManageStatus,
		// 个人资料
		ProfileView, ProfileUpdate, ProfileChangePassword,
	},
	"user": {
		// 个人资料
		ProfileView, ProfileUpdate, ProfileChangePassword,
	},
}

// HasPermission 检查用户是否有指定权限
func HasPermission(userRole, permission string) bool {
	permissions, exists := RolePermissions[userRole]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// HasAllPermissions 检查用户是否有所有指定权限
func HasAllPermissions(userRole string, permissions []string) bool {
	for _, permission := range permissions {
		if !HasPermission(userRole, permission) {
			return false
		}
	}
	return true
}

// HasAnyPermission 检查用户是否有任意一个指定权限
func HasAnyPermission(userRole string, permissions []string) bool {
	for _, permission := range permissions {
		if HasPermission(userRole, permission) {
			return true
		}
	}
	return false
}

// GetUserPermissions 获取用户所有权限
func GetUserPermissions(userRole string) []string {
	permissions, exists := RolePermissions[userRole]
	if !exists {
		return []string{}
	}
	return permissions
}

// GetPermissionsByCategory 按分类获取权限
func GetPermissionsByCategory() map[string][]Permission {
	categories := make(map[string][]Permission)
	for _, permission := range AllPermissions {
		categories[permission.Category] = append(categories[permission.Category], permission)
	}
	return categories
}
