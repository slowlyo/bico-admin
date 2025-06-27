package model

// Role 角色模型
type Role struct {
	BaseModel
	Name        string       `json:"name" gorm:"uniqueIndex;size:50;not null" validate:"required,max=50"`
	Code        string       `json:"code" gorm:"uniqueIndex;size:50;not null" validate:"required,max=50"`
	Description string       `json:"description" gorm:"size:255"`
	Status      RoleStatus   `json:"status" gorm:"default:1"`
	
	// 关联关系
	Users       []User       `json:"users" gorm:"many2many:user_roles;"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}

// RoleStatus 角色状态
type RoleStatus int

const (
	RoleStatusInactive RoleStatus = 0 // 未激活
	RoleStatusActive   RoleStatus = 1 // 激活
)

// RoleCreateRequest 创建角色请求
type RoleCreateRequest struct {
	Name          string `json:"name" validate:"required,max=50"`
	Code          string `json:"code" validate:"required,max=50"`
	Description   string `json:"description" validate:"max=255"`
	PermissionIDs []uint `json:"permission_ids"`
}

// RoleUpdateRequest 更新角色请求
type RoleUpdateRequest struct {
	Name          string     `json:"name" validate:"max=50"`
	Code          string     `json:"code" validate:"max=50"`
	Description   string     `json:"description" validate:"max=255"`
	Status        RoleStatus `json:"status" validate:"oneof=0 1"`
	PermissionIDs []uint     `json:"permission_ids"`
}

// Permission 权限模型
type Permission struct {
	BaseModel
	Name        string           `json:"name" gorm:"uniqueIndex;size:50;not null" validate:"required,max=50"`
	Code        string           `json:"code" gorm:"uniqueIndex;size:100;not null" validate:"required,max=100"`
	Type        PermissionType   `json:"type" gorm:"default:1"`
	Resource    string           `json:"resource" gorm:"size:50"`
	Action      string           `json:"action" gorm:"size:50"`
	Description string           `json:"description" gorm:"size:255"`
	ParentID    *uint            `json:"parent_id" gorm:"index"`
	Sort        int              `json:"sort" gorm:"default:0"`
	Status      PermissionStatus `json:"status" gorm:"default:1"`
	
	// 关联关系
	Parent   *Permission  `json:"parent" gorm:"foreignKey:ParentID"`
	Children []Permission `json:"children" gorm:"foreignKey:ParentID"`
	Roles    []Role       `json:"roles" gorm:"many2many:role_permissions;"`
}

// PermissionType 权限类型
type PermissionType int

const (
	PermissionTypeMenu   PermissionType = 1 // 菜单
	PermissionTypeButton PermissionType = 2 // 按钮
	PermissionTypeAPI    PermissionType = 3 // API
)

// PermissionStatus 权限状态
type PermissionStatus int

const (
	PermissionStatusInactive PermissionStatus = 0 // 未激活
	PermissionStatusActive   PermissionStatus = 1 // 激活
)

// PermissionCreateRequest 创建权限请求
type PermissionCreateRequest struct {
	Name        string         `json:"name" validate:"required,max=50"`
	Code        string         `json:"code" validate:"required,max=100"`
	Type        PermissionType `json:"type" validate:"oneof=1 2 3"`
	Resource    string         `json:"resource" validate:"max=50"`
	Action      string         `json:"action" validate:"max=50"`
	Description string         `json:"description" validate:"max=255"`
	ParentID    *uint          `json:"parent_id"`
	Sort        int            `json:"sort"`
}

// PermissionUpdateRequest 更新权限请求
type PermissionUpdateRequest struct {
	Name        string           `json:"name" validate:"max=50"`
	Code        string           `json:"code" validate:"max=100"`
	Type        PermissionType   `json:"type" validate:"oneof=1 2 3"`
	Resource    string           `json:"resource" validate:"max=50"`
	Action      string           `json:"action" validate:"max=50"`
	Description string           `json:"description" validate:"max=255"`
	ParentID    *uint            `json:"parent_id"`
	Sort        int              `json:"sort"`
	Status      PermissionStatus `json:"status" validate:"oneof=0 1"`
}
