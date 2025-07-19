package service

import (
	"context"
	"errors"
	"strings"

	"bico-admin/internal/admin/definitions"
	"bico-admin/internal/admin/models"
	"bico-admin/internal/admin/repository"
	"bico-admin/internal/admin/types"
	sharedTypes "bico-admin/internal/shared/types"
	"bico-admin/pkg/utils"
)

// AdminRoleService 管理员角色服务接口
type AdminRoleService interface {
	GetRoleList(ctx context.Context, req *types.RoleListRequest) (*sharedTypes.PageResult, error)
	CreateRole(ctx context.Context, req *types.RoleCreateRequest) (*types.RoleResponse, error)
	UpdateRole(ctx context.Context, id uint, req *types.RoleUpdateRequest) (*types.RoleResponse, error)
	DeleteRole(ctx context.Context, id uint) error
	GetRoleByID(ctx context.Context, id uint) (*types.RoleResponse, error)
	GetPermissionTree(ctx context.Context, roleID *uint) ([]types.PermissionTreeNode, error)
	UpdateRolePermissions(ctx context.Context, id uint, req *types.RolePermissionUpdateRequest) error
	AssignRolesToUser(ctx context.Context, req *types.RoleAssignRequest) error

	GetUserPermissions(ctx context.Context, userID uint) ([]string, error)
	GetActiveRoles(ctx context.Context) ([]*models.AdminRole, error)
}

// adminRoleService 管理员角色服务实现
type adminRoleService struct {
	adminRoleRepo repository.AdminRoleRepository
	adminUserRepo repository.AdminUserRepository
}

// NewAdminRoleService 创建管理员角色服务
func NewAdminRoleService(adminRoleRepo repository.AdminRoleRepository, adminUserRepo repository.AdminUserRepository) AdminRoleService {
	return &adminRoleService{
		adminRoleRepo: adminRoleRepo,
		adminUserRepo: adminUserRepo,
	}
}

// GetRoleList 获取角色列表
func (s *adminRoleService) GetRoleList(ctx context.Context, req *types.RoleListRequest) (*sharedTypes.PageResult, error) {
	roles, total, err := s.adminRoleRepo.ListRoles(ctx, req)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var roleResponses []types.RoleResponse
	for _, role := range roles {
		roleResponse := s.convertToRoleResponse(ctx, *role)
		roleResponses = append(roleResponses, roleResponse)
	}

	return sharedTypes.NewPageResult(roleResponses, total, req.GetPage(), req.GetPageSize()), nil
}

// CreateRole 创建角色
func (s *adminRoleService) CreateRole(ctx context.Context, req *types.RoleCreateRequest) (*types.RoleResponse, error) {
	// 检查角色代码是否已存在
	if exists, err := s.adminRoleRepo.ExistsByCode(ctx, req.Code); err != nil {
		return nil, err
	} else if exists {
		return nil, errors.New("角色代码已存在")
	}

	// 验证权限代码
	if err := s.validatePermissions(req.Permissions); err != nil {
		return nil, err
	}

	// 创建角色
	role := &models.AdminRole{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      &req.Status,
	}

	if err := s.adminRoleRepo.Create(ctx, role); err != nil {
		return nil, err
	}

	// 创建角色权限关联
	if err := s.adminRoleRepo.CreateRolePermissions(ctx, nil, role.ID, req.Permissions); err != nil {
		return nil, err
	}

	// 重新查询角色信息
	updatedRole, err := s.adminRoleRepo.GetByID(ctx, role.ID)
	if err != nil {
		return nil, err
	}

	response := s.convertToRoleResponse(ctx, *updatedRole)
	return &response, nil
}

// UpdateRole 更新角色
func (s *adminRoleService) UpdateRole(ctx context.Context, id uint, req *types.RoleUpdateRequest) (*types.RoleResponse, error) {
	// 查找角色
	role, err := s.adminRoleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 检查是否为超级管理员角色
	if role.IsSuperAdminRole() {
		return nil, errors.New("超级管理员角色不可编辑")
	}

	// 验证权限代码
	if err := s.validatePermissions(req.Permissions); err != nil {
		return nil, err
	}

	// 更新角色基本信息
	role.Name = req.Name
	role.Description = req.Description
	role.Status = &req.Status

	if err := s.adminRoleRepo.Update(ctx, role); err != nil {
		return nil, err
	}

	// 删除旧的权限关联
	if err := s.adminRoleRepo.DeleteRolePermissions(ctx, nil, role.ID); err != nil {
		return nil, err
	}

	// 创建新的权限关联
	if err := s.adminRoleRepo.CreateRolePermissions(ctx, nil, role.ID, req.Permissions); err != nil {
		return nil, err
	}

	// 重新查询角色信息
	updatedRole, err := s.adminRoleRepo.GetByID(ctx, role.ID)
	if err != nil {
		return nil, err
	}

	response := s.convertToRoleResponse(ctx, *updatedRole)
	return &response, nil
}

// DeleteRole 删除角色
func (s *adminRoleService) DeleteRole(ctx context.Context, id uint) error {
	// 检查角色是否存在
	role, err := s.adminRoleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查是否为超级管理员角色
	if role.IsSuperAdminRole() {
		return errors.New("超级管理员角色不可删除")
	}

	// 检查是否有用户使用该角色
	userCount, err := s.adminRoleRepo.CountUsersByRoleID(ctx, id)
	if err != nil {
		return err
	}

	if userCount > 0 {
		return errors.New("该角色正在被用户使用，无法删除")
	}

	// 删除角色权限关联
	if err := s.adminRoleRepo.DeleteRolePermissions(ctx, nil, id); err != nil {
		return err
	}

	// 删除角色
	if err := s.adminRoleRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// GetRoleByID 根据ID获取角色
func (s *adminRoleService) GetRoleByID(ctx context.Context, id uint) (*types.RoleResponse, error) {
	role, err := s.adminRoleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := s.convertToRoleResponse(ctx, *role)
	return &response, nil
}

// GetPermissionTree 获取权限树
func (s *adminRoleService) GetPermissionTree(ctx context.Context, roleID *uint) ([]types.PermissionTreeNode, error) {
	// 获取所有权限定义
	permissionTree := definitions.GetAllPermissions()

	// 如果指定了角色ID，获取该角色的权限
	var rolePermissions map[string]bool
	if roleID != nil {
		role, err := s.adminRoleRepo.GetByID(ctx, *roleID)
		if err != nil {
			return nil, err
		}

		rolePermissions = make(map[string]bool)
		for _, permission := range role.Permissions {
			rolePermissions[permission.PermissionCode] = true
		}
	}

	// 构建无限极权限树
	var treeNodes []types.PermissionTreeNode
	for _, rootPerm := range permissionTree {
		node := s.buildPermissionTreeNode(rootPerm, rolePermissions)
		treeNodes = append(treeNodes, node)
	}

	return treeNodes, nil
}

// buildPermissionTreeNode 递归构建权限树节点
func (s *adminRoleService) buildPermissionTreeNode(perm definitions.Permission, rolePermissions map[string]bool) types.PermissionTreeNode {
	node := types.PermissionTreeNode{
		Key:      perm.Code,
		Title:    perm.Name,
		Type:     string(perm.Type),
		Selected: rolePermissions != nil && rolePermissions[perm.Code],
		Children: make([]types.PermissionTreeNode, 0),
	}

	// 递归处理子权限
	for _, child := range perm.Children {
		childNode := s.buildPermissionTreeNode(child, rolePermissions)
		node.Children = append(node.Children, childNode)
	}

	return node
}

// extractModuleFromPermissionCode 从权限代码中提取模块信息
func (s *adminRoleService) extractModuleFromPermissionCode(code string) string {
	// 权限代码格式：system.admin_user:list -> 返回 system.admin_user
	// 或者 system -> 返回 system
	if idx := strings.Index(code, ":"); idx != -1 {
		return code[:idx]
	}
	return code
}

// UpdateRolePermissions 更新角色权限
func (s *adminRoleService) UpdateRolePermissions(ctx context.Context, id uint, req *types.RolePermissionUpdateRequest) error {
	// 检查角色是否存在
	role, err := s.adminRoleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 检查是否为超级管理员角色
	if role.IsSuperAdminRole() {
		return errors.New("超级管理员角色权限不可修改")
	}

	// 验证权限代码
	if err := s.validatePermissions(req.Permissions); err != nil {
		return err
	}

	// 删除旧的权限关联
	if err := s.adminRoleRepo.DeleteRolePermissions(ctx, nil, id); err != nil {
		return err
	}

	// 创建新的权限关联
	if err := s.adminRoleRepo.CreateRolePermissions(ctx, nil, id, req.Permissions); err != nil {
		return err
	}

	return nil
}

// AssignRolesToUser 为管理员用户分配角色
func (s *adminRoleService) AssignRolesToUser(ctx context.Context, req *types.RoleAssignRequest) error {
	// 检查管理员用户是否存在
	_, err := s.adminUserRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return err
	}

	// 检查角色是否都存在且启用
	activeRoles, err := s.adminRoleRepo.ListActiveRoles(ctx)
	if err != nil {
		return err
	}

	activeRoleMap := make(map[uint]bool)
	for _, role := range activeRoles {
		activeRoleMap[role.ID] = true
	}

	for _, roleID := range req.RoleIDs {
		if !activeRoleMap[roleID] {
			return errors.New("部分角色不存在或已禁用")
		}
	}

	// 删除用户现有角色
	if err := s.adminRoleRepo.DeleteUserRoles(ctx, nil, req.UserID); err != nil {
		return err
	}

	// 分配新角色
	if err := s.adminRoleRepo.AssignRolesToUser(ctx, nil, req.UserID, req.RoleIDs); err != nil {
		return err
	}

	return nil
}

// GetUserPermissions 获取用户所有权限
func (s *adminRoleService) GetUserPermissions(ctx context.Context, userID uint) ([]string, error) {
	return s.adminRoleRepo.GetUserPermissions(ctx, userID)
}

// 私有方法

// validatePermissions 验证权限代码
func (s *adminRoleService) validatePermissions(permissionCodes []string) error {
	validPermissions := definitions.GetPermissionCodes()
	validPermissionSet := make(map[string]bool)
	for _, code := range validPermissions {
		validPermissionSet[code] = true
	}

	for _, code := range permissionCodes {
		if !validPermissionSet[code] {
			return errors.New("无效的权限代码: " + code)
		}
	}

	return nil
}

// convertToRoleResponse 转换为角色响应
func (s *adminRoleService) convertToRoleResponse(ctx context.Context, role models.AdminRole) types.RoleResponse {
	var permissions []types.RolePermissionResponse
	for _, permission := range role.Permissions {
		if permissionDef := definitions.GetPermissionByCode(permission.PermissionCode); permissionDef != nil {
			// 从权限代码中提取模块信息
			module := s.extractModuleFromPermissionCode(permissionDef.Code)
			permissions = append(permissions, types.RolePermissionResponse{
				PermissionCode: permission.PermissionCode,
				PermissionName: permissionDef.Name,
				Module:         module,
				Level:          permissionDef.Level,
			})
		}
	}

	// 获取用户数量
	userCount, _ := s.adminRoleRepo.CountUsersByRoleID(ctx, role.ID)

	status := 0
	if role.Status != nil {
		status = *role.Status
	}

	return types.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		Status:      status,
		StatusText:  s.getStatusText(status),
		Permissions: permissions,
		UserCount:   userCount,
		CanEdit:     role.CanBeEdited(),
		CanDelete:   role.CanBeDeleted(),
		CreatedAt:   utils.NewFormattedTime(role.CreatedAt),
		UpdatedAt:   utils.NewFormattedTime(role.UpdatedAt),
	}
}

// getStatusText 获取状态文本
func (s *adminRoleService) getStatusText(status int) string {
	switch status {
	case sharedTypes.StatusActive:
		return "启用"
	case sharedTypes.StatusInactive:
		return "禁用"
	default:
		return "未知"
	}
}

// GetActiveRoles 获取所有启用的角色
func (s *adminRoleService) GetActiveRoles(ctx context.Context) ([]*models.AdminRole, error) {
	return s.adminRoleRepo.ListActiveRoles(ctx)
}
