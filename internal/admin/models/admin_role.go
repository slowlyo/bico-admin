package models

import (
	"time"

	"bico-admin/internal/shared/types"
)

// AdminRole 管理员角色模型
type AdminRole struct {
	types.BaseModel
	Name        string    `json:"name" gorm:"size:100;not null;comment:角色名称"`
	Code        string    `json:"code" gorm:"size:50;uniqueIndex;not null;comment:角色代码"`
	Description string    `json:"description" gorm:"size:500;comment:角色描述"`
	Status      int       `json:"status" gorm:"default:1;comment:状态：1-启用，0-禁用"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联关系
	Permissions []AdminRolePermission `json:"permissions" gorm:"foreignKey:RoleID"`
	Users       []AdminUserRole       `json:"users" gorm:"foreignKey:RoleID"`
}

// AdminRolePermission 管理员角色权限关联模型
type AdminRolePermission struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	RoleID         uint      `json:"role_id" gorm:"not null;comment:角色ID"`
	PermissionCode string    `json:"permission_code" gorm:"size:100;not null;comment:权限代码"`
	CreatedAt      time.Time `json:"created_at"`

	// 关联关系
	Role AdminRole `json:"role" gorm:"foreignKey:RoleID"`
}

// AdminUserRole 管理员用户角色关联模型
type AdminUserRole struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;comment:用户ID"`
	RoleID    uint      `json:"role_id" gorm:"not null;comment:角色ID"`
	CreatedAt time.Time `json:"created_at"`

	// 关联关系
	User AdminUser `json:"user" gorm:"foreignKey:UserID"`
	Role AdminRole `json:"role" gorm:"foreignKey:RoleID"`
}

// TableName 指定表名
func (AdminRole) TableName() string {
	return "admin_roles"
}

// TableName 指定表名
func (AdminRolePermission) TableName() string {
	return "admin_role_permissions"
}

// TableName 指定表名
func (AdminUserRole) TableName() string {
	return "admin_user_roles"
}

// GetPermissionCodes 获取角色的所有权限代码
func (r *AdminRole) GetPermissionCodes() []string {
	var codes []string
	for _, permission := range r.Permissions {
		codes = append(codes, permission.PermissionCode)
	}
	return codes
}

// HasPermission 检查角色是否有指定权限
func (r *AdminRole) HasPermission(permissionCode string) bool {
	for _, permission := range r.Permissions {
		if permission.PermissionCode == permissionCode {
			return true
		}
	}
	return false
}

// IsEnabled 检查角色是否启用
func (r *AdminRole) IsEnabled() bool {
	return r.Status == types.StatusActive
}

// 预定义角色代码常量
const (
	RoleCodeSuperAdmin = "super_admin" // 超级管理员
	RoleCodeAdmin      = "admin"       // 管理员
	RoleCodeOperator   = "operator"    // 操作员
	RoleCodeViewer     = "viewer"      // 查看者
)

// GetDefaultRoles 获取默认角色配置
func GetDefaultRoles() []AdminRole {
	return []AdminRole{
		{
			Name:        "超级管理员",
			Code:        RoleCodeSuperAdmin,
			Description: "拥有系统所有权限的超级管理员",
			Status:      types.StatusActive,
		},
		{
			Name:        "管理员",
			Code:        RoleCodeAdmin,
			Description: "拥有大部分管理权限的管理员",
			Status:      types.StatusActive,
		},
		{
			Name:        "操作员",
			Code:        RoleCodeOperator,
			Description: "拥有基本操作权限的操作员",
			Status:      types.StatusActive,
		},
		{
			Name:        "查看者",
			Code:        RoleCodeViewer,
			Description: "只有查看权限的用户",
			Status:      types.StatusActive,
		},
	}
}
