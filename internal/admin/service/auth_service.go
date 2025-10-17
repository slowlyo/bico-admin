package service

import (
	"errors"
	"time"
	
	"bico-admin/internal/admin/model"
	"bico-admin/internal/shared/password"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("用户不存在")
	ErrInvalidPassword   = errors.New("密码错误")
	ErrUserDisabled      = errors.New("用户已被禁用")
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  *UserInfo   `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Enabled  bool   `json:"enabled"`
}

// IAuthService 认证服务接口
type IAuthService interface {
	Login(req interface{}) (interface{}, error)
	Logout(token string) error
	IsTokenBlacklisted(token string) bool
	GetUserByID(userID uint) (*UserInfo, error)
}

// AuthService 认证服务
type AuthService struct {
	db         *gorm.DB
	jwtManager interface{}
	cache      interface{}
}

// NewAuthService 创建认证服务
func NewAuthService(db *gorm.DB, jwtManager interface{}, cache interface{}) *AuthService {
	return &AuthService{
		db:         db,
		jwtManager: jwtManager,
		cache:      cache,
	}
}

// Login 用户登录
func (s *AuthService) Login(req interface{}) (interface{}, error) {
	loginReq := req.(*struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	})
	
	var user model.AdminUser
	
	err := s.db.Where("username = ?", loginReq.Username).First(&user).Error
	if err != nil {
		return nil, ErrUserNotFound
	}
	
	if !user.Enabled {
		return nil, ErrUserDisabled
	}
	
	if !password.Verify(user.Password, loginReq.Password) {
		return nil, ErrInvalidPassword
	}
	
	jwtMgr := s.jwtManager.(interface {
		GenerateToken(userID uint, username string) (string, error)
	})
	token, err := jwtMgr.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}
	
	return &LoginResponse{
		Token: token,
		User: &UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
			Avatar:   user.Avatar,
			Enabled:  user.Enabled,
		},
	}, nil
}

// Logout 用户退出登录
func (s *AuthService) Logout(token string) error {
	if token == "" {
		return nil
	}
	
	blacklistKey := "token:blacklist:" + token
	// 设置 token 黑名单，过期时间 7 天
	return s.cache.(interface {
		Set(key string, value interface{}, expiration time.Duration) error
	}).Set(blacklistKey, true, 7*24*time.Hour)
}

// IsTokenBlacklisted 检查 token 是否在黑名单中
func (s *AuthService) IsTokenBlacklisted(token string) bool {
	blacklistKey := "token:blacklist:" + token
	return s.cache.(interface {
		Exists(key string) bool
	}).Exists(blacklistKey)
}

// GetUserByID 根据用户ID获取用户信息
func (s *AuthService) GetUserByID(userID uint) (*UserInfo, error) {
	var user model.AdminUser
	
	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, ErrUserNotFound
	}
	
	if !user.Enabled {
		return nil, ErrUserDisabled
	}
	
	return &UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Avatar:   user.Avatar,
		Enabled:  user.Enabled,
	}, nil
}
