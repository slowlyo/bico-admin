package handler

import (
	"bico-admin/internal/admin/model"
	"bico-admin/internal/pkg/crud"
	"bico-admin/internal/pkg/password"
	"errors"

	"gorm.io/gorm"
)

// 权限定义
var userPerms = crud.NewCRUDPerms("system", "admin_user", "用户管理")

// AdminUserHandler 用户管理处理器
type AdminUserHandler struct {
	crud.CRUDHandler[model.AdminUser, userListReq, createUserReq, updateUserReq]
}

func NewAdminUserHandler(db *gorm.DB) *AdminUserHandler {
	h := &AdminUserHandler{}
	h.DB = db
	h.NotFoundMsg = "用户不存在"

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
		return query
	}

	h.BuildGetQuery = func(db *gorm.DB) *gorm.DB {
		return db.Model(&model.AdminUser{}).Preload("Roles")
	}

	h.NewModelFromCreate = func(req *createUserReq) (*model.AdminUser, error) {
		exists, err := crud.Exists(db, &model.AdminUser{}, "username = ?", req.Username)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("用户名已存在")
		}

		hashed, err := password.Hash(req.Password)
		if err != nil {
			return nil, err
		}

		return &model.AdminUser{
			Username: req.Username,
			Password: hashed,
			Name:     req.Name,
			Avatar:   req.Avatar,
			Enabled:  req.Enabled == nil || *req.Enabled,
		}, nil
	}

	h.CreateInTx = func(tx *gorm.DB, item *model.AdminUser, req *createUserReq) error {
		return h.syncRoles(tx, item, req.RoleIDs)
	}
	// 需要返回带 Roles 的用户数据
	h.ReloadAfterCreate = func(tx *gorm.DB, id uint, item *model.AdminUser) error {
		return tx.Preload("Roles").First(item, item.ID).Error
	}

	h.BuildUpdateQuery = func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&model.AdminUser{})
	}
	h.BuildUpdates = func(req *updateUserReq, existing *model.AdminUser) (map[string]interface{}, error) {
		updates := map[string]interface{}{}
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
		}
		if req.Enabled != nil {
			updates["enabled"] = *req.Enabled
		}
		return updates, nil
	}

	h.UpdateInTx = func(tx *gorm.DB, id uint, existing *model.AdminUser, req *updateUserReq) error {
		if req.RoleIDs == nil {
			return nil
		}
		return h.syncRoles(tx, existing, req.RoleIDs)
	}
	h.ReloadAfterUpdate = func(tx *gorm.DB, id uint, existing *model.AdminUser) error {
		return tx.Preload("Roles").First(existing, id).Error
	}

	h.DeleteInTx = func(tx *gorm.DB, id uint) error {
		var user model.AdminUser
		if err := tx.First(&user, id).Error; err != nil {
			return err
		}
		// 删除用户前先清空角色关联，任一步失败都回滚。
		if err := tx.Model(&user).Association("Roles").Clear(); err != nil {
			return err
		}
		return nil
	}

	return h
}

func (h *AdminUserHandler) ModuleConfig() crud.ModuleConfig {
	return crud.ModuleConfig{
		Name:             "admin_user",
		Group:            "/admin-users",
		ParentPermission: PermSystemManage,
		Permissions:      userPerms.Tree,
		Routes:           userPerms.Routes(),
	}
}

// 请求结构
type (
	userListReq struct {
		Username string `form:"username"`
		Name     string `form:"name"`
		Enabled  *bool  `form:"enabled"`
	}
	createUserReq struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name"`
		Avatar   string `json:"avatar"`
		Enabled  *bool  `json:"enabled"`
		RoleIDs  []uint `json:"role_ids"`
	}
	updateUserReq struct {
		Name     string `json:"name"`
		Avatar   string `json:"avatar"`
		Enabled  *bool  `json:"enabled"`
		RoleIDs  []uint `json:"role_ids"`
		Password string `json:"password"`
	}
)

func (h *AdminUserHandler) syncRoles(tx *gorm.DB, user *model.AdminUser, roleIDs []uint) error {
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
		return errors.New("存在无效角色 ID")
	}
	if err := tx.Model(user).Association("Roles").Append(roles); err != nil {
		return err
	}
	return nil
}

var _ crud.Module = (*AdminUserHandler)(nil)
