package migrate

import (
	adminModel "bico-admin/internal/admin/model"
	sharedModel "bico-admin/internal/shared/model"
	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据表
func AutoMigrate(db *gorm.DB) error {
	// 共享模型
	if err := db.AutoMigrate(
		&sharedModel.User{},
		&sharedModel.Role{},
		&sharedModel.Log{},
	); err != nil {
		return err
	}

	// Admin 模块模型
	if err := db.AutoMigrate(
		&adminModel.Menu{},
	); err != nil {
		return err
	}

	// API 模块模型（暂无）

	return nil
}
