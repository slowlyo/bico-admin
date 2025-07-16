package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"bico-admin/internal/devtools/server"
)

const (
	defaultLogLevel = "info"
)

func main() {
	// 解析命令行参数
	var (
		host     = flag.String("host", "127.0.0.1", "HTTP服务器监听地址")
		port     = flag.Int("port", 18901, "HTTP服务器监听端口")
		logLevel = flag.String("log-level", defaultLogLevel, "日志级别 (debug, info, warn, error)")
		version  = flag.Bool("version", false, "显示版本信息")
		help     = flag.Bool("help", false, "显示帮助信息")
	)
	flag.Parse()

	// 显示版本信息
	if *version {
		fmt.Println("Bico Admin DevTools MCP Server v1.0.0")
		fmt.Println("基于 Model Context Protocol 的开发工具服务 (HTTP模式)")
		return
	}

	// 显示帮助信息
	if *help {
		printHelp()
		return
	}

	// 创建服务器配置
	config := server.DefaultConfig()
	config.Transport.HTTP.Host = *host
	config.Transport.HTTP.Port = *port
	config.Log.Level = *logLevel

	// 创建并启动服务器
	if err := runServer(config); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// runServer 运行服务器
func runServer(config *server.Config) error {
	// 创建MCP开发工具服务器
	mcpServer, err := server.NewMCPDevServer(config)
	if err != nil {
		return fmt.Errorf("创建服务器失败: %w", err)
	}

	// 创建上下文和信号处理
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	serverErrChan := make(chan error, 1)
	go func() {
		serverErrChan <- mcpServer.Start(ctx)
	}()

	// 等待信号或错误
	select {
	case sig := <-sigChan:
		log.Printf("收到信号 %v，正在关闭服务器...", sig)
		cancel()
		return mcpServer.Stop()
	case err := <-serverErrChan:
		if err != nil {
			return fmt.Errorf("服务器运行错误: %w", err)
		}
		return nil
	}
}

// printHelp 打印帮助信息
func printHelp() {
	fmt.Println("Bico Admin DevTools MCP Server")
	fmt.Println("基于 Model Context Protocol 的开发工具服务 (HTTP模式)")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  devtools [选项]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -host string")
	fmt.Println("        HTTP服务器监听地址 (默认: 127.0.0.1)")
	fmt.Println("  -port int")
	fmt.Println("        HTTP服务器监听端口 (默认: 18901)")
	fmt.Println("  -log-level string")
	fmt.Printf("        日志级别，支持 debug, info, warn, error (默认: %s)\n", defaultLogLevel)
	fmt.Println("  -version")
	fmt.Println("        显示版本信息")
	fmt.Println("  -help")
	fmt.Println("        显示此帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  # 使用默认配置启动")
	fmt.Println("  devtools")
	fmt.Println()
	fmt.Println("  # 指定监听地址和端口")
	fmt.Println("  devtools -host 0.0.0.0 -port 18901")
	fmt.Println()
	fmt.Println("  # 启用调试日志")
	fmt.Println("  devtools -log-level debug")
	fmt.Println()
	fmt.Println("可用工具:")
	fmt.Println("  - read_config: 读取并解析应用配置文件")
	fmt.Println("  - execute_sql: 执行SQL语句（仅开发环境可用）")
	fmt.Println("  - 更多工具正在开发中...")
	fmt.Println()
	fmt.Println("注意:")
	fmt.Println("  此工具仅用于开发环境，不应在生产环境中使用。")
	fmt.Println("  服务器将通过HTTP接口提供MCP服务，访问地址: http://host:port/mcp")
}
