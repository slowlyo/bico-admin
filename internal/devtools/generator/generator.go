package generator

import (
	"fmt"
)

// CodeGenerator 代码生成器
type CodeGenerator struct {
	modelGenerator       *ModelGenerator
	repositoryGenerator  *RepositoryGenerator
	serviceGenerator     *ServiceGenerator
	handlerGenerator     *HandlerGenerator
	routeGenerator       *RouteGenerator
	wireGenerator        *WireGenerator
	migrationGenerator   *MigrationGenerator
	permissionGenerator  *PermissionGenerator
	frontendAPIGenerator *FrontendAPIGenerator
	historyManager       *HistoryManager
	validator            *Validator
}

// NewCodeGenerator 创建代码生成器
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		modelGenerator:       NewModelGenerator(),
		repositoryGenerator:  NewRepositoryGenerator(),
		serviceGenerator:     NewServiceGenerator(),
		handlerGenerator:     NewHandlerGenerator(),
		routeGenerator:       NewRouteGenerator(),
		wireGenerator:        NewWireGenerator(),
		migrationGenerator:   NewMigrationGenerator(),
		permissionGenerator:  NewPermissionGenerator(),
		frontendAPIGenerator: NewFrontendAPIGenerator(),
		historyManager:       NewHistoryManager(),
		validator:            NewValidator(),
	}
}

// Generate 生成代码
func (g *CodeGenerator) Generate(req *GenerateRequest) (*GenerateResponse, error) {
	// 验证请求
	if errors := g.validator.ValidateRequest(req); errors.HasErrors() {
		return &GenerateResponse{
			Success: false,
			Message: "请求验证失败",
			Errors:  []string{errors.Error()},
		}, nil
	}

	// 根据组件类型生成代码
	switch req.ComponentType {
	case ComponentModel:
		return g.generateModel(req)
	case ComponentRepository:
		return g.generateRepository(req)
	case ComponentService:
		return g.generateService(req)
	case ComponentHandler:
		return g.generateHandler(req)
	case ComponentRoutes:
		return g.generateRoutes(req)
	case ComponentWire:
		return g.generateWire(req)
	case ComponentMigration:
		return g.generateMigration(req)
	case ComponentPermission:
		return g.generatePermission(req)
	case ComponentFrontendAPI:
		return g.generateFrontendAPI(req)
	case ComponentFrontendPage:
		return g.generateFrontendPage(req)
	case ComponentFrontendForm:
		return g.generateFrontendForm(req)
	case ComponentFrontendRoute:
		return g.generateFrontendRoute(req)
	case ComponentAll:
		return g.generateAll(req)
	default:
		return &GenerateResponse{
			Success: false,
			Message: fmt.Sprintf("暂不支持生成组件类型: %s", req.ComponentType),
			Errors:  []string{"该组件类型尚未实现"},
		}, nil
	}
}

// generateModel 生成模型
func (g *CodeGenerator) generateModel(req *GenerateRequest) (*GenerateResponse, error) {
	// 使用模型生成器生成
	response, err := g.modelGenerator.Generate(req)
	if err != nil {
		return nil, err
	}

	// 如果生成成功，更新历史记录
	if response.Success && len(response.GeneratedFiles) > 0 {
		if err := g.historyManager.AddHistory(req, response.GeneratedFiles); err != nil {
			// 历史记录更新失败不影响生成结果，只记录警告
			fmt.Printf("警告: 更新历史记录失败: %v\n", err)
			response.HistoryUpdated = false
		}
	}

	return response, nil
}

// generateRepository 生成Repository
func (g *CodeGenerator) generateRepository(req *GenerateRequest) (*GenerateResponse, error) {
	// 使用Repository生成器生成
	response, err := g.repositoryGenerator.Generate(req)
	if err != nil {
		return nil, err
	}

	// 如果生成成功，更新历史记录
	if response.Success && len(response.GeneratedFiles) > 0 {
		if err := g.historyManager.AddHistory(req, response.GeneratedFiles); err != nil {
			// 历史记录更新失败不影响生成结果，只记录警告
			fmt.Printf("警告: 更新历史记录失败: %v\n", err)
			response.HistoryUpdated = false
		}
	}

	return response, nil
}

// generateService 生成Service
func (g *CodeGenerator) generateService(req *GenerateRequest) (*GenerateResponse, error) {
	// 使用Service生成器生成
	response, err := g.serviceGenerator.Generate(req)
	if err != nil {
		return nil, err
	}

	// 如果生成成功，更新历史记录
	if response.Success && len(response.GeneratedFiles) > 0 {
		if err := g.historyManager.AddHistory(req, response.GeneratedFiles); err != nil {
			// 历史记录更新失败不影响生成结果，只记录警告
			fmt.Printf("警告: 更新历史记录失败: %v\n", err)
			response.HistoryUpdated = false
		}
	}

	return response, nil
}

// generateHandler 生成Handler
func (g *CodeGenerator) generateHandler(req *GenerateRequest) (*GenerateResponse, error) {
	// 使用Handler生成器生成
	response, err := g.handlerGenerator.Generate(req)
	if err != nil {
		return nil, err
	}

	// 如果生成成功，更新历史记录
	if response.Success && len(response.GeneratedFiles) > 0 {
		if err := g.historyManager.AddHistory(req, response.GeneratedFiles); err != nil {
			// 历史记录更新失败不影响生成结果，只记录警告
			fmt.Printf("警告: 更新历史记录失败: %v\n", err)
			response.HistoryUpdated = false
		}
	}

	return response, nil
}

// generateRoutes 生成Routes
func (g *CodeGenerator) generateRoutes(req *GenerateRequest) (*GenerateResponse, error) {
	// 使用Routes生成器生成代码片段
	response, err := g.routeGenerator.GenerateSnippet(req)
	if err != nil {
		return nil, err
	}

	// Routes组件生成代码片段，不需要更新历史记录
	return response, nil
}

// generateWire 生成Wire Provider
func (g *CodeGenerator) generateWire(req *GenerateRequest) (*GenerateResponse, error) {
	// 使用Wire生成器生成代码片段
	response, err := g.wireGenerator.GenerateSnippet(req)
	if err != nil {
		return nil, err
	}

	// Wire组件生成代码片段，不需要更新历史记录
	return response, nil
}

// generateMigration 生成Migration
func (g *CodeGenerator) generateMigration(req *GenerateRequest) (*GenerateResponse, error) {
	// 使用Migration生成器生成代码片段
	response, err := g.migrationGenerator.GenerateSnippet(req)
	if err != nil {
		return nil, err
	}

	// Migration组件生成代码片段，不需要更新历史记录
	return response, nil
}

// generatePermission 生成Permission
func (g *CodeGenerator) generatePermission(req *GenerateRequest) (*GenerateResponse, error) {
	// 使用Permission生成器生成代码片段
	response, err := g.permissionGenerator.GenerateSnippet(req)
	if err != nil {
		return nil, err
	}

	// Permission组件生成代码片段，不需要更新历史记录
	return response, nil
}

// generateFrontendAPI 生成前端API
func (g *CodeGenerator) generateFrontendAPI(req *GenerateRequest) (*GenerateResponse, error) {
	// 使用前端API生成器生成
	response, err := g.frontendAPIGenerator.Generate(req)
	if err != nil {
		return nil, err
	}

	// 如果生成成功，更新历史记录
	if response.Success && len(response.GeneratedFiles) > 0 {
		if err := g.historyManager.AddHistory(req, response.GeneratedFiles); err != nil {
			// 历史记录更新失败不影响生成结果，只记录警告
			fmt.Printf("警告: 更新历史记录失败: %v\n", err)
			response.HistoryUpdated = false
		}
	}

	return response, nil
}

// generateFrontendPage 生成前端页面
func (g *CodeGenerator) generateFrontendPage(req *GenerateRequest) (*GenerateResponse, error) {
	// 使用前端页面生成器生成
	generator := NewFrontendPageGenerator()
	response, err := generator.Generate(req)
	if err != nil {
		return nil, err
	}

	// 如果生成成功，更新历史记录
	if response.Success && len(response.GeneratedFiles) > 0 {
		if err := g.historyManager.AddHistory(req, response.GeneratedFiles); err != nil {
			// 历史记录更新失败不影响生成结果，只记录警告
			fmt.Printf("警告: 更新历史记录失败: %v\n", err)
			response.HistoryUpdated = false
		}
	}

	return response, nil
}

// generateFrontendForm 生成前端表单
func (g *CodeGenerator) generateFrontendForm(req *GenerateRequest) (*GenerateResponse, error) {
	// 使用前端表单生成器生成
	generator := NewFrontendFormGenerator()
	response, err := generator.Generate(req)
	if err != nil {
		return nil, err
	}

	// 如果生成成功，更新历史记录
	if response.Success && len(response.GeneratedFiles) > 0 {
		if err := g.historyManager.AddHistory(req, response.GeneratedFiles); err != nil {
			// 历史记录更新失败不影响生成结果，只记录警告
			fmt.Printf("警告: 更新历史记录失败: %v\n", err)
			response.HistoryUpdated = false
		}
	}

	return response, nil
}

// generateFrontendRoute 生成前端路由
func (g *CodeGenerator) generateFrontendRoute(req *GenerateRequest) (*GenerateResponse, error) {
	// 使用前端路由生成器生成代码片段
	generator := NewFrontendRouteGenerator()
	response, err := generator.GenerateSnippet(req)
	if err != nil {
		return nil, err
	}

	// 前端路由组件生成代码片段，不需要更新历史记录
	return response, nil
}

// generateAll 生成所有组件
func (g *CodeGenerator) generateAll(req *GenerateRequest) (*GenerateResponse, error) {
	var allGeneratedFiles []string
	var allCodeSnippets []CodeSnippet
	var allErrors []string

	// 1. 生成模型
	modelReq := *req
	modelReq.ComponentType = ComponentModel

	modelResponse, err := g.generateModel(&modelReq)
	if err != nil {
		return nil, err
	}

	if modelResponse.Success {
		allGeneratedFiles = append(allGeneratedFiles, modelResponse.GeneratedFiles...)
	} else {
		allErrors = append(allErrors, modelResponse.Errors...)
	}

	// 2. 生成Repository
	repositoryReq := *req
	repositoryReq.ComponentType = ComponentRepository

	repositoryResponse, err := g.generateRepository(&repositoryReq)
	if err != nil {
		return nil, err
	}

	if repositoryResponse.Success {
		allGeneratedFiles = append(allGeneratedFiles, repositoryResponse.GeneratedFiles...)
	} else {
		allErrors = append(allErrors, repositoryResponse.Errors...)
	}

	// 3. 生成Service
	serviceReq := *req
	serviceReq.ComponentType = ComponentService

	serviceResponse, err := g.generateService(&serviceReq)
	if err != nil {
		return nil, err
	}

	if serviceResponse.Success {
		allGeneratedFiles = append(allGeneratedFiles, serviceResponse.GeneratedFiles...)
	} else {
		allErrors = append(allErrors, serviceResponse.Errors...)
	}

	// 4. 生成Handler
	handlerReq := *req
	handlerReq.ComponentType = ComponentHandler

	handlerResponse, err := g.generateHandler(&handlerReq)
	if err != nil {
		return nil, err
	}

	if handlerResponse.Success {
		allGeneratedFiles = append(allGeneratedFiles, handlerResponse.GeneratedFiles...)
	} else {
		allErrors = append(allErrors, handlerResponse.Errors...)
	}

	// 5. 生成Routes（代码片段模式）
	routesReq := *req
	routesReq.ComponentType = ComponentRoutes

	routesResponse, err := g.generateRoutes(&routesReq)
	if err != nil {
		return nil, err
	}

	if routesResponse.Success {
		allCodeSnippets = append(allCodeSnippets, routesResponse.CodeSnippets...)
	} else {
		allErrors = append(allErrors, routesResponse.Errors...)
	}

	// 6. 生成Wire Provider（代码片段模式）
	wireReq := *req
	wireReq.ComponentType = ComponentWire

	wireResponse, err := g.generateWire(&wireReq)
	if err != nil {
		return nil, err
	}

	if wireResponse.Success {
		allCodeSnippets = append(allCodeSnippets, wireResponse.CodeSnippets...)
	} else {
		allErrors = append(allErrors, wireResponse.Errors...)
	}

	// 7. 生成Migration（代码片段模式）
	migrationReq := *req
	migrationReq.ComponentType = ComponentMigration

	migrationResponse, err := g.generateMigration(&migrationReq)
	if err != nil {
		return nil, err
	}

	if migrationResponse.Success {
		allCodeSnippets = append(allCodeSnippets, migrationResponse.CodeSnippets...)
	} else {
		allErrors = append(allErrors, migrationResponse.Errors...)
	}

	// 8. 生成Permission（代码片段模式）
	permissionReq := *req
	permissionReq.ComponentType = ComponentPermission

	permissionResponse, err := g.generatePermission(&permissionReq)
	if err != nil {
		return nil, err
	}

	if permissionResponse.Success {
		allCodeSnippets = append(allCodeSnippets, permissionResponse.CodeSnippets...)
	} else {
		allErrors = append(allErrors, permissionResponse.Errors...)
	}

	// 9. 生成前端API
	frontendAPIReq := *req
	frontendAPIReq.ComponentType = ComponentFrontendAPI
	frontendAPIResponse, err := g.generateFrontendAPI(&frontendAPIReq)
	if err != nil {
		return nil, err
	}

	if frontendAPIResponse.Success {
		allGeneratedFiles = append(allGeneratedFiles, frontendAPIResponse.GeneratedFiles...)
	} else {
		allErrors = append(allErrors, frontendAPIResponse.Errors...)
	}

	// 10. 生成前端页面
	frontendPageReq := *req
	frontendPageReq.ComponentType = ComponentFrontendPage
	frontendPageResponse, err := g.generateFrontendPage(&frontendPageReq)
	if err != nil {
		return nil, err
	}

	if frontendPageResponse.Success {
		allGeneratedFiles = append(allGeneratedFiles, frontendPageResponse.GeneratedFiles...)
	} else {
		allErrors = append(allErrors, frontendPageResponse.Errors...)
	}

	// 11. 生成前端表单
	frontendFormReq := *req
	frontendFormReq.ComponentType = ComponentFrontendForm
	frontendFormResponse, err := g.generateFrontendForm(&frontendFormReq)
	if err != nil {
		return nil, err
	}

	if frontendFormResponse.Success {
		allGeneratedFiles = append(allGeneratedFiles, frontendFormResponse.GeneratedFiles...)
	} else {
		allErrors = append(allErrors, frontendFormResponse.Errors...)
	}

	// 12. 生成前端路由（代码片段模式）
	frontendRouteReq := *req
	frontendRouteReq.ComponentType = ComponentFrontendRoute
	frontendRouteResponse, err := g.generateFrontendRoute(&frontendRouteReq)
	if err != nil {
		return nil, err
	}

	if frontendRouteResponse.Success {
		allCodeSnippets = append(allCodeSnippets, frontendRouteResponse.CodeSnippets...)
	} else {
		allErrors = append(allErrors, frontendRouteResponse.Errors...)
	}

	// 构建最终响应
	success := len(allErrors) == 0
	message := fmt.Sprintf("生成完成，共生成 %d 个文件和 %d 个代码片段", len(allGeneratedFiles), len(allCodeSnippets))
	if !success {
		message = fmt.Sprintf("生成部分完成，共生成 %d 个文件和 %d 个代码片段，%d 个错误", len(allGeneratedFiles), len(allCodeSnippets), len(allErrors))
	}

	response := &GenerateResponse{
		Success:        success,
		GeneratedFiles: allGeneratedFiles,
		CodeSnippets:   allCodeSnippets,
		Message:        message,
		HistoryUpdated: success,
		Errors:         allErrors,
	}

	// 如果有文件生成成功，更新历史记录
	if len(allGeneratedFiles) > 0 {
		if err := g.historyManager.AddHistory(req, allGeneratedFiles); err != nil {
			fmt.Printf("警告: 更新历史记录失败: %v\n", err)
			response.HistoryUpdated = false
		}
	}

	return response, nil
}

// GetHistory 获取生成历史
func (g *CodeGenerator) GetHistory() ([]GenerateHistory, error) {
	return g.historyManager.GetHistory()
}

// GetHistoryByModule 根据模块获取历史
func (g *CodeGenerator) GetHistoryByModule(moduleName string) (*GenerateHistory, error) {
	return g.historyManager.GetHistoryByModule(moduleName)
}

// DeleteHistory 删除历史记录
func (g *CodeGenerator) DeleteHistory(moduleName string) error {
	return g.historyManager.DeleteHistory(moduleName)
}

// ClearHistory 清空历史记录
func (g *CodeGenerator) ClearHistory() error {
	return g.historyManager.ClearHistory()
}
