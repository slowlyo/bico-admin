.PHONY: help build serve migrate clean tidy

help:
	@echo "可用命令:"
	@echo "  make build     - 编译应用"
	@echo "  make serve     - 启动服务"
	@echo "  make migrate   - 执行数据库迁移"
	@echo "  make clean     - 清理构建产物"
	@echo "  make tidy      - 整理依赖"

build:
	@echo "🔨 开始编译..."
	@go build -o bin/bico-admin ./cmd/main.go
	@echo "✅ 编译完成: bin/bico-admin"

serve:
	@go run cmd/main.go serve

migrate:
	@go run cmd/main.go migrate

clean:
	@echo "🧹 清理构建产物..."
	@rm -rf bin/
	@echo "✅ 清理完成"

tidy:
	@echo "📦 整理依赖..."
	@go mod tidy
	@echo "✅ 依赖整理完成"
