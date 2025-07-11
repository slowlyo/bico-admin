package master

import (
	"github.com/google/wire"

	"bico-admin/internal/master/handler"
	"bico-admin/internal/master/routes"
)

// ProviderSet Master端Provider集合
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
