package types

import "time"

// 通用状态枚举
const (
	StatusActive   = 1  // 激活
	StatusInactive = 0  // 未激活
	StatusDeleted  = -1 // 已删除
)

// 用户类型枚举
const (
	UserTypeAdmin  = "admin"  // 管理员
	UserTypeMaster = "master" // 主控用户
	UserTypeNormal = "normal" // 普通用户
)

// 性别枚举
const (
	GenderUnknown = 0 // 未知
	GenderMale    = 1 // 男性
	GenderFemale  = 2 // 女性
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IDRequest 通用ID请求
type IDRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

// StatusRequest 状态更新请求
type StatusRequest struct {
	Status int `json:"status" binding:"oneof=0 1 -1"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserInfo  UserInfo  `json:"user_info"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	UserType string `json:"user_type"`
	Status   int    `json:"status"`
}

// GetStatusText 获取状态文本
func GetStatusText(status int) string {
	switch status {
	case StatusActive:
		return "激活"
	case StatusInactive:
		return "未激活"
	case StatusDeleted:
		return "已删除"
	default:
		return "未知"
	}
}

// GetGenderText 获取性别文本
func GetGenderText(gender int) string {
	switch gender {
	case GenderMale:
		return "男"
	case GenderFemale:
		return "女"
	default:
		return "未知"
	}
}

// GetUserTypeText 获取用户类型文本
func GetUserTypeText(userType string) string {
	switch userType {
	case UserTypeAdmin:
		return "管理员"
	case UserTypeMaster:
		return "主控用户"
	case UserTypeNormal:
		return "普通用户"
	default:
		return "未知"
	}
}
