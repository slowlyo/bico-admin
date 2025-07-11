package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"bico-admin/internal/admin/types"
	"bico-admin/internal/shared/model"
	sharedTypes "bico-admin/internal/shared/types"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	GetList(ctx context.Context, req *types.UserListRequest) ([]*model.User, int64, error)
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	UpdateStatus(ctx context.Context, id uint, status int) error
	UpdatePassword(ctx context.Context, id uint, password string) error
	UpdateLoginInfo(ctx context.Context, id uint, ip string) error
	GetStats(ctx context.Context) (*types.UserStatsResponse, error)
}

// userRepository 用户仓储实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// GetList 获取用户列表
func (r *userRepository) GetList(ctx context.Context, req *types.UserListRequest) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})

	// 添加搜索条件
	if req.Keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.UserType != "" {
		query = query.Where("user_type = ?", req.UserType)
	}

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	if req.Email != "" {
		query = query.Where("email LIKE ?", "%"+req.Email+"%")
	}

	if req.Phone != "" {
		query = query.Where("phone LIKE ?", "%"+req.Phone+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 添加排序
	orderBy := "created_at DESC"
	if req.SortBy != "" {
		orderBy = req.SortBy + " " + req.GetSortOrder()
	}

	// 分页查询
	if err := query.Order(orderBy).
		Offset(req.GetOffset()).
		Limit(req.GetPageSize()).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update 更新用户
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户（软删除）
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// UpdateStatus 更新用户状态
func (r *userRepository) UpdateStatus(ctx context.Context, id uint, status int) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// UpdatePassword 更新用户密码
func (r *userRepository) UpdatePassword(ctx context.Context, id uint, password string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Update("password", password).Error
}

// UpdateLoginInfo 更新登录信息
func (r *userRepository) UpdateLoginInfo(ctx context.Context, id uint, ip string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_at":  now,
			"last_login_ip":  ip,
			"login_count":    gorm.Expr("login_count + 1"),
		}).Error
}

// GetStats 获取用户统计
func (r *userRepository) GetStats(ctx context.Context) (*types.UserStatsResponse, error) {
	var stats types.UserStatsResponse

	// 总用户数
	r.db.WithContext(ctx).Model(&model.User{}).Count(&stats.TotalUsers)

	// 激活用户数
	r.db.WithContext(ctx).Model(&model.User{}).
		Where("status = ?", sharedTypes.StatusActive).
		Count(&stats.ActiveUsers)

	// 管理员用户数
	r.db.WithContext(ctx).Model(&model.User{}).
		Where("user_type = ?", sharedTypes.UserTypeAdmin).
		Count(&stats.AdminUsers)

	// 主控用户数
	r.db.WithContext(ctx).Model(&model.User{}).
		Where("user_type = ?", sharedTypes.UserTypeMaster).
		Count(&stats.MasterUsers)

	// 普通用户数
	r.db.WithContext(ctx).Model(&model.User{}).
		Where("user_type = ?", sharedTypes.UserTypeNormal).
		Count(&stats.NormalUsers)

	// 今日登录数
	today := time.Now().Format("2006-01-02")
	r.db.WithContext(ctx).Model(&model.User{}).
		Where("DATE(last_login_at) = ?", today).
		Count(&stats.TodayLogins)

	// 本周登录数
	weekAgo := time.Now().AddDate(0, 0, -7)
	r.db.WithContext(ctx).Model(&model.User{}).
		Where("last_login_at >= ?", weekAgo).
		Count(&stats.WeeklyLogins)

	return &stats, nil
}
