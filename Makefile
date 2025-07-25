# 项目参数
BINARY_NAME=bico-admin
MAIN_PATH=./cmd/server
DEVTOOLS_PATH=./cmd/devtools
WEB_DIR=web

# 构建目标
.PHONY: run dev build clean test deps wire help build-devtools docker-build docker-dev docker-stop build-web dev-web build-embed build-windows-embed

# 默认目标
default: help

# 启动后端服务
run: wire
	go run $(MAIN_PATH)

# 启动后端服务 (嵌入模式)
run-embed: wire build-web
	BICO_FRONTEND_MODE=embed go run -tags embed $(MAIN_PATH)

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

# 构建前端
build-web:
	@echo "构建前端..."
	cd $(WEB_DIR) && pnpm build

# 开发模式启动前端
dev-web:
	@echo "启动前端开发服务器..."
	cd $(WEB_DIR) && pnpm dev

# 构建二进制文件 (外部文件模式)
build: wire build-web
	@echo "构建应用 (外部文件模式)..."
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o $(BINARY_NAME) -v $(MAIN_PATH)

# 构建二进制文件 (嵌入模式)
build-embed: wire build-web
	@echo "构建应用 (嵌入模式)..."
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -tags embed -o $(BINARY_NAME) -v $(MAIN_PATH)

# 构建Windows二进制文件 (嵌入模式)
build-windows-embed: wire build-web
	@echo "构建Windows应用 (嵌入模式)..."
	GOOS=windows GOARCH=amd64 go build -tags embed -o $(BINARY_NAME).exe -v $(MAIN_PATH)

# 清理
clean:
	@echo "清理文件..."
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe
	rm -f bin/devtools
	rm -rf $(WEB_DIR)/dist

# MCP开发工具相关命令

# 构建MCP开发工具
build-devtools:
	@echo "构建MCP开发工具..."
	@mkdir -p bin
	@go build -o bin/devtools ./cmd/devtools
	@echo ""
	@echo "🎉 MCP开发工具构建完成！"
	@echo ""
	@echo "📋 MCP客户端配置 (复制到你的MCP客户端配置文件中):"
	@echo "{"
	@echo "  \"mcpServers\": {"
	@echo "    \"bico-admin-devtools\": {"
	@echo "      \"command\": \"$(PWD)/bin/devtools\","
	@echo "      \"cwd\": \"$(PWD)\""
	@echo "    }"
	@echo "  }"
	@echo "}"
	@echo ""

# Docker 相关命令

# 构建 Docker 镜像
docker-build:
	@echo "构建 Docker 镜像..."
	./scripts/docker-build.sh

# 启动开发环境 Docker 服务
docker-dev:
	@echo "启动开发环境 Docker 服务..."
	docker-compose up -d

# 停止 Docker 服务
docker-stop:
	@echo "停止 Docker 服务..."
	docker-compose down

# 帮助信息
help:
	@echo "可用命令:"
	@echo ""
	@echo "主应用命令:"
	@echo "  make run         - 启动后端服务"
	@echo "  make run-embed   - 启动后端服务 (嵌入模式)"
	@echo "  make dev         - 开发模式 (run的别名)"
	@echo "  make build       - 构建完整应用 (外部文件模式)"
	@echo "  make build-embed - 构建完整应用 (嵌入模式)"
	@echo "  make build-windows-embed - 构建Windows应用 (嵌入模式)"
	@echo "  make build-web   - 仅构建前端"
	@echo "  make dev-web     - 启动前端开发服务器"
	@echo "  make test        - 运行测试"
	@echo "  make deps        - 安装依赖"
	@echo "  make wire        - 生成Wire代码"
	@echo "  make clean       - 清理文件"
	@echo ""
	@echo "MCP开发工具命令:"
	@echo "  make build-devtools - 构建MCP开发工具并输出配置"
	@echo ""
	@echo "Docker命令:"
	@echo "  make docker-build    - 构建Docker镜像"
	@echo "  make docker-dev      - 启动开发环境Docker服务"
	@echo "  make docker-stop     - 停止Docker服务"
	@echo ""
	@echo "  make help    - 显示帮助信息"