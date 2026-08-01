.PHONY: help serve air dev tidy install migrate build web build-web package package-win clean swagger

export GOROOT :=

help:
	@echo "可用命令:"
	@echo "  make serve     - 启动后端服务"
	@echo "  make air       - 使用 air 热重载启动后端服务"
	@echo "  make dev       - 同时启动前后端开发服务"
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

dev:
	@if nc -z localhost 8080 2>/dev/null; then echo "端口 8080 已被占用，请先停止已有后端服务"; exit 1; fi
	@echo "🚀 启动后端开发服务..."
	@$(MAKE) air & \
	backend_pid=$$!; \
	frontend_pid=; \
	watcher_pid=; \
	stop_process_tree() { \
		for child_pid in $$(pgrep -P $$1 2>/dev/null); do \
			stop_process_tree $$child_pid; \
		done; \
		kill -TERM $$1 2>/dev/null || true; \
	}; \
	cleanup() { \
		kill -TERM $$watcher_pid 2>/dev/null || true; \
		wait $$watcher_pid 2>/dev/null || true; \
		stop_process_tree $$frontend_pid; \
		stop_process_tree $$backend_pid; \
		wait $$backend_pid 2>/dev/null || true; \
	}; \
	trap cleanup EXIT; \
	trap 'exit 130' INT TERM; \
	echo "⏳ 等待后端服务就绪..."; \
	while ! nc -z localhost 8080 2>/dev/null; do \
		if ! kill -0 $$backend_pid 2>/dev/null; then \
			wait $$backend_pid; \
			exit 1; \
		fi; \
		sleep 0.5; \
		done; \
	echo "✅ 后端已就绪，启动前端..."; \
	$(MAKE) web & \
	frontend_pid=$$!; \
	watch_parent() { \
		while kill -0 $$1 2>/dev/null; do sleep 1; done; \
		stop_process_tree $$2; \
		stop_process_tree $$3; \
	}; \
	watch_parent $$PPID $$backend_pid $$frontend_pid & \
	watcher_pid=$$!; \
	wait $$frontend_pid

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
	@swag init -g doc.go -d docs/admin,internal/admin -o docs/admin --instanceName admin
	@swag init -g doc.go -d docs/api,internal/api -o docs/api --instanceName api
	@go run ./cmd/swagger-enhance
	@echo "✅ Swagger 文档生成完成"

build:
	@echo "🔨 编译后端..."
	@go build -o bin/bico-admin ./cmd/main.go
	@echo "✅ 编译完成: bin/bico-admin"

web:
	@echo "🚀 启动前端开发服务器..."
	@cd web && pnpm run dev

build-web:
	@echo "🎨 构建前端..."
	@cd web && pnpm run build
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
