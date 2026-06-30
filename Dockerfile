# 阶段 1: 前端构建
FROM oven/bun:1-alpine AS web-builder

WORKDIR /app/web

COPY web/package.json web/bun.lock ./
RUN bun install --frozen-lockfile --registry=https://registry.npmmirror.com

COPY web/ ./
RUN bun run build


# 阶段 2: 后端构建
FROM golang:1.25-alpine AS go-builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

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
