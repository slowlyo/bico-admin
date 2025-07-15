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
	UpdateLastLoginTime(ctx context.Context, id uint) error
	List(ctx context.Context, req *sharedTypes.BasePageQuery) ([]*models.AdminUser, int64, error)
	ListWithFilter(ctx context.Context, req *types.AdminUserListRequest) ([]*models.AdminUser, int64, error)

	// 权限检查方法
	CanUserBeDeleted(ctx context.Context, userID uint) (bool, error)
	CanUserBeDisabled(ctx context.Context, userID uint) (bool, error)
}

// adminUserService 管理员用户服务实现
type adminUserService struct {
	adminUserRepo repository.AdminUserRepository
	adminRoleRepo repository.AdminRoleRepository
}

// NewAdminUserService 创建管理员用户服务
func NewAdminUserService(adminUserRepo repository.AdminUserRepository, adminRoleRepo repository.AdminRoleRepository) AdminUserService {
	return &adminUserService{
		adminUserRepo: adminUserRepo,
		adminRoleRepo: adminRoleRepo,
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
		Email:    req.Email,
		Phone:    req.Phone,
		Remark:   req.Remark,
		Status:   status,
	}

	if err := s.adminUserRepo.Create(ctx, adminUser); err != nil {
		return nil, err
	}

	// 分配角色
	if len(req.RoleIDs) > 0 {
		if err := s.adminRoleRepo.AssignRolesToUser(ctx, nil, adminUser.ID, req.RoleIDs); err != nil {
			// 如果角色分配失败，删除已创建的用户
			s.adminUserRepo.Delete(ctx, adminUser.ID)
			return nil, err
		}
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
	adminUser.Email = req.Email
	adminUser.Phone = req.Phone
	adminUser.Remark = req.Remark

	// 转换 Enabled 字段为 Status
	// 如果要禁用超级管理员，需要检查是否会导致系统没有可用的超管
	if adminUser.IsSuperAdmin() && !req.Enabled {
		// 统计除了当前用户外的其他超级管理员数量
		otherSuperAdminCount, err := s.adminUserRepo.CountSuperAdminsExcludeID(ctx, id)
		if err != nil {
			return nil, err
		}

		// 如果禁用后系统将没有可用的超级管理员，则不允许禁用
		if otherSuperAdminCount == 0 {
			return nil, errors.New("系统必须保留至少一个可用的超级管理员，无法禁用")
		}
	}

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

	// 更新角色分配
	if req.RoleIDs != nil {
		// 删除用户现有角色
		if err := s.adminRoleRepo.DeleteUserRoles(ctx, nil, id); err != nil {
			return nil, err
		}

		// 分配新角色
		if len(req.RoleIDs) > 0 {
			if err := s.adminRoleRepo.AssignRolesToUser(ctx, nil, id, req.RoleIDs); err != nil {
				return nil, err
			}
		}
	}

	return adminUser, nil
}

// Delete 删除管理员用户
func (s *adminUserService) Delete(ctx context.Context, id uint) error {
	// 检查用户是否存在
	adminUser, err := s.adminUserRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 如果是超级管理员，需要检查是否会导致系统没有超管
	if adminUser.IsSuperAdmin() {
		// 统计除了当前用户外的其他超级管理员数量
		otherSuperAdminCount, err := s.adminUserRepo.CountSuperAdminsExcludeID(ctx, id)
		if err != nil {
			return err
		}

		// 如果删除后系统将没有超级管理员，则不允许删除
		if otherSuperAdminCount == 0 {
			return errors.New("系统必须保留至少一个超级管理员，无法删除")
		}
	}

	// 先删除用户角色关联
	if err := s.adminRoleRepo.DeleteUserRoles(ctx, nil, id); err != nil {
		return err
	}

	// 再删除用户
	return s.adminUserRepo.Delete(ctx, id)
}

// UpdateStatus 更新管理员用户状态
func (s *adminUserService) UpdateStatus(ctx context.Context, id uint, enabled bool) error {
	// 检查用户是否存在
	adminUser, err := s.adminUserRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 如果要禁用超级管理员，需要检查是否会导致系统没有可用的超管
	if adminUser.IsSuperAdmin() && !enabled {
		// 统计除了当前用户外的其他超级管理员数量
		otherSuperAdminCount, err := s.adminUserRepo.CountSuperAdminsExcludeID(ctx, id)
		if err != nil {
			return err
		}

		// 如果禁用后系统将没有可用的超级管理员，则不允许禁用
		if otherSuperAdminCount == 0 {
			return errors.New("系统必须保留至少一个可用的超级管理员，无法禁用")
		}
	}

	return s.adminUserRepo.UpdateStatus(ctx, id, enabled)
}

// UpdateLastLoginTime 更新最后登录时间
func (s *adminUserService) UpdateLastLoginTime(ctx context.Context, id uint) error {
	return s.adminUserRepo.UpdateLastLoginTime(ctx, id)
}

// List 获取管理员用户列表
func (s *adminUserService) List(ctx context.Context, req *sharedTypes.BasePageQuery) ([]*models.AdminUser, int64, error) {
	return s.adminUserRepo.List(ctx, req)
}

// ListWithFilter 获取管理员用户列表（带筛选）
func (s *adminUserService) ListWithFilter(ctx context.Context, req *types.AdminUserListRequest) ([]*models.AdminUser, int64, error) {
	return s.adminUserRepo.ListWithFilter(ctx, req)
}

// CanUserBeDeleted 检查用户是否可以被删除
func (s *adminUserService) CanUserBeDeleted(ctx context.Context, userID uint) (bool, error) {
	// 检查用户是否存在
	adminUser, err := s.adminUserRepo.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}

	// 如果不是超级管理员，可以删除
	if !adminUser.IsSuperAdmin() {
		return true, nil
	}

	// 如果是超级管理员，检查是否还有其他超管
	otherSuperAdminCount, err := s.adminUserRepo.CountSuperAdminsExcludeID(ctx, userID)
	if err != nil {
		return false, err
	}

	// 如果删除后系统将没有超级管理员，则不允许删除
	return otherSuperAdminCount > 0, nil
}

// CanUserBeDisabled 检查用户是否可以被禁用
func (s *adminUserService) CanUserBeDisabled(ctx context.Context, userID uint) (bool, error) {
	// 检查用户是否存在
	adminUser, err := s.adminUserRepo.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}

	// 如果不是超级管理员，可以禁用
	if !adminUser.IsSuperAdmin() {
		return true, nil
	}

	// 如果是超级管理员，检查是否还有其他超管
	otherSuperAdminCount, err := s.adminUserRepo.CountSuperAdminsExcludeID(ctx, userID)
	if err != nil {
		return false, err
	}

	// 如果禁用后系统将没有可用的超级管理员，则不允许禁用
	return otherSuperAdminCount > 0, nil
}
