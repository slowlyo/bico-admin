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
		Message:      fmt.Sprintf("Migration代码片段生成完成，共生成 %d 个片段", len(snippets)),
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

	return snippets, nil
}

// generateMigrationRegistrarSnippet 生成Migration模型片段
func (g *MigrationGenerator) generateMigrationRegistrarSnippet(data *MigrationTemplateData, fields []FieldDefinition) (CodeSnippet, error) {
	tmplContent := `		&models.{{.ModelName}}{},`

	tmpl, err := template.New("migration_model").Parse(tmplContent)
	if err != nil {
		return CodeSnippet{}, fmt.Errorf("解析Migration模型模板失败: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return CodeSnippet{}, fmt.Errorf("执行Migration模型模板失败: %w", err)
	}

	return CodeSnippet{
		Content:      buf.String(),
		TargetFile:   data.PackagePath + "/initializer/database.go",
		InsertPoint:  "在 modelList 数组中，注释之前",
		InsertBefore: "// 注意：生成的模型应该直接添加到上面的 modelList 数组中",
		Description:  fmt.Sprintf("在 modelList 中添加 %s 模型", data.ModelName),
	}, nil
}

// mapFieldTypeToSQL 将Go字段类型映射为SQL类型
func (g *MigrationGenerator) mapFieldTypeToSQL(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int32":
		return "int"
	case "int64":
		return "int64"
	case "uint", "uint32":
		return "uint"
	case "uint64":
		return "uint64"
	case "float32":
		return "float32"
	case "float64":
		return "float64"
	case "bool":
		return "bool"
	case "time.Time", "*time.Time":
		return "time.Time"
	case "[]byte":
		return "[]byte"
	default:
		return "string" // 默认为string类型
	}
}
