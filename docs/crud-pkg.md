# CRUD 包文档

> `internal/pkg/crud` - 声明式 CRUD 框架，让新功能开发只需一个文件。

## 概述

该包提供了一套声明式的 CRUD 框架，核心目标是**减少样板代码**。传统方式开发一个功能需要：

- `handler/xxx_handler.go`
- `service/xxx_service.go`
- `consts/permissions.go` (修改)
- `router.go` (修改)
- `container.go` (修改)

使用本框架后，**只需一个 handler 文件**，路由、权限、DI 全部自动处理。

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
var articlePerms = crud.NewCRUDPerms("article", "文章管理")

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

// 5. 自动注册
func init() {
    crud.RegisterModule(NewArticleHandler)
}

var _ crud.Module = (*ArticleHandler)(nil)
```

完成！不需要修改 `router.go`、`container.go` 或任何其他文件。

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
var perms = crud.NewCRUDPerms("article", "文章管理")

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
var perms = crud.NewCRUDPerms("article", "文章管理").WithExtra(
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

// 获取所有权限
crud.GetAllPermissions()     // []Permission
crud.GetAllPermissionKeys()  // []string
```

---

## 完整示例

参考 `internal/admin/handler/admin_user_handler.go`：

```go
package handler

import (
    "bico-admin/internal/admin/model"
    "bico-admin/internal/pkg/crud"
    "bico-admin/internal/pkg/password"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

var userPerms = crud.NewCRUDPerms("admin_user", "用户管理")

type AdminUserHandler struct {
    crud.BaseHandler
    db *gorm.DB
}

func NewAdminUserHandler(db *gorm.DB) *AdminUserHandler {
    return &AdminUserHandler{db: db}
}

func (h *AdminUserHandler) ModuleConfig() crud.ModuleConfig {
    return crud.ModuleConfig{
        Name:             "admin_user",
        Group:            "/admin-users",
        ParentPermission: PermSystemManage,
        Permissions:      userPerms.Tree,
        Routes:           userPerms.Routes(),
    }
}

type (
    userListReq   struct { Username, Name string; Enabled *bool }
    createUserReq struct { Username, Password, Name, Avatar string; Enabled *bool; RoleIDs []uint }
    updateUserReq struct { Name, Avatar string; Enabled *bool; RoleIDs []uint }
)

func (h *AdminUserHandler) List(c *gin.Context) {
    var req userListReq
    h.BindQuery(c, &req)

    query := h.db.Model(&model.AdminUser{}).Preload("Roles")
    if req.Username != "" {
        query = query.Where("username LIKE ?", "%"+req.Username+"%")
    }
    if req.Name != "" {
        query = query.Where("name LIKE ?", "%"+req.Name+"%")
    }
    if req.Enabled != nil {
        query = query.Where("enabled = ?", *req.Enabled)
    }

    var users []model.AdminUser
    h.QueryList(c, query, &users)
}

func (h *AdminUserHandler) Get(c *gin.Context) {
    id, err := h.ParseID(c)
    if err != nil {
        return
    }
    var user model.AdminUser
    if h.QueryOne(c, h.db.Preload("Roles").Where("id = ?", id), &user, "用户不存在") {
        h.Success(c, user)
    }
}

func (h *AdminUserHandler) Create(c *gin.Context) {
    var req createUserReq
    if err := h.BindJSON(c, &req); err != nil {
        return
    }
    // ... 业务逻辑
    h.ExecTx(c, h.db, func(tx *gorm.DB) error {
        // 创建用户
        return nil
    }, "创建成功", user)
}

func (h *AdminUserHandler) Update(c *gin.Context) {
    // ... 类似逻辑
}

func (h *AdminUserHandler) Delete(c *gin.Context) {
    id, _ := h.ParseID(c)
    h.ExecTx(c, h.db, func(tx *gorm.DB) error {
        return tx.Delete(&model.AdminUser{}, id).Error
    }, "删除成功", nil)
}

func init() {
    crud.RegisterModule(NewAdminUserHandler)
}

var _ crud.Module = (*AdminUserHandler)(nil)
```

---

## 架构说明

```
┌─────────────────────────────────────────────────────────┐
│                    init() 阶段                          │
│  crud.RegisterModule(NewXxxHandler) 注册构造函数         │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                   Router 初始化                         │
│  1. 遍历所有注册的模块                                   │
│  2. 反射调用构造函数，注入 *gorm.DB                       │
│  3. 调用 ModuleConfig() 获取配置                         │
│  4. 注册权限到全局权限树                                  │
│  5. 反射注册路由到 Gin                                   │
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
2. **使用 CRUDPerms** - 标准 CRUD 操作一行搞定
3. **嵌入 BaseHandler** - 减少重复代码
4. **使用 QueryList/QueryOne** - 统一查询逻辑
5. **使用 ExecTx** - 统一事务处理
6. **init() 注册** - 自动发现，无需修改其他文件
