# 项目架构设计

## 🏗️ 整体架构

Bico Admin 采用前后端分离的架构设计，具有高度的模块化和可扩展性。

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   前端 (React)   │    │   后端 (Go)     │    │   数据库 (MySQL) │
│                 │    │                 │    │                 │
│  ┌─────────────┐ │    │  ┌─────────────┐ │    │  ┌─────────────┐ │
│  │ Refine UI   │ │◄──►│  │ Fiber API   │ │◄──►│  │ GORM ORM    │ │
│  └─────────────┘ │    │  └─────────────┘ │    │  └─────────────┘ │
│                 │    │                 │    │                 │
│  ┌─────────────┐ │    │  ┌─────────────┐ │    │  ┌─────────────┐ │
│  │ Ant Design  │ │    │  │ JWT Auth    │ │    │  │ 数据表结构   │ │
│  └─────────────┘ │    │  └─────────────┘ │    │  └─────────────┘ │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🎯 设计原则

### 1. 模块化设计
- **core**: 框架核心功能，可整体更新
- **admin**: 后台管理业务模块
- **api**: 对外API业务模块
- **business**: 通用业务方法封装

### 2. 分层架构
```
┌─────────────────────────────────────┐
│           Presentation Layer        │  ← Handler (控制器层)
├─────────────────────────────────────┤
│            Business Layer           │  ← Service (业务逻辑层)
├─────────────────────────────────────┤
│         Data Access Layer           │  ← Repository (数据访问层)
├─────────────────────────────────────┤
│            Data Layer               │  ← Model (数据模型层)
└─────────────────────────────────────┘
```

### 3. AI友好设计
- 清晰的目录结构和命名规范
- 完善的注释和文档
- 标准化的代码模式
- 类型安全的接口设计

## 🔧 后端架构

### 目录结构
```
backend/
├── cmd/server/           # 统一服务入口
├── core/                 # 框架核心
│   ├── config/          # 配置管理
│   ├── middleware/      # 中间件
│   ├── model/          # 基础模型
│   ├── repository/     # 数据访问
│   ├── service/        # 业务服务
│   ├── handler/        # 请求处理
│   └── router/         # 路由配置
├── modules/            # 业务模块目录
│   ├── admin/          # 后台管理模块
│   └── api/            # 对外API模块
├── business/           # 业务封装
└── pkg/                # 公共包
```

### 请求流程
```
HTTP Request
     ↓
Middleware (认证、日志、限流)
     ↓
Router (路由分发)
     ↓
Handler (请求处理)
     ↓
Service (业务逻辑)
     ↓
Repository (数据访问)
     ↓
Database (数据存储)
```

### 核心组件

#### 1. 配置管理
- 环境变量配置
- 数据库连接配置
- Redis缓存配置
- JWT认证配置

#### 2. 中间件系统
- 认证中间件 (JWT)
- 跨域中间件 (CORS)
- 日志中间件
- 限流中间件
- 错误处理中间件

#### 3. 业务封装
- 基础CRUD操作
- 分页查询封装
- 数据验证封装
- 响应格式统一

## 🎨 前端架构

### 技术栈
- **框架**: React 18 + TypeScript
- **UI库**: UmiJS + Ant Design Pro
- **状态管理**: UmiJS内置状态管理 + Model
- **路由**: UmiJS路由系统
- **构建工具**: UmiJS内置构建工具

### 组件架构
```
App
├── Layout (布局组件)
│   ├── Header (头部)
│   ├── Sidebar (侧边栏)
│   └── Content (内容区)
├── Pages (页面组件)
│   ├── Auth (认证页面)
│   ├── Dashboard (仪表板)
│   └── Management (管理页面)
└── Components (通用组件)
    ├── Form (表单组件)
    ├── Table (表格组件)
    └── Modal (弹窗组件)
```

## 🗄️ 数据库设计

### 核心表结构
```sql
-- 用户表
users
├── id (主键)
├── username (用户名)
├── email (邮箱)
├── password (密码)
├── status (状态)
└── timestamps

-- 角色表
roles
├── id (主键)
├── name (角色名)
├── code (角色代码)
├── description (描述)
└── timestamps

-- 权限表
permissions
├── id (主键)
├── name (权限名)
├── code (权限代码)
├── type (权限类型)
├── parent_id (父权限)
└── timestamps

-- 用户角色关联表
user_roles
├── user_id (用户ID)
└── role_id (角色ID)

-- 角色权限关联表
role_permissions
├── role_id (角色ID)
└── permission_id (权限ID)
```

## 🔐 安全架构

### 认证机制
- JWT Token 认证
- Token 自动刷新
- 登录状态管理

### 权限控制
- RBAC 权限模型
- 路由级权限控制
- 接口级权限验证
- 前端权限显示控制

### 安全防护
- SQL注入防护 (GORM)
- XSS攻击防护
- CSRF攻击防护
- 请求限流保护

## 📊 性能优化

### 后端优化
- 数据库连接池
- Redis缓存机制
- 分页查询优化
- 索引优化

### 前端优化
- 代码分割 (Code Splitting)
- 懒加载 (Lazy Loading)
- 缓存策略 (React Query)
- 构建优化 (Vite)

## 🔄 扩展性设计

### 水平扩展
- 无状态服务设计
- 负载均衡支持
- 微服务架构准备

### 垂直扩展
- 模块化插件系统
- 配置化权限管理
- 主题定制支持

## 📈 监控和日志

### 日志系统
- 结构化日志记录
- 日志级别管理
- 日志文件轮转

### 监控指标
- 接口响应时间
- 错误率统计
- 系统资源使用

## 🚀 部署架构

### 开发环境
```
Developer Machine
├── Frontend (Vite Dev Server)
├── Backend (Go Dev Server)
└── Database (Local MySQL)
```

### 生产环境
```
Load Balancer
├── Frontend (Nginx)
├── Backend (Go Binary)
├── Database (MySQL Cluster)
└── Cache (Redis Cluster)
```

## 🔮 未来规划

### 短期目标
- 完善权限管理系统
- 增加更多业务模块
- 优化性能和用户体验

### 长期目标
- 微服务架构演进
- 多租户支持
- 国际化支持
- 移动端适配
