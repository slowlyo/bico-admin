//go:build !embed
// +build !embed

package web

import (
	"net/http"
)

// GetFileSystem 获取前端文件系统（外部文件模式）
func GetFileSystem() (http.FileSystem, error) {
	// 外部文件模式下返回 nil，由调用方处理
	return nil, nil
}

// IsEmbedded 检查是否使用嵌入模式
func IsEmbedded() bool {
	return false
}
