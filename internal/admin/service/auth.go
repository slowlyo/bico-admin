package service

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"bico-admin/internal/admin/types"
	"bico-admin/internal/shared/model"
	sharedTypes "bico-admin/internal/shared/types"
	"bico-admin/pkg/logger"
)

// AuthService 认证服务接口
type AuthService interface {
	Login(ctx context.Context, req *types.AdminLoginRequest) (*types.AdminLoginResponse, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, req *types.RefreshTokenRequest) (*types.AdminLoginResponse, error)
	GetProfile(ctx context.Context, userID uint) (*types.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req *types.UserUpdateRequest) (*types.UserResponse, error)
}

// authService 认证服务实现
type authService struct {
	userService UserService
}

// NewAuthService 创建认证服务
func NewAuthService(userService UserService) AuthService {
	return &authService{
		userService: userService,
	}
}

// Login 管理员登录
func (s *authService) Login(ctx context.Context, req *types.AdminLoginRequest) (*types.AdminLoginResponse, error) {
	// TODO: 验证验证码
	if req.Captcha != "1234" { // 临时验证码验证
		return nil, errors.New("验证码错误")
	}

	// 查找用户
	user, err := s.userService.GetByUsername(ctx, req.Username)
	if err != nil {
		logger.Error("用户登录失败",
			zap.String("username", req.Username),
			zap.Error(err))
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Error("密码验证失败",
			zap.String("username", req.Username))
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if !user.IsActive() {
		return nil, errors.New("用户已被禁用")
	}

	// 检查用户类型（只允许管理员和主控用户登录）
	if !user.IsAdmin() && !user.IsMaster() {
		return nil, errors.New("权限不足")
	}

	// 生成JWT令牌
	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		logger.Error("生成令牌失败", zap.Error(err))
		return nil, errors.New("登录失败")
	}

	// 更新登录信息
	if err := s.userService.UpdateLoginInfo(ctx, user.ID, "127.0.0.1"); err != nil {
		logger.Error("更新登录信息失败", zap.Error(err))
	}

	// 获取用户权限和菜单
	permissions := s.getUserPermissions(user)
	menus := s.getUserMenus(user)

	return &types.AdminLoginResponse{
		LoginResponse: sharedTypes.LoginResponse{
			Token:     token,
			ExpiresAt: expiresAt,
			UserInfo:  user.ToUserInfo(),
		},
		Permissions: permissions,
		Menus:       menus,
	}, nil
}

// Logout 登出
func (s *authService) Logout(ctx context.Context, token string) error {
	// TODO: 将令牌加入黑名单
	logger.Info("用户登出", zap.String("token", token))
	return nil
}

// RefreshToken 刷新令牌
func (s *authService) RefreshToken(ctx context.Context, req *types.RefreshTokenRequest) (*types.AdminLoginResponse, error) {
	// TODO: 验证刷新令牌并生成新的访问令牌
	return nil, errors.New("功能暂未实现")
}

// GetProfile 获取用户资料
func (s *authService) GetProfile(ctx context.Context, userID uint) (*types.UserResponse, error) {
	return s.userService.GetByID(ctx, userID)
}

// UpdateProfile 更新用户资料
func (s *authService) UpdateProfile(ctx context.Context, userID uint, req *types.UserUpdateRequest) (*types.UserResponse, error) {
	return s.userService.Update(ctx, userID, req)
}

// generateToken 生成JWT令牌
func (s *authService) generateToken(user *model.User) (string, time.Time, error) {
	// TODO: 实现JWT令牌生成
	expiresAt := time.Now().Add(24 * time.Hour)
	token := "mock_token_" + user.Username // 临时mock令牌
	return token, expiresAt, nil
}

// getUserPermissions 获取用户权限
func (s *authService) getUserPermissions(user *model.User) []string {
	permissions := []string{}

	if user.IsAdmin() {
		permissions = append(permissions,
			"user:read", "user:write", "user:delete",
			"system:read", "system:write",
			"config:read", "config:write",
		)
	}

	if user.IsMaster() {
		permissions = append(permissions,
			"user:read", "user:write",
			"system:read",
		)
	}

	return permissions
}

// getUserMenus 获取用户菜单
func (s *authService) getUserMenus(user *model.User) []types.Menu {
	menus := []types.Menu{}

	if user.IsAdmin() || user.IsMaster() {
		menus = append(menus, types.Menu{
			ID:   1,
			Name: "用户管理",
			Path: "/admin/users",
			Icon: "user",
			Sort: 1,
		})
	}

	if user.IsAdmin() {
		menus = append(menus, types.Menu{
			ID:   2,
			Name: "系统管理",
			Path: "/admin/system",
			Icon: "setting",
			Sort: 2,
		})
	}

	return menus
}
