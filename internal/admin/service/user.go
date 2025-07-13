package service

import (
	"context"
	"errors"

	"bico-admin/internal/admin/types"
	sharedTypes "bico-admin/internal/shared/types"
)

// UserService 普通用户服务接口
type UserService interface {
	GetByID(ctx context.Context, id uint) (*types.UserResponse, error)
	Create(ctx context.Context, req *types.UserCreateRequest) (*types.UserResponse, error)
	Update(ctx context.Context, id uint, req *types.UserUpdateRequest) (*types.UserResponse, error)
	Delete(ctx context.Context, id uint) error
	UpdateStatus(ctx context.Context, id uint, status int) error
	List(ctx context.Context, req *types.UserListRequest) ([]*types.UserResponse, int64, error)
	GetList(ctx context.Context, req *types.UserListRequest) (*sharedTypes.PageResult, error)
	ResetPassword(ctx context.Context, id uint, password string) error
	GetStats(ctx context.Context) (*types.UserStatsResponse, error)
}

// userService 普通用户服务实现
type userService struct {
	// TODO: 添加用户仓储依赖
}

// NewUserService 创建普通用户服务
func NewUserService() UserService {
	return &userService{}
}

// GetByID 根据ID获取用户
func (s *userService) GetByID(ctx context.Context, id uint) (*types.UserResponse, error) {
	// TODO: 实现用户查询逻辑
	return nil, errors.New("用户服务暂未实现")
}

// Create 创建用户
func (s *userService) Create(ctx context.Context, req *types.UserCreateRequest) (*types.UserResponse, error) {
	// TODO: 实现用户创建逻辑
	return nil, errors.New("用户服务暂未实现")
}

// Update 更新用户
func (s *userService) Update(ctx context.Context, id uint, req *types.UserUpdateRequest) (*types.UserResponse, error) {
	// TODO: 实现用户更新逻辑
	return nil, errors.New("用户服务暂未实现")
}

// Delete 删除用户
func (s *userService) Delete(ctx context.Context, id uint) error {
	// TODO: 实现用户删除逻辑
	return errors.New("用户服务暂未实现")
}

// UpdateStatus 更新用户状态
func (s *userService) UpdateStatus(ctx context.Context, id uint, status int) error {
	// TODO: 实现用户状态更新逻辑
	return errors.New("用户服务暂未实现")
}

// List 获取用户列表
func (s *userService) List(ctx context.Context, req *types.UserListRequest) ([]*types.UserResponse, int64, error) {
	// TODO: 实现用户列表查询逻辑
	return nil, 0, errors.New("用户服务暂未实现")
}

// GetList 获取用户列表（分页）
func (s *userService) GetList(ctx context.Context, req *types.UserListRequest) (*sharedTypes.PageResult, error) {
	// TODO: 实现用户列表查询逻辑
	return nil, errors.New("用户服务暂未实现")
}

// ResetPassword 重置用户密码
func (s *userService) ResetPassword(ctx context.Context, id uint, password string) error {
	// TODO: 实现用户密码重置逻辑
	return errors.New("用户服务暂未实现")
}

// GetStats 获取用户统计
func (s *userService) GetStats(ctx context.Context) (*types.UserStatsResponse, error) {
	// TODO: 实现用户统计逻辑
	return nil, errors.New("用户服务暂未实现")
}
