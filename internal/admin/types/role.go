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

// PermissionTreeNode 权限树节点（无限极树结构）
type PermissionTreeNode struct {
	Key      string               `json:"key"`      // 权限代码或模块代码
	Title    string               `json:"title"`    // 显示名称
	Type     string               `json:"type"`     // 类型：module 或 action
	Selected bool                 `json:"selected"` // 是否选中
	Children []PermissionTreeNode `json:"children"` // 子节点
}

// RoleStatsResponse 角色统计响应
type RoleStatsResponse struct {
	TotalRoles    int64 `json:"total_roles"`
	ActiveRoles   int64 `json:"active_roles"`
	InactiveRoles int64 `json:"inactive_roles"`
	TotalUsers    int64 `json:"total_users"`
}

// RoleOptionResponse 角色选项响应（用于下拉选择）
type RoleOptionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
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
