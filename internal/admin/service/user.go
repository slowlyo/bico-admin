package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"bico-admin/internal/admin/repository"
	"bico-admin/internal/admin/types"
	"bico-admin/internal/shared/model"
	sharedTypes "bico-admin/internal/shared/types"
)

// UserService 用户服务接口
type UserService interface {
	GetList(ctx context.Context, req *types.UserListRequest) (*sharedTypes.PageResult, error)
	GetByID(ctx context.Context, id uint) (*types.UserResponse, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Create(ctx context.Context, req *types.UserCreateRequest) (*types.UserResponse, error)
	Update(ctx context.Context, id uint, req *types.UserUpdateRequest) (*types.UserResponse, error)
	Delete(ctx context.Context, id uint) error
	UpdateStatus(ctx context.Context, id uint, status int) error
	ResetPassword(ctx context.Context, id uint, password string) error
	UpdateLoginInfo(ctx context.Context, id uint, ip string) error
	GetStats(ctx context.Context) (*types.UserStatsResponse, error)
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetList 获取用户列表
func (s *userService) GetList(ctx context.Context, req *types.UserListRequest) (*sharedTypes.PageResult, error) {
	users, total, err := s.userRepo.GetList(ctx, req)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	userResponses := make([]*types.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = s.toUserResponse(user)
	}

	return sharedTypes.NewPageResult(userResponses, total, req.Page, req.GetPageSize()), nil
}

// GetByID 根据ID获取用户
func (s *userService) GetByID(ctx context.Context, id uint) (*types.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// GetByUsername 根据用户名获取用户
func (s *userService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

// Create 创建用户
func (s *userService) Create(ctx context.Context, req *types.UserCreateRequest) (*types.UserResponse, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(ctx, req.Username); err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Gender:   req.Gender,
		Birthday: req.Birthday,
		UserType: req.UserType,
		Status:   req.Status,
		Remark:   req.Remark,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// Update 更新用户
func (s *userService) Update(ctx context.Context, id uint, req *types.UserUpdateRequest) (*types.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Phone = req.Phone
	user.Avatar = req.Avatar
	user.Gender = req.Gender
	user.Birthday = req.Birthday
	user.UserType = req.UserType
	user.Status = req.Status
	user.Remark = req.Remark

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// Delete 删除用户
func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}

// UpdateStatus 更新用户状态
func (s *userService) UpdateStatus(ctx context.Context, id uint, status int) error {
	return s.userRepo.UpdateStatus(ctx, id, status)
}

// ResetPassword 重置密码
func (s *userService) ResetPassword(ctx context.Context, id uint, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	return s.userRepo.UpdatePassword(ctx, id, string(hashedPassword))
}

// UpdateLoginInfo 更新登录信息
func (s *userService) UpdateLoginInfo(ctx context.Context, id uint, ip string) error {
	return s.userRepo.UpdateLoginInfo(ctx, id, ip)
}

// GetStats 获取用户统计
func (s *userService) GetStats(ctx context.Context) (*types.UserStatsResponse, error) {
	return s.userRepo.GetStats(ctx)
}

// toUserResponse 转换为用户响应格式
func (s *userService) toUserResponse(user *model.User) *types.UserResponse {
	return &types.UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Nickname:     user.Nickname,
		Email:        user.Email,
		Phone:        user.Phone,
		Avatar:       user.Avatar,
		Gender:       user.Gender,
		GenderText:   user.GetGenderText(),
		Birthday:     user.Birthday,
		UserType:     user.UserType,
		UserTypeText: user.GetUserTypeText(),
		Status:       user.Status,
		StatusText:   user.GetStatusText(),
		LastLoginAt:  user.LastLoginAt,
		LastLoginIP:  user.LastLoginIP,
		LoginCount:   user.LoginCount,
		Remark:       user.Remark,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
