package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

// Generate 生成路由代码
func (g *RouteGenerator) Generate(req *GenerateRequest) (*GenerateResponse, error) {
	// 验证请求参数
	if req.ModelName == "" {
		return &GenerateResponse{
			Success: false,
			Message: "模型名称不能为空",
			Errors:  []string{"ModelName is required"},
		}, nil
	}

	if req.PackagePath == "" {
		return &GenerateResponse{
			Success: false,
			Message: "包路径不能为空",
			Errors:  []string{"PackagePath is required"},
		}, nil
	}

	// 准备模板数据
	templateData := g.prepareTemplateData(req)

	// 生成文件
	var generatedFiles []string
	var errors []string

	// 生成路由注册文件
	routeFile, err := g.generateRouteFile(templateData)
	if err != nil {
		errors = append(errors, fmt.Sprintf("生成路由文件失败: %v", err))
	} else {
		generatedFiles = append(generatedFiles, routeFile)
	}

	// 构建响应
	success := len(errors) == 0
	message := fmt.Sprintf("路由生成完成，共生成 %d 个文件", len(generatedFiles))
	if !success {
		message = fmt.Sprintf("路由生成部分完成，共生成 %d 个文件，%d 个错误", len(generatedFiles), len(errors))
	}

	return &GenerateResponse{
		Success:        success,
		GeneratedFiles: generatedFiles,
		Message:        message,
		Errors:         errors,
	}, nil
}

// prepareTemplateData 准备模板数据
func (g *RouteGenerator) prepareTemplateData(req *GenerateRequest) *RouteTemplateData {
	modelName := req.ModelName
	modelNameLower := strings.ToLower(modelName)
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

// generateRouteFile 生成统一的路由文件
func (g *RouteGenerator) generateRouteFile(data *RouteTemplateData) (string, error) {
	// 确定输出目录和文件路径
	outputDir := data.PackagePath + "/routes"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 统一的文件名
	fileName := "routes_gen.go"
	outputPath := filepath.Join(outputDir, fileName)

	// 检查文件是否存在
	existingContent, err := g.readExistingFile(outputPath)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("读取现有文件失败: %w", err)
	}

	// 生成新的路由内容
	newContent, err := g.generateRouteContent(data, existingContent)
	if err != nil {
		return "", fmt.Errorf("生成路由内容失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(outputPath, []byte(newContent), 0644); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	return outputPath, nil
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

// UnifiedRouteData 统一路由文件的模板数据
type UnifiedRouteData struct {
	PackageName string              // 包名
	PackagePath string              // 包路径
	Timestamp   time.Time           // 生成时间戳
	Models      []RouteTemplateData // 所有模型的路由数据
	Imports     []string            // 导入语句
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

// readExistingFile 读取现有文件内容
func (g *RouteGenerator) readExistingFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// generateRouteContent 生成路由内容
func (g *RouteGenerator) generateRouteContent(data *RouteTemplateData, existingContent string) (string, error) {
	if existingContent == "" {
		// 如果文件不存在，创建新文件
		return g.generateNewRouteFile(data)
	}

	// 如果文件存在，合并内容
	return g.mergeRouteContent(data, existingContent)
}

// generateNewRouteFile 生成新的路由文件内容
func (g *RouteGenerator) generateNewRouteFile(data *RouteTemplateData) (string, error) {
	// 创建统一的路由文件模板数据
	unifiedData := &UnifiedRouteData{
		PackageName: data.PackageName,
		PackagePath: data.PackagePath,
		Timestamp:   data.Timestamp,
		Models:      []RouteTemplateData{*data},
		Imports:     g.generateImports(data.PackagePath),
	}

	// 使用统一模板生成内容
	return g.executeUnifiedTemplate(unifiedData)
}

// generateImports 生成导入语句
func (g *RouteGenerator) generateImports(_ string) []string {
	return []string{
		"github.com/gin-gonic/gin",
	}
}

// executeUnifiedTemplate 执行统一模板
func (g *RouteGenerator) executeUnifiedTemplate(data *UnifiedRouteData) (string, error) {
	// 创建统一的路由文件模板
	templateContent := g.getUnifiedTemplate()

	// 解析模板
	tmpl, err := template.New("unified_routes").Funcs(template.FuncMap{
		"lower": strings.ToLower,
	}).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("解析统一模板失败: %w", err)
	}

	// 执行模板
	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("执行统一模板失败: %w", err)
	}

	return result.String(), nil
}

// getUnifiedTemplate 获取统一的路由模板
func (g *RouteGenerator) getUnifiedTemplate() string {
	return `// Code generated by bico-admin code generator. DO NOT EDIT.
// Generated at: {{.Timestamp.Format "2006-01-02 15:04:05"}}

package routes

import (
	"github.com/gin-gonic/gin"
)

{{range .Models}}
// {{.ModelName}}RouteRegistrar {{.ModelName}}路由注册器
type {{.ModelName}}RouteRegistrar struct{}

// RegisterProtectedRoutes 实现 ProtectedRouteRegistrar 接口
func (r *{{.ModelName}}RouteRegistrar) RegisterProtectedRoutes(protectedGroup *gin.RouterGroup, handlers *Handlers) {
	Register{{.ModelName}}Routes(protectedGroup, handlers)
}

// Register{{.ModelName}}Routes 注册{{.ModelName}}路由（需要认证）
func Register{{.ModelName}}Routes(protectedGroup *gin.RouterGroup, handlers *Handlers) {
	// {{.ModelName}}路由组
	{{.ModelNameLower}}Group := protectedGroup.Group("{{.RoutePath}}")
	{
{{$modelNameLower := .ModelNameLower}}{{$handlerName := .HandlerName}}{{range .HandlerMethods}}		{{$modelNameLower}}Group.{{.HTTPMethod}}("{{.Path}}", handlers.{{$handlerName}}.{{.Name}})
{{end}}	}
}
{{end}}

// init 自动注册所有路由
func init() {
{{range .Models}}	// 注册{{.ModelName}}路由
	{{.ModelNameLower}}Registrar := &{{.ModelName}}RouteRegistrar{}
	RegisterProtectedRouteRegistrar({{.ModelNameLower}}Registrar)
{{end}}}
`
}

// mergeRouteContent 合并路由内容
func (g *RouteGenerator) mergeRouteContent(data *RouteTemplateData, existingContent string) (string, error) {
	// 解析现有文件，提取已存在的模型
	existingModels, err := g.parseExistingModels(existingContent)
	if err != nil {
		return "", fmt.Errorf("解析现有模型失败: %w", err)
	}

	// 检查是否已存在相同的模型
	modelExists := false
	for i, model := range existingModels {
		if model.ModelName == data.ModelName {
			// 更新现有模型
			existingModels[i] = *data
			modelExists = true
			break
		}
	}

	// 如果模型不存在，添加新模型
	if !modelExists {
		existingModels = append(existingModels, *data)
	}

	// 创建统一的路由文件模板数据
	unifiedData := &UnifiedRouteData{
		PackageName: data.PackageName,
		PackagePath: data.PackagePath,
		Timestamp:   time.Now(),
		Models:      existingModels,
		Imports:     g.generateImports(data.PackagePath),
	}

	// 使用统一模板生成内容
	return g.executeUnifiedTemplate(unifiedData)
}

// parseExistingModels 解析现有文件中的模型
func (g *RouteGenerator) parseExistingModels(content string) ([]RouteTemplateData, error) {
	var models []RouteTemplateData

	// 使用正则表达式提取模型信息
	// 匹配 "type XXXRouteRegistrar struct{}" 模式
	registrarPattern := regexp.MustCompile(`type\s+(\w+)RouteRegistrar\s+struct\{\}`)
	matches := registrarPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			modelName := match[1]

			// 为每个找到的模型创建基本的模板数据
			// 注意：这里我们只能恢复基本信息，详细的路由配置需要重新生成
			model := RouteTemplateData{
				ModelName:        modelName,
				ModelNameLower:   strings.ToLower(modelName),
				ModelNameSnake:   toSnakeCase(modelName),
				RoutePath:        "/" + strings.ReplaceAll(toSnakeCase(modelName), "_", "-"),
				HandlerName:      modelName + "Handler",
				HandlerFieldName: strings.ToLower(modelName) + "Handler",
				HandlerMethods: []HandlerMethod{
					{Name: "GetList", HTTPMethod: "GET", Path: "", Description: "获取" + modelName + "列表"},
					{Name: "GetByID", HTTPMethod: "GET", Path: "/:id", Description: "根据ID获取" + modelName},
					{Name: "Create", HTTPMethod: "POST", Path: "", Description: "创建" + modelName},
					{Name: "Update", HTTPMethod: "PUT", Path: "/:id", Description: "更新" + modelName},
					{Name: "Delete", HTTPMethod: "DELETE", Path: "/:id", Description: "删除" + modelName},
				},
				Timestamp: time.Now(),
			}

			models = append(models, model)
		}
	}

	return models, nil
}
