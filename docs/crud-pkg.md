# CRUD 包使用文档

> 支持多模块分组（admin/api/自定义），自动注册权限和路由

> `internal/pkg/crud` - 声明式 CRUD 框架，让新功能开发只需一个文件。

## 概述

该包提供了一套声明式的 CRUD 框架，核心目标是**减少样板代码**。传统方式开发一个功能需要：

- `handler/xxx_handler.go`
- `service/xxx_service.go`
- 权限常量定义（修改）
- `router.go`（逐条注册路由，修改）

使用本框架后，**只需一个 handler 文件**，路由和权限树由 `ModuleConfig` 声明后自动注册；模块入口只需把 handler 加入 `NewCRUDModules` 列表。

说明：当前项目采用“core 只装配基础设施、模块自管理 DI”的架构，CRUD 模块实例由模块入口（例如 `internal/admin/module.go`）显式装配并传入路由层。

## 快速开始

### 最小示例

```go
package handler

import (
    "your-project/internal/admin/model"
    "your-project/internal/pkg/crud"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// 1. 定义权限
var articlePerms = crud.NewCRUDPerms("system", "article", "文章管理")

// 2. 定义 Handler
type ArticleHandler struct {
    crud.BaseHandler
    db *gorm.DB
}

func NewArticleHandler(db *gorm.DB) *ArticleHandler {
    return &ArticleHandler{db: db}
}

// 3. 实现 ModuleConfig
func (h *ArticleHandler) ModuleConfig() crud.ModuleConfig {
    return crud.ModuleConfig{
        Name:             "article",
        Group:            "/articles",
        ParentPermission: "system:manage",
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

var _ crud.Module = (*ArticleHandler)(nil)
```

### 最小示例（推荐：CRUDHandler）

当业务是标准 CRUD 且差异点可通过 hook 表达时，推荐直接使用 `CRUDHandler`：

```go
package handler

import (
	"bico-admin/internal/admin/model"
	"bico-admin/internal/pkg/crud"
	"errors"

	"gorm.io/gorm"
)

var articlePerms = crud.NewCRUDPerms("system", "article", "文章管理")

type ArticleHandler struct {
	crud.CRUDHandler[model.Article, articleListReq, articleCreateReq, articleUpdateReq]
}

func NewArticleHandler(db *gorm.DB) *ArticleHandler {
	h := &ArticleHandler{}
	h.DB = db
	h.NotFoundMsg = "文章不存在"

	h.SortableFields = []string{"id", "title", "created_at"}

	h.NewModelFromCreate = func(req *articleCreateReq) (*model.Article, error) {
		if req.Title == "" {
			return nil, crud.NewBizError("标题不能为空")
		}
		return &model.Article{Title: req.Title}, nil
	}

	h.BuildUpdates = func(tx *gorm.DB, req *articleUpdateReq, existing *model.Article) (map[string]interface{}, error) {
		updates := map[string]interface{}{}
		if req.Title != "" {
			updates["title"] = req.Title
		}
		return updates, nil
	}

	return h
}

func (h *ArticleHandler) ModuleConfig() crud.ModuleConfig {
	return crud.ModuleConfig{
		Name:             "article",
		Group:            "/articles",
		ParentPermission: PermSystemManage,
		Permissions:      articlePerms.Tree,
		Routes:           articlePerms.Routes(),
	}
}

var _ crud.Module = (*ArticleHandler)(nil)
```

说明：示例为突出 `CRUDHandler` 的配置方式，省略了请求结构体定义与完整校验。

完成！接下来在 `internal/admin/module.go` 的 `NewCRUDModules` 中把该模块加入列表（运行时路由注册与 Swagger 文档生成共用这份列表，避免接口与文档漂移）：

```go
return []crud.Module{
    // ... 其他模块
    handler.NewArticleHandler(db),
}
```

注意：路由通过反射按 `Route.Handler` 的名字查找 handler 方法，方法必须是导出方法且签名为 `func(*gin.Context)`；名字写错会在启动注册路由时直接 panic。

---

## 多模块分组

支持 `admin`、`api` 及自定义模块分组，不同分组可配置不同的中间件。

### 配置分组中间件

在模块入口（如 `internal/admin/module.go`）或模块 router 中配置：

```go
// admin 分组：JWT + 权限 + 用户状态
adminRouter := crud.NewModuleRouter(jwtAuth, permMiddleware, userStatusMiddleware)

// api 分组：仅 JWT，无权限验证
apiRouter := crud.NewModuleRouterWithConfig(crud.RouterConfig{
    AuthMiddleware: jwtAuth,
    // 不设置 PermMiddleware，则不验证权限
})

// 公开分组：无需认证
publicRouter := crud.NewModuleRouterWithConfig(crud.RouterConfig{})

// 注册路由（示例：模块内自行决定注册哪些 modules）
for _, m := range modules {
    adminRouter.RegisterModule(engine.Group("/admin-api"), m)
}
```

---

## 核心组件

### Module 接口

所有模块必须实现此接口：

```go
type Module interface {
    ModuleConfig() ModuleConfig
}
```

### ModuleConfig 结构

```go
type ModuleConfig struct {
    Name             string       // 模块名称，如 "article"
    Group            string       // 路由分组，如 "/articles"
    Description      string       // 描述（可选）
    ParentPermission string       // 父级权限 key
    Permissions      []Permission // 权限树
    Routes           []Route      // 路由配置
    Swagger          SwaggerConfig // Swagger 类型信息（可选）
}

// SwaggerConfig 仅用于文档生成，不参与运行时请求处理
type SwaggerConfig struct {
    Model         interface{} // 模型
    ListRequest   interface{} // 列表查询参数
    CreateRequest interface{} // 创建请求体
    UpdateRequest interface{} // 更新请求体
}
```

填写 `Swagger` 后，`make swagger` 会通过 `cmd/swagger-enhance` 为该模块的标准 CRUD 路由补全文档。

### Permission 结构

```go
type Permission struct {
    Key      string       `json:"key"`      // 权限标识，如 "system:article:list"
    Label    string       `json:"label"`    // 显示名称
    Children []Permission `json:"children"` // 子权限
}
```

### Route 结构

```go
type Route struct {
    Method     string // HTTP 方法: GET, POST, PUT, DELETE, PATCH
    Path       string // 路由路径: "", "/:id", "/all"
    Handler    string // Handler 方法名: "List", "Get", "Create"
    Permission string // 权限 key，为空则不校验
    Public     bool   // 是否公开（不需要登录）
}
```

---

## BaseHandler

嵌入 `crud.BaseHandler` 获得常用方法：

### 请求处理

| 方法 | 说明 |
|------|------|
| `ParseID(c)` | 从路由参数解析 ID，失败自动返回 400 |
| `BindJSON(c, &req)` | 绑定 JSON 请求体，失败自动返回 400 |
| `BindQuery(c, &req)` | 绑定 Query 参数 |
| `GetPagination(c)` | 获取分页参数 |

### 响应方法

| 方法 | 说明 |
|------|------|
| `Success(c, data)` | 成功响应 |
| `SuccessWithMessage(c, msg, data)` | 带消息的成功响应 |
| `SuccessWithPagination(c, data, total)` | 分页响应 |
| `Error(c, msg)` | 400 错误响应 |
| `NotFound(c, msg)` | 404 响应 |

### 通用操作

| 方法 | 说明 |
|------|------|
| `QueryList(c, query, &dest)` | 通用分页查询，自动处理 count/offset/limit，只使用默认排序（`created_at DESC`） |
| `crud.QueryListWithHook(&h.BaseHandler, c, query, &dest, order, after)` | 通用分页查询 + 排序白名单（`crud.ListOrder`） + 结果二次处理（例如补充权限、统计字段） |

客户端排序（`sortField`、`sortOrder`）必须在 `ListOrder.Sortable` 白名单内才生效，否则忽略并使用 `ListOrder.Default`。不要把密码等敏感列放进白名单。
| `QueryOne(c, query, &dest, notFoundMsg)` | 通用单条查询，返回 bool 表示是否找到 |
| `ExecDelete(c, db, model, id)` | 通用删除操作 |
| `ExecTx(c, db, fn, successMsg, data)` | 通用事务操作 |

### 使用示例

```go
// 分页查询 - 一行搞定
func (h *ArticleHandler) List(c *gin.Context) {
    var req listReq
    h.BindQuery(c, &req)
    
    query := h.db.Model(&model.Article{})
    if req.Title != "" {
        query = query.Where("title LIKE ?", "%"+req.Title+"%")
    }
    
    var articles []model.Article
    h.QueryList(c, query, &articles)  // 自动处理分页和响应
}

// 分页查询 + 二次处理（示例：回填额外字段）
func (h *ArticleHandler) ListWithExtra(c *gin.Context) {
    var req listReq
    h.BindQuery(c, &req)

    query := h.db.Model(&model.Article{})

    var articles []model.Article
    order := crud.ListOrder{Sortable: []string{"id", "title", "created_at"}}
    crud.QueryListWithHook(&h.BaseHandler, c, query, &articles, order, func(items []model.Article) error {
        // 这里做 items 的二次处理，例如：补充统计字段、聚合关联数据等
        // 返回 error 会自动转为统一错误响应
        return nil
    })
}

// 单条查询
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

// 事务操作
func (h *ArticleHandler) Create(c *gin.Context) {
    var req createReq
    if err := h.BindJSON(c, &req); err != nil {
        return
    }
    
    article := &model.Article{Title: req.Title}
    h.ExecTx(c, h.db, func(tx *gorm.DB) error {
        return tx.Create(article).Error
    }, "创建成功", article)
}
```

---

## CRUDPerms 辅助

快速生成标准 CRUD 权限和路由：

```go
// 基本用法
var perms = crud.NewCRUDPerms("system", "article", "文章管理")

// 生成的权限 key:
// - perms.Menu   = "system:article:menu"
// - perms.List   = "system:article:list"
// - perms.Create = "system:article:create"
// - perms.Edit   = "system:article:edit"
// - perms.Delete = "system:article:delete"

// 生成权限树
perms.Tree  // []Permission，包含菜单和子权限

// 生成标准路由
perms.Routes()  // GET /, GET /:id, POST /, PUT /:id, DELETE /batch, DELETE /:id
```

### 添加额外权限

```go
var perms = crud.NewCRUDPerms("system", "article", "文章管理").WithExtra(
    crud.Permission{Key: "system:article:publish", Label: "发布文章"},
    crud.Permission{Key: "system:article:export", Label: "导出文章"},
)
```

### 添加额外路由

```go
Routes: perms.RoutesWithExtra(
    crud.Route{Method: "POST", Path: "/:id/publish", Handler: "Publish", Permission: "system:article:publish"},
    crud.Route{Method: "GET", Path: "/export", Handler: "Export", Permission: "system:article:export"},
),
```

---

## 路由辅助函数

```go
// 公开路由（不需要登录）
crud.PublicRoute("GET", "/public", "PublicList")

// 需要登录但不需要权限
crud.AuthRoute("GET", "/my", "MyList")

// 需要权限
crud.PermRoute("DELETE", "/:id", "Delete", "system:article:delete")
```

---

## 工具函数

```go
// 数组去重
ids := crud.UniqueUints([]uint{1, 2, 2, 3})  // [1, 2, 3]

// 判断记录是否存在（常用于唯一性校验）
exists, err := crud.Exists(db, &model.AdminUser{}, "username = ?", "admin")
if err != nil { /* ... */ }
if exists { /* ... */ }

// 获取所有权限
crud.GetAllPermissions()     // []Permission
crud.GetAllPermissionKeys()  // []string
```

---

## CRUDHandler 完整示例（推荐）

参考 `internal/admin/handler/admin_user_handler.go` + `internal/admin/module.go`。

`CRUDHandler` 是在 `BaseHandler` 之上进一步封装的“配置式 CRUD”。

核心目标：

- **单个 handler 文件尽量少代码**完成完整 CRUD
- CRUD 的通用流程（参数绑定、ID 解析、事务、统一响应）下沉到框架
- 业务差异通过一组 hook 注入（例如 preload、唯一性校验、关联表同步、回填字段）

### 适用场景

- 数据表具备标准 `List/Get/Create/Update/Delete`
- 需要在 CRUD 内做少量业务差异（例如密码加密、关联表写入、二次回填）

不适用场景：

- 认证、上传等强业务流程（建议继续用普通 handler + service）

### 常用配置字段

| 字段 | 说明 |
|------|------|
| `DB` | 必填，GORM DB |
| `NotFoundMsg` | 404 文案，如“用户不存在” |
| `CreateSuccessMsg` / `UpdateSuccessMsg` / `DeleteSuccessMsg` | 成功提示文案（不配默认“创建成功/更新成功/删除成功”） |
| `EnabledField` | 启用字段名（默认 `enabled`） |
| `EnabledSuccessMsg` | 启用/禁用成功提示（默认“更新成功”） |
| `SortableFields` | 允许客户端排序的列名白名单；为空则忽略 `sortField`，只用默认排序 |
| `DefaultOrder` | 默认排序子句，如 `"sort ASC, id DESC"`；为空时使用 `created_at DESC`（表没有该列时必须配置） |
| `UniqueFields` | 需要全表唯一的列，如 `[]crud.Unique{{Column: "username", Err: ErrUsernameExists}}`；创建与更新时在事务内自动查重（更新时排除自身），数据库仍应建唯一索引兜底并发 |

### hook 一览

执行顺序：`Before*` 与 `After*` 都在**同一个事务内**，`*Commit` 在事务提交之后。

| 阶段 | hook | 说明 |
|------|------|------|
| 列表 | `BuildListQuery(db, req)` | 构建列表查询（可选，不配则查询整张表） |
| | `AfterList(items)` | 列表返回前二次处理（回填权限、统计字段等） |
| 详情 | `BuildGetQuery(db)` | 构建详情查询（可选，常用于 `Preload`） |
| | `AfterGet(item)` | 详情返回前二次处理 |
| 创建 | `NewModelFromCreate(req)` | 请求转模型（必填，在事务之外执行，不要在这里查库） |
| | `BeforeCreate(tx, item, req)` | 插入前校验、补字段 |
| | `AfterCreate(tx, item, req)` | 插入后逻辑，如写关联表 |
| | `AfterCreateCommit(item, req)` | 提交后执行，适合缓存失效等外部副作用 |
| 更新 | `BuildUpdateQuery(tx)` | 更新时的查询（可选，用于 `Preload` 或加锁） |
| | `BuildUpdates(tx, req, existing)` | 请求转 `updates map`（必填，运行在事务内） |
| | `BeforeUpdate(tx, id, existing, req, updates)` | 更新前校验，可改写 `updates` |
| | `AfterUpdate(tx, id, existing, req)` | 更新后逻辑，如同步关联表 |
| | `AfterUpdateCommit(id, existing, req, updates)` | 提交后执行；`updates` 是实际写入的字段 |
| 删除 | `BeforeDelete(tx, ids)` | 删除前校验与关联清理，单个与批量共用 |
| | `AfterDelete(tx, ids)` | 删除后逻辑 |
| | `AfterDeleteCommit(ids)` | 提交后执行，单个删除时 `ids` 只有一个元素 |

### 返回值与详情一致

只要配置了 `BuildGetQuery` 或 `AfterGet`，创建和更新成功后会在提交之后按同一口径重新加载记录再返回，无需再写“重新加载”的 hook。未配置时直接返回内存中的记录。

### 错误处理

hook 与 `NewModelFromCreate`、`BuildUpdates` 返回的错误按类型处理：

- `crud.NewBizError("提示")`：业务错误，消息原样返回给前端。
- 其他错误（数据库驱动错误等）：只记日志，前端收到统一的“服务器内部错误”，避免泄露 SQL、表名和约束名。

```go
var ErrNameExists = crud.NewBizError("名称已存在") // 可作为哨兵错误，errors.Is 按指针比较
```

### 注意事项

- **缓存失效放在提交后**：事务提交前失效缓存，并发请求可能把旧值回填进去。依据 `updates` 而不是 `req` 判断变更，这样 `UpdateEnabled` 路径（`req` 为零值）也能正确失效。
- **查重用 `UniqueFields`**：不要在 `NewModelFromCreate` 里查库（它在事务之外）；声明 `UniqueFields` 即可在事务内自动查重。
- **保护性校验放在 `BeforeUpdate`**：`UpdateEnabled` 不会调用 `BuildUpdates`，只调用 `BeforeUpdate`。“内置记录不可修改”之类的保护写在 `BuildUpdates` 里会被 `UpdateEnabled` 绕过。
- **提交后钩子拿不到事务内数据**：删除后关联表已清理，无法再据此定位受影响的对象。需要先查再删的失效逻辑（例如按角色定位用户）只能留在 `BeforeDelete` 的事务内，参考 `admin_role_handler.go`。

### 扩展方法

#### 1) UpdateEnabled（通用启用/禁用）

方法：`UpdateEnabled(c)`

- 请求体：`{"enabled": true}`
- 默认更新字段：`enabled`
- 可通过 `EnabledField` 修改字段名

路由示例：

```go
crud.Route{Method: "PATCH", Path: "/:id/enabled", Handler: "UpdateEnabled", Permission: perms.Edit}
```

#### 2) DeleteBatch（批量删除）

方法：`DeleteBatch(c)`

- 请求体：`{"ids": [1,2,3]}`
- 自动去重
- 要求所有 id 都存在，否则整体回滚并返回不存在
- 与单个删除共用 `BeforeDelete` / `AfterDelete` / `AfterDeleteCommit`，在其中按 `ids` 统一清理关联表（避免循环中执行 I/O）

`perms.Routes()` 已包含该路由，无需手动添加。

### Exists（通用存在性判断）

用于唯一性校验：

```go
exists, err := crud.Exists(db, &model.AdminUser{}, "username = ?", req.Username)
if err != nil { return nil, err }
if exists { return nil, crud.NewBizError("用户名已存在") }
```

```go
package handler

import (
	"bico-admin/internal/admin/model"
	"bico-admin/internal/pkg/crud"
	"bico-admin/internal/pkg/password"

	"gorm.io/gorm"
)

var userPerms = crud.NewCRUDPerms("system", "admin_user", "用户管理")

type AdminUserHandler struct {
	crud.CRUDHandler[model.AdminUser, userListReq, createUserReq, updateUserReq]
}

// 真实代码还会注入 service.AuthCacheInvalidator（传 nil 时使用 NoopAuthCacheInvalidator），
// 并在 AfterUpdateCommit / AfterDeleteCommit 中失效鉴权缓存。
func NewAdminUserHandler(db *gorm.DB) *AdminUserHandler {
	h := &AdminUserHandler{}
	h.DB = db
	h.NotFoundMsg = "用户不存在"
	h.SortableFields = []string{"id", "username", "name", "enabled", "created_at", "updated_at"}
	h.UniqueFields = []crud.Unique{{Column: "username", Err: crud.NewBizError("用户名已存在")}}

	h.BuildListQuery = func(db *gorm.DB, req *userListReq) *gorm.DB {
		query := db.Model(&model.AdminUser{}).Preload("Roles")
		if req.Username != "" {
			query = query.Where("username LIKE ?", "%"+req.Username+"%")
		}
		if req.Name != "" {
			query = query.Where("name LIKE ?", "%"+req.Name+"%")
		}
		if req.Enabled != nil {
			query = query.Where("enabled = ?", *req.Enabled)
		}
		return query
	}

	h.BuildGetQuery = func(db *gorm.DB) *gorm.DB {
		return db.Model(&model.AdminUser{}).Preload("Roles")
	}

	h.NewModelFromCreate = func(req *createUserReq) (*model.AdminUser, error) {
		hashed, err := password.Hash(req.Password)
		if err != nil {
			return nil, err
		}
		return &model.AdminUser{Username: req.Username, Password: hashed}, nil
	}

	return h
}

// 说明：示例为突出 CRUDHandler 的配置方式，省略了请求结构体、关联表同步等完整实现。
// 真实项目请参考现有 AdminUser/AdminRole Handler。

func (h *AdminUserHandler) ModuleConfig() crud.ModuleConfig {
	return crud.ModuleConfig{
		Name:             "admin_user",
		Group:            "/admin-users",
		ParentPermission: PermSystemManage,
		Permissions:      userPerms.Tree,
		Routes:           userPerms.Routes(),
	}
}

var _ crud.Module = (*AdminUserHandler)(nil)
```

---

## 代码生成器

标准 CRUD 模块可以一条命令生成前后端样板：

```bash
make gen NAME=article TITLE=文章 FIELDS="title:string:标题,content:text:内容,sort:int:排序"
```

- `NAME`：下划线单数名，路由分组为 `/{kebab}s`，表名为 `{name}s`，权限码为 `system:{name}:*`。
- `FIELDS`：`列名:类型:标签`，逗号分隔；类型可选 `string`（参与搜索）、`text`、`int`、`bool`；`enabled`、`id`、`created_at`、`updated_at` 为保留列，无需声明。
- 生成内容：`internal/admin/model/{name}.go`、`internal/admin/handler/{name}_handler.go` 与对应的 `{name}_handler_test.go`（覆盖创建、列表、更新零值、启停、删除）、前端 `api/system/{kebab}.ts`、`views/system/{kebab}/index.vue` 和 `modules/form.vue`。
- 自动注册：在 `module.go`、`migrate.go`、`router/routes/modules/system.ts` 的 `// gen:*` 锚点前插入条目，请不要删除这些锚点。
- 之后自动执行 `make swagger web-types`，刷新文档与前端类型。目标文件已存在时整体中止，不会覆盖手写代码。

模板位于 `cmd/gen/templates/`。生成后按需补充校验、`UniqueFields`、关联表等业务逻辑，再执行 `make check`。

---

## 架构说明

```
┌─────────────────────────────────────────────────────────┐
│                  模块装配阶段（module.go）               │
│  1. 模块自行装配依赖（显式创建对象）                       │
│  2. NewCRUDModules 创建 []crud.Module（handler 构造）     │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                   Router 注册阶段                        │
│  1. 遍历 modules                                           │
│  2. 调用 ModuleConfig() 获取配置                           │
│  3. 注册权限到全局权限树                                    │
│  4. 反射注册路由到 Gin                                     │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                    请求处理                             │
│  JWT 认证 → 用户状态检查 → 权限检查 → Handler 方法        │
└─────────────────────────────────────────────────────────┘
```

---

## 最佳实践

1. **一个功能一个文件** - Handler 包含权限、路由、业务逻辑
2. **优先使用 CRUDHandler** - 用配置式 hook 代替重复的 CRUD 方法体
3. **使用 CRUDPerms** - 标准 CRUD 操作一行搞定
4. **复杂业务保持显式** - 超出 CRUD 的接口放在具体 handler 内
5. **批量操作避免循环 I/O** - 批量删除在 `BeforeDelete(tx, ids)` 中按 `ids` 统一清理，批量写入用 `CreateInBatches`
6. **模块显式装配** - 在 `NewCRUDModules` 维护 modules 列表，依赖清晰可控
7. **缓存失效放在提交之后** - 使用 `AfterCreateCommit`、`AfterUpdateCommit`、`AfterDeleteCommit`，不要在事务内失效缓存
