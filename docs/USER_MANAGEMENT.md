# 用户管理功能文档

## 概述

用户管理功能是 Bico Admin 系统的核心功能之一，提供完整的用户生命周期管理，包括用户的创建、查询、更新、删除以及状态管理等功能。

## 功能特性

### 🎯 核心功能

1. **用户列表管理**
   - 分页查询用户列表
   - 支持用户名、邮箱、昵称搜索
   - 实时状态显示和切换
   - 批量操作支持

2. **用户信息管理**
   - 创建新用户
   - 查看用户详情
   - 更新用户资料
   - 删除用户（软删除）

3. **用户状态管理**
   - 启用/禁用用户
   - 实时状态切换
   - 状态变更记录

4. **密码管理**
   - 管理员重置用户密码
   - 用户修改密码（需验证旧密码）
   - 密码加密存储

5. **批量操作**
   - 批量删除用户
   - 批量状态管理（计划中）

## 技术架构

### 后端架构

采用简化的三层架构设计：

```
┌─────────────────┐
│   UserHandler   │  ← HTTP 请求处理，参数验证，响应格式化
├─────────────────┤
│ Database.Ops[T] │  ← 通用数据库操作，业务逻辑
├─────────────────┤
│   User Model    │  ← 数据模型，数据验证，关系定义
├─────────────────┤
│   Database      │  ← 数据存储
└─────────────────┘
```

### 前端架构

基于 UmiJS + Ant Design Pro：

```
┌─────────────────┐
│   Users Page    │  ← 用户管理主页面
├─────────────────┤
│   Components    │  ← 创建/编辑/详情组件
├─────────────────┤
│   Services      │  ← API 调用服务
├─────────────────┤
│   Backend API   │  ← 后端接口
└─────────────────┘
```

## API 接口

### 用户管理接口

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/admin/users` | 获取用户列表 |
| POST | `/admin/users` | 创建用户 |
| GET | `/admin/users/:id` | 获取单个用户 |
| PUT | `/admin/users/:id` | 更新用户 |
| DELETE | `/admin/users/:id` | 删除用户 |

### 批量操作接口

| 方法 | 路径 | 描述 |
|------|------|------|
| DELETE | `/admin/users/batch` | 批量删除用户 |

### 状态管理接口

| 方法 | 路径 | 描述 |
|------|------|------|
| PUT | `/admin/users/:id/status` | 更新用户状态 |

### 密码管理接口

| 方法 | 路径 | 描述 |
|------|------|------|
| PUT | `/admin/users/:id/password` | 修改用户密码 |
| PUT | `/admin/users/:id/reset-password` | 重置用户密码 |

## 数据模型

### User 模型

```go
type User struct {
    BaseModel
    Username    string     `json:"username" gorm:"uniqueIndex;size:50;not null"`
    Email       string     `json:"email" gorm:"uniqueIndex;size:100;not null"`
    Password    string     `json:"-" gorm:"size:255;not null"`
    Nickname    string     `json:"nickname" gorm:"size:50"`
    Avatar      string     `json:"avatar" gorm:"size:255"`
    Phone       string     `json:"phone" gorm:"size:20"`
    Status      UserStatus `json:"status" gorm:"default:1"`
    LastLoginAt *time.Time `json:"last_login_at"`
    LastLoginIP string     `json:"last_login_ip" gorm:"size:45"`
    Roles       []Role     `json:"roles" gorm:"many2many:user_roles;"`
}
```

### 用户状态

```go
const (
    UserStatusInactive UserStatus = 0 // 未激活
    UserStatusActive   UserStatus = 1 // 激活
    UserStatusBlocked  UserStatus = 2 // 被封禁
)
```

## 前端组件

### 主要页面

1. **用户列表页面** (`/src/pages/Users/index.tsx`)
   - ProTable 表格组件
   - 搜索和过滤功能
   - 批量操作工具栏
   - 状态切换开关

2. **创建用户表单** (`/src/pages/Users/components/CreateForm.tsx`)
   - ModalForm 弹窗表单
   - 表单验证
   - 角色选择（计划中）

3. **编辑用户表单** (`/src/pages/Users/components/UpdateForm.tsx`)
   - 用户信息编辑
   - 状态管理
   - 角色管理（计划中）

4. **用户详情** (`/src/pages/Users/components/UserDetail.tsx`)
   - 用户信息展示
   - 操作历史（计划中）

## 安全特性

### 密码安全

1. **密码加密**
   - 使用 bcrypt 算法加密存储
   - 密码不在 JSON 响应中返回

2. **密码策略**
   - 最小长度 6 位
   - 支持复杂度验证（可扩展）

### 权限控制

1. **认证中间件**
   - JWT Token 验证
   - 路由级别权限控制

2. **操作权限**
   - 管理员权限验证
   - 操作日志记录（计划中）

## 使用指南

### 创建用户

1. 点击"新建"按钮
2. 填写用户基本信息
3. 设置初始密码
4. 选择用户状态
5. 提交创建

### 管理用户状态

1. 在用户列表中找到目标用户
2. 点击状态开关进行切换
3. 系统自动保存状态变更

### 重置密码

1. 进入用户详情页面
2. 点击"重置密码"按钮
3. 输入新密码
4. 确认重置操作

### 批量操作

1. 在用户列表中选择多个用户
2. 使用底部工具栏进行批量操作
3. 确认操作后执行

## 开发指南

### 添加新字段

1. **后端**：
   - 更新 User 模型
   - 更新请求/响应结构
   - 添加数据库迁移

2. **前端**：
   - 更新 TypeScript 类型定义
   - 更新表单组件
   - 更新列表显示

### 扩展功能

1. **角色管理集成**
   - 实现角色分配逻辑
   - 更新用户创建/编辑表单
   - 添加角色权限验证

2. **操作日志**
   - 添加日志记录中间件
   - 创建日志查看页面
   - 实现日志搜索功能

## 故障排除

### 常见问题

1. **用户创建失败**
   - 检查用户名/邮箱是否重复
   - 验证密码复杂度
   - 检查网络连接

2. **状态切换失败**
   - 确认用户权限
   - 检查后端服务状态
   - 查看浏览器控制台错误

3. **密码重置失败**
   - 验证新密码格式
   - 检查用户是否存在
   - 确认操作权限

### 调试方法

1. **后端调试**
   - 查看服务器日志
   - 使用 API 测试工具
   - 检查数据库连接

2. **前端调试**
   - 打开浏览器开发者工具
   - 查看网络请求
   - 检查控制台错误信息

## 性能优化

### 后端优化

1. **数据库查询优化**
   - 使用索引优化查询
   - 实现分页查询
   - 预加载关联数据

2. **缓存策略**
   - 用户信息缓存
   - 权限信息缓存
   - 查询结果缓存

### 前端优化

1. **组件优化**
   - 使用 React.memo 优化渲染
   - 实现虚拟滚动
   - 懒加载组件

2. **数据管理**
   - 实现本地状态管理
   - 优化 API 调用
   - 缓存查询结果

## 未来规划

### 短期目标

1. 完善角色管理集成
2. 添加操作日志功能
3. 实现用户导入/导出
4. 优化批量操作性能

### 长期目标

1. 实现用户组管理
2. 添加用户行为分析
3. 集成第三方认证
4. 实现多租户支持
