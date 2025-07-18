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

// PermissionGenerator Permission 生成器
type PermissionGenerator struct {
	templateDir string
}

// NewPermissionGenerator 创建 Permission 生成器
func NewPermissionGenerator() *PermissionGenerator {
	return &PermissionGenerator{
		templateDir: "internal/devtools/generator/templates",
	}
}

// Generate 生成 Permission 代码
func (g *PermissionGenerator) Generate(req *GenerateRequest) (*GenerateResponse, error) {
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

	// 生成 Permission 文件
	permissionFile, err := g.generatePermissionFile(templateData)
	if err != nil {
		errors = append(errors, fmt.Sprintf("生成Permission文件失败: %v", err))
	} else {
		generatedFiles = append(generatedFiles, permissionFile)
	}

	// 构建响应
	success := len(errors) == 0
	message := fmt.Sprintf("Permission生成完成，共生成 %d 个文件", len(generatedFiles))
	if !success {
		message = fmt.Sprintf("Permission生成部分完成，共生成 %d 个文件，%d 个错误", len(generatedFiles), len(errors))
	}

	return &GenerateResponse{
		Success:        success,
		GeneratedFiles: generatedFiles,
		Message:        message,
		Errors:         errors,
	}, nil
}

// prepareTemplateData 准备模板数据
func (g *PermissionGenerator) prepareTemplateData(req *GenerateRequest) *PermissionTemplateData {
	modelName := req.ModelName
	modelNameLower := strings.ToLower(modelName)
	modelNameSnake := toSnakeCase(modelName)
	packageName := getPackageNameFromPath(req.PackagePath)

	// 生成权限定义
	permissions := []PermissionDef{
		{
			Code:        fmt.Sprintf("system.%s", modelNameSnake),
			Name:        modelName,
			Parent:      "system",
			Type:        "module",
			Level:       1,
			Buttons:     "",
			APIs:        "",
			Description: fmt.Sprintf("%s管理", modelName),
		},
		{
			Code:        fmt.Sprintf("system.%s:list", modelNameSnake),
			Name:        fmt.Sprintf("查看%s列表", modelName),
			Parent:      fmt.Sprintf("system.%s", modelNameSnake),
			Type:        "action",
			Level:       1,
			Buttons:     "search,filter",
			APIs:        fmt.Sprintf("/admin-api/%ss,/admin-api/%ss/:id", modelNameSnake, modelNameSnake),
			Description: fmt.Sprintf("查看%s列表权限", modelName),
		},
		{
			Code:        fmt.Sprintf("system.%s:create", modelNameSnake),
			Name:        fmt.Sprintf("创建%s", modelName),
			Parent:      fmt.Sprintf("system.%s", modelNameSnake),
			Type:        "action",
			Level:       3,
			Buttons:     "create",
			APIs:        fmt.Sprintf("/admin-api/%ss", modelNameSnake),
			Description: fmt.Sprintf("创建%s权限", modelName),
		},
		{
			Code:        fmt.Sprintf("system.%s:update", modelNameSnake),
			Name:        fmt.Sprintf("编辑%s", modelName),
			Parent:      fmt.Sprintf("system.%s", modelNameSnake),
			Type:        "action",
			Level:       3,
			Buttons:     "edit,save",
			APIs:        fmt.Sprintf("/admin-api/%ss/:id,/admin-api/%ss/:id/status", modelNameSnake, modelNameSnake),
			Description: fmt.Sprintf("编辑%s权限", modelName),
		},
		{
			Code:        fmt.Sprintf("system.%s:delete", modelNameSnake),
			Name:        fmt.Sprintf("删除%s", modelName),
			Parent:      fmt.Sprintf("system.%s", modelNameSnake),
			Type:        "action",
			Level:       4,
			Buttons:     "delete",
			APIs:        fmt.Sprintf("/admin-api/%ss/:id", modelNameSnake),
			Description: fmt.Sprintf("删除%s权限", modelName),
		},
	}

	return &PermissionTemplateData{
		PackageName:    packageName,
		PackagePath:    req.PackagePath,
		ModelName:      modelName,
		ModelNameLower: modelNameLower,
		ModelNameSnake: modelNameSnake,
		Permissions:    permissions,
		Timestamp:      time.Now(),
	}
}

// generatePermissionFile 生成统一的 Permission 文件
func (g *PermissionGenerator) generatePermissionFile(data *PermissionTemplateData) (string, error) {
	// 确定输出目录和文件路径
	outputDir := data.PackagePath + "/definitions"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 统一的文件名
	fileName := "permission_gen.go"
	outputPath := filepath.Join(outputDir, fileName)

	// 检查文件是否存在
	existingContent, err := g.readExistingFile(outputPath)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("读取现有文件失败: %w", err)
	}

	// 生成新的 Permission 内容
	newContent, err := g.generatePermissionContent(data, existingContent)
	if err != nil {
		return "", fmt.Errorf("生成Permission内容失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(outputPath, []byte(newContent), 0644); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	return outputPath, nil
}

// readExistingFile 读取现有文件内容
func (g *PermissionGenerator) readExistingFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// generatePermissionContent 生成 Permission 内容
func (g *PermissionGenerator) generatePermissionContent(data *PermissionTemplateData, existingContent string) (string, error) {
	if existingContent == "" {
		// 如果文件不存在，创建新文件
		return g.generateNewPermissionFile(data)
	}

	// 如果文件存在，合并内容
	return g.mergePermissionContent(data, existingContent)
}

// generateNewPermissionFile 生成新的 Permission 文件内容
func (g *PermissionGenerator) generateNewPermissionFile(data *PermissionTemplateData) (string, error) {
	// 创建统一的 Permission 文件模板数据
	unifiedData := &UnifiedPermissionData{
		PackageName: data.PackageName,
		PackagePath: data.PackagePath,
		Timestamp:   data.Timestamp,
		Models:      []PermissionTemplateData{*data},
		Imports:     g.generateImports(),
	}

	// 使用统一模板生成内容
	return g.executeUnifiedTemplate(unifiedData)
}

// generateImports 生成导入语句
func (g *PermissionGenerator) generateImports() []string {
	return []string{
		"strings",
	}
}

// PermissionTemplateData Permission 模板数据
type PermissionTemplateData struct {
	PackageName    string          // 包名
	PackagePath    string          // 包路径
	ModelName      string          // 模型名（如User）
	ModelNameLower string          // 模型名小写（如user）
	ModelNameSnake string          // 模型名蛇形命名（如user_info）
	Permissions    []PermissionDef // 权限定义列表
	Timestamp      time.Time       // 生成时间戳
}

// PermissionDef 权限定义
type PermissionDef struct {
	Code        string // 权限代码
	Name        string // 权限名称
	Parent      string // 父权限
	Type        string // 权限类型
	Level       int    // 权限级别
	Buttons     string // 按钮列表
	APIs        string // API列表
	Description string // 描述
}

// UnifiedPermissionData 统一 Permission 文件的模板数据
type UnifiedPermissionData struct {
	PackageName string                   // 包名
	PackagePath string                   // 包路径
	Timestamp   time.Time                // 生成时间戳
	Models      []PermissionTemplateData // 所有模型的 Permission 数据
	Imports     []string                 // 导入语句
}

// executeUnifiedTemplate 执行统一模板
func (g *PermissionGenerator) executeUnifiedTemplate(data *UnifiedPermissionData) (string, error) {
	// 创建统一的 Permission 文件模板
	templateContent := g.getUnifiedTemplate()

	// 解析模板
	tmpl, err := template.New("unified_permission").Parse(templateContent)
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

// mergePermissionContent 合并 Permission 内容
func (g *PermissionGenerator) mergePermissionContent(data *PermissionTemplateData, existingContent string) (string, error) {
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

	// 创建统一的 Permission 文件模板数据
	unifiedData := &UnifiedPermissionData{
		PackageName: data.PackageName,
		PackagePath: data.PackagePath,
		Timestamp:   time.Now(),
		Models:      existingModels,
		Imports:     g.generateImports(),
	}

	// 使用统一模板生成内容
	return g.executeUnifiedTemplate(unifiedData)
}

// parseExistingModels 解析现有文件中的模型
func (g *PermissionGenerator) parseExistingModels(content string) ([]PermissionTemplateData, error) {
	var models []PermissionTemplateData

	// 使用正则表达式提取模型信息
	// 匹配 "type XXXPermissionRegistrar struct{}" 模式
	registrarPattern := regexp.MustCompile(`type\s+(\w+)PermissionRegistrar\s+struct\{\}`)
	matches := registrarPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			modelName := match[1]

			// 为每个找到的模型创建基本的模板数据
			model := PermissionTemplateData{
				ModelName:      modelName,
				ModelNameLower: strings.ToLower(modelName),
				ModelNameSnake: toSnakeCase(modelName),
				Timestamp:      time.Now(),
			}

			models = append(models, model)
		}
	}

	return models, nil
}

// getUnifiedTemplate 获取统一的 Permission 模板
func (g *PermissionGenerator) getUnifiedTemplate() string {
	return `// Code generated by bico-admin code generator. DO NOT EDIT.
// Generated at: {{.Timestamp.Format "2006-01-02 15:04:05"}}

package definitions

import (
{{range .Imports}}	"{{.}}"
{{end}})

{{range .Models}}
// {{.ModelName}}PermissionRegistrar {{.ModelName}} Permission 注册器
type {{.ModelName}}PermissionRegistrar struct{}

// GetPermissions 实现 PermissionRegistrar 接口
func (r *{{.ModelName}}PermissionRegistrar) GetPermissions() []PermissionDef {
	return {{.ModelNameLower}}PermissionDefs
}

// {{.ModelNameLower}}PermissionDefs {{.ModelName}}权限定义
var {{.ModelNameLower}}PermissionDefs = []PermissionDef{
{{range .Permissions}}	{"{{.Code}}", "{{.Name}}", "{{.Parent}}", "{{.Type}}", {{.Level}}, "{{.Buttons}}", "{{.APIs}}"},
{{end}}}

// {{.ModelName}}权限常量
const (
{{$modelName := .ModelName}}{{range .Permissions}}{{if eq .Type "action"}}	// {{.Description}}
	{{$modelName}}_{{.Type}}_{{.Level}} = "{{.Code}}"
{{end}}{{end}})
{{end}}

// init 自动注册所有 Permission
func init() {
{{range .Models}}	// 注册{{.ModelName}} Permission
	{{.ModelNameLower}}PermissionRegistrar := &{{.ModelName}}PermissionRegistrar{}
	RegisterPermissionRegistrar({{.ModelNameLower}}PermissionRegistrar)
{{end}}}
`
}
