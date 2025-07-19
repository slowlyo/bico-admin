package initializer

import (
	"fmt"

	"gorm.io/gorm"

	"bico-admin/internal/admin/models"
	pkgLogger "bico-admin/pkg/logger"
)

// DatabaseInitializer 数据库初始化器
type DatabaseInitializer struct {
	db *gorm.DB
}

// NewDatabaseInitializer 创建数据库初始化器
func NewDatabaseInitializer(db *gorm.DB) *DatabaseInitializer {
	return &DatabaseInitializer{
		db: db,
	}
}

// AutoMigrateAdminModels 自动迁移Admin模块的数据库表
func (d *DatabaseInitializer) AutoMigrateAdminModels() error {
	// 基础模型列表
	modelList := []interface{}{
		&models.AdminUser{},
		&models.AdminRole{},
		&models.AdminUserRole{},
		&models.AdminRolePermission{},
	}

	// 注意：生成的模型应该直接添加到上面的 modelList 数组中
	// 不再使用动态注册模式

	if err := d.db.AutoMigrate(modelList...); err != nil {
		return fmt.Errorf("Admin模块数据库迁移失败: %w", err)
	}

	pkgLogger.Info("Admin模块数据库迁移完成")
	return nil
}

// InitializeDefaultData 初始化默认数据
func (d *DatabaseInitializer) InitializeDefaultData() error {
	// 检查是否已存在管理员用户
	var userCount int64
	if err := d.db.Model(&models.AdminUser{}).Count(&userCount).Error; err != nil {
		return fmt.Errorf("检查管理员用户数量失败: %w", err)
	}

	// 如果已有用户，跳过初始化
	if userCount > 0 {
		pkgLogger.Info("检测到已存在管理员用户，跳过默认数据初始化")
		return nil
	}

	// 创建种子数据
	seeder := NewSeeder(d.db)
	if err := seeder.SeedAll(); err != nil {
		return fmt.Errorf("创建种子数据失败: %w", err)
	}

	pkgLogger.Info("默认数据初始化完成")
	return nil
}

// MigrateAndSeed 执行数据库迁移和种子数据创建
func (d *DatabaseInitializer) MigrateAndSeed() error {
	// 1. 执行数据库迁移
	if err := d.AutoMigrateAdminModels(); err != nil {
		return err
	}

	// 2. 初始化默认数据
	if err := d.InitializeDefaultData(); err != nil {
		return err
	}

	return nil
}
