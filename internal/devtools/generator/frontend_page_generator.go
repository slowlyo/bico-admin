package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// FrontendPageGenerator 前端页面生成器
type FrontendPageGenerator struct {
	templateDir string
}

// NewFrontendPageGenerator 创建前端页面生成器
func NewFrontendPageGenerator() *FrontendPageGenerator {
	return &FrontendPageGenerator{
		templateDir: "internal/devtools/generator/templates",
	}
}

// Generate 生成前端页面文件
func (g *FrontendPageGenerator) Generate(req *GenerateRequest) (*GenerateResponse, error) {
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

	// 生成页面文件
	filePath, content, err := g.generatePageFile(templateData)
	if err != nil {
		return &GenerateResponse{
			Success: false,
			Message: "生成前端页面文件失败",
			Errors:  []string{err.Error()},
		}, nil
	}

	// 写入文件
	if err := g.writeFile(filePath, content, req.Options.OverwriteExisting); err != nil {
		return &GenerateResponse{
			Success: false,
			Message: "写入前端页面文件失败",
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
		Message:        fmt.Sprintf("前端页面文件生成成功: %s", filePath),
	}, nil
}

// FrontendPageTemplateData 前端页面模板数据
type FrontendPageTemplateData struct {
	ModelName        string            // 模型名 (如Product)
	ModelNameLower   string            // 模型名小写 (如product)
	ModelNameSnake   string            // 模型名蛇形命名 (如product)
	ModelNameKebab   string            // 模型名短横线命名 (如product-category)
	ModelNameChinese string            // 模型中文名 (如产品)
	ServiceName      string            // 服务类名 (如ProductService)
	TypeNamespace    string            // 类型命名空间 (如ProductTypes)
	APIImportPath    string            // API导入路径
	TableColumns     []TableColumn     // 表格列定义
	SearchFormItems  []SearchFormItem  // 搜索表单项
	SearchFormFields []SearchFormField // 搜索表单字段
	Fields           []FieldDefinition // 字段定义
	Imports          []string          // 导入语句
	Timestamp        time.Time         // 生成时间戳
}

// TableColumn 表格列定义
type TableColumn struct {
	Label     string // 列标题
	Prop      string // 列属性名
	Width     string // 列宽度
	Sortable  bool   // 是否可排序
	Formatter string // 格式化函数
	Type      string // 列类型
}

// SearchFormItem 搜索表单项
type SearchFormItem struct {
	Label       string   // 标签
	Prop        string   // 属性名
	Type        string   // 组件类型
	Placeholder string   // 占位符
	Options     []Option // 选项（用于select等）
}

// Option 选项
type Option struct {
	Label string      // 显示文本
	Value interface{} // 值
}

// SearchFormField 搜索表单字段
type SearchFormField struct {
	Prop         string // 属性名
	DefaultValue string // 默认值
}

// prepareTemplateData 准备模板数据
func (g *FrontendPageGenerator) prepareTemplateData(req *GenerateRequest) *FrontendPageTemplateData {
	modelName := req.ModelName
	modelNameLower := ToLowerCamelCase(modelName)
	modelNameSnake := ToSnakeCase(modelName)

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

	// 生成表格列
	tableColumns := g.generateTableColumns(req.Fields)

	// 生成搜索表单项
	searchFormItems := g.generateSearchFormItems(req.Fields)

	// 生成导入语句
	imports := []string{
		"import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'",
		"import { useTable } from '@/composables/useTable'",
		"import { useAuth } from '@/composables/useAuth'",
		fmt.Sprintf("import { %s, type %s } from '%s'", serviceName, typeNamespace, apiImportPath),
		fmt.Sprintf("import %sDialog from './modules/%s-dialog.vue'", modelName, ToKebabCase(modelName)),
		"import { ElMessage, ElMessageBox, ElSwitch, ElTag } from 'element-plus'",
		"import type { ColumnOption, SearchFormItem } from '@/types/component'",
	}

	// 生成搜索表单字段
	searchFormFields := g.generateSearchFormFields(req.Fields)

	return &FrontendPageTemplateData{
		ModelName:        modelName,
		ModelNameLower:   modelNameLower,
		ModelNameSnake:   modelNameSnake,
		ModelNameKebab:   ToKebabCase(modelName),
		ModelNameChinese: modelNameChinese,
		ServiceName:      serviceName,
		TypeNamespace:    typeNamespace,
		APIImportPath:    apiImportPath,
		TableColumns:     tableColumns,
		SearchFormItems:  searchFormItems,
		SearchFormFields: searchFormFields,
		Fields:           req.Fields,
		Imports:          imports,
		Timestamp:        time.Now(),
	}
}

// generateTableColumns 生成表格列定义
func (g *FrontendPageGenerator) generateTableColumns(fields []FieldDefinition) []TableColumn {
	columns := []TableColumn{
		{
			Label: "ID",
			Prop:  "id",
			// 不设置宽度，让表格自动调整
			Type: "number",
		},
	}

	// 添加字段列
	for _, field := range fields {
		column := TableColumn{
			Label:    GetDisplayComment(field.Comment), // 使用清理后的注释
			Prop:     ToLowerCamelCase(field.Name),
			Sortable: true,
		}

		// 根据字段类型设置列属性，但不设置固定宽度
		switch field.Type {
		case "string":
			if strings.Contains(strings.ToLower(field.Name), "status") {
				column.Formatter = "statusFormatter"
				// 状态列可以设置较小的宽度
				column.Width = "100"
			}
			// 其他字符串字段不设置宽度
		case "int", "decimal", "float32", "float64":
			// 数字字段不设置宽度
		case "bool":
			column.Formatter = "boolFormatter"
			column.Width = "80" // 布尔值可以设置较小宽度
		case "time.Time":
			column.Formatter = "timeFormatter"
			// 时间字段不设置宽度，使用minWidth
		}

		columns = append(columns, column)
	}

	// 添加时间列
	columns = append(columns,
		TableColumn{
			Label:     "创建时间",
			Prop:      "created_at",
			Formatter: "timeFormatter",
			// 不设置宽度
		},
		TableColumn{
			Label:     "更新时间",
			Prop:      "updated_at",
			Formatter: "timeFormatter",
			// 不设置宽度
		},
	)

	return columns
}

// generateSearchFormItems 生成搜索表单项
func (g *FrontendPageGenerator) generateSearchFormItems(fields []FieldDefinition) []SearchFormItem {
	var items []SearchFormItem

	for _, field := range fields {
		// 只为字符串和状态字段生成搜索项
		if field.Type == "string" || strings.Contains(strings.ToLower(field.Name), "status") {
			displayComment := GetDisplayComment(field.Comment)

			item := SearchFormItem{
				Label:       displayComment, // 使用清理后的注释
				Prop:        ToLowerCamelCase(field.Name),
				Placeholder: fmt.Sprintf("请输入%s", displayComment),
			}

			// 根据字段类型设置搜索组件类型
			if strings.Contains(strings.ToLower(field.Name), "status") {
				item.Type = "select"
				item.Options = []Option{
					{Label: "启用", Value: 1},
					{Label: "禁用", Value: 0},
				}
				item.Placeholder = fmt.Sprintf("请选择%s", displayComment)
			} else {
				item.Type = "input"
			}

			items = append(items, item)
		}
	}

	return items
}

// generateSearchFormFields 生成搜索表单字段
func (g *FrontendPageGenerator) generateSearchFormFields(fields []FieldDefinition) []SearchFormField {
	var formFields []SearchFormField

	// 添加常见的搜索字段
	for _, field := range fields {
		if field.Name == "created_at" || field.Name == "updated_at" {
			continue // 跳过时间字段
		}

		formField := SearchFormField{
			Prop: ToLowerCamelCase(field.Name),
		}

		// 根据字段类型设置默认值
		switch field.Type {
		case "string":
			formField.DefaultValue = "''"
		case "int", "int32", "int64":
			if strings.Contains(strings.ToLower(field.Name), "status") {
				formField.DefaultValue = "undefined"
			} else {
				formField.DefaultValue = "undefined"
			}
		case "bool":
			formField.DefaultValue = "undefined"
		default:
			formField.DefaultValue = "undefined"
		}

		formFields = append(formFields, formField)
	}

	return formFields
}

// generatePageFile 生成页面文件
func (g *FrontendPageGenerator) generatePageFile(data *FrontendPageTemplateData) (string, string, error) {
	// 生成文件路径 - 直接放在views目录下，不放在system子目录
	dirPath := filepath.Join("web/src/views", ToKebabCase(data.ModelNameLower))
	fileName := "index.vue"
	filePath := filepath.Join(dirPath, fileName)

	// 确保目录存在
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 加载模板
	tmplPath := filepath.Join(g.templateDir, "frontend_page.vue.tmpl")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", "", fmt.Errorf("加载前端页面模板失败: %w", err)
	}

	// 执行模板
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("执行前端页面模板失败: %w", err)
	}

	return filePath, buf.String(), nil
}

// writeFile 写入文件
func (g *FrontendPageGenerator) writeFile(filePath, content string, overwrite bool) error {
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

// formatVueFile 格式化Vue文件（占位符实现）
func formatVueFile(filePath string) error {
	// TODO: 实现Vue文件格式化逻辑
	// 可以调用prettier或其他格式化工具
	return nil
}
