package model

// Role 角色模型
type Role struct {
	BaseModel
	Name        string     `json:"name" gorm:"uniqueIndex;size:50;not null" validate:"required,max=50"`
	Code        string     `json:"code" gorm:"uniqueIndex;size:50;not null" validate:"required,max=50"`
	Description string     `json:"description" gorm:"size:255"`
	Status      RoleStatus `json:"status" gorm:"default:1"`

	// 关联关系
	Users []User `json:"users" gorm:"many2many:user_roles;"`
	// 注意：权限关联通过RolePermission表管理，直接存储权限代码
}

// RoleStatus 角色状态
type RoleStatus int

const (
	RoleStatusInactive RoleStatus = 0 // 未激活
	RoleStatusActive   RoleStatus = 1 // 激活
)

// RoleCreateRequest 创建角色请求
type RoleCreateRequest struct {
	Name            string   `json:"name" validate:"required,max=50"`
	Code            string   `json:"code" validate:"required,max=50"`
	Description     string   `json:"description" validate:"max=255"`
	PermissionCodes []string `json:"permission_ids"` // 前端发送permission_ids，实际是权限代码数组
}

// RoleUpdateRequest 更新角色请求
type RoleUpdateRequest struct {
	Name            string     `json:"name" validate:"max=50"`
	Code            string     `json:"code" validate:"max=50"`
	Description     string     `json:"description" validate:"max=255"`
	Status          RoleStatus `json:"status" validate:"oneof=0 1"`
	PermissionCodes []string   `json:"permission_ids"` // 前端发送permission_ids，实际是权限代码数组
}

// 注意：Permission model已移除
// 权限定义现在完全基于代码，位于 backend/core/permission/config.go
// 数据库中只存储角色与权限代码的关联关系

// RolePermission 角色权限关联模型
// 直接存储权限代码，不依赖权限表
type RolePermission struct {
	RoleID         uint   `json:"role_id" gorm:"primaryKey"`
	PermissionCode string `json:"permission_code" gorm:"primaryKey;size:100"`

	// 关联关系
	Role Role `json:"role" gorm:"foreignKey:RoleID"`
}
