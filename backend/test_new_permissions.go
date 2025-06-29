package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"bico-admin/core/config"
	"bico-admin/core/permission"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 初始化配置
	cfg := config.New()

	// 初始化数据库
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("=== 新权限系统测试 ===")

	// 创建权限中间件实例
	permissionMiddleware := permission.NewPermissionMiddleware(db)

	// 测试用户ID（假设存在）
	testUserID := uint(1)

	// 测试权限检查
	fmt.Println("\n1. 测试权限检查:")
	testPermissions := []string{
		"system:view",
		"system:manage",
		"user:view",
		"user:create",
		"user:delete",
		"role:view",
		"role:create",
		"profile:view", // 这个应该总是返回true（权限豁免）
	}

	for _, perm := range testPermissions {
		hasPermission, err := permissionMiddleware.HasPermission(testUserID, perm)
		if err != nil {
			fmt.Printf("  权限 %s: 检查失败 - %v\n", perm, err)
		} else {
			fmt.Printf("  权限 %s: %v\n", perm, hasPermission)
		}
	}

	// 测试超级管理员检查
	fmt.Println("\n2. 测试超级管理员检查:")
	isSuperAdmin, err := permissionMiddleware.IsSuperAdmin(testUserID)
	if err != nil {
		fmt.Printf("  检查失败: %v\n", err)
	} else {
		fmt.Printf("  用户 %d 是否为超级管理员: %v\n", testUserID, isSuperAdmin)
	}

	// 测试个人资料路由豁免
	fmt.Println("\n3. 测试个人资料路由豁免:")
	profileRoutes := []string{
		"/auth/profile",
		"/auth/change-password",
		"profile:view",
		"profile:update",
		"/admin/users", // 这个不应该被豁免
	}

	for _, route := range profileRoutes {
		isProfile := permissionMiddleware.IsProfileRoute(route)
		fmt.Printf("  路由 %s 是否为个人资料路由: %v\n", route, isProfile)
	}

	fmt.Println("\n=== 测试完成 ===")
}

// 为了测试，我们需要添加一些方法到权限中间件
// 这些方法应该是公开的，以便测试使用

// 扩展权限中间件以支持测试
type TestPermissionMiddleware struct {
	*permission.PermissionMiddleware
	db *gorm.DB
}

func NewTestPermissionMiddleware(db *gorm.DB) *TestPermissionMiddleware {
	return &TestPermissionMiddleware{
		PermissionMiddleware: permission.NewPermissionMiddleware(db),
		db:                   db,
	}
}

// HasPermission 公开权限检查方法用于测试
func (pm *TestPermissionMiddleware) HasPermission(userID uint, perm string) (bool, error) {
	// 这里需要调用私有方法，但由于Go的包访问限制，我们需要重新实现
	// 或者在权限中间件中添加公开的测试方法
	
	// 检查是否为个人资料相关权限（无需验证）
	if pm.IsProfileRoute(perm) {
		return true, nil
	}

	// 检查用户是否为超级管理员
	isSuperAdmin, err := pm.IsSuperAdmin(userID)
	if err != nil {
		return false, err
	}
	
	if isSuperAdmin {
		return true, nil
	}
	
	// 查询用户是否有该权限
	var count int64
	err = pm.db.Table("permissions p").
		Joins("JOIN role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN user_roles ur ON rp.role_id = ur.role_id").
		Where("ur.user_id = ? AND p.code = ? AND p.status = ?", 
			userID, perm, 1). // 1 = PermissionStatusActive
		Count(&count).Error
	
	return count > 0, err
}

// IsSuperAdmin 公开超级管理员检查方法用于测试
func (pm *TestPermissionMiddleware) IsSuperAdmin(userID uint) (bool, error) {
	var count int64
	err := pm.db.Table("users u").
		Joins("JOIN user_roles ur ON u.id = ur.user_id").
		Joins("JOIN roles r ON ur.role_id = r.id").
		Where("u.id = ? AND r.code = ? AND r.status = ?", 
			userID, "super_admin", 1). // 1 = RoleStatusActive
		Count(&count).Error
	
	if err != nil {
		// 如果查询失败，尝试通过用户表的role字段检查（兼容旧版本）
		var user struct {
			Role string
		}
		if err := pm.db.Table("users").Select("role").Where("id = ?", userID).First(&user).Error; err != nil {
			return false, err
		}
		return user.Role == "super_admin", nil
	}
	
	return count > 0, nil
}

// IsProfileRoute 公开个人资料路由检查方法用于测试
func (pm *TestPermissionMiddleware) IsProfileRoute(permission string) bool {
	profileRoutes := []string{
		"/auth/profile",
		"/auth/change-password",
		"profile:view",
		"profile:update", 
		"profile:change_password",
	}
	
	for _, route := range profileRoutes {
		if permission == route {
			return true
		}
	}
	return false
}
