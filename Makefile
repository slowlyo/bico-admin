# 项目参数
BINARY_NAME=bico-admin
MAIN_PATH=./cmd/server
WEB_DIR=web

# 构建目标
.PHONY: run dev build clean test deps wire swagger help

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
	@which swag > /dev/null 2>&1 || (echo "安装 swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	@which pnpm > /dev/null 2>&1 || (echo "请先安装 pnpm: npm install -g pnpm" && exit 1)
	cd $(WEB_DIR) && pnpm install

# 生成Wire代码
wire:
	@echo "生成Wire代码..."
	@wire $(MAIN_PATH)

# 生成Swagger文档
swagger:
	@echo "生成Swagger文档..."
	@swag init -g $(MAIN_PATH)/main.go -o ./docs --parseDependency --parseInternal

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

# 帮助信息
help:
	@echo "可用命令:"
	@echo "  make run     - 同时启动前后端服务"
	@echo "  make dev     - 开发模式 (run的别名)"
	@echo "  make build   - 构建二进制文件"
	@echo "  make test    - 运行测试"
	@echo "  make deps    - 安装依赖"
	@echo "  make wire    - 生成Wire代码"
	@echo "  make swagger - 生成Swagger文档"
	@echo "  make clean   - 清理文件"
	@echo "  make help    - 显示帮助信息"
