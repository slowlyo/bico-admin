package service

import (
	"testing"

	"bico-admin/internal/admin/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// TestAssertCanAssignRoles 验证角色分配授权与超级管理员豁免规则。
func TestAssertCanAssignRoles(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("打开数据库失败: %v", err)
	}
	if err := db.AutoMigrate(&model.AdminUser{}, &model.AdminRole{}, &model.AdminUserRole{}); err != nil {
		t.Fatalf("迁移失败: %v", err)
	}

	super := model.AdminRole{Name: "超级管理员", Code: model.SuperAdminRoleCode, Enabled: true}
	normal := model.AdminRole{Name: "运营", Code: "ops", Enabled: true}
	other := model.AdminRole{Name: "财务", Code: "finance", Enabled: true}
	if err := db.Create(&super).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&normal).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&other).Error; err != nil {
		t.Fatal(err)
	}

	opsUser := model.AdminUser{Username: "ops", Password: "x", Enabled: true}
	root := model.AdminUser{Username: "root", Password: "x", Enabled: true}
	if err := db.Create(&opsUser).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&root).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&model.AdminUserRole{UserID: opsUser.ID, RoleID: normal.ID}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&model.AdminUserRole{UserID: root.ID, RoleID: super.ID}).Error; err != nil {
		t.Fatal(err)
	}

	if err := AssertCanAssignRoles(db, opsUser.ID, []uint{other.ID}); err != ErrRolesNotAssignable {
		t.Fatalf("应拒绝未持有角色，实际: %v", err)
	}
	if err := AssertCanAssignRoles(db, opsUser.ID, []uint{normal.ID}); err != nil {
		t.Fatalf("应允许自身角色: %v", err)
	}
	if err := AssertCanAssignRoles(db, root.ID, []uint{super.ID, other.ID}); err != nil {
		t.Fatalf("超管应可分配任意角色: %v", err)
	}
}
