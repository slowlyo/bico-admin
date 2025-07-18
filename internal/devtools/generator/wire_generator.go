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

// WireGenerator Wire Provider 生成器
type WireGenerator struct {
	templateDir string
}

// NewWireGenerator 创建 Wire 生成器
func NewWireGenerator() *WireGenerator {
	return &WireGenerator{
		templateDir: "internal/devtools/generator/templates",
	}
}

// Generate 生成 Wire Provider 代码
func (g *WireGenerator) Generate(req *GenerateRequest) (*GenerateResponse, error) {
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

	// 生成 Wire Provider 文件
	wireFile, err := g.generateWireFile(templateData)
	if err != nil {
		errors = append(errors, fmt.Sprintf("生成Wire文件失败: %v", err))
	} else {
		generatedFiles = append(generatedFiles, wireFile)
	}

	// 构建响应
	success := len(errors) == 0
	message := fmt.Sprintf("Wire Provider生成完成，共生成 %d 个文件", len(generatedFiles))
	if !success {
		message = fmt.Sprintf("Wire Provider生成部分完成，共生成 %d 个文件，%d 个错误", len(generatedFiles), len(errors))
	}

	return &GenerateResponse{
		Success:        success,
		GeneratedFiles: generatedFiles,
		Message:        message,
		Errors:         errors,
	}, nil
}

// prepareTemplateData 准备模板数据
func (g *WireGenerator) prepareTemplateData(req *GenerateRequest) *WireTemplateData {
	modelName := req.ModelName
	packageName := getPackageNameFromPath(req.PackagePath)

	return &WireTemplateData{
		PackageName:    packageName,
		PackagePath:    req.PackagePath,
		ModelName:      modelName,
		ModelNameLower: ToLowerCamelCase(modelName),
		Timestamp:      time.Now(),
	}
}

// generateWireFile 生成统一的 Wire Provider 文件
func (g *WireGenerator) generateWireFile(data *WireTemplateData) (string, error) {
	// 确定输出目录和文件路径
	outputDir := data.PackagePath
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 统一的文件名
	fileName := "provider_gen.go"
	outputPath := filepath.Join(outputDir, fileName)

	// 检查文件是否存在
	existingContent, err := g.readExistingFile(outputPath)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("读取现有文件失败: %w", err)
	}

	// 生成新的 Wire Provider 内容
	newContent, err := g.generateWireContent(data, existingContent)
	if err != nil {
		return "", fmt.Errorf("生成Wire Provider内容失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(outputPath, []byte(newContent), 0644); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	return outputPath, nil
}

// readExistingFile 读取现有文件内容
func (g *WireGenerator) readExistingFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// generateWireContent 生成 Wire Provider 内容
func (g *WireGenerator) generateWireContent(data *WireTemplateData, existingContent string) (string, error) {
	if existingContent == "" {
		// 如果文件不存在，创建新文件
		return g.generateNewWireFile(data)
	}

	// 如果文件存在，合并内容
	return g.mergeWireContent(data, existingContent)
}

// generateNewWireFile 生成新的 Wire Provider 文件内容
func (g *WireGenerator) generateNewWireFile(data *WireTemplateData) (string, error) {
	// 创建统一的 Wire Provider 文件模板数据
	unifiedData := &UnifiedWireData{
		PackageName: data.PackageName,
		PackagePath: data.PackagePath,
		Timestamp:   data.Timestamp,
		Models:      []WireTemplateData{*data},
		Imports:     g.generateImports(data.PackagePath),
	}

	// 使用统一模板生成内容
	return g.executeUnifiedTemplate(unifiedData)
}

// generateImports 生成导入语句
func (g *WireGenerator) generateImports(packagePath string) []string {
	return []string{
		"github.com/google/wire",
		"gorm.io/gorm",
		"bico-admin/" + packagePath + "/repository",
		"bico-admin/" + packagePath + "/service",
		"bico-admin/" + packagePath + "/handler",
	}
}

// WireTemplateData Wire Provider 模板数据
type WireTemplateData struct {
	PackageName    string    // 包名
	PackagePath    string    // 包路径
	ModelName      string    // 模型名（如User）
	ModelNameLower string    // 模型名小写（如user）
	Timestamp      time.Time // 生成时间戳
}

// UnifiedWireData 统一 Wire Provider 文件的模板数据
type UnifiedWireData struct {
	PackageName string             // 包名
	PackagePath string             // 包路径
	Timestamp   time.Time          // 生成时间戳
	Models      []WireTemplateData // 所有模型的 Wire Provider 数据
	Imports     []string           // 导入语句
}

// executeUnifiedTemplate 执行统一模板
func (g *WireGenerator) executeUnifiedTemplate(data *UnifiedWireData) (string, error) {
	// 创建统一的 Wire Provider 文件模板
	templateContent := g.getUnifiedTemplate()

	// 解析模板
	tmpl, err := template.New("unified_wire").Parse(templateContent)
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

// mergeWireContent 合并 Wire Provider 内容
func (g *WireGenerator) mergeWireContent(data *WireTemplateData, existingContent string) (string, error) {
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

	// 创建统一的 Wire Provider 文件模板数据
	unifiedData := &UnifiedWireData{
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
func (g *WireGenerator) parseExistingModels(content string) ([]WireTemplateData, error) {
	var models []WireTemplateData

	// 使用正则表达式提取模型信息
	// 匹配 "func ProvideXXXRepository" 模式
	repositoryPattern := regexp.MustCompile(`func\s+Provide(\w+)Repository\s*\(`)
	matches := repositoryPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			modelName := match[1]

			// 为每个找到的模型创建基本的模板数据
			model := WireTemplateData{
				ModelName:      modelName,
				ModelNameLower: strings.ToLower(modelName),
				Timestamp:      time.Now(),
			}

			models = append(models, model)
		}
	}

	return models, nil
}

// getUnifiedTemplate 获取统一的 Wire Provider 模板
func (g *WireGenerator) getUnifiedTemplate() string {
	return `// Code generated by bico-admin code generator. DO NOT EDIT.
// Generated at: {{.Timestamp.Format "2006-01-02 15:04:05"}}

package {{.PackageName}}

import (
{{range .Imports}}	"{{.}}"
{{end}}	"bico-admin/internal/admin/routes"
)

{{range .Models}}
// {{.ModelName}}ProviderRegistrar {{.ModelName}} Provider 注册器
type {{.ModelName}}ProviderRegistrar struct{}

// GetProviders 实现 ProviderRegistrar 接口
func (r *{{.ModelName}}ProviderRegistrar) GetProviders() []interface{} {
	return []interface{}{
		Provide{{.ModelName}}Repository,
		Provide{{.ModelName}}Service,
		Provide{{.ModelName}}Handler,
	}
}

// Provide{{.ModelName}}Repository 提供{{.ModelName}} Repository
func Provide{{.ModelName}}Repository(db *gorm.DB) repository.{{.ModelName}}Repository {
	return repository.New{{.ModelName}}Repository(db)
}

// Provide{{.ModelName}}Service 提供{{.ModelName}} Service
func Provide{{.ModelName}}Service(repo repository.{{.ModelName}}Repository) service.{{.ModelName}}Service {
	return service.New{{.ModelName}}Service(repo)
}

// Provide{{.ModelName}}Handler 提供{{.ModelName}} Handler
func Provide{{.ModelName}}Handler(svc service.{{.ModelName}}Service) *handler.{{.ModelName}}Handler {
	return handler.New{{.ModelName}}Handler(svc)
}

// {{.ModelName}}ProviderSet {{.ModelName}} Provider 集合
var {{.ModelName}}ProviderSet = wire.NewSet(
	Provide{{.ModelName}}Repository,
	Provide{{.ModelName}}Service,
	Provide{{.ModelName}}Handler,
)
{{end}}

// HandlerExtender 处理器扩展器接口
type HandlerExtender interface {
	ExtendHandlers(base *routes.Handlers) *routes.Handlers
}

// GeneratedHandlerExtender 生成的处理器扩展器
type GeneratedHandlerExtender struct{}

// ExtendHandlers 扩展处理器集合
func (e *GeneratedHandlerExtender) ExtendHandlers(base *routes.Handlers) *routes.Handlers {
	// 创建扩展的处理器集合
	extended := *base // 复制基础处理器

{{range .Models}}	// 添加{{.ModelName}}Handler（如果存在）
	if {{.ModelNameLower}}Handler := GetGenerated{{.ModelName}}Handler(); {{.ModelNameLower}}Handler != nil {
		extended.{{.ModelName}}Handler = {{.ModelNameLower}}Handler
	}
{{end}}

	return &extended
}

{{range .Models}}
// GetGenerated{{.ModelName}}Handler 获取生成的{{.ModelName}}Handler
// 这个函数会在运行时通过反射或其他机制获取实际的handler实例
var GetGenerated{{.ModelName}}Handler = func() *handler.{{.ModelName}}Handler {
	// 这里应该通过某种机制获取实际的handler实例
	// 暂时返回nil，实际实现需要在wire生成后处理
	return nil
}
{{end}}

// init 自动注册所有 Provider 和扩展器
func init() {
{{range .Models}}	// 注册{{.ModelName}} Provider
	{{.ModelNameLower}}ProviderRegistrar := &{{.ModelName}}ProviderRegistrar{}
	RegisterProviderRegistrar({{.ModelNameLower}}ProviderRegistrar)
{{end}}

	// 注册处理器扩展器
	extender := &GeneratedHandlerExtender{}
	RegisterHandlerExtender(extender)
}

// GeneratedProviderSet 生成的 Provider 集合
// 可以在 wire.Build 中使用这个 ProviderSet
var GeneratedProviderSet = wire.NewSet(
{{range .Models}}	{{.ModelName}}ProviderSet,
{{end}})
`
}
