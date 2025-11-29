# Bico Admin

<img src="./web/public/logo.png" width="200" />

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

使用声明式 CRUD 框架，**只需一个文件**：

```go
// internal/admin/handler/article_handler.go
package handler

import (
    "your-project/internal/admin/model"
    "your-project/internal/pkg/crud"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// 1. 定义权限
var articlePerms = crud.NewCRUDPerms("article", "文章管理")

// 2. 定义 Handler
type ArticleHandler struct {
    crud.BaseHandler
    db *gorm.DB
}

func NewArticleHandler(db *gorm.DB) *ArticleHandler {
    return &ArticleHandler{db: db}
}

// 3. 声明模块配置（路由 + 权限）
func (h *ArticleHandler) ModuleConfig() crud.ModuleConfig {
    return crud.ModuleConfig{
        Name:             "article",
        Group:            "/articles",
        ParentPermission: PermSystemManage,
        Permissions:      articlePerms.Tree,
        Routes:           articlePerms.Routes(),
    }
}

// 4. 实现业务方法
func (h *ArticleHandler) List(c *gin.Context)   { /* ... */ }
func (h *ArticleHandler) Get(c *gin.Context)    { /* ... */ }
func (h *ArticleHandler) Create(c *gin.Context) { /* ... */ }
func (h *ArticleHandler) Update(c *gin.Context) { /* ... */ }
func (h *ArticleHandler) Delete(c *gin.Context) { /* ... */ }

// 5. 自动注册
func init() {
    crud.RegisterModule(NewArticleHandler)
}

var _ crud.Module = (*ArticleHandler)(nil)
```

**完成！** 无需修改 `router.go`、`container.go` 或其他文件。

详细文档见 [CRUD 包文档](./docs/crud-pkg.md)

### 前端路由

在 `web/config/routes.ts` 配置：
```ts
{
  path: "/system/articles",
  name: "articles",
  component: "./system/articles",
  access: "system:article:menu"  // 对应后端权限 key
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

- [后端 CRUD 框架](./docs/crud-pkg.md) - 声明式后端开发指南
- [前端 CRUD 组件](./docs/frontend-crud.md) - CrudTable 使用指南
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
