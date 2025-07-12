---
type: "always_apply"
---

# RBAC 权限系统设计规范

## 概述

本项目采用基于角色的访问控制（RBAC）系统，实现用户、角色、权限、菜单、API的完整关联管理。

## 核心设计理念

### 权限流转链路
```
用户(AdminUser) → 角色(AdminRole) → 权限(Permission) → 菜单/API/按钮
```

### 设计原则
- **权限和菜单在程序中静态定义** - 保证系统一致性和安全性
- **用户和角色在界面动态管理** - 提供灵活的权限配置
- **多维度权限控制** - 菜单级别、按钮级别、API级别
- **权限级别分层** - 查看、操作、管理、超级四个级别

## 数据模型设计

### 1. 管理员用户模型 (AdminUser)
```go
// 位置: internal/admin/models/admin_user.go
type AdminUser struct {
    types.BaseModel
    Username    string     // 用户名
    Password    string     // 密码
    Name        string     // 姓名
    Avatar      string     // 头像
    Email       string     // 邮箱
    Phone       string     // 手机号
    Status      int        // 状态：1-启用，0-禁用
    LastLoginAt *time.Time // 最后登录时间
    LastLoginIP string     // 最后登录IP
    LoginCount  int        // 登录次数
    Remark      string     // 备注

    // 关联关系
    Roles []AdminUserRole // 用户角色关联
}
```

### 2. 管理员角色模型 (AdminRole)
```go
// 位置: internal/admin/models/admin_role.go
type AdminRole struct {
    types.BaseModel
    Name        string // 角色名称
    Code        string // 角色代码（唯一）
    Description string // 角色描述
    Status      int    // 状态：1-启用，0-禁用

    // 关联关系
    Permissions []AdminRolePermission // 角色权限关联
    Users       []AdminUserRole       // 角色用户关联
}
```

### 3. 用户角色关联模型 (AdminUserRole)
```go
type AdminUserRole struct {
    ID        uint
    UserID    uint // 用户ID
    RoleID    uint // 角色ID
    CreatedAt time.Time

    // 关联关系
    User AdminUser
    Role AdminRole
}
```

### 4. 角色权限关联模型 (AdminRolePermission)
```go
type AdminRolePermission struct {
    ID             uint
    RoleID         uint   // 角色ID
    PermissionCode string // 权限代码
    CreatedAt      time.Time

    // 关联关系
    Role AdminRole
}
```

## 权限定义系统

### 权限结构 (Permission)
```go
// 位置: internal/admin/definitions/permissions.go
type Permission struct {
    Code        string   // 权限代码，如 "user:list"
    Name        string   // 权限名称，如 "查看用户列表"
    Module      string   // 所属模块，如 "user", "system"
    Description string   // 权限描述
    MenuSigns   []string // 关联的菜单标识列表
    Buttons     []string // 关联的按钮标识列表
    APIs        []string // 关联的API路径列表
    Level       int      // 权限级别：1-查看，2-操作，3-管理，4-超级
}
```

### 权限级别定义
```go
const (
    PermissionLevelView   = 1 // 查看权限 - 只能查看数据
    PermissionLevelAction = 2 // 操作权限 - 可以创建、编辑数据
    PermissionLevelManage = 3 // 管理权限 - 可以删除、管理数据
    PermissionLevelSuper  = 4 // 超级权限 - 系统级操作
)
```

### 权限代码规范
- 格式：`模块:操作`
- 示例：
  - `user:list` - 用户列表查看
  - `user:create` - 用户创建
  - `admin_user:delete` - 管理员删除
  - `system:config` - 系统配置

## 菜单定义系统

### 菜单结构 (MenuDefinition)
```go
// 位置: internal/admin/definitions/menus.go
type MenuDefinition struct {
    Sign        string           // 菜单唯一标识，如 "user.list"
    Name        string           // 菜单名称
    Path        string           // 路由路径
    Icon        string           // 图标
    Sort        int              // 排序
    ParentSign  string           // 父菜单标识，空表示顶级菜单
    Permissions []string         // 访问该菜单需要的权限
    Children    []MenuDefinition // 子菜单
}
```

### 菜单标识规范
- 顶级菜单：`user`, `admin_user`, `system`
- 子菜单：`user.list`, `user.detail`, `admin_user.create`
- 层级关系通过 `.` 分隔

## 权限与菜单关联机制

### 1. 静态关联
在权限定义中直接指定关联的菜单：
```go
{
    Code:      "user:list",
    Name:      "查看用户列表",
    MenuSigns: []string{"user", "user.list"},
    Buttons:   []string{"search", "filter", "export"},
    APIs:      []string{"/api/admin/users", "/api/admin/users/stats"},
}
```

### 2. 动态过滤
根据用户权限动态过滤菜单：
```go
func FilterMenusByPermissions(userPermissions []string) []types.Menu {
    // 1. 获取用户权限对应的所有菜单标识
    allowedMenuSigns := GetMenuSignsByPermissions(userPermissions)

    // 2. 递归过滤菜单树
    // 3. 返回用户可访问的菜单
}
```

## API权限控制

### 1. 权限中间件
```go
// 检查API访问权限
func CheckAPIPermission(userPermissions []string, apiPath string) bool {
    return HasPermissionForAPI(userPermissions, apiPath)
}
```

### 2. API路径匹配
- 支持精确匹配：`/api/admin/users`
- 支持参数路径：`/api/admin/users/:id`
- 支持通配符匹配（待实现）

## 按钮权限控制

### 1. 按钮权限检查
```go
func HasPermissionForButton(userPermissions []string, buttonKey string) bool {
    // 检查用户权限是否包含指定按钮权限
}
```

### 2. 前端按钮控制
```javascript
// 根据权限显示/隐藏按钮
if (hasPermission('user:create')) {
    showCreateButton();
}
```

## 预定义角色

### 角色代码常量
```go
const (
    RoleCodeSuperAdmin = "super_admin" // 超级管理员
    RoleCodeAdmin      = "admin"       // 管理员
    RoleCodeOperator   = "operator"    // 操作员
    RoleCodeViewer     = "viewer"      // 查看者
)
```

### 默认权限分配
- **超级管理员**: 拥有所有权限
- **管理员**: 拥有大部分管理权限，不包括系统级操作
- **操作员**: 拥有基本操作权限，不能删除数据
- **查看者**: 只有查看权限

## 服务层设计

### AdminRoleService 核心方法
```go
// 角色管理
GetRoleList(req *types.RoleListRequest) (*PageResponse[types.RoleResponse], error)
CreateRole(req *types.RoleCreateRequest) (*types.RoleResponse, error)
UpdateRole(id uint, req *types.RoleUpdateRequest) (*types.RoleResponse, error)
DeleteRole(id uint) error

// 权限管理
GetPermissionTree(roleID *uint) ([]types.PermissionTreeNode, error)
AssignRolesToUser(req *types.RoleAssignRequest) error

// 用户权限查询
GetUserRoles(userID uint) (*types.UserRoleResponse, error)
GetUserPermissions(userID uint) ([]string, error)
```

## 使用示例

### 1. 检查用户权限
```go
user := getCurrentUser()
if user.HasPermission("user:create") {
    // 允许创建用户
}
```

### 2. 过滤用户菜单
```go
userPermissions := user.GetPermissionCodes()
menus := FilterMenusByPermissions(userPermissions)
```

### 3. API权限验证
```go
if !HasPermissionForAPI(userPermissions, "/api/admin/users") {
    return errors.New("无权限访问")
}
```

## 数据库表结构

### 表命名规范
- `admin_users` - 管理员用户表
- `admin_roles` - 管理员角色表
- `admin_user_roles` - 用户角色关联表
- `admin_role_permissions` - 角色权限关联表

### 索引设计
- 用户表：`username` 唯一索引
- 角色表：`code` 唯一索引
- 关联表：复合索引 `(user_id, role_id)`, `(role_id, permission_code)`

## 安全考虑

1. **权限最小化原则** - 用户只获得必要的最小权限
2. **权限继承** - 通过角色继承权限，避免直接分配
3. **权限验证** - 每个API调用都需要验证权限
4. **审计日志** - 记录权限变更和敏感操作
5. **会话管理** - 权限变更后需要刷新用户会话

## 扩展性设计

1. **模块化权限** - 按业务模块组织权限
2. **动态权限** - 支持运行时权限检查
3. **权限缓存** - 缓存用户权限提高性能
4. **权限继承** - 支持权限的层级继承
5. **条件权限** - 支持基于条件的权限控制（如数据权限）
