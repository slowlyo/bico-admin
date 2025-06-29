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

	// 2. 创建默认管理员用户
	if err := s.SeedDefaultAdmin(); err != nil {
		return err
	}

	// 注意：权限定义在代码中，不需要在数据库中存储权限数据
	// 角色权限关联通过 role_permissions 表管理，但权限本身在代码中定义

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

// 注意：权限定义在代码中，不需要在数据库中存储权限数据
// SeedPermissions 方法已移除，权限通过代码常量定义

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

// 注意：超级管理员权限分配已移除
// 超级管理员在代码中直接拥有所有权限，无需数据库存储
