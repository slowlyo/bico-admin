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
