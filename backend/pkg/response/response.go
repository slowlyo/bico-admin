package response

import (
	"github.com/gofiber/fiber/v2"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginationResponse 分页响应结构 - 符合Ant Design Pro标准
type PaginationResponse struct {
	Data     interface{} `json:"data"`
	Total    int64       `json:"total"`
	Success  bool        `json:"success"`
	Current  int         `json:"current"`
	PageSize int         `json:"pageSize"`
}

// LegacyPaginationResponse 传统分页响应结构（兼容性）
type LegacyPaginationResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
}

// 响应状态码常量
const (
	CodeSuccess             = 200
	CodeBadRequest          = 400
	CodeUnauthorized        = 401
	CodeForbidden           = 403
	CodeNotFound            = 404
	CodeMethodNotAllowed    = 405
	CodeConflict            = 409
	CodeUnprocessableEntity = 422
	CodeTooManyRequests     = 429
	CodeInternalServerError = 500
)

// Success 成功响应
func Success(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    CodeSuccess,
		Message: "Success",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 400错误
func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, message)
}

// Unauthorized 401错误
func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message)
}

// Forbidden 403错误
func Forbidden(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, message)
}

// NotFound 404错误
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message)
}

// Conflict 409错误
func Conflict(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusConflict, message)
}

// UnprocessableEntity 422错误
func UnprocessableEntity(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnprocessableEntity, message)
}

// InternalServerError 500错误
func InternalServerError(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusInternalServerError, message)
}

// Pagination 分页响应 - 符合Ant Design Pro标准
func Pagination(c *fiber.Ctx, data interface{}, total int64, current, pageSize int) error {
	return c.Status(fiber.StatusOK).JSON(PaginationResponse{
		Data:     data,
		Total:    total,
		Success:  true,
		Current:  current,
		PageSize: pageSize,
	})
}

// LegacyPagination 传统分页响应（兼容性）
func LegacyPagination(c *fiber.Ctx, data interface{}, total int64, page, size int) error {
	return c.Status(fiber.StatusOK).JSON(LegacyPaginationResponse{
		Code:    CodeSuccess,
		Message: "Success",
		Data:    data,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// ValidationError 验证错误响应
func ValidationError(c *fiber.Ctx, errors interface{}) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{
		Code:    CodeUnprocessableEntity,
		Message: "Validation failed",
		Data:    errors,
	})
}
