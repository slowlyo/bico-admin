package repository

import (
	"context"
	"reflect"

	"gorm.io/gorm"

	"bico-admin/internal/shared/types"
)

// BaseRepository 基础仓储，提供通用的 CRUD 操作
type BaseRepository[T any] struct {
	db    *gorm.DB
	model T
}

// NewBaseRepository 创建基础仓储
func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	var model T
	return &BaseRepository[T]{
		db:    db,
		model: model,
	}
}

// WithTx 使用事务
func (r *BaseRepository[T]) WithTx(tx *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		db:    tx,
		model: r.model,
	}
}

// Create 创建记录
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// GetByID 根据ID获取记录
func (r *BaseRepository[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update 更新记录
func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete 删除记录
func (r *BaseRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

// Count 统计总数
func (r *BaseRepository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	var entity T
	err := r.db.WithContext(ctx).Model(&entity).Count(&count).Error
	return count, err
}

// CountByStatus 根据状态统计数量
func (r *BaseRepository[T]) CountByStatus(ctx context.Context, status int) (int64, error) {
	var count int64
	var entity T
	err := r.db.WithContext(ctx).Model(&entity).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

// UpdateStatus 更新状态
func (r *BaseRepository[T]) UpdateStatus(ctx context.Context, id uint, status int) error {
	var entity T
	return r.db.WithContext(ctx).Model(&entity).
		Where("id = ?", id).
		Update("status", status).Error
}

// BatchUpdateStatus 批量更新状态
func (r *BaseRepository[T]) BatchUpdateStatus(ctx context.Context, ids []uint, status int) error {
	var entity T
	return r.db.WithContext(ctx).Model(&entity).
		Where("id IN ?", ids).
		Update("status", status).Error
}

// ExistsByField 检查字段值是否存在
func (r *BaseRepository[T]) ExistsByField(ctx context.Context, field string, value interface{}) (bool, error) {
	var count int64
	var entity T
	err := r.db.WithContext(ctx).Model(&entity).
		Where(field+" = ?", value).
		Count(&count).Error
	return count > 0, err
}

// ExistsByFieldExcludeID 检查字段值是否存在（排除指定ID）
func (r *BaseRepository[T]) ExistsByFieldExcludeID(ctx context.Context, field string, value interface{}, excludeID uint) (bool, error) {
	var count int64
	var entity T
	err := r.db.WithContext(ctx).Model(&entity).
		Where(field+" = ? AND id != ?", value, excludeID).
		Count(&count).Error
	return count > 0, err
}

// ListByStatus 根据状态分页查询
func (r *BaseRepository[T]) ListByStatus(ctx context.Context, status int, req *types.BasePageQuery) ([]*T, int64, error) {
	var entities []*T
	var total int64
	var entity T

	db := r.db.WithContext(ctx).Model(&entity).Where("status = ?", status)

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := req.GetOffset()
	pageSize := req.GetPageSize()
	err := db.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&entities).Error

	return entities, total, err
}

// List 分页查询所有记录
func (r *BaseRepository[T]) List(ctx context.Context, req *types.BasePageQuery) ([]*T, int64, error) {
	var entities []*T
	var total int64
	var entity T

	db := r.db.WithContext(ctx).Model(&entity)

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := req.GetOffset()
	pageSize := req.GetPageSize()
	err := db.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&entities).Error

	return entities, total, err
}

// GetModelType 获取模型类型（用于调试）
func (r *BaseRepository[T]) GetModelType() reflect.Type {
	return reflect.TypeOf(r.model)
}
