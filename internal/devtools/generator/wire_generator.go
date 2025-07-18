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

	// 1. 生成Provider函数片段
	providerFuncSnippet, err := g.generateProviderFuncSnippet(data)
	if err != nil {
		return nil, fmt.Errorf("生成Provider函数片段失败: %w", err)
	}
	snippets = append(snippets, providerFuncSnippet)

	// 2. 生成ProviderSet更新片段
	providerSetUpdateSnippet, err := g.generateProviderSetUpdateSnippet(data)
	if err != nil {
		return nil, fmt.Errorf("生成ProviderSet更新片段失败: %w", err)
	}
	snippets = append(snippets, providerSetUpdateSnippet)

	return snippets, nil
}

// generateProviderFuncSnippet 生成Provider函数片段
func (g *WireGenerator) generateProviderFuncSnippet(data *WireTemplateData) (CodeSnippet, error) {
	tmplContent := `// Provide{{.ModelName}}Repository 提供{{.ModelName}} Repository
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
}`

	tmpl, err := template.New("provider_func").Parse(tmplContent)
	if err != nil {
		return CodeSnippet{}, fmt.Errorf("解析Provider函数模板失败: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return CodeSnippet{}, fmt.Errorf("执行Provider函数模板失败: %w", err)
	}

	return CodeSnippet{
		Content:     buf.String(),
		TargetFile:  data.PackagePath + "/provider.go",
		InsertPoint: "在 ProvidePermissionMiddleware 函数之后",
		InsertAfter: "func ProvidePermissionMiddleware\\(.*\\) gin\\.HandlerFunc \\{[\\s\\S]*?\\}",
		Description: fmt.Sprintf("添加 %s Provider函数", data.ModelName),
	}, nil
}

// generateProviderSetUpdateSnippet 生成GeneratedProviderSet更新片段
func (g *WireGenerator) generateProviderSetUpdateSnippet(data *WireTemplateData) (CodeSnippet, error) {
	tmplContent := `	{{.ModelName}}ProviderSet,`

	tmpl, err := template.New("provider_set_update").Parse(tmplContent)
	if err != nil {
		return CodeSnippet{}, fmt.Errorf("解析ProviderSet更新模板失败: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return CodeSnippet{}, fmt.Errorf("执行ProviderSet更新模板失败: %w", err)
	}

	return CodeSnippet{
		Content:     buf.String(),
		TargetFile:  data.PackagePath + "/provider.go",
		InsertPoint: "在 ProviderSet 中，权限中间件之后",
		InsertAfter: "ProvidePermissionMiddleware,",
		Description: fmt.Sprintf("在 ProviderSet 中添加 %s Provider", data.ModelName),
	}, nil
}
