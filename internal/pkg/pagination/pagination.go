package pagination

import (
	"slices"
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

// GetOrderBy 获取排序子句。
//
// 说明：只有出现在 sortable 白名单中的字段才允许排序，避免客户端按敏感列（如密码哈希）
// 排序来推断其内容。字段不在白名单内时返回空串，由调用方回退默认排序。
func (p *Pagination) GetOrderBy(sortable []string) string {
	// 白名单为空或字段未命中时忽略客户端排序，字段名不会被拼进 SQL。
	if !slices.Contains(sortable, p.SortField) {
		return ""
	}

	order := "DESC"
	// 前端约定 ascend 表示升序，其余值默认降序。
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
