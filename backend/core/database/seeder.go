package database

import (
	"log"

	"gorm.io/gorm"

	"bico-admin/core/model"
)

// Seeder 数据库种子数据管理器
type Seeder struct {
	db *gorm.DB
}

// NewSeeder 创建种子数据管理器
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// SeedAll 执行所有种子数据
func (s *Seeder) SeedAll() error {
	log.Println("开始执行RBAC权限系统种子数据...")

	// 1. 创建角色
	if err := s.SeedRoles(); err != nil {
		return err
	}

	// 2. 创建权限
	if err := s.SeedPermissions(); err != nil {
		return err
	}

	// 3. 创建默认管理员用户
	if err := s.SeedDefaultAdmin(); err != nil {
		return err
	}

	// 4. 分配超级管理员权限
	if err := s.AssignSuperAdminPermissions(); err != nil {
		return err
	}

	log.Println("RBAC权限系统种子数据执行完成")
	return nil
}

// SeedRoles 创建默认角色
func (s *Seeder) SeedRoles() error {
	roles := []model.Role{
		{
			Name:        "超级管理员",
			Code:        "super_admin",
			Description: "系统超级管理员，拥有所有权限，不可删除或编辑",
			Status:      model.RoleStatusActive,
		},
		{
			Name:        "管理员",
			Code:        "manager",
			Description: "系统管理员，拥有大部分管理权限",
			Status:      model.RoleStatusActive,
		},
		{
			Name:        "普通用户",
			Code:        "user",
			Description: "普通用户，拥有基础权限",
			Status:      model.RoleStatusActive,
		},
		// 保留原有的 admin 角色以兼容现有数据
		{
			Name:        "管理员(旧)",
			Code:        "admin",
			Description: "兼容旧版本的管理员角色，建议迁移到 super_admin",
			Status:      model.RoleStatusActive,
		},
	}

	for _, role := range roles {
		var existingRole model.Role
		if err := s.db.Where("code = ?", role.Code).First(&existingRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := s.db.Create(&role).Error; err != nil {
					log.Printf("创建角色失败: %v", err)
					return err
				}
				log.Printf("创建角色成功: %s", role.Name)
			} else {
				return err
			}
		} else {
			log.Printf("角色已存在: %s", role.Name)
		}
	}

	return nil
}

// SeedPermissions 创建默认权限
func (s *Seeder) SeedPermissions() error {
	// 使用权限常量定义权限数据
	permissions := []model.Permission{
		// 系统管理权限
		{
			Name:        "查看系统",
			Code:        "system:view",
			Type:        model.PermissionTypeMenu,
			Resource:    "system",
			Action:      "view",
			Description: "查看系统管理页面",
			Status:      model.PermissionStatusActive,
		},
		{
			Name:        "管理系统",
			Code:        "system:manage",
			Type:        model.PermissionTypeAPI,
			Resource:    "system",
			Action:      "manage",
			Description: "管理系统设置",
			Status:      model.PermissionStatusActive,
		},

		// 用户管理权限
		{
			Name:        "查看用户",
			Code:        "user:view",
			Type:        model.PermissionTypeAPI,
			Resource:    "user",
			Action:      "view",
			Description: "查看用户列表和详情",
			Status:      model.PermissionStatusActive,
		},
		{
			Name:        "创建用户",
			Code:        "user:create",
			Type:        model.PermissionTypeAPI,
			Resource:    "user",
			Action:      "create",
			Description: "创建新用户",
			Status:      model.PermissionStatusActive,
		},
		{
			Name:        "编辑用户",
			Code:        "user:update",
			Type:        model.PermissionTypeAPI,
			Resource:    "user",
			Action:      "update",
			Description: "编辑用户信息",
			Status:      model.PermissionStatusActive,
		},
		{
			Name:        "删除用户",
			Code:        "user:delete",
			Type:        model.PermissionTypeAPI,
			Resource:    "user",
			Action:      "delete",
			Description: "删除用户",
			Status:      model.PermissionStatusActive,
		},
		{
			Name:        "管理用户状态",
			Code:        "user:manage_status",
			Type:        model.PermissionTypeAPI,
			Resource:    "user",
			Action:      "manage_status",
			Description: "启用/禁用用户",
			Status:      model.PermissionStatusActive,
		},
		{
			Name:        "重置密码",
			Code:        "user:reset_password",
			Type:        model.PermissionTypeAPI,
			Resource:    "user",
			Action:      "reset_password",
			Description: "重置用户密码",
			Status:      model.PermissionStatusActive,
		},

		// 角色管理权限
		{
			Name:        "查看角色",
			Code:        "role:view",
			Type:        model.PermissionTypeAPI,
			Resource:    "role",
			Action:      "view",
			Description: "查看角色列表和详情",
			Status:      model.PermissionStatusActive,
		},
		{
			Name:        "创建角色",
			Code:        "role:create",
			Type:        model.PermissionTypeAPI,
			Resource:    "role",
			Action:      "create",
			Description: "创建新角色",
			Status:      model.PermissionStatusActive,
		},
		{
			Name:        "编辑角色",
			Code:        "role:update",
			Type:        model.PermissionTypeAPI,
			Resource:    "role",
			Action:      "update",
			Description: "编辑角色信息",
			Status:      model.PermissionStatusActive,
		},
		{
			Name:        "删除角色",
			Code:        "role:delete",
			Type:        model.PermissionTypeAPI,
			Resource:    "role",
			Action:      "delete",
			Description: "删除角色",
			Status:      model.PermissionStatusActive,
		},
		{
			Name:        "分配权限",
			Code:        "role:assign_permissions",
			Type:        model.PermissionTypeAPI,
			Resource:    "role",
			Action:      "assign_permissions",
			Description: "为角色分配权限",
			Status:      model.PermissionStatusActive,
		},
		{
			Name:        "编辑用户",
			Code:        "user.update",
			Type:        model.PermissionTypeAPI,
			Resource:    "user",
			Action:      "update",
			Description: "编辑用户信息",
			Status:      model.PermissionStatusActive,
		},
		{
			Name:        "删除用户",
			Code:        "user.delete",
			Type:        model.PermissionTypeAPI,
			Resource:    "user",
			Action:      "delete",
			Description: "删除用户",
			Status:      model.PermissionStatusActive,
		},
		// 角色管理权限
		{
			Name:        "角色管理",
			Code:        "role.manage",
			Type:        model.PermissionTypeMenu,
			Resource:    "role",
			Action:      "manage",
			Description: "角色管理菜单权限",
			Status:      model.PermissionStatusActive,
		},
		{
			Name:        "查看角色",
			Code:        "role.view",
			Type:        model.PermissionTypeAPI,
			Resource:    "role",
			Action:      "view",
			Description: "查看角色列表和详情",
			Status:      model.PermissionStatusActive,
		},
		// 系统管理权限
		{
			Name:        "系统管理",
			Code:        "system.manage",
			Type:        model.PermissionTypeMenu,
			Resource:    "system",
			Action:      "manage",
			Description: "系统管理菜单权限",
			Status:      model.PermissionStatusActive,
		},
	}

	for _, permission := range permissions {
		var existingPermission model.Permission
		if err := s.db.Where("code = ?", permission.Code).First(&existingPermission).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := s.db.Create(&permission).Error; err != nil {
					log.Printf("创建权限失败: %v", err)
					return err
				}
				log.Printf("创建权限成功: %s", permission.Name)
			} else {
				return err
			}
		} else {
			log.Printf("权限已存在: %s", permission.Name)
		}
	}

	return nil
}

// SeedDefaultAdmin 创建默认管理员用户
func (s *Seeder) SeedDefaultAdmin() error {
	// 检查是否已存在管理员用户
	var adminUser model.User
	if err := s.db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建默认管理员用户
			admin := model.User{
				Username: "admin",
				Email:    "admin@bico.com",
				Password: "123456", // 默认密码，建议首次登录后修改
				Nickname: "系统管理员",
				Role:     "super_admin", // 使用超级管理员角色
				Status:   model.UserStatusActive,
			}

			// 加密密码
			if err := admin.HashPassword(); err != nil {
				log.Printf("密码加密失败: %v", err)
				return err
			}

			// 创建用户
			if err := s.db.Create(&admin).Error; err != nil {
				log.Printf("创建管理员用户失败: %v", err)
				return err
			}

			// 获取超级管理员角色
			var adminRole model.Role
			if err := s.db.Where("code = ?", "super_admin").First(&adminRole).Error; err != nil {
				log.Printf("获取超级管理员角色失败: %v", err)
				return err
			}

			// 关联角色
			if err := s.db.Model(&admin).Association("Roles").Append(&adminRole); err != nil {
				log.Printf("关联角色失败: %v", err)
				return err
			}

			log.Println("创建默认管理员用户成功: admin/123456")
		} else {
			return err
		}
	} else {
		log.Println("管理员用户已存在")
	}

	return nil
}

// SeedIfEmpty 仅在数据库为空时执行种子数据
func (s *Seeder) SeedIfEmpty() error {
	var userCount int64
	if err := s.db.Model(&model.User{}).Count(&userCount).Error; err != nil {
		return err
	}

	if userCount == 0 {
		log.Println("数据库为空，执行种子数据...")
		return s.SeedAll()
	}

	log.Println("数据库已有数据，跳过种子数据")
	return nil
}

// AssignSuperAdminPermissions 为超级管理员角色分配所有权限
func (s *Seeder) AssignSuperAdminPermissions() error {
	log.Println("为超级管理员角色分配权限...")

	// 获取超级管理员角色
	var superAdminRole model.Role
	if err := s.db.Where("code = ?", "super_admin").First(&superAdminRole).Error; err != nil {
		log.Printf("未找到超级管理员角色: %v", err)
		return err
	}

	// 获取所有权限
	var allPermissions []model.Permission
	if err := s.db.Where("status = ?", model.PermissionStatusActive).Find(&allPermissions).Error; err != nil {
		log.Printf("获取权限列表失败: %v", err)
		return err
	}

	// 检查是否已经分配了权限
	var existingCount int64
	s.db.Model(&model.RolePermission{}).Where("role_id = ?", superAdminRole.ID).Count(&existingCount)

	if existingCount > 0 {
		log.Printf("超级管理员角色已分配 %d 个权限，跳过", existingCount)
		return nil
	}

	// 为超级管理员角色分配所有权限
	successCount := 0
	for _, perm := range allPermissions {
		rolePermission := model.RolePermission{
			RoleID:       superAdminRole.ID,
			PermissionID: perm.ID,
		}

		if err := s.db.Create(&rolePermission).Error; err != nil {
			log.Printf("分配权限 '%s' 失败: %v", perm.Code, err)
		} else {
			successCount++
		}
	}

	log.Printf("为超级管理员角色成功分配了 %d/%d 个权限", successCount, len(allPermissions))
	return nil
}
