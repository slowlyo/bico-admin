# 前端静态文件打包功能

## 功能说明

将前端构建产物嵌入到 Go 二进制文件中，实现单文件部署。

## 配置

在 `config/config.yaml` 中配置：

```yaml
server:
  embed_static: false  # 是否将前端静态文件打包到二进制文件中
```

- `false`（默认）：不嵌入，使用常规模式开发
- `true`：嵌入静态文件，需要使用 embed 编译标签

## 使用步骤

### 1. 构建前端

```bash
cd web
npm run build
```

这会在 `web/dist` 目录生成前端构建产物。

### 2. 启用 embed 配置

修改 `config/config.yaml`：

```yaml
server:
  embed_static: true
```

### 3. 使用 embed 标签编译

```bash
go build -tags embed -o bico-admin cmd/main.go
```

编译完成后，二进制文件已包含前端资源。

### 4. 部署运行

将单个二进制文件和配置文件复制到服务器：

```bash
./bico-admin serve -c config/config.yaml
```

访问 `http://localhost:8080` 即可看到前端页面。

## 工作原理

### 路由模式

- **前端**：使用 hash 路由模式（`#/dashboard`）
- **后端**：通过 `NoRoute` 处理所有非 API 请求，返回前端资源

### 条件编译

项目使用 Go 的 build tags 实现条件编译：

- **`web/embed.go`**（`//go:build embed`）：嵌入模式，包含真实的静态文件
- **`web/embed_stub.go`**（`//go:build !embed`）：开发模式，空的 embed.FS

### API 路由保护

嵌入模式下，以下路由仍走后端 API：

- `/admin-api/*` - 后台管理 API
- `/api/*` - 业务 API
- `/uploads/*` - 文件上传
- `/health` - 健康检查

## 注意事项

1. **开发模式**
   - 关闭 `embed_static`
   - 前端使用 `npm run dev` 独立运行
   - 后端通过 CORS 支持前后端分离开发

2. **生产模式**
   - 开启 `embed_static`
   - 使用 `-tags embed` 编译
   - 前端必须先执行 `npm run build`

3. **文件体积**
   - 嵌入后二进制文件会增大（通常 2-5 MB）
   - 可通过 UPX 压缩二进制文件

4. **更新前端**
   - 每次前端修改后需要重新 build 和编译
   - 开发阶段建议关闭 embed 模式

## Makefile 示例

```makefile
# 开发模式
dev:
	go run cmd/main.go serve

# 生产构建
build-prod:
	cd web && npm run build
	go build -tags embed -ldflags="-s -w" -o dist/bico-admin cmd/main.go

# 清理
clean:
	rm -rf web/dist dist/
```

## 故障排查

### 问题：编译报错 "pattern dist: no matching files found"

**原因**：前端未构建或 `web/dist` 目录不存在

**解决**：先执行 `cd web && npm run build`

### 问题：访问前端显示 404

**原因**：
1. `embed_static` 未开启
2. 编译时未使用 `-tags embed`

**解决**：检查配置并重新编译

### 问题：前端路由刷新后 404

**原因**：未使用 hash 路由模式

**解决**：前端已配置 hash 模式，确保前端重新构建
