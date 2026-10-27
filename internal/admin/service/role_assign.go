package service

import (
	"bico-admin/internal/admin/model"
	"bico-admin/internal/pkg/crud"

	"gorm.io/gorm"
)

// ErrRolesNotAssignable 表示操作者试图绑定自己无权分配的角色。
var ErrRolesNotAssignable = crud.NewBizError("无权分配所选角色")

// IsUserSuperAdmin 判断用户是否持有启用的超级管理员角色。
func IsUserSuperAdmin(db *gorm.DB, userID uint) (bool, error) {
	var count int64
	err := db.Table("admin_user_roles").
		Joins("JOIN admin_roles ON admin_user_roles.role_id = admin_roles.id").
		Where(
			"admin_user_roles.user_id = ? AND admin_roles.code = ? AND admin_roles.enabled = ?",
			userID, model.SuperAdminRoleCode, true,
		).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ListUserEnabledRoleIDs 返回用户当前绑定的启用角色 ID 列表。
func ListUserEnabledRoleIDs(db *gorm.DB, userID uint) ([]uint, error) {
	var roleIDs []uint
	err := db.Table("admin_user_roles").
		Joins("JOIN admin_roles ON admin_user_roles.role_id = admin_roles.id").
		Where("admin_user_roles.user_id = ? AND admin_roles.enabled = ?", userID, true).
		Pluck("admin_user_roles.role_id", &roleIDs).Error
	return roleIDs, err
}

// AssertCanAssignRoles 校验操作者是否有权把 roleIDs 绑定到目标用户。
//
// 超级管理员可分配任意角色；普通管理员只能分配自己已拥有的启用角色，与前端下拉口径一致。
func AssertCanAssignRoles(db *gorm.DB, operatorID uint, roleIDs []uint) error {
	unique := crud.UniqueUints(roleIDs)
	if len(unique) == 0 {
		return nil
	}

	isSuper, err := IsUserSuperAdmin(db, operatorID)
	if err != nil {
		return err
	}
	if isSuper {
		return nil
	}

	assignable, err := ListUserEnabledRoleIDs(db, operatorID)
	if err != nil {
		return err
	}
	assignableSet := make(map[uint]struct{}, len(assignable))
	for _, id := range assignable {
		assignableSet[id] = struct{}{}
	}
	for _, id := range unique {
		if _, ok := assignableSet[id]; !ok {
			return ErrRolesNotAssignable
		}
	}
	return nil
}
