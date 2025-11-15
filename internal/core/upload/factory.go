package upload

import (
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
