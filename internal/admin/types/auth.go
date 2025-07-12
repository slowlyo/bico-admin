package types

import "bico-admin/internal/shared/types"

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
	Menus       []Menu   `json:"menus"`
}

// Menu 菜单结构
type Menu struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Icon     string `json:"icon"`
	Sort     int    `json:"sort"`
	Children []Menu `json:"children,omitempty"`
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

// AdminUserResponse 管理员用户响应
type AdminUserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	Enabled   bool   `json:"enabled"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
