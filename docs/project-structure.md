# Bico Admin 项目结构

## 目录结构概览

```
bico-admin/
├── cmd/                   # 应用程序入口点
│   ├── migrate/           # 数据库迁移工具
│   │   └── main.go
│   └── server/            # Web服务入口
│       ├── main.go        # 主程序入口
│       ├── wire.go        # Wire依赖注入配置
│       └── wire_gen.go    # Wire生成的代码
├── internal/              # 私有应用代码
│   ├── admin/             # 管理端模块
│   │   ├── definitions/   # 权限和菜单定义
│   │   ├── handler/       # 请求处理器
│   │   ├── initializer/   # 初始化器
│   │   ├── middleware/    # 中间件
│   │   ├── models/        # 数据模型
│   │   ├── repository/    # 数据访问层
│   │   ├── routes/        # 路由注册
│   │   ├── service/       # 业务逻辑层
│   │   ├── types/         # 类型定义
│   │   └── provider.go    # Wire Provider
│   ├── api/               # API端模块
│   │   ├── handler/       # 请求处理器
│   │   ├── middleware/    # 中间件
│   │   ├── repository/    # 数据访问层
│   │   ├── routes/        # 路由注册
│   │   ├── service/       # 业务逻辑层
│   │   ├── types/         # 类型定义
│   │   └── provider.go    # Wire Provider
│   ├── master/            # 主控端模块
│   │   ├── handler/       # 请求处理器
│   │   ├── middleware/    # 中间件
│   │   ├── repository/    # 数据访问层
│   │   ├── routes/        # 路由注册
│   │   ├── service/       # 业务逻辑层
│   │   ├── types/         # 类型定义
│   │   └── provider.go    # Wire Provider
│   └── shared/            # 共享模块
│       ├── middleware/    # 通用中间件
│       ├── models/        # 共享数据模型
│       ├── types/         # 通用类型定义
│       ├── cache.go       # 缓存配置
│       └── provider.go    # 共享组件Wire Provider
├── pkg/                   # 可复用的库代码
│   ├── cache/             # 缓存封装
│   │   ├── interface.go   # 缓存接口
│   │   ├── manager.go     # 缓存管理器
│   │   ├── memory.go      # 内存缓存实现
│   │   └── redis.go       # Redis缓存实现
│   ├── config/            # 配置管理
│   │   ├── config.go      # 配置结构定义
│   │   └── loader.go      # 配置加载器
│   ├── database/          # 数据库连接
│   │   ├── mysql.go       # MySQL连接
│   │   ├── postgres.go    # PostgreSQL连接
│   │   ├── redis.go       # Redis连接
│   │   └── sqlite.go      # SQLite连接
│   ├── jwt/               # JWT工具
│   │   └── jwt.go
│   ├── logger/            # 日志工具
│   │   └── logger.go
│   └── response/          # 统一响应封装
│       └── response.go
├── web/                   # 前端应用
│   ├── src/               # 源代码
│   │   ├── components/    # 组件
│   │   ├── config/        # 配置
│   │   ├── models/        # 数据模型
│   │   ├── pages/         # 页面
│   │   ├── services/      # API服务
│   │   ├── utils/         # 工具函数
│   │   ├── access.ts      # 权限配置
│   │   └── app.tsx        # 应用入口
│   ├── dist/              # 构建产物
│   ├── package.json       # 依赖配置
│   └── tsconfig.json      # TypeScript配置
├── config/                # 配置文件
│   ├── app.yml            # 基础配置
│   └── app.dev.yml        # 开发环境配置
├── data/                  # 数据文件
│   └── *.db               # SQLite数据库文件
├── docs/                  # 项目文档
│   ├── database-sqlite.md # 数据库文档
│   ├── file-upload.md     # 文件上传文档
│   ├── project-structure.md # 项目结构文档
│   ├── docs.go            # Swagger配置

├── logs/                  # 日志文件
├── tests/                 # 测试文件
│   ├── config/            # 配置测试
│   └── cache_test.go      # 缓存测试
├── scripts/               # 脚本文件
├── bin/                   # 编译产物
├── go.mod                 # Go模块定义
├── go.sum                 # Go模块校验和
├── Makefile               # 构建命令
└── README.md              # 项目说明
```

## 架构设计

### 多端分离架构
```
┌─────────────┬─────────────┬─────────────┐
│   主控端     │   管理端    │   API端     │
│   /master   │   /admin    │   /api      │
└─────────────┴─────────────┴─────────────┘
              │
              ▼
        共享业务层 /shared
              │
              ▼
           数据库层
```

### 分层架构
```
HTTP请求 → Routes → Middleware → Handler → Service → Repository → Database
```

## 目录说明

### `/cmd` - 应用程序入口点
- **server/**: Web服务入口，包含主程序和Wire配置
- **migrate/**: 数据库迁移工具

### `/internal` - 私有应用代码

#### `/internal/admin` - 管理端模块
- **职责**: 管理员用户管理、角色权限管理、系统配置
- **特点**: 基于RBAC的权限控制，支持多角色管理
- **包含**:
  - `definitions/` - 权限和菜单定义
  - `handler/` - HTTP请求处理器
  - `service/` - 业务逻辑层
  - `repository/` - 数据访问层
  - `middleware/` - 权限中间件
  - `models/` - 数据模型
  - `types/` - 类型定义

#### `/internal/api` - API端模块
- **职责**: 对外API接口，供第三方调用
- **特点**: 轻量级，专注于API服务

#### `/internal/master` - 主控端模块
- **职责**: 系统监控、主控管理
- **特点**: 系统级管理功能

#### `/internal/shared` - 共享模块
- **职责**: 各端共用的组件和服务
- **包含**:
  - `models/` - 共享数据模型
  - `types/` - 通用类型定义
  - `middleware/` - 通用中间件

### `/pkg` - 可复用库代码

#### `/pkg/cache` - 缓存封装
- **功能**: 支持Redis和内存缓存，统一接口
- **特点**: 可插拔设计，便于切换缓存实现

#### `/pkg/config` - 配置管理
- **功能**: 配置文件加载、环境变量覆盖
- **特点**: 支持多环境配置

#### `/pkg/database` - 数据库连接
- **功能**: 支持MySQL、PostgreSQL、SQLite
- **特点**: 统一接口，便于切换数据库

#### `/pkg/jwt` - JWT工具
- **功能**: JWT令牌生成和验证

#### `/pkg/logger` - 日志工具
- **功能**: 基于Zap的结构化日志

#### `/pkg/response` - 统一响应
- **功能**: 标准化API响应格式

### `/web` - 前端应用
- **技术栈**: React + TypeScript + Ant Design + UmiJS
- **特点**: 现代化前端架构，组件化开发

### 其他目录
- `/config/` - 配置文件
- `/data/` - SQLite数据库文件
- `/docs/` - 项目文档
- `/tests/` - 测试文件
- `/scripts/` - 脚本文件

## 开发规范

### 代码组织原则
1. **按端分离**: 不同端点的代码完全分离，避免耦合
2. **分层架构**: Handler → Service → Repository → Database
3. **依赖注入**: 使用Google Wire进行编译时依赖注入
4. **接口抽象**: Repository层使用接口，便于测试和扩展

### 命名规范
- **包名**: 小写，简洁明了
- **文件名**: 小写+下划线，如 `admin_user.go`
- **结构体**: 大驼峰，如 `AdminUser`
- **方法**: 大驼峰（公开）或小驼峰（私有）

### 目录规范
- **handler**: HTTP请求处理，参数验证，调用service
- **service**: 业务逻辑处理，事务管理
- **repository**: 数据访问，数据库操作
- **types**: 请求/响应结构体定义
- **middleware**: 中间件实现