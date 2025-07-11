package handler

import (
	"github.com/gin-gonic/gin"

	"bico-admin/pkg/response"
)

// HelloHandler API端Hello处理器
type HelloHandler struct{}

// NewHelloHandler 创建Hello处理器
func NewHelloHandler() *HelloHandler {
	return &HelloHandler{}
}

// Hello API端Hello接口
// @Summary API端Hello
// @Description API端简单的Hello接口
// @Tags API
// @Produce json
// @Success 200 {object} response.ApiResponse
// @Router /api/hello [get]
func (h *HelloHandler) Hello(c *gin.Context) {
	response.Success(c, gin.H{
		"message": "Hello from API!",
		"module":  "api",
		"status":  "running",
	})
}
