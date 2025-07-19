package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// FrontendAPIGenerator 前端API生成器
type FrontendAPIGenerator struct {
	templateDir string
}

// NewFrontendAPIGenerator 创建前端API生成器
func NewFrontendAPIGenerator() *FrontendAPIGenerator {
	return &FrontendAPIGenerator{
		templateDir: "internal/devtools/generator/templates",
	}
}

// Generate 生成前端API文件
func (g *FrontendAPIGenerator) Generate(req *GenerateRequest) (*GenerateResponse, error) {
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

	// 生成API文件
	filePath, content, err := g.generateAPIFile(templateData)
	if err != nil {
		return &GenerateResponse{
			Success: false,
			Message: "生成前端API文件失败",
			Errors:  []string{err.Error()},
		}, nil
	}

	// 写入文件
	if err := g.writeFile(filePath, content, req.Options.OverwriteExisting); err != nil {
		return &GenerateResponse{
			Success: false,
			Message: "写入前端API文件失败",
			Errors:  []string{err.Error()},
		}, nil
	}

	// 格式化代码（如果需要）
	if req.Options.FormatCode {
		if err := formatTypeScriptFile(filePath); err != nil {
			// 格式化失败不影响生成结果，只记录警告
			fmt.Printf("警告: 格式化TypeScript文件失败: %v\n", err)
		}
	}

	return &GenerateResponse{
		Success:        true,
		GeneratedFiles: []string{filePath},
		Message:        fmt.Sprintf("前端API文件生成成功: %s", filePath),
	}, nil
}

// FrontendAPITemplateData 前端API模板数据
type FrontendAPITemplateData struct {
	ServiceName      string      // 服务类名 (如AdminUserService)
	ServiceNameLower string      // 小写服务名 (如adminUser)
	ModelName        string      // 模型名 (如AdminUser)
	ModelNameLower   string      // 模型名小写 (如adminUser)
	ModelNameSnake   string      // 模型名蛇形命名 (如admin_user)
	APIBasePath      string      // API基础路径 (如/admin-api/admin-users)
	Methods          []APIMethod // API方法列表
	TypeNamespace    string      // 类型命名空间 (如AdminUserTypes)
	TypeDefinitions  []TypeDef   // 类型定义列表
	Imports          []string    // 导入语句
	Timestamp        time.Time   // 生成时间戳
}

// TypeDef 类型定义
type TypeDef struct {
	Name        string      // 类型名称
	Comment     string      // 注释
	Fields      []TypeField // 字段列表
	IsInterface bool        // 是否为接口
}

// TypeField 类型字段
type TypeField struct {
	Name     string // 字段名
	Type     string // 字段类型
	Comment  string // 注释
	Optional bool   // 是否可选
}

// APIMethod API方法定义
type APIMethod struct {
	Name         string // 方法名 (如getAdminUserList)
	Comment      string // 注释
	HTTPMethod   string // HTTP方法 (GET/POST/PUT/DELETE/PATCH)
	URL          string // 请求URL
	ParamsType   string // 参数类型 (可选)
	ResponseType string // 响应类型
	HasParams    bool   // 是否有参数
	HasPathParam bool   // 是否有路径参数
	RequestData  string // 请求数据参数名 (如data)
}

// prepareTemplateData 准备模板数据
func (g *FrontendAPIGenerator) prepareTemplateData(req *GenerateRequest) *FrontendAPITemplateData {
	modelName := req.ModelName
	modelNameLower := ToLowerCamelCase(modelName)
	modelNameSnake := ToSnakeCase(modelName)

	// 生成服务类名
	serviceName := modelName + "Service"

	// 生成API基础路径
	apiBasePath := "/admin-api/" + strings.ReplaceAll(modelNameSnake, "_", "-")

	// 生成类型命名空间（现在使用本地命名空间）
	typeNamespace := modelName + "Types"

	// 生成类型定义
	typeDefinitions := g.generateTypeDefinitions(modelName, req.Fields)

	// 生成API方法
	methods := g.generateAPIMethods(modelName, modelNameLower, apiBasePath, typeNamespace)

	// 生成导入语句
	imports := []string{
		"import request from '@/utils/http'",
	}

	return &FrontendAPITemplateData{
		ServiceName:      serviceName,
		ServiceNameLower: modelNameLower,
		ModelName:        modelName,
		ModelNameLower:   modelNameLower,
		ModelNameSnake:   modelNameSnake,
		APIBasePath:      apiBasePath,
		Methods:          methods,
		TypeNamespace:    typeNamespace,
		TypeDefinitions:  typeDefinitions,
		Imports:          imports,
		Timestamp:        time.Now(),
	}
}

// generateTypeDefinitions 生成类型定义
func (g *FrontendAPIGenerator) generateTypeDefinitions(modelName string, fields []FieldDefinition) []TypeDef {
	// 基础字段
	baseFields := []TypeField{
		{Name: "id", Type: "number", Comment: "ID"},
		{Name: "created_at", Type: "string", Comment: "创建时间"},
		{Name: "updated_at", Type: "string", Comment: "更新时间"},
	}

	// 添加自定义字段
	for _, field := range fields {
		tsType := g.convertGoTypeToTypeScript(field.Type)
		baseFields = append(baseFields, TypeField{
			Name:     g.convertFieldNameToTypeScript(field.Name),
			Type:     tsType,
			Comment:  field.Comment,
			Optional: false,
		})
	}

	// 创建请求字段（排除ID和时间戳字段）
	createFields := []TypeField{}
	updateFields := []TypeField{}

	for _, field := range fields {
		tsType := g.convertGoTypeToTypeScript(field.Type)
		fieldName := g.convertFieldNameToTypeScript(field.Name)

		createFields = append(createFields, TypeField{
			Name:     fieldName,
			Type:     tsType,
			Comment:  field.Comment,
			Optional: false,
		})

		// 更新请求中某些字段可能是可选的
		updateFields = append(updateFields, TypeField{
			Name:     fieldName,
			Type:     tsType,
			Comment:  field.Comment,
			Optional: true,
		})
	}

	return []TypeDef{
		// 基础信息类型
		{
			Name:        modelName + "Info",
			Comment:     modelName + "信息",
			IsInterface: true,
			Fields:      baseFields,
		},
		// 列表数据类型
		{
			Name:        modelName + "ListData",
			Comment:     modelName + "列表数据",
			IsInterface: true,
			Fields: []TypeField{
				{Name: "list", Type: modelName + "Info[]", Comment: "列表数据"},
				{Name: "total", Type: "number", Comment: "总数"},
				{Name: "page", Type: "number", Comment: "当前页"},
				{Name: "page_size", Type: "number", Comment: "每页大小"},
			},
		},
		// 列表查询参数类型
		{
			Name:        modelName + "ListParams",
			Comment:     modelName + "列表查询参数",
			IsInterface: true,
			Fields: []TypeField{
				{Name: "page", Type: "number", Comment: "页码", Optional: true},
				{Name: "page_size", Type: "number", Comment: "每页大小", Optional: true},
				{Name: "[key: string]", Type: "unknown", Comment: "其他参数"},
			},
		},
		// 创建请求类型
		{
			Name:        modelName + "CreateRequest",
			Comment:     modelName + "创建请求",
			IsInterface: true,
			Fields:      createFields,
		},
		// 更新请求类型
		{
			Name:        modelName + "UpdateRequest",
			Comment:     modelName + "更新请求",
			IsInterface: true,
			Fields:      updateFields,
		},
	}
}

// convertGoTypeToTypeScript 将Go类型转换为TypeScript类型
func (g *FrontendAPIGenerator) convertGoTypeToTypeScript(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "number"
	case "float32", "float64", "decimal":
		return "number"
	case "bool":
		return "boolean"
	case "time.Time":
		return "string"
	default:
		return "any"
	}
}

// convertFieldNameToTypeScript 将Go字段名转换为TypeScript字段名
func (g *FrontendAPIGenerator) convertFieldNameToTypeScript(fieldName string) string {
	// 将PascalCase转换为camelCase
	return ToLowerCamelCase(fieldName)
}

// generateAPIMethods 生成基础CRUD API方法列表
func (g *FrontendAPIGenerator) generateAPIMethods(modelName, modelNameLower, apiBasePath, typeNamespace string) []APIMethod {
	return []APIMethod{
		// 获取列表
		{
			Name:         "get" + modelName + "List",
			Comment:      "获取" + modelName + "列表",
			HTTPMethod:   "GET",
			URL:          apiBasePath,
			ParamsType:   typeNamespace + "." + modelName + "ListParams",
			ResponseType: "Api.Http.BaseResponse<" + typeNamespace + "." + modelName + "ListData>",
			HasParams:    true,
		},
		// 根据ID获取
		{
			Name:         "get" + modelName + "ById",
			Comment:      "根据ID获取" + modelName,
			HTTPMethod:   "GET",
			URL:          apiBasePath + "/${id}",
			ResponseType: "Api.Http.BaseResponse<" + typeNamespace + "." + modelName + "Info>",
			HasPathParam: true,
		},
		// 创建
		{
			Name:         "create" + modelName,
			Comment:      "创建" + modelName,
			HTTPMethod:   "POST",
			URL:          apiBasePath,
			ParamsType:   typeNamespace + "." + modelName + "CreateRequest",
			ResponseType: "Api.Http.BaseResponse<" + typeNamespace + "." + modelName + "Info>",
			RequestData:  "data",
		},
		// 更新
		{
			Name:         "update" + modelName,
			Comment:      "更新" + modelName,
			HTTPMethod:   "PUT",
			URL:          apiBasePath + "/${id}",
			ParamsType:   typeNamespace + "." + modelName + "UpdateRequest",
			ResponseType: "Api.Http.BaseResponse<" + typeNamespace + "." + modelName + "Info>",
			HasPathParam: true,
			RequestData:  "data",
		},
		// 删除
		{
			Name:         "delete" + modelName,
			Comment:      "删除" + modelName,
			HTTPMethod:   "DELETE",
			URL:          apiBasePath + "/${id}",
			ResponseType: "Api.Http.BaseResponse<null>",
			HasPathParam: true,
		},
	}
}

// generateAPIFile 生成API文件
func (g *FrontendAPIGenerator) generateAPIFile(data *FrontendAPITemplateData) (string, string, error) {
	// 生成文件路径
	fileName := fmt.Sprintf("%sApi.ts", data.ModelNameLower)
	filePath := filepath.Join("web/src/api", fileName)

	// 加载模板
	tmplPath := filepath.Join(g.templateDir, "frontend_api.ts.tmpl")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		// 如果模板文件不存在，使用内置模板
		return g.generateAPIFileWithBuiltinTemplate(data, filePath)
	}

	// 执行模板
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("执行前端API模板失败: %w", err)
	}

	return filePath, buf.String(), nil
}

// generateAPIFileWithBuiltinTemplate 使用内置模板生成API文件（模板文件不存在时的后备方案）
func (g *FrontendAPIGenerator) generateAPIFileWithBuiltinTemplate(data *FrontendAPITemplateData, filePath string) (string, string, error) {
	return "", "", fmt.Errorf("模板文件不存在，请确保 %s 文件存在", filepath.Join(g.templateDir, "frontend_api.ts.tmpl"))
}

// writeFile 写入文件
func (g *FrontendAPIGenerator) writeFile(filePath, content string, overwrite bool) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); err == nil && !overwrite {
		return fmt.Errorf("文件已存在: %s", filePath)
	}

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 写入内容
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("写入文件内容失败: %w", err)
	}

	return nil
}

// formatTypeScriptFile 格式化TypeScript文件
func formatTypeScriptFile(filePath string) error {
	// TODO: 实现TypeScript文件格式化
	// 可以调用prettier或其他格式化工具
	return nil
}
