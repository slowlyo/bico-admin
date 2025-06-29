# RBAC权限系统设计文档

## 概述

Bico Admin 采用基于角色的访问控制（RBAC）系统，实现了灵活的权限管理机制。系统遵循"权限标识在代码中定义，角色关联在数据库中管理"的设计原则，既保证了权限定义的一致性，又提供了动态权限分配的灵活性。

## 🎯 设计原则

### 1. 权限标识代码化
- **权限常量**：在 `backend/core/permission/config.go` 中定义
- **统一引用**：前后端都使用相同的权限标识
- **类型安全**：通过常量避免权限标识拼写错误

### 2. 角色关联数据库化
- **动态分配**：角色与权限的关联存储在 `role_permissions` 表中
- **界面管理**：管理员可以通过系统界面为角色分配权限
- **实时生效**：权限变更立即生效，无需重启系统

### 3. 简化架构原则
- **无Service层**：遵循项目简化架构，业务逻辑直接在Handler中实现
- **通用工具**：使用 `database.Operations[T]` 进行数据库操作
- **AI友好**：清晰的代码结构，便于AI理解和维护

## 🏗️ 系统架构

### 权限流程图
```
用户登录 → 获取用户角色 → 查询角色权限 → 权限验证 → 访问控制
    ↓           ↓            ↓           ↓          ↓
  JWT Token   user_roles   role_permissions  中间件    API/页面
```

### 数据库设计
```sql
-- 用户表
users (id, username, email, role, ...)

-- 角色表
roles (id, name, code, description, status, ...)

-- 用户角色关联表
user_roles (user_id, role_id)

-- 角色权限关联表（直接存储权限代码）
role_permissions (role_id, permission_code)

-- 注意：不需要permissions表，权限定义在代码中
```

## 📋 权限定义

### 权限常量结构
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
```

### 权限分类
- **系统管理**：系统设置、配置管理等
- **用户管理**：用户增删改查、状态管理等
- **角色管理**：角色增删改查、权限分配等

## 🔧 核心组件

### 1. 权限中间件 (`PermissionMiddleware`)
```go
// 权限验证中间件
func (pm *PermissionMiddleware) RequirePermission(permission string) fiber.Handler

// 任意权限验证
func (pm *PermissionMiddleware) RequireAnyPermission(permissions ...string) fiber.Handler

// 全部权限验证  
func (pm *PermissionMiddleware) RequireAllPermissions(permissions ...string) fiber.Handler
```

### 2. 角色权限管理器 (`RolePermissionHandler`)
```go
// 获取所有权限定义
func (h *RolePermissionHandler) GetAllPermissions(c *fiber.Ctx) error

// 获取角色权限
func (h *RolePermissionHandler) GetRolePermissions(c *fiber.Ctx) error

// 分配角色权限
func (h *RolePermissionHandler) AssignRolePermissions(c *fiber.Ctx) error

// 移除角色权限
func (h *RolePermissionHandler) RemoveRolePermission(c *fiber.Ctx) error
```

## 🔐 特殊权限处理

### 1. 超级管理员
- **角色标识**：`super_admin`
- **权限范围**：拥有所有权限
- **保护机制**：不可删除、不可编辑权限
- **权限检查**：直接返回true，无需查询数据库

### 2. 个人资料功能
- **权限豁免**：所有登录用户都可以访问个人资料功能
- **豁免路由**：`/auth/profile`、`/auth/change-password`
- **豁免权限**：`profile:*` 相关权限已移除

## 📡 API接口

### 权限管理接口
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/admin/permissions` | 获取所有权限定义 |
| GET | `/admin/roles/:id/permissions` | 获取角色权限 |
| PUT | `/admin/roles/:id/permissions` | 分配角色权限 |
| DELETE | `/admin/roles/:roleId/permissions/:permissionCode` | 移除角色权限 |

### 用户权限接口
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/auth/permissions` | 获取当前用户权限 |

## 🎨 前端集成

### 权限检查函数
```typescript
// 检查单个权限
const hasPermission = (userPermissions: string[], permission: string): boolean

// 检查多个权限（全部满足）
const hasAllPermissions = (userPermissions: string[], permissions: string[]): boolean

// 检查多个权限（任意满足）
const hasAnyPermission = (userPermissions: string[], permissions: string[]): boolean
```

### 使用示例
```typescript
// 获取用户权限
const userPermissions = await getUserPermissions();

// 权限检查
if (hasPermission(userPermissions, 'user:create')) {
  // 显示创建用户按钮
}
```

## 🔄 权限分配流程

### 管理员操作流程
1. **登录系统** → 使用超级管理员账号
2. **角色管理** → 创建或编辑角色
3. **权限分配** → 为角色选择权限
4. **用户分配** → 为用户分配角色
5. **权限生效** → 用户重新登录后权限生效

### 系统处理流程
1. **权限选择** → 前端展示所有可用权限（从代码获取）
2. **提交请求** → 发送权限代码列表到后端
3. **权限验证** → 验证权限代码在代码中是否存在
4. **数据库更新** → 直接在 `role_permissions` 表中存储权限代码
5. **权限生效** → 下次权限检查时生效

## 🛡️ 安全特性

### 1. 权限验证
- **中间件保护**：所有需要权限的路由都使用权限中间件
- **数据库查询**：实时查询用户权限，确保权限变更及时生效
- **超级管理员检查**：优先检查超级管理员身份

### 2. 角色保护
- **受保护角色**：超级管理员角色不可删除或修改权限
- **权限验证**：只有有权限的用户才能分配权限
- **操作日志**：权限变更操作记录（计划中）

## 🚀 使用指南

### 开发者指南
1. **添加新权限**：在 `permission/config.go` 中添加权限常量
2. **使用权限**：在路由中使用权限中间件
3. **前端集成**：使用权限检查函数控制UI显示

### 管理员指南
1. **创建角色**：在角色管理页面创建新角色
2. **分配权限**：为角色选择合适的权限
3. **用户管理**：为用户分配合适的角色

## 📈 扩展性

### 1. 权限扩展
- **新增权限**：在代码中添加新的权限常量
- **权限分类**：按功能模块组织权限
- **权限继承**：支持权限层级结构（计划中）

### 2. 功能扩展
- **权限缓存**：Redis缓存用户权限信息
- **操作日志**：记录权限变更操作
- **权限审计**：权限使用情况统计

## 🎉 总结

新的RBAC权限系统具有以下优势：

✅ **架构清晰**：权限定义与角色关联分离，职责明确
✅ **开发友好**：权限常量化，避免硬编码字符串
✅ **管理灵活**：支持动态权限分配，无需重启系统
✅ **安全可靠**：多层权限验证，超级管理员保护
✅ **易于维护**：遵循简化架构，代码结构清晰
✅ **扩展性强**：支持权限扩展和功能增强
✅ **数据库简化**：不存储权限数据，只存储权限代码关联

这个设计完美符合您的需求：权限标识在代码中定义，角色关联在数据库中管理，用户可以在界面上灵活分配权限，同时避免了不必要的权限数据存储。
