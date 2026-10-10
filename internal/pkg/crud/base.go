package crud

import (
	"errors"
	"strconv"

	"bico-admin/internal/core/logger"
	"bico-admin/internal/pkg/pagination"
	"bico-admin/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

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
	if err := c.ShouldBindQuery(req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return err
	}
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

// Fail 按错误类型输出响应。
//
// 说明：BizError 的消息可以展示给用户；其余错误（数据库驱动错误等）可能带有 SQL、
// 表名或约束名，只记录日志并返回通用提示，避免内部细节泄露给客户端。
func (h *BaseHandler) Fail(c *gin.Context, err error) {
	// 业务错误：消息由业务代码编写，原样返回。
	if msg, ok := bizMessage(err); ok {
		h.Error(c, msg)
		return
	}
	// 系统错误：记录原因，对外只给通用提示。
	logger.Error("请求处理失败",
		zap.String("method", c.Request.Method),
		zap.String("path", c.FullPath()),
		zap.Error(err),
	)
	response.ErrorWithCode(c, 500, "服务器内部错误")
}

// NotFound 404 响应
func (h *BaseHandler) NotFound(c *gin.Context, msg string) {
	response.NotFound(c, msg)
}

// defaultListOrder 未配置 ListOrder.Default 时使用的排序。
const defaultListOrder = "created_at DESC"

// ListOrder 描述列表排序策略。
type ListOrder struct {
	// Sortable 允许客户端通过 sortField 排序的列名白名单；为空时忽略客户端排序。
	Sortable []string
	// Default 客户端未指定或指定了不可排序字段时的排序子句；为空时使用 created_at DESC。
	Default string
}

// QueryList 通用分页查询，只使用默认排序。
//
// 说明：需要支持客户端排序时使用 QueryListWithHook 并配置 ListOrder.Sortable。
func (h *BaseHandler) QueryList(c *gin.Context, query *gorm.DB, dest interface{}) {
	QueryListWithHook(h, c, query, dest, ListOrder{}, nil)
}

// QueryListWithHook 通用分页查询，并支持排序白名单与结果二次处理
func QueryListWithHook(
	h *BaseHandler,
	c *gin.Context,
	query *gorm.DB,
	dest interface{},
	order ListOrder,
	after func() error,
) {
	pg := h.GetPagination(c)

	var total int64
	// Count 会改写当前 Statement 的 SELECT，先复制一份再计数，列表查询仍按原字段取出。
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		h.Fail(c, err)
		return
	}

	if orderBy := pg.GetOrderBy(order.Sortable); orderBy != "" {
		query = query.Order(orderBy)
	} else if order.Default != "" {
		// 表没有 created_at 或需要自定义顺序时，由模块指定默认排序。
		query = query.Order(order.Default)
	} else {
		query = query.Order(defaultListOrder)
	}

	if err := query.Offset(pg.GetOffset()).Limit(pg.GetPageSize()).Find(dest).Error; err != nil {
		h.Fail(c, err)
		return
	}

	if after != nil {
		if err := after(); err != nil {
			h.Fail(c, err)
			return
		}
	}

	h.SuccessWithPagination(c, dest, total)
}

// QueryOne 通用单条查询
func (h *BaseHandler) QueryOne(c *gin.Context, query *gorm.DB, dest interface{}, notFoundMsg string) bool {
	err := query.First(dest).Error
	if err == nil {
		return true
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		h.NotFound(c, notFoundMsg)
		return false
	}

	h.Fail(c, err)
	return false
}

// ExecDelete 通用删除操作
func (h *BaseHandler) ExecDelete(c *gin.Context, db *gorm.DB, model interface{}, id uint) {
	if err := db.Delete(model, id).Error; err != nil {
		h.Fail(c, err)
		return
	}
	h.SuccessWithMessage(c, "删除成功", nil)
}

// ExecTx 通用事务操作
func (h *BaseHandler) ExecTx(c *gin.Context, db *gorm.DB, fn func(tx *gorm.DB) error, successMsg string, data interface{}) {
	if err := db.Transaction(fn); err != nil {
		h.Fail(c, err)
		return
	}
	h.SuccessWithMessage(c, successMsg, data)
}
