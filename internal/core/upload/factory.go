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
		// TODO: 实现七牛云上传器
		return nil, ErrInvalidDriver
	case "aliyun":
		// TODO: 实现阿里云OSS上传器
		return nil, ErrInvalidDriver
	case "tencent":
		// TODO: 实现腾讯云COS上传器
		return nil, ErrInvalidDriver
	default:
		return nil, ErrInvalidDriver
	}
}
