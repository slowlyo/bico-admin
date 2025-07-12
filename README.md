# Bico Admin

基于Go语言的多端管理系统，包含master、admin、api三个端点，以admin端为主要开发重点。

## 项目特点

- 🎯 **多端分离架构**: master端(主控)、admin端(管理)、api端(接口)
- 🏗️ **清晰的分层架构**: Handler → Service → Repository → Database
- 🔧 **依赖注入**: 使用Google Wire进行编译时依赖注入
- 📝 **统一响应格式**: 标准化API响应结构
- 🔐 **权限控制**: 不同端点的权限隔离
- 📊 **完善的日志**: 基于Zap的结构化日志
- 🗄️ **数据库支持**: MySQL/PostgreSQL/SQLite + Redis缓存

## 目录结构

```
bico-admin/
├── cmd/                    # 应用程序入口点
│   ├── server/            # Web服务入口
│   │   ├── main.go        # 主程序入口
│   │   ├── wire.go        # Wire依赖注入配置
│   │   └── wire_gen.go    # Wire生成的代码
│   └── migrate/           # 数据库迁移工具
│       └── main.go        # 迁移程序入口
├── internal/              # 私有应用代码
│   ├── master/            # 主控端模块
│   │   ├── handler/       # 主控端处理器
│   │   ├── service/       # 主控端业务逻辑
│   │   ├── repository/    # 主控端数据访问层
│   │   ├── types/         # 主控端类型定义
│   │   ├── routes/        # 主控端路由注册
│   │   ├── middleware/    # 主控端中间件
│   │   └── provider.go    # 主控端Wire Provider
│   ├── admin/             # 管理端模块
│   │   ├── handler/       # 管理端处理器
│   │   ├── service/       # 管理端业务逻辑
│   │   ├── repository/    # 管理端数据访问层
│   │   ├── types/         # 管理端类型定义
│   │   ├── routes/        # 管理端路由注册
│   │   ├── middleware/    # 管理端中间件
│   │   └── provider.go    # 管理端Wire Provider
│   ├── api/               # API端模块
│   │   ├── handler/       # API端处理器
│   │   ├── service/       # API端业务逻辑
│   │   ├── repository/    # API端数据访问层
│   │   ├── types/         # API端类型定义
│   │   ├── routes/        # API端路由注册
│   │   ├── middleware/    # API端中间件
│   │   └── provider.go    # API端Wire Provider
│   └── shared/            # 共享模块
│       ├── model/         # 数据模型
│       ├── types/         # 通用类型定义
│       ├── middleware/    # 通用中间件
│       └── provider.go    # 共享组件Wire Provider
├── pkg/                   # 可被外部应用使用的库代码
│   ├── config/            # 配置读取逻辑
│   │   ├── config.go      # 配置结构定义
│   │   └── loader.go      # 配置加载器
│   ├── database/          # 数据库连接
│   │   ├── mysql.go       # MySQL连接
│   │   ├── postgres.go    # PostgreSQL连接
│   │   ├── sqlite.go      # SQLite连接
│   │   └── redis.go       # Redis连接
│   ├── logger/            # 日志工具
│   │   └── logger.go      # 日志实现
│   └── response/          # 统一响应封装
│       └── response.go    # 响应结构定义
├── config/                # 配置文件目录
│   ├── app.yml            # 基础配置模板
│   └── app.dev.yml        # 开发环境配置
├── data/                  # 数据文件目录
│   ├── bico_admin.db      # 生产环境SQLite数据库
│   └── bico_admin_dev.db  # 开发环境SQLite数据库
├── docs/                  # 项目文档
│   ├── project-structure.md  # 项目结构说明
│   ├── database-sqlite.md    # 数据库文档
│   ├── docs.go              # Swagger文档配置
│   ├── swagger.json         # Swagger JSON文档
│   └── swagger.yaml         # Swagger YAML文档
├── tests/                 # 测试文件
│   └── config/            # 配置相关测试
│       ├── config_test.go     # 配置测试
│       ├── integration_test.go # 集成测试
│       ├── loader_test.go     # 加载器测试
│       ├── switch_test.go     # 切换测试
│       └── testdata/          # 测试数据
├── logs/                  # 日志文件目录
│   └── app.log            # 应用日志文件
├── .gitignore             # Git忽略文件
├── go.mod                 # Go模块定义
├── go.sum                 # Go模块校验和
├── Makefile               # 构建命令
└── README.md              # 项目说明文档
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

## 技术栈

- **Web框架**: Gin
- **ORM框架**: GORM
- **数据库**: MySQL/PostgreSQL/SQLite
- **缓存**: Redis
- **依赖注入**: Google Wire
- **配置管理**: Viper
- **日志**: Zap
- **文档**: Swagger