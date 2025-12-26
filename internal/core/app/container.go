package app

import (
	"errors"

	"go.uber.org/dig"
)

// BuildContainer 构建 DI 容器
//
// 说明：项目已改造为 core 仅负责基础设施装配，业务模块自行处理 DI。
// 该函数保留是为了避免误用时产生隐式循环依赖；请改用 BuildContext。
func BuildContainer(configPath string) (*dig.Container, error) {
	_ = configPath
	return nil, errors.New("BuildContainer 已废弃，请使用 BuildContext 并通过模块 Register 自行装配")
}
