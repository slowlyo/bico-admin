package repository

import (
	"context"

	"gorm.io/gorm"

	"bico-admin/internal/admin/models"
	adminTypes "bico-admin/internal/admin/types"
	"bico-admin/internal/shared/types"
)

// AdminUserRepository 管理员用户仓储接口
type AdminUserRepository interface {
	// 基础CRUD
	Create(ctx context.Context, user *models.AdminUser) error
	GetByID(ctx context.Context, id uint) (*models.AdminUser, error)
	GetByUsername(ctx context.Context, username string) (*models.AdminUser, error)
	Update(ctx context.Context, user *models.AdminUser) error
	Delete(ctx context.Context, id uint) error

	// 查询方法
	List(ctx context.Context, req *types.BasePageQuery) ([]*models.AdminUser, int64, error)
	ListByStatus(ctx context.Context, enabled bool, req *types.BasePageQuery) ([]*models.AdminUser, int64, error)
	ListWithFilter(ctx context.Context, req *adminTypes.AdminUserListRequest) ([]*models.AdminUser, int64, error)

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
	*BaseRepository[models.AdminUser]
}

// NewAdminUserRepository 创建管理员用户仓储
func NewAdminUserRepository(db *gorm.DB) AdminUserRepository {
	return &adminUserRepository{
		BaseRepository: NewBaseRepository[models.AdminUser](db),
	}
}

// GetByUsername 根据用户名获取管理员用户
func (r *adminUserRepository) GetByUsername(ctx context.Context, username string) (*models.AdminUser, error) {
	var user models.AdminUser
	err := r.BaseRepository.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ListByStatus 根据状态分页获取管理员用户列表
func (r *adminUserRepository) ListByStatus(ctx context.Context, enabled bool, req *types.BasePageQuery) ([]*models.AdminUser, int64, error) {
	status := types.StatusInactive
	if enabled {
		status = types.StatusActive
	}
	return r.BaseRepository.ListByStatus(ctx, status, req)
}

// UpdateStatus 更新管理员用户状态
func (r *adminUserRepository) UpdateStatus(ctx context.Context, id uint, enabled bool) error {
	status := types.StatusInactive
	if enabled {
		status = types.StatusActive
	}
	return r.BaseRepository.UpdateStatus(ctx, id, status)
}

// BatchUpdateStatus 批量更新管理员用户状态
func (r *adminUserRepository) BatchUpdateStatus(ctx context.Context, ids []uint, enabled bool) error {
	status := types.StatusInactive
	if enabled {
		status = types.StatusActive
	}
	return r.BaseRepository.BatchUpdateStatus(ctx, ids, status)
}

// CountByStatus 根据状态统计管理员用户数量
func (r *adminUserRepository) CountByStatus(ctx context.Context, enabled bool) (int64, error) {
	status := types.StatusInactive
	if enabled {
		status = types.StatusActive
	}
	return r.BaseRepository.CountByStatus(ctx, status)
}

// ExistsByUsername 检查用户名是否存在
func (r *adminUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return r.BaseRepository.ExistsByField(ctx, "username", username)
}

// ExistsByUsernameExcludeID 检查用户名是否存在（排除指定ID）
func (r *adminUserRepository) ExistsByUsernameExcludeID(ctx context.Context, username string, excludeID uint) (bool, error) {
	return r.BaseRepository.ExistsByFieldExcludeID(ctx, "username", username, excludeID)
}

// ListWithFilter 根据条件分页获取管理员用户列表
func (r *adminUserRepository) ListWithFilter(ctx context.Context, req *adminTypes.AdminUserListRequest) ([]*models.AdminUser, int64, error) {
	var users []*models.AdminUser
	var total int64

	db := r.BaseRepository.db.WithContext(ctx).Model(&models.AdminUser{})

	// 条件过滤
	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

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
