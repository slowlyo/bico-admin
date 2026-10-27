package handler

import (
	"net/http"
	"testing"

	"bico-admin/internal/admin/model"
)

// TestUserRoleAssignPrivilegeEscalation 验证普通管理员不能通过篡改 role_ids 绑定超管或自身未持有的角色。
func TestUserRoleAssignPrivilegeEscalation(t *testing.T) {
	db := newAdminTestDB(t)
	h := NewAdminUserHandler(db, &recordingInvalidator{})

	super := seedRole(t, db, "超级管理员", model.SuperAdminRoleCode)
	normal := seedRole(t, db, "运营", "ops")
	other := seedRole(t, db, "财务", "finance")

	operator := seedUser(t, db, "ops-user", normal)
	victim := seedUser(t, db, "victim")

	// 创建用户时不能绑定超级管理员角色。
	expectBizError(t, callHandlerAs(t, http.MethodPost, "", `{"username":"evil","password":"password123","role_ids":[`+
		itoa(super.ID)+`]}`, operator.ID, h.Create), "无权分配所选角色")

	// 创建用户时不能绑定自己未持有的角色。
	expectBizError(t, callHandlerAs(t, http.MethodPost, "", `{"username":"evil2","password":"password123","role_ids":[`+
		itoa(other.ID)+`]}`, operator.ID, h.Create), "无权分配所选角色")

	// 只能绑定自己拥有的角色。
	expectOK(t, callHandlerAs(t, http.MethodPost, "", `{"username":"ok-user","password":"password123","role_ids":[`+
		itoa(normal.ID)+`]}`, operator.ID, h.Create))

	// 更新他人用户时不能提权为超级管理员。
	expectBizError(t, callHandlerAs(t, http.MethodPut, itoa(victim.ID), `{"role_ids":[`+itoa(super.ID)+`]}`, operator.ID, h.Update),
		"无权分配所选角色")
	if n := countRows(t, db, &model.AdminUserRole{}, "user_id = ? AND role_id = ?", victim.ID, super.ID); n != 0 {
		t.Fatal("越权绑定超管角色应被回滚")
	}

	// 超级管理员可为他人绑定任意角色。
	admin := seedUser(t, db, "root", super)
	expectOK(t, callHandlerAs(t, http.MethodPut, itoa(victim.ID), `{"role_ids":[`+itoa(super.ID)+`]}`, admin.ID, h.Update))
	if n := countRows(t, db, &model.AdminUserRole{}, "user_id = ? AND role_id = ?", victim.ID, super.ID); n != 1 {
		t.Fatalf("超管应能绑定超管角色，实际关联数 %d", n)
	}
}
