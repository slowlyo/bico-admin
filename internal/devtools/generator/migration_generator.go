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

// Generate 生成 Migration 代码
func (g *MigrationGenerator) Generate(req *GenerateRequest) (*GenerateResponse, error) {
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

	// 生成 Migration 文件
	migrationFile, err := g.generateMigrationFile(templateData)
	if err != nil {
		errors = append(errors, fmt.Sprintf("生成Migration文件失败: %v", err))
	} else {
		generatedFiles = append(generatedFiles, migrationFile)
	}

	// 构建响应
	success := len(errors) == 0
	message := fmt.Sprintf("Migration生成完成，共生成 %d 个文件", len(generatedFiles))
	if !success {
		message = fmt.Sprintf("Migration生成部分完成，共生成 %d 个文件，%d 个错误", len(generatedFiles), len(errors))
	}

	return &GenerateResponse{
		Success:        success,
		GeneratedFiles: generatedFiles,
		Message:        message,
		Errors:         errors,
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

// generateMigrationFile 生成统一的 Migration 文件
func (g *MigrationGenerator) generateMigrationFile(data *MigrationTemplateData) (string, error) {
	// 确定输出目录和文件路径
	outputDir := data.PackagePath + "/initializer"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 统一的文件名
	fileName := "migration_gen.go"
	outputPath := filepath.Join(outputDir, fileName)

	// 检查文件是否存在
	existingContent, err := g.readExistingFile(outputPath)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("读取现有文件失败: %w", err)
	}

	// 生成新的 Migration 内容
	newContent, err := g.generateMigrationContent(data, existingContent)
	if err != nil {
		return "", fmt.Errorf("生成Migration内容失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(outputPath, []byte(newContent), 0644); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	return outputPath, nil
}

// readExistingFile 读取现有文件内容
func (g *MigrationGenerator) readExistingFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// generateMigrationContent 生成 Migration 内容
func (g *MigrationGenerator) generateMigrationContent(data *MigrationTemplateData, existingContent string) (string, error) {
	if existingContent == "" {
		// 如果文件不存在，创建新文件
		return g.generateNewMigrationFile(data)
	}

	// 如果文件存在，合并内容
	return g.mergeMigrationContent(data, existingContent)
}

// generateNewMigrationFile 生成新的 Migration 文件内容
func (g *MigrationGenerator) generateNewMigrationFile(data *MigrationTemplateData) (string, error) {
	// 创建统一的 Migration 文件模板数据
	unifiedData := &UnifiedMigrationData{
		PackageName: data.PackageName,
		PackagePath: data.PackagePath,
		Timestamp:   data.Timestamp,
		Models:      []MigrationTemplateData{*data},
		Imports:     g.generateImports(data.PackagePath),
	}

	// 使用统一模板生成内容
	return g.executeUnifiedTemplate(unifiedData)
}

// generateImports 生成导入语句
func (g *MigrationGenerator) generateImports(packagePath string) []string {
	return []string{
		"gorm.io/gorm",
		"bico-admin/internal/shared/models",
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

// UnifiedMigrationData 统一 Migration 文件的模板数据
type UnifiedMigrationData struct {
	PackageName string                  // 包名
	PackagePath string                  // 包路径
	Timestamp   time.Time               // 生成时间戳
	Models      []MigrationTemplateData // 所有模型的 Migration 数据
	Imports     []string                // 导入语句
}

// executeUnifiedTemplate 执行统一模板
func (g *MigrationGenerator) executeUnifiedTemplate(data *UnifiedMigrationData) (string, error) {
	// 创建统一的 Migration 文件模板
	templateContent := g.getUnifiedTemplate()

	// 解析模板
	tmpl, err := template.New("unified_migration").Parse(templateContent)
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

// mergeMigrationContent 合并 Migration 内容
func (g *MigrationGenerator) mergeMigrationContent(data *MigrationTemplateData, existingContent string) (string, error) {
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

	// 创建统一的 Migration 文件模板数据
	unifiedData := &UnifiedMigrationData{
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
func (g *MigrationGenerator) parseExistingModels(content string) ([]MigrationTemplateData, error) {
	var models []MigrationTemplateData

	// 使用正则表达式提取模型信息
	// 匹配 "type XXXMigrationRegistrar struct{}" 模式
	registrarPattern := regexp.MustCompile(`type\s+(\w+)MigrationRegistrar\s+struct\{\}`)
	matches := registrarPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			modelName := match[1]

			// 为每个找到的模型创建基本的模板数据
			model := MigrationTemplateData{
				ModelName:      modelName,
				ModelNameLower: strings.ToLower(modelName),
				Timestamp:      time.Now(),
			}

			models = append(models, model)
		}
	}

	return models, nil
}

// getUnifiedTemplate 获取统一的 Migration 模板
func (g *MigrationGenerator) getUnifiedTemplate() string {
	return `// Code generated by bico-admin code generator. DO NOT EDIT.
// Generated at: {{.Timestamp.Format "2006-01-02 15:04:05"}}

package initializer

import (
{{range .Imports}}	"{{.}}"
{{end}})

{{range .Models}}
// {{.ModelName}}MigrationRegistrar {{.ModelName}} Migration 注册器
type {{.ModelName}}MigrationRegistrar struct{}

// GetMigrations 实现 MigrationRegistrar 接口
func (r *{{.ModelName}}MigrationRegistrar) GetMigrations() []interface{} {
	return []interface{}{
		&models.{{.ModelName}}{},
	}
}

// Migrate{{.ModelName}}Table 迁移{{.ModelName}}表
func Migrate{{.ModelName}}Table(db *gorm.DB) error {
	return db.AutoMigrate(&models.{{.ModelName}}{})
}
{{end}}

// init 自动注册所有 Migration
func init() {
{{range .Models}}	// 注册{{.ModelName}} Migration
	{{.ModelNameLower}}MigrationRegistrar := &{{.ModelName}}MigrationRegistrar{}
	RegisterMigrationRegistrar({{.ModelNameLower}}MigrationRegistrar)
{{end}}}
`
}
