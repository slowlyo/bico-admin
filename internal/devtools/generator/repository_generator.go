package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// RepositoryGenerator Repository生成器
type RepositoryGenerator struct {
	validator *Validator
}

// NewRepositoryGenerator 创建Repository生成器
func NewRepositoryGenerator() *RepositoryGenerator {
	return &RepositoryGenerator{
		validator: NewValidator(),
	}
}

// Generate 生成Repository代码
func (g *RepositoryGenerator) Generate(req *GenerateRequest) (*GenerateResponse, error) {
	// 准备模板数据
	templateData, err := g.prepareTemplateData(req)
	if err != nil {
		return &GenerateResponse{
			Success: false,
			Message: "准备模板数据失败",
			Errors:  []string{err.Error()},
		}, nil
	}

	// 生成Repository文件
	outputPath, err := g.generateRepositoryFile(req, templateData)
	if err != nil {
		return &GenerateResponse{
			Success: false,
			Message: "生成Repository文件失败",
			Errors:  []string{err.Error()},
		}, nil
	}

	return &GenerateResponse{
		Success:        true,
		GeneratedFiles: []string{outputPath},
		Message:        fmt.Sprintf("Repository生成成功: %s", outputPath),
		HistoryUpdated: false,
	}, nil
}

// prepareTemplateData 准备模板数据
func (g *RepositoryGenerator) prepareTemplateData(req *GenerateRequest) (*TemplateData, error) {
	// 确定表名
	tableName := req.TableName
	if tableName == "" {
		tableName = ToSnakeCase(req.ModelName) + "s"
	}

	// 确定包路径
	packagePath := req.PackagePath
	if packagePath == "" {
		// 根据模型名推断包路径
		if strings.Contains(strings.ToLower(req.ModelName), "admin") {
			packagePath = "bico-admin/internal/admin"
		} else if strings.Contains(strings.ToLower(req.ModelName), "master") {
			packagePath = "bico-admin/internal/master"
		} else {
			packagePath = "bico-admin/internal/shared"
		}
	}

	// 确定包名
	packageName := "repository"

	templateData := &TemplateData{
		PackageName:    packageName,
		PackagePath:    packagePath,
		ModelName:      req.ModelName,
		ModelNameLower: strings.ToLower(req.ModelName),
		ModelNameSnake: ToSnakeCase(req.ModelName),
		TableName:      tableName,
		Fields:         req.Fields,
		Imports:        []string{}, // Repository不需要额外的导入
		HasTimeField:   false,      // Repository不关心字段类型
		HasValidation:  false,      // Repository不需要验证
		Timestamp:      time.Now(),
	}

	return templateData, nil
}

// generateRepositoryFile 生成Repository文件
func (g *RepositoryGenerator) generateRepositoryFile(req *GenerateRequest, templateData *TemplateData) (string, error) {
	// 确定输出目录
	var outputDir string
	if req.PackagePath != "" {
		// 从包路径推断输出目录
		if strings.Contains(req.PackagePath, "admin") {
			outputDir = "internal/admin/repository"
		} else if strings.Contains(req.PackagePath, "master") {
			outputDir = "internal/master/repository"
		} else {
			outputDir = "internal/shared/repository"
		}
	} else {
		// 默认输出到shared
		outputDir = "internal/shared/repository"
	}

	// 构建输出路径
	fileName := fmt.Sprintf("%s.go", ToSnakeCase(req.ModelName))
	outputPath := filepath.Join(outputDir, fileName)

	// 检查文件冲突
	if err := g.validator.CheckFileConflict(req, outputPath); err != nil {
		return "", err
	}

	// 验证输出路径
	if err := g.validator.ValidateOutputPath(outputPath); err != nil {
		return "", err
	}

	// 确保输出目录存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 加载模板
	tmplContent, err := g.loadTemplate("repository.go.tmpl")
	if err != nil {
		return "", fmt.Errorf("加载模板失败: %w", err)
	}

	// 解析模板
	tmpl, err := template.New("repository").Funcs(g.getTemplateFuncs()).Parse(tmplContent)
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
func (g *RepositoryGenerator) loadTemplate(templateName string) (string, error) {
	templatePath := filepath.Join("internal/devtools/generator/templates", templateName)
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("读取模板文件'%s'失败: %w", templatePath, err)
	}
	return string(content), nil
}

// getTemplateFuncs 获取模板函数
func (g *RepositoryGenerator) getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"ToLower":     strings.ToLower,
		"ToUpper":     strings.ToUpper,
		"ToTitle":     strings.Title,
		"ToSnakeCase": ToSnakeCase,
		"ToCamelCase": ToCamelCase,
		"Contains":    strings.Contains,
		"HasPrefix":   strings.HasPrefix,
		"HasSuffix":   strings.HasSuffix,
		"Replace":     strings.Replace,
		"Split":       strings.Split,
		"Join":        strings.Join,
	}
}

// formatGoFile 格式化Go文件
func (g *RepositoryGenerator) formatGoFile(filePath string) error {
	// 这里可以调用 gofmt 或 goimports 来格式化代码
	// 为了简化，暂时跳过格式化
	return nil
}
