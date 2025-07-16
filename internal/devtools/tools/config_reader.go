package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"bico-admin/internal/devtools/types"
	"bico-admin/internal/devtools/utils"
	"bico-admin/pkg/config"
)

// ConfigReaderTool 配置读取工具
type ConfigReaderTool struct {
	parser *utils.ConfigParser
}

// NewConfigReaderTool 创建配置读取工具
func NewConfigReaderTool() *ConfigReaderTool {
	return &ConfigReaderTool{
		parser: utils.NewConfigParser(),
	}
}

// GetTool 获取MCP工具定义
func (t *ConfigReaderTool) GetTool() mcp.Tool {
	return mcp.NewTool("read_config",
		mcp.WithDescription("读取并解析当前应用的配置文件，返回结构化的配置信息"),
		mcp.WithString("config_path",
			mcp.Description("配置文件路径，为空时使用默认路径"),
		),
		mcp.WithString("format",
			mcp.Description("输出格式"),
			mcp.Enum("json", "summary"),
		),
		mcp.WithBoolean("validate",
			mcp.Description("是否验证配置的完整性"),
		),
	)
}

// Handle 处理工具调用
func (t *ConfigReaderTool) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 解析参数
	configPath := request.GetString("config_path", "")
	format := request.GetString("format", "json")
	validate := request.GetBool("validate", false)

	// 加载配置
	cfg, err := config.Load(configPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("加载配置失败: %v", err)), nil
	}

	// 验证配置（如果需要）
	if validate {
		if err := t.parser.ValidateConfig(cfg); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("配置验证失败: %v", err)), nil
		}
	}

	// 解析配置
	configInfo, err := t.parser.ParseConfig(cfg)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("解析配置失败: %v", err)), nil
	}

	// 根据格式返回结果
	switch format {
	case "json":
		return t.handleJSONFormat(configInfo)
	case "summary":
		return t.handleSummaryFormat(configInfo)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("不支持的输出格式: %s", format)), nil
	}
}

// handleJSONFormat 处理JSON格式输出
func (t *ConfigReaderTool) handleJSONFormat(configInfo *types.ConfigInfo) (*mcp.CallToolResult, error) {
	response := types.ToolResponse{
		Success: true,
		Data:    configInfo,
		Message: "配置读取成功",
	}

	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("序列化响应失败: %v", err)), nil
	}
	return mcp.NewToolResultText(string(responseJSON)), nil
}

// handleSummaryFormat 处理摘要格式输出
func (t *ConfigReaderTool) handleSummaryFormat(configInfo *types.ConfigInfo) (*mcp.CallToolResult, error) {
	summary := fmt.Sprintf(`配置摘要信息:

应用信息:
  名称: %s
  版本: %s
  环境: %s
  调试模式: %t

服务器配置:
  监听地址: %s:%d
  读取超时: %v
  写入超时: %v
  空闲超时: %v

数据库配置:
  驱动: %s
  数据库: %s
  最大空闲连接: %d
  最大打开连接: %d
  连接最大生存时间: %v

Redis配置:
  地址: %s:%d
  数据库: %d
  连接池大小: %d
  最小空闲连接: %d

日志配置:
  级别: %s
  格式: %s
  输出: %s
  文件: %s

JWT配置:
  签发者: %s
  过期时间: %v

缓存配置:
  驱动: %s`,
		configInfo.App.Name,
		configInfo.App.Version,
		configInfo.App.Environment,
		configInfo.App.Debug,
		configInfo.Server.Host,
		configInfo.Server.Port,
		configInfo.Server.ReadTimeout,
		configInfo.Server.WriteTimeout,
		configInfo.Server.IdleTimeout,
		configInfo.Database.Driver,
		configInfo.Database.Database,
		configInfo.Database.MaxIdleConns,
		configInfo.Database.MaxOpenConns,
		configInfo.Database.ConnMaxLifetime,
		configInfo.Redis.Host,
		configInfo.Redis.Port,
		configInfo.Redis.Database,
		configInfo.Redis.PoolSize,
		configInfo.Redis.MinIdleConns,
		configInfo.Log.Level,
		configInfo.Log.Format,
		configInfo.Log.Output,
		configInfo.Log.Filename,
		configInfo.JWT.Issuer,
		configInfo.JWT.ExpireTime,
		configInfo.Cache.Driver,
	)

	response := types.ToolResponse{
		Success: true,
		Data:    summary,
		Message: "配置摘要生成成功",
	}

	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	return mcp.NewToolResultText(string(responseJSON)), nil
}
