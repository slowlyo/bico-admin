# 前端静态资源嵌入说明

## 功能

将 `web/dist` 嵌入 Go 二进制，实现单文件部署。

对应文件：

- `web/embed.go`（`//go:build embed`）
- `web/embed_stub.go`（`//go:build !embed`）

## 相关配置

```yaml
server:
  embed_static: true
  admin_path: /admin
```

说明：

- `embed_static=true` 时，后端会注册前端静态资源路由
- `admin_path=/admin` 时，前端入口为 `/admin/`

## 使用方式

### 1. 分离开发（推荐本地）

```yaml
server:
  embed_static: false
```

启动：

```bash
make serve
cd web && pnpm dev
```

### 2. 嵌入发布

```bash
# 方式一：直接命令
cd web && pnpm build
cd ..
go build -tags embed -o bin/bico-admin ./cmd/main.go

# 方式二：Makefile
make package
```

运行：

```bash
./bin/bico-admin serve -c config/config.yaml
```

访问：`http://localhost:8080/admin/`

## 路由处理说明

嵌入模式下：

1. `/admin-api/*`、`/api/*` 仍走后端接口。
2. `admin_path` 下的前端路由由静态资源处理器兜底到 `index.html`。
3. 访问 `/admin` 会重定向到 `/admin/`。

## 常见问题

### 编译报错：`pattern dist: no matching files found`

先执行前端构建：

```bash
cd web && pnpm build
```

### 访问前端 404

检查三项：

1. `server.embed_static` 是否为 `true`
2. 是否使用 `-tags embed` 编译
3. `server.admin_path` 与访问路径是否一致
