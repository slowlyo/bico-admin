package repository

import (
	"gorm.io/gorm"
)

// BaseRepository 基础仓储接口
type BaseRepository[T any] interface {
	Create(data *T) error
	GetByID(id uint) (*T, error)
	Update(id uint, data *T) error
	Delete(id uint) error
	List(offset, limit int) ([]T, error)
	Count() (int64, error)
}

// BaseRepositoryImpl 基础仓储实现
type BaseRepositoryImpl[T any] struct {
	DB *gorm.DB
}

// NewBaseRepository 创建基础仓储实例
func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &BaseRepositoryImpl[T]{
		DB: db,
	}
}

// Create 创建记录
func (r *BaseRepositoryImpl[T]) Create(data *T) error {
	return r.DB.Create(data).Error
}

// GetByID 根据ID获取记录
func (r *BaseRepositoryImpl[T]) GetByID(id uint) (*T, error) {
	var data T
	err := r.DB.Where("id = ?", id).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// Update 更新记录
func (r *BaseRepositoryImpl[T]) Update(id uint, data *T) error {
	return r.DB.Model(data).Where("id = ?", id).Updates(data).Error
}

// Delete 删除记录
func (r *BaseRepositoryImpl[T]) Delete(id uint) error {
	var model T
	return r.DB.Where("id = ?", id).Delete(&model).Error
}

// List 获取列表
func (r *BaseRepositoryImpl[T]) List(offset, limit int) ([]T, error) {
	var data []T
	err := r.DB.Offset(offset).Limit(limit).Find(&data).Error
	return data, err
}

// Count 统计数量
func (r *BaseRepositoryImpl[T]) Count() (int64, error) {
	var count int64
	var model T
	err := r.DB.Model(&model).Count(&count).Error
	return count, err
}
