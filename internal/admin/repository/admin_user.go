package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"bico-admin/internal/admin/models"
	adminTypes "bico-admin/internal/admin/types"
	"bico-admin/internal/shared/repository"
	"bico-admin/internal/shared/types"
)

// AdminUserRepository 管理员用户仓储接口
type AdminUserRepository interface {
	// 继承基础仓储接口
	repository.BaseRepositoryInterface[models.AdminUser]

	// 自定义查询方法
	GetByUsername(ctx context.Context, username string) (*models.AdminUser, error)
	ListWithFilter(ctx context.Context, req *adminTypes.AdminUserListRequest) ([]*models.AdminUser, int64, error)

	// 状态管理
	UpdateLastLoginTime(ctx context.Context, id uint) error
	UpdatePassword(ctx context.Context, id uint, hashedPassword string) error

	// 验证方法
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByUsernameExcludeID(ctx context.Context, username string, excludeID uint) (bool, error)

	// 超级管理员相关
	CountSuperAdmins(ctx context.Context) (int64, error)
	CountSuperAdminsExcludeID(ctx context.Context, excludeID uint) (int64, error)
}

// adminUserRepository 管理员用户仓储实现
type adminUserRepository struct {
	repository.BaseRepositoryInterface[models.AdminUser]
}

// NewAdminUserRepository 创建管理员用户仓储
func NewAdminUserRepository(db *gorm.DB) AdminUserRepository {
	return &adminUserRepository{
		BaseRepositoryInterface: repository.NewBaseRepository[models.AdminUser](db),
	}
}

// GetByID 根据ID获取管理员用户（预加载角色和权限）
func (r *adminUserRepository) GetByID(ctx context.Context, id uint) (*models.AdminUser, error) {
	var user models.AdminUser
	err := r.DB().WithContext(ctx).
		Preload("Roles.Role.Permissions").
		First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取管理员用户（预加载角色和权限）
func (r *adminUserRepository) GetByUsername(ctx context.Context, username string) (*models.AdminUser, error) {
	var user models.AdminUser
	err := r.DB().WithContext(ctx).
		Preload("Roles.Role.Permissions").
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateLastLoginTime 更新最后登录时间
func (r *adminUserRepository) UpdateLastLoginTime(ctx context.Context, id uint) error {
	now := time.Now()
	return r.DB().WithContext(ctx).
		Model(&models.AdminUser{}).
		Where("id = ?", id).
		Update("last_login_at", now).Error
}

// UpdatePassword 更新用户密码
func (r *adminUserRepository) UpdatePassword(ctx context.Context, id uint, hashedPassword string) error {
	return r.DB().WithContext(ctx).
		Model(&models.AdminUser{}).
		Where("id = ?", id).
		Update("password", hashedPassword).Error
}

// ExistsByUsername 检查用户名是否存在
func (r *adminUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return r.ExistsByField(ctx, "username", username)
}

// ExistsByUsernameExcludeID 检查用户名是否存在（排除指定ID）
func (r *adminUserRepository) ExistsByUsernameExcludeID(ctx context.Context, username string, excludeID uint) (bool, error) {
	return r.ExistsByFieldExcludeID(ctx, "username", username, excludeID)
}

// ListWithFilter 根据条件分页获取管理员用户列表
func (r *adminUserRepository) ListWithFilter(ctx context.Context, req *adminTypes.AdminUserListRequest) ([]*models.AdminUser, int64, error) {
	var users []*models.AdminUser
	var total int64

	db := r.DB().WithContext(ctx).Model(&models.AdminUser{})

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
	if req.RoleID != nil {
		// 通过角色筛选用户
		db = db.Joins("JOIN admin_user_roles ON admin_users.id = admin_user_roles.user_id").
			Where("admin_user_roles.role_id = ?", *req.RoleID)
	}

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := req.GetOffset()
	pageSize := req.GetPageSize()

	// 构建排序条件
	orderClause := r.buildOrderClause(req.SortBy, req.SortDesc)

	err := db.Preload("Roles.Role").
		Offset(offset).Limit(pageSize).
		Order(orderClause).
		Find(&users).Error

	return users, total, err
}

// CountSuperAdmins 统计超级管理员数量
func (r *adminUserRepository) CountSuperAdmins(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB().WithContext(ctx).
		Table("admin_users").
		Joins("JOIN admin_user_roles ON admin_users.id = admin_user_roles.user_id").
		Joins("JOIN admin_roles ON admin_user_roles.role_id = admin_roles.id").
		Where("admin_roles.code = ? AND admin_users.status = ?", models.RoleCodeSuperAdmin, types.StatusActive).
		Count(&count).Error
	return count, err
}

// CountSuperAdminsExcludeID 统计超级管理员数量（排除指定ID）
func (r *adminUserRepository) CountSuperAdminsExcludeID(ctx context.Context, excludeID uint) (int64, error) {
	var count int64
	err := r.DB().WithContext(ctx).
		Table("admin_users").
		Joins("JOIN admin_user_roles ON admin_users.id = admin_user_roles.user_id").
		Joins("JOIN admin_roles ON admin_user_roles.role_id = admin_roles.id").
		Where("admin_roles.code = ? AND admin_users.status = ? AND admin_users.id != ?",
			models.RoleCodeSuperAdmin, types.StatusActive, excludeID).
		Count(&count).Error
	return count, err
}

// buildOrderClause 构建排序条件
func (r *adminUserRepository) buildOrderClause(sortBy string, sortDesc bool) string {
	// 定义允许排序的字段映射
	allowedSortFields := map[string]string{
		"created_at":    "admin_users.created_at",
		"last_login_at": "admin_users.last_login_at",
		"username":      "admin_users.username",
		"name":          "admin_users.name",
		"status":        "admin_users.status",
	}

	// 检查排序字段是否允许
	dbField, exists := allowedSortFields[sortBy]
	if !exists {
		// 默认按创建时间降序排序
		return "admin_users.created_at DESC"
	}

	// 构建排序方向
	direction := "ASC"
	if sortDesc {
		direction = "DESC"
	}

	return dbField + " " + direction
}
