# =============================================================================
# 多阶段构建 Dockerfile for bico-admin
# 阶段1: 前端构建
# 阶段2: 后端构建  
# 阶段3: 生产运行环境
# =============================================================================

# -----------------------------------------------------------------------------
# 阶段1: 前端构建
# -----------------------------------------------------------------------------
FROM node:22-alpine AS frontend-builder

# 设置工作目录
WORKDIR /app/web

# 安装 pnpm
RUN npm install -g pnpm

# 复制前端依赖文件
COPY web/package.json web/pnpm-lock.yaml ./

# 安装前端依赖 (利用 Docker 缓存层)
RUN pnpm install --frozen-lockfile

# 复制前端源码
COPY web/ ./

# 构建前端项目
RUN pnpm run build

# -----------------------------------------------------------------------------
# 阶段2: 后端构建
# -----------------------------------------------------------------------------
FROM golang:1.23-alpine AS backend-builder

# 安装必要的构建工具
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 设置 Go 环境变量 (纯 Go SQLite 驱动，无需 CGO)
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# 复制 Go modules 文件 (利用 Docker 缓存层)
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 安装 wire 工具
RUN go install github.com/google/wire/cmd/wire@latest

# 复制源码
COPY . .

# 生成 Wire 代码
RUN wire ./cmd/server

# 构建应用
RUN go build -ldflags="-s -w" -o bico-admin ./cmd/server

# -----------------------------------------------------------------------------
# 阶段3: 生产运行环境
# -----------------------------------------------------------------------------
FROM alpine:3.19

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata curl && \
    update-ca-certificates

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 创建非 root 用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 创建必要的目录结构
RUN mkdir -p /app/config \
             /app/data/uploads \
             /app/logs \
             /app/web/dist && \
    chown -R appuser:appgroup /app

# 从构建阶段复制文件
COPY --from=backend-builder /app/bico-admin /app/
COPY --from=frontend-builder /app/web/dist /app/web/dist/
COPY --chown=appuser:appgroup config/app.yml /app/config/

# 设置文件权限
RUN chmod +x /app/bico-admin && \
    chown -R appuser:appgroup /app

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8899

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8899/admin-api/auth/login || exit 1

# 设置环境变量
ENV BICO_APP_ENVIRONMENT=production \
    BICO_APP_DEBUG=false \
    BICO_LOG_LEVEL=info \
    BICO_LOG_FORMAT=json \
    BICO_LOG_OUTPUT=stdout \
    BICO_SERVER_HOST=0.0.0.0 \
    BICO_SERVER_PORT=8899

# 数据卷
VOLUME ["/app/data", "/app/logs", "/app/config"]

# 启动命令
CMD ["/app/bico-admin", "-config", "/app/config/app.yml"]
