package server

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"bico-admin/internal/devtools/tools"
)

// MCPDevServer MCP开发工具服务器
type MCPDevServer struct {
	config    *Config
	mcpServer *server.MCPServer
	tools     map[string]ToolHandler
}

// ToolHandler 工具处理器接口
type ToolHandler interface {
	GetTool() mcp.Tool
	Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

// NewMCPDevServer 创建MCP开发工具服务器
func NewMCPDevServer(config *Config) (*MCPDevServer, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	// 创建MCP服务器
	mcpServer := server.NewMCPServer(
		config.Name,
		config.Version,
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	devServer := &MCPDevServer{
		config:    config,
		mcpServer: mcpServer,
		tools:     make(map[string]ToolHandler),
	}

	// 注册工具
	if err := devServer.registerTools(); err != nil {
		return nil, fmt.Errorf("注册工具失败: %w", err)
	}

	return devServer, nil
}

// registerTools 注册所有工具
func (s *MCPDevServer) registerTools() error {
	// 注册配置读取工具
	if s.config.Tools.ConfigReader {
		configTool := tools.NewConfigReaderTool()
		if err := s.registerTool("read_config", configTool); err != nil {
			return fmt.Errorf("注册配置读取工具失败: %w", err)
		}
		log.Printf("已注册工具: read_config")
	}

	// 注册数据库操作工具
	if s.config.Tools.DatabaseTool {
		dbTool := tools.NewDatabaseTool()
		if err := s.registerTool("execute_sql", dbTool); err != nil {
			return fmt.Errorf("注册数据库操作工具失败: %w", err)
		}
		log.Printf("已注册工具: execute_sql")
	}

	// 注册代码生成工具
	if s.config.Tools.CodeGenerator {
		codeTool := tools.NewCodeGeneratorTool()
		if err := s.registerTool("generate_code", codeTool); err != nil {
			return fmt.Errorf("注册代码生成工具失败: %w", err)
		}
		log.Printf("已注册工具: generate_code")
	}

	// 注册目录树工具
	if s.config.Tools.DirectoryTree {
		treeTool := tools.NewDirectoryTreeTool()
		if err := s.registerTool("read_directory_tree", treeTool); err != nil {
			return fmt.Errorf("注册目录树工具失败: %w", err)
		}
		log.Printf("已注册工具: read_directory_tree")
	}

	// 注册表结构查看工具
	if s.config.Tools.TableSchema {
		schemaTool := tools.NewTableSchemaTool()
		if err := s.registerTool("inspect_table_schema", schemaTool); err != nil {
			return fmt.Errorf("注册表结构查看工具失败: %w", err)
		}
		log.Printf("已注册工具: inspect_table_schema")
	}

	// TODO: 后续添加其他工具

	return nil
}

// registerTool 注册单个工具
func (s *MCPDevServer) registerTool(name string, handler ToolHandler) error {
	tool := handler.GetTool()
	s.mcpServer.AddTool(tool, handler.Handle)
	s.tools[name] = handler
	return nil
}

// Start 启动服务器
func (s *MCPDevServer) Start(ctx context.Context) error {
	log.Printf("启动 MCP 开发工具服务器...")
	log.Printf("服务器名称: %s", s.config.Name)
	log.Printf("服务器版本: %s", s.config.Version)
	log.Printf("传输方式: %s", s.config.Transport.Mode)

	// 只支持 stdio 模式
	return s.startStdio(ctx)
}

// startStdio 启动Stdio传输
func (s *MCPDevServer) startStdio(ctx context.Context) error {
	log.Printf("使用 Stdio 传输启动服务器")
	log.Printf("Stdio 服务器已启动，等待客户端连接...")

	// 使用 ServeStdio 启动服务器
	return server.ServeStdio(s.mcpServer)
}

// Stop 停止服务器
func (s *MCPDevServer) Stop() error {
	log.Printf("正在停止 MCP 开发工具服务器...")
	// TODO: 实现优雅关闭逻辑
	return nil
}

// GetTools 获取已注册的工具列表
func (s *MCPDevServer) GetTools() []string {
	tools := make([]string, 0, len(s.tools))
	for name := range s.tools {
		tools = append(tools, name)
	}
	return tools
}
