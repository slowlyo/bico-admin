# 项目结构说明

本文档描述当前代码目录与职责边界。

## 目录总览

```text
bico-admin/
├── cmd/
│   └── main.go                # Cobra 入口（serve/migrate）
├── config/
│   ├── config.yaml            # 默认配置
│   └── config.prod.yaml       # 生产配置模板
├── internal/
│   ├── admin/                 # 后台模块（认证、用户、角色、上传、权限）
│   ├── api/                   # 预留 API 模块（当前未注册业务路由）
│   ├── core/                  # 基础设施层
│   │   ├── app/               # 应用上下文与生命周期
│   │   ├── cache/             # 缓存实现（memory/redis）
│   │   ├── config/            # 配置加载与监听
│   │   ├── db/                # GORM 初始化
│   │   ├── logger/            # Zap 与 GORM 日志适配
│   │   ├── middleware/        # CORS/JWT/限流
│   │   ├── model/             # 基础模型
│   │   ├── scheduler/         # 定时调度器
│   │   ├── server/            # Gin 引擎与核心路由
│   │   └── upload/            # 上传驱动
│   ├── job/                   # 定时任务注册与任务实现
│   ├── migrate/               # AutoMigrate 与初始化数据
│   └── pkg/                   # 通用包（crud/response/jwt/password/excel/...）
├── web/                       # 前端工程（Vue 3 + Vite）
├── docs/                      # 文档
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```

## 分层职责

### 1) core 层

负责基础设施创建，不承载具体业务逻辑：

- 配置、日志、数据库、缓存、上传、调度器、HTTP 服务。
- 通过 `AppContext` 向模块暴露基础设施对象。

### 2) 业务模块层

- `admin`：已实现，负责后台业务接口。
- `api`：预留模块，当前 `Register()` 为空。
- `job`：注册并运行定时任务。

### 3) pkg 层

无状态工具与可复用组件，例如：

- `response` 统一响应
- `crud` 声明式 CRUD 框架
- `jwt`、`password`、`pagination`、`excel`

## 启动装配流程

1. `BuildContext` 创建基础设施对象。
2. `RegisterCoreRoutes` 注册 `/health`、`/swagger`、静态资源。
3. `RegisterModules` 依次注册 `admin/api/job`。
4. `Run` 启动 HTTP 服务与调度器，并处理优雅退出。

## 数据模型约定

- 基础字段统一复用 `internal/core/model/base.go`。
- GORM 命名策略启用 `SingularTable=true`，默认表名为单数形式。

## 命令行

```bash
bico-admin serve
bico-admin migrate
```

全局参数：

```bash
-c, --config string
```

## 扩展建议

1. 新增后台资源优先走 `internal/pkg/crud`。
2. 在 `internal/admin/module.go` 注册新 CRUD 模块。
3. 新模型加入 `internal/migrate/migrate.go`。
4. 前端页面路由放在 `web/src/router/modules/`。
