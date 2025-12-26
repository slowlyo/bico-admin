package api

import "bico-admin/internal/core/app"

// Module api 模块
type Module struct{}

// NewModule 创建 api 模块
func NewModule() *Module {
	return &Module{}
}

// Name 模块名称
func (m *Module) Name() string {
	return "api"
}

// Register 注册 api 路由
func (m *Module) Register(ctx *app.AppContext) error {
	return nil
}
