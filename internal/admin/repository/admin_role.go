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

	Update(ctx context.Context, role *models.AdminRole) error
	Delete(ctx context.Context, id uint) error

	// 查询方法
	List(ctx context.Context, req *types.RoleListRequest) ([]*models.AdminRole, int64, error)

	ListActiveRoles(ctx context.Context) ([]*models.AdminRole, error)

	// 验证方法
	ExistsByCode(ctx context.Context, code string) (bool, error)

	// 权限相关
	CreateRolePermissions(ctx context.Context, tx *gorm.DB, roleID uint, permissionCodes []string) error
	DeleteRolePermissions(ctx context.Context, tx *gorm.DB, roleID uint) error

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
	*BaseRepository[models.AdminRole]
}

// NewAdminRoleRepository 创建管理员角色仓储
func NewAdminRoleRepository(db *gorm.DB) AdminRoleRepository {
	return &adminRoleRepository{
		BaseRepository: NewBaseRepository[models.AdminRole](db),
	}
}

// WithTx 使用事务
func (r *adminRoleRepository) WithTx(tx *gorm.DB) AdminRoleRepository {
	return &adminRoleRepository{
		BaseRepository: NewBaseRepository[models.AdminRole](tx),
	}
}

// GetByID 根据ID获取管理员角色
func (r *adminRoleRepository) GetByID(ctx context.Context, id uint) (*models.AdminRole, error) {
	var role models.AdminRole
	err := r.BaseRepository.db.WithContext(ctx).Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// List 分页获取管理员角色列表
func (r *adminRoleRepository) List(ctx context.Context, req *types.RoleListRequest) ([]*models.AdminRole, int64, error) {
	var roles []*models.AdminRole
	var total int64

	db := r.BaseRepository.db.WithContext(ctx).Model(&models.AdminRole{})

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

	// 构建排序条件
	orderClause := r.buildOrderClause(req.SortBy, req.SortDesc)

	err := db.Preload("Permissions").
		Offset(offset).Limit(pageSize).
		Order(orderClause).
		Find(&roles).Error

	return roles, total, err
}

// ListActiveRoles 获取所有启用的角色
func (r *adminRoleRepository) ListActiveRoles(ctx context.Context) ([]*models.AdminRole, error) {
	var roles []*models.AdminRole
	err := r.BaseRepository.db.WithContext(ctx).
		Where("status = ?", sharedTypes.StatusActive).
		Order("created_at DESC").
		Find(&roles).Error
	return roles, err
}

// 这些方法现在直接继承自 BaseRepository，无需重新实现

// ExistsByCode 检查角色代码是否存在
func (r *adminRoleRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	return r.BaseRepository.ExistsByField(ctx, "code", code)
}

// CreateRolePermissions 创建角色权限关联
func (r *adminRoleRepository) CreateRolePermissions(ctx context.Context, tx *gorm.DB, roleID uint, permissionCodes []string) error {
	db := r.BaseRepository.db
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
	db := r.BaseRepository.db
	if tx != nil {
		db = tx
	}

	return db.WithContext(ctx).Where("role_id = ?", roleID).Delete(&models.AdminRolePermission{}).Error
}

// AssignRolesToUser 为用户分配角色
func (r *adminRoleRepository) AssignRolesToUser(ctx context.Context, tx *gorm.DB, userID uint, roleIDs []uint) error {
	db := r.BaseRepository.db
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
	db := r.BaseRepository.db
	if tx != nil {
		db = tx
	}

	return db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.AdminUserRole{}).Error
}

// GetUserRoles 获取用户角色列表
func (r *adminRoleRepository) GetUserRoles(ctx context.Context, userID uint) ([]*models.AdminRole, error) {
	var userRoles []models.AdminUserRole
	err := r.BaseRepository.db.WithContext(ctx).
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
	err := r.BaseRepository.db.WithContext(ctx).
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
	err := r.BaseRepository.db.WithContext(ctx).Model(&models.AdminUserRole{}).
		Where("role_id = ?", roleID).
		Count(&count).Error
	return count, err
}

// CountUsersWithRoles 统计拥有角色的用户数量
func (r *adminRoleRepository) CountUsersWithRoles(ctx context.Context) (int64, error) {
	var count int64
	err := r.BaseRepository.db.WithContext(ctx).Model(&models.AdminUserRole{}).
		Distinct("user_id").
		Count(&count).Error
	return count, err
}

// buildOrderClause 构建排序条件
func (r *adminRoleRepository) buildOrderClause(sortBy string, sortDesc bool) string {
	// 定义允许排序的字段映射
	allowedSortFields := map[string]string{
		"created_at": "created_at",
		"name":       "name",
		"code":       "code",
		"status":     "status",
	}

	// 检查排序字段是否允许
	dbField, exists := allowedSortFields[sortBy]
	if !exists {
		// 默认按创建时间降序排序
		return "created_at DESC"
	}

	// 构建排序方向
	direction := "ASC"
	if sortDesc {
		direction = "DESC"
	}

	return dbField + " " + direction
}
