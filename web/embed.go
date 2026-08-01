//go:build embed
// +build embed

package web

import "embed"

// DistFS 前端构建产物
//
// all: 保留 Umi 生成的 _*.async.js 分包，默认嵌入规则会忽略以下划线开头的文件。
//go:embed all:dist
var DistFS embed.FS
