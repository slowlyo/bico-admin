package handler

import (
	"bico-admin/internal/admin/model"
	"bico-admin/internal/pkg/crud"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 权限定义
var rolePerms = crud.NewCRUDPerms("system", "admin_role", "角色管理").WithExtra(
	crud.Permission{Key: "system:admin_role:permission", Label: "配置权限"},
)

// AdminRoleHandler 角色管理处理器
type AdminRoleHandler struct {
	crud.BaseHandler
	db *gorm.DB
}

func NewAdminRoleHandler(db *gorm.DB) *AdminRoleHandler {
	return &AdminRoleHandler{db: db}
}

func (h *AdminRoleHandler) ModuleConfig() crud.ModuleConfig {
	return crud.ModuleConfig{
		Name:             "admin_role",
		Group:            "/admin-roles",
		ParentPermission: PermSystemManage,
		Permissions:      rolePerms.Tree,
		Routes: rolePerms.RoutesWithExtra(
			crud.Route{Method: "GET", Path: "/all", Handler: "GetAll", Permission: rolePerms.List},
			crud.Route{Method: "GET", Path: "/permissions", Handler: "GetAllPermissions", Permission: rolePerms.List},
			crud.Route{Method: "PUT", Path: "/:id/permissions", Handler: "UpdatePermissions", Permission: "system:admin_role:permission"},
		),
	}
}

// 请求结构
type (
	roleListReq struct {
		Name    string `form:"name"`
		Code    string `form:"code"`
		Enabled *bool  `form:"enabled"`
	}
	createRoleReq struct {
		Name, Code, Description string
		Enabled                 *bool
		Permissions             []string
	}
	updateRoleReq struct {
		Name, Description string
		Enabled           *bool
	}
	updateRolePermReq struct {
		Permissions []string `json:"permissions" binding:"required"`
	}
)

func (h *AdminRoleHandler) List(c *gin.Context) {
	var req roleListReq
	h.BindQuery(c, &req)

	query := h.db.Model(&model.AdminRole{})
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		query = query.Where("code LIKE ?", "%"+req.Code+"%")
	}
	if req.Enabled != nil {
		query = query.Where("enabled = ?", *req.Enabled)
	}

	pg := h.GetPagination(c)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		h.Error(c, err.Error())
		return
	}

	if orderBy := pg.GetOrderBy(); orderBy != "" {
		query = query.Order(orderBy)
	} else {
		query = query.Order("created_at DESC")
	}

	var roles []model.AdminRole
	if err := query.Offset(pg.GetOffset()).Limit(pg.GetPageSize()).Find(&roles).Error; err != nil {
		h.Error(c, err.Error())
		return
	}

	roleIDs := make([]uint, 0, len(roles))
	for i := range roles {
		roleIDs = append(roleIDs, roles[i].ID)
	}
	permsMap, err := h.getPermsMap(roleIDs)
	if err != nil {
		h.Error(c, err.Error())
		return
	}
	for i := range roles {
		roles[i].Permissions = permsMap[roles[i].ID]
	}

	h.SuccessWithPagination(c, roles, total)
}

func (h *AdminRoleHandler) Get(c *gin.Context) {
	id, err := h.ParseID(c)
	if err != nil {
		return
	}
	if role := h.findRole(c, id); role != nil {
		h.Success(c, role)
	}
}

func (h *AdminRoleHandler) Create(c *gin.Context) {
	var req createRoleReq
	if err := h.BindJSON(c, &req); err != nil {
		return
	}

	if h.exists("code = ?", req.Code) {
		h.Error(c, "角色代码已存在")
		return
	}
	if h.exists("name = ?", req.Name) {
		h.Error(c, "角色名称已存在")
		return
	}

	role := &model.AdminRole{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Enabled:     req.Enabled == nil || *req.Enabled,
	}

	h.ExecTx(c, h.db, func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err
		}
		role.Permissions = req.Permissions
		return h.savePerms(tx, role.ID, req.Permissions)
	}, "创建成功", role)
}

func (h *AdminRoleHandler) Update(c *gin.Context) {
	id, err := h.ParseID(c)
	if err != nil {
		return
	}

	var req updateRoleReq
	if err := h.BindJSON(c, &req); err != nil {
		return
	}

	role := h.findRole(c, id)
	if role == nil {
		return
	}

	if req.Name != "" && req.Name != role.Name && h.exists("name = ? AND id != ?", req.Name, id) {
		h.Error(c, "角色名称已存在")
		return
	}

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

	if len(updates) > 0 {
		if err := h.db.Model(role).Updates(updates).Error; err != nil {
			h.Error(c, err.Error())
			return
		}
	}

	h.SuccessWithMessage(c, "更新成功", h.findRole(c, id))
}

func (h *AdminRoleHandler) Delete(c *gin.Context) {
	id, err := h.ParseID(c)
	if err != nil {
		return
	}

	h.ExecTx(c, h.db, func(tx *gorm.DB) error {
		tx.Where("role_id = ?", id).Delete(&model.AdminRolePermission{})
		tx.Where("role_id = ?", id).Delete(&model.AdminUserRole{})
		return tx.Delete(&model.AdminRole{}, id).Error
	}, "删除成功", nil)
}

func (h *AdminRoleHandler) UpdatePermissions(c *gin.Context) {
	id, err := h.ParseID(c)
	if err != nil {
		return
	}

	var req updateRolePermReq
	if err := h.BindJSON(c, &req); err != nil {
		return
	}

	if h.findRole(c, id) == nil {
		return
	}

	h.ExecTx(c, h.db, func(tx *gorm.DB) error {
		tx.Where("role_id = ?", id).Delete(&model.AdminRolePermission{})
		return h.savePerms(tx, id, req.Permissions)
	}, "权限配置成功", nil)
}

func (h *AdminRoleHandler) GetAllPermissions(c *gin.Context) {
	h.Success(c, crud.GetAllPermissions())
}

func (h *AdminRoleHandler) GetAll(c *gin.Context) {
	var roles []model.AdminRole
	h.db.Where("enabled = ?", true).Find(&roles)
	h.Success(c, roles)
}

// 私有方法
func (h *AdminRoleHandler) findRole(c *gin.Context, id uint) *model.AdminRole {
	var role model.AdminRole
	if !h.QueryOne(c, h.db.Where("id = ?", id), &role, "角色不存在") {
		return nil
	}
	role.Permissions, _ = h.getPerms(id)
	return &role
}

func (h *AdminRoleHandler) getPerms(roleID uint) ([]string, error) {
	var perms []string
	err := h.db.Model(&model.AdminRolePermission{}).Where("role_id = ?", roleID).Pluck("permission", &perms).Error
	return perms, err
}

func (h *AdminRoleHandler) getPermsMap(roleIDs []uint) (map[uint][]string, error) {
	permsMap := make(map[uint][]string)
	if len(roleIDs) == 0 {
		return permsMap, nil
	}

	type row struct {
		RoleID     uint
		Permission string
	}
	var rows []row
	err := h.db.Model(&model.AdminRolePermission{}).
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
	for _, p := range perms {
		if err := tx.Create(&model.AdminRolePermission{RoleID: roleID, Permission: p}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (h *AdminRoleHandler) exists(query string, args ...interface{}) bool {
	var count int64
	h.db.Model(&model.AdminRole{}).Where(query, args...).Count(&count)
	return count > 0
}

func init() {
	crud.RegisterModule(NewAdminRoleHandler)
}

var _ crud.Module = (*AdminRoleHandler)(nil)
