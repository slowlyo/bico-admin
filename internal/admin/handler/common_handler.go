package handler

import (
	"bico-admin/internal/admin/service"
	"bico-admin/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// CommonHandler 公共处理器
type CommonHandler struct {
	configService service.IConfigService
}

// NewCommonHandler 创建公共处理器
func NewCommonHandler(configService service.IConfigService) *CommonHandler {
	return &CommonHandler{
		configService: configService,
	}
}

// GetAppConfig 获取应用配置
func (h *CommonHandler) GetAppConfig(c *gin.Context) {
	config := h.configService.GetAppConfig()
	response.SuccessWithData(c, config)
}
