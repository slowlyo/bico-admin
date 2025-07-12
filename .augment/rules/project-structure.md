---
type: "always_apply"
---

# 项目结构设计规范

## 项目概述

Bico Admin 是基于Go语言的多端管理系统，采用清晰的分层架构和模块化设计，包含master、admin、api三个端点，以admin端为主要开发重点。

## 核心设计原则

### 1. 多端分离架构
- **master端(主控)**: 系统主控制台，负责全局配置和监控
- **admin端(管理)**: 管理后台，主要开发重点，包含完整的RBAC权限系统
- **api端(接口)**: 对外API服务，提供标准化接口

### 2. 分层架构
```
Handler → Service → Repository → Database
```
- **Handler**: 处理HTTP请求，参数验证，响应格式化
- **Service**: 业务逻辑处理，事务管理
- **Repository**: 数据访问层，数据库操作
- **Database**: 数据存储层

### 3. 依赖注入
- 使用Google Wire进行编译时依赖注入
- 每个模块都有独立的Provider
- 清晰的依赖关系管理

## 目录结构规范

### 根目录结构
```
bico-admin/
├── cmd/                    # 应用程序入口点
├── internal/              # 私有应用代码
├── pkg/                   # 可复用的库代码
├── configs/               # 配置文件
├── docs/                  # 项目文档
├── scripts/               # 构建和部署脚本
├── web/                   # 前端资源（如果有）
├── .augment/              # Augment配置和规则
├── go.mod                 # Go模块定义
├── go.sum                 # Go模块校验和
├── Makefile              # 构建脚本
└── README.md             # 项目说明
```

### cmd/ 目录结构
```
cmd/
├── server/                # Web服务入口
│   ├── main.go           # 主程序入口
│   ├── wire.go           # Wire依赖注入配置
│   └── wire_gen.go       # Wire生成的代码
└── migrate/              # 数据库迁移工具
    └── main.go           # 迁移程序入口
```

### internal/ 目录结构
```
internal/
├── master/               # 主控端模块
├── admin/                # 管理端模块（主要开发重点）
├── api/                  # API端模块
└── shared/               # 共享组件
```

### 端点模块结构（以admin为例）
```
internal/admin/
├── handler/              # 处理器层
│   ├── admin_user.go     # 管理员用户处理器
│   ├── admin_role.go     # 管理员角色处理器
│   ├── system.go         # 系统管理处理器
│   └── auth.go           # 认证处理器
├── service/              # 服务层
│   ├── admin_user.go     # 管理员用户服务
│   ├── admin_role.go     # 管理员角色服务
│   ├── system.go         # 系统管理服务
│   └── auth.go           # 认证服务
├── repository/           # 数据访问层
│   ├── admin_user.go     # 管理员用户数据访问
│   ├── admin_role.go     # 管理员角色数据访问
│   └── system_log.go     # 系统日志数据访问
├── models/               # 数据模型（特定于admin端）
│   ├── admin_user.go     # 管理员用户模型
│   ├── admin_role.go     # 管理员角色模型
│   └── system_log.go     # 系统日志模型
├── types/                # 类型定义
│   ├── auth.go           # 认证相关类型
│   ├── user.go           # 用户管理相关类型
│   ├── role.go           # 角色管理相关类型
│   ├── system.go         # 系统管理相关类型
│   └── common.go         # 通用类型定义
├── routes/               # 路由注册
│   └── routes.go         # 路由定义和注册
├── middleware/           # 中间件
│   ├── auth.go           # 认证中间件
│   ├── permission.go     # 权限中间件
│   ├── cors.go           # 跨域中间件
│   └── logger.go         # 日志中间件
├── definitions/          # 静态定义（admin端特有）
│   ├── permissions.go    # 权限定义
│   └── menus.go          # 菜单定义
└── provider.go           # Wire Provider配置
```

### shared/ 目录结构
```
internal/shared/
├── model/                # 共享数据模型
│   ├── user.go           # 普通用户模型
│   ├── admin_user.go     # 管理员用户模型（简化版）
│   └── base.go           # 基础模型
├── types/                # 共享类型定义
│   ├── common.go         # 通用类型
│   ├── response.go       # 响应类型
│   └── pagination.go    # 分页类型
├── config/               # 配置管理
│   ├── config.go         # 配置结构定义
│   └── loader.go         # 配置加载器
├── database/             # 数据库相关
│   ├── connection.go     # 数据库连接
│   ├── migration.go      # 迁移管理
│   └── redis.go          # Redis连接
├── middleware/           # 共享中间件
│   ├── recovery.go       # 恢复中间件
│   ├── timeout.go        # 超时中间件
│   └── rate_limit.go     # 限流中间件
├── utils/                # 工具函数
│   ├── crypto.go         # 加密工具
│   ├── validator.go      # 验证工具
│   ├── time.go           # 时间工具
│   └── string.go         # 字符串工具
├── response/             # 响应处理
│   ├── response.go       # 统一响应格式
│   └── error.go          # 错误处理
└── logger/               # 日志管理
    ├── logger.go         # 日志配置
    └── middleware.go     # 日志中间件
```

## 文件命名规范

### 1. Go文件命名
- 使用小写字母和下划线：`admin_user.go`, `system_config.go`
- 避免使用驼峰命名：~~`adminUser.go`~~
- 文件名应该反映其主要功能

### 2. 包命名
- 使用小写字母，避免下划线：`handler`, `service`, `repository`
- 包名应该简洁且有意义
- 避免使用复数形式：`handler` 而不是 `handlers`

### 3. 结构体命名
- 使用大驼峰命名：`AdminUser`, `SystemConfig`
- 接口名通常以 `er` 结尾：`UserService`, `ConfigLoader`

## 模块间依赖规则

### 1. 依赖方向
```
Handler → Service → Repository → Model
```
- 上层可以依赖下层，下层不能依赖上层
- 同层之间避免循环依赖

### 2. 跨端点依赖
- 各端点模块相互独立，避免直接依赖
- 共享功能放在 `shared` 目录
- 通过接口定义进行解耦

### 3. 外部依赖
- 第三方库的使用应该封装在 `shared` 或具体模块内
- 避免在多个地方直接使用同一个外部库

## 数据模型设计规范

### 1. 模型分类
- **共享模型**: 放在 `internal/shared/model/`，多个端点都会使用
- **特定模型**: 放在各端点的 `models/` 目录，只在该端点使用

### 2. 表命名规范
- 使用小写字母和下划线：`admin_users`, `admin_roles`
- 管理员相关表加 `admin_` 前缀：`admin_users`, `admin_roles`
- 关联表使用两个表名组合：`admin_user_roles`

### 3. 字段命名规范
- 使用小写字母和下划线：`user_id`, `created_at`
- 外键字段以 `_id` 结尾：`user_id`, `role_id`
- 时间字段使用标准命名：`created_at`, `updated_at`, `deleted_at`

## API设计规范

### 1. 路由规范
- 使用RESTful风格：`GET /api/admin/users`, `POST /api/admin/users`
- 端点前缀：`/api/master/`, `/api/admin/`, `/api/api/`
- 版本控制：`/api/v1/admin/users`（如需要）

### 2. 响应格式
```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

### 3. 错误处理
- 使用标准HTTP状态码
- 提供详细的错误信息
- 统一的错误响应格式

## 配置管理规范

### 1. 配置文件结构
```
configs/
├── config.yaml           # 主配置文件
├── config.dev.yaml       # 开发环境配置
├── config.prod.yaml      # 生产环境配置
└── config.test.yaml      # 测试环境配置
```

### 2. 配置优先级
1. 环境变量
2. 命令行参数
3. 配置文件
4. 默认值

## 测试规范

### 1. 测试文件命名
- 测试文件以 `_test.go` 结尾
- 与被测试文件在同一目录：`user.go` → `user_test.go`

### 2. 测试分类
- 单元测试：测试单个函数或方法
- 集成测试：测试模块间交互
- 端到端测试：测试完整的业务流程

## 部署和构建规范

### 1. Makefile规范
- 提供常用的构建、测试、部署命令
- 使用 `.PHONY` 声明伪目标
- 命令应该具有良好的可读性

### 2. Docker支持
- 提供 `Dockerfile` 用于容器化部署
- 使用多阶段构建优化镜像大小
- 提供 `docker-compose.yml` 用于本地开发

## 文档规范

### 1. 代码注释
- 公开的函数、结构体必须有注释
- 注释应该说明功能、参数、返回值
- 使用标准的Go注释格式

### 2. API文档
- 使用Swagger/OpenAPI规范
- 提供详细的接口说明和示例
- 保持文档与代码同步

### 3. 项目文档
- README.md：项目概述和快速开始
- docs/：详细的设计文档和使用指南
- CHANGELOG.md：版本变更记录

## 特殊规范说明

### 1. Admin端特殊目录
- `definitions/`: 存放权限和菜单的静态定义
- `models/`: 存放admin端特有的数据模型（如AdminUser, AdminRole）
- 这些目录不应该出现在其他端点中

### 2. Shared目录使用原则
- 只存放真正通用的组件
- 管理员相关的模型和逻辑不应该放在shared中
- 普通用户模型可以放在shared中供多个端点使用

### 3. 权限系统特殊性
- 权限定义必须在代码中静态定义（definitions/permissions.go）
- 菜单定义必须在代码中静态定义（definitions/menus.go）
- 用户和角色数据在数据库中动态管理
- 权限验证逻辑应该集中在middleware中

