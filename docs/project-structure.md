# Go Web应用项目结构设计

## 目录结构概览

```
go-web-app/
├── cmd/                    # 应用程序入口点
│   ├── server/
│   │   ├── main.go        # Web服务入口
│   │   ├── wire.go        # Wire依赖注入配置
│   │   └── wire_gen.go    # Wire生成的代码
│   └── scheduler/
│       ├── main.go        # 定时任务入口
│       ├── wire.go        # Wire依赖注入配置
│       └── wire_gen.go    # Wire生成的代码
├── internal/              # 私有应用代码
│   ├── admin/             # 超管端模块
│   │   ├── handler/       # 超管端处理器
│   │   │   ├── user.go
│   │   │   ├── merchant.go
│   │   │   └── system.go
│   │   ├── service/       # 超管端业务逻辑
│   │   │   ├── user.go
│   │   │   ├── merchant.go
│   │   │   └── system.go
│   │   ├── repository/    # 超管端数据访问层
│   │   │   ├── system_log.go
│   │   │   └── admin_config.go
│   │   ├── types/         # 超管端类型定义
│   │   │   ├── user.go    # 用户管理相关类型
│   │   │   ├── merchant.go # 商户管理相关类型
│   │   │   ├── system.go  # 系统管理相关类型
│   │   │   └── common.go  # 通用类型定义
│   │   ├── routes/        # 超管端路由注册
│   │   │   └── routes.go
│   │   ├── middleware/    # 超管端中间件
│   │   │   ├── auth.go
│   │   │   └── permission.go
│   │   └── provider.go    # 超管端Wire Provider
│   ├── merchant/          # 商户端模块
│   │   ├── handler/       # 商户端处理器
│   │   ├── service/       # 商户端业务逻辑
│   │   ├── repository/    # 商户端数据访问层
│   │   ├── types/         # 商户端类型定义
│   │   ├── routes/        # 商户端路由注册
│   │   │   └── routes.go
│   │   ├── middleware/    # 商户端中间件
│   │   └── provider.go    # 商户端Wire Provider
│   ├── user/              # 用户端模块
│   │   ├── handler/       # 用户端处理器
│   │   ├── service/       # 用户端业务逻辑
│   │   ├── repository/    # 用户端数据访问层
│   │   ├── types/         # 用户端类型定义
│   │   ├── routes/        # 用户端路由注册
│   │   │   └── routes.go
│   │   ├── middleware/    # 用户端中间件
│   │   └── provider.go    # 用户端Wire Provider
│   ├── scheduler/         # 定时任务模块
│   │   ├── job/           # 任务定义
│   │   │   ├── user_cleanup.go
│   │   │   ├── order_sync.go
│   │   │   └── data_backup.go
│   │   ├── service/       # 任务业务逻辑
│   │   │   ├── cleanup.go
│   │   │   ├── sync.go
│   │   │   └── backup.go
│   │   ├── cron/          # 定时调度器
│   │   │   ├── scheduler.go
│   │   │   └── manager.go
│   │   └── provider.go    # 定时任务Wire Provider
│   ├── shared/            # 共享模块
│   │   ├── repository/    # 数据访问层
│   │   │   ├── user.go
│   │   │   ├── merchant.go
│   │   │   ├── product.go
│   │   │   ├── order.go
│   │   │   └── interfaces.go
│   │   ├── model/         # 数据模型
│   │   │   ├── user.go
│   │   │   ├── merchant.go
│   │   │   ├── product.go
│   │   │   └── order.go
│   │   ├── types/         # 通用类型定义
│   │   │   ├── pagination.go # 分页相关类型
│   │   │   └── common.go     # 其他通用类型
│   │   ├── middleware/    # 通用中间件
│   │   │   ├── cors.go
│   │   │   ├── logging.go
│   │   │   └── recovery.go
│   │   └── provider.go    # 共享组件Wire Provider
├── config/                # 配置文件目录
│   ├── app.yml            # 基础配置模板
│   ├── app.dev.yml        # 开发环境配置
│   ├── app.test.yml       # 测试环境配置
│   ├── app.staging.yml    # 预发布环境配置
│   └── app.prod.yml       # 生产环境配置
├── pkg/                   # 可被外部应用使用的库代码
│   ├── config/            # 配置读取逻辑
│   │   ├── config.go      # 配置结构定义
│   │   └── loader.go      # 配置加载器
│   ├── logger/            # 日志工具
│   │   └── logger.go
│   ├── database/          # 数据库连接
│   │   ├── postgres.go
│   │   └── redis.go
│   ├── cache/             # 缓存封装
│   │   ├── redis.go       # Redis缓存实现
│   │   ├── memory.go      # 内存缓存实现
│   │   └── interface.go   # 缓存接口定义
│   ├── response/          # 统一响应封装
│   │   ├── response.go    # 响应结构定义
│   │   └── helper.go      # 响应辅助函数
│   ├── validator/         # 验证工具
│   │   └── validator.go
│   └── utils/             # 通用工具
│       └── crypto.go
├── api/                   # API定义文件
│   └── openapi/           # OpenAPI/Swagger规范
│       └── api.yaml
├── web/                   # Web资源
│   ├── static/            # 静态文件
│   │   ├── css/
│   │   ├── js/
│   │   └── images/
│   └── templates/         # 模板文件
│       ├── layout/
│       └── pages/
├── scripts/               # 构建和部署脚本
│   ├── build.sh
│   ├── deploy.sh
│   └── migrate.sh
├── deployments/           # 部署配置
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── k8s/
│       ├── deployment.yaml
│       └── service.yaml
├── tests/                 # 测试文件
│   ├── integration/
│   └── fixtures/
├── docs/                  # 项目文档
│   ├── api.md
│   └── deployment.md
├── .env.example           # 环境变量示例
├── .gitignore
├── go.mod
├── go.sum
├── Makefile               # 构建命令
├── Dockerfile
└── README.md
```

## 目录说明

### `/cmd`
- **用途**: 应用程序的主要入口点
- **原则**: 每个应用一个子目录，保持main函数简洁
- **包含**:
  - `/cmd/server` - Web服务入口，包含Wire配置
  - `/cmd/scheduler` - 定时任务服务入口，包含Wire配置
- **Wire文件**:
  - `wire.go` - 依赖注入配置，定义Provider和Injector
  - `wire_gen.go` - Wire自动生成的代码，不可手动编辑

### `/internal`
- **用途**: 私有应用和库代码，不能被其他应用导入
- **优势**: 强制代码封装，防止外部依赖
- **架构**: 按业务端分离，共享通用组件

#### `/internal/admin` - 超管端模块
- **职责**: 系统管理、用户管理、商户管理
- **权限**: 最高权限，可管理所有资源
- **特点**: 功能全面，安全性要求最高
- **repository目录**: 超管端专用表（如系统日志、管理员配置等）
- **types目录**: 按业务模块组织类型定义，每个文件包含相关的请求、响应结构
- **routes目录**: 超管端路由注册，包含所有管理接口路由
- **provider.go**: 超管端Wire Provider，定义该模块的依赖注入

#### `/internal/merchant` - 商户端模块
- **职责**: 商品管理、订单处理、店铺运营
- **权限**: 仅能操作自己的商户数据
- **特点**: 业务功能丰富，数据隔离
- **repository目录**: 商户端专用表（如商品库存、店铺配置等）
- **types目录**: 按业务功能组织类型定义（如product.go、order.go等）
- **routes目录**: 商户端路由注册，包含商户业务接口路由
- **provider.go**: 商户端Wire Provider，定义该模块的依赖注入

#### `/internal/user` - 用户端模块
- **职责**: 用户注册登录、浏览商品、下单购买
- **权限**: 仅能操作自己的用户数据
- **特点**: 面向C端，性能要求高
- **repository目录**: 用户端专用表（如用户偏好、浏览记录等）
- **types目录**: 按业务功能组织类型定义（如auth.go、profile.go等）
- **routes目录**: 用户端路由注册，包含用户业务接口路由
- **provider.go**: 用户端Wire Provider，定义该模块的依赖注入

#### `/internal/scheduler` - 定时任务模块
- **职责**: 定时任务的定义、调度和执行
- **特点**: 独立运行，可单独部署
- **包含**:
  - `job/` - 具体任务实现（数据清理、同步、备份等）
  - `service/` - 任务业务逻辑
  - `cron/` - 任务调度管理

#### `/internal/shared` - 共享模块
- **职责**: 各端共用的组件和服务
- **包含**: 数据模型、仓储层、通用中间件
- **原则**: 高内聚低耦合，便于复用

##### `/internal/shared/repository`
- **职责**: 共享数据访问抽象层
- **原则**: 接口定义，便于mock和测试
- **实现**: 基于ORM框架的数据库操作逻辑
- **内容**: 各端都需要的核心表（如用户表、商户表、订单表等）

##### `/internal/shared/model`
- **职责**: 共享数据结构定义
- **包含**: 数据库实体模型、ORM结构体

##### `/internal/shared/types`
- **职责**: 业务通用类型定义
- **包含**:
  - `pagination.go` - 分页查询结构（BasePageQuery等）
  - `common.go` - 业务通用类型（状态枚举、常量等）

##### `/internal/shared/middleware`
- **职责**: 通用HTTP中间件实现
- **功能**: CORS、日志、恢复、限流等

##### `/internal/shared/provider.go`
- **职责**: 共享组件Wire Provider
- **功能**: 提供数据库、缓存、日志等基础设施的Provider
- **内容**: 数据库连接、Redis连接、Repository等基础组件的Provider函数

### `/config`
- **用途**: 配置文件存储，按环境分离
- **原则**: 环境变量优先，配置文件补充
- **组织方式**:
  - `app.yml` - 基础配置模板，包含所有配置项
  - `app.{env}.yml` - 各环境特定配置，覆盖基础配置
- **支持**: 多环境配置管理，便于部署

### `/pkg`
- **用途**: 可以被外部应用使用的库代码
- **原则**: 通用性强，独立性好

#### `/pkg/config`
- **职责**: 配置读取和解析逻辑
- **功能**: 加载yml文件、环境变量覆盖、配置验证
- **特点**: 与具体配置文件解耦，便于测试

#### `/pkg/cache`
- **职责**: 缓存抽象层实现
- **功能**: 支持Redis、内存缓存，统一接口
- **特点**: 可插拔设计，便于切换缓存实现

#### `/pkg/response`
- **职责**: 统一响应格式封装
- **功能**: 标准化API响应结构，错误码管理
- **特点**: 统一响应格式，便于前端处理
- **内容**: ApiResponse、ErrorResponse、Success/Error辅助函数

#### `/pkg/validator`
- **职责**: 数据验证工具封装
- **功能**: 请求参数验证、业务规则校验
- **特点**: 统一验证逻辑，支持自定义验证规则

#### `/pkg/utils`
- **职责**: 通用工具函数库
- **功能**: 加密解密、字符串处理、时间工具等
- **特点**: 无状态工具函数，便于复用

### `/api`
- **用途**: API契约定义
- **格式**: OpenAPI规范、Swagger文档

### `/web`
- **用途**: Web应用特定的组件
- **内容**: 静态文件、模板、前端资源

### `/scripts`
- **用途**: 执行各种构建、安装、分析等操作的脚本

### `/deployments`
- **用途**: 部署配置和模板
- **内容**: Docker、Kubernetes、CI/CD配置

## 架构设计原则

### 1. 多端分离架构
```
┌─────────────┬─────────────┬─────────────┐
│   超管端     │   商户端    │   用户端    │
│   /admin    │ /merchant   │   /user     │
└─────────────┴─────────────┴─────────────┘
              │
              ▼
        共享业务层 /shared
              │
              ▼
           数据库层
```

### 2. 路由分层架构
```
HTTP请求 → Routes → Middleware → Handler → Service → Repository → Database
```

### 3. 分层架构
```
HTTP请求 → Handler → Service → Repository → Database
```

### 3. 依赖注入 (Wire)
- 使用Google Wire进行编译时依赖注入
- 自动生成依赖关系代码，避免运行时反射
- 便于单元测试和mock
- 提高代码可维护性和性能

### 4. 统一响应格式
- 标准化API响应结构
- 统一错误码和消息
- 便于前端统一处理

### 5. 缓存策略
- 多级缓存支持（Redis + 内存）
- 统一缓存接口，便于切换实现
- 支持缓存预热和失效策略

### 6. 数据隔离
- 超管端：全局数据访问权限
- 商户端：基于商户ID的数据隔离
- 用户端：基于用户ID的数据隔离

### 7. 配置外部化
- 环境变量优先
- 配置文件作为补充
- 支持多环境部署

### 8. 错误处理
- 统一错误处理机制
- 结构化错误信息
- 适当的日志记录

### 9. 安全考虑
- 输入验证
- 认证授权
- 敏感信息保护
- 不同端的权限控制

## 最佳实践

1. **遵循Go惯用法**: 简洁、清晰、高效
2. **接口设计**: 小而专一的接口定义
3. **错误处理**: 显式错误处理，避免panic
4. **并发安全**: 合理使用goroutine和channel
5. **性能优化**: 适当的缓存和连接池
6. **可观测性**: 完善的日志、监控和追踪

## 技术栈

- **Web框架**: Gin
- **ORM框架**: GORM
- **数据库**: MySQL
- **缓存**: Redis
- **定时任务**: Cron
- **依赖注入**: Google Wire
- **配置管理**: Viper
- **日志**: Zap
- **测试**: Testify
- **文档**: Swagger

这个结构设计兼顾了可维护性、可扩展性和Go语言的最佳实践，适合中小型到大型的Web应用开发。