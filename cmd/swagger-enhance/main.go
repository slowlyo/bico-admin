package main

import (
	"encoding/json"
	"os"

	adminmodule "bico-admin/internal/admin"
	swaggercrud "bico-admin/internal/pkg/swagger"

	"go.yaml.in/yaml/v3"
)

// main 增强 swag 生成的静态文档文件。
//
// 说明：运行时 Swagger UI 会读取增强后的模板；这里同步处理落盘 json/yaml，便于外部工具直接消费。
func main() {
	if err := enhanceAdminDocs(); err != nil {
		panic(err)
	}
}

// enhanceAdminDocs 将后台 CRUD 模块写入 admin Swagger 文件。
func enhanceAdminDocs() error {
	content, err := os.ReadFile("docs/admin/admin_swagger.json")
	if err != nil {
		return err
	}

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

	doc, err := swaggercrud.ApplyCRUDModules(string(content), crudModules)
	if err != nil {
		return err
	}
	if err := os.WriteFile("docs/admin/admin_swagger.json", []byte(doc), 0o644); err != nil {
		return err
	}

	var data interface{}
	if err := json.Unmarshal([]byte(doc), &data); err != nil {
		return err
	}
	yamlContent, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile("docs/admin/admin_swagger.yaml", yamlContent, 0o644)
}
