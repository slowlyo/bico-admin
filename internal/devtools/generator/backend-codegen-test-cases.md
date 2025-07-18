# 后端代码生成器功能实现文档

## 概述

本文档基于现有代码实现，详细描述 Bico Admin 后端代码生成器的实际功能和生成效果。

## 功能配置

- **服务地址**: `http://127.0.0.1:18901`
- **工具名称**: `generate_code`
- **支持组件**: `model`, `repository`, `service`, `handler`, `routes`, `wire`, `migration`, `permission`, `all`
- **默认包路径**: `internal/admin`

## 1. 模型生成功能

### 1.1 模型文件生成

**功能描述**: 生成 GORM 数据模型文件

**实现效果**:
- 生成位置：`internal/shared/models/{model_name}.go`
- 继承 `types.BaseModel`（包含 ID、CreatedAt、UpdatedAt）
- 自动生成 `TableName()` 方法，支持自定义表名
- 智能字段类型映射：
  - `string` → `string`
  - `int/int32/int64` → 对应 Go 类型
  - `decimal` → `float64`
  - `time/date/datetime/timestamp` → `*time.Time`
  - `bool` → `bool`
  - `text/json` → `string`

### 1.2 字段处理能力

**功能描述**: 处理各种字段定义和标签

**实现效果**:
- 自动处理 GORM 标签（索引、约束、默认值等）
- 自动生成 JSON 标签（支持自定义或蛇形命名）
- 字段名称清理（处理 Go 关键字冲突）
- 自动导入必要的包（time 包、shared/types 包）
- 支持字段注释生成

### 1.3 表名处理

**功能描述**: 智能处理数据库表名

**实现效果**:
- 默认使用模型名的蛇形复数形式（如 User → users）
- 支持自定义表名覆盖
- 生成正确的 TableName() 方法

## 2. Repository 数据访问层

### 2.1 Repository 接口和实现

**功能描述**: 生成基于泛型的数据访问层

**实现效果**:
- 生成位置：`internal/admin/repository/{model_name}_repository.go`
- 继承 `repository.BaseRepositoryInterface[models.{ModelName}]`
- 自动生成接口定义和实现结构体
- 提供 `New{ModelName}Repository` 构造函数
- 集成 GORM 数据库操作

### 2.2 查询功能

**功能描述**: 提供标准的数据查询能力

**实现效果**:
- 继承基础仓储的所有 CRUD 方法
- 实现 `ListWithFilter` 方法支持分页和过滤
- 支持 `types.BasePageQuery` 分页查询
- 自动处理关键词搜索和排序
- 基于泛型实现类型安全

## 3. Service 业务逻辑层

### 3.1 Service 接口和实现

**功能描述**: 生成业务逻辑服务层

**实现效果**:
- 生成位置：`internal/admin/service/{model_name}_service.go`
- 继承 `service.BaseServiceInterface[models.{ModelName}, repository.{ModelName}Repository]`
- 自动生成接口定义和实现结构体
- 提供 `New{ModelName}Service` 构造函数
- 依赖注入对应的 Repository

### 3.2 业务方法实现

**功能描述**: 提供完整的业务操作方法

**实现效果**:
- 实现 `ListWithFilter` 分页查询方法
- 重写 `Create/Update/Delete` 方法添加业务验证
- 实现 `UpdateStatus` 状态管理方法
- 提供参数验证（分页参数、ID 验证、状态验证）
- 统一错误处理和响应格式

### 3.3 验证框架

**功能描述**: 内置业务验证框架

**实现效果**:
- 生成 `validate{ModelName}` 实体验证方法
- 生成 `validateDelete{ModelName}` 删除验证方法
- 生成 `validateStatusUpdate{ModelName}` 状态更新验证方法
- 提供 TODO 注释指导具体验证逻辑实现
- 支持自定义业务规则扩展

## 4. Handler HTTP 处理层

### 4.1 Handler 结构生成

**功能描述**: 生成基于泛型的 HTTP 处理器

**实现效果**:
- 生成位置：`internal/admin/handler/{model_name}_handler.go`
- 继承 `BaseHandler[models.{ModelName}, types.{CreateRequest}, types.{UpdateRequest}, types.{ListRequest}, types.{Response}]`
- 自动生成处理器结构体和构造函数
- 配置处理器选项（软删除、状态管理、批量操作）
- 依赖注入对应的 Service

### 4.2 数据转换方法

**功能描述**: 实现请求响应数据转换

**实现效果**:
- `ConvertToResponse`: 实体转响应格式，处理状态文本和时间格式化
- `ConvertCreateRequest`: 创建请求转实体
- `ConvertUpdateRequest`: 更新请求转实体
- `ConvertListRequest`: 列表请求转分页查询，包含参数验证
- `ConvertListToResponse`: 列表转响应格式

### 4.3 辅助功能

**功能描述**: 提供状态和时间处理辅助方法

**实现效果**:
- `getStatusValue`: 处理状态值的空值情况
- `getStatusText`: 状态值转文本描述（启用/禁用/已删除）
- `formatTime`: 时间格式化（仅在有时间字段时生成）
- 自动检测字段类型生成对应的处理逻辑
- 支持状态字段的特殊处理

## 5. Routes 路由配置

### 5.1 路由文件生成

**功能描述**: 生成 RESTful 路由配置

**实现效果**:
- 生成位置：`internal/admin/routes/{model_name}_routes_gen.go`
- 基于模板生成标准的 CRUD 路由
- 自动配置路由分组和前缀
- 集成处理器方法绑定
- 支持中间件配置

## 6. Wire 依赖注入

### 6.1 Provider 生成

**功能描述**: 生成 Wire 依赖注入配置

**实现效果**:
- 生成位置：`internal/admin/wire/{model_name}_wire_gen.go`
- 自动生成 Provider 函数
- 处理 Repository、Service、Handler 的依赖关系
- 支持依赖注入链的自动配置
- 集成到主 Wire 配置中

## 7. Migration 数据库迁移

### 7.1 迁移文件生成

**功能描述**: 生成数据库迁移脚本

**实现效果**:
- 基于模型字段生成建表语句
- 处理字段类型和约束映射
- 支持索引和外键创建
- 生成迁移版本控制
- 提供回滚脚本

## 8. Permission 权限配置

### 8.1 权限定义生成

**功能描述**: 生成权限配置文件

**实现效果**:
- 生成标准的 CRUD 权限定义
- 配置资源和操作权限
- 集成到权限管理系统
- 支持角色权限绑定
- 提供权限检查中间件配置

## 9. 全量模块生成

### 9.1 All 组件生成

**功能描述**: 一次性生成完整的 CRUD 功能模块

**实现效果**:
- 按顺序生成：Model → Repository → Service → Handler → Routes → Wire → Migration → Permission
- 确保组件间依赖关系正确
- 统一的错误处理和成功响应
- 收集所有生成的文件路径
- 更新生成历史记录

### 9.2 生成流程控制

**功能描述**: 智能的生成流程管理

**实现效果**:
- 单个组件失败不影响其他组件生成
- 收集所有错误信息统一返回
- 支持部分成功的场景处理
- 提供详细的生成报告

## 10. 参数验证和边界处理

### 10.1 输入参数验证

**功能描述**: 严格的参数验证机制

**实现效果**:
- 组件类型枚举验证
- 模型名称格式验证（Go 标识符、非关键字、帕斯卡命名）
- 字段定义验证（名称、类型、重复检查）
- 表名格式验证（蛇形命名）
- 字段数量限制（最多 50 个）

### 10.2 文件冲突处理

**功能描述**: 智能的文件覆盖策略

**实现效果**:
- `overwrite_existing` 参数控制覆盖行为
- 文件存在时的冲突检测
- 目录自动创建
- 路径验证和规范化

## 11. 代码质量保证

### 11.1 代码格式化

**功能描述**: 自动代码格式化和优化

**实现效果**:
- `format_code` 参数控制是否格式化
- `optimize_imports` 参数控制导入优化
- 使用 Go 官方格式化工具
- 格式化失败不影响生成结果

### 11.2 模板系统

**功能描述**: 基于模板的代码生成

**实现效果**:
- 模板位置：`internal/devtools/generator/templates/`
- 支持条件渲染（如时间字段检测）
- 模板函数支持（类型转换、命名转换等）
- 模板数据预处理和验证

## 12. 历史记录管理

### 12.1 生成历史追踪

**功能描述**: 完整的代码生成历史管理

**实现效果**:
- 历史文件：`data/code-generate-history.json`
- 记录生成时间、组件类型、文件列表
- 支持按模块查询历史记录
- 提供历史记录清理功能

### 12.2 历史记录操作

**功能描述**: 历史记录的管理操作

**实现效果**:
- `GetHistory()`: 获取所有历史记录
- `GetHistoryByModule()`: 按模块获取历史
- `DeleteHistory()`: 删除指定模块历史
- `ClearHistory()`: 清空所有历史记录

## 总结

基于现有代码实现，后端代码生成器提供了完整的 CRUD 模块生成能力：

**核心特性**:
- 8 种组件类型支持（model、repository、service、handler、routes、wire、migration、permission）
- 基于泛型的类型安全实现
- 继承共享基础组件减少重复代码
- 智能的字段类型映射和处理
- 完整的参数验证和错误处理

**生成文件结构**:
- Models: `internal/shared/models/`
- Repository: `internal/admin/repository/`
- Service: `internal/admin/service/`
- Handler: `internal/admin/handler/`
- Routes: `internal/admin/routes/`
- Wire: `internal/admin/wire/`

**质量保证**:
- 自动代码格式化
- 导入优化
- 文件冲突处理
- 生成历史追踪