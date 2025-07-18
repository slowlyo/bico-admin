package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"bico-admin/internal/shared/types"
	"bico-admin/pkg/response"
)

// BaseHandler 基础处理器
type BaseHandler[T any, CreateReq any, UpdateReq any, ListReq any, Response any] struct {
	service  ServiceInterface[T]
	options  *HandlerOptions
	metadata *HandlerMetadata
	cache    CacheManager
}

// NewBaseHandler 创建基础处理器
func NewBaseHandler[T any, CreateReq any, UpdateReq any, ListReq any, Response any](
	service ServiceInterface[T],
	options *HandlerOptions,
) *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response] {
	if options == nil {
		options = DefaultHandlerOptions()
	}

	return &BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]{
		service: service,
		options: options,
	}
}

// SetCache 设置缓存管理器
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) SetCache(cache CacheManager) {
	h.cache = cache
}

// SetMetadata 设置元数据
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) SetMetadata(metadata *HandlerMetadata) {
	h.metadata = metadata
}

// GetByID 根据ID获取记录
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) GetByID(c *gin.Context) {
	var req types.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 检查缓存
	if h.options.EnableCache && h.cache != nil {
		cacheKey := fmt.Sprintf("entity:%d", req.ID)
		if cached, found := h.cache.Get(cacheKey); found {
			if entity, ok := cached.(*T); ok {
				resp := h.ConvertToResponse(c, entity)
				response.Success(c, resp)
				return
			}
		}
	}

	entity, err := h.service.GetByID(c.Request.Context(), req.ID)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeNotFound, err.Error())
		return
	}

	// 设置缓存
	if h.options.EnableCache && h.cache != nil {
		cacheKey := fmt.Sprintf("entity:%d", req.ID)
		h.cache.Set(cacheKey, entity, h.options.CacheExpiration)
	}

	resp := h.ConvertToResponse(c, entity)
	response.Success(c, resp)
}

// Create 创建记录
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) Create(c *gin.Context) {
	var req CreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 转换请求为实体
	entity := h.ConvertCreateRequest(c, &req)
	if entity == nil {
		response.ErrorWithMessage(c, response.CodeBadRequest, "请求转换失败")
		return
	}

	// 执行创建前事件
	if eventHandler, ok := h.service.(EventHandler[T]); ok {
		if err := eventHandler.BeforeCreate(c.Request.Context(), entity); err != nil {
			response.ErrorWithMessage(c, response.CodeBadRequest, err.Error())
			return
		}
	}

	if err := h.service.Create(c.Request.Context(), entity); err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	// 执行创建后事件
	if eventHandler, ok := h.service.(EventHandler[T]); ok {
		if err := eventHandler.AfterCreate(c.Request.Context(), entity); err != nil {
			// 记录警告，但不影响响应
			fmt.Printf("警告: 创建后事件处理失败: %v\n", err)
		}
	}

	// 清除相关缓存
	if h.options.EnableCache && h.cache != nil {
		h.cache.Clear() // 简单实现，清除所有缓存
	}

	resp := h.ConvertToResponse(c, entity)
	response.Success(c, resp)
}

// Update 更新记录
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) Update(c *gin.Context) {
	var idReq types.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	var req UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 转换请求为实体
	entity := h.ConvertUpdateRequest(c, idReq.ID, &req)
	if entity == nil {
		response.ErrorWithMessage(c, response.CodeBadRequest, "请求转换失败")
		return
	}

	// 执行更新前事件
	if eventHandler, ok := h.service.(EventHandler[T]); ok {
		if err := eventHandler.BeforeUpdate(c.Request.Context(), entity); err != nil {
			response.ErrorWithMessage(c, response.CodeBadRequest, err.Error())
			return
		}
	}

	if err := h.service.Update(c.Request.Context(), entity); err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	// 执行更新后事件
	if eventHandler, ok := h.service.(EventHandler[T]); ok {
		if err := eventHandler.AfterUpdate(c.Request.Context(), entity); err != nil {
			fmt.Printf("警告: 更新后事件处理失败: %v\n", err)
		}
	}

	// 清除缓存
	if h.options.EnableCache && h.cache != nil {
		cacheKey := fmt.Sprintf("entity:%d", idReq.ID)
		h.cache.Delete(cacheKey)
	}

	resp := h.ConvertToResponse(c, entity)
	response.Success(c, resp)
}

// Delete 删除记录
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) Delete(c *gin.Context) {
	var req types.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 执行删除前事件
	if eventHandler, ok := h.service.(EventHandler[T]); ok {
		if err := eventHandler.BeforeDelete(c.Request.Context(), req.ID); err != nil {
			response.ErrorWithMessage(c, response.CodeBadRequest, err.Error())
			return
		}
	}

	if err := h.service.Delete(c.Request.Context(), req.ID); err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	// 执行删除后事件
	if eventHandler, ok := h.service.(EventHandler[T]); ok {
		if err := eventHandler.AfterDelete(c.Request.Context(), req.ID); err != nil {
			fmt.Printf("警告: 删除后事件处理失败: %v\n", err)
		}
	}

	// 清除缓存
	if h.options.EnableCache && h.cache != nil {
		cacheKey := fmt.Sprintf("entity:%d", req.ID)
		h.cache.Delete(cacheKey)
	}

	response.Success(c, nil)
}

// UpdateStatus 更新状态
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) UpdateStatus(c *gin.Context) {
	if !h.options.EnableStatusManagement {
		response.ErrorWithMessage(c, response.CodeForbidden, "状态管理功能未启用")
		return
	}

	var idReq types.IDRequest
	if err := c.ShouldBindUri(&idReq); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	var statusReq types.StatusRequest
	if err := c.ShouldBindJSON(&statusReq); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 获取当前状态（用于事件处理）
	var oldStatus int
	if _, err := h.service.GetByID(c.Request.Context(), idReq.ID); err == nil {
		// 这里需要根据实际模型结构获取状态，暂时设为0
		oldStatus = 0
	}

	if err := h.service.UpdateStatus(c.Request.Context(), idReq.ID, statusReq.Status); err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	// 执行状态变更事件
	if eventHandler, ok := h.service.(EventHandler[T]); ok {
		if err := eventHandler.OnStatusChange(c.Request.Context(), idReq.ID, oldStatus, statusReq.Status); err != nil {
			fmt.Printf("警告: 状态变更事件处理失败: %v\n", err)
		}
	}

	// 清除缓存
	if h.options.EnableCache && h.cache != nil {
		cacheKey := fmt.Sprintf("entity:%d", idReq.ID)
		h.cache.Delete(cacheKey)
	}

	response.Success(c, nil)
}

// List 获取列表
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) List(c *gin.Context) {
	var req ListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 转换为基础分页查询
	pageQuery := h.ConvertListRequest(c, &req)
	if pageQuery == nil {
		response.ErrorWithMessage(c, response.CodeBadRequest, "请求转换失败")
		return
	}

	// 应用默认分页设置
	if pageQuery.PageSize <= 0 {
		pageQuery.PageSize = h.options.DefaultPageSize
	}
	if pageQuery.PageSize > h.options.MaxPageSize {
		pageQuery.PageSize = h.options.MaxPageSize
	}

	result, err := h.service.ListWithFilter(c.Request.Context(), pageQuery)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	// 转换列表项为响应格式
	responseList := h.ConvertListToResponse(c, result.List)

	response.Page(c, responseList, result.Total, result.Page, result.PageSize)
}

// 以下方法需要在具体的 handler 中实现

// ConvertToResponse 转换实体为响应格式（需要子类实现）
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) ConvertToResponse(c *gin.Context, entity *T) Response {
	panic("ConvertToResponse must be implemented by concrete handler")
}

// ConvertCreateRequest 转换创建请求为实体（需要子类实现）
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) ConvertCreateRequest(c *gin.Context, req *CreateReq) *T {
	panic("ConvertCreateRequest must be implemented by concrete handler")
}

// ConvertUpdateRequest 转换更新请求为实体（需要子类实现）
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) ConvertUpdateRequest(c *gin.Context, id uint, req *UpdateReq) *T {
	panic("ConvertUpdateRequest must be implemented by concrete handler")
}

// ConvertListRequest 转换列表请求为基础分页查询（需要子类实现）
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) ConvertListRequest(c *gin.Context, req *ListReq) *types.BasePageQuery {
	panic("ConvertListRequest must be implemented by concrete handler")
}

// ConvertListToResponse 转换列表为响应格式（需要子类实现）
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) ConvertListToResponse(c *gin.Context, list any) []Response {
	panic("ConvertListToResponse must be implemented by concrete handler")
}

// BatchDelete 批量删除
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) BatchDelete(c *gin.Context) {
	if !h.options.EnableBatchOperations {
		response.ErrorWithMessage(c, response.CodeForbidden, "批量操作功能未启用")
		return
	}

	var req BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if len(req.IDs) == 0 {
		response.ErrorWithMessage(c, response.CodeBadRequest, "ID列表不能为空")
		return
	}

	// 逐个删除（因为基础服务接口没有批量删除方法）
	for _, id := range req.IDs {
		if err := h.service.Delete(c.Request.Context(), id); err != nil {
			response.ErrorWithMessage(c, response.CodeInternalServerError, fmt.Sprintf("删除ID %d 失败: %v", id, err))
			return
		}
	}

	// 清除缓存
	if h.options.EnableCache && h.cache != nil {
		for _, id := range req.IDs {
			cacheKey := fmt.Sprintf("entity:%d", id)
			h.cache.Delete(cacheKey)
		}
	}

	response.Success(c, nil)
}

// BatchUpdateStatus 批量更新状态
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) BatchUpdateStatus(c *gin.Context) {
	if !h.options.EnableBatchOperations || !h.options.EnableStatusManagement {
		response.ErrorWithMessage(c, response.CodeForbidden, "批量状态更新功能未启用")
		return
	}

	var req struct {
		IDs    []uint `json:"ids" binding:"required,min=1"`
		Status int    `json:"status" binding:"oneof=0 1 -1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if len(req.IDs) == 0 {
		response.ErrorWithMessage(c, response.CodeBadRequest, "ID列表不能为空")
		return
	}

	// 逐个更新状态（因为基础服务接口没有批量更新方法）
	for _, id := range req.IDs {
		if err := h.service.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
			response.ErrorWithMessage(c, response.CodeInternalServerError, fmt.Sprintf("更新ID %d 状态失败: %v", id, err))
			return
		}
	}

	// 清除缓存
	if h.options.EnableCache && h.cache != nil {
		for _, id := range req.IDs {
			cacheKey := fmt.Sprintf("entity:%d", id)
			h.cache.Delete(cacheKey)
		}
	}

	response.Success(c, nil)
}

// GetMetadata 获取处理器元数据
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) GetMetadata() *HandlerMetadata {
	return h.metadata
}

// HealthCheck 健康检查
func (h *BaseHandler[T, CreateReq, UpdateReq, ListReq, Response]) HealthCheck(c *gin.Context) {
	// 简单的健康检查，尝试调用列表方法
	pageQuery := &types.BasePageQuery{Page: 1, PageSize: 1}
	if _, err := h.service.ListWithFilter(c.Request.Context(), pageQuery); err != nil {
		response.ErrorWithMessage(c, response.CodeServiceUnavailable, "服务不可用")
		return
	}

	// 检查缓存连接（如果启用）
	if h.options.EnableCache && h.cache != nil {
		if err := h.cache.Set("health_check", "ok", 1); err != nil {
			response.ErrorWithMessage(c, response.CodeServiceUnavailable, "缓存服务不可用")
			return
		}
		h.cache.Delete("health_check")
	}

	response.Success(c, gin.H{
		"status":  "healthy",
		"service": "available",
		"cache":   h.options.EnableCache,
	})
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}

// BatchStatusRequest 批量状态更新请求
type BatchStatusRequest struct {
	IDs    []uint `json:"ids" binding:"required,min=1"`
	Status int    `json:"status" binding:"oneof=0 1 -1"`
}
