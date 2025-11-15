package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(data interface{}) *Response {
	return &Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
}

// Error 错误响应
func Error(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
	}
}

// PageData 分页数据
type PageData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// SuccessWithData 成功响应（直接输出到context）
func SuccessWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Success(data))
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, &Response{
		Code: 0,
		Msg:  msg,
		Data: data,
	})
}

// ErrorWithCode 错误响应（直接输出到context）
func ErrorWithCode(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Error(code, msg))
}

// ErrorWithStatus 错误响应（自定义HTTP状态码）
func ErrorWithStatus(c *gin.Context, httpStatus int, code int, msg string) {
	c.JSON(httpStatus, Error(code, msg))
}

// BadRequest 400错误
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Error(400, msg))
}

// NotFound 404错误
func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Error(404, msg))
}

// TooManyRequests 429限流错误
func TooManyRequests(c *gin.Context, msg string) {
	c.JSON(http.StatusTooManyRequests, Error(429, msg))
}

// SuccessWithPagination 分页成功响应（data直接是列表，total在外层）
func SuccessWithPagination(c *gin.Context, list interface{}, total int64) {
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "success",
		"data":  list,
		"total": total,
	})
}
