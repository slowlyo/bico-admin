package service

import (
	"context"
	"fmt"

	"bico-admin/internal/shared/types"
)

// ServiceError 服务层错误
type ServiceError struct {
	Code    string
	Message string
	Err     error
}

func (e *ServiceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *ServiceError) Unwrap() error {
	return e.Err
}

// 常见错误代码
const (
	ErrCodeValidation   = "VALIDATION_ERROR"
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeDuplicate    = "DUPLICATE"
	ErrCodePermission   = "PERMISSION_DENIED"
	ErrCodeBusiness     = "BUSINESS_ERROR"
	ErrCodeInternal     = "INTERNAL_ERROR"
)

// NewValidationError 创建验证错误
func NewValidationError(message string, err error) *ServiceError {
	return &ServiceError{
		Code:    ErrCodeValidation,
		Message: message,
		Err:     err,
	}
}

// NewNotFoundError 创建未找到错误
func NewNotFoundError(message string) *ServiceError {
	return &ServiceError{
		Code:    ErrCodeNotFound,
		Message: message,
	}
}

// NewDuplicateError 创建重复错误
func NewDuplicateError(message string) *ServiceError {
	return &ServiceError{
		Code:    ErrCodeDuplicate,
		Message: message,
	}
}

// NewPermissionError 创建权限错误
func NewPermissionError(message string) *ServiceError {
	return &ServiceError{
		Code:    ErrCodePermission,
		Message: message,
	}
}

// NewBusinessError 创建业务错误
func NewBusinessError(message string) *ServiceError {
	return &ServiceError{
		Code:    ErrCodeBusiness,
		Message: message,
	}
}

// NewInternalError 创建内部错误
func NewInternalError(message string, err error) *ServiceError {
	return &ServiceError{
		Code:    ErrCodeInternal,
		Message: message,
		Err:     err,
	}
}

// ValidatePageQuery 验证分页查询参数
func ValidatePageQuery(req *types.BasePageQuery) error {
	if req == nil {
		return NewValidationError("分页参数不能为空", nil)
	}

	if req.Page < 1 {
		return NewValidationError("页码必须大于0", nil)
	}

	if req.PageSize < 1 || req.PageSize > 100 {
		return NewValidationError("每页数量必须在1-100之间", nil)
	}

	return nil
}

// ValidateID 验证ID参数
func ValidateID(id uint, fieldName string) error {
	if id == 0 {
		return NewValidationError(fmt.Sprintf("%s不能为空", fieldName), nil)
	}
	return nil
}

// ValidateIDs 验证ID列表参数
func ValidateIDs(ids []uint, fieldName string) error {
	if len(ids) == 0 {
		return NewValidationError(fmt.Sprintf("%s列表不能为空", fieldName), nil)
	}

	for i, id := range ids {
		if id == 0 {
			return NewValidationError(fmt.Sprintf("%s列表中第%d个ID不能为空", fieldName, i+1), nil)
		}
	}

	return nil
}

// ValidateStatus 验证状态值
func ValidateStatus(status int) error {
	if status != types.StatusActive && status != types.StatusInactive && status != types.StatusDeleted {
		return NewValidationError("状态值无效", nil)
	}
	return nil
}

// BuildPageResult 构建分页结果
func BuildPageResult[T any](list []*T, total int64, page, pageSize int) *types.PageResult {
	return types.NewPageResult(list, total, page, pageSize)
}

// TransactionWrapper 事务包装器
type TransactionWrapper interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// ExecuteInTransaction 在事务中执行操作
func ExecuteInTransaction(ctx context.Context, wrapper TransactionWrapper, fn func(ctx context.Context) error) error {
	if wrapper == nil {
		// 如果没有事务包装器，直接执行
		return fn(ctx)
	}

	return wrapper.WithTransaction(ctx, fn)
}

// ConvertToMap 将结构体转换为map（用于更新操作）
func ConvertToMap(v interface{}) map[string]interface{} {
	// 这里可以使用反射或者其他方式将结构体转换为map
	// 为了简化，这里返回空map，具体实现可以根据需要添加
	return make(map[string]interface{})
}

// FilterNilPointers 过滤空指针
func FilterNilPointers[T any](items []*T) []*T {
	var result []*T
	for _, item := range items {
		if item != nil {
			result = append(result, item)
		}
	}
	return result
}

// ContainsID 检查ID列表是否包含指定ID
func ContainsID(ids []uint, targetID uint) bool {
	for _, id := range ids {
		if id == targetID {
			return true
		}
	}
	return false
}

// RemoveID 从ID列表中移除指定ID
func RemoveID(ids []uint, targetID uint) []uint {
	var result []uint
	for _, id := range ids {
		if id != targetID {
			result = append(result, id)
		}
	}
	return result
}

// UniqueIDs 去重ID列表
func UniqueIDs(ids []uint) []uint {
	seen := make(map[uint]bool)
	var result []uint
	
	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}
	
	return result
}

// ChunkIDs 将ID列表分块处理
func ChunkIDs(ids []uint, chunkSize int) [][]uint {
	if chunkSize <= 0 {
		chunkSize = 100 // 默认块大小
	}

	var chunks [][]uint
	for i := 0; i < len(ids); i += chunkSize {
		end := i + chunkSize
		if end > len(ids) {
			end = len(ids)
		}
		chunks = append(chunks, ids[i:end])
	}

	return chunks
}
