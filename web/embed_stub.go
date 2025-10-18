//go:build !embed
// +build !embed

package web

import "embed"

// DistFS 前端构建产物（stub 版本，用于非嵌入模式）
var DistFS embed.FS
