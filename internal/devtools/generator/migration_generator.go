package generator

import (
	"fmt"
	"strings"
	"text/template"
	"time"
)

// MigrationGenerator Migration 生成器
type MigrationGenerator struct {
	templateDir string
}

// NewMigrationGenerator 创建 Migration 生成器
func NewMigrationGenerator() *MigrationGenerator {
	return &MigrationGenerator{
		templateDir: "internal/devtools/generator/templates",
	}
}

// GenerateSnippet 生成Migration代码片段
func (g *MigrationGenerator) GenerateSnippet(req *GenerateRequest) (*GenerateResponse, error) {
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
	snippets, err := g.generateMigrationSnippets(templateData, req.Fields)
	if err != nil {
		return &GenerateResponse{
			Success: false,
			Message: "生成Migration代码片段失败",
			Errors:  []string{err.Error()},
		}, nil
	}

	return &GenerateResponse{
		Success:      true,
		CodeSnippets: snippets,
		Message:      fmt.Sprintf("Migration代码片段生成完成，共生成 %d 个片段 (注意严格遵守: 需要先插入代码内容, 最后处理导入)", len(snippets)),
	}, nil
}

// prepareTemplateData 准备模板数据
func (g *MigrationGenerator) prepareTemplateData(req *GenerateRequest) *MigrationTemplateData {
	modelName := req.ModelName
	packageName := getPackageNameFromPath(req.PackagePath)

	return &MigrationTemplateData{
		PackageName:    packageName,
		PackagePath:    req.PackagePath,
		ModelName:      modelName,
		ModelNameLower: ToLowerCamelCase(modelName),
		Timestamp:      time.Now(),
	}
}

// MigrationTemplateData Migration 模板数据
type MigrationTemplateData struct {
	PackageName    string    // 包名
	PackagePath    string    // 包路径
	ModelName      string    // 模型名（如User）
	ModelNameLower string    // 模型名小写（如user）
	Timestamp      time.Time // 生成时间戳
}

// generateMigrationSnippets 生成Migration代码片段
func (g *MigrationGenerator) generateMigrationSnippets(data *MigrationTemplateData, fields []FieldDefinition) ([]CodeSnippet, error) {
	var snippets []CodeSnippet

	// 1. 生成Migration模型片段
	modelSnippet, err := g.generateMigrationRegistrarSnippet(data, fields)
	if err != nil {
		return nil, fmt.Errorf("生成Migration模型片段失败: %w", err)
	}
	snippets = append(snippets, modelSnippet)

	// 2. 生成导入片段
	importSnippet, err := g.generateMigrationImportSnippet(data, fields)
	if err != nil {
		return nil, fmt.Errorf("生成Migration导入片段失败: %w", err)
	}
	snippets = append(snippets, importSnippet)

	return snippets, nil
}

// generateMigrationRegistrarSnippet 生成Migration模型片段
func (g *MigrationGenerator) generateMigrationRegistrarSnippet(data *MigrationTemplateData, fields []FieldDefinition) (CodeSnippet, error) {
	// 使用带别名的模型引用，避免IDE自动移除导入
	tmplContent := `		&sharedModels.{{.ModelName}}{},`

	tmpl, err := template.New("migration_model").Parse(tmplContent)
	if err != nil {
		return CodeSnippet{}, fmt.Errorf("解析Migration模型模板失败: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return CodeSnippet{}, fmt.Errorf("执行Migration模型模板失败: %w", err)
	}

	return CodeSnippet{
		ID:           fmt.Sprintf("migration_model_%s", strings.ToLower(data.ModelName)),
		Content:      buf.String(),
		TargetFile:   data.PackagePath + "/initializer/database.go",
		InsertPoint:  "在 modelList 数组中，注释之前",
		InsertBefore: "// 注意：生成的模型应该直接添加到上面的 modelList 数组中, model 就应该在 shared 模块下, 这是正确的",
		Description:  fmt.Sprintf("在 modelList 中添加 %s 模型 (⚠️ 重要：请先插入此代码片段，然后再添加导入 sharedModels \"bico-admin/internal/shared/models\")", data.ModelName),
		Priority:     1,
		Category:     "migration_model",
	}, nil
}

// generateMigrationImportSnippet 生成Migration导入片段
func (g *MigrationGenerator) generateMigrationImportSnippet(data *MigrationTemplateData, fields []FieldDefinition) (CodeSnippet, error) {
	// 导入内容
	importContent := `	sharedModels "bico-admin/internal/shared/models"`

	return CodeSnippet{
		ID:          fmt.Sprintf("migration_import_%s", strings.ToLower(data.ModelName)),
		Content:     importContent,
		TargetFile:  data.PackagePath + "/initializer/database.go",
		InsertPoint: "在导入部分添加 shared models 导入",
		InsertAfter: `"bico-admin/internal/admin/models"`,
		Description: fmt.Sprintf("为 %s 模型添加 shared models 导入", data.ModelName),
		Priority:    2, // 优先级较低，在模型片段之后处理
		Category:    "migration_import",
	}, nil
}
