# Go参数
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
GOMOD=$(GOCMD) mod
GOGET=$(GOCMD) get

# 项目参数
BINARY_NAME=bico-admin
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=./cmd/server
DATA_DIR=data
LOGS_DIR=logs

# 构建目标
.PHONY: all build build-linux clean test deps wire run dev prod help init default

# 默认目标
default: help

# 完整构建流程
all: deps wire test build

# 初始化项目环境
init:
	@echo "初始化项目环境..."
	@mkdir -p $(DATA_DIR) $(LOGS_DIR)
	@echo "创建目录: $(DATA_DIR), $(LOGS_DIR)"

# 安装依赖
deps:
	@echo "下载依赖..."
	$(GOMOD) download
	$(GOMOD) tidy

# 生成Wire依赖注入代码
wire:
	@echo "生成Wire代码..."
	@wire $(MAIN_PATH)

# 运行测试
test:
	@echo "运行测试..."
	$(GOTEST) -v ./...

# 开发模式运行 (使用go run)
run: wire init
	@echo "开发模式启动..."
	$(GORUN) $(MAIN_PATH)

# 开发模式 (别名)
dev: run

# 构建二进制文件
build: wire
	@echo "构建应用..."
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

# Linux构建
build-linux: wire
	@echo "构建Linux版本..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PATH)

# 生产模式运行 (先构建再运行)
prod: build init
	@echo "生产模式启动..."
	./$(BINARY_NAME)

# 清理
clean:
	@echo "清理文件..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -rf $(DATA_DIR)/*.db
	rm -rf $(LOGS_DIR)/*.log

# 开发工具
install-tools:
	@echo "安装开发工具..."
	$(GOGET) github.com/google/wire/cmd/wire@latest
	$(GOGET) github.com/swaggo/swag/cmd/swag@latest

# 生成API文档
swag:
	@echo "生成API文档..."
	swag init -g $(MAIN_PATH)/main.go

# 格式化代码
fmt:
	@echo "格式化代码..."
	$(GOCMD) fmt ./...

# 代码检查
vet:
	@echo "代码检查..."
	$(GOCMD) vet ./...

# 完整检查
check: fmt vet test
	@echo "代码检查完成"

# Docker相关
docker-build:
	@echo "构建Docker镜像..."
	docker build -t $(BINARY_NAME) .

docker-run:
	@echo "运行Docker容器..."
	docker run -p 8080:8080 $(BINARY_NAME)

# 日志相关
logs-clean:
	@echo "清理日志..."
	rm -f $(LOGS_DIR)/*.log

# 帮助信息
help:
	@echo "可用命令:"
	@echo "  make          - 显示帮助信息 (默认)"
	@echo "  make all      - 完整构建流程 (deps + wire + test + build)"
	@echo "  make run      - 开发模式运行 (使用go run)"
	@echo "  make dev      - 开发模式运行 (run的别名)"
	@echo "  make prod     - 生产模式运行 (先构建再运行)"
	@echo "  make build    - 构建二进制文件"
	@echo "  make test     - 运行测试"
	@echo "  make deps     - 安装依赖"
	@echo "  make wire     - 生成Wire代码"
	@echo "  make clean    - 清理文件"
	@echo "  make init     - 初始化项目环境"
	@echo "  make check    - 完整代码检查"
	@echo "  make fmt      - 格式化代码"
	@echo "  make vet      - 代码静态检查"
	@echo "  make swag     - 生成API文档"
	@echo "  make help     - 显示帮助信息"
