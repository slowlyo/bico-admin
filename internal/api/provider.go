package api

import (
	"github.com/google/wire"

	"bico-admin/internal/api/handler"
	"bico-admin/internal/api/routes"
)

// ProviderSet API端Provider集合
var ProviderSet = wire.NewSet(
	// Handler层
	handler.NewHelloHandler,

	// 路由处理器集合
	ProvideHandlers,
)

// ProvideHandlers 提供处理器集合
func ProvideHandlers(
	helloHandler *handler.HelloHandler,
) *routes.Handlers {
	return &routes.Handlers{
		HelloHandler: helloHandler,
	}
}
