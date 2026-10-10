# AI Agent 项目说明文档

> 本文档面向 AI 助手，提供项目架构、开发规范和关键信息的快速索引

## 项目概览

**Bico Admin** 是一个前后端分离的现代化后台管理系统：
- **后端**: Go + Gin + GORM（模块内显式装配依赖）
- **前端**: Vue 3 + Vben Admin + antdv-next + TypeScript
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
├── pkg/            # 工具包（无状态，跨模块复用）
│   ├── crud/       # 声明式 CRUD 框架
│   ├── response/   # 统一响应格式
│   ├── jwt/        # JWT令牌管理
│   ├── password/   # 密码加密（bcrypt）
│   ├── pagination/ # 分页工具
│   └── excel/      # Excel 导入导出
│
├── admin/          # 后台管理模块
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
├── apps/admin/     # 后台应用
│   └── src/
│       ├── router/routes/modules/  # 路由
│       ├── views/                  # 页面
│       ├── api/                    # 接口
│       └── layouts/                # 布局
├── packages/       # Vben 框架包
├── embed.go        # 生产构建嵌入 web/dist
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

**位置**: 基础权限在 `internal/admin/handler/permissions.go`，各模块权限由 `crud.NewCRUDPerms` 生成。

```go
// 权限命名规范: 命名空间:资源:操作
const (
    PermDashboardMenu = "dashboard:menu" // 菜单权限
    PermSystemManage  = "system:manage"  // 模块权限
)

// 自动生成 menu/list/create/edit/delete 五个权限，如 system:admin_user:list
var userPerms = crud.NewCRUDPerms("system", "admin_user", "用户管理")
```

权限树由 `crud` 包汇总，`GetAllPermissions()` 返回完整树，超级管理员角色自动拥有全部权限。

#### 路由权限绑定

**后端**: 在 `ModuleConfig()` 的 `Routes` 中声明，路由注册时自动挂载 `RequirePermission`：
```go
Routes: userPerms.Routes(), // 标准 CRUD 路由，各自绑定 list/create/edit/delete 权限
```

**前端** (`web/apps/admin/src/router/routes/modules/`):
```ts
{
  name: 'AdminUsers',
  path: '/system/users',
  component: () => import('#/views/system/users/index.vue'),
  meta: {
    authority: ['system:admin_user:menu'],
    title: '管理员',
  },
}
```

#### 权限验证流程

1. 用户登录后，前端调用 `/auth/current-user` 获取权限列表
2. 前端将权限列表写入 Vben 的 access store（`store/auth.ts`）
3. 路由 `meta.authority` 与权限列表比对，决定菜单和路由是否可见（`accessMode: 'frontend'`）
4. 按钮使用 `v-access:code="'system:admin_user:create'"` 控制显示
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

// 需要权限的 CRUD 路由由 crud.ModuleRouter 统一挂载：
// jwtAuth -> userStatus -> RequirePermission(route.Permission) -> handler
```

### 4. 统一响应格式

**位置**: `internal/pkg/response/response.go`

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

**BaseModel** (`internal/core/model/base.go`):
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

然后在 `internal/admin/module.go` 的 `NewCRUDModules` 中把新模块加入列表（运行时路由与 Swagger 生成共用这份列表）：

```go
return []crud.Module{
    handler.NewAdminUserHandler(db, cacheInvalidator),
    handler.NewAdminRoleHandler(db, cacheInvalidator),
    handler.NewArticleHandler(db),
}
```

> 现有的用户、角色模块使用 `crud.CRUDHandler` 泛型（通过 hook 声明差异逻辑），可参考 `admin_user_handler.go`；详见 [crud-pkg.md](./crud-pkg.md)。

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
func AutoMigrate(db *gorm.DB, mode string) error {
    if err := db.AutoMigrate(
        // ...
        &adminModel.Article{},
    ); err != nil {
        return err
    }
    // ...
}
```

#### 4. 前端路由 (`web/apps/admin/src/router/routes/modules/articles.ts`)

```ts
{
  meta: {
    authority: ['content:article:menu'],
    icon: 'lucide:file-text',
    title: '文章管理',
  },
  name: 'Articles',
  path: '/content/articles',
  component: () => import('#/views/content/articles/index.vue'),
}
```

#### 5. 前端页面 (`web/apps/admin/src/views/content/articles/index.vue`)

页面用 Vue 单文件组件，放在 `web/apps/admin/src/views`，接口放在 `web/apps/admin/src/api`。写法参考 [frontend-crud.md](./frontend-crud.md) 和 [frontend-services.md](./frontend-services.md)。

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
- `web/apps/admin/src/router/routes/modules/` - 前端路由
- `web/apps/admin/vite.config.ts` - 开发代理配置

### 核心文件
- `internal/core/app/context.go` - AppContext 构建 + 模块注册
- `internal/admin/handler/permissions.go` - 基础权限定义
- `internal/admin/router.go` - 后台路由
- `internal/migrate/migrate.go` - 数据库迁移

### 中间件
- `internal/core/middleware/jwt.go` - JWT认证
- `internal/admin/middleware/permission.go` - 权限验证
- `internal/admin/middleware/user_status.go` - 用户状态检查

### 前端核心
- `web/apps/admin/src/main.ts` - 应用入口
- `web/apps/admin/src/router/access.ts` - 权限控制
- `web/apps/admin/src/api/request.ts` - 请求拦截器

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

1. 使用 TypeScript，充分利用类型系统
2. 组件拆分合理，单个文件不超过 500 行
3. 组件不应该太大：可复用的通用组件提取到 `src/components/`，仅页面内使用的组件提取到对应页面目录下的 `./components/`
4. API 调用统一在 `src/api/` 目录
5. 权限控制使用路由 `meta.authority` 和 `v-access:code` 指令

### 命名规范

**Go**:
- 包名: 小写单词 `package handler`
- 文件名: 蛇形 `admin_user.go`
- 类型: 大驼峰 `AdminUser`
- 方法: 大驼峰（公开） `GetUser` / 小驼峰（私有） `getUserID`
- 常量: 大驼峰 `PermAdminUserList`

**TypeScript**:
- 文件名: 短横线 `admin-users`（页面目录）、`avatar-upload.vue`（组件）
- 组件: 大驼峰 `AdminUserList`
- 函数: 小驼峰 `getUserList`
- 常量: 大写蛇形 `API_BASE_URL`

## 常见问题

### 如何调试权限问题？

1. 检查用户是否有该权限: 查看 `admin_role_permissions` 表中该角色的权限记录
2. 检查前端路由配置: `router/routes/modules/*.ts` 的 `meta.authority`
3. 检查后端路由: 模块 `ModuleConfig().Routes` 中的 `Permission`
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
