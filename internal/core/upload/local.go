package upload

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalUploader 本地存储上传器
type LocalUploader struct {
	basePath     string
	urlPrefix    string
	maxSize      int64
	allowedTypes []string
}

// NewLocalUploader 创建本地上传器
func NewLocalUploader(basePath, urlPrefix string, maxSize int64, allowedTypes []string) *LocalUploader {
	return &LocalUploader{
		basePath:     basePath,
		urlPrefix:    urlPrefix,
		maxSize:      maxSize,
		allowedTypes: allowedTypes,
	}
}

// Upload 上传文件
func (u *LocalUploader) Upload(file *multipart.FileHeader, subPath string) (string, error) {
	if err := u.ValidateFile(file); err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", ErrUploadFailed
	}
	defer src.Close()

	filename := u.generateFilename(file.Filename)
	
	fullSubPath := filepath.Join(u.basePath, subPath)
	if err := os.MkdirAll(fullSubPath, 0755); err != nil {
		return "", ErrUploadFailed
	}

	filePath := filepath.Join(fullSubPath, filename)
	
	dst, err := os.Create(filePath)
	if err != nil {
		return "", ErrUploadFailed
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", ErrUploadFailed
	}

	// 构建返回的 URL（支持完整 URL 和相对路径）
	urlPath := filepath.ToSlash(filepath.Join(subPath, filename))
	if u.urlPrefix == "" {
		return "/" + urlPath, nil
	}
	
	// 如果 urlPrefix 以 / 结尾，去掉末尾的 /
	prefix := u.urlPrefix
	if len(prefix) > 0 && prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}
	
	return prefix + "/" + urlPath, nil
}

// Delete 删除文件
func (u *LocalUploader) Delete(url string) error {
	if url == "" {
		return nil
	}

	relPath := strings.TrimPrefix(url, u.urlPrefix)
	filePath := filepath.Join(u.basePath, relPath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(filePath)
}

// ValidateFile 验证文件
func (u *LocalUploader) ValidateFile(file *multipart.FileHeader) error {
	if file.Size > u.maxSize {
		return ErrFileTooLarge
	}

	contentType := file.Header.Get("Content-Type")
	if !u.isAllowedType(contentType) {
		return ErrInvalidFileType
	}

	return nil
}

// isAllowedType 检查文件类型是否允许
func (u *LocalUploader) isAllowedType(contentType string) bool {
	for _, allowed := range u.allowedTypes {
		if allowed == contentType {
			return true
		}
	}
	return false
}

// generateFilename 生成唯一文件名
func (u *LocalUploader) generateFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	
	timestamp := time.Now().Format("20060102150405")
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%s%d", originalFilename, time.Now().UnixNano())))
	hashString := hex.EncodeToString(hash.Sum(nil))[:8]
	
	return fmt.Sprintf("%s_%s%s", timestamp, hashString, ext)
}
