package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ApiResponse 统一API响应结构
type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageResponse 分页响应结构
type PageResponse struct {
	List       interface{} `json:"list"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// 响应码定义
const (
	// 成功
	CodeSuccess = 200

	// 客户端错误 4xx
	CodeBadRequest          = 400
	CodeUnauthorized        = 401
	CodeForbidden           = 403
	CodeNotFound            = 404
	CodeMethodNotAllowed    = 405
	CodeConflict            = 409
	CodeUnprocessableEntity = 422
	CodeTooManyRequests     = 429

	// 服务器错误 5xx
	CodeInternalServerError = 500
	CodeBadGateway          = 502
	CodeServiceUnavailable  = 503
	CodeGatewayTimeout      = 504

	// 业务错误码 1xxx
	CodeValidationError = 1001
	CodeDuplicateError  = 1002
	CodeNotExistsError  = 1003
	CodePermissionError = 1004
	CodeTokenError      = 1005
	CodePasswordError   = 1006
)

// 错误消息映射
var codeMessages = map[int]string{
	CodeSuccess:             "操作成功",
	CodeBadRequest:          "请求参数错误",
	CodeUnauthorized:        "未授权访问",
	CodeForbidden:           "禁止访问",
	CodeNotFound:            "资源不存在",
	CodeMethodNotAllowed:    "请求方法不允许",
	CodeConflict:            "资源冲突",
	CodeUnprocessableEntity: "请求参数验证失败",
	CodeTooManyRequests:     "请求过于频繁",
	CodeInternalServerError: "服务器内部错误",
	CodeBadGateway:          "网关错误",
	CodeServiceUnavailable:  "服务不可用",
	CodeGatewayTimeout:      "网关超时",
	CodeValidationError:     "数据验证失败",
	CodeDuplicateError:      "数据已存在",
	CodeNotExistsError:      "数据不存在",
	CodePermissionError:     "权限不足",
	CodeTokenError:          "令牌无效",
	CodePasswordError:       "密码错误",
}

// GetMessage 获取错误码对应的消息
func GetMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ApiResponse{
		Code:    CodeSuccess,
		Message: GetMessage(CodeSuccess),
		Data:    data,
	})
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, ApiResponse{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int) {
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, ApiResponse{
		Code:    code,
		Message: GetMessage(code),
	})
}

// ErrorWithMessage 带自定义消息的错误响应
func ErrorWithMessage(c *gin.Context, code int, message string) {
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, ApiResponse{
		Code:    code,
		Message: message,
	})
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, ApiResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Page 分页响应
func Page(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, ApiResponse{
		Code:    CodeSuccess,
		Message: GetMessage(CodeSuccess),
		Data: PageResponse{
			List:       list,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	})
}

// getHTTPStatus 根据业务错误码获取HTTP状态码
func getHTTPStatus(code int) int {
	switch {
	case code >= 400 && code < 500:
		return code
	case code >= 500 && code < 600:
		return code
	case code >= 1000 && code < 2000:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// BadRequest 400错误
func BadRequest(c *gin.Context, message ...string) {
	msg := GetMessage(CodeBadRequest)
	if len(message) > 0 {
		msg = message[0]
	}
	ErrorWithMessage(c, CodeBadRequest, msg)
}

// Unauthorized 401错误
func Unauthorized(c *gin.Context, message ...string) {
	msg := GetMessage(CodeUnauthorized)
	if len(message) > 0 {
		msg = message[0]
	}
	ErrorWithMessage(c, CodeUnauthorized, msg)
}

// Forbidden 403错误
func Forbidden(c *gin.Context, message ...string) {
	msg := GetMessage(CodeForbidden)
	if len(message) > 0 {
		msg = message[0]
	}
	ErrorWithMessage(c, CodeForbidden, msg)
}

// NotFound 404错误
func NotFound(c *gin.Context, message ...string) {
	msg := GetMessage(CodeNotFound)
	if len(message) > 0 {
		msg = message[0]
	}
	ErrorWithMessage(c, CodeNotFound, msg)
}

// InternalServerError 500错误
func InternalServerError(c *gin.Context, message ...string) {
	msg := GetMessage(CodeInternalServerError)
	if len(message) > 0 {
		msg = message[0]
	}
	ErrorWithMessage(c, CodeInternalServerError, msg)
}

// ValidationError 验证错误
func ValidationError(c *gin.Context, message ...string) {
	msg := GetMessage(CodeValidationError)
	if len(message) > 0 {
		msg = message[0]
	}
	ErrorWithMessage(c, CodeValidationError, msg)
}
