# Bico Admin

基于 Go + React 构建的现代化后台管理系统。

## 技术栈

### 后端
- **[Gin](https://github.com/gin-gonic/gin)** - HTTP Web 框架
- **[GORM](https://gorm.io/)** - ORM 数据库操作
- **[Viper](https://github.com/spf13/viper)** - 配置管理
- **[Dig](https://github.com/uber-go/dig)** - 依赖注入容器
- **[Cobra](https://github.com/spf13/cobra)** - 命令行框架

### 前端
- **[React 19](https://react.dev/)** - UI 框架
- **[Ant Design Pro](https://pro.ant.design/)** - 企业级中后台解决方案
- **[UmiJS 4](https://umijs.org/)** - 企业级前端框架
- **[TypeScript](https://www.typescriptlang.org/)** - 类型安全
- **[pnpm](https://pnpm.io/)** - 包管理器

## 快速开始

### 后端

```bash
# 修改配置文件
vim config/config.yaml

# 执行数据库迁移
make migrate

# 启动服务
make serve
```

### 前端

```bash
cd web

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 构建生产版本
pnpm build
```

## 开发指南

### 新增后台功能流程

1. **定义权限常量** - 在 `internal/admin/consts/permissions.go` 添加权限定义
2. **编写 Service** - 在 `internal/admin/service/` 创建业务逻辑
3. **编写 Handler** - 在 `internal/admin/handler/` 创建处理器
4. **注册 DI** - 在 `internal/core/app/container.go` 注册到容器
5. **配置路由** - 在 `internal/admin/router.go` 添加路由和权限中间件
6. **前端路由** - 在 `web/config/routes.ts` 配置前端路由和 access 权限
7. **实现页面** - 在 `web/src/pages/` 编写页面组件

### 路由配置

**后端路由** (`internal/admin/router.go`)
```go
// 需要权限的路由
users.GET("", r.permMiddleware.RequirePermission(consts.PermUserList), r.userHandler.List)
```

**前端路由** (`web/config/routes.ts`)
```ts
{
  path: "/system/users",
  name: "users",
  component: "./system/users",
  access: "system:user:menu"  // 对应后端权限 key
}
```

### 权限定义

权限采用树形结构，定义在 `internal/admin/consts/permissions.go`：

```go
const (
    PermModuleManage = "module:manage"        // 模块菜单
    PermModuleList   = "module:list"          // 查看列表
    PermModuleCreate = "module:create"        // 创建
    PermModuleEdit   = "module:edit"          // 编辑
    PermModuleDelete = "module:delete"        // 删除
)
```

### 依赖注入 (DI)

使用 Uber Dig 管理依赖，在 `internal/core/app/container.go` 注册：

```go
providers := []interface{}{
    // 基础设施层
    provideDatabase,
    provideCache,
    
    // 服务层
    service.NewUserService,
    
    // 处理层
    handler.NewUserHandler,
}
```

Handler 会自动注入所需依赖：
```go
func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}
```

## 常用命令

```bash
make help      # 查看所有可用命令
make build     # 编译应用
make serve     # 启动服务
make migrate   # 数据库迁移
make clean     # 清理构建产物
make tidy      # 整理依赖
```

## 项目文档

详细文档位于 `docs/` 目录：

- [项目结构说明](./docs/structure.md)
- [认证 API](./docs/auth-api.md)
- [缓存机制](./docs/cache.md)
- [限流中间件](./docs/rate-limit.md)
- [配置热更新](./docs/config-hot-reload.md)

## 开发环境

- Go 1.21+
- MySQL 5.7+
- Node.js 20+
- pnpm 9+ (前端包管理器)

## License

MIT
