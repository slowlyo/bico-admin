package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"bico-admin/internal/admin/model"
	"bico-admin/internal/core/cache"
	"bico-admin/internal/pkg/crud"
	"bico-admin/internal/pkg/jwt"
	"bico-admin/internal/pkg/password"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound     = errors.New("用户不存在")
	ErrInvalidPassword  = errors.New("密码错误")
	ErrUserDisabled     = errors.New("用户已被禁用")
	ErrOldPasswordWrong = errors.New("原密码错误")
)

const (
	permissionCacheTTL = 5 * time.Minute
	userStatusCacheTTL = 1 * time.Minute
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID          uint     `json:"id"`
	Username    string   `json:"username"`
	Name        string   `json:"name"`
	Avatar      string   `json:"avatar"`
	Enabled     bool     `json:"enabled"`
	Permissions []string `json:"permissions"`
}

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// IAuthService 认证服务接口
type IAuthService interface {
	Login(req *LoginRequest) (*LoginResponse, error)
	Logout(token string) error
	IsTokenBlacklisted(token string) bool
	GetUserByID(userID uint) (*UserInfo, error)
	UpdateProfile(userID uint, req *UpdateProfileRequest) (*UserInfo, error)
	ChangePassword(userID uint, req *ChangePasswordRequest) error
	GetUserPermissions(userID uint) ([]string, error)
	IsUserEnabled(userID uint) (bool, error)
}

// AuthCacheInvalidator 认证缓存失效接口
type AuthCacheInvalidator interface {
	InvalidateUserPermissionCache(userID uint)
	InvalidateRoleUsersPermissionCache(roleID uint)
	InvalidateUserStatusCache(userID uint)
}

// AuthService 认证服务
type AuthService struct {
	db         *gorm.DB
	jwtManager *jwt.JWTManager
	cache      cache.Cache
}

// NewAuthService 创建认证服务
func NewAuthService(db *gorm.DB, jwtManager *jwt.JWTManager, cache cache.Cache) *AuthService {
	return &AuthService{
		db:         db,
		jwtManager: jwtManager,
		cache:      cache,
	}
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	var user model.AdminUser

	err := s.db.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.Enabled {
		return nil, ErrUserDisabled
	}

	if !password.Verify(user.Password, req.Password) {
		return nil, ErrInvalidPassword
	}

	token, err := s.jwtManager.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{Token: token}, nil
}

// Logout 用户退出登录
func (s *AuthService) Logout(token string) error {
	if token == "" {
		return nil
	}

	blacklistKey := "token:blacklist:" + token
	// 设置 token 黑名单，过期时间 7 天
	return s.cache.Set(blacklistKey, true, 7*24*time.Hour)
}

// IsTokenBlacklisted 检查 token 是否在黑名单中
func (s *AuthService) IsTokenBlacklisted(token string) bool {
	blacklistKey := "token:blacklist:" + token
	return s.cache.Exists(blacklistKey)
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

	permissions, err := s.GetUserPermissions(userID)
	// 权限读取失败时直接返回，避免返回不完整用户信息。
	if err != nil {
		return nil, err
	}

	// 命中用户查询后同步刷新状态缓存，减少后续中间件查库频率。
	s.setUserStatusCache(userID, user.Enabled)

	return &UserInfo{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Avatar:      user.Avatar,
		Enabled:     user.Enabled,
		Permissions: permissions,
	}, nil
}

// UpdateProfile 更新用户资料
func (s *AuthService) UpdateProfile(userID uint, req *UpdateProfileRequest) (*UserInfo, error) {
	var user model.AdminUser

	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.Enabled {
		return nil, ErrUserDisabled
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if len(updates) > 0 {
		err = s.db.Model(&user).Updates(updates).Error
		if err != nil {
			return nil, err
		}
	}

	permissions, err := s.GetUserPermissions(userID)
	// 权限读取失败时返回错误，避免返回不完整数据。
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Avatar:      user.Avatar,
		Enabled:     user.Enabled,
		Permissions: permissions,
	}, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	var user model.AdminUser

	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return ErrUserNotFound
	}

	if !user.Enabled {
		return ErrUserDisabled
	}

	if !password.Verify(user.Password, req.OldPassword) {
		return ErrOldPasswordWrong
	}

	hashedPassword, err := password.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	err = s.db.Model(&user).Update("password", hashedPassword).Error
	if err != nil {
		return err
	}

	return nil
}

// GetUserPermissions 获取用户的所有权限
func (s *AuthService) GetUserPermissions(userID uint) ([]string, error) {
	// 先读缓存，命中则直接返回，降低鉴权链路查库频率。
	if cachedPerms, ok := s.getPermissionsCache(userID); ok {
		return cachedPerms, nil
	}

	// 获取用户信息
	var user model.AdminUser
	if err := s.db.Select("username").First(&user, userID).Error; err != nil {
		return nil, err
	}

	// 如果是默认管理员账户，返回所有权限
	if user.Username == "admin" {
		adminPerms := crud.GetAllPermissionKeys()
		s.setPermissionsCache(userID, adminPerms)
		return adminPerms, nil
	}

	// 从数据库查询普通用户权限
	var permissions []string
	err := s.db.Table("admin_user_roles").
		Select("DISTINCT admin_role_permissions.permission").
		Joins("JOIN admin_role_permissions ON admin_user_roles.role_id = admin_role_permissions.role_id").
		Joins("JOIN admin_roles ON admin_role_permissions.role_id = admin_roles.id").
		Where("admin_user_roles.user_id = ? AND admin_roles.enabled = ?", userID, true).
		Pluck("permission", &permissions).Error

	if err != nil {
		return nil, err
	}

	s.setPermissionsCache(userID, permissions)
	return permissions, nil
}

// IsUserEnabled 获取用户启用状态（优先读取缓存）
func (s *AuthService) IsUserEnabled(userID uint) (bool, error) {
	if cachedEnabled, ok := s.getUserStatusCache(userID); ok {
		return cachedEnabled, nil
	}

	var user model.AdminUser
	if err := s.db.Select("enabled").First(&user, userID).Error; err != nil {
		// 用户不存在时返回统一业务错误。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, ErrUserNotFound
		}
		return false, err
	}

	s.setUserStatusCache(userID, user.Enabled)
	return user.Enabled, nil
}

// InvalidateUserPermissionCache 失效指定用户权限缓存
func (s *AuthService) InvalidateUserPermissionCache(userID uint) {
	_ = s.cache.Delete(permissionCacheKey(userID))
}

// InvalidateRoleUsersPermissionCache 失效指定角色下所有用户的权限缓存
func (s *AuthService) InvalidateRoleUsersPermissionCache(roleID uint) {
	var userIDs []uint
	if err := s.db.Table("admin_user_roles").
		Where("role_id = ?", roleID).
		Distinct("user_id").
		Pluck("user_id", &userIDs).Error; err != nil {
		return
	}
	for _, userID := range userIDs {
		_ = s.cache.Delete(permissionCacheKey(userID))
	}
}

// InvalidateUserStatusCache 失效指定用户状态缓存
func (s *AuthService) InvalidateUserStatusCache(userID uint) {
	_ = s.cache.Delete(userStatusCacheKey(userID))
}

// getPermissionsCache 获取用户权限缓存
func (s *AuthService) getPermissionsCache(userID uint) ([]string, bool) {
	value, err := s.cache.Get(permissionCacheKey(userID))
	if err != nil {
		return nil, false
	}
	return parseStringSlice(value)
}

// setPermissionsCache 写入用户权限缓存
func (s *AuthService) setPermissionsCache(userID uint, permissions []string) {
	_ = s.cache.Set(permissionCacheKey(userID), permissions, permissionCacheTTL)
}

// getUserStatusCache 获取用户状态缓存
func (s *AuthService) getUserStatusCache(userID uint) (bool, bool) {
	value, err := s.cache.Get(userStatusCacheKey(userID))
	if err != nil {
		return false, false
	}
	switch v := value.(type) {
	case bool:
		return v, true
	case string:
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return false, false
		}
		return parsed, true
	case float64:
		// 某些序列化场景可能将布尔值落成 0/1。
		return v != 0, true
	default:
		return false, false
	}
}

// setUserStatusCache 写入用户状态缓存
func (s *AuthService) setUserStatusCache(userID uint, enabled bool) {
	_ = s.cache.Set(userStatusCacheKey(userID), enabled, userStatusCacheTTL)
}

// permissionCacheKey 生成权限缓存 key
func permissionCacheKey(userID uint) string {
	return fmt.Sprintf("auth:user:%d:permissions", userID)
}

// userStatusCacheKey 生成用户状态缓存 key
func userStatusCacheKey(userID uint) string {
	return fmt.Sprintf("auth:user:%d:enabled", userID)
}

// parseStringSlice 将缓存值安全转换为字符串数组
func parseStringSlice(value interface{}) ([]string, bool) {
	switch v := value.(type) {
	case []string:
		return v, true
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			str, ok := item.(string)
			// 缓存内容出现非字符串时认为数据损坏，直接回源查询。
			if !ok {
				return nil, false
			}
			result = append(result, str)
		}
		return result, true
	default:
		return nil, false
	}
}
