package service

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"bico-admin/internal/admin/definitions"
	"bico-admin/internal/admin/types"
	sharedTypes "bico-admin/internal/shared/types"
	"bico-admin/pkg/config"
	"bico-admin/pkg/jwt"
	"bico-admin/pkg/logger"
)

// AuthService 认证服务接口
type AuthService interface {
	Login(ctx context.Context, req *types.AdminLoginRequest) (*types.AdminLoginResponse, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, req *types.RefreshTokenRequest) (*types.AdminLoginResponse, error)
	GetProfile(ctx context.Context, userID uint) (*types.AdminUserResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req *types.AdminUserUpdateRequest) (*types.AdminUserResponse, error)
}

// authService 认证服务实现
type authService struct {
	adminUserService AdminUserService
	jwtManager       *jwt.JWTManager
}

// NewAuthService 创建认证服务
func NewAuthService(adminUserService AdminUserService) AuthService {
	cfg := config.Get()
	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.ExpireTime)

	return &authService{
		adminUserService: adminUserService,
		jwtManager:       jwtManager,
	}
}

// Login 管理员登录
func (s *authService) Login(ctx context.Context, req *types.AdminLoginRequest) (*types.AdminLoginResponse, error) {
	// TODO: 验证验证码
	if req.Captcha != "1234" { // 临时验证码验证
		return nil, errors.New("验证码错误")
	}

	// 查找管理员用户
	adminUser, err := s.adminUserService.GetByUsername(ctx, req.Username)
	if err != nil {
		logger.Error("管理员登录失败",
			zap.String("username", req.Username),
			zap.Error(err))
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(adminUser.Password), []byte(req.Password)); err != nil {
		logger.Error("密码验证失败",
			zap.String("username", req.Username))
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if !adminUser.IsEnabled() {
		return nil, errors.New("用户已被禁用")
	}

	// 生成JWT令牌
	token, expiresAt, err := s.jwtManager.GenerateToken(adminUser.ID, adminUser.Username, sharedTypes.UserTypeAdmin)
	if err != nil {
		logger.Error("生成令牌失败", zap.Error(err))
		return nil, errors.New("登录失败")
	}

	// 获取用户权限和菜单
	permissions := s.getAdminPermissions()
	menus := s.getAdminMenus()

	return &types.AdminLoginResponse{
		LoginResponse: sharedTypes.LoginResponse{
			Token:     token,
			ExpiresAt: expiresAt,
			UserInfo:  adminUser.ToUserInfo(),
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
	// 验证并刷新令牌
	token, expiresAt, err := s.jwtManager.RefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("刷新令牌失败: " + err.Error())
	}

	// 解析令牌获取用户信息
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, errors.New("令牌解析失败")
	}

	// 获取管理员用户信息
	adminUser, err := s.adminUserService.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户状态
	if !adminUser.IsEnabled() {
		return nil, errors.New("用户已被禁用")
	}

	// 获取用户权限和菜单
	permissions := s.getAdminPermissions()
	menus := s.getAdminMenus()

	return &types.AdminLoginResponse{
		LoginResponse: sharedTypes.LoginResponse{
			Token:     token,
			ExpiresAt: expiresAt,
			UserInfo:  adminUser.ToUserInfo(),
		},
		Permissions: permissions,
		Menus:       menus,
	}, nil
}

// GetProfile 获取用户资料
func (s *authService) GetProfile(ctx context.Context, userID uint) (*types.AdminUserResponse, error) {
	adminUser, err := s.adminUserService.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &types.AdminUserResponse{
		ID:          adminUser.ID,
		Username:    adminUser.Username,
		Name:        adminUser.Name,
		Avatar:      adminUser.Avatar,
		Email:       adminUser.Email,
		Phone:       adminUser.Phone,
		Status:      adminUser.Status,
		StatusText:  adminUser.GetStatusText(),
		LastLoginAt: adminUser.LastLoginAt,
		LastLoginIP: adminUser.LastLoginIP,
		LoginCount:  adminUser.LoginCount,
		Remark:      adminUser.Remark,
		CreatedAt:   adminUser.CreatedAt,
		UpdatedAt:   adminUser.UpdatedAt,
	}, nil
}

// UpdateProfile 更新用户资料
func (s *authService) UpdateProfile(ctx context.Context, userID uint, req *types.AdminUserUpdateRequest) (*types.AdminUserResponse, error) {
	adminUser, err := s.adminUserService.Update(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	return &types.AdminUserResponse{
		ID:          adminUser.ID,
		Username:    adminUser.Username,
		Name:        adminUser.Name,
		Avatar:      adminUser.Avatar,
		Email:       adminUser.Email,
		Phone:       adminUser.Phone,
		Status:      adminUser.Status,
		StatusText:  adminUser.GetStatusText(),
		LastLoginAt: adminUser.LastLoginAt,
		LastLoginIP: adminUser.LastLoginIP,
		LoginCount:  adminUser.LoginCount,
		Remark:      adminUser.Remark,
		CreatedAt:   adminUser.CreatedAt,
		UpdatedAt:   adminUser.UpdatedAt,
	}, nil
}

// getAdminPermissions 获取管理员权限
func (s *authService) getAdminPermissions() []string {
	// TODO: 根据用户角色或具体权限配置返回权限
	// 目前返回所有权限，后续可以根据用户角色进行过滤
	return definitions.GetPermissionCodes()
}

// getAdminMenus 获取管理员菜单
func (s *authService) getAdminMenus() []types.Menu {
	// 获取用户权限
	userPermissions := s.getAdminPermissions()

	// 根据权限过滤菜单
	return definitions.FilterMenusByPermissions(userPermissions)
}
