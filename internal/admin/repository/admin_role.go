package repository

import (
	"context"

	"gorm.io/gorm"

	"bico-admin/internal/admin/models"
	"bico-admin/internal/admin/types"
	sharedTypes "bico-admin/internal/shared/types"
)

// AdminRoleRepository 管理员角色仓储接口
type AdminRoleRepository interface {
	// 基础CRUD
	Create(ctx context.Context, role *models.AdminRole) error
	GetByID(ctx context.Context, id uint) (*models.AdminRole, error)
	GetByCode(ctx context.Context, code string) (*models.AdminRole, error)
	Update(ctx context.Context, role *models.AdminRole) error
	Delete(ctx context.Context, id uint) error

	// 查询方法
	List(ctx context.Context, req *types.RoleListRequest) ([]*models.AdminRole, int64, error)
	ListByStatus(ctx context.Context, status int, req *sharedTypes.BasePageQuery) ([]*models.AdminRole, int64, error)
	ListActiveRoles(ctx context.Context) ([]*models.AdminRole, error)

	// 状态管理
	UpdateStatus(ctx context.Context, id uint, status int) error
	BatchUpdateStatus(ctx context.Context, ids []uint, status int) error

	// 统计方法
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status int) (int64, error)

	// 验证方法
	ExistsByCode(ctx context.Context, code string) (bool, error)
	ExistsByCodeExcludeID(ctx context.Context, code string, excludeID uint) (bool, error)

	// 权限相关
	CreateRolePermissions(ctx context.Context, tx *gorm.DB, roleID uint, permissionCodes []string) error
	DeleteRolePermissions(ctx context.Context, tx *gorm.DB, roleID uint) error
	GetRolePermissions(ctx context.Context, roleID uint) ([]string, error)

	// 用户角色关联
	AssignRolesToUser(ctx context.Context, tx *gorm.DB, userID uint, roleIDs []uint) error
	DeleteUserRoles(ctx context.Context, tx *gorm.DB, userID uint) error
	GetUserRoles(ctx context.Context, userID uint) ([]*models.AdminRole, error)
	GetUserPermissions(ctx context.Context, userID uint) ([]string, error)
	CountUsersByRoleID(ctx context.Context, roleID uint) (int64, error)
	CountUsersWithRoles(ctx context.Context) (int64, error)

	// 事务支持
	WithTx(tx *gorm.DB) AdminRoleRepository
}

// adminRoleRepository 管理员角色仓储实现
type adminRoleRepository struct {
	db *gorm.DB
}

// NewAdminRoleRepository 创建管理员角色仓储
func NewAdminRoleRepository(db *gorm.DB) AdminRoleRepository {
	return &adminRoleRepository{
		db: db,
	}
}

// WithTx 使用事务
func (r *adminRoleRepository) WithTx(tx *gorm.DB) AdminRoleRepository {
	return &adminRoleRepository{
		db: tx,
	}
}

// Create 创建管理员角色
func (r *adminRoleRepository) Create(ctx context.Context, role *models.AdminRole) error {
	return r.db.WithContext(ctx).Create(role).Error
}

// GetByID 根据ID获取管理员角色
func (r *adminRoleRepository) GetByID(ctx context.Context, id uint) (*models.AdminRole, error) {
	var role models.AdminRole
	err := r.db.WithContext(ctx).Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByCode 根据代码获取管理员角色
func (r *adminRoleRepository) GetByCode(ctx context.Context, code string) (*models.AdminRole, error) {
	var role models.AdminRole
	err := r.db.WithContext(ctx).Preload("Permissions").Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// Update 更新管理员角色
func (r *adminRoleRepository) Update(ctx context.Context, role *models.AdminRole) error {
	return r.db.WithContext(ctx).Save(role).Error
}

// Delete 删除管理员角色
func (r *adminRoleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.AdminRole{}, id).Error
}

// List 分页获取管理员角色列表
func (r *adminRoleRepository) List(ctx context.Context, req *types.RoleListRequest) ([]*models.AdminRole, int64, error) {
	var roles []*models.AdminRole
	var total int64

	db := r.db.WithContext(ctx).Model(&models.AdminRole{})

	// 条件过滤
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		db = db.Where("code LIKE ?", "%"+req.Code+"%")
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
	err := db.Preload("Permissions").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&roles).Error

	return roles, total, err
}

// ListByStatus 根据状态分页获取管理员角色列表
func (r *adminRoleRepository) ListByStatus(ctx context.Context, status int, req *sharedTypes.BasePageQuery) ([]*models.AdminRole, int64, error) {
	var roles []*models.AdminRole
	var total int64

	db := r.db.WithContext(ctx).Model(&models.AdminRole{}).Where("status = ?", status)

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := req.GetOffset()
	pageSize := req.GetPageSize()
	err := db.Preload("Permissions").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&roles).Error

	return roles, total, err
}

// ListActiveRoles 获取所有启用的角色
func (r *adminRoleRepository) ListActiveRoles(ctx context.Context) ([]*models.AdminRole, error) {
	var roles []*models.AdminRole
	err := r.db.WithContext(ctx).
		Where("status = ?", sharedTypes.StatusActive).
		Order("created_at DESC").
		Find(&roles).Error
	return roles, err
}

// UpdateStatus 更新管理员角色状态
func (r *adminRoleRepository) UpdateStatus(ctx context.Context, id uint, status int) error {
	return r.db.WithContext(ctx).Model(&models.AdminRole{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// BatchUpdateStatus 批量更新管理员角色状态
func (r *adminRoleRepository) BatchUpdateStatus(ctx context.Context, ids []uint, status int) error {
	return r.db.WithContext(ctx).Model(&models.AdminRole{}).
		Where("id IN ?", ids).
		Update("status", status).Error
}

// Count 统计管理员角色总数
func (r *adminRoleRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AdminRole{}).Count(&count).Error
	return count, err
}

// CountByStatus 根据状态统计管理员角色数量
func (r *adminRoleRepository) CountByStatus(ctx context.Context, status int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AdminRole{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

// ExistsByCode 检查角色代码是否存在
func (r *adminRoleRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AdminRole{}).
		Where("code = ?", code).
		Count(&count).Error
	return count > 0, err
}

// ExistsByCodeExcludeID 检查角色代码是否存在（排除指定ID）
func (r *adminRoleRepository) ExistsByCodeExcludeID(ctx context.Context, code string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AdminRole{}).
		Where("code = ? AND id != ?", code, excludeID).
		Count(&count).Error
	return count > 0, err
}

// CreateRolePermissions 创建角色权限关联
func (r *adminRoleRepository) CreateRolePermissions(ctx context.Context, tx *gorm.DB, roleID uint, permissionCodes []string) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	for _, code := range permissionCodes {
		rolePermission := models.AdminRolePermission{
			RoleID:         roleID,
			PermissionCode: code,
		}
		if err := db.WithContext(ctx).Create(&rolePermission).Error; err != nil {
			return err
		}
	}
	return nil
}

// DeleteRolePermissions 删除角色权限关联
func (r *adminRoleRepository) DeleteRolePermissions(ctx context.Context, tx *gorm.DB, roleID uint) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	return db.WithContext(ctx).Where("role_id = ?", roleID).Delete(&models.AdminRolePermission{}).Error
}

// GetRolePermissions 获取角色权限代码列表
func (r *adminRoleRepository) GetRolePermissions(ctx context.Context, roleID uint) ([]string, error) {
	var permissions []models.AdminRolePermission
	err := r.db.WithContext(ctx).Where("role_id = ?", roleID).Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	var codes []string
	for _, permission := range permissions {
		codes = append(codes, permission.PermissionCode)
	}
	return codes, nil
}

// AssignRolesToUser 为用户分配角色
func (r *adminRoleRepository) AssignRolesToUser(ctx context.Context, tx *gorm.DB, userID uint, roleIDs []uint) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	for _, roleID := range roleIDs {
		userRole := models.AdminUserRole{
			UserID: userID,
			RoleID: roleID,
		}
		if err := db.WithContext(ctx).Create(&userRole).Error; err != nil {
			return err
		}
	}
	return nil
}

// DeleteUserRoles 删除用户角色关联
func (r *adminRoleRepository) DeleteUserRoles(ctx context.Context, tx *gorm.DB, userID uint) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	return db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.AdminUserRole{}).Error
}

// GetUserRoles 获取用户角色列表
func (r *adminRoleRepository) GetUserRoles(ctx context.Context, userID uint) ([]*models.AdminRole, error) {
	var userRoles []models.AdminUserRole
	err := r.db.WithContext(ctx).
		Preload("Role.Permissions").
		Where("user_id = ?", userID).
		Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	var roles []*models.AdminRole
	for _, userRole := range userRoles {
		roles = append(roles, &userRole.Role)
	}
	return roles, nil
}

// GetUserPermissions 获取用户所有权限代码
func (r *adminRoleRepository) GetUserPermissions(ctx context.Context, userID uint) ([]string, error) {
	var userRoles []models.AdminUserRole
	err := r.db.WithContext(ctx).
		Preload("Role.Permissions").
		Where("user_id = ?", userID).
		Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	permissionSet := make(map[string]bool)
	for _, userRole := range userRoles {
		if userRole.Role.IsEnabled() {
			for _, permission := range userRole.Role.Permissions {
				permissionSet[permission.PermissionCode] = true
			}
		}
	}

	var permissions []string
	for permission := range permissionSet {
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

// CountUsersByRoleID 统计使用指定角色的用户数量
func (r *adminRoleRepository) CountUsersByRoleID(ctx context.Context, roleID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AdminUserRole{}).
		Where("role_id = ?", roleID).
		Count(&count).Error
	return count, err
}

// CountUsersWithRoles 统计拥有角色的用户数量
func (r *adminRoleRepository) CountUsersWithRoles(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AdminUserRole{}).
		Distinct("user_id").
		Count(&count).Error
	return count, err
}
