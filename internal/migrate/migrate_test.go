package migrate

import (
	"testing"

	adminModel "bico-admin/internal/admin/model"
	"bico-admin/internal/pkg/password"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// TestAutoMigrateInitializesDefaultAdmin 验证首次迁移自动创建默认管理员。
func TestAutoMigrateInitializesDefaultAdmin(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		// 数据库不可用时无法验证初始化事务。
		t.Fatalf("创建测试数据库失败: %v", err)
	}
	t.Setenv("BICO_ADMIN_INITIAL_PASSWORD", "")

	if err := AutoMigrate(database, "release"); err != nil {
		// 初始密码不再依赖环境变量，生产迁移应能直接完成。
		t.Fatalf("执行迁移失败: %v", err)
	}

	var admin adminModel.AdminUser
	if err := database.Where("username = ?", "admin").First(&admin).Error; err != nil {
		// 查询不到账号说明初始化没有实际落库。
		t.Fatalf("查询默认管理员失败: %v", err)
	}
	if !password.Verify(admin.Password, "admin") {
		// 初始密码必须保持与系统默认登录凭据一致。
		t.Fatal("默认管理员密码错误")
	}
}

// TestAutoMigrateAssignsSuperAdminRole 验证首个管理员通过保留角色获得权限。
func TestAutoMigrateAssignsSuperAdminRole(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		// 数据库不可用时无法验证角色关系。
		t.Fatalf("创建测试数据库失败: %v", err)
	}
	if err := AutoMigrate(database, "release"); err != nil {
		// 合法初始密码应完成全部迁移。
		t.Fatalf("执行迁移失败: %v", err)
	}

	var relationCount int64
	if err := database.Table("admin_user_roles").
		Joins("JOIN admin_roles ON admin_user_roles.role_id = admin_roles.id").
		Joins("JOIN admin_users ON admin_user_roles.user_id = admin_users.id").
		Where("admin_users.username = ? AND admin_roles.code = ?", "admin", adminModel.SuperAdminRoleCode).
		Count(&relationCount).Error; err != nil {
		// 关联查询失败说明迁移结构不完整。
		t.Fatalf("查询超级管理员关联失败: %v", err)
	}
	if relationCount != 1 {
		// 首个管理员必须且只能关联一次保留角色。
		t.Fatalf("超级管理员关联数量错误: %d", relationCount)
	}
}
