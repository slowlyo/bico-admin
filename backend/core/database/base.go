package database

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// BaseOperations 基础数据库操作工具 - 提供通用的CRUD方法
type BaseOperations[T any] struct {
	DB *gorm.DB
}

// NewBaseOperations 创建基础操作实例
func NewBaseOperations[T any](db *gorm.DB) *BaseOperations[T] {
	return &BaseOperations[T]{
		DB: db,
	}
}

// CreateOne 创建单个记录
func (ops *BaseOperations[T]) CreateOne(data *T) error {
	if err := ops.DB.Create(data).Error; err != nil {
		return fmt.Errorf("创建记录失败: %w", err)
	}
	return nil
}

// GetById 根据ID获取记录
func (ops *BaseOperations[T]) GetById(id uint) (*T, error) {
	var data T
	if err := ops.DB.Where("id = ?", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("记录不存在: ID=%d", id)
		}
		return nil, fmt.Errorf("查询记录失败: %w", err)
	}
	return &data, nil
}

// GetByIds 根据ID列表获取记录
func (ops *BaseOperations[T]) GetByIds(ids []uint) ([]T, error) {
	var data []T
	if err := ops.DB.Where("id IN ?", ids).Find(&data).Error; err != nil {
		return nil, fmt.Errorf("批量查询记录失败: %w", err)
	}
	return data, nil
}

// UpdateById 根据ID更新记录
func (ops *BaseOperations[T]) UpdateById(id uint, data *T) error {
	result := ops.DB.Model(data).Where("id = ?", id).Updates(data)
	if result.Error != nil {
		return fmt.Errorf("更新记录失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("记录不存在或无需更新: ID=%d", id)
	}
	return nil
}

// UpdateByIdWithMap 根据ID更新记录（使用map）
func (ops *BaseOperations[T]) UpdateByIdWithMap(id uint, updates map[string]interface{}) error {
	var model T
	result := ops.DB.Model(&model).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("更新记录失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("记录不存在或无需更新: ID=%d", id)
	}
	return nil
}

// DeleteById 根据ID删除记录（软删除）
func (ops *BaseOperations[T]) DeleteById(id uint) error {
	var model T
	result := ops.DB.Where("id = ?", id).Delete(&model)
	if result.Error != nil {
		return fmt.Errorf("删除记录失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("记录不存在: ID=%d", id)
	}
	return nil
}

// HardDeleteById 根据ID硬删除记录
func (ops *BaseOperations[T]) HardDeleteById(id uint) error {
	var model T
	result := ops.DB.Unscoped().Where("id = ?", id).Delete(&model)
	if result.Error != nil {
		return fmt.Errorf("硬删除记录失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("记录不存在: ID=%d", id)
	}
	return nil
}

// Exists 检查记录是否存在
func (ops *BaseOperations[T]) Exists(id uint) (bool, error) {
	var count int64
	var model T
	if err := ops.DB.Model(&model).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("检查记录存在性失败: %w", err)
	}
	return count > 0, nil
}

// Count 统计记录总数
func (ops *BaseOperations[T]) Count() (int64, error) {
	var count int64
	var model T
	if err := ops.DB.Model(&model).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("统计记录数量失败: %w", err)
	}
	return count, nil
}

// CountWithCondition 根据条件统计记录数量
func (ops *BaseOperations[T]) CountWithCondition(condition string, args ...interface{}) (int64, error) {
	var count int64
	var model T
	if err := ops.DB.Model(&model).Where(condition, args...).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("条件统计记录数量失败: %w", err)
	}
	return count, nil
}

// CountWithConditions 根据多个条件统计记录数量
func (ops *BaseOperations[T]) CountWithConditions(conditions map[string]interface{}) (int64, error) {
	var count int64
	var model T
	query := ops.DB.Model(&model)

	for key, value := range conditions {
		query = query.Where(key, value)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("条件统计记录数量失败: %w", err)
	}
	return count, nil
}

// GetByCondition 根据条件获取单个记录
func (ops *BaseOperations[T]) GetByCondition(condition string, args ...interface{}) (*T, error) {
	var data T
	if err := ops.DB.Where(condition, args...).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("记录不存在")
		}
		return nil, fmt.Errorf("查询记录失败: %w", err)
	}
	return &data, nil
}

// GetByConditions 根据多个条件获取单个记录
func (ops *BaseOperations[T]) GetByConditions(conditions map[string]interface{}) (*T, error) {
	var data T
	query := ops.DB

	for key, value := range conditions {
		query = query.Where(key, value)
	}

	if err := query.First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("记录不存在")
		}
		return nil, fmt.Errorf("查询记录失败: %w", err)
	}
	return &data, nil
}

// FindByCondition 根据条件获取多个记录
func (ops *BaseOperations[T]) FindByCondition(condition string, args ...interface{}) ([]T, error) {
	var data []T
	if err := ops.DB.Where(condition, args...).Find(&data).Error; err != nil {
		return nil, fmt.Errorf("查询记录失败: %w", err)
	}
	return data, nil
}

// FindByConditions 根据多个条件获取多个记录
func (ops *BaseOperations[T]) FindByConditions(conditions map[string]interface{}) ([]T, error) {
	var data []T
	query := ops.DB

	for key, value := range conditions {
		query = query.Where(key, value)
	}

	if err := query.Find(&data).Error; err != nil {
		return nil, fmt.Errorf("查询记录失败: %w", err)
	}
	return data, nil
}
