//go:build embed
// +build embed

package web

import "embed"

// DistFS 前端构建产物
//
//go:embed all:dist
var DistFS embed.FS
