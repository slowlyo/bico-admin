package service

import "bico-admin/internal/core/config"

// IConfigService 配置服务接口
type IConfigService interface {
	GetAppConfig() *AppConfigResponse
}

// ConfigService 配置服务
type ConfigService struct {
	configManager *config.ConfigManager
}

// AppConfigResponse 应用配置响应
type AppConfigResponse struct {
	Name  string `json:"name"`
	Logo  string `json:"logo"`
	Debug bool   `json:"debug"`
}

// NewConfigService 创建配置服务
func NewConfigService(cm *config.ConfigManager) *ConfigService {
	return &ConfigService{
		configManager: cm,
	}
}

// GetAppConfig 获取应用配置（支持热更新）
func (s *ConfigService) GetAppConfig() *AppConfigResponse {
	cfg := s.configManager.GetConfig()
	return &AppConfigResponse{
		Name:  cfg.App.Name,
		Logo:  cfg.App.Logo,
		Debug: cfg.Server.Mode == "debug",
	}
}
