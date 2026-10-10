package handler

import (
	"net/http"
	"testing"

	"bico-admin/internal/admin/model"
)

// TestUserCreateSyncsRolesAndReloads 验证创建时同步角色、返回带角色的详情，且非法角色整体回滚。
func TestUserCreateSyncsRolesAndReloads(t *testing.T) {
	db := newAdminTestDB(t)
	h := NewAdminUserHandler(db, &recordingInvalidator{})
	r1 := seedRole(t, db, "角色一", "r1")
	r2 := seedRole(t, db, "角色二", "r2")

	resp := callHandler(t, http.MethodPost, "", `{"username":"bob","password":"password123","role_ids":[`+
		itoa(r1.ID)+`,`+itoa(r2.ID)+`]}`, h.Create)
	expectOK(t, resp)
	data := resp["data"].(map[string]interface{})
	if roles, _ := data["roles"].([]interface{}); len(roles) != 2 {
		t.Fatalf("返回值应带 2 个角色，实际: %v", data["roles"])
	}
	if n := countRows(t, db, &model.AdminUserRole{}, "role_id IN ?", []uint{r1.ID, r2.ID}); n != 2 {
		t.Fatalf("应写入 2 条角色关联，实际 %d", n)
	}

	// 非法角色 ID：返回业务错误，用户记录随事务回滚。
	resp = callHandler(t, http.MethodPost, "", `{"username":"carl","password":"password123","role_ids":[9999]}`, h.Create)
	expectBizError(t, resp, "存在无效角色 ID")
	if n := countRows(t, db, &model.AdminUser{}, "username = ?", "carl"); n != 0 {
		t.Fatalf("角色非法时不应创建用户，实际 %d", n)
	}
}

// TestUserPasswordUpdateBumpsTokenVersion 验证改密码会递增 token 版本并失效对应缓存。
func TestUserPasswordUpdateBumpsTokenVersion(t *testing.T) {
	db := newAdminTestDB(t)
	inv := &recordingInvalidator{}
	h := NewAdminUserHandler(db, inv)
	user := seedUser(t, db, "alice")

	expectOK(t, callHandler(t, http.MethodPut, itoa(user.ID), `{"password":"newpassword1"}`, h.Update))

	var got model.AdminUser
	if err := db.First(&got, user.ID).Error; err != nil {
		t.Fatal(err)
	}
	if got.TokenVersion != 1 || got.Password == "old-hash" {
		t.Fatalf("改密后应递增 token 版本并更新哈希，实际 version=%d", got.TokenVersion)
	}
	if len(inv.userToken) != 1 || inv.userToken[0] != user.ID {
		t.Fatalf("应失效 token 版本缓存，实际 %v", inv.userToken)
	}

	// 不改密码的更新不应影响 token 版本缓存。
	expectOK(t, callHandler(t, http.MethodPut, itoa(user.ID), `{"name":"Alice"}`, h.Update))
	if len(inv.userToken) != 1 {
		t.Fatalf("未改密码不应再次失效 token 缓存，实际 %v", inv.userToken)
	}
}

// TestUserLastSuperAdminProtected 验证最后一个启用的超级管理员不能被禁用或删除。
func TestUserLastSuperAdminProtected(t *testing.T) {
	db := newAdminTestDB(t)
	inv := &recordingInvalidator{}
	h := NewAdminUserHandler(db, inv)
	super := seedRole(t, db, "超级管理员", model.SuperAdminRoleCode)
	first := seedUser(t, db, "root1", super)

	const msg = "必须保留至少一个启用的超级管理员"
	expectBizError(t, callHandler(t, http.MethodPut, itoa(first.ID), `{"enabled":false}`, h.UpdateEnabled), msg)
	expectBizError(t, callHandler(t, http.MethodDelete, itoa(first.ID), "", h.Delete), msg)
	if n := countRows(t, db, &model.AdminUser{}, "id = ? AND enabled = ?", first.ID, true); n != 1 {
		t.Fatal("被拒绝的操作不应改变用户")
	}

	// 有第二个启用的超级管理员后，禁用第一个应成功，并失效状态缓存。
	seedUser(t, db, "root2", super)
	expectOK(t, callHandler(t, http.MethodPut, itoa(first.ID), `{"enabled":false}`, h.UpdateEnabled))
	if len(inv.userStatus) != 1 || inv.userStatus[0] != first.ID {
		t.Fatalf("禁用后应失效状态缓存，实际 %v", inv.userStatus)
	}
}

// TestUserDeleteCleansRelationsAndInvalidates 验证删除会清理角色关联，并在提交后批量失效缓存。
func TestUserDeleteCleansRelationsAndInvalidates(t *testing.T) {
	db := newAdminTestDB(t)
	inv := &recordingInvalidator{}
	h := NewAdminUserHandler(db, inv)
	role := seedRole(t, db, "普通角色", "normal")
	u1 := seedUser(t, db, "u1", role)
	u2 := seedUser(t, db, "u2", role)
	u3 := seedUser(t, db, "u3", role)

	// 单个删除。
	expectOK(t, callHandler(t, http.MethodDelete, itoa(u1.ID), "", h.Delete))
	if n := countRows(t, db, &model.AdminUserRole{}, "user_id = ?", u1.ID); n != 0 {
		t.Fatalf("删除后角色关联应清理，实际 %d", n)
	}
	if len(inv.usersAuth) != 1 || len(inv.usersAuth[0]) != 1 || inv.usersAuth[0][0] != u1.ID {
		t.Fatalf("单个删除应失效该用户缓存，实际 %v", inv.usersAuth)
	}

	// 批量删除共用同一组 hook。
	expectOK(t, callHandler(t, http.MethodDelete, "", `{"ids":[`+itoa(u2.ID)+`,`+itoa(u3.ID)+`]}`, h.DeleteBatch))
	if n := countRows(t, db, &model.AdminUserRole{}, "role_id = ?", role.ID); n != 0 {
		t.Fatalf("批量删除后角色关联应清理，实际 %d", n)
	}
	if len(inv.usersAuth) != 2 || len(inv.usersAuth[1]) != 2 {
		t.Fatalf("批量删除应一次失效 2 个用户缓存，实际 %v", inv.usersAuth)
	}
}
