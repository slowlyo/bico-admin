//go:build embed
// +build embed

package web

import "embed"

// DistFS 前端构建产物。
//
// all: 保留构建生成的下划线文件，默认嵌入规则会忽略这些文件。
// 产物目录由 apps/admin 的 vite outDir 固定到仓库的 web/dist。
//go:embed all:dist
var DistFS embed.FS
