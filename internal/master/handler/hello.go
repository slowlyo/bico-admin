package handler

import (
	"github.com/gin-gonic/gin"

	"bico-admin/pkg/response"
)

// HelloHandler 主控端Hello处理器
type HelloHandler struct{}

// NewHelloHandler 创建Hello处理器
func NewHelloHandler() *HelloHandler {
	return &HelloHandler{}
}

// Hello 主控端Hello接口
// @Summary 主控端Hello
// @Description 主控端简单的Hello接口
// @Tags 主控端
// @Produce json
// @Success 200 {object} response.ApiResponse
// @Router /master/hello [get]
func (h *HelloHandler) Hello(c *gin.Context) {
	response.Success(c, gin.H{
		"message": "Hello from Master!",
		"module":  "master",
		"status":  "running",
	})
}
