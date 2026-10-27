package handler

import (
	"context"
	"strconv"
	"strings"

	"bico-admin/internal/admin/model"
	"bico-admin/internal/admin/service"
	"bico-admin/internal/pkg/crud"
	"bico-admin/internal/pkg/password"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 权限定义
var userPerms = crud.NewCRUDPerms("system", "admin_user", "用户管理")

// AdminUserHandler 用户管理处理器
type AdminUserHandler struct {
	crud.CRUDHandler[model.AdminUser, userListReq, createUserReq, updateUserReq]
	cacheInvalidator service.AuthCacheInvalidator
}

// NewAdminUserHandler 创建用户管理处理器。
//
// 说明：cacheInvalidator 为空时使用空实现，业务代码无需逐处判空。
func NewAdminUserHandler(db *gorm.DB, cacheInvalidator service.AuthCacheInvalidator) *AdminUserHandler {
	if cacheInvalidator == nil {
		cacheInvalidator = service.NoopAuthCacheInvalidator{}
	}
	h := &AdminUserHandler{cacheInvalidator: cacheInvalidator}
	h.DB = db
	h.NotFoundMsg = "用户不存在"
	// 不包含 password 等敏感列，避免按密码哈希排序推断其内容。
	h.SortableFields = []string{"id", "username", "name", "enabled", "created_at", "updated_at"}
	// 用户名是全局唯一的登录标识。
	h.UniqueFields = []crud.Unique{{Column: "username", Err: service.ErrUsernameExists}}

	h.BuildListQuery = func(db *gorm.DB, req *userListReq) *gorm.DB {
		query := db.Model(&model.AdminUser{}).Preload("Roles")
		if req.Username != "" {
			query = query.Where("username LIKE ?", "%"+req.Username+"%")
		}
		if req.Name != "" {
			query = query.Where("name LIKE ?", "%"+req.Name+"%")
		}
		if req.Enabled != nil {
			query = query.Where("enabled = ?", *req.Enabled)
		}
		roleIDs := parseUserRoleIDs(req.RoleIDs)
		if len(roleIDs) > 0 {
			// 角色筛选走关联表 EXISTS，避免 JOIN 影响用户列表分页总数。
			query = query.Where(
				"EXISTS (SELECT 1 FROM admin_user_roles WHERE admin_user_roles.user_id = admin_users.id AND admin_user_roles.role_id IN ?)",
				roleIDs,
			)
		}
		return query
	}

	// 详情与写操作返回值都需要带 Roles。
	h.BuildGetQuery = func(db *gorm.DB) *gorm.DB {
		return db.Model(&model.AdminUser{}).Preload("Roles")
	}

	h.NewModelFromCreate = func(req *createUserReq) (*model.AdminUser, error) {
		username := strings.TrimSpace(req.Username)
		if username == "" {
			// 去除首尾空格后为空的用户名不能用于登录。
			return nil, service.ErrUsernameRequired
		}
		hashed, err := password.Hash(req.Password)
		if err != nil {
			return nil, err
		}

		return &model.AdminUser{
			Username: username,
			Password: hashed,
			Name:     req.Name,
			Avatar:   req.Avatar,
			Enabled:  req.Enabled == nil || *req.Enabled,
		}, nil
	}
	h.AfterCreate = func(tx *gorm.DB, item *model.AdminUser, req *createUserReq) error {
		return h.syncRoles(tx, item, req.RoleIDs)
	}

	h.BuildUpdates = func(tx *gorm.DB, req *updateUserReq, existing *model.AdminUser) (map[string]interface{}, error) {
		updates := map[string]interface{}{}
		if req.Username != nil {
			username := strings.TrimSpace(*req.Username)
			if username == "" {
				// 显式提交空用户名时拒绝保存。
				return nil, service.ErrUsernameRequired
			}
			updates["username"] = username
		}
		if req.Name != "" {
			updates["name"] = req.Name
		}
		if req.Avatar != "" {
			updates["avatar"] = req.Avatar
		}
		if req.Password != "" {
			hashed, err := password.Hash(req.Password)
			if err != nil {
				return nil, err
			}
			updates["password"] = hashed
			updates["token_version"] = gorm.Expr("token_version + 1")
		}
		if req.Enabled != nil {
			updates["enabled"] = *req.Enabled
		}
		return updates, nil
	}
	h.BeforeUpdate = func(tx *gorm.DB, id uint, existing *model.AdminUser, req *updateUserReq, updates map[string]interface{}) error {
		// 禁用账号前，至少保留一个启用的超级管理员，避免后台永久失去管理入口。
		if enabled, ok := updates["enabled"].(bool); ok && !enabled {
			return h.ensureSuperAdminRemains(tx, existing.ID)
		}
		return nil
	}
	h.AfterUpdate = func(tx *gorm.DB, id uint, existing *model.AdminUser, req *updateUserReq) error {
		// 请求未携带角色时保持原有角色关联不变。
		if req.RoleIDs == nil {
			return nil
		}
		return h.syncRoles(tx, existing, req.RoleIDs)
	}
	h.AfterUpdateCommit = func(id uint, existing *model.AdminUser, req *updateUserReq, updates map[string]interface{}) {
		// 提交后再清缓存，防止并发请求在提交前把旧值回填进缓存。
		// 依据 updates 而不是 req 判断变更，这样 UpdateEnabled 路径也能正确失效。
		_, enabledChanged := updates["enabled"]
		if req.RoleIDs != nil {
			h.cacheInvalidator.InvalidateUserPermissionCache(existing.ID)
		}
		// 角色或启用状态变化都会影响状态缓存。
		if req.RoleIDs != nil || enabledChanged {
			h.cacheInvalidator.InvalidateUserStatusCache(existing.ID)
		}
		if _, passwordChanged := updates["password"]; passwordChanged {
			h.cacheInvalidator.InvalidateUserTokenVersionCache(existing.ID)
		}
	}

	h.BeforeDelete = func(tx *gorm.DB, ids []uint) error {
		if err := h.ensureSuperAdminsRemain(tx, ids); err != nil {
			return err
		}
		// 删除用户前先清理角色关联，后续任一步失败都会随事务回滚。
		return tx.Where("user_id IN ?", ids).Delete(&model.AdminUserRole{}).Error
	}
	h.AfterDeleteCommit = func(ids []uint) {
		// 删除提交后再批量失效权限与状态缓存，单个与批量删除共用。
		h.cacheInvalidator.InvalidateUsersAuthCache(ids)
	}

	return h
}

// Create 创建用户并在请求上下文中注入操作者，供角色绑定授权校验使用。
func (h *AdminUserHandler) Create(c *gin.Context) {
	operatorID, ok := mustOperatorUserID(c)
	if !ok {
		return
	}
	c.Request = c.Request.WithContext(withOperatorUserID(c.Request.Context(), operatorID))
	h.CRUDHandler.Create(c)
}

// Update 更新用户并在请求上下文中注入操作者，供角色绑定授权校验使用。
func (h *AdminUserHandler) Update(c *gin.Context) {
	operatorID, ok := mustOperatorUserID(c)
	if !ok {
		return
	}
	c.Request = c.Request.WithContext(withOperatorUserID(c.Request.Context(), operatorID))
	h.CRUDHandler.Update(c)
}

// parseUserRoleIDs 解析用户列表角色筛选参数，忽略非法片段以保持筛选接口容错。
func parseUserRoleIDs(value string) []uint {
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	ids := make([]uint, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		// 空片段来自多余逗号，不参与筛选。
		if trimmed == "" {
			continue
		}
		id, err := strconv.ParseUint(trimmed, 10, 64)
		// 非法角色 ID 直接忽略，避免一个脏值导致整个列表不可用。
		if err != nil || id == 0 {
			continue
		}
		ids = append(ids, uint(id))
	}
	return crud.UniqueUints(ids)
}

func (h *AdminUserHandler) ModuleConfig() crud.ModuleConfig {
	return crud.ModuleConfig{
		Name:             "admin_user",
		Group:            "/admin-users",
		Description:      "用户管理",
		ParentPermission: PermSystemManage,
		Permissions:      userPerms.Tree,
		Routes:           userPerms.Routes(),
		Swagger: crud.SwaggerConfig{
			Model:         model.AdminUser{},
			ListRequest:   userListReq{},
			CreateRequest: createUserReq{},
			UpdateRequest: updateUserReq{},
		},
	}
}

// 请求结构
type (
	userListReq struct {
		Username string `form:"username"`
		Name     string `form:"name"`
		Enabled  *bool  `form:"enabled"`
		RoleIDs  string `form:"role_ids"`
	}
	createUserReq struct {
		Username string `json:"username" binding:"required,max=64"`
		Password string `json:"password" binding:"required,min=8"`
		Name     string `json:"name"`
		Avatar   string `json:"avatar"`
		Enabled  *bool  `json:"enabled"`
		RoleIDs  []uint `json:"role_ids"`
	}
	updateUserReq struct {
		Username *string `json:"username" binding:"omitempty,max=64"`
		Name     string  `json:"name"`
		Avatar   string  `json:"avatar"`
		Enabled  *bool   `json:"enabled"`
		RoleIDs  []uint  `json:"role_ids"`
		Password string  `json:"password" binding:"omitempty,min=8"`
	}
)

func (h *AdminUserHandler) syncRoles(tx *gorm.DB, user *model.AdminUser, roleIDs []uint) error {
	var ctx context.Context
	if tx.Statement != nil {
		ctx = tx.Statement.Context
	}
	operatorID := operatorUserIDFromContext(ctx)
	if operatorID == 0 {
		// 角色绑定必须能追溯到当前登录用户，缺失时拒绝写入，避免绕过授权校验。
		return crud.NewBizError("未授权")
	}
	if err := service.AssertCanAssignRoles(tx, operatorID, roleIDs); err != nil {
		return err
	}
	if err := h.ensureRoleChangeKeepsSuperAdmin(tx, user.ID, roleIDs); err != nil {
		return err
	}
	// 先清空旧关联，确保写入结果与请求保持一致。
	if err := tx.Model(user).Association("Roles").Clear(); err != nil {
		return err
	}
	if len(roleIDs) == 0 {
		return nil
	}

	uniqueRoleIDs := crud.UniqueUints(roleIDs)
	var roles []*model.AdminRole
	if err := tx.Where("id IN ?", uniqueRoleIDs).Find(&roles).Error; err != nil {
		return err
	}
	// 如果查询数量不一致，说明请求中存在无效角色 ID。
	if len(roles) != len(uniqueRoleIDs) {
		return crud.NewBizError("存在无效角色 ID")
	}
	if err := tx.Model(user).Association("Roles").Append(roles); err != nil {
		return err
	}
	return nil
}

// ensureRoleChangeKeepsSuperAdmin 校验角色调整不会移除最后一个超级管理员。
func (h *AdminUserHandler) ensureRoleChangeKeepsSuperAdmin(db *gorm.DB, userID uint, roleIDs []uint) error {
	var currentSuperCount int64
	if err := db.Table("admin_user_roles").
		Joins("JOIN admin_roles ON admin_user_roles.role_id = admin_roles.id").
		Where("admin_user_roles.user_id = ? AND admin_roles.code = ?", userID, model.SuperAdminRoleCode).
		Count(&currentSuperCount).Error; err != nil {
		return err
	}
	if currentSuperCount == 0 {
		// 普通用户的角色调整不影响超级管理员连续性。
		return nil
	}

	var requestedSuperCount int64
	if len(roleIDs) > 0 {
		if err := db.Model(&model.AdminRole{}).
			Where("id IN ? AND code = ?", crud.UniqueUints(roleIDs), model.SuperAdminRoleCode).
			Count(&requestedSuperCount).Error; err != nil {
			return err
		}
	}
	if requestedSuperCount > 0 {
		// 请求仍保留超级管理员角色，无需额外校验。
		return nil
	}
	return h.ensureSuperAdminRemains(db, userID)
}

// ensureSuperAdminRemains 确保目标用户之外仍有启用的超级管理员。
func (h *AdminUserHandler) ensureSuperAdminRemains(db *gorm.DB, userID uint) error {
	return h.ensureSuperAdminsRemain(db, []uint{userID})
}

// ensureSuperAdminsRemain 确保排除目标用户后仍有启用的超级管理员。
func (h *AdminUserHandler) ensureSuperAdminsRemain(db *gorm.DB, excludedUserIDs []uint) error {
	var targetSuperCount int64
	if err := db.Table("admin_user_roles").
		Joins("JOIN admin_roles ON admin_user_roles.role_id = admin_roles.id").
		Where("admin_user_roles.user_id IN ? AND admin_roles.code = ?", excludedUserIDs, model.SuperAdminRoleCode).
		Count(&targetSuperCount).Error; err != nil {
		return err
	}
	if targetSuperCount == 0 {
		// 目标不是超级管理员时可直接变更。
		return nil
	}

	var remaining int64
	if err := db.Table("admin_user_roles").
		Joins("JOIN admin_roles ON admin_user_roles.role_id = admin_roles.id").
		Joins("JOIN admin_users ON admin_user_roles.user_id = admin_users.id").
		Where("admin_roles.code = ? AND admin_roles.enabled = ? AND admin_users.enabled = ? AND admin_users.id NOT IN ?", model.SuperAdminRoleCode, true, true, excludedUserIDs).
		Count(&remaining).Error; err != nil {
		return err
	}
	if remaining == 0 {
		return crud.NewBizError("必须保留至少一个启用的超级管理员")
	}
	return nil
}

var _ crud.Module = (*AdminUserHandler)(nil)
