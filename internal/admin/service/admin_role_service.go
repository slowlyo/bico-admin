package service

import (
	"errors"

	"bico-admin/internal/admin/model"
	"bico-admin/internal/pkg/pagination"
	"gorm.io/gorm"
)

var (
	ErrRoleCodeExists = errors.New("角色代码已存在")
	ErrRoleNameExists = errors.New("角色名称已存在")
)

// AdminRoleService 角色管理服务
type AdminRoleService struct {
	db *gorm.DB
}

// NewAdminRoleService 创建角色管理服务
func NewAdminRoleService(db *gorm.DB) *AdminRoleService {
	return &AdminRoleService{db: db}
}

// RoleListRequest 角色列表查询请求
type RoleListRequest struct {
	pagination.Pagination
	Name    string `form:"name" json:"name"`
	Code    string `form:"code" json:"code"`
	Enabled *bool  `form:"enabled" json:"enabled"`
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Code        string   `json:"code" binding:"required"`
	Description string   `json:"description"`
	Enabled     *bool    `json:"enabled"`
	Permissions []string `json:"permissions"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     *bool  `json:"enabled"`
}

// UpdateRolePermissionsRequest 更新角色权限请求
type UpdateRolePermissionsRequest struct {
	Permissions []string `json:"permissions" binding:"required"`
}

// ListRoles 获取角色列表
func (s *AdminRoleService) ListRoles(req *RoleListRequest) (*pagination.Response, error) {
	query := s.db.Model(&model.AdminRole{})
	
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		query = query.Where("code LIKE ?", "%"+req.Code+"%")
	}
	if req.Enabled != nil {
		query = query.Where("enabled = ?", *req.Enabled)
	}
	
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	
	// 应用排序
	if orderBy := req.GetOrderBy(); orderBy != "" {
		query = query.Order(orderBy)
	} else {
		query = query.Order("created_at DESC")
	}
	
	var roles []model.AdminRole
	if err := query.Offset(req.GetOffset()).Limit(req.GetPageSize()).Find(&roles).Error; err != nil {
		return nil, err
	}
	
	for i := range roles {
		permissions, _ := s.GetRolePermissions(roles[i].ID)
		roles[i].Permissions = permissions
	}
	
	return &pagination.Response{
		Total: total,
		Data:  roles,
	}, nil
}

// GetRole 获取角色详情
func (s *AdminRoleService) GetRole(id uint) (*model.AdminRole, error) {
	var role model.AdminRole
	if err := s.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	
	permissions, err := s.GetRolePermissions(id)
	if err != nil {
		return nil, err
	}
	role.Permissions = permissions
	
	return &role, nil
}

// CreateRole 创建角色
func (s *AdminRoleService) CreateRole(req *CreateRoleRequest) (*model.AdminRole, error) {
	var existingRole model.AdminRole
	if err := s.db.Where("code = ?", req.Code).First(&existingRole).Error; err == nil {
		return nil, ErrRoleCodeExists
	}
	if err := s.db.Where("name = ?", req.Name).First(&existingRole).Error; err == nil {
		return nil, ErrRoleNameExists
	}
	
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	
	role := &model.AdminRole{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Enabled:     enabled,
	}
	
	return role, s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err
		}
		
		if len(req.Permissions) > 0 {
			for _, perm := range req.Permissions {
				rolePermission := &model.AdminRolePermission{
					RoleID:     role.ID,
					Permission: perm,
				}
				if err := tx.Create(rolePermission).Error; err != nil {
					return err
				}
			}
		}
		
		role.Permissions = req.Permissions
		return nil
	})
}

// UpdateRole 更新角色
func (s *AdminRoleService) UpdateRole(id uint, req *UpdateRoleRequest) (*model.AdminRole, error) {
	var role model.AdminRole
	if err := s.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	
	if req.Name != "" && req.Name != role.Name {
		var existingRole model.AdminRole
		if err := s.db.Where("name = ? AND id != ?", req.Name, id).First(&existingRole).Error; err == nil {
			return nil, ErrRoleNameExists
		}
	}
	
	updates := make(map[string]interface{})
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
		if err := s.db.Model(&role).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	
	permissions, err := s.GetRolePermissions(id)
	if err != nil {
		return nil, err
	}
	role.Permissions = permissions
	
	return &role, nil
}

// DeleteRole 删除角色
func (s *AdminRoleService) DeleteRole(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&model.AdminRolePermission{}).Error; err != nil {
			return err
		}
		
		if err := tx.Where("role_id = ?", id).Delete(&model.AdminUserRole{}).Error; err != nil {
			return err
		}
		
		return tx.Delete(&model.AdminRole{}, id).Error
	})
}

// UpdateRolePermissions 更新角色权限
func (s *AdminRoleService) UpdateRolePermissions(id uint, req *UpdateRolePermissionsRequest) error {
	var role model.AdminRole
	if err := s.db.First(&role, id).Error; err != nil {
		return err
	}
	
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&model.AdminRolePermission{}).Error; err != nil {
			return err
		}
		
		for _, perm := range req.Permissions {
			rolePermission := &model.AdminRolePermission{
				RoleID:     id,
				Permission: perm,
			}
			if err := tx.Create(rolePermission).Error; err != nil {
				return err
			}
		}
		
		return nil
	})
}

// GetRolePermissions 获取角色权限
func (s *AdminRoleService) GetRolePermissions(roleID uint) ([]string, error) {
	var permissions []string
	err := s.db.Model(&model.AdminRolePermission{}).
		Where("role_id = ?", roleID).
		Pluck("permission", &permissions).Error
	return permissions, err
}

// GetAllRoles 获取所有角色（用于下拉选择）
func (s *AdminRoleService) GetAllRoles() ([]model.AdminRole, error) {
	var roles []model.AdminRole
	err := s.db.Where("enabled = ?", true).Find(&roles).Error
	return roles, err
}
