package types

// BasePageQuery 基础分页查询参数
type BasePageQuery struct {
	Page     int    `form:"page" json:"page" binding:"min=1"`                   // 页码，从1开始
	PageSize int    `form:"page_size" json:"page_size" binding:"min=1,max=100"` // 每页数量，最大100
	Keyword  string `form:"keyword" json:"keyword"`                             // 搜索关键词
	SortBy   string `form:"sort_by" json:"sort_by"`                             // 排序字段
	SortDesc bool   `form:"sort_desc" json:"sort_desc"`                         // 是否降序
}

// GetOffset 获取偏移量
func (p *BasePageQuery) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.GetPageSize()
}

// GetPageSize 获取每页数量
func (p *BasePageQuery) GetPageSize() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}

// GetSortOrder 获取排序方式
func (p *BasePageQuery) GetSortOrder() string {
	if p.SortDesc {
		return "DESC"
	}
	return "ASC"
}

// PageResult 分页结果
type PageResult struct {
	List       interface{} `json:"list"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// NewPageResult 创建分页结果
func NewPageResult(list interface{}, total int64, page, pageSize int) *PageResult {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PageResult{
		List:       list,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
