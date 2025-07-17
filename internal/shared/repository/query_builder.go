package repository

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"bico-admin/internal/shared/types"
)

// QueryBuilder 查询构建器
type QueryBuilder[T any] struct {
	db         *gorm.DB
	conditions []string
	args       []interface{}
	joins      []string
	preloads   []string
	groupBy    string
	having     string
	havingArgs []interface{}
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder[T any](db *gorm.DB) *QueryBuilder[T] {
	return &QueryBuilder[T]{
		db: db,
	}
}

// Where 添加 WHERE 条件
func (qb *QueryBuilder[T]) Where(condition string, args ...interface{}) *QueryBuilder[T] {
	qb.conditions = append(qb.conditions, condition)
	qb.args = append(qb.args, args...)
	return qb
}

// WhereIn 添加 IN 条件
func (qb *QueryBuilder[T]) WhereIn(field string, values []interface{}) *QueryBuilder[T] {
	if len(values) > 0 {
		placeholders := make([]string, len(values))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		condition := fmt.Sprintf("%s IN (%s)", field, strings.Join(placeholders, ","))
		qb.conditions = append(qb.conditions, condition)
		qb.args = append(qb.args, values...)
	}
	return qb
}

// WhereBetween 添加 BETWEEN 条件
func (qb *QueryBuilder[T]) WhereBetween(field string, start, end interface{}) *QueryBuilder[T] {
	condition := fmt.Sprintf("%s BETWEEN ? AND ?", field)
	qb.conditions = append(qb.conditions, condition)
	qb.args = append(qb.args, start, end)
	return qb
}

// WhereLike 添加 LIKE 条件
func (qb *QueryBuilder[T]) WhereLike(field string, value string) *QueryBuilder[T] {
	condition := fmt.Sprintf("%s LIKE ?", field)
	qb.conditions = append(qb.conditions, condition)
	qb.args = append(qb.args, "%"+value+"%")
	return qb
}

// WhereNotNull 添加 NOT NULL 条件
func (qb *QueryBuilder[T]) WhereNotNull(field string) *QueryBuilder[T] {
	condition := fmt.Sprintf("%s IS NOT NULL", field)
	qb.conditions = append(qb.conditions, condition)
	return qb
}

// WhereNull 添加 NULL 条件
func (qb *QueryBuilder[T]) WhereNull(field string) *QueryBuilder[T] {
	condition := fmt.Sprintf("%s IS NULL", field)
	qb.conditions = append(qb.conditions, condition)
	return qb
}

// Join 添加 JOIN
func (qb *QueryBuilder[T]) Join(joinClause string) *QueryBuilder[T] {
	qb.joins = append(qb.joins, joinClause)
	return qb
}

// LeftJoin 添加 LEFT JOIN
func (qb *QueryBuilder[T]) LeftJoin(table string, condition string) *QueryBuilder[T] {
	joinClause := fmt.Sprintf("LEFT JOIN %s ON %s", table, condition)
	qb.joins = append(qb.joins, joinClause)
	return qb
}

// InnerJoin 添加 INNER JOIN
func (qb *QueryBuilder[T]) InnerJoin(table string, condition string) *QueryBuilder[T] {
	joinClause := fmt.Sprintf("INNER JOIN %s ON %s", table, condition)
	qb.joins = append(qb.joins, joinClause)
	return qb
}

// Preload 添加预加载
func (qb *QueryBuilder[T]) Preload(association string) *QueryBuilder[T] {
	qb.preloads = append(qb.preloads, association)
	return qb
}

// GroupBy 添加 GROUP BY
func (qb *QueryBuilder[T]) GroupBy(fields string) *QueryBuilder[T] {
	qb.groupBy = fields
	return qb
}

// Having 添加 HAVING 条件
func (qb *QueryBuilder[T]) Having(condition string, args ...interface{}) *QueryBuilder[T] {
	qb.having = condition
	qb.havingArgs = args
	return qb
}

// buildQuery 构建查询
func (qb *QueryBuilder[T]) buildQuery(ctx context.Context) *gorm.DB {
	var entity T
	db := qb.db.WithContext(ctx).Model(&entity)

	// 添加 WHERE 条件
	if len(qb.conditions) > 0 {
		whereClause := strings.Join(qb.conditions, " AND ")
		db = db.Where(whereClause, qb.args...)
	}

	// 添加 JOIN
	for _, join := range qb.joins {
		db = db.Joins(join)
	}

	// 添加预加载
	for _, preload := range qb.preloads {
		db = db.Preload(preload)
	}

	// 添加 GROUP BY
	if qb.groupBy != "" {
		db = db.Group(qb.groupBy)
	}

	// 添加 HAVING
	if qb.having != "" {
		db = db.Having(qb.having, qb.havingArgs...)
	}

	return db
}

// Find 查询多条记录
func (qb *QueryBuilder[T]) Find(ctx context.Context) ([]*T, error) {
	var entities []*T
	db := qb.buildQuery(ctx)
	err := db.Find(&entities).Error
	return entities, err
}

// First 查询第一条记录
func (qb *QueryBuilder[T]) First(ctx context.Context) (*T, error) {
	var entity T
	db := qb.buildQuery(ctx)
	err := db.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Count 统计数量
func (qb *QueryBuilder[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	db := qb.buildQuery(ctx)
	err := db.Count(&count).Error
	return count, err
}

// Paginate 分页查询
func (qb *QueryBuilder[T]) Paginate(ctx context.Context, req *types.BasePageQuery) ([]*T, int64, error) {
	// 先统计总数
	total, err := qb.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 再查询数据
	var entities []*T
	db := qb.buildQuery(ctx)

	// 添加排序
	if req.SortBy != "" {
		sortField := toSnakeCase(req.SortBy)
		sortOrder := req.GetSortOrder()
		db = db.Order(fmt.Sprintf("%s %s", sortField, sortOrder))
	} else {
		db = db.Order("created_at DESC")
	}

	// 分页
	offset := req.GetOffset()
	pageSize := req.GetPageSize()
	err = db.Offset(offset).Limit(pageSize).Find(&entities).Error

	return entities, total, err
}

// Exists 检查是否存在
func (qb *QueryBuilder[T]) Exists(ctx context.Context) (bool, error) {
	count, err := qb.Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update 更新记录
func (qb *QueryBuilder[T]) Update(ctx context.Context, updates map[string]interface{}) error {
	db := qb.buildQuery(ctx)
	return db.Updates(updates).Error
}

// Delete 删除记录
func (qb *QueryBuilder[T]) Delete(ctx context.Context) error {
	var entity T
	db := qb.buildQuery(ctx)
	return db.Delete(&entity).Error
}

// Clone 克隆查询构建器
func (qb *QueryBuilder[T]) Clone() *QueryBuilder[T] {
	newQB := &QueryBuilder[T]{
		db:         qb.db,
		conditions: make([]string, len(qb.conditions)),
		args:       make([]interface{}, len(qb.args)),
		joins:      make([]string, len(qb.joins)),
		preloads:   make([]string, len(qb.preloads)),
		groupBy:    qb.groupBy,
		having:     qb.having,
		havingArgs: make([]interface{}, len(qb.havingArgs)),
	}

	copy(newQB.conditions, qb.conditions)
	copy(newQB.args, qb.args)
	copy(newQB.joins, qb.joins)
	copy(newQB.preloads, qb.preloads)
	copy(newQB.havingArgs, qb.havingArgs)

	return newQB
}

// Reset 重置查询构建器
func (qb *QueryBuilder[T]) Reset() *QueryBuilder[T] {
	qb.conditions = qb.conditions[:0]
	qb.args = qb.args[:0]
	qb.joins = qb.joins[:0]
	qb.preloads = qb.preloads[:0]
	qb.groupBy = ""
	qb.having = ""
	qb.havingArgs = qb.havingArgs[:0]
	return qb
}
