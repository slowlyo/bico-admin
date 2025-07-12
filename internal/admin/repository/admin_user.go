package repository

import (
	"context"

	"gorm.io/gorm"

	"bico-admin/internal/shared/model"
	"bico-admin/internal/shared/types"
)

// AdminUserRepository 管理员用户仓储接口
type AdminUserRepository interface {
	// 基础CRUD
	Create(ctx context.Context, user *model.AdminUser) error
	GetByID(ctx context.Context, id uint) (*model.AdminUser, error)
	GetByUsername(ctx context.Context, username string) (*model.AdminUser, error)
	Update(ctx context.Context, user *model.AdminUser) error
	Delete(ctx context.Context, id uint) error

	// 查询方法
	List(ctx context.Context, req *types.BasePageQuery) ([]*model.AdminUser, int64, error)
	ListByStatus(ctx context.Context, enabled bool, req *types.BasePageQuery) ([]*model.AdminUser, int64, error)

	// 状态管理
	UpdateStatus(ctx context.Context, id uint, enabled bool) error
	BatchUpdateStatus(ctx context.Context, ids []uint, enabled bool) error

	// 统计方法
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, enabled bool) (int64, error)

	// 验证方法
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByUsernameExcludeID(ctx context.Context, username string, excludeID uint) (bool, error)
}

// adminUserRepository 管理员用户仓储实现
type adminUserRepository struct {
	db *gorm.DB
}

// NewAdminUserRepository 创建管理员用户仓储
func NewAdminUserRepository(db *gorm.DB) AdminUserRepository {
	return &adminUserRepository{
		db: db,
	}
}

// Create 创建管理员用户
func (r *adminUserRepository) Create(ctx context.Context, user *model.AdminUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID 根据ID获取管理员用户
func (r *adminUserRepository) GetByID(ctx context.Context, id uint) (*model.AdminUser, error) {
	var user model.AdminUser
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取管理员用户
func (r *adminUserRepository) GetByUsername(ctx context.Context, username string) (*model.AdminUser, error) {
	var user model.AdminUser
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新管理员用户
func (r *adminUserRepository) Update(ctx context.Context, user *model.AdminUser) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除管理员用户
func (r *adminUserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.AdminUser{}, id).Error
}

// List 分页获取管理员用户列表
func (r *adminUserRepository) List(ctx context.Context, req *types.BasePageQuery) ([]*model.AdminUser, int64, error) {
	var users []*model.AdminUser
	var total int64

	db := r.db.WithContext(ctx).Model(&model.AdminUser{})

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := req.GetOffset()
	pageSize := req.GetPageSize()
	err := db.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error

	return users, total, err
}

// ListByStatus 根据状态分页获取管理员用户列表
func (r *adminUserRepository) ListByStatus(ctx context.Context, enabled bool, req *types.BasePageQuery) ([]*model.AdminUser, int64, error) {
	var users []*model.AdminUser
	var total int64

	db := r.db.WithContext(ctx).Model(&model.AdminUser{}).Where("enabled = ?", enabled)

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := req.GetOffset()
	pageSize := req.GetPageSize()
	err := db.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error

	return users, total, err
}

// UpdateStatus 更新管理员用户状态
func (r *adminUserRepository) UpdateStatus(ctx context.Context, id uint, enabled bool) error {
	return r.db.WithContext(ctx).Model(&model.AdminUser{}).
		Where("id = ?", id).
		Update("enabled", enabled).Error
}

// BatchUpdateStatus 批量更新管理员用户状态
func (r *adminUserRepository) BatchUpdateStatus(ctx context.Context, ids []uint, enabled bool) error {
	return r.db.WithContext(ctx).Model(&model.AdminUser{}).
		Where("id IN ?", ids).
		Update("enabled", enabled).Error
}

// Count 统计管理员用户总数
func (r *adminUserRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.AdminUser{}).Count(&count).Error
	return count, err
}

// CountByStatus 根据状态统计管理员用户数量
func (r *adminUserRepository) CountByStatus(ctx context.Context, enabled bool) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.AdminUser{}).
		Where("enabled = ?", enabled).
		Count(&count).Error
	return count, err
}

// ExistsByUsername 检查用户名是否存在
func (r *adminUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.AdminUser{}).
		Where("username = ?", username).
		Count(&count).Error
	return count > 0, err
}

// ExistsByUsernameExcludeID 检查用户名是否存在（排除指定ID）
func (r *adminUserRepository) ExistsByUsernameExcludeID(ctx context.Context, username string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.AdminUser{}).
		Where("username = ? AND id != ?", username, excludeID).
		Count(&count).Error
	return count > 0, err
}
