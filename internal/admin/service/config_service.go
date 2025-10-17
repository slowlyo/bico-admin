package service

import "bico-admin/internal/core/config"

// IConfigService 配置服务接口
type IConfigService interface {
	GetAppConfig() *AppConfigResponse
}

// ConfigService 配置服务
type ConfigService struct {
	config *config.Config
}

// AppConfigResponse 应用配置响应
type AppConfigResponse struct {
	Name  string `json:"name"`
	Logo  string `json:"logo"`
	Debug bool   `json:"debug"`
}

// NewConfigService 创建配置服务
func NewConfigService(config *config.Config) *ConfigService {
	return &ConfigService{
		config: config,
	}
}

// GetAppConfig 获取应用配置
func (s *ConfigService) GetAppConfig() *AppConfigResponse {
	return &AppConfigResponse{
		Name:  s.config.App.Name,
		Logo:  s.config.App.Logo,
		Debug: s.config.Server.Mode == "debug",
	}
}
