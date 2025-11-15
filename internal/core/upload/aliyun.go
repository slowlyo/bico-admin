package upload

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// AliyunUploader 阿里云OSS上传器
type AliyunUploader struct {
	client       *oss.Client
	bucket       *oss.Bucket
	bucketName   string
	domain       string
	useHTTPS     bool
	maxSize      int64
	allowedTypes []string
}

// AliyunConfig 阿里云OSS配置
type AliyunConfig struct {
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	Endpoint        string
	Domain          string
	UseHTTPS        bool
}

// NewAliyunUploader 创建阿里云OSS上传器
func NewAliyunUploader(config AliyunConfig, maxSize int64, allowedTypes []string) (*AliyunUploader, error) {
	if config.AccessKeyId == "" || config.AccessKeySecret == "" {
		return nil, fmt.Errorf("阿里云 AccessKeyId 和 AccessKeySecret 不能为空")
	}
	if config.Bucket == "" {
		return nil, fmt.Errorf("阿里云 Bucket 不能为空")
	}
	if config.Endpoint == "" {
		return nil, fmt.Errorf("阿里云 Endpoint 不能为空")
	}

	client, err := oss.New(config.Endpoint, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建阿里云OSS客户端失败: %w", err)
	}

	bucket, err := client.Bucket(config.Bucket)
	if err != nil {
		return nil, fmt.Errorf("获取阿里云OSS存储空间失败: %w", err)
	}

	domain := config.Domain
	if domain == "" {
		domain = fmt.Sprintf("%s.%s", config.Bucket, config.Endpoint)
	}

	return &AliyunUploader{
		client:       client,
		bucket:       bucket,
		bucketName:   config.Bucket,
		domain:       domain,
		useHTTPS:     config.UseHTTPS,
		maxSize:      maxSize,
		allowedTypes: allowedTypes,
	}, nil
}

// Upload 上传文件
func (u *AliyunUploader) Upload(file *multipart.FileHeader, subPath string) (string, error) {
	if err := u.ValidateFile(file); err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", ErrUploadFailed
	}
	defer src.Close()

	filename := u.generateFilename(file.Filename)
	key := filepath.ToSlash(filepath.Join(subPath, filename))

	err = u.bucket.PutObject(key, src)
	if err != nil {
		return "", fmt.Errorf("阿里云OSS上传失败: %w", err)
	}

	return u.buildURL(key), nil
}

// Delete 删除文件
func (u *AliyunUploader) Delete(url string) error {
	key := u.extractKeyFromURL(url)
	if key == "" {
		return fmt.Errorf("无效的文件 URL")
	}

	err := u.bucket.DeleteObject(key)
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

// ValidateFile 验证文件
func (u *AliyunUploader) ValidateFile(file *multipart.FileHeader) error {
	if file.Size > u.maxSize {
		return ErrFileTooLarge
	}

	if len(u.allowedTypes) > 0 {
		contentType := file.Header.Get("Content-Type")
		allowed := false
		for _, t := range u.allowedTypes {
			if t == contentType {
				allowed = true
				break
			}
		}
		if !allowed {
			return ErrInvalidFileType
		}
	}

	return nil
}

// generateFilename 生成文件名
func (u *AliyunUploader) generateFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d%s", timestamp, ext)
}

// buildURL 构建完整的访问 URL
func (u *AliyunUploader) buildURL(key string) string {
	domain := u.domain
	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		if u.useHTTPS {
			domain = "https://" + domain
		} else {
			domain = "http://" + domain
		}
	}

	domain = strings.TrimRight(domain, "/")
	return fmt.Sprintf("%s/%s", domain, key)
}

// extractKeyFromURL 从 URL 中提取 key
func (u *AliyunUploader) extractKeyFromURL(url string) string {
	domain := strings.TrimRight(u.domain, "/")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")

	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	if strings.HasPrefix(url, domain+"/") {
		return strings.TrimPrefix(url, domain+"/")
	}

	return ""
}
