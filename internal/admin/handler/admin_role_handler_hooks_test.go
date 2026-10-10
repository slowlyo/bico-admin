package handler

import (
	"net/http"
	"strconv"
	"testing"

	"bico-admin/internal/admin/model"
	"bico-admin/internal/pkg/crud"
)

// itoa 把 uint 转成十进制字符串，便于拼接请求体和路径参数。
func itoa(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

// TestRoleSuperAdminProtected 验证保留角色不可修改、删除，权限也不可改。
func TestRoleSuperAdminProtected(t *testing.T) {
	db := newAdminTestDB(t)
	h := NewAdminRoleHandler(db, &recordingInvalidator{})
	super := seedRole(t, db, "超级管理员", model.SuperAdminRoleCode)

	expectBizError(t, callHandler(t, http.MethodPut, itoa(super.ID), `{"name":"改名"}`, h.Update), "超级管理员角色不可修改")
	expectBizError(t, callHandler(t, http.MethodDelete, itoa(super.ID), "", h.Delete), "超级管理员角色不可删除")
	expectBizError(t, callHandler(t, http.MethodPut, itoa(super.ID), `{"permissions":[]}`, h.UpdatePermissions), "超级管理员角色权限不可修改")
}

// TestRoleUniqueName 验证角色名称重复被拒绝。
func TestRoleUniqueName(t *testing.T) {
	db := newAdminTestDB(t)
	h := NewAdminRoleHandler(db, &recordingInvalidator{})

	expectOK(t, callHandler(t, http.MethodPost, "", `{"name":"运营"}`, h.Create))
	expectBizError(t, callHandler(t, http.MethodPost, "", `{"name":"运营"}`, h.Create), "角色名称已存在")
}

// TestRoleUpdatePermissions 验证覆盖权限、失败回滚和提交后失效缓存。
func TestRoleUpdatePermissions(t *testing.T) {
	// 临时替换全局权限树，保证用例不依赖其他模块是否已注册。
	original := crud.GetAllPermissions()
	crud.SetBasePermissions([]crud.Permission{
		{Key: "test:menu", Label: "测试", Children: []crud.Permission{
			{Key: "test:list", Label: "列表"},
			{Key: "test:edit", Label: "编辑"},
		}},
	})
	t.Cleanup(func() { crud.SetBasePermissions(original) })

	db := newAdminTestDB(t)
	inv := &recordingInvalidator{}
	h := NewAdminRoleHandler(db, inv)
	role := seedRole(t, db, "运营", "ops")
	id := itoa(role.ID)

	expectOK(t, callHandler(t, http.MethodPut, id, `{"permissions":["test:list","test:edit","test:list"]}`, h.UpdatePermissions))
	if n := countRows(t, db, &model.AdminRolePermission{}, "role_id = ?", role.ID); n != 2 {
		t.Fatalf("重复权限应去重，期望 2 条，实际 %d", n)
	}
	if len(inv.rolePermission) != 1 || inv.rolePermission[0] != role.ID {
		t.Fatalf("成功后应失效该角色用户缓存，实际 %v", inv.rolePermission)
	}

	// 非法权限：整体回滚，旧权限保留，且不失效缓存。
	resp := callHandler(t, http.MethodPut, id, `{"permissions":["bogus"]}`, h.UpdatePermissions)
	expectBizError(t, resp, "存在无效权限标识: bogus")
	if n := countRows(t, db, &model.AdminRolePermission{}, "role_id = ?", role.ID); n != 2 {
		t.Fatalf("失败应回滚并保留旧权限，实际 %d 条", n)
	}
	if len(inv.rolePermission) != 1 {
		t.Fatalf("失败时不应失效缓存，实际 %v", inv.rolePermission)
	}

	// 缺少 permissions 字段：拒绝，避免误清空权限。
	expectBizError(t, callHandler(t, http.MethodPut, id, `{}`, h.UpdatePermissions), "permissions 字段不能为空")
	if n := countRows(t, db, &model.AdminRolePermission{}, "role_id = ?", role.ID); n != 2 {
		t.Fatalf("缺少字段不应改动权限，实际 %d 条", n)
	}
}

// TestRoleDeleteCleansRelationsAndInvalidates 验证删除角色清理权限与用户关联，并失效关联用户缓存。
func TestRoleDeleteCleansRelationsAndInvalidates(t *testing.T) {
	db := newAdminTestDB(t)
	inv := &recordingInvalidator{}
	h := NewAdminRoleHandler(db, inv)
	role := seedRole(t, db, "运营", "ops")
	seedUser(t, db, "u1", role)
	if err := db.Create(&model.AdminRolePermission{RoleID: role.ID, Permission: "test:list"}).Error; err != nil {
		t.Fatal(err)
	}

	expectOK(t, callHandler(t, http.MethodDelete, itoa(role.ID), "", h.Delete))
	if n := countRows(t, db, &model.AdminRolePermission{}, "role_id = ?", role.ID); n != 0 {
		t.Fatalf("权限关联应清理，实际 %d", n)
	}
	if n := countRows(t, db, &model.AdminUserRole{}, "role_id = ?", role.ID); n != 0 {
		t.Fatalf("用户关联应清理，实际 %d", n)
	}
	if len(inv.rolesUsers) != 1 || len(inv.rolesUsers[0]) != 1 || inv.rolesUsers[0][0] != role.ID {
		t.Fatalf("应失效该角色关联用户的缓存，实际 %v", inv.rolesUsers)
	}
}

// TestRoleEnabledChangeInvalidates 验证只有启用状态变化才失效角色用户的权限缓存。
func TestRoleEnabledChangeInvalidates(t *testing.T) {
	db := newAdminTestDB(t)
	inv := &recordingInvalidator{}
	h := NewAdminRoleHandler(db, inv)
	role := seedRole(t, db, "运营", "ops")
	id := itoa(role.ID)

	// 只改描述不影响权限结果。
	expectOK(t, callHandler(t, http.MethodPut, id, `{"description":"说明"}`, h.Update))
	if len(inv.rolePermission) != 0 {
		t.Fatalf("改描述不应失效缓存，实际 %v", inv.rolePermission)
	}

	expectOK(t, callHandler(t, http.MethodPut, id, `{"enabled":false}`, h.UpdateEnabled))
	if len(inv.rolePermission) != 1 || inv.rolePermission[0] != role.ID {
		t.Fatalf("禁用角色应失效缓存，实际 %v", inv.rolePermission)
	}
}
