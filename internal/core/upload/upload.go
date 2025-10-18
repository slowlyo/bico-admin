package upload

import (
	"errors"
	"mime/multipart"
)

var (
	ErrFileTooLarge      = errors.New("文件大小超过限制")
	ErrInvalidFileType   = errors.New("不支持的文件类型")
	ErrUploadFailed      = errors.New("文件上传失败")
)

// Uploader 文件上传接口
type Uploader interface {
	// Upload 上传文件
	Upload(file *multipart.FileHeader, subPath string) (string, error)
	// Delete 删除文件
	Delete(url string) error
	// ValidateFile 验证文件
	ValidateFile(file *multipart.FileHeader) error
}

// Config 上传配置接口
type Config interface {
	GetDriver() string
	GetMaxSize() int64
	GetAllowedTypes() []string
}
