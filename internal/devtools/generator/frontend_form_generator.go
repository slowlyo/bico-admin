package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// FrontendFormGenerator 前端表单生成器
type FrontendFormGenerator struct {
	templateDir string
}

// NewFrontendFormGenerator 创建前端表单生成器
func NewFrontendFormGenerator() *FrontendFormGenerator {
	return &FrontendFormGenerator{
		templateDir: "internal/devtools/generator/templates",
	}
}

// Generate 生成前端表单文件
func (g *FrontendFormGenerator) Generate(req *GenerateRequest) (*GenerateResponse, error) {
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

	// 生成表单文件
	filePath, content, err := g.generateFormFile(templateData)
	if err != nil {
		return &GenerateResponse{
			Success: false,
			Message: "生成前端表单文件失败",
			Errors:  []string{err.Error()},
		}, nil
	}

	// 写入文件
	if err := g.writeFile(filePath, content, req.Options.OverwriteExisting); err != nil {
		return &GenerateResponse{
			Success: false,
			Message: "写入前端表单文件失败",
			Errors:  []string{err.Error()},
		}, nil
	}

	// 格式化代码（如果需要）
	if req.Options.FormatCode {
		if err := formatVueFile(filePath); err != nil {
			// 格式化失败不影响生成结果，只记录警告
			fmt.Printf("警告: 格式化Vue文件失败: %v\n", err)
		}
	}

	return &GenerateResponse{
		Success:        true,
		GeneratedFiles: []string{filePath},
		Message:        fmt.Sprintf("前端表单文件生成成功: %s", filePath),
	}, nil
}

// FrontendFormTemplateData 前端表单模板数据
type FrontendFormTemplateData struct {
	ModelName        string            // 模型名 (如Product)
	ModelNameLower   string            // 模型名小写 (如product)
	ModelNameKebab   string            // 模型名短横线命名 (如product-category)
	ModelNameChinese string            // 模型中文名 (如产品)
	ServiceName      string            // 服务类名 (如ProductService)
	TypeNamespace    string            // 类型命名空间 (如ProductTypes)
	APIImportPath    string            // API导入路径
	FormFields       []FormField       // 表单字段定义
	ValidationRules  []ValidationRule  // 验证规则
	Fields           []FieldDefinition // 字段定义
	Imports          []string          // 导入语句
	Timestamp        time.Time         // 生成时间戳
}

// FormField 表单字段定义
type FormField struct {
	Label       string   // 字段标签
	Prop        string   // 字段属性名
	Type        string   // 字段类型 (input, select, textarea, switch, etc.)
	Component   string   // 组件名称
	Placeholder string   // 占位符
	Required    bool     // 是否必填
	Options     []Option // 选项（用于select等）
	Rules       []string // 验证规则
	ColSpan     int      // 栅格占位格数
}

// ValidationRule 验证规则
type ValidationRule struct {
	Field string // 字段名
	Rules string // 规则代码
}

// prepareTemplateData 准备模板数据
func (g *FrontendFormGenerator) prepareTemplateData(req *GenerateRequest) *FrontendFormTemplateData {
	modelName := req.ModelName
	modelNameLower := ToLowerCamelCase(modelName)

	// 使用传入的中文名，如果没有则使用英文名
	modelNameChinese := req.ModelNameCN
	if modelNameChinese == "" {
		modelNameChinese = modelName
	}

	// 生成服务类名和类型命名空间
	serviceName := modelName + "Service"
	typeNamespace := modelName + "Types"

	// 生成API导入路径
	apiImportPath := fmt.Sprintf("@/api/%sApi", modelNameLower)

	// 生成表单字段
	formFields := g.generateFormFields(req.Fields)

	// 生成验证规则
	validationRules := g.generateValidationRules(req.Fields)

	// 生成导入语句
	imports := []string{
		"import type { FormInstance, FormRules } from 'element-plus'",
		"import { ElMessage } from 'element-plus'",
		fmt.Sprintf("import { %s, type %s } from '%s'", serviceName, typeNamespace, apiImportPath),
	}

	return &FrontendFormTemplateData{
		ModelName:        modelName,
		ModelNameLower:   modelNameLower,
		ModelNameKebab:   ToKebabCase(modelName),
		ModelNameChinese: modelNameChinese,
		ServiceName:      serviceName,
		TypeNamespace:    typeNamespace,
		APIImportPath:    apiImportPath,
		FormFields:       formFields,
		ValidationRules:  validationRules,
		Fields:           req.Fields,
		Imports:          imports,
		Timestamp:        time.Now(),
	}
}

// generateFieldLabel 生成字段标签
func (g *FrontendFormGenerator) generateFieldLabel(field FieldDefinition) string {
	// 如果注释包含逗号，取第一部分作为标签
	if strings.Contains(field.Comment, "，") {
		return strings.Split(field.Comment, "，")[0]
	}
	if strings.Contains(field.Comment, ",") {
		return strings.Split(field.Comment, ",")[0]
	}

	// 移除常见的后缀
	comment := field.Comment
	comment = strings.TrimSuffix(comment, "字段")
	comment = strings.TrimSuffix(comment, "信息")

	return comment
}

// generateFormFields 生成表单字段定义
func (g *FrontendFormGenerator) generateFormFields(fields []FieldDefinition) []FormField {
	var formFields []FormField

	for _, field := range fields {
		formField := FormField{
			Label:    g.generateFieldLabel(field),
			Prop:     ToLowerCamelCase(field.Name),
			Required: strings.Contains(field.Validate, "required"),
			ColSpan:  12, // 默认占一半宽度
		}

		// 根据字段类型设置表单组件
		fieldType := strings.TrimPrefix(field.Type, "*") // 移除指针标记
		switch fieldType {
		case "string":
			if strings.Contains(strings.ToLower(field.Name), "password") {
				formField.Type = "password"
				formField.Component = "el-input"
				formField.Placeholder = fmt.Sprintf("请输入%s", field.Comment)
			} else if strings.Contains(strings.ToLower(field.Name), "email") {
				formField.Type = "email"
				formField.Component = "el-input"
				formField.Placeholder = fmt.Sprintf("请输入%s", field.Comment)
			} else if strings.Contains(strings.ToLower(field.Name), "description") ||
				strings.Contains(strings.ToLower(field.Name), "remark") ||
				strings.Contains(strings.ToLower(field.Name), "content") {
				formField.Type = "textarea"
				formField.Component = "el-input"
				formField.Placeholder = fmt.Sprintf("请输入%s", field.Comment)
				formField.ColSpan = 24 // 文本域占满一行
			} else {
				formField.Type = "text"
				formField.Component = "el-input"
				formField.Placeholder = fmt.Sprintf("请输入%s", field.Comment)
			}

		case "int":
			if strings.Contains(strings.ToLower(field.Name), "status") {
				formField.Type = "select"
				formField.Component = "el-select"
				formField.Placeholder = fmt.Sprintf("请选择%s", field.Comment)
				formField.Options = []Option{
					{Label: "启用", Value: 1},
					{Label: "禁用", Value: 0},
				}
			} else {
				formField.Type = "number"
				formField.Component = "el-input-number"
				formField.Placeholder = fmt.Sprintf("请输入%s", field.Comment)
			}

		case "uint":
			formField.Type = "number"
			formField.Component = "el-input-number"
			formField.Placeholder = fmt.Sprintf("请输入%s", field.Comment)

		case "decimal", "float32", "float64":
			formField.Type = "number"
			formField.Component = "el-input-number"
			formField.Placeholder = fmt.Sprintf("请输入%s", field.Comment)

		case "bool":
			formField.Type = "switch"
			formField.Component = "el-switch"
			formField.ColSpan = 24

		case "time.Time":
			if strings.Contains(strings.ToLower(field.Name), "date") {
				formField.Type = "date"
				formField.Component = "el-date-picker"
			} else {
				formField.Type = "datetime"
				formField.Component = "el-date-picker"
			}
			formField.Placeholder = fmt.Sprintf("请选择%s", field.Comment)

		default:
			formField.Type = "text"
			formField.Component = "el-input"
			formField.Placeholder = fmt.Sprintf("请输入%s", field.Comment)
		}

		// 生成验证规则
		if formField.Required {
			formField.Rules = append(formField.Rules, "required")
		}

		// 根据validate标签添加其他规则
		if field.Validate != "" {
			if strings.Contains(field.Validate, "email") {
				formField.Rules = append(formField.Rules, "email")
			}
			if strings.Contains(field.Validate, "min=") {
				formField.Rules = append(formField.Rules, "min")
			}
			if strings.Contains(field.Validate, "max=") {
				formField.Rules = append(formField.Rules, "max")
			}
		}

		formFields = append(formFields, formField)
	}

	return formFields
}

// generateValidationRules 生成验证规则
func (g *FrontendFormGenerator) generateValidationRules(fields []FieldDefinition) []ValidationRule {
	var rules []ValidationRule

	for _, field := range fields {
		fieldName := ToLowerCamelCase(field.Name)
		var ruleLines []string

		// 必填验证
		if strings.Contains(field.Validate, "required") {
			ruleLines = append(ruleLines, "{ required: true, message: '请输入"+field.Comment+"', trigger: 'blur' }")
		}

		// 邮箱验证
		if strings.Contains(field.Validate, "email") {
			ruleLines = append(ruleLines, "{ type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }")
		}

		// 长度验证
		if strings.Contains(field.Validate, "min=") || strings.Contains(field.Validate, "max=") {
			// 简单处理，实际可能需要解析具体的min/max值
			if field.Type == "string" {
				ruleLines = append(ruleLines, "{ min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }")
			}
		}

		if len(ruleLines) > 0 {
			rules = append(rules, ValidationRule{
				Field: fieldName,
				Rules: strings.Join(ruleLines, ",\n    "),
			})
		}
	}

	return rules
}

// generateFormFile 生成表单文件
func (g *FrontendFormGenerator) generateFormFile(data *FrontendFormTemplateData) (string, string, error) {
	// 生成文件路径
	dirPath := filepath.Join("web/src/views/system", ToKebabCase(data.ModelNameLower), "modules")
	fileName := fmt.Sprintf("%s-dialog.vue", ToKebabCase(data.ModelNameLower))
	filePath := filepath.Join(dirPath, fileName)

	// 确保目录存在
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 加载模板
	tmplPath := filepath.Join(g.templateDir, "frontend_form.vue.tmpl")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", "", fmt.Errorf("加载前端表单模板失败: %w", err)
	}

	// 执行模板
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("执行前端表单模板失败: %w", err)
	}

	return filePath, buf.String(), nil
}

// writeFile 写入文件
func (g *FrontendFormGenerator) writeFile(filePath, content string, overwrite bool) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); err == nil && !overwrite {
		return fmt.Errorf("文件已存在: %s", filePath)
	}

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}
