package swagger

import (
	"encoding/json"
	"reflect"
	"strings"

	"bico-admin/internal/pkg/crud"
)

const (
	swaggerTypeObject  = "object"
	swaggerTypeArray   = "array"
	swaggerTypeString  = "string"
	swaggerTypeInteger = "integer"
	swaggerTypeNumber  = "number"
	swaggerTypeBoolean = "boolean"
)

// CRUDModule 描述一个 CRUD 模块生成 Swagger 所需的元数据。
type CRUDModule struct {
	BasePath string
	Config   crud.ModuleConfig
}

// ApplyCRUDModules 将声明式 CRUD 模块补充到 Swagger JSON 文档。
//
// 说明：swag 只能解析静态注释，无法识别运行时反射注册的 CRUD 路由，因此这里在生成文档后补齐 paths 和 definitions。
func ApplyCRUDModules(doc string, modules []CRUDModule) (string, error) {
	var spec map[string]interface{}
	if err := json.Unmarshal([]byte(doc), &spec); err != nil {
		return "", err
	}

	paths := objectMap(spec, "paths")
	definitions := objectMap(spec, "definitions")
	for _, module := range modules {
		// 模块配置不完整时跳过，避免生成不可访问的路径。
		if module.Config.Group == "" {
			continue
		}
		applyCRUDModule(paths, definitions, module)
	}

	updated, err := json.Marshal(spec)
	if err != nil {
		return "", err
	}

	return string(updated), nil
}

// applyCRUDModule 将单个模块的路由和类型定义写入 Swagger spec。
func applyCRUDModule(paths map[string]interface{}, definitions map[string]interface{}, module CRUDModule) {
	modelName := schemaName(module.Config.Swagger.Model)
	listReqName := schemaName(module.Config.Swagger.ListRequest)
	createReqName := schemaName(module.Config.Swagger.CreateRequest)
	updateReqName := schemaName(module.Config.Swagger.UpdateRequest)

	addDefinition(definitions, module.Config.Swagger.Model)
	addDefinition(definitions, module.Config.Swagger.ListRequest)
	addDefinition(definitions, module.Config.Swagger.CreateRequest)
	addDefinition(definitions, module.Config.Swagger.UpdateRequest)
	addSharedDefinitions(definitions)

	tag := moduleTag(module.Config)
	security := routeSecurity(module.Config.Routes)

	for _, route := range module.Config.Routes {
		fullPath := swaggerPath(module.BasePath, module.Config.Group, route.Path)
		pathItem := objectMap(paths, fullPath)
		method := strings.ToLower(route.Method)
		// 已有静态注释生成的 operation 时保留原文档，避免 CRUD 自动增强覆盖手写说明。
		if _, exists := pathItem[method]; exists {
			continue
		}
		operation := buildOperation(tag, module.Config, route, modelName, listReqName, createReqName, updateReqName, security)
		pathItem[method] = operation
	}
}

// buildOperation 根据 CRUD 路由语义生成单个 operation。
func buildOperation(
	tag string,
	config crud.ModuleConfig,
	route crud.Route,
	modelName string,
	listReqName string,
	createReqName string,
	updateReqName string,
	security []map[string][]string,
) map[string]interface{} {
	summary := routeSummary(config, route)
	operation := map[string]interface{}{
		"tags":        []string{tag},
		"summary":     summary,
		"description": routeDescription(route),
		"responses": map[string]interface{}{
			"200": map[string]interface{}{
				"description": "OK",
				"schema":      responseSchemaForRoute(route, modelName),
			},
		},
	}

	parameters := parametersForRoute(route, listReqName, createReqName, updateReqName)
	if len(parameters) > 0 {
		operation["parameters"] = parameters
	}
	if !route.Public && len(security) > 0 {
		// 私有路由复用全局 BearerAuth 定义，公开路由不追加安全声明。
		operation["security"] = security
	}
	return operation
}

// parametersForRoute 生成路径、查询和 body 参数。
func parametersForRoute(route crud.Route, listReqName string, createReqName string, updateReqName string) []map[string]interface{} {
	parameters := make([]map[string]interface{}, 0, 4)
	if strings.Contains(route.Path, ":id") {
		parameters = append(parameters, map[string]interface{}{
			"name":        "id",
			"in":          "path",
			"description": "记录 ID",
			"required":    true,
			"type":        swaggerTypeInteger,
			"format":      "uint",
		})
	}

	if route.Handler == "List" {
		parameters = append(parameters, paginationParameters()...)
		parameters = append(parameters, queryParametersFromSchema(listReqName)...)
	}

	if route.Handler == "Create" && createReqName != "" {
		parameters = append(parameters, bodyParameter("body", "创建参数", createReqName))
	}
	if route.Handler == "Update" && updateReqName != "" {
		parameters = append(parameters, bodyParameter("body", "更新参数", updateReqName))
	}
	return parameters
}

// paginationParameters 返回 CRUD 列表接口约定的分页参数。
func paginationParameters() []map[string]interface{} {
	return []map[string]interface{}{
		{"name": "page", "in": "query", "description": "页码", "required": false, "type": swaggerTypeInteger},
		{"name": "pageSize", "in": "query", "description": "每页数量", "required": false, "type": swaggerTypeInteger},
		{"name": "sortField", "in": "query", "description": "排序字段", "required": false, "type": swaggerTypeString},
		{"name": "sortOrder", "in": "query", "description": "排序方向：ascend 为升序，其余为降序", "required": false, "type": swaggerTypeString},
	}
}

// queryParametersFromSchema 从请求结构体 definition 中提取查询参数。
func queryParametersFromSchema(schemaName string) []map[string]interface{} {
	if schemaName == "" {
		return nil
	}
	t := typeBySchemaName(schemaName)
	if t == nil {
		return nil
	}

	parameters := make([]map[string]interface{}, 0, t.NumField())
	collectQueryParameters(t, &parameters)
	return parameters
}

// collectQueryParameters 遍历结构体字段，将 form tag 转成 Swagger query 参数。
func collectQueryParameters(t reflect.Type, parameters *[]map[string]interface{}) {
	t = indirectType(t)
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		if field.Anonymous {
			collectQueryParameters(field.Type, parameters)
			continue
		}
		name := tagName(field.Tag.Get("form"))
		if name == "" {
			continue
		}
		schema := schemaFromType(field.Type)
		parameter := map[string]interface{}{
			"name":        name,
			"in":          "query",
			"description": fieldDescription(field),
			"required":    hasRequiredBinding(field),
			"type":        schema["type"],
		}
		if format, ok := schema["format"]; ok {
			parameter["format"] = format
		}
		*parameters = append(*parameters, parameter)
	}
}

// bodyParameter 创建 body 参数引用。
func bodyParameter(name string, description string, schemaName string) map[string]interface{} {
	return map[string]interface{}{
		"name":        name,
		"in":          "body",
		"description": description,
		"required":    true,
		"schema":      refSchema(schemaName),
	}
}

// responseSchemaForRoute 根据标准 CRUD handler 选择响应 schema。
func responseSchemaForRoute(route crud.Route, modelName string) map[string]interface{} {
	if route.Handler == "List" {
		return refSchema("swagger.PageResponse")
	}
	if modelName == "" || route.Handler == "Delete" {
		return refSchema("swagger.Response")
	}
	return map[string]interface{}{
		"allOf": []map[string]interface{}{
			refSchema("swagger.Response"),
			{
				"type": swaggerTypeObject,
				"properties": map[string]interface{}{
					"data": refSchema(modelName),
				},
			},
		},
	}
}

// addDefinition 将 Go 类型转换为 Swagger definition。
func addDefinition(definitions map[string]interface{}, sample interface{}) {
	name := schemaName(sample)
	if name == "" {
		return
	}
	if _, exists := definitions[name]; exists {
		return
	}

	t := indirectType(reflect.TypeOf(sample))
	if t.Kind() != reflect.Struct {
		return
	}

	typeRegistry[name] = t
	definitions[name] = schemaForStruct(t)
}

// addSharedDefinitions 写入 CRUD 通用响应定义。
func addSharedDefinitions(definitions map[string]interface{}) {
	if _, exists := definitions["swagger.Response"]; !exists {
		definitions["swagger.Response"] = map[string]interface{}{
			"type": swaggerTypeObject,
			"properties": map[string]interface{}{
				"code": map[string]interface{}{"type": swaggerTypeInteger},
				"msg":  map[string]interface{}{"type": swaggerTypeString},
				"data": map[string]interface{}{"type": swaggerTypeObject},
			},
		}
	}
	if _, exists := definitions["swagger.PageResponse"]; !exists {
		definitions["swagger.PageResponse"] = map[string]interface{}{
			"type": swaggerTypeObject,
			"properties": map[string]interface{}{
				"code": map[string]interface{}{"type": swaggerTypeInteger},
				"msg":  map[string]interface{}{"type": swaggerTypeString},
				"data": map[string]interface{}{
					"type": swaggerTypeObject,
					"properties": map[string]interface{}{
						"list":  map[string]interface{}{"type": swaggerTypeArray, "items": map[string]interface{}{"type": swaggerTypeObject}},
						"total": map[string]interface{}{"type": swaggerTypeInteger, "format": "int64"},
					},
				},
			},
		}
	}
}

// schemaForStruct 将结构体字段转换为 Swagger object schema。
func schemaForStruct(t reflect.Type) map[string]interface{} {
	properties := make(map[string]interface{})
	required := make([]string, 0)
	collectStructProperties(t, properties, &required)

	schema := map[string]interface{}{
		"type":       swaggerTypeObject,
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

// collectStructProperties 收集结构体字段，支持匿名嵌入字段展开。
func collectStructProperties(t reflect.Type, properties map[string]interface{}, required *[]string) {
	t = indirectType(t)
	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		if field.Anonymous {
			collectStructProperties(field.Type, properties, required)
			continue
		}

		name := jsonFieldName(field)
		if name == "" {
			continue
		}

		properties[name] = schemaFromType(field.Type)
		if hasRequiredBinding(field) {
			*required = append(*required, name)
		}
	}
}

// schemaFromType 将 Go 字段类型映射为 Swagger schema。
func schemaFromType(t reflect.Type) map[string]interface{} {
	t = indirectType(t)
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		return map[string]interface{}{
			"type":  swaggerTypeArray,
			"items": schemaFromType(t.Elem()),
		}
	}
	if t.Kind() == reflect.Struct {
		if isTimeLikeType(t) {
			return map[string]interface{}{"type": swaggerTypeString, "format": "date-time"}
		}
		return schemaForStruct(t)
	}

	switch t.Kind() {
	case reflect.Bool:
		return map[string]interface{}{"type": swaggerTypeBoolean}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return map[string]interface{}{"type": swaggerTypeInteger, "format": "int32"}
	case reflect.Int64:
		return map[string]interface{}{"type": swaggerTypeInteger, "format": "int64"}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return map[string]interface{}{"type": swaggerTypeInteger, "format": "uint"}
	case reflect.Uint64:
		return map[string]interface{}{"type": swaggerTypeInteger, "format": "uint64"}
	case reflect.Float32:
		return map[string]interface{}{"type": swaggerTypeNumber, "format": "float"}
	case reflect.Float64:
		return map[string]interface{}{"type": swaggerTypeNumber, "format": "double"}
	default:
		return map[string]interface{}{"type": swaggerTypeString}
	}
}

// objectMap 获取或创建 map[string]interface{} 子节点。
func objectMap(parent map[string]interface{}, key string) map[string]interface{} {
	if value, ok := parent[key].(map[string]interface{}); ok {
		return value
	}
	value := make(map[string]interface{})
	parent[key] = value
	return value
}

// refSchema 创建 definition 引用。
func refSchema(name string) map[string]interface{} {
	return map[string]interface{}{"$ref": "#/definitions/" + name}
}

// routeSecurity 为私有路由生成 BearerAuth 安全声明。
func routeSecurity(routes []crud.Route) []map[string][]string {
	for _, route := range routes {
		if !route.Public {
			return []map[string][]string{{"BearerAuth": []string{}}}
		}
	}
	return nil
}

// swaggerPath 拼接 gin 路由路径并转换 :id 为 Swagger path 参数。
func swaggerPath(basePath string, group string, path string) string {
	full := "/" + strings.Trim(strings.Trim(basePath, "/")+"/"+strings.Trim(group, "/")+"/"+strings.Trim(path, "/"), "/")
	full = strings.ReplaceAll(full, "/:id", "/{id}")
	if full == "/" {
		return full
	}
	return strings.TrimRight(full, "/")
}

// moduleTag 返回模块文档标签。
func moduleTag(config crud.ModuleConfig) string {
	if config.Description != "" {
		return config.Description
	}
	if len(config.Permissions) > 0 {
		return config.Permissions[0].Label
	}
	if config.Name != "" {
		return config.Name
	}
	return strings.Trim(config.Group, "/")
}

// routeSummary 根据 handler 名称生成简洁摘要。
func routeSummary(config crud.ModuleConfig, route crud.Route) string {
	module := moduleTag(config)
	switch route.Handler {
	case "List":
		return "获取" + module + "列表"
	case "Get":
		return "获取" + module + "详情"
	case "Create":
		return "创建" + module
	case "Update":
		return "更新" + module
	case "Delete":
		return "删除" + module
	default:
		return route.Handler
	}
}

// routeDescription 返回路由权限和公开状态说明。
func routeDescription(route crud.Route) string {
	if route.Public {
		return "公开接口"
	}
	if route.Permission != "" {
		return "需要权限：" + route.Permission
	}
	return "需要登录"
}

// jsonFieldName 读取 JSON 字段名。
func jsonFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "-" {
		return ""
	}
	name := tagName(jsonTag)
	if name != "" {
		return name
	}
	if tagName(field.Tag.Get("form")) != "" {
		return tagName(field.Tag.Get("form"))
	}
	return field.Name
}

// tagName 解析 struct tag 中逗号前的字段名。
func tagName(tag string) string {
	if tag == "-" {
		return ""
	}
	if idx := strings.Index(tag, ","); idx >= 0 {
		tag = tag[:idx]
	}
	return tag
}

// hasRequiredBinding 判断字段是否带 required 约束。
func hasRequiredBinding(field reflect.StructField) bool {
	return strings.Contains(field.Tag.Get("binding"), "required")
}

// fieldDescription 返回字段说明，优先使用 comment tag。
func fieldDescription(field reflect.StructField) string {
	if comment := field.Tag.Get("comment"); comment != "" {
		return comment
	}
	return field.Name
}

// schemaName 生成稳定的 definition 名称。
func schemaName(sample interface{}) string {
	if sample == nil {
		return ""
	}
	t := indirectType(reflect.TypeOf(sample))
	if t == nil || t.Name() == "" {
		return ""
	}
	if t.PkgPath() == "" {
		return t.Name()
	}
	parts := strings.Split(t.PkgPath(), "/")
	return parts[len(parts)-1] + "." + t.Name()
}

// indirectType 解开指针类型。
func indirectType(t reflect.Type) reflect.Type {
	if t == nil {
		return nil
	}
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

// isTimeLikeType 将时间类字段按字符串处理，避免展开内部实现细节。
func isTimeLikeType(t reflect.Type) bool {
	return (t.PkgPath() == "time" && t.Name() == "Time") || t.Name() == "JSONTime"
}

var typeRegistry = map[string]reflect.Type{}

// typeBySchemaName 读取已注册 definition 对应的 Go 类型。
func typeBySchemaName(name string) reflect.Type {
	return typeRegistry[name]
}
