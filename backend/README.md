# Bico Admin Backend

基于Go Fiber + GORM + MySQL的后端服务，采用模块化设计，AI友好的代码结构。

## 🏗️ 项目结构

```
backend/
├── cmd/server/           # 统一服务入口
│   └── main.go          # 主程序入口
├── core/                # 框架核心（系统默认功能）
│   ├── config/          # 配置管理
│   ├── middleware/      # 核心中间件
│   ├── model/          # 系统基础模型
│   ├── repository/     # 核心数据访问层
│   ├── service/        # 核心业务服务
│   ├── handler/        # 核心处理器
│   └── router/         # 核心路由
├── modules/            # 业务模块目录
│   ├── admin/          # 后台管理模块
│   │   ├── handler/    # 后台管理处理器
│   │   ├── service/    # 后台管理业务服务
│   │   ├── model/      # 后台管理数据模型
│   │   ├── repository/ # 后台管理数据访问
│   │   └── router/     # 后台管理路由
│   └── api/            # 对外API模块
│       ├── handler/    # API处理器
│       ├── service/    # API业务服务
│       ├── model/      # API数据模型
│       ├── repository/ # API数据访问
│       └── router/     # API路由
├── business/           # 业务方法封装
│   ├── base.go         # 基础业务方法
│   ├── crud.go         # CRUD操作封装
│   └── ...             # 其他业务封装
├── pkg/                # 公共包
│   ├── utils/          # 工具函数
│   ├── validator/      # 数据验证
│   ├── response/       # 响应格式
│   └── constants/      # 常量定义
├── docs/               # API文档
├── migrations/         # 数据库迁移文件
├── storage/            # 文件存储目录
├── go.mod              # Go模块文件
└── .env.example        # 环境变量示例
```

## 🚀 快速开始

### 环境要求
- Go 1.21+
- MySQL 8.0+
- Redis 6.0+ (可选)

### 安装步骤

1. **安装依赖**
   ```bash
   go mod tidy
   ```

2. **配置环境变量**
   ```bash
   cp .env.example .env
   # 编辑 .env 文件配置数据库连接等信息
   ```

3. **启动服务**
   ```bash
   go run cmd/server/main.go
   ```

### 访问地址
- 后台管理API: http://localhost:8080/admin/api
- 对外API: http://localhost:8080/api
- 健康检查: http://localhost:8080/health

## 📋 技术栈

- **Web框架**: [Go Fiber](https://gofiber.io/) - 高性能Web框架
- **ORM**: [GORM](https://gorm.io/) - Go语言ORM库
- **数据库**: MySQL 8.0+
- **缓存**: Redis (可选)
- **认证**: JWT
- **验证**: go-playground/validator

## 🔧 核心特性

### 模块化设计
- **core**: 框架核心功能，可整体更新
- **admin**: 后台管理业务模块
- **api**: 对外API业务模块
- **business**: 通用业务方法封装

### 业务方法封装
提供标准化的CRUD操作：
```go
// 基础CRUD操作
service.CreateOne(data)           // 创建
service.UpdateById(id, data)      // 更新
service.DeleteById(id)            // 删除
service.GetById(id)               // 获取单个
service.List(params)              // 列表查询

// 批量操作
service.BatchCreate(data)         // 批量创建
service.BatchUpdate(ids, updates) // 批量更新
service.BatchDelete(ids)          // 批量删除
```

### AI友好设计
- 清晰的目录结构和命名规范
- 完善的注释和文档
- 标准化的代码模式
- 类型安全的接口设计

## 📖 开发指南

### 添加新的业务模块
1. 在对应模块目录下创建model、repository、service、handler文件
2. 在router中注册路由
3. 在main.go中引入路由

### 使用业务封装
```go
// 创建CRUD服务
crudService := business.NewCRUDService[model.YourModel](db)

// 使用封装的方法
result, err := crudService.List(business.ListParams{
    Page:     1,
    PageSize: 10,
    Sort:     "id",
    Order:    "desc",
})
```

### 数据库迁移
数据库表会在启动时自动迁移，新增模型需要在`core/config/database.go`的`autoMigrate`函数中添加。

## 🔒 安全特性

- JWT认证机制
- 请求参数验证
- SQL注入防护
- CORS跨域配置
- 请求限流
- 错误处理中间件

## 📝 注意事项

1. **框架更新**: `core/` 目录包含框架核心功能，可通过脚本整体更新
2. **业务隔离**: `admin/` 和 `api/` 模块的业务代码不会被框架更新影响
3. **配置管理**: 所有配置通过环境变量管理，支持不同环境部署
4. **日志记录**: 内置日志中间件，支持文件和控制台输出
5. **性能优化**: 使用连接池、缓存等机制优化性能

## 🤝 贡献指南

1. 遵循Go官方代码规范
2. 添加适当的注释和文档
3. 编写单元测试
4. 提交前运行`go fmt`和`go vet`

## 📄 许可证

MIT License
