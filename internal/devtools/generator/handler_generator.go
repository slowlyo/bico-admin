package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// HandlerGenerator handler生成器
type HandlerGenerator struct {
	validator *Validator
}

// NewHandlerGenerator 创建handler生成器
func NewHandlerGenerator() *HandlerGenerator {
	return &HandlerGenerator{
		validator: NewValidator(),
	}
}

// Generate 生成handler文件
func (g *HandlerGenerator) Generate(req *GenerateRequest) (*GenerateResponse, error) {
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
		Message:        fmt.Sprintf("Handler '%s' 生成成功", req.ModelName),
		HistoryUpdated: true,
	}, nil
}

// prepareTemplateData 准备模板数据
func (g *HandlerGenerator) prepareTemplateData(req *GenerateRequest) (*HandlerTemplateData, error) {
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

		// 转换字段名为PascalCase（Go结构体字段命名规范）
		fields[i].Name = ToPascalCase(fields[i].Name)
	}

	// 确定导入包
	imports := g.determineImports(fields)

	// 包名固定为handler
	packageName := "handler"

	templateData := &HandlerTemplateData{
		PackageName:    packageName,
		ModelName:      req.ModelName,
		ModelNameLower: ToLowerCamelCase(req.ModelName),
		ModelNameSnake: ToSnakeCase(req.ModelName),
		Fields:         fields,
		Imports:        imports, // 现在为空，但保留以兼容模板
		HasTimeField:   NeedsTimeImport(req.Fields),
		HasValidation:  NeedsValidationImport(fields),
		Timestamp:      time.Now(),

		// Handler特有的字段
		CreateRequestName: fmt.Sprintf("%sCreateRequest", req.ModelName),
		UpdateRequestName: fmt.Sprintf("%sUpdateRequest", req.ModelName),
		ListRequestName:   fmt.Sprintf("%sListRequest", req.ModelName),
		ResponseName:      fmt.Sprintf("%sResponse", req.ModelName),
		HandlerName:       fmt.Sprintf("%sHandler", req.ModelName),
		ServiceName:       fmt.Sprintf("%sService", req.ModelName),
		ServiceInterface:  fmt.Sprintf("%sService", req.ModelName),

		// 路由相关
		BasePath:    fmt.Sprintf("/%s", ToKebabCase(req.ModelName)),
		RoutePrefix: ToKebabCase(req.ModelName),

		// 权限相关
		PermissionPrefix: ToSnakeCase(req.ModelName),
	}

	return templateData, nil
}

// determineImports 确定需要导入的包（现在使用固定导入）
func (g *HandlerGenerator) determineImports(fields []FieldDefinition) []string {
	// 新的 Handler 模板使用固定的导入，不需要动态计算
	// 所有导入都在模板中硬编码
	return []string{}
}

// generateFiles 生成文件
func (g *HandlerGenerator) generateFiles(req *GenerateRequest, templateData *HandlerTemplateData) ([]string, error) {
	var generatedFiles []string

	// 生成handler文件
	handlerFile, err := g.generateHandlerFile(req, templateData)
	if err != nil {
		return nil, err
	}
	generatedFiles = append(generatedFiles, handlerFile)

	// 生成types文件（请求和响应类型）
	typesFile, err := g.generateTypesFile(req, templateData)
	if err != nil {
		return nil, err
	}
	generatedFiles = append(generatedFiles, typesFile)

	return generatedFiles, nil
}

// generateHandlerFile 生成handler文件
func (g *HandlerGenerator) generateHandlerFile(req *GenerateRequest, templateData *HandlerTemplateData) (string, error) {
	// 构建输出路径
	fileName := fmt.Sprintf("%s.go", ToSnakeCase(req.ModelName))
	outputPath := filepath.Join("internal/admin/handler", fileName)

	// 检查文件冲突
	if err := g.validator.CheckFileConflict(req, outputPath); err != nil {
		return "", err
	}

	// 验证输出路径
	if err := g.validator.ValidateOutputPath(outputPath); err != nil {
		return "", err
	}

	// 加载模板
	tmplContent, err := g.loadTemplate("handler.go.tmpl")
	if err != nil {
		return "", fmt.Errorf("加载模板失败: %w", err)
	}

	// 解析模板
	tmpl, err := template.New("handler").Funcs(g.getTemplateFuncs()).Parse(tmplContent)
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

// generateTypesFile 生成types文件
func (g *HandlerGenerator) generateTypesFile(req *GenerateRequest, templateData *HandlerTemplateData) (string, error) {
	// 构建输出路径
	fileName := fmt.Sprintf("%s_types.go", ToSnakeCase(req.ModelName))
	outputPath := filepath.Join("internal/admin/types", fileName)

	// 检查文件冲突
	if err := g.validator.CheckFileConflict(req, outputPath); err != nil {
		return "", err
	}

	// 验证输出路径
	if err := g.validator.ValidateOutputPath(outputPath); err != nil {
		return "", err
	}

	// 加载模板
	tmplContent, err := g.loadTemplate("handler_types.go.tmpl")
	if err != nil {
		return "", fmt.Errorf("加载模板失败: %w", err)
	}

	// 解析模板
	tmpl, err := template.New("handler_types").Funcs(g.getTemplateFuncs()).Parse(tmplContent)
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
func (g *HandlerGenerator) loadTemplate(templateName string) (string, error) {
	templatePath := filepath.Join("internal/devtools/generator/templates", templateName)
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("读取模板文件'%s'失败: %w", templatePath, err)
	}
	return string(content), nil
}

// getTemplateFuncs 获取模板函数
func (g *HandlerGenerator) getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"contains":         strings.Contains,
		"hasPrefix":        strings.HasPrefix,
		"hasSuffix":        strings.HasSuffix,
		"toLower":          strings.ToLower,
		"toUpper":          strings.ToUpper,
		"toSnakeCase":      ToSnakeCase,
		"toCamelCase":      ToCamelCase,
		"toLowerCamelCase": ToLowerCamelCase,
		"toPascalCase":     ToPascalCase,
		"toKebabCase":      ToKebabCase,
		"GetGoType":        GetGoType,
		"hasStatusField": func(fields []FieldDefinition) bool {
			for _, field := range fields {
				if field.Name == "Status" {
					return true
				}
			}
			return false
		},
		"hasPointerStatusField": func(fields []FieldDefinition) bool {
			for _, field := range fields {
				if field.Name == "Status" && strings.HasPrefix(field.Type, "*") {
					return true
				}
			}
			return false
		},
		"getValidationTag": func(field FieldDefinition) string {
			if field.Validate != "" {
				return fmt.Sprintf(` binding:"%s"`, field.Validate)
			}
			return ""
		},
	}
}

// formatGoFile 格式化Go文件
func (g *HandlerGenerator) formatGoFile(filePath string) error {
	// 这里可以调用gofmt或goimports来格式化代码
	// 为了简化，暂时跳过实际的格式化实现
	return nil
}

// HandlerTemplateData handler模板数据
type HandlerTemplateData struct {
	PackageName    string            // 包名
	ModelName      string            // 模型名（如User）
	ModelNameLower string            // 模型名小写（如user）
	ModelNameSnake string            // 模型名蛇形命名（如user_info）
	Fields         []FieldDefinition // 字段列表
	Imports        []string          // 导入包列表
	HasTimeField   bool              // 是否包含时间字段
	HasValidation  bool              // 是否包含验证
	Timestamp      time.Time         // 生成时间戳

	// Handler特有字段
	CreateRequestName string // 创建请求类型名
	UpdateRequestName string // 更新请求类型名
	ListRequestName   string // 列表请求类型名
	ResponseName      string // 响应类型名
	HandlerName       string // Handler名
	ServiceName       string // Service名
	ServiceInterface  string // Service接口名

	// 路由相关
	BasePath    string // 基础路径
	RoutePrefix string // 路由前缀

	// 权限相关
	PermissionPrefix string // 权限前缀
}
