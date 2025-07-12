package model

import (
	"time"

	"gorm.io/gorm"

	"bico-admin/internal/shared/types"
)

// AdminUser 管理员用户模型
type AdminUser struct {
	ID        uint      `json:"id" gorm:"primarykey;autoIncrement;comment:主键ID"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;comment:更新时间"`
	Username  string    `json:"username" gorm:"uniqueIndex:uni_admin_users_username;size:255;not null;comment:用户名"`
	Password  string    `json:"-" gorm:"size:255;not null;comment:密码"`
	Name      string    `json:"name" gorm:"size:255;not null;comment:姓名"`
	Avatar    string    `json:"avatar" gorm:"size:255;not null;comment:头像"`
	Enabled   bool      `json:"enabled" gorm:"default:true;not null;comment:是否启用"`
}

// TableName 表名
func (AdminUser) TableName() string {
	return "admin_users"
}

// IsEnabled 是否启用
func (u *AdminUser) IsEnabled() bool {
	return u.Enabled
}

// ToUserInfo 转换为用户信息
func (u *AdminUser) ToUserInfo() types.UserInfo {
	status := types.StatusActive
	if !u.Enabled {
		status = types.StatusInactive
	}

	return types.UserInfo{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Name,
		Email:    "", // AdminUser没有email字段
		Avatar:   u.Avatar,
		UserType: types.UserTypeAdmin,
		Status:   status,
	}
}

// GetStatusText 获取状态文本
func (u *AdminUser) GetStatusText() string {
	if u.Enabled {
		return "启用"
	}
	return "禁用"
}

// BeforeCreate GORM钩子：创建前
func (u *AdminUser) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑，比如密码加密等
	return nil
}

// BeforeUpdate GORM钩子：更新前
func (u *AdminUser) BeforeUpdate(tx *gorm.DB) error {
	// 可以在这里添加更新前的逻辑
	return nil
}
