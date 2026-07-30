package migrate

import (
	"fmt"

	adminModel "bico-admin/internal/admin/model"
	"bico-admin/internal/core/logger"
	"bico-admin/internal/pkg/password"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据表
func AutoMigrate(db *gorm.DB, mode string) error {
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

	// 初始化超级管理员角色与首个管理员账户。
	if err := initSuperAdmin(db, mode); err != nil {
		return err
	}

	// API 模块模型（暂无）

	return nil
}

// initSuperAdmin 初始化超级管理员角色和账号关联。
// 既有数据库会将原 admin 账号迁移到保留角色，保持升级前权限。
func initSuperAdmin(db *gorm.DB, mode string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		role := adminModel.AdminRole{
			Name:        "超级管理员",
			Code:        adminModel.SuperAdminRoleCode,
			Description: "系统保留角色，自动拥有全部权限",
			Enabled:     true,
		}
		if err := tx.Where("code = ?", adminModel.SuperAdminRoleCode).FirstOrCreate(&role).Error; err != nil {
			return fmt.Errorf("初始化超级管理员角色失败: %w", err)
		}
		if !role.Enabled {
			// 系统保留角色必须保持启用，否则所有超级管理员会同时失权。
			if err := tx.Model(&role).Update("enabled", true).Error; err != nil {
				return fmt.Errorf("启用超级管理员角色失败: %w", err)
			}
		}

		admin, err := ensureInitialAdmin(tx, mode)
		if err != nil {
			return err
		}
		if admin == nil {
			// 已有用户但不存在历史 admin 账号时，不擅自提升任意用户。
			return nil
		}

		relation := adminModel.AdminUserRole{UserID: admin.ID, RoleID: role.ID}
		if err := tx.Where("user_id = ? AND role_id = ?", admin.ID, role.ID).FirstOrCreate(&relation).Error; err != nil {
			return fmt.Errorf("授予超级管理员角色失败: %w", err)
		}
		return nil
	})
}

// ensureInitialAdmin 确保空数据库拥有首个管理员账号。
// 初始账号沿用系统默认凭据，避免迁移依赖部署环境变量。
func ensureInitialAdmin(db *gorm.DB, mode string) (*adminModel.AdminUser, error) {
	var count int64
	if err := db.Model(&adminModel.AdminUser{}).Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		var admin adminModel.AdminUser
		if err := db.Where("username = ?", "admin").First(&admin).Error; err != nil {
			// 历史 admin 不存在时保持现有账号不变，避免意外提权。
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, err
		}
		return &admin, nil
	}

	hashedPassword, err := password.Hash("admin")
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	admin := adminModel.AdminUser{
		Username: "admin",
		Password: hashedPassword,
		Name:     "系统管理员",
		Avatar:   "https://api.dicebear.com/9.x/thumbs/png?seed=slowlyo",
		Enabled:  true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return nil, fmt.Errorf("创建管理员失败: %w", err)
	}

	logger.Info("初始化管理员账户成功", zap.String("username", "admin"), zap.String("password", "admin"))
	return &admin, nil
}
