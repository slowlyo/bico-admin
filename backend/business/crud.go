package business

import (
	"fmt"

	"gorm.io/gorm"

	"bico-admin/core/model"
)

// CRUDService CRUD服务结构
type CRUDService[T any] struct {
	*BaseService[T]
}

// NewCRUDService 创建CRUD服务实例
func NewCRUDService[T any](db *gorm.DB) *CRUDService[T] {
	return &CRUDService[T]{
		BaseService: NewBaseService[T](db),
	}
}

// ListParams 列表查询参数
type ListParams struct {
	Page     int                    `json:"page" query:"page"`
	PageSize int                    `json:"page_size" query:"page_size"`
	Sort     string                 `json:"sort" query:"sort"`
	Order    string                 `json:"order" query:"order"`
	Search   string                 `json:"search" query:"search"`
	Filters  map[string]interface{} `json:"filters"`
}

// ListResult 列表查询结果
type ListResult[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// List 分页列表查询
func (s *CRUDService[T]) List(params ListParams) (*ListResult[T], error) {
	// 验证和设置默认值
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}
	if params.Sort == "" {
		params.Sort = "id"
	}
	if params.Order != "asc" && params.Order != "desc" {
		params.Order = "desc"
	}

	var data []T
	var total int64

	// 构建查询
	query := s.DB.Model(new(T))

	// 应用过滤条件
	if params.Filters != nil {
		for key, value := range params.Filters {
			query = query.Where(key, value)
		}
	}

	// 应用搜索条件（需要在具体实现中重写）
	if params.Search != "" {
		query = s.applySearch(query, params.Search)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 应用排序和分页
	offset := (params.Page - 1) * params.PageSize
	orderClause := fmt.Sprintf("%s %s", params.Sort, params.Order)
	
	if err := query.Order(orderClause).Offset(offset).Limit(params.PageSize).Find(&data).Error; err != nil {
		return nil, err
	}

	// 计算总页数
	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &ListResult[T]{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

// applySearch 应用搜索条件（默认实现，可在具体服务中重写）
func (s *CRUDService[T]) applySearch(query *gorm.DB, search string) *gorm.DB {
	// 默认不应用搜索，具体实现需要重写此方法
	return query
}

// BatchCreate 批量创建
func (s *CRUDService[T]) BatchCreate(data []T) error {
	return s.DB.Create(&data).Error
}

// BatchUpdate 批量更新
func (s *CRUDService[T]) BatchUpdate(ids []uint, updates map[string]interface{}) error {
	var model T
	return s.DB.Model(&model).Where("id IN ?", ids).Updates(updates).Error
}

// BatchDelete 批量删除
func (s *CRUDService[T]) BatchDelete(ids []uint) error {
	var model T
	return s.DB.Where("id IN ?", ids).Delete(&model).Error
}

// SoftDelete 软删除
func (s *CRUDService[T]) SoftDelete(id uint) error {
	var model T
	return s.DB.Where("id = ?", id).Delete(&model).Error
}

// Restore 恢复软删除的记录
func (s *CRUDService[T]) Restore(id uint) error {
	var model T
	return s.DB.Unscoped().Model(&model).Where("id = ?", id).Update("deleted_at", nil).Error
}

// HardDelete 硬删除
func (s *CRUDService[T]) HardDelete(id uint) error {
	var model T
	return s.DB.Unscoped().Where("id = ?", id).Delete(&model).Error
}

// GetWithRelations 获取记录及其关联数据
func (s *CRUDService[T]) GetWithRelations(id uint, relations []string) (*T, error) {
	var data T
	query := s.DB.Where("id = ?", id)
	
	for _, relation := range relations {
		query = query.Preload(relation)
	}
	
	err := query.First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}
