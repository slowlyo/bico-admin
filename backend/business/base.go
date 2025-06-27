package business

import (
	"gorm.io/gorm"

	"bico-admin/core/model"
)

// BaseService 基础服务结构
type BaseService[T any] struct {
	DB *gorm.DB
}

// NewBaseService 创建基础服务实例
func NewBaseService[T any](db *gorm.DB) *BaseService[T] {
	return &BaseService[T]{
		DB: db,
	}
}

// CreateOne 创建单个记录
func (s *BaseService[T]) CreateOne(data *T) error {
	return s.DB.Create(data).Error
}

// UpdateById 根据ID更新记录
func (s *BaseService[T]) UpdateById(id uint, data *T) error {
	return s.DB.Model(data).Where("id = ?", id).Updates(data).Error
}

// DeleteById 根据ID删除记录（软删除）
func (s *BaseService[T]) DeleteById(id uint) error {
	var model T
	return s.DB.Where("id = ?", id).Delete(&model).Error
}

// GetById 根据ID获取记录
func (s *BaseService[T]) GetById(id uint) (*T, error) {
	var data T
	err := s.DB.Where("id = ?", id).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// GetByIds 根据ID列表获取记录
func (s *BaseService[T]) GetByIds(ids []uint) ([]T, error) {
	var data []T
	err := s.DB.Where("id IN ?", ids).Find(&data).Error
	return data, err
}

// Exists 检查记录是否存在
func (s *BaseService[T]) Exists(id uint) (bool, error) {
	var count int64
	var model T
	err := s.DB.Model(&model).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// Count 统计记录数量
func (s *BaseService[T]) Count(conditions map[string]interface{}) (int64, error) {
	var count int64
	var model T
	query := s.DB.Model(&model)
	
	for key, value := range conditions {
		query = query.Where(key, value)
	}
	
	err := query.Count(&count).Error
	return count, err
}

// FindOne 根据条件查找单个记录
func (s *BaseService[T]) FindOne(conditions map[string]interface{}) (*T, error) {
	var data T
	query := s.DB
	
	for key, value := range conditions {
		query = query.Where(key, value)
	}
	
	err := query.First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// FindAll 根据条件查找所有记录
func (s *BaseService[T]) FindAll(conditions map[string]interface{}) ([]T, error) {
	var data []T
	query := s.DB
	
	for key, value := range conditions {
		query = query.Where(key, value)
	}
	
	err := query.Find(&data).Error
	return data, err
}

// Transaction 执行事务
func (s *BaseService[T]) Transaction(fn func(*gorm.DB) error) error {
	return s.DB.Transaction(fn)
}
