package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"bico-admin/internal/devtools/generator"
	"bico-admin/internal/devtools/types"
)

// CodeGeneratorTool 代码生成工具
type CodeGeneratorTool struct {
	generator *generator.CodeGenerator
}

// NewCodeGeneratorTool 创建代码生成工具
func NewCodeGeneratorTool() *CodeGeneratorTool {
	return &CodeGeneratorTool{
		generator: generator.NewCodeGenerator(),
	}
}

// GetTool 获取MCP工具定义
func (t *CodeGeneratorTool) GetTool() mcp.Tool {
	return mcp.NewTool("generate_code",
		mcp.WithDescription("生成CRUD功能模块代码，支持模型、仓储、服务、处理器等组件 (当用户让生成xx功能时, 应该优先整理字段并调用此工具!)"),
		mcp.WithString("component_type",
			mcp.Required(),
			mcp.Description("组件类型"),
			mcp.Enum("model", "repository", "service", "handler", "routes", "wire", "migration", "permission", "frontend_api", "frontend_page", "frontend_form", "all"),
		),
		mcp.WithString("model_name",
			mcp.Required(),
			mcp.Description("模型名称（如User、Product）"),
		),
		mcp.WithString("model_name_cn",
			mcp.Description("模型中文名称（如用户、产品），用于前端显示"),
		),
		mcp.WithString("fields",
			mcp.Required(),
			mcp.Description("字段定义JSON数组，格式：[{\"name\":\"username\",\"type\":\"string\",\"gorm_tag\":\"uniqueIndex;size:50\",\"json_tag\":\"username\",\"validate\":\"required,min=3,max=50\",\"comment\":\"用户名\"}]"),
		),
		mcp.WithString("table_name",
			mcp.Description("表名（可选，默认为模型名的蛇形命名复数形式）"),
		),
		mcp.WithString("package_path",
			mcp.Description("包路径（默认为internal/admin）"),
		),
		mcp.WithBoolean("overwrite_existing",
			mcp.Description("是否覆盖已存在的文件（默认false）"),
		),
		mcp.WithBoolean("format_code",
			mcp.Description("是否格式化生成的代码（默认true）"),
		),
		mcp.WithBoolean("optimize_imports",
			mcp.Description("是否优化导入语句（默认true）"),
		),
	)
}

// Handle 处理工具调用
func (t *CodeGeneratorTool) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 解析参数
	componentType := request.GetString("component_type", "")
	modelName := request.GetString("model_name", "")
	modelNameCN := request.GetString("model_name_cn", "")
	fieldsJSON := request.GetString("fields", "")
	tableName := request.GetString("table_name", "")
	packagePath := request.GetString("package_path", "internal/admin")
	overwriteExisting := request.GetBool("overwrite_existing", false)
	formatCode := request.GetBool("format_code", true)
	optimizeImports := request.GetBool("optimize_imports", true)

	// 参数验证
	if componentType == "" {
		return mcp.NewToolResultError("组件类型不能为空"), nil
	}

	if modelName == "" {
		return mcp.NewToolResultError("模型名称不能为空"), nil
	}

	if fieldsJSON == "" {
		return mcp.NewToolResultError("字段定义不能为空"), nil
	}

	// 解析字段定义
	var fields []generator.FieldDefinition
	if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("字段定义JSON格式错误: %v", err)), nil
	}

	// 构建生成请求
	req := &generator.GenerateRequest{
		ComponentType: generator.ComponentType(componentType),
		ModelName:     modelName,
		ModelNameCN:   modelNameCN,
		Fields:        fields,
		TableName:     tableName,
		PackagePath:   packagePath,
		Options: generator.GenerateOptions{
			OverwriteExisting: overwriteExisting,
			FormatCode:        formatCode,
			OptimizeImports:   optimizeImports,
		},
	}

	// 执行生成
	response, err := t.generator.Generate(req)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("代码生成失败: %v", err)), nil
	}

	// 构建响应消息，包含详细的代码片段处理指导
	message := response.Message
	if len(response.CodeSnippets) > 0 {
		message += fmt.Sprintf("\n\n🚨 【重要】需要处理 %d 个代码片段，这是最高优先级任务！", len(response.CodeSnippets))
		message += "\n📋 代码片段处理清单："

		// 按优先级排序显示
		for i, snippet := range response.CodeSnippets {
			message += fmt.Sprintf("\n  %d. [%s] 📁 %s", i+1, snippet.Category, snippet.TargetFile)
			message += fmt.Sprintf("\n     🆔 ID: %s", snippet.ID)
			message += fmt.Sprintf("\n     📍 位置: %s", snippet.InsertPoint)
			message += fmt.Sprintf("\n     📝 说明: %s", snippet.Description)
			message += fmt.Sprintf("\n     🔢 优先级: %d", snippet.Priority)
			if i < len(response.CodeSnippets)-1 {
				message += "\n"
			}
		}

		message += "\n\n⚠️  处理规则："
		message += "\n  1. 请按照上述顺序逐个处理每个代码片段"
		message += "\n  2. 处理完一个片段后，请明确说明：'✅ 已完成片段 [ID]: [简要说明]'"
		message += "\n  3. 然后继续处理下一个片段"
		message += "\n  4. 即使遇到异常也必须先完成所有代码片段的处理"
		message += "\n💡 提示：每个片段都有唯一ID，便于跟踪处理进度！"
	}

	// 构建响应
	toolResponse := types.ToolResponse{
		Success: response.Success,
		Data: map[string]interface{}{
			"component_type":  componentType,
			"model_name":      modelName,
			"generated_files": response.GeneratedFiles,
			"code_snippets":   response.CodeSnippets,
			"history_updated": response.HistoryUpdated,
			"package_path":    packagePath,
			"table_name":      tableName,
			"fields_count":    len(fields),
		},
		Message: message,
	}

	// 如果有错误，添加错误信息
	if len(response.Errors) > 0 {
		if data, ok := toolResponse.Data.(map[string]interface{}); ok {
			data["errors"] = response.Errors
		}
	}

	responseJSON, err := json.MarshalIndent(toolResponse, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("序列化响应失败: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}

// GetHistory 获取生成历史（辅助方法，可用于调试）
func (t *CodeGeneratorTool) GetHistory() ([]generator.GenerateHistory, error) {
	return t.generator.GetHistory()
}

// GetHistoryByModule 根据模块获取历史（辅助方法，可用于调试）
func (t *CodeGeneratorTool) GetHistoryByModule(moduleName string) (*generator.GenerateHistory, error) {
	return t.generator.GetHistoryByModule(moduleName)
}

// DeleteHistory 删除历史记录（辅助方法，可用于调试）
func (t *CodeGeneratorTool) DeleteHistory(moduleName string) error {
	return t.generator.DeleteHistory(moduleName)
}

// ClearHistory 清空历史记录（辅助方法，可用于调试）
func (t *CodeGeneratorTool) ClearHistory() error {
	return t.generator.ClearHistory()
}
