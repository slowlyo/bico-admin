# CRUD 包使用文档

> 支持多模块分组（admin/api/自定义），自动注册权限和路由

> `internal/pkg/crud` - 声明式 CRUD 框架，让新功能开发只需一个文件。

## 概述

该包提供了一套声明式的 CRUD 框架，核心目标是**减少样板代码**。传统方式开发一个功能需要：

- `handler/xxx_handler.go`
- `service/xxx_service.go`
- `consts/permissions.go` (修改)
- `router.go` (修改)

使用本框架后，**只需一个 handler 文件**，路由、权限、DI 全部自动处理。

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

	h.BuildListQuery = func(db *gorm.DB, req *articleListReq) *gorm.DB {
		return db.Model(&model.Article{})
	}

	h.NewModelFromCreate = func(req *articleCreateReq) (*model.Article, error) {
		if req.Title == "" {
			return nil, errors.New("标题不能为空")
		}
		return &model.Article{Title: req.Title}, nil
	}

	h.BuildUpdates = func(req *articleUpdateReq, existing *model.Article) (map[string]interface{}, error) {
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

完成！接下来在模块入口（例如 `internal/admin/module.go`）把该模块加入模块列表：

```go
return []crud.Module{
    // ... 其他模块
    NewArticleHandler(db),
}
```

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
}
```

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
| `QueryList(c, query, &dest)` | 通用分页查询，自动处理 count/order/offset/limit |
| `crud.QueryListWithHook(&h.BaseHandler, c, query, &dest, after)` | 通用分页查询 + 结果二次处理（例如补充权限、统计字段） |
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
    crud.QueryListWithHook(&h.BaseHandler, c, query, &articles, func(items []model.Article) error {
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
perms.Routes()  // GET /, GET /:id, POST /, PUT /:id, DELETE /:id
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

### 常用 hook（按使用频率）

| hook | 说明 |
|------|------|
| `BuildListQuery(db, req)` | 构建列表查询（必填） |
| `BuildGetQuery(db)` | 构建详情查询（可选，常用于 `Preload`） |
| `NewModelFromCreate(req)` | Create 请求转模型 + 业务校验（必填） |
| `BuildUpdates(req, existing)` | Update 请求转 `updates map` + 业务校验（必填） |
| `CreateInTx(tx, item, req)` | 创建事务内扩展逻辑（写关联表等） |
| `UpdateInTx(tx, id, existing, req)` | 更新事务内扩展逻辑（同步关联表等） |
| `DeleteInTx(tx, id)` | 删除事务内扩展逻辑（清理关联表） |
| `AfterList(items)` | 列表返回前二次处理（回填权限、统计字段等） |
| `AfterGet(item)` | 详情返回前二次处理 |
| `ReloadAfterCreate(tx, id, item)` | 创建后重新加载（需要 preload 返回） |
| `ReloadAfterUpdate(tx, id, existing)` | 更新后重新加载（需要 preload 返回） |
| `DeleteBatchInTx(tx, ids)` | 批量删除事务内扩展逻辑（可选，避免循环 I/O） |

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
- 可通过 `DeleteBatchInTx` 统一清理关联表（避免循环中执行 I/O）

路由示例：

```go
crud.Route{Method: "DELETE", Path: "/batch", Handler: "DeleteBatch", Permission: perms.Delete}
```

### Exists（通用存在性判断）

用于唯一性校验：

```go
exists, err := crud.Exists(db, &model.AdminUser{}, "username = ?", req.Username)
if err != nil { return nil, err }
if exists { return nil, errors.New("用户名已存在") }
```

```go
package handler

import (
	"bico-admin/internal/admin/model"
	"bico-admin/internal/pkg/crud"
	"bico-admin/internal/pkg/password"
	"errors"

	"gorm.io/gorm"
)

var userPerms = crud.NewCRUDPerms("system", "admin_user", "用户管理")

type AdminUserHandler struct {
	crud.CRUDHandler[model.AdminUser, userListReq, createUserReq, updateUserReq]
}

func NewAdminUserHandler(db *gorm.DB) *AdminUserHandler {
	h := &AdminUserHandler{}
	h.DB = db
	h.NotFoundMsg = "用户不存在"

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
		exists, err := crud.Exists(db, &model.AdminUser{}, "username = ?", req.Username)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("用户名已存在")
		}

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

## 架构说明

```
┌─────────────────────────────────────────────────────────┐
│                  模块装配阶段（module.go）               │
│  1. 模块自行装配依赖（显式创建对象）                       │
│  2. 显式创建 []crud.Module（handler 构造）                │
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
5. **批量操作避免循环 I/O** - 批量删除用 `DeleteBatchInTx`，批量写入用 `CreateInBatches`
6. **模块显式装配** - 在模块入口维护 modules 列表，依赖清晰可控
