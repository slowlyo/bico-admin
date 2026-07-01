package admin

import (
	adminmodule "bico-admin/internal/admin"
	swaggercrud "bico-admin/internal/pkg/swagger"
)

// init 在 Swagger 文档注册后补齐声明式 CRUD 路由。
//
// 说明：swag 注册的是 SwaggerInfoadmin 指针，更新模板字段后 UI 会读取增强后的内容。
func init() {
	modules := adminmodule.NewCRUDModules(nil, nil)
	crudModules := make([]swaggercrud.CRUDModule, 0, len(modules))
	for _, module := range modules {
		if module == nil {
			continue
		}
		crudModules = append(crudModules, swaggercrud.CRUDModule{
			BasePath: "",
			Config:   module.ModuleConfig(),
		})
	}

	doc, err := swaggercrud.ApplyCRUDModules(SwaggerInfoadmin.ReadDoc(), crudModules)
	if err != nil {
		panic(err)
	}
	SwaggerInfoadmin.SwaggerTemplate = doc
}
