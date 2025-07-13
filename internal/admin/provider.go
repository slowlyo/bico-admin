package admin

import (
	"fmt"

	"github.com/google/wire"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"bico-admin/internal/admin/handler"
	"bico-admin/internal/admin/models"
	"bico-admin/internal/admin/repository"
	"bico-admin/internal/admin/routes"
	"bico-admin/internal/admin/service"
	pkgLogger "bico-admin/pkg/logger"
)

// ProviderSet Admin端Provider集合
var ProviderSet = wire.NewSet(
	// Repository层
	repository.NewAdminUserRepository,

	// Service层
	service.NewUserService,
	service.NewAdminUserService,
	service.NewAuthService,

	// Handler层
	handler.NewAuthHandler,
	handler.NewUserHandler,
	handler.NewAdminUserHandler,
	handler.NewSystemHandler,

	// 路由处理器集合
	ProvideHandlers,
)

// ProvideHandlers 提供处理器集合
func ProvideHandlers(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	adminUserHandler *handler.AdminUserHandler,
	systemHandler *handler.SystemHandler,
) *routes.Handlers {
	return &routes.Handlers{
		AuthHandler:      authHandler,
		UserHandler:      userHandler,
		AdminUserHandler: adminUserHandler,
		SystemHandler:    systemHandler,
	}
}

// AutoMigrateAdminModels 自动迁移Admin模块的数据库表
func AutoMigrateAdminModels(db *gorm.DB) error {
	modelList := []interface{}{
		&models.AdminUser{},
		&models.AdminRole{},
		&models.AdminUserRole{},
		&models.AdminRolePermission{},
	}

	if err := db.AutoMigrate(modelList...); err != nil {
		return fmt.Errorf("Admin模块数据库迁移失败: %w", err)
	}

	// 插入默认数据
	if err := insertDefaultAdminData(db); err != nil {
		return fmt.Errorf("插入默认Admin数据失败: %w", err)
	}

	pkgLogger.Info("Admin模块数据库迁移完成")
	return nil
}

// insertDefaultAdminData 插入默认管理员数据
func insertDefaultAdminData(db *gorm.DB) error {
	// 检查是否已存在管理员用户
	var count int64
	if err := db.Model(&models.AdminUser{}).Count(&count).Error; err != nil {
		return err
	}

	// 如果已有用户，跳过插入
	if count > 0 {
		return nil
	}

	// 生成正确的密码哈希值 (密码: admin123)
	// 使用bcrypt生成新的哈希值
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("生成密码哈希失败: %w", err)
	}

	// 创建默认管理员用户
	adminUser := &models.AdminUser{
		Username: "admin",
		Password: string(hashedPassword),
		Name:     "系统管理员",
		Status:   1, // 启用状态
	}

	if err := db.Create(adminUser).Error; err != nil {
		return err
	}

	pkgLogger.Info("已创建默认管理员用户: admin/admin123")
	return nil
}
