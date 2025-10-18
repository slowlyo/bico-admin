package service

import (
	"errors"

	"bico-admin/internal/admin/model"
	"bico-admin/internal/shared/pagination"
	"bico-admin/internal/shared/password"
	"gorm.io/gorm"
)

var (
	ErrUsernameExists = errors.New("用户名已存在")
)

// AdminUserService 用户管理服务
type AdminUserService struct {
	db *gorm.DB
}

// NewAdminUserService 创建用户管理服务
func NewAdminUserService(db *gorm.DB) *AdminUserService {
	return &AdminUserService{db: db}
}

// ListRequest 列表查询请求
type ListRequest struct {
	pagination.Pagination
	Username string `form:"username" json:"username"`
	Name     string `form:"name" json:"name"`
	Enabled  *bool  `form:"enabled" json:"enabled"`
}

// CreateAdminUserRequest 创建用户请求
type CreateAdminUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Enabled  *bool  `json:"enabled"`
	RoleIDs  []uint `json:"roleIds"`
}

// UpdateAdminUserRequest 更新用户请求
type UpdateAdminUserRequest struct {
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	Enabled *bool  `json:"enabled"`
	RoleIDs []uint `json:"roleIds"`
}

// List 获取用户列表
func (s *AdminUserService) List(req *ListRequest) (*pagination.Response, error) {
	query := s.db.Model(&model.AdminUser{}).Preload("Roles")
	
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Enabled != nil {
		query = query.Where("enabled = ?", *req.Enabled)
	}
	
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	
	var users []model.AdminUser
	if err := query.Offset(req.GetOffset()).Limit(req.GetPageSize()).Find(&users).Error; err != nil {
		return nil, err
	}
	
	return &pagination.Response{
		Total: total,
		Data:  users,
	}, nil
}

// Get 获取用户详情
func (s *AdminUserService) Get(id uint) (*model.AdminUser, error) {
	var user model.AdminUser
	if err := s.db.Preload("Roles").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (s *AdminUserService) Create(req *CreateAdminUserRequest) (*model.AdminUser, error) {
	var existingUser model.AdminUser
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, ErrUsernameExists
	}
	
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	
	user := &model.AdminUser{
		Username: req.Username,
		Password: hashedPassword,
		Name:     req.Name,
		Avatar:   req.Avatar,
		Enabled:  enabled,
	}
	
	return user, s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		
		if len(req.RoleIDs) > 0 {
			// 去重 RoleIDs
			uniqueRoleIDs := make([]uint, 0)
			seen := make(map[uint]bool)
			for _, roleID := range req.RoleIDs {
				if !seen[roleID] {
					seen[roleID] = true
					uniqueRoleIDs = append(uniqueRoleIDs, roleID)
				}
			}
			
			var roles []*model.AdminRole
			if err := tx.Where("id IN ?", uniqueRoleIDs).Find(&roles).Error; err != nil {
				return err
			}
			if err := tx.Model(user).Association("Roles").Append(roles); err != nil {
				return err
			}
		}
		
		return tx.Preload("Roles").First(user, user.ID).Error
	})
}

// Update 更新用户
func (s *AdminUserService) Update(id uint, req *UpdateAdminUserRequest) (*model.AdminUser, error) {
	var user model.AdminUser
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	
	return &user, s.db.Transaction(func(tx *gorm.DB) error {
		updates := make(map[string]interface{})
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
			// 先清除旧的角色关联
			if err := tx.Model(&user).Association("Roles").Clear(); err != nil {
				return err
			}
			
			// 去重 RoleIDs
			uniqueRoleIDs := make([]uint, 0)
			seen := make(map[uint]bool)
			for _, roleID := range req.RoleIDs {
				if !seen[roleID] {
					seen[roleID] = true
					uniqueRoleIDs = append(uniqueRoleIDs, roleID)
				}
			}
			
			// 添加新的角色关联
			if len(uniqueRoleIDs) > 0 {
				var roles []*model.AdminRole
				if err := tx.Where("id IN ?", uniqueRoleIDs).Find(&roles).Error; err != nil {
					return err
				}
				if err := tx.Model(&user).Association("Roles").Append(roles); err != nil {
					return err
				}
			}
		}
		
		return tx.Preload("Roles").First(&user, id).Error
	})
}

// Delete 删除用户
func (s *AdminUserService) Delete(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var user model.AdminUser
		if err := tx.First(&user, id).Error; err != nil {
			return err
		}
		
		if err := tx.Model(&user).Association("Roles").Clear(); err != nil {
			return err
		}
		
		return tx.Delete(&user).Error
	})
}
