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

	// TODO: 后续添加其他工具
	// if s.config.Tools.CRUDGenerator {
	//     crudTool := tools.NewCRUDGeneratorTool()
	//     if err := s.registerTool("generate_crud", crudTool); err != nil {
	//         return fmt.Errorf("注册CRUD生成工具失败: %w", err)
	//     }
	//     log.Printf("已注册工具: generate_crud")
	// }

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
	log.Printf("传输方式: HTTP")

	return s.startHTTP(ctx)
}

// startHTTP 启动HTTP传输
func (s *MCPDevServer) startHTTP(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.config.Transport.HTTP.Host, s.config.Transport.HTTP.Port)
	log.Printf("使用 HTTP 传输启动服务器，监听地址: %s", addr)

	// 创建 StreamableHTTPServer
	httpServer := server.NewStreamableHTTPServer(s.mcpServer,
		server.WithEndpointPath("/mcp"),
	)

	// 启动服务器
	go func() {
		<-ctx.Done()
		log.Printf("正在关闭 HTTP 服务器...")
		httpServer.Shutdown(context.Background())
	}()

	log.Printf("HTTP 服务器已启动，访问地址: http://%s/mcp", addr)
	return httpServer.Start(addr)
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
