package handler

import (
	"net/http"
	
	"bico-admin/internal/admin/service"
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
	
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": config,
	})
}
