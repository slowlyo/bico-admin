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
// @Summary 获取应用配置
// @Description 获取后台名称、Logo 和调试模式状态
// @Tags 公共
// @Produce json
// @Success 200 {object} adminResponse{data=appConfigDocResponse}
// @Router /app-config [get]
func (h *CommonHandler) GetAppConfig(c *gin.Context) {
	config := h.configService.GetAppConfig()
	response.SuccessWithData(c, config)
}
