package generator

import (
	"fmt"
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

// GenerateSnippet 生成Wire Provider代码片段
func (g *WireGenerator) GenerateSnippet(req *GenerateRequest) (*GenerateResponse, error) {
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
	snippets, err := g.generateWireSnippets(templateData)
	if err != nil {
		return &GenerateResponse{
			Success: false,
			Message: "生成Wire Provider代码片段失败",
			Errors:  []string{err.Error()},
		}, nil
	}

	return &GenerateResponse{
		Success:      true,
		CodeSnippets: snippets,
		Message:      fmt.Sprintf("Wire Provider代码片段生成完成，共生成 %d 个片段", len(snippets)),
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

// WireTemplateData Wire Provider 模板数据
type WireTemplateData struct {
	PackageName    string    // 包名
	PackagePath    string    // 包路径
	ModelName      string    // 模型名（如User）
	ModelNameLower string    // 模型名小写（如user）
	Timestamp      time.Time // 生成时间戳
}

// generateWireSnippets 生成Wire Provider代码片段
func (g *WireGenerator) generateWireSnippets(data *WireTemplateData) ([]CodeSnippet, error) {
	var snippets []CodeSnippet

	// 1. 生成ProviderSet更新片段
	providerSetUpdateSnippet, err := g.generateProviderSetUpdateSnippet(data)
	if err != nil {
		return nil, fmt.Errorf("生成ProviderSet更新片段失败: %w", err)
	}
	snippets = append(snippets, providerSetUpdateSnippet)

	// 2. 生成ProvideHandlers函数更新片段
	handlersUpdateSnippet, err := g.generateHandlersUpdateSnippet(data)
	if err != nil {
		return nil, fmt.Errorf("生成ProvideHandlers更新片段失败: %w", err)
	}
	snippets = append(snippets, handlersUpdateSnippet)

	return snippets, nil
}

// generateProviderSetUpdateSnippet 生成ProviderSet更新片段
func (g *WireGenerator) generateProviderSetUpdateSnippet(data *WireTemplateData) (CodeSnippet, error) {
	// 根据provider.go的结构，生成需要添加到ProviderSet中的构造函数调用
	tmplContent := `
	// {{.ModelName}} 相关Provider
	repository.New{{.ModelName}}Repository,
	service.New{{.ModelName}}Service,
	handler.New{{.ModelName}}Handler,`

	tmpl, err := template.New("provider_set_update").Parse(tmplContent)
	if err != nil {
		return CodeSnippet{}, fmt.Errorf("解析ProviderSet更新模板失败: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return CodeSnippet{}, fmt.Errorf("执行ProviderSet更新模板失败: %w", err)
	}

	return CodeSnippet{
		ID:          fmt.Sprintf("wire_provider_set_%s", strings.ToLower(data.ModelName)),
		Content:     buf.String(),
		TargetFile:  data.PackagePath + "/provider.go",
		InsertPoint: "在 ProviderSet 中的 Handler层 部分",
		InsertAfter: "handler\\.NewCommonHandler,",
		Description: fmt.Sprintf("在 ProviderSet 中添加 %s 相关的Repository、Service和Handler", data.ModelName),
		Priority:    1,
		Category:    "wire_provider",
	}, nil
}

// generateHandlersUpdateSnippet 生成ProvideHandlers函数更新片段
func (g *WireGenerator) generateHandlersUpdateSnippet(data *WireTemplateData) (CodeSnippet, error) {
	// 生成需要添加到ProvideHandlers函数的参数和字段
	parameterSnippet := fmt.Sprintf("\t\t%sHandler *handler.%sHandler,", ToLowerCamelCase(data.ModelName), data.ModelName)
	fieldSnippet := fmt.Sprintf("\t\t\t%sHandler: %sHandler,", data.ModelName, ToLowerCamelCase(data.ModelName))

	content := fmt.Sprintf(`
需要手动更新 ProvideHandlers 函数：

1. 在函数参数中添加：
%s

2. 在返回的 routes.Handlers 结构体中添加：
%s

3. 确保 routes.Handlers 结构体中也定义了对应的字段

4. wire 命令将于程序启动时自动运行, 你无需处理!你无需处理!你无需处理!`, parameterSnippet, fieldSnippet)

	return CodeSnippet{
		ID:          fmt.Sprintf("wire_handlers_%s", strings.ToLower(data.ModelName)),
		Content:     content,
		TargetFile:  data.PackagePath + "/provider.go",
		InsertPoint: "ProvideHandlers 函数需要手动更新",
		Description: fmt.Sprintf("更新 ProvideHandlers 函数以包含 %s Handler", data.ModelName),
		Priority:    2,
		Category:    "wire_handlers",
	}, nil
}
