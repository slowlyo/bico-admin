package upload

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// QiniuUploader 七牛云上传器
type QiniuUploader struct {
	accessKey    string
	secretKey    string
	bucket       string
	domain       string
	zone         *storage.Zone
	useHTTPS     bool
	useCDNDomain bool
	maxSize      int64
	allowedTypes []string
}

// QiniuConfig 七牛云配置
type QiniuConfig struct {
	AccessKey    string
	SecretKey    string
	Bucket       string
	Domain       string
	Zone         string
	UseHTTPS     bool
	UseCDNDomain bool
}

// NewQiniuUploader 创建七牛云上传器
func NewQiniuUploader(config QiniuConfig, maxSize int64, allowedTypes []string) (*QiniuUploader, error) {
	if config.AccessKey == "" || config.SecretKey == "" {
		return nil, fmt.Errorf("七牛云 AccessKey 和 SecretKey 不能为空")
	}
	if config.Bucket == "" {
		return nil, fmt.Errorf("七牛云 Bucket 不能为空")
	}
	if config.Domain == "" {
		return nil, fmt.Errorf("七牛云 Domain 不能为空")
	}

	zone := parseZone(config.Zone)

	return &QiniuUploader{
		accessKey:    config.AccessKey,
		secretKey:    config.SecretKey,
		bucket:       config.Bucket,
		domain:       config.Domain,
		zone:         zone,
		useHTTPS:     config.UseHTTPS,
		useCDNDomain: config.UseCDNDomain,
		maxSize:      maxSize,
		allowedTypes: allowedTypes,
	}, nil
}

// parseZone 解析存储区域
func parseZone(zone string) *storage.Zone {
	switch zone {
	case "z0":
		return &storage.ZoneHuadong
	case "z1":
		return &storage.ZoneHuabei
	case "z2":
		return &storage.ZoneHuanan
	case "na0":
		return &storage.ZoneBeimei
	case "as0":
		return &storage.ZoneXinjiapo
	default:
		return &storage.ZoneHuadong
	}
}

// Upload 上传文件
func (u *QiniuUploader) Upload(file *multipart.FileHeader, subPath string) (string, error) {
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

	mac := qbox.NewMac(u.accessKey, u.secretKey)
	putPolicy := storage.PutPolicy{
		Scope: u.bucket,
	}
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{
		Zone:          u.zone,
		UseHTTPS:      u.useHTTPS,
		UseCdnDomains: u.useCDNDomain,
	}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	putExtra := storage.PutExtra{}

	err = formUploader.Put(context.Background(), &ret, upToken, key, src, file.Size, &putExtra)
	if err != nil {
		return "", fmt.Errorf("七牛云上传失败: %w", err)
	}

	return u.buildURL(key), nil
}

// Delete 删除文件
func (u *QiniuUploader) Delete(url string) error {
	key := u.extractKeyFromURL(url)
	if key == "" {
		return fmt.Errorf("无效的文件 URL")
	}

	mac := qbox.NewMac(u.accessKey, u.secretKey)
	cfg := storage.Config{
		Zone:          u.zone,
		UseHTTPS:      u.useHTTPS,
		UseCdnDomains: u.useCDNDomain,
	}

	bucketManager := storage.NewBucketManager(mac, &cfg)
	err := bucketManager.Delete(u.bucket, key)
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

// ValidateFile 验证文件
func (u *QiniuUploader) ValidateFile(file *multipart.FileHeader) error {
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
func (u *QiniuUploader) generateFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d%s", timestamp, ext)
}

// buildURL 构建完整的访问 URL
func (u *QiniuUploader) buildURL(key string) string {
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
func (u *QiniuUploader) extractKeyFromURL(url string) string {
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
