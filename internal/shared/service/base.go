package service

import (
	"context"
	"fmt"

	"bico-admin/internal/shared/repository"
	"bico-admin/internal/shared/types"
)

// BaseServiceInterface 基础服务接口
type BaseServiceInterface[T any, R repository.BaseRepositoryInterface[T]] interface {
	// 基础 CRUD 操作
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id uint) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error

	// 批量操作
	BatchCreate(ctx context.Context, entities []*T) error
	BatchDelete(ctx context.Context, ids []uint) error

	// 查询操作
	List(ctx context.Context, req *types.BasePageQuery) ([]*T, int64, error)
	ListByStatus(ctx context.Context, status int, req *types.BasePageQuery) ([]*T, int64, error)

	// 统计操作
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status int) (int64, error)

	// 状态管理
	UpdateStatus(ctx context.Context, id uint, status int) error
	BatchUpdateStatus(ctx context.Context, ids []uint, status int) error

	// 验证操作
	ExistsByID(ctx context.Context, id uint) (bool, error)

	// 获取 Repository 实例
	Repository() R
}

// BaseService 基础服务实现
type BaseService[T any, R repository.BaseRepositoryInterface[T]] struct {
	repo R
}

// NewBaseService 创建基础服务
func NewBaseService[T any, R repository.BaseRepositoryInterface[T]](repo R) BaseServiceInterface[T, R] {
	return &BaseService[T, R]{
		repo: repo,
	}
}

// Repository 获取 Repository 实例
func (s *BaseService[T, R]) Repository() R {
	return s.repo
}

// Create 创建实体
func (s *BaseService[T, R]) Create(ctx context.Context, entity *T) error {
	// 可以在这里添加通用的业务验证逻辑
	if err := s.validateEntity(ctx, entity); err != nil {
		return err
	}

	return s.repo.Create(ctx, entity)
}

// BatchCreate 批量创建实体
func (s *BaseService[T, R]) BatchCreate(ctx context.Context, entities []*T) error {
	if len(entities) == 0 {
		return nil
	}

	// 验证所有实体
	for _, entity := range entities {
		if err := s.validateEntity(ctx, entity); err != nil {
			return fmt.Errorf("验证实体失败: %w", err)
		}
	}

	return s.repo.BatchCreate(ctx, entities)
}

// GetByID 根据ID获取实体
func (s *BaseService[T, R]) GetByID(ctx context.Context, id uint) (*T, error) {
	if id == 0 {
		return nil, fmt.Errorf("ID不能为空")
	}

	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取实体失败: %w", err)
	}

	return entity, nil
}

// Update 更新实体
func (s *BaseService[T, R]) Update(ctx context.Context, entity *T) error {
	// 可以在这里添加通用的业务验证逻辑
	if err := s.validateEntity(ctx, entity); err != nil {
		return err
	}

	return s.repo.Update(ctx, entity)
}

// Delete 删除实体
func (s *BaseService[T, R]) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return fmt.Errorf("ID不能为空")
	}

	// 检查实体是否存在
	exists, err := s.ExistsByID(ctx, id)
	if err != nil {
		return fmt.Errorf("检查实体存在性失败: %w", err)
	}
	if !exists {
		return fmt.Errorf("实体不存在")
	}

	// 可以在这里添加删除前的业务验证
	if err := s.validateDelete(ctx, id); err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

// BatchDelete 批量删除实体
func (s *BaseService[T, R]) BatchDelete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	// 验证所有ID
	for _, id := range ids {
		if id == 0 {
			return fmt.Errorf("ID不能为空")
		}
		
		// 可以在这里添加删除前的业务验证
		if err := s.validateDelete(ctx, id); err != nil {
			return fmt.Errorf("验证删除条件失败 (ID: %d): %w", id, err)
		}
	}

	return s.repo.BatchDelete(ctx, ids)
}

// List 分页查询实体列表
func (s *BaseService[T, R]) List(ctx context.Context, req *types.BasePageQuery) ([]*T, int64, error) {
	if req == nil {
		req = &types.BasePageQuery{Page: 1, PageSize: 10}
	}

	return s.repo.List(ctx, req)
}

// ListByStatus 根据状态分页查询实体列表
func (s *BaseService[T, R]) ListByStatus(ctx context.Context, status int, req *types.BasePageQuery) ([]*T, int64, error) {
	if req == nil {
		req = &types.BasePageQuery{Page: 1, PageSize: 10}
	}

	return s.repo.ListByStatus(ctx, status, req)
}

// Count 统计实体总数
func (s *BaseService[T, R]) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// CountByStatus 根据状态统计实体数量
func (s *BaseService[T, R]) CountByStatus(ctx context.Context, status int) (int64, error) {
	return s.repo.CountByStatus(ctx, status)
}

// UpdateStatus 更新实体状态
func (s *BaseService[T, R]) UpdateStatus(ctx context.Context, id uint, status int) error {
	if id == 0 {
		return fmt.Errorf("ID不能为空")
	}

	// 检查实体是否存在
	exists, err := s.ExistsByID(ctx, id)
	if err != nil {
		return fmt.Errorf("检查实体存在性失败: %w", err)
	}
	if !exists {
		return fmt.Errorf("实体不存在")
	}

	// 可以在这里添加状态更新的业务验证
	if err := s.validateStatusUpdate(ctx, id, status); err != nil {
		return err
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

// BatchUpdateStatus 批量更新实体状态
func (s *BaseService[T, R]) BatchUpdateStatus(ctx context.Context, ids []uint, status int) error {
	if len(ids) == 0 {
		return nil
	}

	// 验证所有ID和状态更新条件
	for _, id := range ids {
		if id == 0 {
			return fmt.Errorf("ID不能为空")
		}
		
		if err := s.validateStatusUpdate(ctx, id, status); err != nil {
			return fmt.Errorf("验证状态更新条件失败 (ID: %d): %w", id, err)
		}
	}

	return s.repo.BatchUpdateStatus(ctx, ids, status)
}

// ExistsByID 检查实体是否存在
func (s *BaseService[T, R]) ExistsByID(ctx context.Context, id uint) (bool, error) {
	if id == 0 {
		return false, nil
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// 如果是记录不存在的错误，返回 false
		// 这里可以根据具体的 ORM 错误类型进行判断
		return false, nil
	}

	return true, nil
}

// validateEntity 验证实体（子类可以重写）
func (s *BaseService[T, R]) validateEntity(ctx context.Context, entity *T) error {
	// 基础验证逻辑，子类可以重写
	if entity == nil {
		return fmt.Errorf("实体不能为空")
	}
	return nil
}

// validateDelete 验证删除条件（子类可以重写）
func (s *BaseService[T, R]) validateDelete(ctx context.Context, id uint) error {
	// 基础删除验证逻辑，子类可以重写
	return nil
}

// validateStatusUpdate 验证状态更新条件（子类可以重写）
func (s *BaseService[T, R]) validateStatusUpdate(ctx context.Context, id uint, status int) error {
	// 基础状态更新验证逻辑，子类可以重写
	return nil
}
