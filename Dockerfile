# 阶段 1: 前端构建
FROM node:20-alpine AS web-builder

WORKDIR /app/web

# 安装 pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

COPY web/package.json web/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile --registry=https://registry.npmmirror.com

COPY web/ ./
RUN pnpm build


# 阶段 2: 后端构建
FROM golang:1.24-alpine AS go-builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=web-builder /app/web/dist ./web/dist

RUN CGO_ENABLED=0 GOOS=linux go build \
    -tags embed \
    -ldflags="-s -w" \
    -o /app/bin/bico-admin \
    ./cmd/main.go


# 阶段 3: 运行时镜像
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

COPY --from=go-builder /app/bin/bico-admin .

EXPOSE 8080

CMD ["./bico-admin", "serve"]
