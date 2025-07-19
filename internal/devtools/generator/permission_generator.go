package generator

import (
	"fmt"
	"strings"
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

// GenerateSnippet 生成Permission代码片段
func (g *PermissionGenerator) GenerateSnippet(req *GenerateRequest) (*GenerateResponse, error) {
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
	snippets, err := g.generatePermissionSnippets(templateData)
	if err != nil {
		return &GenerateResponse{
			Success: false,
			Message: "生成Permission代码片段失败",
			Errors:  []string{err.Error()},
		}, nil
	}

	return &GenerateResponse{
		Success:      true,
		CodeSnippets: snippets,
		Message:      fmt.Sprintf("Permission代码片段生成完成，共生成 %d 个片段", len(snippets)),
	}, nil
}

// prepareTemplateData 准备模板数据
func (g *PermissionGenerator) prepareTemplateData(req *GenerateRequest) *PermissionTemplateData {
	modelName := req.ModelName
	modelNameSnake := toSnakeCase(modelName)
	packageName := getPackageNameFromPath(req.PackagePath)

	// 使用传入的中文名，如果没有则使用英文名作为后备
	modelNameChinese := req.ModelNameCN
	if modelNameChinese == "" {
		modelNameChinese = modelName
	}

	// 生成权限定义 - 作为顶级功能，不放在 system 下
	permissions := []PermissionDef{
		{
			Code:        modelNameSnake,
			Name:        fmt.Sprintf("%s管理", modelNameChinese),
			Parent:      "",
			Type:        "module",
			Level:       1,
			Buttons:     "",
			APIs:        "",
			Description: fmt.Sprintf("%s管理", modelNameChinese),
		},
		{
			Code:        fmt.Sprintf("%s:list", modelNameSnake),
			Name:        fmt.Sprintf("查看%s列表", modelNameChinese),
			Parent:      modelNameSnake,
			Type:        "action",
			Level:       1,
			Buttons:     "search,filter",
			APIs:        fmt.Sprintf("/admin-api/%ss,/admin-api/%ss/:id", modelNameSnake, modelNameSnake),
			Description: fmt.Sprintf("查看%s列表权限", modelNameChinese),
		},
		{
			Code:        fmt.Sprintf("%s:create", modelNameSnake),
			Name:        fmt.Sprintf("创建%s", modelNameChinese),
			Parent:      modelNameSnake,
			Type:        "action",
			Level:       3,
			Buttons:     "create",
			APIs:        fmt.Sprintf("/admin-api/%ss", modelNameSnake),
			Description: fmt.Sprintf("创建%s权限", modelNameChinese),
		},
		{
			Code:        fmt.Sprintf("%s:update", modelNameSnake),
			Name:        fmt.Sprintf("编辑%s", modelNameChinese),
			Parent:      modelNameSnake,
			Type:        "action",
			Level:       3,
			Buttons:     "edit,save",
			APIs:        fmt.Sprintf("/admin-api/%ss/:id,/admin-api/%ss/:id/status", modelNameSnake, modelNameSnake),
			Description: fmt.Sprintf("编辑%s权限", modelNameChinese),
		},
		{
			Code:        fmt.Sprintf("%s:delete", modelNameSnake),
			Name:        fmt.Sprintf("删除%s", modelNameChinese),
			Parent:      modelNameSnake,
			Type:        "action",
			Level:       4,
			Buttons:     "delete",
			APIs:        fmt.Sprintf("/admin-api/%ss/:id", modelNameSnake),
			Description: fmt.Sprintf("删除%s权限", modelNameChinese),
		},
	}

	return &PermissionTemplateData{
		PackageName:    packageName,
		PackagePath:    req.PackagePath,
		ModelName:      modelName,
		ModelNameLower: ToLowerCamelCase(modelName),
		ModelNameSnake: modelNameSnake,
		Permissions:    permissions,
		Timestamp:      time.Now(),
	}
}

// generatePermissionSnippets 生成Permission代码片段

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

// generatePermissionSnippets 生成Permission代码片段
func (g *PermissionGenerator) generatePermissionSnippets(data *PermissionTemplateData) ([]CodeSnippet, error) {
	var snippets []CodeSnippet

	// 1. 生成Permission定义片段
	defSnippet, err := g.generatePermissionDefSnippet(data)
	if err != nil {
		return nil, fmt.Errorf("生成Permission定义片段失败: %w", err)
	}
	snippets = append(snippets, defSnippet)

	return snippets, nil
}

// generatePermissionDefSnippet 生成Permission定义片段
func (g *PermissionGenerator) generatePermissionDefSnippet(data *PermissionTemplateData) (CodeSnippet, error) {
	// 构建权限定义
	var permDefs []string
	for _, perm := range data.Permissions {
		permDefs = append(permDefs, fmt.Sprintf(`	{
		Code:        "%s",
		Name:        "%s",
		Parent:      "%s",
		Type:        "%s",
		Level:       %d,
		Buttons:     "%s",
		APIs:        "%s",
		Description: "%s",
	}`, perm.Code, perm.Name, perm.Parent, perm.Type, perm.Level, perm.Buttons, perm.APIs, perm.Description))
	}

	tmplContent := `// {{.ModelNameLower}}PermissionDefs {{.ModelName}} 权限定义
var {{.ModelNameLower}}PermissionDefs = []PermissionDef{
` + strings.Join(permDefs, ",\n") + `,
}`

	return CodeSnippet{
		ID:           fmt.Sprintf("permission_def_%s", strings.ToLower(data.ModelName)),
		Content:      tmplContent,
		TargetFile:   data.PackagePath + "/definitions/permissions.go",
		InsertPoint:  "在 baseDefs 数组中，注释之前",
		InsertBefore: "// 注意：生成的权限定义应该直接添加到上面的 baseDefs 数组中",
		Description:  fmt.Sprintf("添加 %s 权限定义", data.ModelName),
		Priority:     1,
		Category:     "permission_def",
	}, nil
}
