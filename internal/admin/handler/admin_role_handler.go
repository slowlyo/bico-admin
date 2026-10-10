package handler

import (
	"bico-admin/internal/admin/model"
	"bico-admin/internal/admin/service"
	"bico-admin/internal/pkg/crud"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 权限定义
var rolePerms = crud.NewCRUDPerms("system", "admin_role", "角色管理").WithExtra(
	crud.Permission{Key: "system:admin_role:permission", Label: "配置权限"},
)

// AdminRoleHandler 角色管理处理器
type AdminRoleHandler struct {
	crud.CRUDHandler[model.AdminRole, roleListReq, createRoleReq, updateRoleReq]
	cacheInvalidator service.AuthCacheInvalidator
}

// NewAdminRoleHandler 创建角色管理处理器。
//
// 说明：cacheInvalidator 为空时使用空实现，业务代码无需逐处判空。
func NewAdminRoleHandler(db *gorm.DB, cacheInvalidator service.AuthCacheInvalidator) *AdminRoleHandler {
	if cacheInvalidator == nil {
		cacheInvalidator = service.NoopAuthCacheInvalidator{}
	}
	h := &AdminRoleHandler{cacheInvalidator: cacheInvalidator}
	h.DB = db
	h.NotFoundMsg = "角色不存在"
	// 不包含 code 等内部标识列。
	h.SortableFields = []string{"id", "name", "enabled", "created_at", "updated_at"}
	h.UniqueFields = []crud.Unique{{Column: "name", Err: crud.NewBizError("角色名称已存在")}}

	h.BuildListQuery = func(db *gorm.DB, req *roleListReq) *gorm.DB {
		query := db.Model(&model.AdminRole{})
		if req.Name != "" {
			query = query.Where("name LIKE ?", "%"+req.Name+"%")
		}
		if req.Description != "" {
			query = query.Where("description LIKE ?", "%"+req.Description+"%")
		}
		if req.Enabled != nil {
			query = query.Where("enabled = ?", *req.Enabled)
		}
		return query
	}

	h.AfterList = func(items []model.AdminRole) error {
		roleIDs := make([]uint, 0, len(items))
		for i := range items {
			roleIDs = append(roleIDs, items[i].ID)
		}
		permsMap, err := h.getPermsMap(db, roleIDs)
		if err != nil {
			return err
		}
		for i := range items {
			items[i].Permissions = permsMap[items[i].ID]
			items[i].System = items[i].Code == model.SuperAdminRoleCode
		}
		return nil
	}

	// 详情与写操作返回值共用，创建和更新后会按此口径重新加载。
	h.AfterGet = func(item *model.AdminRole) error {
		perms, err := h.getPerms(db, item.ID)
		if err != nil {
			return err
		}
		item.Permissions = perms
		item.System = item.Code == model.SuperAdminRoleCode
		return nil
	}

	h.NewModelFromCreate = func(req *createRoleReq) (*model.AdminRole, error) {
		code, err := generateRoleCode()
		if err != nil {
			// 内部标识生成失败时禁止创建不完整角色。
			return nil, err
		}
		return &model.AdminRole{
			Name:        req.Name,
			Code:        code,
			Description: req.Description,
			Enabled:     req.Enabled == nil || *req.Enabled,
		}, nil
	}
	h.AfterCreate = func(tx *gorm.DB, item *model.AdminRole, req *createRoleReq) error {
		return h.savePerms(tx, item.ID, req.Permissions)
	}

	h.BuildUpdates = func(tx *gorm.DB, req *updateRoleReq, existing *model.AdminRole) (map[string]interface{}, error) {
		updates := map[string]interface{}{}
		if req.Name != "" {
			updates["name"] = req.Name
		}
		if req.Description != "" {
			updates["description"] = req.Description
		}
		if req.Enabled != nil {
			updates["enabled"] = *req.Enabled
		}
		return updates, nil
	}
	// 保护性校验放在 BeforeUpdate，UpdateEnabled 等不经过 BuildUpdates 的入口同样生效。
	h.BeforeUpdate = func(tx *gorm.DB, id uint, existing *model.AdminRole, req *updateRoleReq, updates map[string]interface{}) error {
		if existing.Code == model.SuperAdminRoleCode {
			// 超级管理员角色必须保持启用且不可改名。
			return crud.NewBizError("超级管理员角色不可修改")
		}
		return nil
	}
	h.AfterUpdateCommit = func(id uint, existing *model.AdminRole, req *updateRoleReq, updates map[string]interface{}) {
		// 角色启用状态会直接影响用户权限结果；提交后再失效该角色下用户缓存，
		// 避免并发请求在提交前把旧权限回填进缓存。角色与用户的关联不随更新变化，提交后仍可定位用户。
		if _, enabledChanged := updates["enabled"]; enabledChanged {
			h.cacheInvalidator.InvalidateRoleUsersPermissionCache(id)
		}
	}

	h.BeforeDelete = func(tx *gorm.DB, ids []uint) error {
		var count int64
		if err := tx.Model(&model.AdminRole{}).
			Where("id IN ? AND code = ?", ids, model.SuperAdminRoleCode).
			Count(&count).Error; err != nil {
			// 保留角色查询失败时不执行删除。
			return err
		}
		if count > 0 {
			// 删除保留角色会导致系统失去最高权限入口。
			return crud.NewBizError("超级管理员角色不可删除")
		}
		// 必须在清理关联之前失效缓存：关联清理后就无法再定位该角色下的用户，
		// 所以这里不能改成提交后失效。
		h.cacheInvalidator.InvalidateRolesUsersPermissionCache(ids)
		if err := tx.Where("role_id IN ?", ids).Delete(&model.AdminRolePermission{}).Error; err != nil {
			return err
		}
		return tx.Where("role_id IN ?", ids).Delete(&model.AdminUserRole{}).Error
	}

	return h
}

func (h *AdminRoleHandler) ModuleConfig() crud.ModuleConfig {
	return crud.ModuleConfig{
		Name:             "admin_role",
		Group:            "/admin-roles",
		Description:      "角色管理",
		ParentPermission: PermSystemManage,
		Permissions:      rolePerms.Tree,
		Routes: rolePerms.RoutesWithExtra(
			crud.AuthRoute("GET", "/all", "GetAll"),
			crud.AuthRoute("GET", "/permissions", "GetAllPermissions"),
			crud.Route{Method: "GET", Path: "/:id/permissions", Handler: "GetPermissions", Permission: rolePerms.List},
			crud.Route{Method: "PUT", Path: "/:id/permissions", Handler: "UpdatePermissions", Permission: "system:admin_role:permission"},
		),
		Swagger: crud.SwaggerConfig{
			Model:         model.AdminRole{},
			ListRequest:   roleListReq{},
			CreateRequest: createRoleReq{},
			UpdateRequest: updateRoleReq{},
		},
	}
}

// 请求结构
type (
	roleListReq struct {
		Name        string `form:"name"`
		Description string `form:"description"`
		Enabled     *bool  `form:"enabled"`
	}
	createRoleReq struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		Enabled     *bool    `json:"enabled"`
		Permissions []string `json:"permissions"`
	}
	updateRoleReq struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Enabled     *bool  `json:"enabled"`
	}
	updateRolePermReq struct {
		Permissions []string `json:"permissions"`
	}
)

// GetPermissions 获取角色权限。
// @Summary 获取角色权限
// @Description 获取指定角色已配置的权限 key 列表
// @Tags 角色管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色 ID"
// @Success 200 {object} adminResponse{data=rolePermissionsResponse}
// @Router /admin-roles/{id}/permissions [get]
func (h *AdminRoleHandler) GetPermissions(c *gin.Context) {
	id, err := h.ParseID(c)
	if err != nil {
		return
	}

	perms, err := h.getPerms(h.DB, id)
	if err != nil {
		h.Fail(c, err)
		return
	}

	h.Success(c, gin.H{"permissions": perms})
}

// UpdatePermissions 更新角色权限。
// @Summary 更新角色权限
// @Description 覆盖指定角色的权限 key 列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色 ID"
// @Param body body updateRolePermReq true "权限列表"
// @Success 200 {object} adminResponse
// @Router /admin-roles/{id}/permissions [put]
func (h *AdminRoleHandler) UpdatePermissions(c *gin.Context) {
	id, err := h.ParseID(c)
	if err != nil {
		return
	}

	var req updateRolePermReq
	if err := h.BindJSON(c, &req); err != nil {
		return
	}
	// permissions 字段缺失时直接返回，避免误清空权限。
	if req.Permissions == nil {
		h.Error(c, "permissions 字段不能为空")
		return
	}
	// 只需要校验角色是否存在，不需要加载权限
	var role model.AdminRole
	if !h.QueryOne(c, h.DB.Where("id = ?", id), &role, "角色不存在") {
		return
	}
	if role.Code == model.SuperAdminRoleCode {
		// 保留角色动态拥有全部权限，无需持久化权限明细。
		h.Error(c, "超级管理员角色权限不可修改")
		return
	}

	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		// 先清空旧权限，再写入新权限，任一步失败都回滚。
		if err := tx.Where("role_id = ?", id).Delete(&model.AdminRolePermission{}).Error; err != nil {
			return err
		}
		return h.savePerms(tx, id, req.Permissions)
	}); err != nil {
		h.Fail(c, err)
		return
	}
	// 提交后再失效缓存，避免并发请求在提交前把旧权限回填进缓存。
	h.cacheInvalidator.InvalidateRoleUsersPermissionCache(id)
	h.SuccessWithMessage(c, "权限配置成功", nil)
}

// GetAllPermissions 获取完整权限树。
// @Summary 获取完整权限树
// @Description 获取后台所有菜单和按钮权限
// @Tags 角色管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} adminResponse{data=[]permissionDocItem}
// @Router /admin-roles/permissions [get]
func (h *AdminRoleHandler) GetAllPermissions(c *gin.Context) {
	h.Success(c, crud.GetAllPermissions())
}

// GetAll 获取全部启用角色。
// @Summary 获取全部启用角色
// @Description 获取下拉选择使用的启用角色列表
// @Tags 角色管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} adminResponse{data=[]adminRoleDocItem}
// @Router /admin-roles/all [get]
func (h *AdminRoleHandler) GetAll(c *gin.Context) {
	var roles []model.AdminRole
	if err := h.DB.Where("enabled = ?", true).Find(&roles).Error; err != nil {
		h.Fail(c, err)
		return
	}
	for i := range roles {
		roles[i].System = roles[i].Code == model.SuperAdminRoleCode
	}
	h.Success(c, roles)
}

// generateRoleCode 生成仅供系统内部使用的角色标识。
// 随机标识不依赖角色名称，避免重命名影响既有关系。
func generateRoleCode() (string, error) {
	random := make([]byte, 8)
	if _, err := rand.Read(random); err != nil {
		// 安全随机源不可用时返回错误，避免产生可预测或重复标识。
		return "", fmt.Errorf("生成角色标识失败: %w", err)
	}
	return "role_" + hex.EncodeToString(random), nil
}

// 私有方法
func (h *AdminRoleHandler) getPerms(db *gorm.DB, roleID uint) ([]string, error) {
	var perms []string
	err := db.Model(&model.AdminRolePermission{}).Where("role_id = ?", roleID).Pluck("permission", &perms).Error
	return perms, err
}

func (h *AdminRoleHandler) getPermsMap(db *gorm.DB, roleIDs []uint) (map[uint][]string, error) {
	permsMap := make(map[uint][]string)
	if len(roleIDs) == 0 {
		return permsMap, nil
	}

	type row struct {
		RoleID     uint
		Permission string
	}
	var rows []row
	err := db.Model(&model.AdminRolePermission{}).
		Select("role_id, permission").
		Where("role_id IN ?", roleIDs).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		permsMap[r.RoleID] = append(permsMap[r.RoleID], r.Permission)
	}
	return permsMap, nil
}

func (h *AdminRoleHandler) savePerms(tx *gorm.DB, roleID uint, perms []string) error {
	if len(perms) == 0 {
		return nil
	}

	validPerms := make(map[string]struct{})
	for _, key := range crud.GetAllPermissionKeys() {
		validPerms[key] = struct{}{}
	}

	dedupPerms := make([]string, 0, len(perms))
	seen := make(map[string]struct{}, len(perms))
	for _, raw := range perms {
		p := strings.TrimSpace(raw)
		// 空权限直接拒绝，避免脏数据入库。
		if p == "" {
			return crud.NewBizError("权限标识不能为空")
		}
		// 非法权限直接拒绝，避免越权配置。
		if _, ok := validPerms[p]; !ok {
			return crud.NewBizError("存在无效权限标识: " + p)
		}
		// 重复权限自动去重，降低唯一约束冲突概率。
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		dedupPerms = append(dedupPerms, p)
	}

	items := make([]model.AdminRolePermission, 0, len(dedupPerms))
	for _, p := range dedupPerms {
		items = append(items, model.AdminRolePermission{RoleID: roleID, Permission: p})
	}
	return tx.CreateInBatches(items, 100).Error
}

var _ crud.Module = (*AdminRoleHandler)(nil)
