package crud

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Unique 声明一个需要全表唯一的列。
//
// 说明：框架在创建和更新时于事务内自动查重（更新时排除自身），业务无需再手写 Exists 校验。
// 数据库仍应为该列建唯一索引，用来兜底并发写入。
type Unique struct {
	// Column 数据库列名
	Column string
	// Err 冲突时返回的业务错误（可选，为空时返回 "<Column> 已存在"）
	Err error
}

// CRUDHandler 通用 CRUD Handler，通过配置 hook 来适配不同业务模型。
// 目标：让单个 handler 文件只保留少量业务差异代码。
// T: 模型类型；L: List 查询参数类型；C: Create 请求体类型；U: Update 请求体类型。
//
// 说明：
// - List/Get/Create/Update/Delete 走统一的参数绑定、错误处理与响应
// - 差异点通过 BuildQuery/Before/After 等函数注入
// - 默认使用 gorm 的 First/Delete/Updates
// - Before*/After* 都运行在同一个事务内，*Commit 运行在事务提交之后
// - 业务校验失败请返回 BizError，其余错误只记日志，不会原样暴露给客户端
//
// 注意：为了保持 KISS，这里只覆盖标准 CRUD；复杂额外接口仍放在具体 handler 内。
type CRUDHandler[T any, L any, C any, U any] struct {
	BaseHandler

	DB          *gorm.DB
	NotFoundMsg string

	EnabledField      string
	EnabledSuccessMsg string

	CreateSuccessMsg string
	UpdateSuccessMsg string
	DeleteSuccessMsg string

	// SortableFields 允许客户端通过 sortField 排序的列名白名单（可选）。
	// 为空时忽略客户端排序，避免按密码哈希等敏感列排序泄露信息。
	SortableFields []string
	// DefaultOrder 默认排序子句（可选，为空时使用 created_at DESC）
	DefaultOrder string
	// UniqueFields 需要全表唯一的列（可选），创建与更新时自动在事务内查重
	UniqueFields []Unique

	// BuildListQuery 构建列表查询（可选，为 nil 则查询整张表）
	BuildListQuery func(db *gorm.DB, req *L) *gorm.DB
	// AfterList 列表查询后对结果二次处理（可选）
	AfterList func(items []T) error

	// BuildGetQuery 构建详情查询（可选，为 nil 则使用 DB.Model(&T{})）。
	// 配置后，创建和更新成功会按同一口径重新加载，返回值与详情接口保持一致。
	BuildGetQuery func(db *gorm.DB) *gorm.DB
	// AfterGet 查询详情后对结果二次处理（可选，创建和更新后同样会执行）
	AfterGet func(item *T) error

	// NewModelFromCreate Create 请求转模型（必须提供）。
	// 在事务之外执行，只适合做请求到模型的转换；依赖数据库的校验请放 BeforeCreate 或 UniqueFields。
	NewModelFromCreate func(req *C) (*T, error)
	// BeforeCreate 事务内、插入前的校验/补字段（可选）
	BeforeCreate func(tx *gorm.DB, item *T, req *C) error
	// AfterCreate 事务内、插入后的逻辑（可选，例如写关联表）
	AfterCreate func(tx *gorm.DB, item *T, req *C) error
	// AfterCreateCommit 创建事务提交后的逻辑（可选，不得返回业务错误）
	AfterCreateCommit func(item *T, req *C)

	// BuildUpdateQuery 构建更新时查询（可选，为 nil 则 tx.First）
	BuildUpdateQuery func(tx *gorm.DB) *gorm.DB
	// BuildUpdates Update 请求转 updates map（必须提供）。运行在事务内，查库请使用传入的 tx。
	BuildUpdates func(tx *gorm.DB, req *U, existing *T) (map[string]interface{}, error)
	// BeforeUpdate 事务内、更新前的校验/补字段（可选，可改写 updates）。
	// 保护性校验（如内置记录不可修改）请放这里：UpdateEnabled 不会调用 BuildUpdates。
	BeforeUpdate func(tx *gorm.DB, id uint, existing *T, req *U, updates map[string]interface{}) error
	// AfterUpdate 事务内、更新后的逻辑（可选，例如同步关联表）
	AfterUpdate func(tx *gorm.DB, id uint, existing *T, req *U) error
	// AfterUpdateCommit 更新事务提交后的逻辑（可选，不得返回业务错误）。
	// updates 为实际写入的字段，UpdateEnabled 路径下 req 是零值，应依据 updates 判断变更。
	AfterUpdateCommit func(id uint, existing *T, req *U, updates map[string]interface{})

	// BeforeDelete 事务内、删除前的校验与关联清理（可选）。
	// 单个删除与批量删除共用，单个删除时 ids 只有一个元素。
	BeforeDelete func(tx *gorm.DB, ids []uint) error
	// AfterDelete 事务内、删除后的逻辑（可选）
	AfterDelete func(tx *gorm.DB, ids []uint) error
	// AfterDeleteCommit 删除事务提交后的逻辑（可选，不得返回业务错误）
	AfterDeleteCommit func(ids []uint)
}

// List 获取列表
func (h *CRUDHandler[T, L, C, U]) List(c *gin.Context) {
	var req L
	if err := h.BindQuery(c, &req); err != nil {
		return
	}

	// 未提供列表查询构造时直接查询整张表，简单表无需任何配置。
	query := h.DB.Model(new(T))
	if h.BuildListQuery != nil {
		query = h.BuildListQuery(h.DB, &req)
	}

	var items []T
	order := ListOrder{Sortable: h.SortableFields, Default: h.DefaultOrder}
	QueryListWithHook(&h.BaseHandler, c, query, &items, order, func() error {
		if h.AfterList == nil {
			return nil
		}
		return h.AfterList(items)
	})
}

// Get 获取详情
func (h *CRUDHandler[T, L, C, U]) Get(c *gin.Context) {
	id, err := h.ParseID(c)
	// ID 解析失败直接返回，由 BaseHandler 已输出错误响应
	if err != nil {
		return
	}

	item, err := h.loadDetail(id)
	if err != nil {
		h.handleRecordError(c, err)
		return
	}

	h.Success(c, item)
}

// loadDetail 按详情接口的口径加载记录，供 Get 与写操作返回值共用。
func (h *CRUDHandler[T, L, C, U]) loadDetail(id uint) (*T, error) {
	query := h.DB
	// 允许业务自定义详情查询（例如 preload 关联）
	if h.BuildGetQuery != nil {
		query = h.BuildGetQuery(h.DB)
	}

	var item T
	if err := query.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	// 需要二次补齐/转换字段时走 AfterGet
	if h.AfterGet != nil {
		if err := h.AfterGet(&item); err != nil {
			return nil, err
		}
	}
	return &item, nil
}

// hasDetailLoader 判断是否配置了详情加载逻辑。
// 说明：未配置时写操作直接返回内存中的记录，避免为简单表多一次查询。
func (h *CRUDHandler[T, L, C, U]) hasDetailLoader() bool {
	return h.BuildGetQuery != nil || h.AfterGet != nil
}

// Create 创建记录
func (h *CRUDHandler[T, L, C, U]) Create(c *gin.Context) {
	var req C
	if err := h.BindJSON(c, &req); err != nil {
		return
	}
	// 未配置创建映射逻辑时无法执行标准 Create
	if h.NewModelFromCreate == nil {
		h.Error(c, "创建逻辑未配置")
		return
	}

	item, err := h.NewModelFromCreate(&req)
	// 请求转模型失败（例如参数不合法）直接返回错误
	if err != nil {
		h.Fail(c, err)
		return
	}

	successMsg := h.CreateSuccessMsg
	// 未配置时使用默认提示
	if successMsg == "" {
		successMsg = "创建成功"
	}

	if err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		// 创建前 hook 适合做跨字段校验或补充审计字段。
		if h.BeforeCreate != nil {
			if err := h.BeforeCreate(tx, item, &req); err != nil {
				return err
			}
		}
		// 唯一性校验与插入处于同一事务。
		if err := h.checkCreateUnique(tx, item); err != nil {
			return err
		}
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		// 创建后 hook 只在主记录成功插入后执行（例如写关联表）。
		if h.AfterCreate != nil {
			return h.AfterCreate(tx, item, &req)
		}
		return nil
	}); err != nil {
		h.Fail(c, err)
		return
	}
	if h.AfterCreateCommit != nil {
		// 缓存失效等外部副作用必须晚于事务提交，避免并发请求回填旧值。
		h.AfterCreateCommit(item, &req)
	}

	result, ok := h.reloadAfterWrite(c, getID(item), item)
	if !ok {
		return
	}
	h.SuccessWithMessage(c, successMsg, result)
}

// Update 更新记录
func (h *CRUDHandler[T, L, C, U]) Update(c *gin.Context) {
	id, err := h.ParseID(c)
	// ID 解析失败直接返回，由 BaseHandler 已输出错误响应
	if err != nil {
		return
	}

	var req U
	if err := h.BindJSON(c, &req); err != nil {
		return
	}
	// 未配置更新字段映射逻辑时无法执行标准 Update
	if h.BuildUpdates == nil {
		h.Error(c, "更新逻辑未配置")
		return
	}

	h.updateWithRequest(c, id, &req, func(tx *gorm.DB, existing *T) (map[string]interface{}, error) {
		return h.BuildUpdates(tx, &req, existing)
	}, h.defaultUpdateSuccessMsg())
}

// Delete 删除记录
func (h *CRUDHandler[T, L, C, U]) Delete(c *gin.Context) {
	id, err := h.ParseID(c)
	// ID 解析失败直接返回，由 BaseHandler 已输出错误响应
	if err != nil {
		return
	}

	h.deleteByIDs(c, []uint{id})
}

// DeleteBatch 批量删除。
// 说明：用于后台表格的批量操作，避免每个 handler 重复实现 ids 解析和 IN 删除。
func (h *CRUDHandler[T, L, C, U]) DeleteBatch(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := h.BindJSON(c, &req); err != nil {
		return
	}
	// ids 为空时没有意义，直接返回
	if len(req.IDs) == 0 {
		h.Error(c, "ids 不能为空")
		return
	}

	h.deleteByIDs(c, UniqueUints(req.IDs))
}

// deleteByIDs 执行标准删除事务。
// 说明：单个删除与批量删除共用，业务只需要写一套 hook；所有 id 必须全部命中，否则整体回滚。
func (h *CRUDHandler[T, L, C, U]) deleteByIDs(c *gin.Context, ids []uint) {
	successMsg := h.DeleteSuccessMsg
	// 未配置时使用默认提示
	if successMsg == "" {
		successMsg = "删除成功"
	}

	if err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		// 删除前 hook 适合做业务保护（例如禁止删除内置记录）和关联表清理，整体处理避免逐条 I/O。
		if h.BeforeDelete != nil {
			if err := h.BeforeDelete(tx, ids); err != nil {
				return err
			}
		}
		var item T
		result := tx.Where("id IN ?", ids).Delete(&item)
		if result.Error != nil {
			return result.Error
		}
		// 必须全部命中，否则返回不存在，避免调用方误以为全部删除成功。
		if result.RowsAffected != int64(len(ids)) {
			return gorm.ErrRecordNotFound
		}
		// 删除后 hook 只在主记录删除成功后执行。
		if h.AfterDelete != nil {
			return h.AfterDelete(tx, ids)
		}
		return nil
	}); err != nil {
		h.handleRecordError(c, err)
		return
	}
	if h.AfterDeleteCommit != nil {
		// 缓存失效等外部副作用必须晚于事务提交，避免并发请求回填旧值。
		h.AfterDeleteCommit(ids)
	}

	h.SuccessWithMessage(c, successMsg, nil)
}

// UpdateEnabled 更新启用状态。
// 说明：适用于存在 enabled 字段的场景，避免每个 handler 重复写启用/禁用逻辑。
func (h *CRUDHandler[T, L, C, U]) UpdateEnabled(c *gin.Context) {
	id, err := h.ParseID(c)
	// ID 解析失败直接返回，由 BaseHandler 已输出错误响应
	if err != nil {
		return
	}

	var req struct {
		Enabled *bool `json:"enabled" binding:"required"`
	}
	if err := h.BindJSON(c, &req); err != nil {
		return
	}

	field := h.EnabledField
	// 未配置时默认使用 enabled 字段
	if field == "" {
		field = "enabled"
	}

	successMsg := h.EnabledSuccessMsg
	// 未配置时使用默认提示
	if successMsg == "" {
		successMsg = "更新成功"
	}

	h.updateFields(c, id, map[string]interface{}{field: *req.Enabled}, successMsg)
}

// updateFields 使用指定字段更新记录，并复用标准更新生命周期。
// 主要服务于启用/禁用这类固定字段更新，避免绕过 AfterUpdate 和缓存失效逻辑。
func (h *CRUDHandler[T, L, C, U]) updateFields(
	c *gin.Context,
	id uint,
	updates map[string]interface{},
	successMsg string,
) {
	var req U
	h.updateWithRequest(c, id, &req, func(tx *gorm.DB, existing *T) (map[string]interface{}, error) {
		return updates, nil
	}, successMsg)
}

// updateWithRequest 执行标准更新事务。
// 所有更新入口都走这里，确保查询、字段更新、扩展逻辑、重载返回和错误响应一致。
func (h *CRUDHandler[T, L, C, U]) updateWithRequest(
	c *gin.Context,
	id uint,
	req *U,
	buildUpdates func(tx *gorm.DB, existing *T) (map[string]interface{}, error),
	successMsg string,
) {
	var updated T
	var updates map[string]interface{}
	if err := h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		q := tx
		// 允许业务自定义更新前的查询（例如 preload/锁）。
		if h.BuildUpdateQuery != nil {
			q = h.BuildUpdateQuery(tx)
		}
		if err := q.Where("id = ?", id).First(&updated).Error; err != nil {
			return err
		}

		var err error
		updates, err = buildUpdates(tx, &updated)
		if err != nil {
			return err
		}
		if err := h.checkUpdateUnique(tx, id, updates); err != nil {
			return err
		}

		// 更新前 hook 可以追加校验或改写 updates。
		if h.BeforeUpdate != nil {
			if err := h.BeforeUpdate(tx, id, &updated, req, updates); err != nil {
				return err
			}
		}
		if len(updates) > 0 {
			if err := tx.Model(&updated).Updates(updates).Error; err != nil {
				return err
			}
		}

		// 更新后 hook 只在字段更新成功后执行（例如同步关联表）。
		if h.AfterUpdate != nil {
			return h.AfterUpdate(tx, id, &updated, req)
		}
		return nil
	}); err != nil {
		h.handleRecordError(c, err)
		return
	}
	if h.AfterUpdateCommit != nil {
		// 缓存失效必须晚于事务提交，避免并发请求回填旧值。
		h.AfterUpdateCommit(id, &updated, req, updates)
	}

	result, ok := h.reloadAfterWrite(c, id, &updated)
	if !ok {
		return
	}
	h.SuccessWithMessage(c, successMsg, result)
}

// reloadAfterWrite 写操作提交后按详情口径重新加载返回值。
//
// 说明：放在提交之后执行，读取的是已落库的最终数据，且不会在事务未结束时占用第二条连接。
// 返回 false 表示已输出错误响应。
func (h *CRUDHandler[T, L, C, U]) reloadAfterWrite(c *gin.Context, id uint, fallback *T) (*T, bool) {
	// 未配置详情加载逻辑，或模型没有标准 ID 字段时，直接返回内存中的记录。
	if !h.hasDetailLoader() || id == 0 {
		return fallback, true
	}
	item, err := h.loadDetail(id)
	if err != nil {
		h.handleRecordError(c, err)
		return nil, false
	}
	return item, true
}

// checkCreateUnique 在创建前对 UniqueFields 逐项查重。
func (h *CRUDHandler[T, L, C, U]) checkCreateUnique(tx *gorm.DB, item *T) error {
	if len(h.UniqueFields) == 0 {
		return nil
	}

	stmt := &gorm.Statement{DB: tx}
	if err := stmt.Parse(item); err != nil {
		return err
	}
	for _, unique := range h.UniqueFields {
		field := stmt.Schema.LookUpField(unique.Column)
		// 列名写错属于编码错误，直接报错而不是静默跳过校验。
		if field == nil {
			return fmt.Errorf("唯一字段 %s 不存在于模型", unique.Column)
		}
		value, isZero := field.ValueOf(tx.Statement.Context, reflect.ValueOf(item).Elem())
		// 零值不参与查重，必填校验由请求绑定负责。
		if isZero {
			continue
		}
		if err := h.ensureUnique(tx, unique, value, 0); err != nil {
			return err
		}
	}
	return nil
}

// checkUpdateUnique 在更新前对本次涉及的唯一列查重，排除记录自身。
func (h *CRUDHandler[T, L, C, U]) checkUpdateUnique(tx *gorm.DB, id uint, updates map[string]interface{}) error {
	for _, unique := range h.UniqueFields {
		value, ok := updates[unique.Column]
		// 本次没有修改该列，无需查重。
		if !ok {
			continue
		}
		if err := h.ensureUnique(tx, unique, value, id); err != nil {
			return err
		}
	}
	return nil
}

// ensureUnique 查询是否已有其他记录占用该值，excludeID 为 0 表示不排除任何记录。
func (h *CRUDHandler[T, L, C, U]) ensureUnique(tx *gorm.DB, unique Unique, value interface{}, excludeID uint) error {
	query := unique.Column + " = ?"
	args := []interface{}{value}
	if excludeID != 0 {
		query += " AND id <> ?"
		args = append(args, excludeID)
	}

	exists, err := Exists(tx, new(T), query, args...)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	if unique.Err != nil {
		return unique.Err
	}
	return NewBizError(unique.Column + " 已存在")
}

// handleRecordError 统一处理 CRUD 记录级错误。
// 目前最重要的是把 gorm.ErrRecordNotFound 统一转成模块自定义 404 文案。
func (h *CRUDHandler[T, L, C, U]) handleRecordError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		h.NotFound(c, h.defaultNotFoundMsg())
		return
	}
	h.Fail(c, err)
}

// defaultNotFoundMsg 返回模块未配置时的默认不存在文案。
func (h *CRUDHandler[T, L, C, U]) defaultNotFoundMsg() string {
	if h.NotFoundMsg != "" {
		return h.NotFoundMsg
	}
	return "记录不存在"
}

// defaultUpdateSuccessMsg 返回模块未配置时的默认更新成功文案。
func (h *CRUDHandler[T, L, C, U]) defaultUpdateSuccessMsg() string {
	if h.UpdateSuccessMsg != "" {
		return h.UpdateSuccessMsg
	}
	return "更新成功"
}

// getID 从模型中尽力读取 ID 字段。
// 说明：为了保持 KISS，这里仅支持约定字段名为 ID、类型为 uint 的场景。
func getID[T any](item *T) uint {
	v := indirectValue(item)
	if !v.IsValid() {
		return 0
	}
	f := v.FieldByName("ID")
	if !f.IsValid() || f.Kind() != reflect.Uint {
		return 0
	}
	return uint(f.Uint())
}

// indirectValue 解引用指针，拿到最终的结构体 Value。
func indirectValue[T any](item *T) reflect.Value {
	v := reflect.ValueOf(item)
	if !v.IsValid() {
		return reflect.Value{}
	}
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	return v
}
