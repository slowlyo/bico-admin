package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// Pagination 分页参数
type Pagination struct {
	Page      int    `json:"page" form:"page"`
	PageSize  int    `json:"pageSize" form:"pageSize"`
	SortField string `json:"sortField" form:"sortField"`
	SortOrder string `json:"sortOrder" form:"sortOrder"`
}

// GetOffset 获取偏移量
func (p *Pagination) GetOffset() int {
	if p.Page < 1 {
		p.Page = DefaultPage
	}
	return (p.Page - 1) * p.GetPageSize()
}

// GetPageSize 获取每页数量（带默认值和最大值限制）
func (p *Pagination) GetPageSize() int {
	if p.PageSize < 1 {
		return DefaultPageSize
	}
	if p.PageSize > MaxPageSize {
		return MaxPageSize
	}
	return p.PageSize
}

// GetPage 获取页码（带默认值）
func (p *Pagination) GetPage() int {
	if p.Page < 1 {
		return DefaultPage
	}
	return p.Page
}

// FromContext 从上下文获取分页参数
func FromContext(c *gin.Context) *Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(DefaultPage)))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", strconv.Itoa(DefaultPageSize)))

	p := &Pagination{
		Page:      page,
		PageSize:  pageSize,
		SortField: c.Query("sortField"),
		SortOrder: c.Query("sortOrder"),
	}

	return p
}

// GetOrderBy 获取排序子句
func (p *Pagination) GetOrderBy() string {
	if p.SortField == "" {
		return ""
	}
	
	order := "DESC"
	if p.SortOrder == "ascend" {
		order = "ASC"
	}
	
	return p.SortField + " " + order
}

// Response 分页响应
type Response struct {
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}
