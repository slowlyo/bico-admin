package repository

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"

	"bico-admin/internal/shared/types"
)

// BaseRepositoryInterface 基础仓储接口
type BaseRepositoryInterface[T any] interface {
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
	ListWithCondition(ctx context.Context, condition string, args []interface{}, req *types.BasePageQuery) ([]*T, int64, error)

	// 统计操作
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status int) (int64, error)
	CountWithCondition(ctx context.Context, condition string, args []interface{}) (int64, error)

	// 状态管理
	UpdateStatus(ctx context.Context, id uint, status int) error
	BatchUpdateStatus(ctx context.Context, ids []uint, status int) error

	// 字段检查
	ExistsByField(ctx context.Context, field string, value interface{}) (bool, error)
	ExistsByFieldExcludeID(ctx context.Context, field string, value interface{}, excludeID uint) (bool, error)

	// 事务支持
	WithTx(tx *gorm.DB) BaseRepositoryInterface[T]

	// 数据库实例
	DB() *gorm.DB
}

// BaseRepository 基础仓储实现
type BaseRepository[T any] struct {
	db    *gorm.DB
	model T
}

// NewBaseRepository 创建基础仓储
func NewBaseRepository[T any](db *gorm.DB) BaseRepositoryInterface[T] {
	var model T
	return &BaseRepository[T]{
		db:    db,
		model: model,
	}
}

// WithTx 使用事务
func (r *BaseRepository[T]) WithTx(tx *gorm.DB) BaseRepositoryInterface[T] {
	return &BaseRepository[T]{
		db:    tx,
		model: r.model,
	}
}

// DB 获取数据库实例
func (r *BaseRepository[T]) DB() *gorm.DB {
	return r.db
}

// Create 创建记录
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// BatchCreate 批量创建记录
func (r *BaseRepository[T]) BatchCreate(ctx context.Context, entities []*T) error {
	if len(entities) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(entities).Error
}

// GetByID 根据ID获取记录
func (r *BaseRepository[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update 更新记录
func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete 删除记录
func (r *BaseRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

// BatchDelete 批量删除记录
func (r *BaseRepository[T]) BatchDelete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, ids).Error
}

// Count 统计总数
func (r *BaseRepository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	var entity T
	err := r.db.WithContext(ctx).Model(&entity).Count(&count).Error
	return count, err
}

// CountByStatus 根据状态统计数量
func (r *BaseRepository[T]) CountByStatus(ctx context.Context, status int) (int64, error) {
	var count int64
	var entity T
	err := r.db.WithContext(ctx).Model(&entity).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

// CountWithCondition 根据条件统计数量
func (r *BaseRepository[T]) CountWithCondition(ctx context.Context, condition string, args []interface{}) (int64, error) {
	var count int64
	var entity T
	db := r.db.WithContext(ctx).Model(&entity)
	if condition != "" {
		db = db.Where(condition, args...)
	}
	err := db.Count(&count).Error
	return count, err
}

// UpdateStatus 更新状态
func (r *BaseRepository[T]) UpdateStatus(ctx context.Context, id uint, status int) error {
	var entity T
	return r.db.WithContext(ctx).Model(&entity).
		Where("id = ?", id).
		Update("status", status).Error
}

// BatchUpdateStatus 批量更新状态
func (r *BaseRepository[T]) BatchUpdateStatus(ctx context.Context, ids []uint, status int) error {
	if len(ids) == 0 {
		return nil
	}
	var entity T
	return r.db.WithContext(ctx).Model(&entity).
		Where("id IN ?", ids).
		Update("status", status).Error
}

// ExistsByField 检查字段值是否存在
func (r *BaseRepository[T]) ExistsByField(ctx context.Context, field string, value interface{}) (bool, error) {
	var count int64
	var entity T
	err := r.db.WithContext(ctx).Model(&entity).
		Where(field+" = ?", value).
		Count(&count).Error
	return count > 0, err
}

// ExistsByFieldExcludeID 检查字段值是否存在（排除指定ID）
func (r *BaseRepository[T]) ExistsByFieldExcludeID(ctx context.Context, field string, value interface{}, excludeID uint) (bool, error) {
	var count int64
	var entity T
	err := r.db.WithContext(ctx).Model(&entity).
		Where(field+" = ? AND id != ?", value, excludeID).
		Count(&count).Error
	return count > 0, err
}

// ListByStatus 根据状态分页查询
func (r *BaseRepository[T]) ListByStatus(ctx context.Context, status int, req *types.BasePageQuery) ([]*T, int64, error) {
	return r.ListWithCondition(ctx, "status = ?", []interface{}{status}, req)
}

// List 分页查询所有记录
func (r *BaseRepository[T]) List(ctx context.Context, req *types.BasePageQuery) ([]*T, int64, error) {
	return r.ListWithCondition(ctx, "", nil, req)
}

// ListWithCondition 根据条件分页查询
func (r *BaseRepository[T]) ListWithCondition(ctx context.Context, condition string, args []interface{}, req *types.BasePageQuery) ([]*T, int64, error) {
	var entities []*T
	var total int64
	var entity T

	db := r.db.WithContext(ctx).Model(&entity)

	// 添加条件
	if condition != "" {
		db = db.Where(condition, args...)
	}

	// 添加关键词搜索（如果支持）
	if req.Keyword != "" {
		db = r.addKeywordSearch(db, req.Keyword)
	}

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 添加排序
	db = r.addSorting(db, req)

	// 分页查询
	offset := req.GetOffset()
	pageSize := req.GetPageSize()
	err := db.Offset(offset).Limit(pageSize).Find(&entities).Error

	return entities, total, err
}

// addKeywordSearch 添加关键词搜索
func (r *BaseRepository[T]) addKeywordSearch(db *gorm.DB, keyword string) *gorm.DB {
	// 获取模型的反射类型
	modelType := reflect.TypeOf(r.model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	var searchFields []string
	
	// 遍历字段，找到可搜索的字符串字段
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		
		// 跳过嵌入字段和非导出字段
		if field.Anonymous || !field.IsExported() {
			continue
		}
		
		// 只搜索字符串类型字段
		if field.Type.Kind() == reflect.String {
			// 获取数据库字段名
			gormTag := field.Tag.Get("gorm")
			var dbFieldName string
			
			if strings.Contains(gormTag, "column:") {
				// 从 gorm 标签中提取列名
				parts := strings.Split(gormTag, ";")
				for _, part := range parts {
					if strings.HasPrefix(part, "column:") {
						dbFieldName = strings.TrimPrefix(part, "column:")
						break
					}
				}
			} else {
				// 使用字段名的蛇形命名
				dbFieldName = toSnakeCase(field.Name)
			}
			
			searchFields = append(searchFields, dbFieldName)
		}
	}

	// 构建搜索条件
	if len(searchFields) > 0 {
		var conditions []string
		var args []interface{}
		
		for _, field := range searchFields {
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", field))
			args = append(args, "%"+keyword+"%")
		}
		
		if len(conditions) > 0 {
			db = db.Where(strings.Join(conditions, " OR "), args...)
		}
	}

	return db
}

// addSorting 添加排序
func (r *BaseRepository[T]) addSorting(db *gorm.DB, req *types.BasePageQuery) *gorm.DB {
	if req.SortBy != "" {
		sortField := toSnakeCase(req.SortBy)
		sortOrder := req.GetSortOrder()
		db = db.Order(fmt.Sprintf("%s %s", sortField, sortOrder))
	} else {
		// 默认按创建时间降序
		db = db.Order("created_at DESC")
	}
	return db
}

// toSnakeCase 将驼峰命名转换为蛇形命名
func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// GetModelType 获取模型类型（用于调试）
func (r *BaseRepository[T]) GetModelType() reflect.Type {
	return reflect.TypeOf(r.model)
}
