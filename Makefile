# 项目参数
BINARY_NAME=bico-admin
MAIN_PATH=./cmd/server
DEVTOOLS_PATH=./cmd/devtools
WEB_DIR=web

# 构建目标
.PHONY: run dev build clean test deps wire help devtools devtools-stdio devtools-sse devtools-help devtools-version

# 默认目标
default: help

# 启动后端服务
run: wire
	go run $(MAIN_PATH) 

# 开发模式 (别名)
dev: run

# 安装依赖
deps:
	@echo "安装依赖..."
	go mod download && go mod tidy
	@echo "检查并安装开发工具..."
	@which wire > /dev/null 2>&1 || (echo "安装 wire..." && go install github.com/google/wire/cmd/wire@latest)

	@echo "安装 MCP 依赖..."
	go get github.com/mark3labs/mcp-go@latest
	@which pnpm > /dev/null 2>&1 || (echo "请先安装 pnpm: npm install -g pnpm" && exit 1)
	cd $(WEB_DIR) && pnpm install

# 生成Wire代码
wire:
	@echo "生成Wire代码..."
	@wire $(MAIN_PATH)



# 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

# 构建二进制文件
build: wire
	@echo "构建应用..."
	go build -o $(BINARY_NAME) -v $(MAIN_PATH)

# 清理
clean:
	@echo "清理文件..."
	go clean
	rm -f $(BINARY_NAME)
	rm -f bin/devtools

# MCP开发工具相关命令

# 启动MCP开发工具 (HTTP模式)
devtools:
	@echo "启动MCP开发工具 (HTTP模式)..."
	go run $(DEVTOOLS_PATH) -host 0.0.0.0 -port 18901

# 显示MCP开发工具帮助
devtools-help:
	@echo "MCP开发工具帮助..."
	go run $(DEVTOOLS_PATH) -help

# 显示MCP开发工具版本
devtools-version:
	@echo "MCP开发工具版本..."
	go run $(DEVTOOLS_PATH) -version



# 帮助信息
help:
	@echo "可用命令:"
	@echo ""
	@echo "主应用命令:"
	@echo "  make run     - 同时启动前后端服务"
	@echo "  make dev     - 开发模式 (run的别名)"
	@echo "  make build   - 构建二进制文件"
	@echo "  make test    - 运行测试"
	@echo "  make deps    - 安装依赖"
	@echo "  make wire    - 生成Wire代码"

	@echo "  make clean   - 清理文件"
	@echo ""
	@echo "MCP开发工具命令:"
	@echo "  make devtools        - 启动MCP开发工具 (HTTP模式，默认)"
	@echo "  make devtools-sse    - 启动MCP开发工具 (SSE模式)"
	@echo "  make devtools-stdio  - 启动MCP开发工具 (stdio模式)"
	@echo "  make devtools-help   - 显示MCP开发工具帮助"
	@echo "  make devtools-version - 显示MCP开发工具版本"
	@echo ""
	@echo "注意: 推荐使用HTTP模式，通过URL连接更简单可靠"
	@echo ""
	@echo "  make help    - 显示帮助信息"
