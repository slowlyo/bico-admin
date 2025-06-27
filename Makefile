# Bico Admin Makefile
# 用于管理前后端服务的统一命令

.PHONY: help install dev build clean start stop restart logs test

# 默认目标
.DEFAULT_GOAL := help

# 颜色定义
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

# 帮助信息
help: ## 显示帮助信息
	@echo "$(GREEN)Bico Admin 项目管理命令$(RESET)"
	@echo ""
	@echo "$(YELLOW)可用命令:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# 安装依赖
install: ## 安装前后端依赖
	@echo "$(GREEN)安装后端依赖...$(RESET)"
	cd backend && go mod tidy
	@echo "$(GREEN)安装前端依赖...$(RESET)"
	cd frontend && pnpm install

# 开发环境
dev: ## 启动开发环境（前后端同时启动）
	@echo "$(GREEN)启动开发环境...$(RESET)"
	@make -j2 dev-backend dev-frontend

dev-backend: ## 启动后端开发服务
	@echo "$(GREEN)启动后端服务...$(RESET)"
	cd backend && go run cmd/server/main.go

dev-frontend: ## 启动前端开发服务
	@echo "$(GREEN)启动前端服务...$(RESET)"
	cd frontend && pnpm dev

# 构建
build: ## 构建前后端项目
	@echo "$(GREEN)构建项目...$(RESET)"
	@make build-backend
	@make build-frontend

build-backend: ## 构建后端
	@echo "$(GREEN)构建后端...$(RESET)"
	cd backend && go build -o bin/server cmd/server/main.go
	@echo "$(GREEN)后端构建完成: backend/bin/server$(RESET)"

build-frontend: ## 构建前端
	@echo "$(GREEN)构建前端...$(RESET)"
	cd frontend && pnpm build
	@echo "$(GREEN)前端构建完成: frontend/dist$(RESET)"

# 生产环境启动
start: build ## 启动生产环境服务
	@echo "$(GREEN)启动生产环境...$(RESET)"
	@make start-backend &
	@make start-frontend &

start-backend: ## 启动后端生产服务
	@echo "$(GREEN)启动后端生产服务...$(RESET)"
	cd backend && ./bin/server

start-frontend: ## 启动前端生产服务
	@echo "$(GREEN)启动前端生产服务...$(RESET)"
	cd frontend && pnpm preview

# 停止服务
stop: ## 停止所有服务
	@echo "$(RED)停止所有服务...$(RESET)"
	@pkill -f "go run cmd/server/main.go" || true
	@pkill -f "pnpm dev" || true
	@pkill -f "pnpm preview" || true
	@pkill -f "bin/server" || true

# 重启服务
restart: stop dev ## 重启开发服务

# 清理
clean: ## 清理构建文件
	@echo "$(RED)清理构建文件...$(RESET)"
	rm -rf backend/bin
	rm -rf frontend/dist
	rm -rf frontend/node_modules/.vite

# 测试
test: ## 运行测试
	@echo "$(GREEN)运行后端测试...$(RESET)"
	cd backend && go test ./...
	@echo "$(GREEN)运行前端测试...$(RESET)"
	cd frontend && pnpm test || echo "前端测试未配置"

# 代码检查
lint: ## 代码检查
	@echo "$(GREEN)后端代码检查...$(RESET)"
	cd backend && go fmt ./...
	cd backend && go vet ./...
	@echo "$(GREEN)前端代码检查...$(RESET)"
	cd frontend && pnpm lint

# 数据库相关
db-migrate: ## 运行数据库迁移
	@echo "$(GREEN)运行数据库迁移...$(RESET)"
	cd backend && go run migrations/migrate.go

db-seed: ## 填充测试数据
	@echo "$(GREEN)填充测试数据...$(RESET)"
	cd backend && go run migrations/seed.go

# Docker相关
docker-build: ## 构建Docker镜像
	@echo "$(GREEN)构建Docker镜像...$(RESET)"
	docker-compose build

docker-up: ## 启动Docker服务
	@echo "$(GREEN)启动Docker服务...$(RESET)"
	docker-compose up -d

docker-down: ## 停止Docker服务
	@echo "$(RED)停止Docker服务...$(RESET)"
	docker-compose down

docker-logs: ## 查看Docker日志
	docker-compose logs -f

# 日志查看
logs: ## 查看服务日志
	@echo "$(GREEN)查看服务日志...$(RESET)"
	tail -f backend/storage/logs/*.log || echo "日志文件不存在"

# 状态检查
status: ## 检查服务状态
	@echo "$(GREEN)检查服务状态...$(RESET)"
	@echo "后端服务:"
	@curl -s http://localhost:8080/health || echo "后端服务未运行"
	@echo ""
	@echo "前端服务:"
	@curl -s http://localhost:5174/ > /dev/null && echo "前端服务运行中" || echo "前端服务未运行"

# 快速命令别名
up: dev ## 启动开发环境（dev的别名）
down: stop ## 停止服务（stop的别名）

# 环境配置
env: ## 创建环境配置文件
	@echo "$(GREEN)创建环境配置文件...$(RESET)"
	@if [ ! -f backend/.env ]; then \
		cp backend/.env.example backend/.env; \
		echo "已创建 backend/.env，请编辑配置"; \
	else \
		echo "backend/.env 已存在"; \
	fi
	@if [ ! -f frontend/.env ]; then \
		cp frontend/.env.example frontend/.env; \
		echo "已创建 frontend/.env，请编辑配置"; \
	else \
		echo "frontend/.env 已存在"; \
	fi

# 项目初始化
init: install env ## 初始化项目（安装依赖+创建配置文件）
	@echo "$(GREEN)项目初始化完成！$(RESET)"
	@echo "$(YELLOW)下一步:$(RESET)"
	@echo "1. 编辑 backend/.env 配置数据库连接"
	@echo "2. 编辑 frontend/.env 配置API地址"
	@echo "3. 运行 make dev 启动开发环境"
