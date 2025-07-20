//go:build embed
// +build embed

package web

import (
	"embed"
	"io/fs"
	"net/http"
)

// 嵌入前端文件（仅在 embed 构建标签时生效）
//go:embed all:dist
var EmbeddedFiles embed.FS

// GetFileSystem 获取前端文件系统
func GetFileSystem() (http.FileSystem, error) {
	// 获取 dist 子目录
	distFS, err := fs.Sub(EmbeddedFiles, "dist")
	if err != nil {
		return nil, err
	}
	return http.FS(distFS), nil
}

// IsEmbedded 检查是否使用嵌入模式
func IsEmbedded() bool {
	return true
}
