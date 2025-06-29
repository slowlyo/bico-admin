# 权限系统文档

## 概述

本系统采用基于代码定义的权限管理方式，权限代码在后端代码中定义，通过数据库存储角色与权限的关联关系。

## 权限代码定义

权限代码采用 `模块:操作` 的格式，例如 `user:view`、`role:create` 等。

### 系统管理权限 (system:*)

| 权限代码 | 权限名称 | 描述 | 分类 |
|---------|---------|------|------|
| `system:view` | 查看系统 | 查看系统管理页面 | 系统管理 |
| `system:manage` | 管理系统 | 管理系统设置 | 系统管理 |

### 用户管理权限 (user:*)

| 权限代码 | 权限名称 | 描述 | 分类 |
|---------|---------|------|------|
| `user:view` | 查看用户 | 查看用户列表和详情 | 用户管理 |
| `user:create` | 创建用户 | 创建新用户 | 用户管理 |
| `user:update` | 编辑用户 | 编辑用户信息 | 用户管理 |
| `user:delete` | 删除用户 | 删除用户 | 用户管理 |
| `user:manage_status` | 管理用户状态 | 启用/禁用用户 | 用户管理 |
| `user:reset_password` | 重置密码 | 重置用户密码 | 用户管理 |

### 角色管理权限 (role:*)

| 权限代码 | 权限名称 | 描述 | 分类 |
|---------|---------|------|------|
| `role:view` | 查看角色 | 查看角色列表和详情 | 角色管理 |
| `role:create` | 创建角色 | 创建新角色 | 角色管理 |
| `role:update` | 编辑角色 | 编辑角色信息 | 角色管理 |
| `role:delete` | 删除角色 | 删除角色 | 角色管理 |
| `role:assign_permissions` | 分配权限 | 为角色分配权限 | 角色管理 |

## 权限系统架构

### 后端权限定义

权限在后端代码中定义，位置：`backend/core/permission/config.go`

```go
// 权限常量定义
const (
    // 系统管理权限
    SystemView   = "system:view"
    SystemManage = "system:manage"

    // 用户管理权限
    UserView          = "user:view"
    UserCreate        = "user:create"
    UserUpdate        = "user:update"
    UserDelete        = "user:delete"
    UserManageStatus  = "user:manage_status"
    UserResetPassword = "user:reset_password"

    // 角色管理权限
    RoleView              = "role:view"
    RoleCreate            = "role:create"
    RoleUpdate            = "role:update"
    RoleDelete            = "role:delete"
    RoleAssignPermissions = "role:assign_permissions"
)

// AllPermissions 所有权限定义
var AllPermissions = []Permission{
    // 系统管理
    {Code: SystemView, Name: "查看系统", Description: "查看系统管理页面", Category: "系统管理"},
    {Code: SystemManage, Name: "管理系统", Description: "管理系统设置", Category: "系统管理"},
    
    // 用户管理
    {Code: UserView, Name: "查看用户", Description: "查看用户列表和详情", Category: "用户管理"},
    {Code: UserCreate, Name: "创建用户", Description: "创建新用户", Category: "用户管理"},
    {Code: UserUpdate, Name: "编辑用户", Description: "编辑用户信息", Category: "用户管理"},
    {Code: UserDelete, Name: "删除用户", Description: "删除用户", Category: "用户管理"},
    {Code: UserManageStatus, Name: "管理用户状态", Description: "启用/禁用用户", Category: "用户管理"},
    {Code: UserResetPassword, Name: "重置密码", Description: "重置用户密码", Category: "用户管理"},
    
    // 角色管理
    {Code: RoleView, Name: "查看角色", Description: "查看角色列表和详情", Category: "角色管理"},
    {Code: RoleCreate, Name: "创建角色", Description: "创建新角色", Category: "角色管理"},
    {Code: RoleUpdate, Name: "编辑角色", Description: "编辑角色信息", Category: "角色管理"},
    {Code: RoleDelete, Name: "删除角色", Description: "删除角色", Category: "角色管理"},
    {Code: RoleAssignPermissions, Name: "分配权限", Description: "为角色分配权限", Category: "角色管理"},
}
```

### 前端权限定义

前端权限常量定义，位置：`frontend/src/constants/permissions.ts`

```typescript
export const PERMISSIONS = {
  // 系统管理权限
  SYSTEM: {
    VIEW: 'system:view',
    MANAGE: 'system:manage',
  },

  // 用户管理权限
  USER: {
    VIEW: 'user:view',
    CREATE: 'user:create',
    UPDATE: 'user:update',
    DELETE: 'user:delete',
    MANAGE_STATUS: 'user:manage_status',
    RESET_PASSWORD: 'user:reset_password',
  },

  // 角色管理权限
  ROLE: {
    VIEW: 'role:view',
    CREATE: 'role:create',
    UPDATE: 'role:update',
    DELETE: 'role:delete',
    ASSIGN_PERMISSIONS: 'role:assign_permissions',
  },
};
```

## 数据库结构

### 角色权限关联表 (role_permissions)

```sql
CREATE TABLE role_permissions (
    role_id INT NOT NULL,
    permission_code VARCHAR(100) NOT NULL,
    PRIMARY KEY (role_id, permission_code),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);
```

### 用户角色关联表 (user_roles)

```sql
CREATE TABLE user_roles (
    user_id INT NOT NULL,
    role_id INT NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);
```

## 权限验证流程

### 1. 用户登录后获取权限

```go
// 获取用户权限的API: GET /admin/auth/permissions
func (h *AuthHandler) GetUserPermissions(c *fiber.Ctx) error {
    userID := middleware.GetUserID(c)
    
    // 检查是否为超级管理员
    isSuperAdmin, err := h.isSuperAdmin(userID)
    if isSuperAdmin {
        // 返回所有权限代码
        allPermissionCodes := []string{
            "system:view", "system:manage",
            "user:view", "user:create", "user:update", "user:delete", 
            "user:manage_status", "user:reset_password",
            "role:view", "role:create", "role:update", "role:delete", 
            "role:assign_permissions",
        }
        return response.Success(c, allPermissionCodes)
    }
    
    // 普通用户通过数据库查询获取权限
    // ...
}
```

### 2. 权限中间件验证

```go
// 权限验证中间件
func (pm *PermissionMiddleware) RequirePermission(permission string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID := middleware.GetUserID(c)
        
        // 检查权限
        hasPermission, err := pm.hasPermission(userID, permission)
        if !hasPermission {
            return response.Forbidden(c, "Insufficient permissions")
        }
        
        return c.Next()
    }
}
```

### 3. 前端权限检查

```typescript
// 前端权限检查函数
export function hasPermission(userPermissions: string[], requiredPermission: string): boolean {
  return userPermissions.includes(requiredPermission);
}

// 在access.ts中使用
export default (initialState: any) => {
  const { userPermissions } = initialState ?? {};
  const permissions = userPermissions || [];

  return {
    canSeeAdmin: hasPermission(permissions, PERMISSIONS.SYSTEM.VIEW),
    canManageUsers: hasPermission(permissions, PERMISSIONS.USER.VIEW),
    canManageRoles: hasPermission(permissions, PERMISSIONS.ROLE.VIEW),
    // ...
  };
};
```

## 超级管理员

### 特殊处理

- **角色标识**: `super_admin`
- **权限**: 拥有所有权限，无需数据库查询
- **保护机制**: 不能删除、不能修改权限

### 检查逻辑

```go
func (pm *PermissionMiddleware) isSuperAdmin(userID uint) (bool, error) {
    var count int64
    err := pm.db.Table("users u").
        Joins("JOIN user_roles ur ON u.id = ur.user_id").
        Joins("JOIN roles r ON ur.role_id = r.id").
        Where("u.id = ? AND r.code = ? AND r.status = ?",
            userID, "super_admin", model.RoleStatusActive).
        Count(&count).Error

    if err != nil {
        // 兼容旧版本：检查用户表的role字段
        var user model.User
        if err := pm.db.First(&user, userID).Error; err != nil {
            return false, err
        }
        return user.Role == "super_admin", nil
    }

    return count > 0, nil
}
```

## 个人资料权限豁免

个人资料相关功能对所有登录用户免权限验证：

- 查看个人资料
- 修改个人资料  
- 修改密码

这些功能不需要特定权限，所有登录用户都可以访问。

## 权限分配

### API接口

```
PUT /admin/roles/:id/permissions
```

### 请求格式

```json
{
  "permission_codes": [
    "user:view",
    "user:create",
    "role:view"
  ]
}
```

### 验证逻辑

1. 验证角色是否存在
2. 检查是否为受保护角色（超级管理员）
3. 验证权限代码是否有效
4. 清除现有权限关联
5. 创建新的权限关联

## 添加新权限

### 1. 后端添加权限常量

在 `backend/core/permission/config.go` 中：

```go
const (
    // 添加新的权限常量
    NewModuleView = "new_module:view"
)

// 在 AllPermissions 中添加权限定义
var AllPermissions = []Permission{
    // ...
    {Code: NewModuleView, Name: "查看新模块", Description: "查看新模块页面", Category: "新模块管理"},
}
```

### 2. 前端添加权限常量

在 `frontend/src/constants/permissions.ts` 中：

```typescript
export const PERMISSIONS = {
  // ...
  NEW_MODULE: {
    VIEW: 'new_module:view',
  },
};
```

### 3. 更新超级管理员权限列表

在 `backend/core/handler/auth.go` 的 `GetUserPermissions` 方法中添加新权限代码。

## 注意事项

1. **权限代码格式**: 必须遵循 `模块:操作` 的格式
2. **权限验证**: 所有API都应该添加适当的权限验证中间件
3. **前后端同步**: 前后端的权限常量定义必须保持一致
4. **超级管理员**: 自动拥有所有权限，添加新权限时需要更新权限列表
5. **个人资料豁免**: 个人资料相关功能无需权限验证

## 相关文件

- 后端权限配置: `backend/core/permission/config.go`
- 前端权限常量: `frontend/src/constants/permissions.ts`
- 权限中间件: `backend/core/permission/middleware.go`
- 角色权限处理器: `backend/modules/admin/handler/role_permission.go`
- 前端权限检查: `frontend/src/access.ts`
- RBAC系统文档: `docs/RBAC_SYSTEM.md`
