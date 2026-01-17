package handler

import (
	"bico-admin/internal/admin/model"
	"bico-admin/internal/pkg/crud"
	"errors"

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
}

func NewAdminRoleHandler(db *gorm.DB) *AdminRoleHandler {
	h := &AdminRoleHandler{}
	h.DB = db
	h.NotFoundMsg = "角色不存在"

	h.BuildListQuery = func(db *gorm.DB, req *roleListReq) *gorm.DB {
		query := db.Model(&model.AdminRole{})
		if req.Name != "" {
			query = query.Where("name LIKE ?", "%"+req.Name+"%")
		}
		if req.Code != "" {
			query = query.Where("code LIKE ?", "%"+req.Code+"%")
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
		}
		return nil
	}

	h.AfterGet = func(item *model.AdminRole) error {
		perms, err := h.getPerms(db, item.ID)
		if err != nil {
			return err
		}
		item.Permissions = perms
		return nil
	}

	h.NewModelFromCreate = func(req *createRoleReq) (*model.AdminRole, error) {
		exists, err := crud.Exists(db, &model.AdminRole{}, "code = ?", req.Code)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("角色代码已存在")
		}

		exists, err = crud.Exists(db, &model.AdminRole{}, "name = ?", req.Name)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("角色名称已存在")
		}
		return &model.AdminRole{
			Name:        req.Name,
			Code:        req.Code,
			Description: req.Description,
			Enabled:     req.Enabled == nil || *req.Enabled,
			Permissions: req.Permissions,
		}, nil
	}

	h.CreateInTx = func(tx *gorm.DB, item *model.AdminRole, req *createRoleReq) error {
		item.Permissions = req.Permissions
		return h.savePerms(tx, item.ID, req.Permissions)
	}

	h.BuildUpdates = func(req *updateRoleReq, existing *model.AdminRole) (map[string]interface{}, error) {
		// 这里校验名称唯一性：只有在传入 name 且发生变更时才检查
		if req.Name != "" && req.Name != existing.Name {
			exists, err := crud.Exists(db, &model.AdminRole{}, "name = ? AND id != ?", req.Name, existing.ID)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, errors.New("角色名称已存在")
			}
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
		return updates, nil
	}

	h.ReloadAfterUpdate = func(tx *gorm.DB, id uint, existing *model.AdminRole) error {
		perms, err := h.getPerms(tx, id)
		if err != nil {
			return err
		}
		existing.Permissions = perms
		return nil
	}

	h.DeleteInTx = func(tx *gorm.DB, id uint) error {
		tx.Where("role_id = ?", id).Delete(&model.AdminRolePermission{})
		tx.Where("role_id = ?", id).Delete(&model.AdminUserRole{})
		return nil
	}

	return h
}

func (h *AdminRoleHandler) ModuleConfig() crud.ModuleConfig {
	return crud.ModuleConfig{
		Name:             "admin_role",
		Group:            "/admin-roles",
		ParentPermission: PermSystemManage,
		Permissions:      rolePerms.Tree,
		Routes: rolePerms.RoutesWithExtra(
			crud.AuthRoute("GET", "/all", "GetAll"),
			crud.AuthRoute("GET", "/permissions", "GetAllPermissions"),
			crud.Route{Method: "GET", Path: "/:id/permissions", Handler: "GetPermissions", Permission: rolePerms.List},
			crud.Route{Method: "PUT", Path: "/:id/permissions", Handler: "UpdatePermissions", Permission: "system:admin_role:permission"},
		),
	}
}

// 请求结构
type (
	roleListReq struct {
		Name        string `form:"name"`
		Code        string `form:"code"`
		Description string `form:"description"`
		Enabled     *bool  `form:"enabled"`
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

func (h *AdminRoleHandler) GetPermissions(c *gin.Context) {
	id, err := h.ParseID(c)
	if err != nil {
		return
	}

	perms, err := h.getPerms(h.DB, id)
	if err != nil {
		h.Error(c, err.Error())
		return
	}

	h.Success(c, gin.H{"permissions": perms})
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
	// 只需要校验角色是否存在，不需要加载权限
	var role model.AdminRole
	if !h.QueryOne(c, h.DB.Where("id = ?", id), &role, "角色不存在") {
		return
	}

	h.ExecTx(c, h.DB, func(tx *gorm.DB) error {
		tx.Where("role_id = ?", id).Delete(&model.AdminRolePermission{})
		return h.savePerms(tx, id, req.Permissions)
	}, "权限配置成功", nil)
}

func (h *AdminRoleHandler) GetAllPermissions(c *gin.Context) {
	h.Success(c, crud.GetAllPermissions())
}

func (h *AdminRoleHandler) GetAll(c *gin.Context) {
	var roles []model.AdminRole
	h.DB.Where("enabled = ?", true).Find(&roles)
	h.Success(c, roles)
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

	items := make([]model.AdminRolePermission, 0, len(perms))
	for _, p := range perms {
		items = append(items, model.AdminRolePermission{RoleID: roleID, Permission: p})
	}
	return tx.CreateInBatches(items, 100).Error
}

var _ crud.Module = (*AdminRoleHandler)(nil)
