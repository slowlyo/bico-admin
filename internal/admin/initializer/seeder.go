package initializer

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"bico-admin/internal/admin/models"
	pkgLogger "bico-admin/pkg/logger"
)

// Seeder 种子数据管理器
type Seeder struct {
	db *gorm.DB
}

// NewSeeder 创建种子数据管理器
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{
		db: db,
	}
}

// SeedAll 创建所有种子数据
func (s *Seeder) SeedAll() error {
	// 1. 创建超级管理员角色
	superAdminRole, err := s.createSuperAdminRole()
	if err != nil {
		return fmt.Errorf("创建超级管理员角色失败: %w", err)
	}

	// 2. 创建超级管理员用户
	superAdminUser, err := s.createSuperAdminUser()
	if err != nil {
		return fmt.Errorf("创建超级管理员用户失败: %w", err)
	}

	// 3. 为超级管理员分配角色
	if err := s.assignRoleToUser(superAdminUser.ID, superAdminRole.ID); err != nil {
		return fmt.Errorf("为超级管理员分配角色失败: %w", err)
	}

	pkgLogger.Info("种子数据创建完成 - 超级管理员: admin/admin123")
	return nil
}

// createSuperAdminRole 创建超级管理员角色
func (s *Seeder) createSuperAdminRole() (*models.AdminRole, error) {
	// 检查是否已存在超级管理员角色
	var existingRole models.AdminRole
	err := s.db.Where("code = ?", models.RoleCodeSuperAdmin).First(&existingRole).Error
	if err == nil {
		return &existingRole, nil // 已存在，直接返回
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 创建超级管理员角色
	superAdminRole := &models.AdminRole{
		Name:        "超级管理员",
		Code:        models.RoleCodeSuperAdmin,
		Description: "系统超级管理员，拥有所有权限且不可修改删除",
		Status:      1,
	}

	if err := s.db.Create(superAdminRole).Error; err != nil {
		return nil, err
	}

	// 注意：超级管理员不需要分配具体权限，在权限检查时会特殊处理
	pkgLogger.Info("已创建超级管理员角色")
	return superAdminRole, nil
}

// createSuperAdminUser 创建超级管理员用户
func (s *Seeder) createSuperAdminUser() (*models.AdminUser, error) {
	// 检查是否已存在超级管理员用户
	var existingUser models.AdminUser
	err := s.db.Where("username = ?", "admin").First(&existingUser).Error
	if err == nil {
		return &existingUser, nil // 已存在，直接返回
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 生成密码哈希值 (密码: admin123)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("生成密码哈希失败: %w", err)
	}

	// 创建超级管理员用户
	superAdminUser := &models.AdminUser{
		Username: "admin",
		Password: string(hashedPassword),
		Name:     "超级管理员",
		Status:   1, // 启用状态
		Remark:   "系统默认超级管理员，不可删除",
	}

	if err := s.db.Create(superAdminUser).Error; err != nil {
		return nil, err
	}

	pkgLogger.Info("已创建超级管理员用户")
	return superAdminUser, nil
}

// assignRoleToUser 为用户分配角色
func (s *Seeder) assignRoleToUser(userID, roleID uint) error {
	// 检查是否已存在关联
	var existingUserRole models.AdminUserRole
	err := s.db.Where("user_id = ? AND role_id = ?", userID, roleID).First(&existingUserRole).Error
	if err == nil {
		return nil // 已存在关联，直接返回
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// 创建用户角色关联
	userRole := &models.AdminUserRole{
		UserID: userID,
		RoleID: roleID,
	}

	if err := s.db.Create(userRole).Error; err != nil {
		return err
	}

	pkgLogger.Info("已为超级管理员分配角色")
	return nil
}
