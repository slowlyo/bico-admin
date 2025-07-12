package models

import (
	"time"

	"bico-admin/internal/shared/types"
)

// AdminUser 管理员用户模型
type AdminUser struct {
	types.BaseModel
	Username    string     `json:"username" gorm:"size:50;uniqueIndex;not null;comment:用户名"`
	Password    string     `json:"-" gorm:"size:255;not null;comment:密码"`
	Name        string     `json:"name" gorm:"size:100;not null;comment:姓名"`
	Avatar      string     `json:"avatar" gorm:"size:255;comment:头像"`
	Email       string     `json:"email" gorm:"size:100;comment:邮箱"`
	Phone       string     `json:"phone" gorm:"size:20;comment:手机号"`
	Status      int        `json:"status" gorm:"default:1;comment:状态：1-启用，0-禁用"`
	LastLoginAt *time.Time `json:"last_login_at" gorm:"comment:最后登录时间"`
	LastLoginIP string     `json:"last_login_ip" gorm:"size:45;comment:最后登录IP"`
	LoginCount  int        `json:"login_count" gorm:"default:0;comment:登录次数"`
	Remark      string     `json:"remark" gorm:"size:500;comment:备注"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// 关联关系
	Roles []AdminUserRole `json:"roles" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (AdminUser) TableName() string {
	return "admin_users"
}

// IsEnabled 检查用户是否启用
func (u *AdminUser) IsEnabled() bool {
	return u.Status == types.StatusActive
}

// GetStatusText 获取状态文本
func (u *AdminUser) GetStatusText() string {
	switch u.Status {
	case types.StatusActive:
		return "启用"
	case types.StatusInactive:
		return "禁用"
	default:
		return "未知"
	}
}

// GetRoleCodes 获取用户的所有角色代码
func (u *AdminUser) GetRoleCodes() []string {
	var codes []string
	for _, userRole := range u.Roles {
		if userRole.Role.IsEnabled() {
			codes = append(codes, userRole.Role.Code)
		}
	}
	return codes
}

// GetPermissionCodes 获取用户的所有权限代码
func (u *AdminUser) GetPermissionCodes() []string {
	permissionSet := make(map[string]bool)

	for _, userRole := range u.Roles {
		if userRole.Role.IsEnabled() {
			for _, permission := range userRole.Role.Permissions {
				permissionSet[permission.PermissionCode] = true
			}
		}
	}

	var permissions []string
	for permission := range permissionSet {
		permissions = append(permissions, permission)
	}

	return permissions
}

// HasRole 检查用户是否有指定角色
func (u *AdminUser) HasRole(roleCode string) bool {
	for _, userRole := range u.Roles {
		if userRole.Role.Code == roleCode && userRole.Role.IsEnabled() {
			return true
		}
	}
	return false
}

// HasPermission 检查用户是否有指定权限
func (u *AdminUser) HasPermission(permissionCode string) bool {
	for _, userRole := range u.Roles {
		if userRole.Role.IsEnabled() && userRole.Role.HasPermission(permissionCode) {
			return true
		}
	}
	return false
}

// IsSuperAdmin 检查是否为超级管理员
func (u *AdminUser) IsSuperAdmin() bool {
	return u.HasRole(RoleCodeSuperAdmin)
}

// GetRoleNames 获取用户的所有角色名称
func (u *AdminUser) GetRoleNames() []string {
	var names []string
	for _, userRole := range u.Roles {
		if userRole.Role.IsEnabled() {
			names = append(names, userRole.Role.Name)
		}
	}
	return names
}
