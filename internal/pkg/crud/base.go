package crud

import (
	"strconv"

	"bico-admin/internal/pkg/pagination"
	"bico-admin/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BaseService 基础 CRUD 服务
type BaseService[T any] struct {
	DB *gorm.DB
}

// NewBaseService 创建基础服务
func NewBaseService[T any](db *gorm.DB) *BaseService[T] {
	return &BaseService[T]{DB: db}
}

// List 获取列表
func (s *BaseService[T]) List(query *gorm.DB, pg *pagination.Pagination) (*pagination.Response, error) {
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 应用排序
	if orderBy := pg.GetOrderBy(); orderBy != "" {
		query = query.Order(orderBy)
	} else {
		query = query.Order("created_at DESC")
	}

	var items []T
	if err := query.Offset(pg.GetOffset()).Limit(pg.GetPageSize()).Find(&items).Error; err != nil {
		return nil, err
	}

	return &pagination.Response{
		Total: total,
		Data:  items,
	}, nil
}

// Get 获取单条记录
func (s *BaseService[T]) Get(id uint) (*T, error) {
	var item T
	if err := s.DB.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Create 创建记录
func (s *BaseService[T]) Create(item *T) error {
	return s.DB.Create(item).Error
}

// Update 更新记录
func (s *BaseService[T]) Update(id uint, updates map[string]interface{}) error {
	var item T
	return s.DB.Model(&item).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除记录
func (s *BaseService[T]) Delete(id uint) error {
	var item T
	return s.DB.Delete(&item, id).Error
}

// Exists 检查记录是否存在
func (s *BaseService[T]) Exists(query *gorm.DB) (bool, error) {
	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

// BaseHandler 基础 CRUD Handler，可嵌入使用
type BaseHandler struct{}

// ParseID 从路由参数解析 ID
func (h *BaseHandler) ParseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return 0, err
	}
	return uint(id), nil
}

// BindJSON 绑定 JSON 请求体
func (h *BaseHandler) BindJSON(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return err
	}
	return nil
}

// BindQuery 绑定 Query 参数
func (h *BaseHandler) BindQuery(c *gin.Context, req interface{}) error {
	_ = c.ShouldBindQuery(req)
	return nil
}

// GetPagination 获取分页参数
func (h *BaseHandler) GetPagination(c *gin.Context) *pagination.Pagination {
	return pagination.FromContext(c)
}

// Success 成功响应
func (h *BaseHandler) Success(c *gin.Context, data interface{}) {
	response.SuccessWithData(c, data)
}

// SuccessWithMessage 带消息的成功响应
func (h *BaseHandler) SuccessWithMessage(c *gin.Context, msg string, data interface{}) {
	response.SuccessWithMessage(c, msg, data)
}

// SuccessWithPagination 带分页的成功响应
func (h *BaseHandler) SuccessWithPagination(c *gin.Context, data interface{}, total int64) {
	response.SuccessWithPagination(c, data, total)
}

// Error 错误响应
func (h *BaseHandler) Error(c *gin.Context, msg string) {
	response.ErrorWithCode(c, 400, msg)
}

// NotFound 404 响应
func (h *BaseHandler) NotFound(c *gin.Context, msg string) {
	response.NotFound(c, msg)
}

// QueryList 通用分页查询
func (h *BaseHandler) QueryList(c *gin.Context, query *gorm.DB, dest interface{}) {
	pg := h.GetPagination(c)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		h.Error(c, err.Error())
		return
	}

	if orderBy := pg.GetOrderBy(); orderBy != "" {
		query = query.Order(orderBy)
	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Offset(pg.GetOffset()).Limit(pg.GetPageSize()).Find(dest).Error; err != nil {
		h.Error(c, err.Error())
		return
	}

	h.SuccessWithPagination(c, dest, total)
}

// QueryOne 通用单条查询
func (h *BaseHandler) QueryOne(c *gin.Context, query *gorm.DB, dest interface{}, notFoundMsg string) bool {
	if err := query.First(dest).Error; err != nil {
		h.NotFound(c, notFoundMsg)
		return false
	}
	return true
}

// ExecDelete 通用删除操作
func (h *BaseHandler) ExecDelete(c *gin.Context, db *gorm.DB, model interface{}, id uint) {
	if err := db.Delete(model, id).Error; err != nil {
		h.Error(c, err.Error())
		return
	}
	h.SuccessWithMessage(c, "删除成功", nil)
}

// ExecTx 通用事务操作
func (h *BaseHandler) ExecTx(c *gin.Context, db *gorm.DB, fn func(tx *gorm.DB) error, successMsg string, data interface{}) {
	if err := db.Transaction(fn); err != nil {
		h.Error(c, err.Error())
		return
	}
	h.SuccessWithMessage(c, successMsg, data)
}
