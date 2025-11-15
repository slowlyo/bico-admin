package upload

import (
	"bico-admin/internal/core/config"
	"errors"
)

var (
	ErrInvalidDriver = errors.New("不支持的上传驱动")
)

// UploaderConfig 上传器配置
type UploaderConfig struct {
	Driver       string
	MaxSize      int64
	AllowedTypes []string
	LocalConfig  LocalConfig
	QiniuConfig  QiniuConfig
	AliyunConfig AliyunConfig
}

// LocalConfig 本地存储配置
type LocalConfig struct {
	BasePath  string
	URLPrefix string
}

// ConfigFromAppConfig 从应用配置转换为上传器配置
func ConfigFromAppConfig(cfg *config.Config) *UploaderConfig {
	return &UploaderConfig{
		Driver:       cfg.Upload.Driver,
		MaxSize:      cfg.Upload.MaxSize,
		AllowedTypes: cfg.Upload.AllowedTypes,
		LocalConfig: LocalConfig{
			BasePath:  cfg.Upload.Local.BasePath,
			URLPrefix: cfg.Upload.Local.URLPrefix,
		},
		QiniuConfig: QiniuConfig{
			AccessKey:    cfg.Upload.Qiniu.AccessKey,
			SecretKey:    cfg.Upload.Qiniu.SecretKey,
			Bucket:       cfg.Upload.Qiniu.Bucket,
			Domain:       cfg.Upload.Qiniu.Domain,
			Zone:         cfg.Upload.Qiniu.Zone,
			UseHTTPS:     cfg.Upload.Qiniu.UseHTTPS,
			UseCDNDomain: cfg.Upload.Qiniu.UseCDNDomain,
		},
		AliyunConfig: AliyunConfig{
			AccessKeyId:     cfg.Upload.Aliyun.AccessKeyId,
			AccessKeySecret: cfg.Upload.Aliyun.AccessKeySecret,
			Bucket:          cfg.Upload.Aliyun.Bucket,
			Endpoint:        cfg.Upload.Aliyun.Endpoint,
			Domain:          cfg.Upload.Aliyun.Domain,
			UseHTTPS:        cfg.Upload.Aliyun.UseHTTPS,
		},
	}
}

// NewUploader 创建上传器
func NewUploader(config *UploaderConfig) (Uploader, error) {
	switch config.Driver {
	case "local":
		return NewLocalUploader(
			config.LocalConfig.BasePath,
			config.LocalConfig.URLPrefix,
			config.MaxSize,
			config.AllowedTypes,
		), nil
	case "qiniu":
		return NewQiniuUploader(
			config.QiniuConfig,
			config.MaxSize,
			config.AllowedTypes,
		)
	case "aliyun":
		return NewAliyunUploader(
			config.AliyunConfig,
			config.MaxSize,
			config.AllowedTypes,
		)
	default:
		return nil, ErrInvalidDriver
	}
}
