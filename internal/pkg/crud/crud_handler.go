package crud

import (
	"errors"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CRUDHandler 通用 CRUD Handler，通过配置 hook 来适配不同业务模型。
// 目标：让单个 handler 文件只保留少量业务差异代码。
// T: 模型类型；L: List 查询参数类型；C: Create 请求体类型；U: Update 请求体类型。
//
// 说明：
// - List/Get/Create/Update/Delete 走统一的参数绑定、错误处理与响应
// - 差异点通过 BuildQuery/After/TxHook 等函数注入
// - 默认使用 gorm 的 First/Delete/Updates
// - 需要 preload/关联清理/二次加工时在 hook 内实现
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

	// BuildListQuery 构建列表查询（必须提供）
	BuildListQuery func(db *gorm.DB, req *L) *gorm.DB
	// AfterList 列表查询后对结果二次处理（可选）
	AfterList func(items []T) error

	// BuildGetQuery 构建详情查询（可选，为 nil 则使用 DB.Model(&T{})）
	BuildGetQuery func(db *gorm.DB) *gorm.DB
	// AfterGet 查询详情后对结果二次处理（可选）
	AfterGet func(item *T) error

	// NewModelFromCreate Create 请求转模型（必须提供）
	NewModelFromCreate func(req *C) (*T, error)
	// CreateInTx 创建事务内的额外逻辑（可选）
	CreateInTx func(tx *gorm.DB, item *T, req *C) error
	// ReloadAfterCreate 创建后重新加载返回数据（可选）
	ReloadAfterCreate func(tx *gorm.DB, id uint, item *T) error

	// BuildUpdateQuery 构建更新时查询（可选，为 nil 则 tx.First）
	BuildUpdateQuery func(tx *gorm.DB) *gorm.DB
	// BuildUpdates Update 请求转 updates map（必须提供）
	BuildUpdates func(req *U, existing *T) (map[string]interface{}, error)
	// UpdateInTx 更新事务内额外逻辑（可选）
	UpdateInTx func(tx *gorm.DB, id uint, existing *T, req *U) error
	// ReloadAfterUpdate 更新后重新加载返回数据（可选）
	ReloadAfterUpdate func(tx *gorm.DB, id uint, existing *T) error

	// DeleteInTx 删除事务内的额外逻辑（可选）
	DeleteInTx func(tx *gorm.DB, id uint) error
	// DeleteBatchInTx 批量删除事务内的额外逻辑（可选）
	DeleteBatchInTx func(tx *gorm.DB, ids []uint) error
}

// List 获取列表
func (h *CRUDHandler[T, L, C, U]) List(c *gin.Context) {
	var req L
	_ = h.BindQuery(c, &req)
	// 未配置列表查询函数时无法执行标准 List
	if h.BuildListQuery == nil {
		h.Error(c, "列表查询未配置")
		return
	}

	query := h.BuildListQuery(h.DB, &req)
	var items []T
	QueryListWithHook(&h.BaseHandler, c, query, &items, h.AfterList)
}

// Get 获取详情
func (h *CRUDHandler[T, L, C, U]) Get(c *gin.Context) {
	id, err := h.ParseID(c)
	// ID 解析失败直接返回，由 BaseHandler 已输出错误响应
	if err != nil {
		return
	}

	var item T
	query := h.DB
	// 允许业务自定义详情查询（例如 preload 关联）
	if h.BuildGetQuery != nil {
		query = h.BuildGetQuery(h.DB)
	}

	notFoundMsg := h.NotFoundMsg
	// 未配置时提供默认提示
	if notFoundMsg == "" {
		notFoundMsg = "记录不存在"
	}

	if !h.QueryOne(c, query.Where("id = ?", id), &item, notFoundMsg) {
		return
	}
	// 需要二次补齐/转换字段时走 AfterGet
	if h.AfterGet != nil {
		if err := h.AfterGet(&item); err != nil {
			h.Error(c, err.Error())
			return
		}
	}

	h.Success(c, item)
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
	// 业务校验失败（例如唯一性）直接返回错误
	if err != nil {
		h.Error(c, err.Error())
		return
	}

	successMsg := h.CreateSuccessMsg
	// 未配置时使用默认提示
	if successMsg == "" {
		successMsg = "创建成功"
	}

	h.ExecTx(c, h.DB, func(tx *gorm.DB) error {
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		// 事务内扩展创建逻辑（例如写关联表）
		if h.CreateInTx != nil {
			if err := h.CreateInTx(tx, item, &req); err != nil {
				return err
			}
		}
		// 需要返回 preload 后的数据时，在这里重新加载
		if h.ReloadAfterCreate != nil {
			return h.ReloadAfterCreate(tx, getID(item), item)
		}
		return nil
	}, successMsg, item)
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

	var updated T
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		q := tx
		// 允许业务自定义更新前的查询（例如 preload/锁）
		if h.BuildUpdateQuery != nil {
			q = h.BuildUpdateQuery(tx)
		}
		if err := q.Where("id = ?", id).First(&updated).Error; err != nil {
			return err
		}

		updates, err := h.BuildUpdates(&req, &updated)
		// 业务校验失败（例如唯一性）直接回滚事务
		if err != nil {
			return err
		}
		// updates 为空时跳过更新，仍允许执行 UpdateInTx/ReloadAfterUpdate
		if len(updates) > 0 {
			if err := tx.Model(&updated).Updates(updates).Error; err != nil {
				return err
			}
		}

		// 事务内扩展更新逻辑（例如同步关联表）
		if h.UpdateInTx != nil {
			if err := h.UpdateInTx(tx, id, &updated, &req); err != nil {
				return err
			}
		}

		// 需要返回 preload 后的数据时，在这里重新加载
		if h.ReloadAfterUpdate != nil {
			if err := h.ReloadAfterUpdate(tx, id, &updated); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		// 记录不存在时返回 404，保持与 QueryOne 行为一致
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFoundMsg := h.NotFoundMsg
			if notFoundMsg == "" {
				notFoundMsg = "记录不存在"
			}
			h.NotFound(c, notFoundMsg)
			return
		}
		h.Error(c, err.Error())
		return
	}

	successMsg := h.UpdateSuccessMsg
	// 未配置时使用默认提示
	if successMsg == "" {
		successMsg = "更新成功"
	}
	h.SuccessWithMessage(c, successMsg, updated)
}

// Delete 删除记录
func (h *CRUDHandler[T, L, C, U]) Delete(c *gin.Context) {
	id, err := h.ParseID(c)
	// ID 解析失败直接返回，由 BaseHandler 已输出错误响应
	if err != nil {
		return
	}

	successMsg := h.DeleteSuccessMsg
	// 未配置时使用默认提示
	if successMsg == "" {
		successMsg = "删除成功"
	}

	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		// 事务内扩展删除逻辑（例如清理关联表）
		if h.DeleteInTx != nil {
			if err := h.DeleteInTx(tx, id); err != nil {
				return err
			}
		}
		var item T
		result := tx.Delete(&item, id)
		if result.Error != nil {
			return result.Error
		}
		// 软删场景中 rowsAffected=0 代表记录不存在。
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	}); err != nil {
		// 删除不存在记录时返回 404，保持与 Get/Update 行为一致。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFoundMsg := h.NotFoundMsg
			if notFoundMsg == "" {
				notFoundMsg = "记录不存在"
			}
			h.NotFound(c, notFoundMsg)
			return
		}
		h.Error(c, err.Error())
		return
	}

	h.SuccessWithMessage(c, successMsg, nil)
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

	ids := UniqueUints(req.IDs)
	successMsg := h.DeleteSuccessMsg
	// 未配置时使用默认提示
	if successMsg == "" {
		successMsg = "删除成功"
	}

	h.ExecTx(c, h.DB, func(tx *gorm.DB) error {
		// 批量场景优先使用专用 hook，避免在循环中做 I/O
		if h.DeleteBatchInTx != nil {
			if err := h.DeleteBatchInTx(tx, ids); err != nil {
				return err
			}
		}
		var item T
		return tx.Where("id IN ?", ids).Delete(&item).Error
	}, successMsg, nil)
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
		Enabled bool `json:"enabled" binding:"required"`
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

	var updated T
	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		q := tx
		// 允许业务自定义查询（例如 preload 关联）
		if h.BuildUpdateQuery != nil {
			q = h.BuildUpdateQuery(tx)
		}
		if err := q.Where("id = ?", id).First(&updated).Error; err != nil {
			return err
		}

		if err := tx.Model(&updated).Update(field, req.Enabled).Error; err != nil {
			return err
		}

		// 需要返回 preload 后的数据时，在这里重新加载
		if h.ReloadAfterUpdate != nil {
			return h.ReloadAfterUpdate(tx, id, &updated)
		}
		return nil
	}); err != nil {
		// 记录不存在时返回 404，保持与 QueryOne 行为一致
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFoundMsg := h.NotFoundMsg
			if notFoundMsg == "" {
				notFoundMsg = "记录不存在"
			}
			h.NotFound(c, notFoundMsg)
			return
		}
		h.Error(c, err.Error())
		return
	}

	h.SuccessWithMessage(c, successMsg, updated)
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
