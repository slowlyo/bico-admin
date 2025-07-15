package types

import (
	"bico-admin/internal/shared/types"
	"time"
)

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Captcha  string `json:"captcha" binding:"required,len=4"`
}

// AdminLoginResponse 管理员登录响应
type AdminLoginResponse struct {
	types.LoginResponse
	Permissions []string `json:"permissions"`
}

// AdminProfileResponse 管理员资料响应（包含权限）
type AdminProfileResponse struct {
	UserInfo    types.UserInfo `json:"user_info"`
	Permissions []string       `json:"permissions"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest 登出请求
type LogoutRequest struct {
	Token string `json:"token" binding:"required"`
}

// AdminUserCreateRequest 创建管理员用户请求
type AdminUserCreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Avatar   string `json:"avatar" binding:"max=255"`
	Enabled  bool   `json:"enabled"`
}

// AdminUserUpdateRequest 更新管理员用户请求
type AdminUserUpdateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"omitempty,min=6,max=100"` // 可选，为空则不更新密码
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Avatar   string `json:"avatar" binding:"max=255"`
	Enabled  bool   `json:"enabled"`
}

// AdminUserListRequest 管理员用户列表请求
type AdminUserListRequest struct {
	types.BasePageQuery
	Username string `form:"username" json:"username"` // 用户名
	Name     string `form:"name" json:"name"`         // 姓名
	Status   *int   `form:"status" json:"status"`     // 状态
}

// AdminUserResponse 管理员用户响应
type AdminUserResponse struct {
	ID          uint       `json:"id"`
	Username    string     `json:"username"`
	Name        string     `json:"name"`
	Avatar      string     `json:"avatar"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	Status      int        `json:"status"`
	StatusText  string     `json:"status_text"`
	LastLoginAt *time.Time `json:"last_login_at"`
	Remark      string     `json:"remark"`
	CanDelete   bool       `json:"can_delete"`  // 是否可删除
	CanDisable  bool       `json:"can_disable"` // 是否可禁用
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
