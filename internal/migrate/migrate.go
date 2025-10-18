package migrate

import (
	"fmt"

	adminModel "bico-admin/internal/admin/model"
	"bico-admin/internal/pkg/password"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据表
func AutoMigrate(db *gorm.DB) error {
	// Admin 模块模型
	if err := db.AutoMigrate(
		&adminModel.Menu{},
		&adminModel.AdminUser{},
		&adminModel.AdminRole{},
		&adminModel.AdminRolePermission{},
		&adminModel.AdminUserRole{},
	); err != nil {
		return err
	}

	// 初始化管理员账户
	if err := initAdminUser(db); err != nil {
		return err
	}

	// API 模块模型（暂无）

	return nil
}

// initAdminUser 初始化管理员账户
func initAdminUser(db *gorm.DB) error {
	var count int64
	if err := db.Model(&adminModel.AdminUser{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		fmt.Println("⏭️  管理员账户已存在，跳过初始化")
		return nil
	}

	hashedPassword, err := password.Hash("admin")
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	admin := adminModel.AdminUser{
		Username: "admin",
		Password: hashedPassword,
		Name:     "系统管理员",
		Avatar:   "https://api.dicebear.com/9.x/thumbs/png?seed=slowlyo",
		Enabled:  true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return fmt.Errorf("创建管理员失败: %w", err)
	}

	fmt.Printf("✅ 初始化管理员账户成功 (用户名: admin, 密码: admin)\n")
	fmt.Println("💡 admin 账户自动拥有所有权限，后续新增权限无需手动分配")
	return nil
}
