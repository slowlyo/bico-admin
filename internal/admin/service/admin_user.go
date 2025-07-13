package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"bico-admin/internal/admin/models"
	"bico-admin/internal/admin/repository"
	"bico-admin/internal/admin/types"
	sharedTypes "bico-admin/internal/shared/types"
)

// AdminUserService 管理员用户服务接口
type AdminUserService interface {
	GetByID(ctx context.Context, id uint) (*models.AdminUser, error)
	GetByUsername(ctx context.Context, username string) (*models.AdminUser, error)
	Create(ctx context.Context, req *types.AdminUserCreateRequest) (*models.AdminUser, error)
	Update(ctx context.Context, id uint, req *types.AdminUserUpdateRequest) (*models.AdminUser, error)
	Delete(ctx context.Context, id uint) error
	UpdateStatus(ctx context.Context, id uint, enabled bool) error
	List(ctx context.Context, req *sharedTypes.BasePageQuery) ([]*models.AdminUser, int64, error)
}

// adminUserService 管理员用户服务实现
type adminUserService struct {
	adminUserRepo repository.AdminUserRepository
}

// NewAdminUserService 创建管理员用户服务
func NewAdminUserService(adminUserRepo repository.AdminUserRepository) AdminUserService {
	return &adminUserService{
		adminUserRepo: adminUserRepo,
	}
}

// GetByID 根据ID获取管理员用户
func (s *adminUserService) GetByID(ctx context.Context, id uint) (*models.AdminUser, error) {
	return s.adminUserRepo.GetByID(ctx, id)
}

// GetByUsername 根据用户名获取管理员用户
func (s *adminUserService) GetByUsername(ctx context.Context, username string) (*models.AdminUser, error) {
	return s.adminUserRepo.GetByUsername(ctx, username)
}

// Create 创建管理员用户
func (s *adminUserService) Create(ctx context.Context, req *types.AdminUserCreateRequest) (*models.AdminUser, error) {
	// 检查用户名是否已存在
	if exists, err := s.adminUserRepo.ExistsByUsername(ctx, req.Username); err != nil {
		return nil, err
	} else if exists {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 转换 Enabled 字段为 Status
	status := sharedTypes.StatusInactive
	if req.Enabled {
		status = sharedTypes.StatusActive
	}

	adminUser := &models.AdminUser{
		Username: req.Username,
		Password: string(hashedPassword),
		Name:     req.Name,
		Avatar:   req.Avatar,
		Status:   status,
	}

	if err := s.adminUserRepo.Create(ctx, adminUser); err != nil {
		return nil, err
	}

	return adminUser, nil
}

// Update 更新管理员用户
func (s *adminUserService) Update(ctx context.Context, id uint, req *types.AdminUserUpdateRequest) (*models.AdminUser, error) {
	adminUser, err := s.adminUserRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 检查用户名是否被其他用户使用
	if req.Username != adminUser.Username {
		if exists, err := s.adminUserRepo.ExistsByUsernameExcludeID(ctx, req.Username, id); err != nil {
			return nil, err
		} else if exists {
			return nil, errors.New("用户名已存在")
		}
	}

	// 更新字段
	adminUser.Username = req.Username
	adminUser.Name = req.Name
	adminUser.Avatar = req.Avatar

	// 转换 Enabled 字段为 Status
	if req.Enabled {
		adminUser.Status = sharedTypes.StatusActive
	} else {
		adminUser.Status = sharedTypes.StatusInactive
	}

	// 如果提供了新密码，则更新密码
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("密码加密失败")
		}
		adminUser.Password = string(hashedPassword)
	}

	if err := s.adminUserRepo.Update(ctx, adminUser); err != nil {
		return nil, err
	}

	return adminUser, nil
}

// Delete 删除管理员用户
func (s *adminUserService) Delete(ctx context.Context, id uint) error {
	return s.adminUserRepo.Delete(ctx, id)
}

// UpdateStatus 更新管理员用户状态
func (s *adminUserService) UpdateStatus(ctx context.Context, id uint, enabled bool) error {
	return s.adminUserRepo.UpdateStatus(ctx, id, enabled)
}

// List 获取管理员用户列表
func (s *adminUserService) List(ctx context.Context, req *sharedTypes.BasePageQuery) ([]*models.AdminUser, int64, error) {
	return s.adminUserRepo.List(ctx, req)
}
