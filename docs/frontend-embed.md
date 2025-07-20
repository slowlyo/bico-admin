# 前端嵌入模式实现

## 概述

本项目实现了前端文件的嵌入模式，支持将前端静态文件直接嵌入到 Go 二进制文件中，实现单文件部署。

## 架构设计

### 1. 构建标签 (Build Tags)

使用 Go 的构建标签来实现条件编译：

- **嵌入模式**: `//go:build embed`
- **外部文件模式**: `//go:build !embed`

### 2. 文件结构

```
web/
├── embed.go      # 嵌入模式实现 (embed 标签)
├── external.go   # 外部文件模式实现 (!embed 标签)
└── dist/         # 前端构建产物
```

### 3. 核心组件

#### web 包

- `web.GetFileSystem()`: 获取前端文件系统
- `web.IsEmbedded()`: 检查是否为嵌入模式

#### frontend 包

- `frontend.Service`: 前端服务
- `setupEmbeddedRoutes()`: 设置嵌入模式路由
- `setupExternalRoutes()`: 设置外部文件模式路由

## 使用方法

### 1. 开发模式 (外部文件)

```bash
# 启动后端服务
make run

# 或直接运行
go run ./cmd/server
```

### 2. 嵌入模式

```bash
# 构建嵌入模式二进制文件
make build-embed

# 运行嵌入模式
make run-embed

# 或手动运行
BICO_FRONTEND_MODE=embed go run -tags embed ./cmd/server
```

### 3. 配置选项

通过环境变量 `BICO_FRONTEND_MODE` 控制前端模式：

- `external`: 外部文件模式 (默认)
- `embed`: 嵌入模式

## 技术实现

### 1. 嵌入文件 (web/embed.go)

```go
//go:build embed

//go:embed all:dist
var EmbeddedFiles embed.FS

func GetFileSystem() (http.FileSystem, error) {
    distFS, err := fs.Sub(EmbeddedFiles, "dist")
    if err != nil {
        return nil, err
    }
    return http.FS(distFS), nil
}

func IsEmbedded() bool {
    return true
}
```

### 2. 外部文件 (web/external.go)

```go
//go:build !embed

func GetFileSystem() (http.FileSystem, error) {
    return nil, nil
}

func IsEmbedded() bool {
    return false
}
```

### 3. 前端服务路由

根据模式自动选择路由设置：

- **嵌入模式**: 从嵌入的文件系统提供静态文件
- **外部模式**: 从磁盘文件系统提供静态文件

## 优势

### 嵌入模式
- ✅ 单文件部署
- ✅ 无需额外文件依赖
- ✅ 便于分发和部署
- ❌ 二进制文件较大
- ❌ 前端更新需要重新编译

### 外部文件模式
- ✅ 二进制文件较小
- ✅ 前端可独立更新
- ✅ 开发调试方便
- ❌ 需要额外的文件依赖
- ❌ 部署时需要多个文件

## 构建流程

### 1. 前端构建

```bash
cd web && pnpm build
```

### 2. 后端构建

```bash
# 外部文件模式
go build -o bico-admin ./cmd/server

# 嵌入模式
go build -tags embed -o bico-admin ./cmd/server
```

## 注意事项

1. **构建顺序**: 嵌入模式必须先构建前端，再构建后端
2. **文件路径**: 嵌入路径必须相对于包含 `//go:embed` 指令的文件
3. **构建标签**: 确保使用正确的构建标签来选择模式
4. **配置优先级**: 环境变量 > 配置文件设置

## 故障排除

### 1. 嵌入文件未找到

检查前端是否已构建：
```bash
ls -la web/dist/
```

### 2. 路径错误

确保 `//go:embed` 指令中的路径正确：
```go
//go:embed all:dist  // 相对于当前文件的路径
```

### 3. 构建标签问题

确保文件头部有正确的构建标签：
```go
//go:build embed
// +build embed
```

## 相关命令

```bash
# 查看帮助
make help

# 构建和运行
make build-embed    # 构建嵌入模式
make run-embed      # 运行嵌入模式
make build          # 构建外部文件模式
make run            # 运行外部文件模式

# 清理
make clean          # 清理构建文件
```
