# 项目结构说明

## 目录树

```
bico-admin/
├── cmd/                      # 应用入口
│   └── main.go              # 主程序，定义 cobra 命令
│
├── config/                   # 配置文件
│   └── config.yaml          # 应用配置（服务器、数据库、日志）
│
├── internal/                 # 内部代码（不对外暴露）
│   ├── core/                # 核心层
│   │   ├── app/            # 应用生命周期管理
│   │   │   ├── app.go      # App 实例，处理启动和优雅关闭
│   │   │   └── container.go # DI 容器构建
│   │   ├── config/         # 配置管理
│   │   │   └── config.go   # 配置结构体 + Viper 加载
│   │   ├── db/             # 数据库层
│   │   │   └── database.go # GORM 初始化和连接池配置
│   │   ├── cache/          # 缓存层
│   │   │   ├── cache.go    # 缓存接口
│   │   │   ├── memory.go   # 内存缓存实现
│   │   │   └── redis.go    # Redis缓存实现
│   │   ├── server/         # 服务器层
│   │   │   └── server.go   # Gin 引擎创建 + 路由注册
│   │   ├── middleware/     # 核心中间件
│   │   │   └── jwt.go      # JWT认证中间件
│   │   ├── upload/         # 文件上传
│   │   │   └── uploader.go # 上传器（支持本地/七牛云等）
│   │   └── scheduler/      # 定时任务调度
│   │
│   ├── shared/              # 共享层（跨模块复用）
│   │   ├── model/          # 公共数据模型
│   │   │   ├── base.go     # BaseModel（ID, CreatedAt, UpdatedAt）
│   │   │   ├── admin_user.go  # 管理员用户模型
│   │   │   └── admin_role.go  # 管理员角色模型
│   │   ├── response/       # 统一响应
│   │   │   └── response.go # Success/Error 响应结构
│   │   ├── jwt/            # JWT令牌管理
│   │   │   └── jwt.go      # 令牌生成和验证
│   │   ├── password/       # 密码加密
│   │   │   └── password.go # bcrypt 密码处理
│   │   ├── pagination/     # 分页工具
│   │   │   └── pagination.go
│   │   ├── logger/         # 日志工具
│   │   └── util/           # 工具函数
│   │
│   ├── admin/               # 后台管理模块
│   │   ├── consts/         # 常量定义
│   │   │   └── permissions.go # 权限树定义
│   │   ├── handler/        # HTTP 处理器
│   │   │   ├── auth_handler.go      # 认证处理器
│   │   │   ├── common_handler.go    # 通用处理器
│   │   │   ├── admin_user_handler.go # 用户管理
│   │   │   └── admin_role_handler.go # 角色管理
│   │   ├── service/        # 业务逻辑
│   │   │   ├── auth_service.go       # 认证服务
│   │   │   ├── config_service.go     # 配置服务
│   │   │   ├── admin_user_service.go # 用户服务
│   │   │   └── admin_role_service.go # 角色服务
│   │   ├── middleware/     # 模块中间件
│   │   │   ├── permission.go    # 权限验证中间件
│   │   │   └── user_status.go   # 用户状态检查中间件
│   │   ├── model/          # 模块专属模型
│   │   └── router.go       # 路由注册
│   │
│   ├── api/                 # 前台 API 模块
│   │   ├── handler/        # HTTP 处理器
│   │   ├── service/        # 业务逻辑
│   │   ├── model/          # 模块专属模型
│   │   └── router.go       # 路由注册
│   │
│   ├── job/                 # 定时任务
│   │   ├── task/           # 任务实现
│   │   └── register.go     # 任务注册器
│   │
│   └── migrate/             # 数据库迁移
│       └── migrate.go       # 统一管理所有模型的 AutoMigrate
│
├── web/                      # 前端项目
│   ├── config/              # 配置文件
│   │   ├── routes.ts        # 路由配置
│   │   ├── config.ts        # 构建配置
│   │   └── proxy.ts         # 代理配置
│   ├── src/                 # 源代码
│   │   ├── app.tsx          # 应用入口
│   │   ├── access.ts        # 权限控制
│   │   ├── components/      # 公共组件
│   │   ├── pages/          # 页面组件
│   │   ├── services/       # API服务
│   │   └── locales/        # 国际化
│   └── package.json
│
├── docs/                     # 项目文档
│   ├── structure.md         # 本文档
│   ├── auth-api.md          # 认证API文档
│   ├── cache.md             # 缓存模块文档
│   └── AGENT.md             # AI助手项目说明
│
├── go.mod                    # Go 模块定义
├── go.sum                    # 依赖校验
├── Makefile                  # 构建脚本
└── README.md                 # 项目说明
```

## 分层架构

### 1. 核心层 (core)

**职责：** 提供基础设施和框架能力

- **app**: DI 容器管理、应用生命周期（启动、关闭）
- **config**: 配置文件加载和解析
- **db**: 数据库连接和连接池管理
- **server**: HTTP 服务器初始化、中间件、路由注册
- **scheduler**: 定时任务调度器

### 2. 共享层 (shared)

**职责：** 提供跨模块复用的通用能力

- **model**: 公共数据模型（User、Role、Log 等）
- **response**: 统一的 API 响应格式
- **logger**: 日志组件
- **util**: 通用工具函数

### 3. 业务模块 (admin / api)

**职责：** 实现具体业务功能

每个模块采用经典三层结构：

- **handler**: 处理 HTTP 请求，参数验证，调用 service
- **service**: 业务逻辑层，处理复杂业务规则
- **model**: 模块专属数据模型
- **router.go**: 路由注册，实现 `server.Router` 接口

### 4. 任务层 (job)

**职责：** 定时任务和后台作业

- **task**: 具体任务实现（清理、同步等）
- **register.go**: 统一注册所有定时任务

### 5. 迁移层 (migrate)

**职责：** 数据库表结构管理

- 统一注册所有模型的 `AutoMigrate`
- 支持通过命令行一键迁移

## 依赖注入流程

### 容器构建 (container.go)

```go
BuildContainer(configPath) -> dig.Container
  ├── Provide Config (从配置文件加载)
  ├── Provide *gorm.DB (数据库连接)
  ├── Provide *gin.Engine (HTTP 引擎)
  ├── Provide Routers (各模块路由)
  └── Provide *app.App (应用实例)
```

### 依赖注入优势

- 解耦：模块间通过接口交互
- 可测试：方便 Mock 依赖
- 可维护：依赖关系清晰可见

## 路由注册机制

### Router 接口

```go
type Router interface {
    Register(engine *gin.Engine)
}
```

### 注册流程

1. 各模块实现 `Router` 接口
2. 在 `container.go` 中 Provide 到 DI 容器
3. 在 `server.RegisterRoutes()` 中统一注册
4. 支持路由分组（`/admin`, `/api`）

## 配置管理

### 配置文件 (config/config.yaml)

```yaml
server:
  port: 8080
  mode: debug

database:
  driver: mysql
  host: localhost
  # ...

log:
  level: info
```

### 配置结构体 (config/config.go)

- 使用 Viper 加载 YAML 配置
- 通过 `mapstructure` tag 映射字段
- 支持多环境配置（dev、prod）

## 数据模型设计

### BaseModel (shared/model/base.go)

所有模型继承 BaseModel：

```go
type BaseModel struct {
    ID        uint
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 模型分类

- **共享模型** (shared/model): User, Role, Log
  - 跨模块使用的通用模型
  
- **模块模型** (module/model): Menu, Order, Product
  - 模块专属，仅在模块内使用

## 命令行工具

### 可用命令

```bash
bico-admin serve    # 启动 HTTP 服务
bico-admin migrate  # 执行数据库迁移
```

### 全局参数

```bash
-c, --config string   # 指定配置文件路径
```

## 扩展指南

### 添加新模块

1. 在 `internal/` 下创建模块目录
2. 按照 `handler/service/model` 组织代码
3. 创建 `router.go` 实现路由注册
4. 在 `container.go` 中注册到 DI 容器
5. 在 `migrate.go` 中添加模型迁移

### 添加新的定时任务

1. 在 `job/task/` 下创建任务文件
2. 实现任务逻辑
3. 在 `job/register.go` 中注册任务

### 添加新的配置项

1. 在 `config/config.yaml` 中添加配置
2. 在 `config/config.go` 中添加结构体字段
3. 使用 `mapstructure` tag 映射

## 最佳实践

### 编码规范

- 遵循 SOLID、DRY、KISS、YAGNI 原则
- 使用依赖注入，避免全局变量
- 错误处理：向上传递，统一处理
- 避免在循环中执行 I/O 操作

### 项目约定

- 模型命名：单数形式（User 而非 Users）
- 表名命名：复数形式（users 而非 user）
- 包名：小写单词，不使用下划线
- 接口命名：通常以 `-er` 结尾

### Git 管理

- `.gitkeep` 用于占位空目录
- `bin/` 目录不提交到版本库
- 敏感配置使用环境变量
