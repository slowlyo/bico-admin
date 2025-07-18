package handler

import (
	"context"

	"github.com/gin-gonic/gin"

	"bico-admin/internal/shared/types"
)

// CRUDHandler 定义标准的 CRUD 操作接口
type CRUDHandler interface {
	// 基础 CRUD 操作
	GetByID(c *gin.Context)   // GET /:id
	Create(c *gin.Context)    // POST /
	Update(c *gin.Context)    // PUT /:id
	Delete(c *gin.Context)    // DELETE /:id
	List(c *gin.Context)      // GET /
	UpdateStatus(c *gin.Context) // PUT /:id/status
}

// HandlerConverter 定义数据转换接口
type HandlerConverter[T any, CreateReq any, UpdateReq any, ListReq any, Response any] interface {
	// 数据转换方法
	ConvertToResponse(c *gin.Context, entity *T) Response
	ConvertCreateRequest(c *gin.Context, req *CreateReq) *T
	ConvertUpdateRequest(c *gin.Context, id uint, req *UpdateReq) *T
	ConvertListRequest(c *gin.Context, req *ListReq) *types.BasePageQuery
	ConvertListToResponse(c *gin.Context, list any) []Response
}

// ExtendedHandler 扩展的 handler 接口，包含更多业务方法
type ExtendedHandler interface {
	CRUDHandler
	
	// 批量操作
	BatchDelete(c *gin.Context)      // DELETE /batch
	BatchUpdateStatus(c *gin.Context) // PUT /batch/status
	
	// 导入导出
	Export(c *gin.Context)  // GET /export
	Import(c *gin.Context)  // POST /import
	
	// 其他常用操作
	GetOptions(c *gin.Context) // GET /options - 获取选项列表（用于下拉框等）
}

// ServiceInterface 定义服务层接口，与 BaseService 保持一致
type ServiceInterface[T any] interface {
	GetByID(ctx context.Context, id uint) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
	UpdateStatus(ctx context.Context, id uint, status int) error
	ListWithFilter(ctx context.Context, req *types.BasePageQuery) (*types.PageResult, error)
}

// ExtendedServiceInterface 扩展的服务接口
type ExtendedServiceInterface[T any] interface {
	ServiceInterface[T]
	
	// 批量操作
	BatchDelete(ctx context.Context, ids []uint) error
	BatchUpdateStatus(ctx context.Context, ids []uint, status int) error
	
	// 验证方法
	ExistsByField(ctx context.Context, field string, value any) (bool, error)
	
	// 统计方法
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status int) (int64, error)
}

// HandlerOptions 处理器选项配置
type HandlerOptions struct {
	// 是否启用软删除
	EnableSoftDelete bool
	
	// 是否启用状态管理
	EnableStatusManagement bool
	
	// 是否启用批量操作
	EnableBatchOperations bool
	
	// 是否启用导入导出
	EnableImportExport bool
	
	// 默认分页大小
	DefaultPageSize int
	
	// 最大分页大小
	MaxPageSize int
	
	// 是否启用缓存
	EnableCache bool
	
	// 缓存过期时间（秒）
	CacheExpiration int
}

// DefaultHandlerOptions 返回默认的处理器选项
func DefaultHandlerOptions() *HandlerOptions {
	return &HandlerOptions{
		EnableSoftDelete:       true,
		EnableStatusManagement: true,
		EnableBatchOperations:  true,
		EnableImportExport:     false,
		DefaultPageSize:        10,
		MaxPageSize:            100,
		EnableCache:            false,
		CacheExpiration:        300, // 5分钟
	}
}

// ValidationRule 验证规则定义
type ValidationRule struct {
	Field    string `json:"field"`    // 字段名
	Rule     string `json:"rule"`     // 验证规则
	Message  string `json:"message"`  // 错误消息
	Required bool   `json:"required"` // 是否必填
}

// FieldPermission 字段权限定义
type FieldPermission struct {
	Field      string   `json:"field"`       // 字段名
	ReadRoles  []string `json:"read_roles"`  // 可读角色
	WriteRoles []string `json:"write_roles"` // 可写角色
	HideRoles  []string `json:"hide_roles"`  // 隐藏角色
}

// HandlerMetadata 处理器元数据
type HandlerMetadata struct {
	// 基础信息
	Name        string `json:"name"`         // 处理器名称
	Description string `json:"description"`  // 描述
	Version     string `json:"version"`      // 版本
	
	// 路由信息
	BasePath    string            `json:"base_path"`    // 基础路径
	RoutePrefix string            `json:"route_prefix"` // 路由前缀
	Middlewares []string          `json:"middlewares"`  // 中间件列表
	
	// 权限信息
	RequiredPermissions []string          `json:"required_permissions"` // 必需权限
	FieldPermissions    []FieldPermission `json:"field_permissions"`    // 字段权限
	
	// 验证规则
	ValidationRules []ValidationRule `json:"validation_rules"` // 验证规则
	
	// 配置选项
	Options *HandlerOptions `json:"options"` // 处理器选项
}

// GetMetadata 获取处理器元数据（需要具体 handler 实现）
type MetadataProvider interface {
	GetMetadata() *HandlerMetadata
}

// HealthChecker 健康检查接口
type HealthChecker interface {
	HealthCheck(c *gin.Context)
}

// CacheManager 缓存管理接口
type CacheManager interface {
	Get(key string) (any, bool)
	Set(key string, value any, expiration int) error
	Delete(key string) error
	Clear() error
}

// EventHandler 事件处理接口
type EventHandler[T any] interface {
	// 生命周期事件
	BeforeCreate(ctx context.Context, entity *T) error
	AfterCreate(ctx context.Context, entity *T) error
	
	BeforeUpdate(ctx context.Context, entity *T) error
	AfterUpdate(ctx context.Context, entity *T) error
	
	BeforeDelete(ctx context.Context, id uint) error
	AfterDelete(ctx context.Context, id uint) error
	
	// 状态变更事件
	OnStatusChange(ctx context.Context, id uint, oldStatus, newStatus int) error
}
