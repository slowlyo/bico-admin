package admin

import (
	"fmt"

	"github.com/google/wire"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"bico-admin/internal/admin/definitions"
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
	repository.NewAdminRoleRepository,

	// Service层
	service.NewUserService,
	service.NewAdminUserService,
	service.NewAdminRoleService,
	service.NewAuthService,

	// Handler层
	handler.NewAuthHandler,
	handler.NewUserHandler,
	handler.NewAdminUserHandler,
	handler.NewAdminRoleHandler,
	handler.NewSystemHandler,

	// 路由处理器集合
	ProvideHandlers,
)

// ProvideHandlers 提供处理器集合
func ProvideHandlers(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	adminUserHandler *handler.AdminUserHandler,
	adminRoleHandler *handler.AdminRoleHandler,
	systemHandler *handler.SystemHandler,
) *routes.Handlers {
	return &routes.Handlers{
		AuthHandler:      authHandler,
		UserHandler:      userHandler,
		AdminUserHandler: adminUserHandler,
		AdminRoleHandler: adminRoleHandler,
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
	var userCount int64
	if err := db.Model(&models.AdminUser{}).Count(&userCount).Error; err != nil {
		return err
	}

	// 如果已有用户，跳过插入
	if userCount > 0 {
		return nil
	}

	// 1. 创建默认角色
	if err := createDefaultRoles(db); err != nil {
		return fmt.Errorf("创建默认角色失败: %w", err)
	}

	// 2. 创建默认管理员用户
	adminUser, err := createDefaultAdminUser(db)
	if err != nil {
		return fmt.Errorf("创建默认管理员用户失败: %w", err)
	}

	// 3. 为默认管理员分配超级管理员角色
	if err := assignSuperAdminRole(db, adminUser.ID); err != nil {
		return fmt.Errorf("分配超级管理员角色失败: %w", err)
	}

	pkgLogger.Info("已创建默认管理员用户: admin/admin123")
	return nil
}

// createDefaultRoles 创建默认角色
func createDefaultRoles(db *gorm.DB) error {
	// 检查是否已存在角色
	var roleCount int64
	if err := db.Model(&models.AdminRole{}).Count(&roleCount).Error; err != nil {
		return err
	}

	// 如果已有角色，跳过创建
	if roleCount > 0 {
		return nil
	}

	// 创建超级管理员角色
	superAdminRole := &models.AdminRole{
		Name:        "超级管理员",
		Code:        models.RoleCodeSuperAdmin,
		Description: "拥有系统所有权限的超级管理员",
		Status:      1,
	}

	if err := db.Create(superAdminRole).Error; err != nil {
		return err
	}

	// 为超级管理员角色分配所有权限
	if err := assignAllPermissionsToRole(db, superAdminRole.ID); err != nil {
		return fmt.Errorf("为超级管理员角色分配权限失败: %w", err)
	}

	pkgLogger.Info("已创建默认角色")
	return nil
}

// createDefaultAdminUser 创建默认管理员用户
func createDefaultAdminUser(db *gorm.DB) (*models.AdminUser, error) {
	// 生成密码哈希值 (密码: admin123)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("生成密码哈希失败: %w", err)
	}

	// 创建默认管理员用户
	adminUser := &models.AdminUser{
		Username: "admin",
		Password: string(hashedPassword),
		Name:     "系统管理员",
		Status:   1, // 启用状态
	}

	if err := db.Create(adminUser).Error; err != nil {
		return nil, err
	}

	return adminUser, nil
}

// assignSuperAdminRole 为用户分配超级管理员角色
func assignSuperAdminRole(db *gorm.DB, userID uint) error {
	// 查找超级管理员角色
	var superAdminRole models.AdminRole
	if err := db.Where("code = ?", models.RoleCodeSuperAdmin).First(&superAdminRole).Error; err != nil {
		return fmt.Errorf("查找超级管理员角色失败: %w", err)
	}

	// 创建用户角色关联
	userRole := &models.AdminUserRole{
		UserID: userID,
		RoleID: superAdminRole.ID,
	}

	if err := db.Create(userRole).Error; err != nil {
		return err
	}

	return nil
}

// assignAllPermissionsToRole 为角色分配所有权限
func assignAllPermissionsToRole(db *gorm.DB, roleID uint) error {
	// 获取所有权限代码（只分配操作类型的权限）
	allPermissions := definitions.GetAllPermissionsFlat()

	var rolePermissions []models.AdminRolePermission
	for _, perm := range allPermissions {
		// 只为操作类型的权限创建角色权限关联
		if perm.Type == definitions.PermissionTypeAction {
			rolePermissions = append(rolePermissions, models.AdminRolePermission{
				RoleID:         roleID,
				PermissionCode: perm.Code,
			})
		}
	}

	// 批量插入权限
	if len(rolePermissions) > 0 {
		if err := db.Create(&rolePermissions).Error; err != nil {
			return err
		}
	}

	return nil
}
