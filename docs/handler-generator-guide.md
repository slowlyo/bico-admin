# Handler 生成器使用指南

## 概述

Handler 生成器是 bico-admin 项目中的代码生成工具，用于自动生成标准的 CRUD Handler 代码。它基于现有的 model 定义，自动生成包含完整 CRUD 操作的 Handler 和相关类型定义。

## 特性

### 🚀 核心功能
- **自动生成 Handler**: 基于 model 定义自动生成完整的 Handler 代码
- **类型安全**: 使用 Go 泛型确保类型安全
- **标准化 CRUD**: 提供标准的增删改查操作
- **响应格式统一**: 使用项目统一的响应格式
- **状态管理**: 智能处理状态字段（如果存在）
- **时间格式化**: 自动处理时间字段的格式化

### 🛠️ 高级功能
- **批量操作**: 支持批量删除和状态更新
- **缓存支持**: 可选的缓存功能
- **事件处理**: 支持生命周期事件钩子
- **权限控制**: 集成权限验证
- **健康检查**: 内置健康检查端点

## 架构设计

### BaseHandler 基类

BaseHandler 是所有生成的 Handler 的基类，提供以下功能：

```go
type BaseHandler[T any, CreateReq any, UpdateReq any, ListReq any, Response any] struct {
    service  ServiceInterface[T]
    options  *HandlerOptions
    metadata *HandlerMetadata
    cache    CacheManager
}
```

### 核心接口

#### CRUDHandler 接口
```go
type CRUDHandler interface {
    GetByID(c *gin.Context)      // GET /:id
    Create(c *gin.Context)       // POST /
    Update(c *gin.Context)       // PUT /:id
    Delete(c *gin.Context)       // DELETE /:id
    List(c *gin.Context)         // GET /
    UpdateStatus(c *gin.Context) // PUT /:id/status
}
```

#### HandlerConverter 接口
```go
type HandlerConverter[T any, CreateReq any, UpdateReq any, ListReq any, Response any] interface {
    ConvertToResponse(c *gin.Context, entity *T) Response
    ConvertCreateRequest(c *gin.Context, req *CreateReq) *T
    ConvertUpdateRequest(c *gin.Context, id uint, req *UpdateReq) *T
    ConvertListRequest(c *gin.Context, req *ListReq) *types.BasePageQuery
    ConvertListToResponse(c *gin.Context, list any) []Response
}
```

## 使用方法

### 1. 基本用法

使用代码生成器生成 Handler：

```go
package main

import (
    "bico-admin/internal/devtools/generator"
)

func main() {
    // 创建代码生成器
    codeGen := generator.NewCodeGenerator()

    // 定义字段
    fields := []generator.FieldDefinition{
        {
            Name:     "Name",
            Type:     "string",
            JsonTag:  "name",
            Validate: "required,min=1,max=100",
            Comment:  "产品名称",
        },
        {
            Name:     "Price",
            Type:     "float64",
            JsonTag:  "price",
            Validate: "required,min=0",
            Comment:  "价格",
        },
        {
            Name:     "Status",
            Type:     "*int",
            JsonTag:  "status",
            Validate: "oneof=0 1",
            Comment:  "状态",
        },
    }

    // 生成单个 Handler
    req := &generator.GenerateRequest{
        ModelName:     "Product",
        ComponentType: generator.ComponentHandler,
        Fields:        fields,
        Options: generator.GenerateOptions{
            FormatCode:        true,
            OverwriteExisting: true,
        },
    }

    response, err := codeGen.Generate(req)
    if err != nil {
        panic(err)
    }

    // 生成完整的 CRUD 模块
    req.ComponentType = generator.ComponentAll
    response, err = codeGen.Generate(req)
}
```

### 2. 字段定义

字段定义支持以下属性：

```json
{
  "name": "字段名（Go 标识符）",
  "type": "字段类型（Go 类型）",
  "json_tag": "JSON 标签（可选，默认为蛇形命名）",
  "validate": "验证规则（gin 验证标签）",
  "gorm_tag": "GORM 标签（数据库相关）",
  "comment": "字段注释"
}
```

#### 支持的字段类型

| 类型 | Go 类型 | 说明 |
|------|---------|------|
| string | string | 字符串 |
| int | int | 整数 |
| int64 | int64 | 64位整数 |
| uint | uint | 无符号整数 |
| float64 | float64 | 浮点数 |
| bool | bool | 布尔值 |
| time | *time.Time | 时间（指针类型） |
| datetime | *time.Time | 日期时间 |

### 3. 生成的文件结构

生成器会创建以下文件：

```
internal/admin/
├── handler/
│   └── product.go              # Handler 实现
└── types/
    └── product_types.go        # 请求/响应类型定义
```

### 4. 生成的代码示例

#### Handler 文件 (product.go)
```go
package handler

import (
    "github.com/gin-gonic/gin"
    "bico-admin/internal/admin/models"
    "bico-admin/internal/admin/service"
    "bico-admin/internal/admin/types"
    // ...
)

type ProductHandler struct {
    *BaseHandler[models.Product, types.ProductCreateRequest, types.ProductUpdateRequest, types.ProductListRequest, types.ProductResponse]
    productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
    options := DefaultHandlerOptions()
    options.EnableStatusManagement = true
    
    baseHandler := NewBaseHandler[models.Product, types.ProductCreateRequest, types.ProductUpdateRequest, types.ProductListRequest, types.ProductResponse](productService, options)
    
    return &ProductHandler{
        BaseHandler: baseHandler,
        productService: productService,
    }
}

// 实现数据转换方法
func (h *ProductHandler) ConvertToResponse(c *gin.Context, entity *models.Product) types.ProductResponse {
    return types.ProductResponse{
        ID:         entity.ID,
        Name:       entity.Name,
        Price:      entity.Price,
        Status:     h.getStatusValue(entity.Status),
        StatusText: h.getStatusText(h.getStatusValue(entity.Status)),
        CreatedAt:  utils.NewFormattedTime(entity.CreatedAt),
        UpdatedAt:  utils.NewFormattedTime(entity.UpdatedAt),
    }
}

// ... 其他转换方法
```

#### 类型定义文件 (product_types.go)
```go
package types

import (
    "bico-admin/internal/shared/types"
    "bico-admin/pkg/utils"
)

// ProductCreateRequest 创建产品请求
type ProductCreateRequest struct {
    Name   string  `json:"name" binding:"required,min=1,max=100"`
    Price  float64 `json:"price" binding:"required,min=0"`
    Status int     `json:"status" binding:"oneof=0 1"`
}

// ProductUpdateRequest 更新产品请求
type ProductUpdateRequest struct {
    Name   string  `json:"name" binding:"required,min=1,max=100"`
    Price  float64 `json:"price" binding:"required,min=0"`
    Status int     `json:"status" binding:"oneof=0 1"`
}

// ProductListRequest 产品列表请求
type ProductListRequest struct {
    types.BasePageQuery
    Name   string `form:"name" json:"name"`
    Status *int   `form:"status" json:"status"`
}

// ProductResponse 产品响应
type ProductResponse struct {
    ID         uint                `json:"id"`
    Name       string              `json:"name"`
    Price      float64             `json:"price"`
    Status     int                 `json:"status"`
    StatusText string              `json:"status_text"`
    CreatedAt  utils.FormattedTime `json:"created_at"`
    UpdatedAt  utils.FormattedTime `json:"updated_at"`
}

// ... 其他类型定义
```

## 配置选项

### HandlerOptions

```go
type HandlerOptions struct {
    EnableSoftDelete       bool // 是否启用软删除
    EnableStatusManagement bool // 是否启用状态管理
    EnableBatchOperations  bool // 是否启用批量操作
    EnableImportExport     bool // 是否启用导入导出
    DefaultPageSize        int  // 默认分页大小
    MaxPageSize           int  // 最大分页大小
    EnableCache           bool // 是否启用缓存
    CacheExpiration       int  // 缓存过期时间（秒）
}
```

### 默认配置

```go
func DefaultHandlerOptions() *HandlerOptions {
    return &HandlerOptions{
        EnableSoftDelete:       true,
        EnableStatusManagement: true,
        EnableBatchOperations:  true,
        EnableImportExport:     false,
        DefaultPageSize:        10,
        MaxPageSize:            100,
        EnableCache:            false,
        CacheExpiration:        300,
    }
}
```

## 最佳实践

### 1. 字段命名规范
- 使用 PascalCase 命名字段
- 状态字段统一命名为 `Status`
- 时间字段使用 `*time.Time` 类型

### 2. 验证规则
- 必填字段使用 `required` 标签
- 字符串长度使用 `min` 和 `max` 标签
- 枚举值使用 `oneof` 标签

### 3. 状态管理
- 使用 `shared/types` 中定义的状态常量
- 状态字段类型为 `*int`（指针类型）
- 提供状态文本转换

### 4. 扩展自定义方法
生成的 Handler 可以添加自定义业务方法：

```go
// 自定义业务方法
func (h *ProductHandler) GetByCategory(c *gin.Context) {
    category := c.Param("category")
    // 实现业务逻辑
}
```

## 注意事项

1. **依赖关系**: Handler 依赖对应的 Service 接口，确保先生成 Service
2. **类型安全**: 使用泛型确保编译时类型检查
3. **接口实现**: 生成的 Handler 必须实现所有转换方法
4. **状态字段**: 如果模型没有状态字段，会提供默认实现
5. **缓存管理**: 缓存功能需要额外配置缓存管理器

## 故障排除

### 常见问题

1. **编译错误**: 检查字段类型定义是否正确
2. **验证失败**: 检查验证标签语法
3. **接口不匹配**: 确保 Service 接口已正确实现
4. **导入错误**: 检查包路径是否正确

### 调试技巧

1. 使用 `--format_code=false` 查看原始生成代码
2. 检查生成历史记录
3. 逐步生成各个组件进行调试
