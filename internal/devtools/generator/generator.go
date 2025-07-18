package generator

import (
	"fmt"
)

// CodeGenerator 代码生成器
type CodeGenerator struct {
	modelGenerator      *ModelGenerator
	repositoryGenerator *RepositoryGenerator
	serviceGenerator    *ServiceGenerator
	handlerGenerator    *HandlerGenerator
	historyManager      *HistoryManager
	validator           *Validator
}

// NewCodeGenerator 创建代码生成器
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		modelGenerator:      NewModelGenerator(),
		repositoryGenerator: NewRepositoryGenerator(),
		serviceGenerator:    NewServiceGenerator(),
		handlerGenerator:    NewHandlerGenerator(),
		historyManager:      NewHistoryManager(),
		validator:           NewValidator(),
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

// generateAll 生成所有组件
func (g *CodeGenerator) generateAll(req *GenerateRequest) (*GenerateResponse, error) {
	var allGeneratedFiles []string
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

	// TODO: 5. 生成Routes
	// routesReq := *req
	// routesReq.ComponentType = ComponentRoutes
	// routesResponse, err := g.generateRoutes(&routesReq)
	// ...

	// TODO: 6. 生成Wire配置
	// wireReq := *req
	// wireReq.ComponentType = ComponentWire
	// wireResponse, err := g.generateWire(&wireReq)
	// ...

	// TODO: 7. 生成Migration
	// migrationReq := *req
	// migrationReq.ComponentType = ComponentMigration
	// migrationResponse, err := g.generateMigration(&migrationReq)
	// ...

	// TODO: 8. 生成Permission
	// permissionReq := *req
	// permissionReq.ComponentType = ComponentPermission
	// permissionResponse, err := g.generatePermission(&permissionReq)
	// ...

	// TODO: 9. 生成前端API
	// frontendAPIReq := *req
	// frontendAPIReq.ComponentType = ComponentFrontendAPI
	// frontendAPIResponse, err := g.generateFrontendAPI(&frontendAPIReq)
	// ...

	// TODO: 10. 生成前端页面
	// frontendPageReq := *req
	// frontendPageReq.ComponentType = ComponentFrontendPage
	// frontendPageResponse, err := g.generateFrontendPage(&frontendPageReq)
	// ...

	// TODO: 11. 生成前端表单
	// frontendFormReq := *req
	// frontendFormReq.ComponentType = ComponentFrontendForm
	// frontendFormResponse, err := g.generateFrontendForm(&frontendFormReq)
	// ...

	// 构建最终响应
	success := len(allErrors) == 0
	message := fmt.Sprintf("生成完成，共生成 %d 个文件", len(allGeneratedFiles))
	if !success {
		message = fmt.Sprintf("生成部分完成，共生成 %d 个文件，%d 个错误", len(allGeneratedFiles), len(allErrors))
	}

	response := &GenerateResponse{
		Success:        success,
		GeneratedFiles: allGeneratedFiles,
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

// TODO: 后续实现的生成器方法

// generateService 生成Service
// func (g *CodeGenerator) generateService(req *GenerateRequest) (*GenerateResponse, error) {
//     // TODO: 实现Service生成逻辑
//     return nil, fmt.Errorf("Service生成器尚未实现")
// }

// generateHandler 生成Handler
// func (g *CodeGenerator) generateHandler(req *GenerateRequest) (*GenerateResponse, error) {
//     // TODO: 实现Handler生成逻辑
//     return nil, fmt.Errorf("Handler生成器尚未实现")
// }

// generateRoutes 生成Routes
// func (g *CodeGenerator) generateRoutes(req *GenerateRequest) (*GenerateResponse, error) {
//     // TODO: 实现Routes生成逻辑
//     return nil, fmt.Errorf("Routes生成器尚未实现")
// }

// generateWire 生成Wire配置
// func (g *CodeGenerator) generateWire(req *GenerateRequest) (*GenerateResponse, error) {
//     // TODO: 实现Wire生成逻辑
//     return nil, fmt.Errorf("Wire生成器尚未实现")
// }

// generateMigration 生成Migration
// func (g *CodeGenerator) generateMigration(req *GenerateRequest) (*GenerateResponse, error) {
//     // TODO: 实现Migration生成逻辑
//     return nil, fmt.Errorf("Migration生成器尚未实现")
// }

// generatePermission 生成Permission
// func (g *CodeGenerator) generatePermission(req *GenerateRequest) (*GenerateResponse, error) {
//     // TODO: 实现Permission生成逻辑
//     return nil, fmt.Errorf("Permission生成器尚未实现")
// }

// generateFrontendAPI 生成前端API
// func (g *CodeGenerator) generateFrontendAPI(req *GenerateRequest) (*GenerateResponse, error) {
//     // TODO: 实现前端API生成逻辑
//     return nil, fmt.Errorf("前端API生成器尚未实现")
// }

// generateFrontendPage 生成前端页面
// func (g *CodeGenerator) generateFrontendPage(req *GenerateRequest) (*GenerateResponse, error) {
//     // TODO: 实现前端页面生成逻辑
//     return nil, fmt.Errorf("前端页面生成器尚未实现")
// }

// generateFrontendForm 生成前端表单
// func (g *CodeGenerator) generateFrontendForm(req *GenerateRequest) (*GenerateResponse, error) {
//     // TODO: 实现前端表单生成逻辑
//     return nil, fmt.Errorf("前端表单生成器尚未实现")
// }
