package model

import "bico-admin/internal/core/model"

// AdminRole 角色模型
type AdminRole struct {
	model.BaseModel
	Name        string   `gorm:"size:64;uniqueIndex;not null" json:"name"`
	Code        string   `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Description string   `gorm:"size:255" json:"description"`
	Enabled     bool     `gorm:"default:true" json:"enabled"`
	Permissions []string `gorm:"-" json:"permissions"`
}

// TableName 指定表名
func (AdminRole) TableName() string {
	return "admin_roles"
}

// AdminRolePermission 角色权限关联表
type AdminRolePermission struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	RoleID     uint   `gorm:"not null;index" json:"role_id"`
	Permission string `gorm:"size:128;not null" json:"permission"`
}

// TableName 指定表名
func (AdminRolePermission) TableName() string {
	return "admin_role_permissions"
}

// AdminUserRole 用户角色关联表
type AdminUserRole struct {
	ID     uint `gorm:"primarykey" json:"id"`
	UserID uint `gorm:"not null;index:idx_user_role,unique" json:"user_id"`
	RoleID uint `gorm:"not null;index:idx_user_role,unique" json:"role_id"`
}

// TableName 指定表名
func (AdminUserRole) TableName() string {
	return "admin_user_roles"
}
