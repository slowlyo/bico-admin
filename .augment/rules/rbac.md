---
type: "always_apply"
---

# RBAC 权限系统设计规范

## 概述

本项目采用基于角色的访问控制（RBAC）系统，实现用户、角色、权限、API的完整关联管理。菜单由前端根据权限动态控制显示。

## 核心设计理念

### 权限流转链路
```
用户(AdminUser) → 角色(AdminRole) → 权限(Permission) → API/按钮
前端菜单 ← 权限控制 ← 用户权限
```

### 设计原则
- **权限在程序中静态定义** - 保证系统一致性和安全性
- **菜单在前端静态定义** - 通过权限控制显示/隐藏
- **用户和角色在界面动态管理** - 提供灵活的权限配置
- **多维度权限控制** - 按钮级别、API级别
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
    Code    string   // 权限代码，如 "user:list"
    Name    string   // 权限名称，如 "查看用户列表"
    Module  string   // 所属模块，如 "user", "system"
    Buttons []string // 关联的按钮标识列表
    APIs    []string // 关联的API路径列表
    Level   int      // 权限级别：1-查看，2-操作，3-管理，4-超级
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

## 前端菜单控制机制

### 菜单权限控制原则
- **菜单在前端静态定义** - 保证路由和菜单的一致性
- **权限控制菜单显示** - 根据用户权限动态显示/隐藏菜单项
- **权限粒度控制** - 支持菜单级别和按钮级别的权限控制

### 前端菜单权限检查
```javascript
// 检查用户是否有权限访问某个菜单
function hasMenuPermission(requiredPermissions, userPermissions) {
    return requiredPermissions.some(permission =>
        userPermissions.includes(permission)
    );
}

// 过滤菜单项
function filterMenuItems(menuItems, userPermissions) {
    return menuItems.filter(item => {
        if (item.permissions && item.permissions.length > 0) {
            return hasMenuPermission(item.permissions, userPermissions);
        }
        return true; // 无权限要求的菜单项默认显示
    }).map(item => ({
        ...item,
        children: item.children ? filterMenuItems(item.children, userPermissions) : []
    }));
}
```

### 菜单权限配置示例
```javascript
// 前端菜单配置示例
const menuConfig = [
    {
        key: 'user',
        name: '用户管理',
        icon: 'UserOutlined',
        permissions: ['user:list'], // 需要的权限
        children: [
            {
                key: 'user-list',
                name: '用户列表',
                path: '/user/list',
                permissions: ['user:list']
            },
            {
                key: 'user-create',
                name: '创建用户',
                path: '/user/create',
                permissions: ['user:create']
            }
        ]
    }
];
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

// 用户权限查询（供前端使用）
GetUserRoles(userID uint) (*types.UserRoleResponse, error)
GetUserPermissions(userID uint) ([]string, error) // 返回权限代码列表
```

## 使用示例

### 1. 检查用户权限
```go
user := getCurrentUser()
if user.HasPermission("user:create") {
    // 允许创建用户
}
```

### 2. 获取用户权限（前端使用）
```go
// 后端只提供权限列表，不处理菜单
userPermissions := user.GetPermissionCodes()
// 返回给前端: ["user:list", "user:create", "admin_user:list"]
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

## 前后端分工

### 后端职责
1. **权限定义和管理** - 在代码中定义权限，提供权限管理API
2. **用户角色管理** - 管理用户和角色的关联关系
3. **API权限验证** - 验证每个API请求的权限
4. **权限数据提供** - 向前端提供用户权限列表

### 前端职责
1. **菜单定义和渲染** - 在前端定义菜单结构和路由
2. **菜单权限控制** - 根据用户权限动态显示/隐藏菜单
3. **按钮权限控制** - 根据权限控制页面按钮的显示
4. **路由权限守卫** - 在路由层面进行权限检查

## 扩展性设计

1. **模块化权限** - 按业务模块组织权限
2. **动态权限** - 支持运行时权限检查
3. **权限缓存** - 缓存用户权限提高性能
4. **权限继承** - 支持权限的层级继承
5. **条件权限** - 支持基于条件的权限控制（如数据权限）
6. **前端菜单缓存** - 前端缓存菜单配置和权限状态
