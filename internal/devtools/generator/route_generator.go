package generator

import (
	"fmt"
	"strings"
	"text/template"
	"time"
)

// RouteGenerator 路由生成器
type RouteGenerator struct {
	templateDir string
}

// NewRouteGenerator 创建路由生成器
func NewRouteGenerator() *RouteGenerator {
	return &RouteGenerator{
		templateDir: "internal/devtools/generator/templates",
	}
}

// GenerateSnippet 生成路由代码片段
func (g *RouteGenerator) GenerateSnippet(req *GenerateRequest) (*GenerateResponse, error) {
	// 验证请求参数
	if req.ModelName == "" {
		return &GenerateResponse{
			Success: false,
			Message: "模型名称不能为空",
			Errors:  []string{"ModelName is required"},
		}, nil
	}

	// 准备模板数据
	templateData := g.prepareTemplateData(req)

	// 生成代码片段
	snippets, err := g.generateRouteSnippets(templateData)
	if err != nil {
		return &GenerateResponse{
			Success: false,
			Message: "生成路由代码片段失败",
			Errors:  []string{err.Error()},
		}, nil
	}

	return &GenerateResponse{
		Success:      true,
		CodeSnippets: snippets,
		Message:      fmt.Sprintf("路由代码片段生成完成，共生成 %d 个片段", len(snippets)),
	}, nil
}

// prepareTemplateData 准备模板数据
func (g *RouteGenerator) prepareTemplateData(req *GenerateRequest) *RouteTemplateData {
	modelName := req.ModelName
	modelNameLower := ToLowerCamelCase(modelName)
	modelNameSnake := toSnakeCase(modelName)
	packageName := getPackageNameFromPath(req.PackagePath)

	// 生成路由路径
	routePath := "/" + strings.ReplaceAll(modelNameSnake, "_", "-")

	// 生成处理器方法
	handlerMethods := []HandlerMethod{
		{Name: "GetList", HTTPMethod: "GET", Path: "", Description: "获取" + modelName + "列表"},
		{Name: "GetByID", HTTPMethod: "GET", Path: "/:id", Description: "根据ID获取" + modelName},
		{Name: "Create", HTTPMethod: "POST", Path: "", Description: "创建" + modelName},
		{Name: "Update", HTTPMethod: "PUT", Path: "/:id", Description: "更新" + modelName},
		{Name: "Delete", HTTPMethod: "DELETE", Path: "/:id", Description: "删除" + modelName},
	}

	return &RouteTemplateData{
		PackageName:      packageName,
		PackagePath:      req.PackagePath,
		ModelName:        modelName,
		ModelNameLower:   modelNameLower,
		ModelNameSnake:   modelNameSnake,
		RoutePath:        routePath,
		HandlerMethods:   handlerMethods,
		HandlerName:      modelName + "Handler",
		HandlerFieldName: modelNameLower + "Handler",
		Timestamp:        time.Now(),
	}
}

// RouteTemplateData 路由模板数据
type RouteTemplateData struct {
	PackageName      string          // 包名
	PackagePath      string          // 包路径
	ModelName        string          // 模型名（如User）
	ModelNameLower   string          // 模型名小写（如user）
	ModelNameSnake   string          // 模型名蛇形命名（如user_info）
	RoutePath        string          // 路由路径（如/user-info）
	HandlerMethods   []HandlerMethod // 处理器方法列表
	HandlerName      string          // 处理器名称（如UserHandler）
	HandlerFieldName string          // 处理器字段名（如userHandler）
	Timestamp        time.Time       // 生成时间戳
}

// HandlerMethod 处理器方法
type HandlerMethod struct {
	Name        string // 方法名
	HTTPMethod  string // HTTP方法
	Path        string // 路径
	Description string // 描述
}

// getPackageNameFromPath 从包路径获取包名
func getPackageNameFromPath(packagePath string) string {
	parts := strings.Split(packagePath, "/")
	if len(parts) == 0 {
		return "main"
	}
	return parts[len(parts)-1]
}

// toSnakeCase 转换为蛇形命名
func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// generateRouteSnippets 生成路由代码片段
func (g *RouteGenerator) generateRouteSnippets(data *RouteTemplateData) ([]CodeSnippet, error) {
	var snippets []CodeSnippet

	// 1. 生成路由注册函数片段
	registerFuncSnippet, err := g.generateRouteRegisterFuncSnippet(data)
	if err != nil {
		return nil, fmt.Errorf("生成路由注册函数片段失败: %w", err)
	}
	snippets = append(snippets, registerFuncSnippet)

	return snippets, nil
}

// generateRouteRegisterFuncSnippet 生成路由注册函数片段
func (g *RouteGenerator) generateRouteRegisterFuncSnippet(data *RouteTemplateData) (CodeSnippet, error) {
	tmplContent := `// Register{{.ModelName}}Routes 注册{{.ModelName}}路由
func Register{{.ModelName}}Routes(router *gin.Engine, handlers *Handlers) {
	// {{.ModelName}}路由组
	{{.ModelNameLower}}Group := router.Group("{{.RoutePath}}")
	{
{{$modelNameLower := .ModelNameLower}}{{$handlerName := .HandlerName}}{{range .HandlerMethods}}		{{$modelNameLower}}Group.{{.HTTPMethod}}("{{.Path}}", handlers.{{$handlerName}}.{{.Name}})
{{end}}	}
}

// Register{{.ModelName}}ProtectedRoutes 注册{{.ModelName}}受保护路由（需要认证）
func Register{{.ModelName}}ProtectedRoutes(protectedGroup *gin.RouterGroup, handlers *Handlers) {
	// {{.ModelName}}路由组
	{{.ModelNameLower}}Group := protectedGroup.Group("{{.RoutePath}}")
	{
{{$modelNameLower := .ModelNameLower}}{{$handlerName := .HandlerName}}{{range .HandlerMethods}}		{{$modelNameLower}}Group.{{.HTTPMethod}}("{{.Path}}", handlers.{{$handlerName}}.{{.Name}})
{{end}}	}
}`

	tmpl, err := template.New("route_register_func").Parse(tmplContent)
	if err != nil {
		return CodeSnippet{}, fmt.Errorf("解析路由注册函数模板失败: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return CodeSnippet{}, fmt.Errorf("执行路由注册函数模板失败: %w", err)
	}

	return CodeSnippet{
		Content:      buf.String(),
		TargetFile:   data.PackagePath + "/routes/routes.go",
		InsertPoint:  "在 RegisterRoutes 函数的末尾，注释之前",
		InsertBefore: "// 注意：生成的路由代码应该直接添加到上面的相应位置",
		Description:  fmt.Sprintf("添加 %s 路由注册函数", data.ModelName),
	}, nil
}
