package service

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"bico-admin/internal/admin/types"
	sharedTypes "bico-admin/internal/shared/types"
	"bico-admin/pkg/cache"
	"bico-admin/pkg/config"
	"bico-admin/pkg/jwt"
	"bico-admin/pkg/logger"
	"bico-admin/pkg/utils"
)

// AuthService 认证服务接口
type AuthService interface {
	Login(ctx context.Context, req *types.AdminLoginRequest) (*types.AdminLoginResponse, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, req *types.RefreshTokenRequest) (*types.AdminLoginResponse, error)

	GetProfileWithPermissions(ctx context.Context, userID uint) (*types.AdminProfileResponse, error)
	UpdateProfileInfo(ctx context.Context, userID uint, req *types.ProfileUpdateRequest) (*types.AdminUserResponse, error)
	ChangePassword(ctx context.Context, userID uint, req *types.ChangePasswordRequest) error
}

// authService 认证服务实现
type authService struct {
	adminUserService AdminUserService
	jwtManager       *jwt.JWTManager
}

// NewAuthService 创建认证服务
func NewAuthService(adminUserService AdminUserService, cache cache.Cache) AuthService {
	cfg := config.Get()
	jwtManager := jwt.NewJWTManagerWithCache(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.ExpireTime, cache)

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

	// 更新最后登录时间
	if err := s.adminUserService.UpdateLastLoginTime(ctx, adminUser.ID); err != nil {
		logger.Error("更新最后登录时间失败", zap.Error(err))
		// 不影响登录流程，只记录错误
	}

	// 生成JWT令牌
	token, expiresAt, err := s.jwtManager.GenerateToken(adminUser.ID, adminUser.Username, sharedTypes.UserTypeAdmin)
	if err != nil {
		logger.Error("生成令牌失败", zap.Error(err))
		return nil, errors.New("登录失败")
	}

	// 获取用户权限（超级管理员会自动获得所有权限）
	permissions := adminUser.GetPermissionCodes()

	return &types.AdminLoginResponse{
		LoginResponse: sharedTypes.LoginResponse{
			Token:     token,
			ExpiresAt: expiresAt,
			UserInfo:  adminUser.ToUserInfo(),
		},
		Permissions: permissions,
	}, nil
}

// Logout 登出
func (s *authService) Logout(ctx context.Context, token string) error {
	// 将令牌加入黑名单
	if err := s.jwtManager.AddToBlacklist(ctx, token); err != nil {
		logger.Error("将令牌加入黑名单失败",
			zap.String("token", token),
			zap.Error(err))
		return errors.New("登出失败")
	}

	logger.Info("用户登出成功", zap.String("token", token))
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

	// 获取用户权限（超级管理员会自动获得所有权限）
	permissions := adminUser.GetPermissionCodes()

	return &types.AdminLoginResponse{
		LoginResponse: sharedTypes.LoginResponse{
			Token:     token,
			ExpiresAt: expiresAt,
			UserInfo:  adminUser.ToUserInfo(),
		},
		Permissions: permissions,
	}, nil
}

// GetProfileWithPermissions 获取用户资料和权限
func (s *authService) GetProfileWithPermissions(ctx context.Context, userID uint) (*types.AdminProfileResponse, error) {
	adminUser, err := s.adminUserService.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 获取用户权限（超级管理员会自动获得所有权限）
	permissions := adminUser.GetPermissionCodes()

	// 检查权限
	canDelete, _ := s.adminUserService.CanUserBeDeleted(ctx, adminUser.ID)
	canDisable, _ := s.adminUserService.CanUserBeDisabled(ctx, adminUser.ID)

	// 转换LastLoginAt
	var lastLoginAt *utils.FormattedTime
	if adminUser.LastLoginAt != nil {
		ft := utils.NewFormattedTime(*adminUser.LastLoginAt)
		lastLoginAt = &ft
	}

	status := 0
	if adminUser.Status != nil {
		status = *adminUser.Status
	}

	userInfo := types.AdminUserResponse{
		ID:          adminUser.ID,
		Username:    adminUser.Username,
		Name:        adminUser.Name,
		Avatar:      adminUser.Avatar,
		Email:       adminUser.Email,
		Phone:       adminUser.Phone,
		Status:      status,
		StatusText:  adminUser.GetStatusText(),
		LastLoginAt: lastLoginAt,
		Remark:      adminUser.Remark,
		CanDelete:   canDelete,
		CanDisable:  canDisable,
		CreatedAt:   utils.NewFormattedTime(adminUser.CreatedAt),
		UpdatedAt:   utils.NewFormattedTime(adminUser.UpdatedAt),
	}

	return &types.AdminProfileResponse{
		UserInfo:    userInfo,
		Permissions: permissions,
	}, nil
}

// UpdateProfileInfo 更新个人信息
func (s *authService) UpdateProfileInfo(ctx context.Context, userID uint, req *types.ProfileUpdateRequest) (*types.AdminUserResponse, error) {
	// 获取当前用户信息
	adminUser, err := s.adminUserService.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 更新用户信息
	adminUser.Name = req.Name
	adminUser.Avatar = req.Avatar
	adminUser.Email = req.Email
	adminUser.Phone = req.Phone

	// 保存更新
	updatedUser, err := s.adminUserService.UpdateProfileInfo(ctx, adminUser)
	if err != nil {
		return nil, err
	}

	// 检查权限
	canDelete, _ := s.adminUserService.CanUserBeDeleted(ctx, updatedUser.ID)
	canDisable, _ := s.adminUserService.CanUserBeDisabled(ctx, updatedUser.ID)

	// 转换LastLoginAt
	var lastLoginAt *utils.FormattedTime
	if updatedUser.LastLoginAt != nil {
		ft := utils.NewFormattedTime(*updatedUser.LastLoginAt)
		lastLoginAt = &ft
	}

	status := 0
	if updatedUser.Status != nil {
		status = *updatedUser.Status
	}

	return &types.AdminUserResponse{
		ID:          updatedUser.ID,
		Username:    updatedUser.Username,
		Name:        updatedUser.Name,
		Avatar:      updatedUser.Avatar,
		Email:       updatedUser.Email,
		Phone:       updatedUser.Phone,
		Status:      status,
		StatusText:  updatedUser.GetStatusText(),
		LastLoginAt: lastLoginAt,
		Remark:      updatedUser.Remark,
		CanDelete:   canDelete,
		CanDisable:  canDisable,
		CreatedAt:   utils.NewFormattedTime(updatedUser.CreatedAt),
		UpdatedAt:   utils.NewFormattedTime(updatedUser.UpdatedAt),
	}, nil
}

// ChangePassword 修改密码
func (s *authService) ChangePassword(ctx context.Context, userID uint, req *types.ChangePasswordRequest) error {
	// 获取用户信息
	adminUser, err := s.adminUserService.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 验证原密码
	if err := bcrypt.CompareHashAndPassword([]byte(adminUser.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	// 检查新密码是否与原密码相同
	if req.OldPassword == req.NewPassword {
		return errors.New("新密码不能与原密码相同")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	return s.adminUserService.UpdatePassword(ctx, userID, string(hashedPassword))
}
