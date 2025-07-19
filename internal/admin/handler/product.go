package handler

import (
	"time"

	"github.com/gin-gonic/gin"

	"bico-admin/internal/admin/service"
	"bico-admin/internal/admin/types"
	"bico-admin/internal/shared/models"
	sharedTypes "bico-admin/internal/shared/types"
	"bico-admin/pkg/response"
	"bico-admin/pkg/utils"
)

// ProductHandler Product处理器
type ProductHandler struct {
	productService service.ProductService
}

// NewProductHandler 创建Product处理器
func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// Create 创建Product
func (h *ProductHandler) Create(c *gin.Context) {
	var req types.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 自定义验证
	if err := req.Validate(); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用服务层
	result, err := h.productService.Create(c.Request.Context(), &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	// 转换响应
	responseData := h.convertToResponse(result)
	response.Success(c, responseData)
}

// GetByID 根据ID获取Product
func (h *ProductHandler) GetByID(c *gin.Context) {
	var req sharedTypes.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用服务层
	entity, err := h.productService.GetByID(c.Request.Context(), req.ID)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeNotFound, err.Error())
		return
	}

	// 转换响应
	responseData := h.convertToResponse(entity)
	response.Success(c, responseData)
}

// Update 更新Product
func (h *ProductHandler) Update(c *gin.Context) {
	var uriReq sharedTypes.IDRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	var req types.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 自定义验证
	if err := req.Validate(); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用服务层
	result, err := h.productService.Update(c.Request.Context(), uriReq.ID, &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	// 转换响应
	responseData := h.convertToResponse(result)
	response.Success(c, responseData)
}

// Delete 删除Product
func (h *ProductHandler) Delete(c *gin.Context) {
	var req sharedTypes.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用服务层
	if err := h.productService.Delete(c.Request.Context(), req.ID); err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetList 获取Product列表
func (h *ProductHandler) GetList(c *gin.Context) {
	var req types.ProductListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 自定义验证
	if err := req.Validate(); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// 调用服务层
	entities, total, err := h.productService.ListWithFilter(c.Request.Context(), &req)
	if err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	// 转换响应
	var responses []types.ProductResponse
	for _, entity := range entities {
		responses = append(responses, h.convertToResponse(entity))
	}

	response.Page(c, responses, total, req.GetPage(), req.GetPageSize())
}

// UpdateStatus 更新Product状态
func (h *ProductHandler) UpdateStatus(c *gin.Context) {
	var uriReq sharedTypes.IDRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	var req sharedTypes.StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.productService.UpdateStatus(c.Request.Context(), uriReq.ID, req.Status); err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// convertToResponse 转换为响应格式
func (h *ProductHandler) convertToResponse(entity *models.Product) types.ProductResponse {

	return types.ProductResponse{
		ID:          entity.ID,
		Name:        entity.Name,
		SKU:         entity.SKU,
		Description: entity.Description,
		Price:       entity.Price,
		Stock:       entity.Stock,
		CategoryID:  entity.CategoryID,
		BrandID:     entity.BrandID,
		Images:      entity.Images,
		Attributes:  entity.Attributes,
		Status:      entity.Status,
		StatusText:  h.getStatusText(entity.Status),
		Weight:      entity.Weight,
		PublishedAt: h.formatTime(entity.PublishedAt),
		ExpiredAt:   h.formatTime(entity.ExpiredAt),
		CreatedAt:   utils.NewFormattedTime(entity.CreatedAt),
		UpdatedAt:   utils.NewFormattedTime(entity.UpdatedAt),
	}
}

// getStatusText 获取状态文本
func (h *ProductHandler) getStatusText(status int) string {
	switch status {
	case sharedTypes.StatusActive:
		return "启用"
	case sharedTypes.StatusInactive:
		return "禁用"
	case sharedTypes.StatusDeleted:
		return "已删除"
	default:
		return "未知"
	}
}

// formatTime 格式化时间
func (h *ProductHandler) formatTime(t *time.Time) *utils.FormattedTime {
	if t == nil {
		return nil
	}
	ft := utils.NewFormattedTime(*t)
	return &ft
}
