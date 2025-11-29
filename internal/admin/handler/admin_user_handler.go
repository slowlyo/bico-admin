package handler

import (
	"bico-admin/internal/admin/model"
	"bico-admin/internal/pkg/crud"
	"bico-admin/internal/pkg/password"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 权限定义
var userPerms = crud.NewCRUDPerms("admin_user", "用户管理")

// AdminUserHandler 用户管理处理器
type AdminUserHandler struct {
	crud.BaseHandler
	db *gorm.DB
}

func NewAdminUserHandler(db *gorm.DB) *AdminUserHandler {
	return &AdminUserHandler{db: db}
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
		RoleIDs  []uint `json:"roleIds"`
	}
	updateUserReq struct {
		Name    string `json:"name"`
		Avatar  string `json:"avatar"`
		Enabled *bool  `json:"enabled"`
		RoleIDs []uint `json:"roleIds"`
	}
)

func (h *AdminUserHandler) List(c *gin.Context) {
	var req userListReq
	h.BindQuery(c, &req)

	query := h.db.Model(&model.AdminUser{}).Preload("Roles")
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Enabled != nil {
		query = query.Where("enabled = ?", *req.Enabled)
	}

	var users []model.AdminUser
	h.QueryList(c, query, &users)
}

func (h *AdminUserHandler) Get(c *gin.Context) {
	id, err := h.ParseID(c)
	if err != nil {
		return
	}

	var user model.AdminUser
	if h.QueryOne(c, h.db.Preload("Roles").Where("id = ?", id), &user, "用户不存在") {
		h.Success(c, user)
	}
}

func (h *AdminUserHandler) Create(c *gin.Context) {
	var req createUserReq
	if err := h.BindJSON(c, &req); err != nil {
		return
	}

	if h.exists("username = ?", req.Username) {
		h.Error(c, "用户名已存在")
		return
	}

	hashed, err := password.Hash(req.Password)
	if err != nil {
		h.Error(c, err.Error())
		return
	}

	user := &model.AdminUser{
		Username: req.Username,
		Password: hashed,
		Name:     req.Name,
		Avatar:   req.Avatar,
		Enabled:  req.Enabled == nil || *req.Enabled,
	}

	h.ExecTx(c, h.db, func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		if err := h.syncRoles(tx, user, req.RoleIDs); err != nil {
			return err
		}
		return tx.Preload("Roles").First(user, user.ID).Error
	}, "创建成功", user)
}

func (h *AdminUserHandler) Update(c *gin.Context) {
	id, err := h.ParseID(c)
	if err != nil {
		return
	}

	var req updateUserReq
	if err := h.BindJSON(c, &req); err != nil {
		return
	}

	var user model.AdminUser
	if !h.QueryOne(c, h.db.Where("id = ?", id), &user, "用户不存在") {
		return
	}

	h.ExecTx(c, h.db, func(tx *gorm.DB) error {
		updates := map[string]interface{}{}
		if req.Name != "" {
			updates["name"] = req.Name
		}
		if req.Avatar != "" {
			updates["avatar"] = req.Avatar
		}
		if req.Enabled != nil {
			updates["enabled"] = *req.Enabled
		}
		if len(updates) > 0 {
			if err := tx.Model(&user).Updates(updates).Error; err != nil {
				return err
			}
		}
		if req.RoleIDs != nil {
			if err := h.syncRoles(tx, &user, req.RoleIDs); err != nil {
				return err
			}
		}
		return tx.Preload("Roles").First(&user, id).Error
	}, "更新成功", &user)
}

func (h *AdminUserHandler) Delete(c *gin.Context) {
	id, err := h.ParseID(c)
	if err != nil {
		return
	}

	h.ExecTx(c, h.db, func(tx *gorm.DB) error {
		var user model.AdminUser
		if err := tx.First(&user, id).Error; err != nil {
			return err
		}
		_ = tx.Model(&user).Association("Roles").Clear()
		return tx.Delete(&user).Error
	}, "删除成功", nil)
}

// 私有方法
func (h *AdminUserHandler) exists(query string, args ...interface{}) bool {
	var count int64
	h.db.Model(&model.AdminUser{}).Where(query, args...).Count(&count)
	return count > 0
}

func (h *AdminUserHandler) syncRoles(tx *gorm.DB, user *model.AdminUser, roleIDs []uint) error {
	_ = tx.Model(user).Association("Roles").Clear()
	if len(roleIDs) == 0 {
		return nil
	}
	var roles []*model.AdminRole
	if err := tx.Where("id IN ?", crud.UniqueUints(roleIDs)).Find(&roles).Error; err != nil {
		return err
	}
	return tx.Model(user).Association("Roles").Append(roles)
}

func init() {
	crud.RegisterModule(NewAdminUserHandler)
}

var _ crud.Module = (*AdminUserHandler)(nil)
