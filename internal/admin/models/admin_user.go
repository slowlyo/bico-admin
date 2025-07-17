package models

import (
	"time"

	"bico-admin/internal/admin/definitions"
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
	Status      *int       `json:"status" gorm:"default:1;comment:状态：1-启用，0-禁用"`
	LastLoginAt *time.Time `json:"last_login_at" gorm:"comment:最后登录时间"`
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
	return u.Status != nil && *u.Status == types.StatusActive
}

// GetStatusText 获取状态文本
func (u *AdminUser) GetStatusText() string {
	if u.Status == nil {
		return "未知"
	}
	switch *u.Status {
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
	// 超级管理员拥有所有权限
	if u.IsSuperAdmin() {
		return definitions.GetPermissionCodes()
	}

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
	// 超级管理员拥有所有权限
	if u.IsSuperAdmin() {
		return true
	}

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

// IsSystemSuperAdmin 检查是否为系统默认超级管理员（不可删除）
func (u *AdminUser) IsSystemSuperAdmin() bool {
	return u.Username == "admin" && u.IsSuperAdmin()
}

// CanBeDeleted 检查用户是否可以被删除
func (u *AdminUser) CanBeDeleted() bool {
	// 系统默认超级管理员不可删除
	return !u.IsSystemSuperAdmin()
}

// CanBeModified 检查用户是否可以被修改（某些字段）
func (u *AdminUser) CanBeModified() bool {
	// 系统默认超级管理员的用户名和角色不可修改
	return !u.IsSystemSuperAdmin()
}

// CanBeDisabled 检查用户是否可以被禁用（需要运行时检查）
func (u *AdminUser) CanBeDisabled() bool {
	// 如果不是超级管理员，可以禁用
	if !u.IsSuperAdmin() {
		return true
	}
	// 超级管理员需要运行时检查是否还有其他超管
	// 这里返回false，具体检查在service层进行
	return false
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

// ToUserInfo 转换为用户信息
func (u *AdminUser) ToUserInfo() types.UserInfo {
	status := 0
	if u.Status != nil {
		status = *u.Status
	}
	return types.UserInfo{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Name, // AdminUser 使用 Name 字段作为昵称
		Email:    u.Email,
		Avatar:   u.Avatar,
		UserType: types.UserTypeAdmin, // 管理员用户类型
		Status:   status,
	}
}
