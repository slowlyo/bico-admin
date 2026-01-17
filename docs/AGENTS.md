# AI Agent 项目说明文档

> 本文档面向 AI 助手，提供项目架构、开发规范和关键信息的快速索引

## 项目概览

**Bico Admin** 是一个前后端分离的现代化后台管理系统：
- **后端**: Go + Gin + GORM（模块内显式装配依赖）
- **前端**: Vue 3 + Artd Pro + Element Plus + TypeScript
- **数据库**: SQLite（开发）/ MySQL（生产）
- **缓存**: Memory / Redis

## 技术架构

### 后端分层架构

```
├── core/           # 核心基础设施层
│   ├── app/        # AppContext 构建 + 模块注册 + 应用生命周期
│   ├── config/     # 配置管理（Viper）
│   ├── db/         # 数据库连接（GORM）
│   ├── cache/      # 缓存接口（Memory/Redis）
│   ├── server/     # Gin引擎 + 框架级路由注册
│   ├── middleware/ # 核心中间件（JWT认证）
│   ├── scheduler/  # 定时调度器（robfig/cron）
│   └── upload/     # 文件上传（支持多种驱动）
│
├── shared/         # 共享层（跨模块复用）
│   ├── model/      # 基础模型（BaseModel）
│   ├── response/   # 统一响应格式
│   ├── jwt/        # JWT令牌管理
│   ├── password/   # 密码加密（bcrypt）
│   ├── pagination/ # 分页工具
│   └── logger/     # 日志工具
│
├── admin/          # 后台管理模块
│   ├── consts/     # 常量定义（权限树）
│   ├── model/      # 数据模型（AdminUser, AdminRole）
│   ├── service/    # 业务逻辑
│   ├── handler/    # HTTP处理器
│   ├── middleware/ # 模块中间件（权限验证、用户状态）
│   └── router.go   # 路由注册
│
├── api/            # 前台API模块
│   ├── handler/
│   ├── service/
│   └── router.go
│
├── job/            # 定时任务
│   ├── task/       # 任务实现
│   └── register.go # 任务注册
│
└── migrate/        # 数据库迁移
    └── migrate.go  # 统一管理模型迁移
```

### 前端架构

```
web/
├── src/
│   ├── api/             # API 接口请求
│   ├── assets/          # 静态资源
│   ├── components/      # 公共组件
│   ├── config/          # 全局配置
│   ├── hooks/           # 自定义 Hook (如 useTable)
│   ├── router/          # 路由配置
│   ├── store/           # 状态管理 (Pinia)
│   ├── views/           # 页面组件
│   ├── App.vue          # 根组件
│   └── main.ts          # 入口文件
└── package.json
```

## 核心概念

### 1. 依赖注入（DI）

**原则**：core 只装配基础设施，业务模块自行装配依赖。

#### Core 侧（基础设施装配）

**位置**: `internal/core/app/context.go`

core 通过 `BuildContext` 构建 `AppContext`，只包含基础设施对象：

```go
ctx, _ := app.BuildContext(configPath)
// ctx.Cfg / ctx.DB / ctx.Cache / ctx.Engine / ctx.JWT / ctx.Scheduler ...
```

#### 模块侧（模块内 DI）

每个模块通过 `module.go` 作为入口，在 `Register(ctx)` 内自行装配依赖（显式创建 service/handler/router）：

```go
func (m *Module) Register(ctx *app.AppContext) error {
    // 1) 创建 service
    // 2) 创建 handler
    // 3) 创建 router 并注册到 ctx.Engine
    return nil
}
```

### 2. 权限系统

#### 后端权限定义

**位置**: `internal/admin/consts/permissions.go`

```go
// 权限命名规范: 模块:资源:操作
const (
    PermDashboardMenu      = "dashboard:menu"          // 菜单权限
    PermSystemManage       = "system:manage"           // 模块权限
    PermAdminUserList      = "system:admin_user:list"  // 操作权限
    PermAdminUserCreate    = "system:admin_user:create"
    PermAdminUserEdit      = "system:admin_user:edit"
    PermAdminUserDelete    = "system:admin_user:delete"
)

// 权限树结构
var AllPermissions = []Permission{
    {
        Key: PermSystemManage,
        Label: "系统管理",
        Children: []Permission{...},
    },
}
```

#### 路由权限绑定

**后端** (`internal/admin/router.go`):
```go
users.GET("", r.permMiddleware.RequirePermission(consts.PermUserList), r.userHandler.List)
```

**前端** (`web/src/router/`):
```ts
{
  path: 'users',
  name: 'User',
  component: () => import('@/views/system/user/index.vue'),
  meta: {
    title: '用户管理',
    permissions: ['system:admin_user:menu'] // 对应后端权限 key
  }
}
```

#### 权限验证流程

1. 用户登录后，后端返回用户的权限列表
2. 前端将权限列表存储在 Pinia Store (`useUserStore`)
3. `v-auth` 指令或 `useAuth()` Hook 进行权限校验
4. 路由守卫根据 `meta.permissions` 控制页面访问
5. 后端中间件 `RequirePermission` 验证接口访问权限

### 3. 认证流程

#### JWT认证

**配置**: `config/config.yaml`
```yaml
jwt:
  secret: "bico-admin-secret-key-change-in-production"
  expire_hours: 168  # 7天
```

**流程**:
1. 用户登录 -> 生成JWT token
2. 前端存储 token 到 localStorage
3. 请求时带上 `Authorization: Bearer {token}`
4. JWT中间件验证 token -> 提取 user_id
5. 后续中间件可从 `c.Get("user_id")` 获取用户ID

**Token黑名单**:
- 退出登录时将 token 加入黑名单（缓存）
- JWT中间件会检查黑名单
- 7天后自动过期清除

#### 中间件链

```go
// 认证路由的中间件顺序
authorized := admin.Group("", 
    r.jwtAuth,                    // 1. JWT认证
    r.userStatusMiddleware.Check(), // 2. 用户状态检查（是否禁用）
)

// 需要权限的路由
users.GET("", 
    r.permMiddleware.RequirePermission(consts.PermUserList), // 3. 权限验证
    r.userHandler.List,
)
```

### 4. 统一响应格式

**位置**: `internal/shared/response/response.go`

```go
// 成功响应
{
  "code": 0,
  "msg": "success",
  "data": {...}
}

// 错误响应
{
  "code": 400,  // 业务错误码
  "msg": "参数错误"
}

// 分页响应
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 5. 数据模型规范

**BaseModel** (`internal/shared/model/base.go`):
```go
type BaseModel struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**命名约定**:
- 模型: 单数形式 `AdminUser`
- 表名: 复数形式 `admin_users`
- JSON字段: 蛇形命名 `created_at`

## 开发流程

### 新增后台功能（声明式 CRUD 框架）

使用 `internal/pkg/crud` 包，**只需一个文件** 即可完成后端功能开发。

假设要添加"文章管理"功能：

#### 1. 创建 Handler (`internal/admin/handler/article_handler.go`)

```go
package handler

import (
    "bico-admin/internal/admin/model"
    "bico-admin/internal/pkg/crud"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// 定义权限（自动生成 menu/list/create/edit/delete）
var articlePerms = crud.NewCRUDPerms("system", "article", "文章管理")

type ArticleHandler struct {
    crud.BaseHandler
    db *gorm.DB
}

func NewArticleHandler(db *gorm.DB) *ArticleHandler {
    return &ArticleHandler{db: db}
}

// 声明模块配置（路由 + 权限 自动注册）
func (h *ArticleHandler) ModuleConfig() crud.ModuleConfig {
    return crud.ModuleConfig{
        Name:             "article",
        Group:            "/articles",
        ParentPermission: PermSystemManage,
        Permissions:      articlePerms.Tree,
        Routes:           articlePerms.Routes(),
    }
}

// 实现业务方法
func (h *ArticleHandler) List(c *gin.Context) {
    var req struct{ Title string }
    h.BindQuery(c, &req)
    
    query := h.db.Model(&model.Article{})
    if req.Title != "" {
        query = query.Where("title LIKE ?", "%"+req.Title+"%")
    }
    
    var articles []model.Article
    h.QueryList(c, query, &articles)  // 自动处理分页和响应
}

func (h *ArticleHandler) Get(c *gin.Context) {
    id, err := h.ParseID(c)
    if err != nil {
        return
    }
    var article model.Article
    if h.QueryOne(c, h.db.Where("id = ?", id), &article, "文章不存在") {
        h.Success(c, article)
    }
}

func (h *ArticleHandler) Create(c *gin.Context) {
    var req struct{ Title, Content string }
    if err := h.BindJSON(c, &req); err != nil {
        return
    }
    article := &model.Article{Title: req.Title, Content: req.Content}
    h.ExecTx(c, h.db, func(tx *gorm.DB) error {
        return tx.Create(article).Error
    }, "创建成功", article)
}

func (h *ArticleHandler) Update(c *gin.Context) { /* 类似逻辑 */ }
func (h *ArticleHandler) Delete(c *gin.Context) { /* 类似逻辑 */ }

var _ crud.Module = (*ArticleHandler)(nil)
```

然后在 `internal/admin/module.go` 中把新模块加入 `[]crud.Module` 列表（由模块显式装配）：

```go
return []crud.Module{
    handler.NewAdminUserHandler(db),
    handler.NewAdminRoleHandler(db),
    NewArticleHandler(db),
}
```

#### 2. 创建模型 (`internal/admin/model/article.go`)

```go
package model

import "bico-admin/internal/core/model"

type Article struct {
    model.BaseModel
    Title   string `gorm:"size:200;not null" json:"title"`
    Content string `gorm:"type:text" json:"content"`
    Status  int    `gorm:"default:1" json:"status"`
}
```

#### 3. 数据库迁移 (`internal/migrate/migrate.go`)

```go
func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        // ...
        &adminModel.Article{},
    )
}
```

#### 4. 前端路由 (`web/src/router/modules/content.ts`)

```ts
{
  path: 'articles',
  name: 'Article',
  component: () => import('@/views/content/article/index.vue'),
  meta: {
    title: '文章管理',
    permissions: ['content:article:menu']
  }
}
```

#### 9. 前端页面 (`web/src/views/content/article/index.vue`)

```vue
<template>
  <div class="art-full-height">
    <ArtTable
      :loading="loading"
      :data="data"
      :columns="columns"
      :pagination="pagination"
      @refresh="refreshData"
    />
  </div>
</template>

<script setup lang="ts">
import { useTable } from '@/hooks/core/useTable'
import { fetchGetArticleList } from '@/api/content'

const { data, loading, columns, pagination, refreshData } = useTable({
  core: {
    apiFn: fetchGetArticleList,
    columnsFactory: () => [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'title', label: '标题' }
    ]
  }
})
</script>
```

### 运行和测试

```bash
# 后端
make migrate   # 执行迁移
make serve     # 启动服务

# 前端
cd web
pnpm dev

# 测试
curl http://localhost:8080/admin-api/articles \
  -H "Authorization: Bearer {token}"
```

## 关键文件位置

### 配置文件
- `config/config.yaml` - 主配置文件
- `web/config/routes.ts` - 前端路由配置
- `web/config/proxy.ts` - 开发代理配置

### 核心文件
- `internal/core/app/context.go` - AppContext 构建 + 模块注册
- `internal/admin/consts/permissions.go` - 权限定义
- `internal/admin/router.go` - 后台路由
- `internal/migrate/migrate.go` - 数据库迁移

### 中间件
- `internal/core/middleware/jwt.go` - JWT认证
- `internal/admin/middleware/permission.go` - 权限验证
- `internal/admin/middleware/user_status.go` - 用户状态检查

### 前端核心
- `web/src/main.ts` - 应用入口
- `web/src/router/` - 路由与权限拦截
- `web/src/utils/http/` - 请求拦截器
- `web/src/hooks/core/useTable.ts` - 核心表格 Hook

## 编码规范

### 后端规范

**严格遵守**:
- SOLID 原则（单一职责、开闭原则等）
- DRY（不重复代码）
- KISS（保持简单）
- YAGNI（不过度设计）

**具体要求**:
1. **先实现功能，后整理 import**
2. **禁止在循环中执行 I/O 操作**（数据库查询、API调用等）
3. **使用 import 导入 class**，避免直接使用命名空间/路径
4. **错误处理**: 向上传递，在 Handler 层统一处理
5. **避免全局变量**: 使用 DI 注入依赖

### 前端规范

1. 使用 Vue 3 组合式 API (`<script setup>`)
2. 复杂表格页面优先使用 `useTable` Hook
3. 组件拆分合理，通用组件提取到 `src/components/core/`，业务组件提取到对应视图目录下的 `modules/`
4. API 调用统一在 `src/api/` 目录，命名以 `fetch` 开头
5. 样式使用 TailwindCSS + SCSS

### 命名规范

**Go**:
- 包名: 小写单词 `package handler`
- 文件名: 蛇形 `admin_user.go`
- 类型: 大驼峰 `AdminUser`
- 方法: 大驼峰（公开） `GetUser` / 小驼峰（私有） `getUserID`
- 常量: 大驼峰 `PermAdminUserList`

**TypeScript**:
- 文件名: 短横线 `admin-users.tsx` 或小驼峰 `adminUsers.tsx`
- 组件: 大驼峰 `AdminUserList`
- 函数: 小驼峰 `getUserList`
- 常量: 大写蛇形 `API_BASE_URL`

## 常见问题

### 如何调试权限问题？

1. 检查用户是否有该权限: 查看 `admin_roles` 表的 `permissions` 字段
2. 检查前端路由配置: `routes.ts` 的 `access` 字段
3. 检查后端路由: `router.go` 的 `RequirePermission` 参数
4. 查看浏览器 Network: `/auth/current-user` 接口返回的权限列表

### 如何添加新的中间件？

```go
// 1. 创建中间件函数
func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置处理
        c.Next()
        // 后置处理
    }
}

// 2. 在路由中使用
admin.Use(MyMiddleware())
```

### 如何处理文件上传？

```go
// 后端
func (h *Handler) Upload(c *gin.Context) {
    file, _ := c.FormFile("file")
    
    // 使用 uploader
    url, err := h.uploader.Upload(file)
    if err != nil {
        response.Error(c, 500, err.Error())
        return
    }
    
    response.SuccessWithData(c, gin.H{"url": url})
}

// 前端
const formData = new FormData();
formData.append('file', file);
await uploadFile(formData);
```

### 如何使用缓存？

```go
// 注入缓存
type MyService struct {
    cache cache.Cache
}

// 使用
s.cache.Set("key", value, 10*time.Minute)
val, err := s.cache.Get("key")
```

## 数据库注意事项

### SQLite vs MySQL

**开发环境** (SQLite):
```yaml
database:
  driver: sqlite
  sqlite:
    path: storage/data.db
```

**生产环境** (MySQL):
```yaml
database:
  driver: mysql
  mysql:
    host: localhost
    port: 3306
    username: root
    password: your_password
    database: bico_admin
```

### 执行迁移

```bash
# 第一次运行时必须执行
make migrate

# 或
go run cmd/main.go migrate
```

**注意**: 迁移会自动创建默认管理员账户（admin/admin）

## 部署说明

### 后端编译

```bash
# 编译
make build

# 运行
./bin/bico-admin serve -c config/prod.yaml
```

### 前端构建

```bash
cd web
pnpm build
# 产物在 web/dist 目录
```

### 生产环境检查清单

- [ ] 修改 `jwt.secret` 为强随机字符串
- [ ] 修改数据库为 MySQL
- [ ] 配置 Redis 缓存（可选）
- [ ] 修改默认管理员密码
- [ ] 设置 `server.mode` 为 `release`
- [ ] 配置文件上传（七牛云/阿里云/腾讯云）

## 总结

这是一个**标准的企业级后台管理系统**，核心特点：

1. **前后端分离**: 清晰的 API 接口设计
2. **权限体系完善**: 树形权限 + 前后端双重验证
3. **依赖注入**: 代码解耦，易于测试和维护
4. **统一规范**: 响应格式、错误处理、命名约定
5. **可扩展性强**: 模块化设计，易于添加新功能

开发时优先参考现有代码，保持风格一致！
