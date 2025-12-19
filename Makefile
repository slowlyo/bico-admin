.PHONY: help serve air tidy install migrate build web build-web package package-win clean swagger

help:
	@echo "可用命令:"
	@echo "  make serve     - 启动后端服务"
	@echo "  make air       - 使用 air 热重载启动后端服务"
	@echo "  make web       - 启动前端开发服务器"
	@echo "  make build     - 编译后端"
	@echo "  make build-web - 编译前端"
	@echo "  make package     - 构建生产版本（嵌入前端）"
	@echo "  make package-win - 构建 Windows 版本（嵌入前端）"
	@echo "  make install     - 安装前端依赖"
	@echo "  make migrate   - 执行数据库迁移"
	@echo "  make swagger   - 生成 Swagger 文档"
	@echo "  make tidy      - 整理后端依赖"
	@echo "  make clean     - 清理构建产物"

serve:
	@go run cmd/main.go serve

air:
	@command -v air >/dev/null 2>&1 || (echo "❌ 未安装 air，请先安装: go install github.com/air-verse/air@latest" && exit 1)
	@air -c .air.toml

tidy:
	@echo "📦 整理依赖..."
	@go mod tidy
	@echo "✅ 依赖整理完成"

install:
	@echo "📦 安装前端依赖..."
	@cd web && pnpm install
	@echo "✅ 前端依赖安装完成"

migrate:
	@go run cmd/main.go migrate

swagger:
	@echo "📝 生成 Swagger 文档..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@swag init -g cmd/main.go -o docs
	@echo "✅ Swagger 文档生成完成"

build:
	@echo "🔨 编译后端..."
	@go build -o bin/bico-admin ./cmd/main.go
	@echo "✅ 编译完成: bin/bico-admin"

web:
	@echo "🚀 启动前端开发服务器..."
	@cd web && npm run dev

build-web:
	@echo "🎨 构建前端..."
	@cd web && npm run build
	@echo "✅ 前端构建完成"

package: build-web
	@echo "🔨 构建生产版本（嵌入前端）..."
	@go build -tags embed -ldflags="-s -w" -o bin/bico-admin ./cmd/main.go
	@echo "✅ 构建完成: bin/bico-admin"

package-win: build-web
	@echo "🔨 构建 Windows 版本（嵌入前端）..."
	@GOOS=windows GOARCH=amd64 go build -tags embed -ldflags="-s -w" -o bin/bico-admin.exe ./cmd/main.go
	@echo "✅ 构建完成: bin/bico-admin.exe"

clean:
	@echo "🧹 清理构建产物..."
	@rm -rf bin/ web/dist/ web/node_modules/.cache
	@echo "✅ 清理完成"
