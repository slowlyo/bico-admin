# 开发流程指南

AI Agent 开发 CRUD 功能模块的标准化流程，通过引用现有代码示例快速上手。

## 开发准备

```bash
make deps  # 安装所有依赖和工具
```

## 后端开发流程

### 1. 数据模型 (Model)
**参考**: `internal/admin/models/admin_user.go`
- 定义结构体，包含 GORM 标签和 JSON 标签
- 实现 `TableName()` 方法
- 包含软删除字段 `DeletedAt`

### 2. 数据访问层 (Repository)
**参考**: `internal/admin/repository/admin_user.go`
- 定义接口：Create, GetByID, List, Update, Delete
- 实现结构体，注入 `*gorm.DB`
- 使用 context，支持分页和过滤

### 3. 业务逻辑层 (Service)
**参考**: `internal/admin/service/admin_user.go`
- 定义接口和请求/响应结构体
- 实现业务逻辑验证和数据转换
- 注入 Repository 依赖

### 4. HTTP 处理层 (Handler)
**参考**: `internal/admin/handler/admin_user.go`
- 实现 CRUD 方法，包含完整 Swagger 注释
- 参数验证和错误处理
- 使用统一响应格式

### 5. 路由注册
**参考**: `internal/admin/routes/routes.go`
- 在路由注册函数中添加路由组
- 使用 `middleware.Permission()` 权限中间件

### 6. 依赖注入
**参考**: `cmd/server/wire.go` 和 `internal/admin/provider.go`
- 在 `ProviderSet` 中添加 Provider
- 在 `Handlers` 结构体中添加字段
- 在 `NewHandlers` 中注入依赖

### 7. 数据库迁移配置
**参考**: `internal/admin/initializer/database.go`
- 在 `AutoMigrateAdminModels()` 中添加新模型
- 配置自动迁移和种子数据
- 应用启动时自动执行迁移

## 权限配置

### 1. 定义权限
**参考**: `internal/admin/definitions/permissions.go`
- 在 `GetPermissions()` 中添加权限定义
- 包含：list, create, update, delete 四个基本权限
- 配置对应的 API 路径

## 前端开发流程

### 1. API 服务
**参考**: `web/src/services/admin.ts`
- 定义接口类型和请求函数
- 使用统一的 request 方法

### 2. 页面组件
**参考**: `web/src/pages/AdminUser/index.tsx`
- 使用 ProTable 组件实现列表
- 集成搜索、分页、操作按钮
- 使用 `useAccess` 进行权限控制

### 3. 表单组件
**参考**: `web/src/pages/AdminUser/components/AdminUserForm.tsx`
- 使用 Modal + Form 实现弹窗表单
- 支持创建和编辑模式

### 4. 路由和菜单配置
**参考**: `web/.umirc.ts`
- 在 `routes` 数组中添加路由配置
- 设置权限检查和菜单显示

## 数据库迁移配置

### 自动迁移
**参考**: `internal/admin/initializer/database.go`
- 应用启动时自动执行数据库迁移
- 在 `AutoMigrateAdminModels()` 中添加新模型
- 自动创建默认超级管理员账户

### 手动迁移工具
**参考**: `cmd/migrate/main.go`
```bash
# 执行迁移
go run cmd/migrate/main.go -action=migrate

# 回滚迁移
go run cmd/migrate/main.go -action=rollback

# 重新创建数据库
go run cmd/migrate/main.go -action=fresh
```

### 种子数据
**参考**: `internal/admin/initializer/seeder.go`
- 创建默认超级管理员用户 (admin/admin123)
- 创建超级管理员角色
- 建立用户角色关联

## 开发工具

```bash
make wire     # 生成依赖注入代码

make dev      # 启动开发服务
```

## 检查清单

### 后端
- [ ] Model 定义完整
- [ ] Repository 实现 CRUD
- [ ] Service 包含业务逻辑
- [ ] Handler 包含 Swagger 注释
- [ ] 路由注册并使用权限中间件
- [ ] Wire 依赖注入配置
- [ ] 数据库迁移配置正确

### 权限
- [ ] 权限定义包含所有操作
- [ ] API 路径与权限匹配
- [ ] 菜单配置正确

### 前端
- [ ] API 服务定义完整
- [ ] 页面组件功能完整
- [ ] 权限控制正确实现
- [ ] 路由配置正确

### 测试
- [ ] 功能正常运行
- [ ] 权限控制生效
- [ ] API 文档生成正确

## 常见问题

### Wire 生成失败
- 检查 Provider 是否正确注册
- 确认导入路径正确
- 运行 `make wire` 重新生成

### API 调用失败
- 检查路由注册是否正确
- 确认权限配置是否匹配
- 查看服务器日志排查错误

### 前端权限不生效
- 检查 access.ts 中权限判断逻辑
- 确认权限代码与后端定义一致
- 验证用户是否拥有相应权限

## 开发步骤总结

1. **后端开发**: Model → Repository → Service → Handler → 路由 → 依赖注入 → 数据库迁移
2. **权限配置**: 定义权限 → 配置菜单
3. **前端开发**: API 服务 → 页面组件 → 表单组件 → 路由配置
4. **测试验证**: 功能测试 → 权限验证 → 文档生成

完成以上步骤后，一个完整的 CRUD 功能模块就开发完成了！
