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
│   ├── core/                # 核心基础设施层（有状态、需配置）
│   │   ├── app/            # 应用生命周期管理
│   │   │   ├── app.go      # App 实例，处理启动和优雅关闭
│   │   │   ├── context.go   # AppContext + Module 接口 + BuildContext
│   │   │   └── container.go # 旧 DI 容器入口（已废弃，保留避免误用）
│   │   ├── config/         # 配置管理
│   │   │   └── config.go   # 配置结构体 + Viper 加载
│   │   ├── db/             # 数据库层
│   │   │   └── database.go # GORM 初始化和连接池配置
│   │   ├── cache/          # 缓存层
│   │   │   ├── cache.go    # 缓存接口
│   │   │   ├── factory.go  # 缓存工厂
│   │   │   ├── memory.go   # 内存缓存实现
│   │   │   └── redis.go    # Redis缓存实现
│   │   ├── logger/         # 日志系统
│   │   │   ├── logger.go   # Zap日志封装
│   │   │   └── gorm_logger.go # GORM日志适配器
│   │   ├── model/          # 基础模型
│   │   │   └── base.go     # BaseModel（ID, CreatedAt, UpdatedAt）
│   │   ├── server/         # 服务器层
│   │   │   └── server.go   # Gin 引擎创建 + 框架级路由注册
│   │   ├── middleware/     # 通用中间件
│   │   │   ├── jwt.go      # JWT认证中间件
│   │   │   └── cors.go     # 跨域中间件
│   │   ├── scheduler/      # 定时调度器（框架能力）
│   │   │   └── scheduler.go # 基于 robfig/cron 的调度器封装
│   │   └── upload/         # 文件上传
│   │       ├── upload.go   # 上传接口
│   │       ├── factory.go  # 上传工厂
│   │       └── local.go    # 本地存储实现
│   │
│   ├── pkg/                 # 工具包（无状态、零依赖）
│   │   ├── response/       # 统一响应
│   │   │   └── response.go # Success/Error 响应结构
│   │   ├── jwt/            # JWT令牌管理
│   │   │   ├── jwt.go      # 令牌生成和验证
│   │   │   └── token.go    # token 工具
│   │   ├── password/       # 密码加密
│   │   │   └── password.go # bcrypt 密码处理
│   │   └── pagination/     # 分页工具
│   │       └── pagination.go
│   │
│   ├── admin/               # 后台管理模块
│   │   ├── handler/        # HTTP 处理器（声明式 CRUD）
│   │   │   ├── permissions.go       # 基础权限常量
│   │   │   ├── auth_handler.go      # 认证处理器
│   │   │   ├── common_handler.go    # 通用处理器
│   │   │   ├── admin_user_handler.go # 用户管理（含权限/路由/业务）
│   │   │   └── admin_role_handler.go # 角色管理（含权限/路由/业务）
│   │   ├── service/        # 核心服务（仅保留复杂业务）
│   │   │   ├── auth_service.go       # 认证服务
│   │   │   └── config_service.go     # 配置服务
│   │   ├── middleware/     # 业务中间件
│   │   │   ├── permission.go    # 权限验证中间件
│   │   │   └── user_status.go   # 用户状态检查中间件
│   │   ├── model/          # 模块专属模型
│   │   │   ├── admin_user.go   # 后台用户模型
│   │   │   ├── admin_role.go   # 后台角色模型
│   │   │   └── menu.go         # 菜单模型
│   │   ├── router.go       # 路由注册（接收模块显式提供的 CRUD modules）
│   │   └── module.go       # 模块入口：模块内 DI 装配 + 注册路由
│   │
│   ├── api/                 # 前台 API 模块
│   │   ├── handler/        # HTTP 处理器
│   │   ├── service/        # 业务逻辑
│   │   ├── model/          # 模块专属模型
│   │   ├── router.go       # 路由注册
│   │   └── module.go       # 模块入口：模块内装配 + 注册路由
│   │
│   ├── job/                 # 定时任务模块
│   │   ├── register.go      # 任务注册器
│   │   ├── module.go        # 模块入口：注册任务到 ctx.Scheduler
│   │   └── task/            # 任务实现
│   │       ├── clean.go     # 清理任务
│   │       └── sync.go      # 同步任务
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
│   ├── structure.md         # 本文档（项目结构）
│   ├── logger.md            # 日志功能文档
│   ├── job.md               # 定时任务文档
│   ├── auth-api.md          # 认证API文档
│   ├── cache.md             # 缓存模块文档
│   └── improvements.md      # 优化建议
│
├── go.mod                    # Go 模块定义
├── go.sum                    # 依赖校验
├── Makefile                  # 构建脚本
└── README.md                 # 项目说明
```

## 分层架构

### 1. 核心层 (core)

**职责：** 提供基础设施和框架能力（有状态、需配置、单例）

- **app**: AppContext 构建、模块注册、应用生命周期（启动、关闭）
- **config**: 配置文件加载和解析
- **db**: 数据库连接和连接池管理
- **cache**: 缓存驱动（Memory/Redis）
- **logger**: 日志系统（Zap + GORM 日志适配）
- **model**: 基础模型（BaseModel）
- **server**: HTTP 服务器初始化、框架级路由注册
- **middleware**: 通用中间件（JWT、CORS）
- **scheduler**: 定时调度器
- **upload**: 文件上传驱动

### 2. 工具层 (pkg)

**职责：** 提供无状态的工具函数（零依赖、可复用）

- **response**: 统一的 API 响应格式
- **jwt**: JWT token 生成和解析
- **password**: 密码加密（bcrypt）
- **pagination**: 分页工具

### 3. 业务模块 (admin / api)

**职责：** 实现具体业务功能

每个模块采用经典三层结构，并通过 `module.go` 作为模块入口：

- **handler**: 处理 HTTP 请求，参数验证，调用 service
- **service**: 业务逻辑层，处理复杂业务规则
- **model**: 模块专属数据模型
- **router.go**: 路由注册
- **module.go**: 模块入口（模块内 DI 装配 + 注册路由/任务）

### 4. 任务层 (job)

**职责：** 注册定时任务（调度器属于 core）

- **register.go**: 任务注册器，统一创建任务并注册到调度器
- **module.go**: 模块入口，将任务注册到 `ctx.Scheduler`
- **task/**: 具体任务实现
  - `clean.go`: 清理过期数据（每天凌晨 3 点）
  - `sync.go`: 同步统计数据（每小时）

**特点**：
- 支持 6 位 cron 表达式（秒 分 时 日 月 周）
- 集成 Zap 日志，记录任务执行状态
- 自动随应用启动/关闭
- 支持依赖注入（DB、Cache、Logger）

### 5. 迁移层 (migrate)

**职责：** 数据库表结构管理

- 统一注册所有模型的 `AutoMigrate`
- 支持通过命令行一键迁移

## 启动与模块装配流程

### 启动流程

```go
ctx, _ := app.BuildContext(configPath)

server.RegisterCoreRoutes(ctx.Engine, ctx.Cfg, web.DistFS)

_ = app.RegisterModules(
    ctx,
    admin.NewModule(),
    api.NewModule(),
    job.NewModule(),
)

_ = app.Run(ctx)
```

### 依赖注入约定

- core 只负责创建基础设施并放入 `AppContext`
- 业务模块在 `module.go` 内自行装配依赖（可使用 dig）
- 业务路由在模块 `Register()` 内直接注册到 `ctx.Engine`

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

### BaseModel (core/model/base.go)

所有模型继承 BaseModel：

```go
type BaseModel struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 模型使用

```go
import "bico-admin/internal/core/model"

type AdminUser struct {
    model.BaseModel
    Username string `json:"username"`
    // ...
}
```

### 模型分类

- **基础模型** (core/model): BaseModel
  - 所有业务模型的基础，提供公共字段
  
- **模块模型** (module/model): AdminUser, AdminRole, Menu
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

### 添加新模块（推荐：声明式 CRUD）

使用 `internal/pkg/crud` 包，只需创建一个 Handler 文件：

1. 在 `internal/admin/handler/` 创建 `xxx_handler.go`
2. 嵌入 `crud.BaseHandler`，实现 `ModuleConfig()` 方法
3. 在 `internal/admin/module.go` 中把模块构造加入 `[]crud.Module`
4. 在 `migrate.go` 中添加模型迁移

详见 [CRUD 包文档](./crud-pkg.md)

### 添加新模块（传统方式）

1. 在 `internal/` 下创建模块目录
2. 按照 `handler/service/model` 组织代码
3. 创建 `router.go` 实现路由注册
4. 创建 `module.go` 并在 `Register()` 中装配依赖、注册路由
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
