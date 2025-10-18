.PHONY: help build build-embed build-web serve migrate clean tidy

help:
	@echo "可用命令:"
	@echo "  make build        - 编译应用（开发模式）"
	@echo "  make build-embed  - 编译应用（嵌入前端资源）"
	@echo "  make build-web    - 构建前端"
	@echo "  make serve        - 启动服务"
	@echo "  make migrate      - 执行数据库迁移"
	@echo "  make clean        - 清理构建产物"
	@echo "  make tidy         - 整理依赖"

build:
	@echo "🔨 开始编译（开发模式）..."
	@go build -o bin/bico-admin ./cmd/main.go
	@echo "✅ 编译完成: bin/bico-admin"

build-web:
	@echo "🎨 构建前端..."
	@cd web && npm run build
	@echo "✅ 前端构建完成"

build-embed: build-web
	@echo "🔨 开始编译（嵌入模式）..."
	@go build -tags embed -ldflags="-s -w" -o bin/bico-admin ./cmd/main.go
	@echo "✅ 编译完成: bin/bico-admin（已嵌入前端资源）"

serve:
	@go run cmd/main.go serve

migrate:
	@go run cmd/main.go migrate

clean:
	@echo "🧹 清理构建产物..."
	@rm -rf bin/ web/dist/
	@echo "✅ 清理完成"

tidy:
	@echo "📦 整理依赖..."
	@go mod tidy
	@echo "✅ 依赖整理完成"
