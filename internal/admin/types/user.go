package types

import (
	"time"

	"bico-admin/internal/shared/types"
)

// StatusRequest 状态更新请求 (别名)
type StatusRequest = types.StatusRequest

// UserListRequest 用户列表请求
type UserListRequest struct {
	types.BasePageQuery
	UserType string `form:"user_type" json:"user_type"`
	Status   *int   `form:"status" json:"status"`
	Email    string `form:"email" json:"email"`
	Phone    string `form:"phone" json:"phone"`
}

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Username string     `json:"username" binding:"required,min=3,max=50"`
	Password string     `json:"password" binding:"required,min=6,max=100"`
	Nickname string     `json:"nickname" binding:"max=100"`
	Email    string     `json:"email" binding:"omitempty,email,max=100"`
	Phone    string     `json:"phone" binding:"omitempty,max=20"`
	Avatar   string     `json:"avatar" binding:"omitempty,max=255"`
	Gender   int        `json:"gender" binding:"oneof=0 1 2"`
	Birthday *time.Time `json:"birthday"`
	UserType string     `json:"user_type" binding:"required,oneof=admin master normal"`
	Status   int        `json:"status" binding:"oneof=0 1"`
	Remark   string     `json:"remark" binding:"max=500"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	Nickname string     `json:"nickname" binding:"max=100"`
	Email    string     `json:"email" binding:"omitempty,email,max=100"`
	Phone    string     `json:"phone" binding:"omitempty,max=20"`
	Avatar   string     `json:"avatar" binding:"omitempty,max=255"`
	Gender   int        `json:"gender" binding:"oneof=0 1 2"`
	Birthday *time.Time `json:"birthday"`
	UserType string     `json:"user_type" binding:"oneof=admin master normal"`
	Status   int        `json:"status" binding:"oneof=0 1"`
	Remark   string     `json:"remark" binding:"max=500"`
}

// UserPasswordRequest 重置用户密码请求
type UserPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6,max=100"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID           uint       `json:"id"`
	Username     string     `json:"username"`
	Nickname     string     `json:"nickname"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone"`
	Avatar       string     `json:"avatar"`
	Gender       int        `json:"gender"`
	GenderText   string     `json:"gender_text"`
	Birthday     *time.Time `json:"birthday"`
	UserType     string     `json:"user_type"`
	UserTypeText string     `json:"user_type_text"`
	Status       int        `json:"status"`
	StatusText   string     `json:"status_text"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	LastLoginIP  string     `json:"last_login_ip"`
	LoginCount   int        `json:"login_count"`
	Remark       string     `json:"remark"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// UserStatsResponse 用户统计响应
type UserStatsResponse struct {
	TotalUsers   int64 `json:"total_users"`
	ActiveUsers  int64 `json:"active_users"`
	AdminUsers   int64 `json:"admin_users"`
	MasterUsers  int64 `json:"master_users"`
	NormalUsers  int64 `json:"normal_users"`
	TodayLogins  int64 `json:"today_logins"`
	WeeklyLogins int64 `json:"weekly_logins"`
}
