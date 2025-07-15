package types

import (
	"time"

	"bico-admin/internal/shared/types"
)

// RoleListRequest 角色列表请求
type RoleListRequest struct {
	types.BasePageQuery
	Name   string `form:"name" json:"name"`     // 角色名称
	Code   string `form:"code" json:"code"`     // 角色代码
	Status *int   `form:"status" json:"status"` // 状态
}

// RoleCreateRequest 创建角色请求
type RoleCreateRequest struct {
	Name        string   `json:"name" binding:"required,min=1,max=100"`
	Code        string   `json:"code" binding:"required,min=1,max=50"`
	Description string   `json:"description" binding:"max=500"`
	Status      int      `json:"status" binding:"oneof=0 1"`
	Permissions []string `json:"permissions"` // 权限代码列表
}

// RoleUpdateRequest 更新角色请求
type RoleUpdateRequest struct {
	Name        string   `json:"name" binding:"required,min=1,max=100"`
	Description string   `json:"description" binding:"max=500"`
	Status      int      `json:"status" binding:"oneof=0 1"`
	Permissions []string `json:"permissions"` // 权限代码列表
}

// RoleResponse 角色响应
type RoleResponse struct {
	ID          uint                     `json:"id"`
	Name        string                   `json:"name"`
	Code        string                   `json:"code"`
	Description string                   `json:"description"`
	Status      int                      `json:"status"`
	StatusText  string                   `json:"status_text"`
	Permissions []RolePermissionResponse `json:"permissions"`
	UserCount   int64                    `json:"user_count"` // 拥有该角色的用户数量
	CanEdit     bool                     `json:"can_edit"`   // 是否可编辑
	CanDelete   bool                     `json:"can_delete"` // 是否可删除
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

// RolePermissionResponse 角色权限响应
type RolePermissionResponse struct {
	PermissionCode string `json:"permission_code"`
	PermissionName string `json:"permission_name"`
	Module         string `json:"module"`
	Level          int    `json:"level"`
}

// RoleAssignRequest 分配角色请求
type RoleAssignRequest struct {
	UserID  uint   `json:"user_id" binding:"required"`
	RoleIDs []uint `json:"role_ids" binding:"required,min=1"`
}

// RolePermissionUpdateRequest 更新角色权限请求
type RolePermissionUpdateRequest struct {
	Permissions []string `json:"permissions" binding:"required"` // 权限代码列表
}

// UserRoleResponse 用户角色响应
type UserRoleResponse struct {
	UserID    uint           `json:"user_id"`
	Username  string         `json:"username"`
	Name      string         `json:"name"`
	Roles     []RoleResponse `json:"roles"`
	CreatedAt time.Time      `json:"created_at"`
}

// PermissionTreeNode 权限树节点
type PermissionTreeNode struct {
	Module      string               `json:"module"`
	Name        string               `json:"name"`
	Permissions []PermissionTreeItem `json:"permissions"`
}

// PermissionTreeItem 权限树项目
type PermissionTreeItem struct {
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Level       int      `json:"level"`
	LevelText   string   `json:"level_text"`
	MenuSigns   []string `json:"menu_signs"`
	Buttons     []string `json:"buttons"`
	APIs        []string `json:"apis"`
	Selected    bool     `json:"selected"` // 是否已选中
}

// RoleStatsResponse 角色统计响应
type RoleStatsResponse struct {
	TotalRoles    int64 `json:"total_roles"`
	ActiveRoles   int64 `json:"active_roles"`
	InactiveRoles int64 `json:"inactive_roles"`
	TotalUsers    int64 `json:"total_users"`
}

// GetStatusText 获取状态文本
func (r *RoleResponse) GetStatusText() string {
	switch r.Status {
	case types.StatusActive:
		return "启用"
	case types.StatusInactive:
		return "禁用"
	default:
		return "未知"
	}
}

// GetLevelText 获取权限级别文本
func (p *PermissionTreeItem) GetLevelText() string {
	switch p.Level {
	case 1:
		return "查看"
	case 2:
		return "操作"
	case 3:
		return "管理"
	case 4:
		return "超级"
	default:
		return "未知"
	}
}
