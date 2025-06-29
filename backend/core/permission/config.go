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
	RoleView              = "role:view"
	RoleCreate            = "role:create"
	RoleUpdate            = "role:update"
	RoleDelete            = "role:delete"
	RoleAssignPermissions = "role:assign_permissions"

	// 注意：个人资料相关功能无需权限验证，所有登录用户都可以访问
	// 已移除 profile:* 相关权限定义
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

	// 注意：个人资料相关功能无需权限验证，已移除相关权限定义
}

// 注意：角色权限映射已移除，现在完全基于数据库的 role_permissions 关联表
// 权限验证通过数据库查询实现，支持动态权限分配

// 注意：旧的静态权限检查函数已移除
// 现在权限验证完全通过 PermissionMiddleware 的数据库查询实现

// GetPermissionsByCategory 按分类获取权限
func GetPermissionsByCategory() map[string][]Permission {
	categories := make(map[string][]Permission)
	for _, permission := range AllPermissions {
		categories[permission.Category] = append(categories[permission.Category], permission)
	}
	return categories
}

// GetAllPermissionCodes 获取所有权限代码列表
func GetAllPermissionCodes() []string {
	var codes []string
	for _, permission := range AllPermissions {
		codes = append(codes, permission.Code)
	}
	return codes
}

// GetPermissionByCode 根据权限代码获取权限信息
func GetPermissionByCode(code string) *Permission {
	for _, permission := range AllPermissions {
		if permission.Code == code {
			return &permission
		}
	}
	return nil
}

// IsSuperAdmin 检查角色是否为超级管理员
func IsSuperAdmin(roleCode string) bool {
	return roleCode == "super_admin"
}

// IsProtectedRole 检查角色是否为受保护角色（不可删除或编辑）
func IsProtectedRole(roleCode string) bool {
	return roleCode == "super_admin"
}

// IsProfileRoute 检查路由是否为个人资料相关路由（无需权限验证）
func IsProfileRoute(path string) bool {
	profileRoutes := []string{
		"/auth/profile",
		"/auth/change-password",
	}

	for _, route := range profileRoutes {
		if path == route {
			return true
		}
	}
	return false
}
