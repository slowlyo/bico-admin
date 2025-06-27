package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含通用字段
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// PaginationQuery 分页查询参数
type PaginationQuery struct {
	Page     int    `json:"page" query:"page" validate:"min=1"`
	PageSize int    `json:"page_size" query:"page_size" validate:"min=1,max=100"`
	Sort     string `json:"sort" query:"sort"`
	Order    string `json:"order" query:"order" validate:"oneof=asc desc"`
	Search   string `json:"search" query:"search"`
}

// PaginationResult 分页结果
type PaginationResult struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// DefaultPagination 默认分页参数
func DefaultPagination() PaginationQuery {
	return PaginationQuery{
		Page:     1,
		PageSize: 10,
		Sort:     "id",
		Order:    "desc",
	}
}

// GetOffset 计算偏移量
func (p *PaginationQuery) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取限制数量
func (p *PaginationQuery) GetLimit() int {
	return p.PageSize
}

// Validate 验证分页参数
func (p *PaginationQuery) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	if p.Sort == "" {
		p.Sort = "id"
	}
	if p.Order != "asc" && p.Order != "desc" {
		p.Order = "desc"
	}
}
