package database

import (
	"gorm.io/gorm"
)

// Operations 完整的数据库操作工具 - 包含所有通用方法
type Operations[T any] struct {
	*PaginationOperations[T]
}

// NewOperations 创建完整的数据库操作实例
func NewOperations[T any](db *gorm.DB) *Operations[T] {
	return &Operations[T]{
		PaginationOperations: NewPaginationOperations[T](db),
	}
}

// Transaction 执行事务操作
func (ops *Operations[T]) Transaction(fn func(*gorm.DB) error) error {
	return ops.DB.Transaction(fn)
}

// GetDB 获取数据库连接
func (ops *Operations[T]) GetDB() *gorm.DB {
	return ops.DB
}

// WithDB 使用指定的数据库连接创建新的操作实例
func (ops *Operations[T]) WithDB(db *gorm.DB) *Operations[T] {
	return NewOperations[T](db)
}

// Raw 执行原生SQL查询
func (ops *Operations[T]) Raw(sql string, args ...interface{}) *gorm.DB {
	return ops.DB.Raw(sql, args...)
}

// Exec 执行原生SQL命令
func (ops *Operations[T]) Exec(sql string, args ...interface{}) error {
	return ops.DB.Exec(sql, args...).Error
}

// 便捷方法 - 为常用操作提供简化接口

// Create 创建记录（别名）
func (ops *Operations[T]) Create(data *T) error {
	return ops.CreateOne(data)
}

// Get 获取记录（别名）
func (ops *Operations[T]) Get(id uint) (*T, error) {
	return ops.GetById(id)
}

// Update 更新记录（别名）
func (ops *Operations[T]) Update(id uint, data *T) error {
	return ops.UpdateById(id, data)
}

// UpdateFields 更新指定字段（别名）
func (ops *Operations[T]) UpdateFields(id uint, updates map[string]interface{}) error {
	return ops.UpdateByIdWithMap(id, updates)
}

// Delete 删除记录（别名）
func (ops *Operations[T]) Delete(id uint) error {
	return ops.DeleteById(id)
}

// List 列表查询（别名）
func (ops *Operations[T]) List(params PaginationParams) (*PaginationResult[T], error) {
	return ops.Paginate(params)
}

// Find 查找记录（别名）
func (ops *Operations[T]) Find(conditions map[string]interface{}) ([]T, error) {
	return ops.FindByConditions(conditions)
}

// First 获取第一个记录（别名）
func (ops *Operations[T]) First(conditions map[string]interface{}) (*T, error) {
	return ops.GetByConditions(conditions)
}

// 工具函数

// ValidatePaginationParams 验证分页参数
func ValidatePaginationParams(params *PaginationParams) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}
	if params.Order != "ASC" && params.Order != "DESC" {
		params.Order = "DESC"
	}
}

// BuildSearchParams 构建搜索参数
func BuildSearchParams(page, pageSize int, search string, searchFields []string) PaginationParams {
	return PaginationParams{
		Page:         page,
		PageSize:     pageSize,
		Search:       search,
		SearchFields: searchFields,
	}
}

// BuildFilterParams 构建过滤参数
func BuildFilterParams(page, pageSize int, filters map[string]interface{}) PaginationParams {
	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}
}

// BuildCompleteParams 构建完整参数
func BuildCompleteParams(page, pageSize int, sort, order, search string, searchFields []string, filters map[string]interface{}, preloads []string) PaginationParams {
	return PaginationParams{
		Page:         page,
		PageSize:     pageSize,
		Sort:         sort,
		Order:        order,
		Search:       search,
		SearchFields: searchFields,
		Filters:      filters,
		Preloads:     preloads,
	}
}
