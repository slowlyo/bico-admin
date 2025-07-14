package service

import (
	"context"
	"errors"

	"bico-admin/internal/admin/definitions"
	"bico-admin/internal/admin/models"
	"bico-admin/internal/admin/repository"
	"bico-admin/internal/admin/types"
	sharedTypes "bico-admin/internal/shared/types"
)

// AdminRoleService 管理员角色服务接口
type AdminRoleService interface {
	GetRoleList(ctx context.Context, req *types.RoleListRequest) (*sharedTypes.PageResult, error)
	CreateRole(ctx context.Context, req *types.RoleCreateRequest) (*types.RoleResponse, error)
	UpdateRole(ctx context.Context, id uint, req *types.RoleUpdateRequest) (*types.RoleResponse, error)
	DeleteRole(ctx context.Context, id uint) error
	GetRoleByID(ctx context.Context, id uint) (*types.RoleResponse, error)
	GetPermissionTree(ctx context.Context, roleID *uint) ([]types.PermissionTreeNode, error)
	AssignRolesToUser(ctx context.Context, req *types.RoleAssignRequest) error
	GetUserRoles(ctx context.Context, userID uint) (*types.UserRoleResponse, error)
	GetUserPermissions(ctx context.Context, userID uint) ([]string, error)
	GetRoleStats(ctx context.Context) (*types.RoleStatsResponse, error)
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
	roles, total, err := s.adminRoleRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	// 转换响应
	var roleResponses []types.RoleResponse
	for _, role := range roles {
		roleResponse := s.convertToRoleResponse(ctx, *role)
		roleResponses = append(roleResponses, roleResponse)
	}

	return sharedTypes.NewPageResult(roleResponses, total, req.Page, req.PageSize), nil
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
		Status:      req.Status,
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

	// 验证权限代码
	if err := s.validatePermissions(req.Permissions); err != nil {
		return nil, err
	}

	// 更新角色基本信息
	role.Name = req.Name
	role.Description = req.Description
	role.Status = req.Status

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
	_, err := s.adminRoleRepo.GetByID(ctx, id)
	if err != nil {
		return err
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
	permissionGroups := definitions.GetAllPermissions()

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

	// 构建权限树
	var treeNodes []types.PermissionTreeNode
	for _, group := range permissionGroups {
		node := types.PermissionTreeNode{
			Module: group.Module,
			Name:   group.Name,
		}

		for _, permission := range group.Permissions {
			item := types.PermissionTreeItem{
				Code:      permission.Code,
				Name:      permission.Name,
				Level:     permission.Level,
				LevelText: s.getLevelText(permission.Level),
				Buttons:   permission.Buttons,
				APIs:      permission.APIs,
				Selected:  rolePermissions != nil && rolePermissions[permission.Code],
			}
			node.Permissions = append(node.Permissions, item)
		}

		treeNodes = append(treeNodes, node)
	}

	return treeNodes, nil
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

// GetUserRoles 获取管理员用户角色
func (s *adminRoleService) GetUserRoles(ctx context.Context, userID uint) (*types.UserRoleResponse, error) {
	user, err := s.adminUserRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	roles, err := s.adminRoleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	var roleResponses []types.RoleResponse
	for _, role := range roles {
		roleResponse := s.convertToRoleResponse(ctx, *role)
		roleResponses = append(roleResponses, roleResponse)
	}

	return &types.UserRoleResponse{
		UserID:   user.ID,
		Username: user.Username,
		Name:     user.Name,
		Roles:    roleResponses,
	}, nil
}

// GetUserPermissions 获取用户所有权限
func (s *adminRoleService) GetUserPermissions(ctx context.Context, userID uint) ([]string, error) {
	return s.adminRoleRepo.GetUserPermissions(ctx, userID)
}

// GetRoleStats 获取角色统计
func (s *adminRoleService) GetRoleStats(ctx context.Context) (*types.RoleStatsResponse, error) {
	var stats types.RoleStatsResponse

	// 总角色数
	totalRoles, err := s.adminRoleRepo.Count(ctx)
	if err != nil {
		return nil, err
	}
	stats.TotalRoles = totalRoles

	// 启用角色数
	activeRoles, err := s.adminRoleRepo.CountByStatus(ctx, sharedTypes.StatusActive)
	if err != nil {
		return nil, err
	}
	stats.ActiveRoles = activeRoles

	// 禁用角色数
	stats.InactiveRoles = stats.TotalRoles - stats.ActiveRoles

	// 拥有角色的用户数
	totalUsers, err := s.adminRoleRepo.CountUsersWithRoles(ctx)
	if err != nil {
		return nil, err
	}
	stats.TotalUsers = totalUsers

	return &stats, nil
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
			permissions = append(permissions, types.RolePermissionResponse{
				PermissionCode: permission.PermissionCode,
				PermissionName: permissionDef.Name,
				Module:         permissionDef.Module,
				Level:          permissionDef.Level,
			})
		}
	}

	// 获取用户数量
	userCount, _ := s.adminRoleRepo.CountUsersByRoleID(ctx, role.ID)

	return types.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		Status:      role.Status,
		StatusText:  s.getStatusText(role.Status),
		Permissions: permissions,
		UserCount:   userCount,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
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

// getLevelText 获取权限级别文本
func (s *adminRoleService) getLevelText(level int) string {
	switch level {
	case definitions.PermissionLevelView:
		return "查看"
	case definitions.PermissionLevelAction:
		return "操作"
	case definitions.PermissionLevelManage:
		return "管理"
	case definitions.PermissionLevelSuper:
		return "超级"
	default:
		return "未知"
	}
}
