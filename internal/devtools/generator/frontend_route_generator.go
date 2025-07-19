package generator

import (
	"fmt"
	"time"
)

// FrontendRouteGenerator 前端路由生成器
type FrontendRouteGenerator struct {
	templateDir string
}

// NewFrontendRouteGenerator 创建前端路由生成器
func NewFrontendRouteGenerator() *FrontendRouteGenerator {
	return &FrontendRouteGenerator{
		templateDir: "internal/devtools/generator/templates",
	}
}

// GenerateSnippet 生成前端路由代码片段
func (g *FrontendRouteGenerator) GenerateSnippet(req *GenerateRequest) (*GenerateResponse, error) {
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
	snippets, err := g.generateRouteSnippets(templateData)
	if err != nil {
		return &GenerateResponse{
			Success: false,
			Message: "生成前端路由代码片段失败",
			Errors:  []string{err.Error()},
		}, nil
	}

	return &GenerateResponse{
		Success:      true,
		CodeSnippets: snippets,
		Message:      fmt.Sprintf("前端路由代码片段生成完成，共生成 %d 个片段", len(snippets)),
	}, nil
}

// FrontendRouteTemplateData 前端路由模板数据
type FrontendRouteTemplateData struct {
	ModelName        string // 模型名 (如Product)
	ModelNameLower   string // 模型名小写 (如product)
	ModelNameSnake   string // 模型名蛇形命名 (如product)
	ModelNameKebab   string // 模型名短横线命名 (如product)
	ModelNameChinese string // 模型中文名 (如产品)
	RoutePath        string // 路由路径 (如/system/product)
	ComponentPath    string // 组件路径 (如() => import('@/views/system/product/index.vue'))
	PermissionPrefix string // 权限前缀 (如system.product)
	Icon             string // 图标
	Timestamp        time.Time
}

// prepareTemplateData 准备模板数据
func (g *FrontendRouteGenerator) prepareTemplateData(req *GenerateRequest) *FrontendRouteTemplateData {
	modelName := req.ModelName
	modelNameLower := ToLowerCamelCase(modelName)
	modelNameSnake := ToSnakeCase(modelName)
	modelNameKebab := ToKebabCase(modelName)

	// 使用传入的中文名，如果没有则使用英文名
	modelNameChinese := req.ModelNameCN
	if modelNameChinese == "" {
		modelNameChinese = modelName
	}

	// 生成路由路径
	routePath := fmt.Sprintf("/system/%s", modelNameKebab)

	// 生成组件路径（使用 RoutesAlias）
	componentPath := fmt.Sprintf("RoutesAlias.%s", modelName)

	// 生成权限前缀
	permissionPrefix := fmt.Sprintf("system.%s", modelNameSnake)

	// 生成图标（简单映射）
	icon := g.generateIcon(modelName)

	return &FrontendRouteTemplateData{
		ModelName:        modelName,
		ModelNameLower:   modelNameLower,
		ModelNameSnake:   modelNameSnake,
		ModelNameKebab:   modelNameKebab,
		ModelNameChinese: modelNameChinese,
		RoutePath:        routePath,
		ComponentPath:    componentPath,
		PermissionPrefix: permissionPrefix,
		Icon:             icon,
		Timestamp:        time.Now(),
	}
}

// generateIcon 生成图标占位符
func (g *FrontendRouteGenerator) generateIcon(modelName string) string {
	// 返回占位符，引导用户自行修改
	return "请修改图标"
}

// generateRouteSnippets 生成路由代码片段
func (g *FrontendRouteGenerator) generateRouteSnippets(data *FrontendRouteTemplateData) ([]CodeSnippet, error) {
	var snippets []CodeSnippet

	// 1. 生成路由配置片段
	routeSnippet := g.generateRouteConfigSnippet(data)
	snippets = append(snippets, routeSnippet)

	// 2. 生成路由别名片段（如果需要）
	aliasSnippet := g.generateRouteAliasSnippet(data)
	snippets = append(snippets, aliasSnippet)

	return snippets, nil
}

// generateRouteConfigSnippet 生成路由配置片段
func (g *FrontendRouteGenerator) generateRouteConfigSnippet(data *FrontendRouteTemplateData) CodeSnippet {
	content := fmt.Sprintf(`  // %s管理路由 - 请添加到 System 路由的 children 数组中
  // 注意：请根据实际需求修改 title（菜单名称）和 icon（图标）
  {
    path: '%s',
    name: '%s',
    component: %s,
    meta: {
      title: '【请修改】%s管理',  // 请修改为实际的菜单名称
      icon: '【请修改】%s',      // 请修改为实际的图标代码
      permissions: ['%s:list']
    }
  }`,
		data.ModelNameChinese,
		data.ModelNameKebab,
		data.ModelName,
		data.ComponentPath,
		data.ModelNameChinese,
		data.Icon,
		data.PermissionPrefix,
	)

	return CodeSnippet{
		Content:     content,
		TargetFile:  "web/src/router/routes/asyncRoutes.ts",
		InsertPoint: fmt.Sprintf("在 System 路由的 children 数组中添加 %s 路由配置", data.ModelNameChinese),
		Description: "将此代码片段添加到 web/src/router/routes/asyncRoutes.ts 文件中 System 路由的 children 数组中",
	}
}

// generateRouteAliasSnippet 生成路由别名片段
func (g *FrontendRouteGenerator) generateRouteAliasSnippet(data *FrontendRouteTemplateData) CodeSnippet {
	content := fmt.Sprintf(`  // %s管理页面
  %s = '/system/%s', // %s管理`,
		data.ModelNameChinese,
		data.ModelName,
		data.ModelNameKebab,
		data.ModelNameChinese,
	)

	return CodeSnippet{
		Content:     content,
		TargetFile:  "web/src/router/routesAlias.ts",
		InsertPoint: fmt.Sprintf("在 RoutesAlias 对象中添加 %s 路由别名（可选）", data.ModelNameChinese),
		Description: "将此代码片段添加到 web/src/router/routesAlias.ts 文件中的 RoutesAlias 对象中（可选）",
	}
}
