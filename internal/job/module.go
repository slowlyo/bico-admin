package job

import (
	"fmt"

	"bico-admin/internal/core/app"
)

// Module job 模块
type Module struct{}

// NewModule 创建 job 模块
func NewModule() *Module {
	return &Module{}
}

// Name 模块名称
func (m *Module) Name() string {
	return "job"
}

// Register 注册任务
func (m *Module) Register(ctx *app.AppContext) error {
	if err := RegisterJobs(ctx.Scheduler, ctx.DB, ctx.Cache, ctx.Logger); err != nil {
		return fmt.Errorf("register jobs failed: %w", err)
	}
	return nil
}
