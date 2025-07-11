package model

import (
	"time"

	"bico-admin/internal/shared/types"
)

// User 用户模型
type User struct {
	types.BaseModel
	Username     string     `json:"username" gorm:"uniqueIndex;size:50;not null;comment:用户名"`
	Password     string     `json:"-" gorm:"size:255;not null;comment:密码"`
	Nickname     string     `json:"nickname" gorm:"size:100;comment:昵称"`
	Email        string     `json:"email" gorm:"size:100;comment:邮箱"`
	Phone        string     `json:"phone" gorm:"size:20;comment:手机号"`
	Avatar       string     `json:"avatar" gorm:"size:255;comment:头像"`
	Gender       int        `json:"gender" gorm:"default:0;comment:性别 0:未知 1:男 2:女"`
	Birthday     *time.Time `json:"birthday" gorm:"comment:生日"`
	UserType     string     `json:"user_type" gorm:"size:20;default:'normal';comment:用户类型"`
	Status       int        `json:"status" gorm:"default:1;comment:状态 1:激活 0:未激活 -1:已删除"`
	LastLoginAt  *time.Time `json:"last_login_at" gorm:"comment:最后登录时间"`
	LastLoginIP  string     `json:"last_login_ip" gorm:"size:45;comment:最后登录IP"`
	LoginCount   int        `json:"login_count" gorm:"default:0;comment:登录次数"`
	Remark       string     `json:"remark" gorm:"size:500;comment:备注"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

// IsActive 是否激活
func (u *User) IsActive() bool {
	return u.Status == types.StatusActive
}

// IsAdmin 是否管理员
func (u *User) IsAdmin() bool {
	return u.UserType == types.UserTypeAdmin
}

// IsMaster 是否主控用户
func (u *User) IsMaster() bool {
	return u.UserType == types.UserTypeMaster
}

// GetGenderText 获取性别文本
func (u *User) GetGenderText() string {
	return types.GetGenderText(u.Gender)
}

// GetStatusText 获取状态文本
func (u *User) GetStatusText() string {
	return types.GetStatusText(u.Status)
}

// GetUserTypeText 获取用户类型文本
func (u *User) GetUserTypeText() string {
	return types.GetUserTypeText(u.UserType)
}

// ToUserInfo 转换为用户信息
func (u *User) ToUserInfo() types.UserInfo {
	return types.UserInfo{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Email:    u.Email,
		Avatar:   u.Avatar,
		UserType: u.UserType,
		Status:   u.Status,
	}
}
