package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// ModelGenerator 模型生成器
type ModelGenerator struct {
	validator *Validator
}

// NewModelGenerator 创建模型生成器
func NewModelGenerator() *ModelGenerator {
	return &ModelGenerator{
		validator: NewValidator(),
	}
}

// Generate 生成模型文件
func (g *ModelGenerator) Generate(req *GenerateRequest) (*GenerateResponse, error) {
	// 验证请求
	if errors := g.validator.ValidateRequest(req); errors.HasErrors() {
		return &GenerateResponse{
			Success: false,
			Message: "请求验证失败",
			Errors:  []string{errors.Error()},
		}, nil
	}

	// 准备模板数据
	templateData, err := g.prepareTemplateData(req)
	if err != nil {
		return &GenerateResponse{
			Success: false,
			Message: "准备模板数据失败",
			Errors:  []string{err.Error()},
		}, nil
	}

	// 生成文件
	generatedFiles, err := g.generateFiles(req, templateData)
	if err != nil {
		return &GenerateResponse{
			Success: false,
			Message: "生成文件失败",
			Errors:  []string{err.Error()},
		}, nil
	}

	return &GenerateResponse{
		Success:        true,
		GeneratedFiles: generatedFiles,
		Message:        fmt.Sprintf("模型'%s'生成成功", req.ModelName),
		HistoryUpdated: true,
	}, nil
}

// prepareTemplateData 准备模板数据
func (g *ModelGenerator) prepareTemplateData(req *GenerateRequest) (*TemplateData, error) {
	// 处理表名
	tableName := req.TableName
	if tableName == "" {
		// 默认使用模型名的蛇形命名复数形式
		tableName = ToPlural(ToSnakeCase(req.ModelName))
	}

	// 处理字段
	fields := make([]FieldDefinition, len(req.Fields))
	copy(fields, req.Fields)

	// 处理字段类型和标签
	for i := range fields {
		// 转换为Go类型
		fields[i].Type = GetGoType(fields[i].Type)

		// 如果没有JSON标签，使用字段名的蛇形命名
		if fields[i].JsonTag == "" {
			fields[i].JsonTag = ToSnakeCase(fields[i].Name)
		}

		// 清理字段名
		fields[i].Name = SanitizeGoIdentifier(fields[i].Name)
	}

	// 确定导入包（在字段类型转换之后检查）
	hasTimeField := NeedsTimeImport(fields) // 使用处理后的字段检查
	imports := g.determineImports(fields)

	// 包名固定为models（shared/models目录）
	packageName := "models"

	templateData := &TemplateData{
		PackageName:    packageName,
		PackagePath:    "", // Model 不需要 PackagePath
		ModelName:      req.ModelName,
		ModelNameLower: ToLowerCamelCase(req.ModelName),
		ModelNameSnake: ToSnakeCase(req.ModelName),
		TableName:      tableName,
		Fields:         fields,
		Imports:        imports,
		HasTimeField:   hasTimeField,
		HasValidation:  NeedsValidationImport(fields),
		Timestamp:      time.Now(),
	}

	return templateData, nil
}

// determineImports 确定需要导入的包
func (g *ModelGenerator) determineImports(fields []FieldDefinition) []string {
	var imports []string

	// 检查是否需要time包
	if NeedsTimeImport(fields) {
		imports = append(imports, "time")
	}

	// 总是导入shared/types包
	imports = append(imports, "bico-admin/internal/shared/types")

	return imports
}

// generateFiles 生成文件
func (g *ModelGenerator) generateFiles(req *GenerateRequest, templateData *TemplateData) ([]string, error) {
	var generatedFiles []string

	// 生成模型文件
	modelFile, err := g.generateModelFile(req, templateData)
	if err != nil {
		return nil, err
	}
	generatedFiles = append(generatedFiles, modelFile)

	return generatedFiles, nil
}

// generateModelFile 生成模型文件
func (g *ModelGenerator) generateModelFile(req *GenerateRequest, templateData *TemplateData) (string, error) {
	// 构建输出路径 - 统一放到shared/models目录
	fileName := fmt.Sprintf("%s.go", ToSnakeCase(req.ModelName))
	outputPath := filepath.Join("internal/shared/models", fileName)

	// 检查文件冲突
	if err := g.validator.CheckFileConflict(req, outputPath); err != nil {
		return "", err
	}

	// 验证输出路径
	if err := g.validator.ValidateOutputPath(outputPath); err != nil {
		return "", err
	}

	// 加载模板
	tmplContent, err := g.loadTemplate("model.go.tmpl")
	if err != nil {
		return "", fmt.Errorf("加载模板失败: %w", err)
	}

	// 解析模板
	tmpl, err := template.New("model").Funcs(g.getTemplateFuncs()).Parse(tmplContent)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %w", err)
	}

	// 创建输出文件
	file, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 执行模板
	if err := tmpl.Execute(file, templateData); err != nil {
		return "", fmt.Errorf("执行模板失败: %w", err)
	}

	// 格式化代码（如果启用）
	if req.Options.FormatCode {
		if err := g.formatGoFile(outputPath); err != nil {
			// 格式化失败不影响生成，只记录警告
			fmt.Printf("警告: 格式化文件'%s'失败: %v\n", outputPath, err)
		}
	}

	return outputPath, nil
}

// loadTemplate 加载模板文件
func (g *ModelGenerator) loadTemplate(templateName string) (string, error) {
	templatePath := filepath.Join("internal/devtools/generator/templates", templateName)
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("读取模板文件'%s'失败: %w", templatePath, err)
	}
	return string(content), nil
}

// getTemplateFuncs 获取模板函数
func (g *ModelGenerator) getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"contains":     strings.Contains,
		"hasPrefix":    strings.HasPrefix,
		"hasSuffix":    strings.HasSuffix,
		"toLower":      strings.ToLower,
		"toUpper":      strings.ToUpper,
		"toSnakeCase":  ToSnakeCase,
		"toCamelCase":  ToCamelCase,
		"toPascalCase": ToPascalCase,
		"GetGoType":    GetGoType,
		"hasStatusField": func(fields []FieldDefinition) bool {
			for _, field := range fields {
				if field.Name == "Status" {
					return true
				}
			}
			return false
		},
	}
}

// formatGoFile 格式化Go文件
func (g *ModelGenerator) formatGoFile(filePath string) error {
	// 这里可以调用gofmt或goimports来格式化代码
	// 为了简化，暂时跳过实际的格式化实现
	return nil
}
