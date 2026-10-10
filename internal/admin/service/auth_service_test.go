package service

import (
	"errors"
	"testing"

	"bico-admin/internal/admin/model"
	"bico-admin/internal/core/cache"
	"bico-admin/internal/pkg/jwt"
	"bico-admin/internal/pkg/password"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// TestChangePasswordInvalidatesExistingTokens 验证修改密码会废弃全部旧令牌。
func TestChangePasswordInvalidatesExistingTokens(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		// 数据库不可用时无法验证持久化令牌版本。
		t.Fatalf("创建测试数据库失败: %v", err)
	}
	if err := database.AutoMigrate(&model.AdminUser{}); err != nil {
		// 缺少用户表时后续认证断言没有意义。
		t.Fatalf("迁移测试数据库失败: %v", err)
	}

	hashed, err := password.Hash("old-password")
	if err != nil {
		// 测试用户必须使用与生产一致的哈希算法。
		t.Fatalf("生成密码哈希失败: %v", err)
	}
	user := model.AdminUser{Username: "tester", Password: hashed, Enabled: true}
	if err := database.Create(&user).Error; err != nil {
		// 用户未落库时无法签发有效令牌。
		t.Fatalf("创建测试用户失败: %v", err)
	}

	memoryCache := cache.NewMemoryCache()
	defer memoryCache.Close()
	manager := jwt.NewJWTManager("0123456789abcdef0123456789abcdef", 1)
	service := NewAuthService(database, manager, memoryCache)
	login, err := service.Login(&LoginRequest{Username: user.Username, Password: "old-password"})
	if err != nil {
		// 登录失败说明测试令牌未按真实流程签发。
		t.Fatalf("登录失败: %v", err)
	}
	claims, err := manager.ParseToken(login.Token)
	if err != nil {
		// 无法解析签发令牌时不能验证版本失效。
		t.Fatalf("解析令牌失败: %v", err)
	}
	if !service.IsTokenVersionValid(user.ID, claims.Version) {
		// 修改密码前令牌必须保持有效。
		t.Fatalf("新签发令牌应有效")
	}

	err = service.ChangePassword(user.ID, &ChangePasswordRequest{OldPassword: "old-password", NewPassword: "new-password"})
	if err != nil {
		// 密码修改失败时无法验证版本递增。
		t.Fatalf("修改密码失败: %v", err)
	}
	if service.IsTokenVersionValid(user.ID, claims.Version) {
		// 旧版本仍有效会让泄露令牌继续访问系统。
		t.Fatalf("修改密码后旧令牌仍有效")
	}
	session, err := service.LoadSession(user.ID, claims.Version)
	if err != nil {
		t.Fatalf("读取会话失败: %v", err)
	}
	if session.VersionOK {
		t.Fatalf("修改密码后会话版本仍有效")
	}
}

// TestLoadSessionUsesOneRow 验证版本和启用状态能从同一条用户记录读出。
func TestLoadSessionUsesOneRow(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("创建测试数据库失败: %v", err)
	}
	if err := database.AutoMigrate(&model.AdminUser{}); err != nil {
		t.Fatalf("迁移测试数据库失败: %v", err)
	}

	hashed, err := password.Hash("password")
	if err != nil {
		t.Fatalf("生成密码哈希失败: %v", err)
	}
	user := model.AdminUser{Username: "session-user", Password: hashed, Enabled: true}
	if err := database.Create(&user).Error; err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}

	memoryCache := cache.NewMemoryCache()
	defer memoryCache.Close()
	manager := jwt.NewJWTManager("0123456789abcdef0123456789abcdef", 1)
	authService := NewAuthService(database, manager, memoryCache)
	login, err := authService.Login(&LoginRequest{Username: user.Username, Password: "password"})
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	claims, err := manager.ParseToken(login.Token)
	if err != nil {
		t.Fatalf("解析令牌失败: %v", err)
	}

	session, err := authService.LoadSession(user.ID, claims.Version)
	if err != nil {
		t.Fatalf("读取会话失败: %v", err)
	}
	if !session.Found || !session.Enabled || !session.VersionOK {
		t.Fatalf("有效用户会话不正确: %+v", session)
	}

	if err := database.Model(&model.AdminUser{}).Where("id = ?", user.ID).Update("enabled", false).Error; err != nil {
		t.Fatalf("禁用用户失败: %v", err)
	}
	authService.InvalidateUserStatusCache(user.ID)
	session, err = authService.LoadSession(user.ID, claims.Version)
	if err != nil {
		t.Fatalf("刷新会话失败: %v", err)
	}
	if !session.Found || session.Enabled || !session.VersionOK {
		t.Fatalf("禁用后会话不正确: %+v", session)
	}
}

// TestUpdateProfileUsername 验证用户名可修改且不允许重复。
func TestUpdateProfileUsername(t *testing.T) {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		// 数据库不可用时无法验证唯一性约束。
		t.Fatalf("创建测试数据库失败: %v", err)
	}
	if err := database.AutoMigrate(
		&model.AdminUser{},
		&model.AdminRole{},
		&model.AdminUserRole{},
		&model.AdminRolePermission{},
	); err != nil {
		// 认证服务会读取用户与角色关联表，需要完整迁移。
		t.Fatalf("迁移测试数据库失败: %v", err)
	}

	users := []model.AdminUser{
		{Username: "current-user", Password: "unused", Name: "Current", Enabled: true},
		{Username: "existing-user", Password: "unused", Name: "Existing", Enabled: true},
	}
	if err := database.Create(&users).Error; err != nil {
		// 测试用户未落库时无法验证重复用户名。
		t.Fatalf("创建测试用户失败: %v", err)
	}

	memoryCache := cache.NewMemoryCache()
	defer memoryCache.Close()
	manager := jwt.NewJWTManager("0123456789abcdef0123456789abcdef", 1)
	authService := NewAuthService(database, manager, memoryCache)

	newUsername := "  renamed-user  "
	updated, err := authService.UpdateProfile(users[0].ID, &UpdateProfileRequest{Username: &newUsername})
	if err != nil {
		// 合法的新用户名应正常保存。
		t.Fatalf("修改用户名失败: %v", err)
	}
	if updated.Username != "renamed-user" {
		// 响应必须返回去除首尾空格后的新用户名。
		t.Fatalf("用户名未正确更新: %q", updated.Username)
	}

	duplicateUsername := users[1].Username
	_, err = authService.UpdateProfile(users[0].ID, &UpdateProfileRequest{Username: &duplicateUsername})
	if !errors.Is(err, ErrUsernameExists) {
		// 重复用户名必须返回稳定的业务错误。
		t.Fatalf("期望用户名已存在错误，实际为: %v", err)
	}
}
