package database

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// PaginationParams 分页查询参数
type PaginationParams struct {
	Page         int                    `json:"page" query:"page"`
	PageSize     int                    `json:"page_size" query:"page_size"`
	Sort         string                 `json:"sort" query:"sort"`
	Order        string                 `json:"order" query:"order"`
	Search       string                 `json:"search" query:"search"`
	SearchFields []string               `json:"search_fields"`
	Filters      map[string]interface{} `json:"filters"`
	Preloads     []string               `json:"preloads"`
}

// PaginationResult 分页查询结果
type PaginationResult[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// BatchCreateParams 批量创建参数
type BatchCreateParams[T any] struct {
	Data      []T `json:"data"`
	BatchSize int `json:"batch_size"`
}

// BatchUpdateParams 批量更新参数
type BatchUpdateParams struct {
	IDs     []uint                 `json:"ids"`
	Updates map[string]interface{} `json:"updates"`
}

// BatchDeleteParams 批量删除参数
type BatchDeleteParams struct {
	IDs        []uint `json:"ids"`
	HardDelete bool   `json:"hard_delete"`
}

// PaginationOperations 分页操作工具
type PaginationOperations[T any] struct {
	*BaseOperations[T]
}

// NewPaginationOperations 创建分页操作实例
func NewPaginationOperations[T any](db *gorm.DB) *PaginationOperations[T] {
	return &PaginationOperations[T]{
		BaseOperations: NewBaseOperations[T](db),
	}
}

// Paginate 分页查询
func (ops *PaginationOperations[T]) Paginate(params PaginationParams) (*PaginationResult[T], error) {
	var data []T
	var total int64

	// 验证和设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	if params.PageSize > 100 {
		params.PageSize = 100 // 限制最大页面大小
	}

	// 构建查询
	query := ops.DB.Model(new(T))

	// 应用过滤条件
	if len(params.Filters) > 0 {
		for key, value := range params.Filters {
			query = query.Where(key, value)
		}
	}

	// 应用搜索条件
	if params.Search != "" && len(params.SearchFields) > 0 {
		query = ops.applySearch(query, params.Search, params.SearchFields)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("统计记录总数失败: %w", err)
	}

	// 应用排序
	if params.Sort != "" {
		order := "ASC"
		if strings.ToUpper(params.Order) == "DESC" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", params.Sort, order))
	} else {
		query = query.Order("id DESC") // 默认按ID降序
	}

	// 应用预加载
	for _, preload := range params.Preloads {
		query = query.Preload(preload)
	}

	// 应用分页
	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Find(&data).Error; err != nil {
		return nil, fmt.Errorf("查询分页数据失败: %w", err)
	}

	// 计算总页数
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &PaginationResult[T]{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// applySearch 应用搜索条件
func (ops *PaginationOperations[T]) applySearch(query *gorm.DB, search string, searchFields []string) *gorm.DB {
	if search == "" || len(searchFields) == 0 {
		return query
	}

	// 构建搜索条件
	searchCondition := ""
	searchArgs := make([]interface{}, 0, len(searchFields))

	for i, field := range searchFields {
		if i > 0 {
			searchCondition += " OR "
		}
		searchCondition += fmt.Sprintf("%s LIKE ?", field)
		searchArgs = append(searchArgs, "%"+search+"%")
	}

	return query.Where(searchCondition, searchArgs...)
}

// BatchCreate 批量创建
func (ops *PaginationOperations[T]) BatchCreate(params BatchCreateParams[T]) error {
	if len(params.Data) == 0 {
		return fmt.Errorf("批量创建数据不能为空")
	}

	batchSize := params.BatchSize
	if batchSize <= 0 {
		batchSize = 100 // 默认批次大小
	}

	return ops.DB.CreateInBatches(params.Data, batchSize).Error
}

// BatchUpdate 批量更新
func (ops *PaginationOperations[T]) BatchUpdate(params BatchUpdateParams) error {
	if len(params.IDs) == 0 {
		return fmt.Errorf("批量更新ID列表不能为空")
	}
	if len(params.Updates) == 0 {
		return fmt.Errorf("批量更新数据不能为空")
	}

	var model T
	result := ops.DB.Model(&model).Where("id IN ?", params.IDs).Updates(params.Updates)
	if result.Error != nil {
		return fmt.Errorf("批量更新失败: %w", result.Error)
	}
	return nil
}

// BatchDelete 批量删除
func (ops *PaginationOperations[T]) BatchDelete(params BatchDeleteParams) error {
	if len(params.IDs) == 0 {
		return fmt.Errorf("批量删除ID列表不能为空")
	}

	var model T
	var result *gorm.DB

	if params.HardDelete {
		result = ops.DB.Unscoped().Where("id IN ?", params.IDs).Delete(&model)
	} else {
		result = ops.DB.Where("id IN ?", params.IDs).Delete(&model)
	}

	if result.Error != nil {
		return fmt.Errorf("批量删除失败: %w", result.Error)
	}
	return nil
}

// GetWithPagination 获取分页数据（简化版本）
func (ops *PaginationOperations[T]) GetWithPagination(page, pageSize int, conditions map[string]interface{}) (*PaginationResult[T], error) {
	params := PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Filters:  conditions,
	}
	return ops.Paginate(params)
}

// SearchWithPagination 搜索分页数据
func (ops *PaginationOperations[T]) SearchWithPagination(page, pageSize int, search string, searchFields []string) (*PaginationResult[T], error) {
	params := PaginationParams{
		Page:         page,
		PageSize:     pageSize,
		Search:       search,
		SearchFields: searchFields,
	}
	return ops.Paginate(params)
}
