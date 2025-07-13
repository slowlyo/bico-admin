package service

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"bico-admin/internal/admin/definitions"
	"bico-admin/internal/admin/models"
	"bico-admin/internal/admin/types"
	sharedTypes "bico-admin/internal/shared/types"
)

// AdminRoleService 管理员角色服务
type AdminRoleService struct {
	db *gorm.DB
}

// NewAdminRoleService 创建管理员角色服务
func NewAdminRoleService(db *gorm.DB) *AdminRoleService {
	return &AdminRoleService{db: db}
}

// GetRoleList 获取角色列表
func (s *AdminRoleService) GetRoleList(req *types.RoleListRequest) (*sharedTypes.PageResult, error) {
	var roles []models.AdminRole
	var total int64

	query := s.db.Model(&models.AdminRole{})

	// 条件过滤
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		query = query.Where("code LIKE ?", "%"+req.Code+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取角色总数失败: %w", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("Permissions").
		Offset(offset).
		Limit(req.PageSize).
		Order("created_at DESC").
		Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("获取角色列表失败: %w", err)
	}

	// 转换响应
	var roleResponses []types.RoleResponse
	for _, role := range roles {
		roleResponse := s.convertToRoleResponse(role)
		roleResponses = append(roleResponses, roleResponse)
	}

	return sharedTypes.NewPageResult(roleResponses, total, req.Page, req.PageSize), nil
}

// CreateRole 创建角色
func (s *AdminRoleService) CreateRole(req *types.RoleCreateRequest) (*types.RoleResponse, error) {
	// 检查角色代码是否已存在
	var existingRole models.AdminRole
	if err := s.db.Where("code = ?", req.Code).First(&existingRole).Error; err == nil {
		return nil, errors.New("角色代码已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("检查角色代码失败: %w", err)
	}

	// 验证权限代码
	if err := s.validatePermissions(req.Permissions); err != nil {
		return nil, err
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建角色
	role := models.AdminRole{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := tx.Create(&role).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建角色失败: %w", err)
	}

	// 创建角色权限关联
	if err := s.createRolePermissions(tx, role.ID, req.Permissions); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	// 重新查询角色信息
	if err := s.db.Preload("Permissions").First(&role, role.ID).Error; err != nil {
		return nil, fmt.Errorf("查询角色信息失败: %w", err)
	}

	response := s.convertToRoleResponse(role)
	return &response, nil
}

// UpdateRole 更新角色
func (s *AdminRoleService) UpdateRole(id uint, req *types.RoleUpdateRequest) (*types.RoleResponse, error) {
	// 查找角色
	var role models.AdminRole
	if err := s.db.Preload("Permissions").First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}

	// 验证权限代码
	if err := s.validatePermissions(req.Permissions); err != nil {
		return nil, err
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新角色基本信息
	role.Name = req.Name
	role.Description = req.Description
	role.Status = req.Status

	if err := tx.Save(&role).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新角色失败: %w", err)
	}

	// 删除旧的权限关联
	if err := tx.Where("role_id = ?", role.ID).Delete(&models.AdminRolePermission{}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("删除旧权限关联失败: %w", err)
	}

	// 创建新的权限关联
	if err := s.createRolePermissions(tx, role.ID, req.Permissions); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	// 重新查询角色信息
	if err := s.db.Preload("Permissions").First(&role, role.ID).Error; err != nil {
		return nil, fmt.Errorf("查询角色信息失败: %w", err)
	}

	response := s.convertToRoleResponse(role)
	return &response, nil
}

// DeleteRole 删除角色
func (s *AdminRoleService) DeleteRole(id uint) error {
	// 检查角色是否存在
	var role models.AdminRole
	if err := s.db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		return fmt.Errorf("查询角色失败: %w", err)
	}

	// 检查是否有用户使用该角色
	var userCount int64
	if err := s.db.Model(&models.AdminUserRole{}).Where("role_id = ?", id).Count(&userCount).Error; err != nil {
		return fmt.Errorf("检查角色使用情况失败: %w", err)
	}

	if userCount > 0 {
		return errors.New("该角色正在被用户使用，无法删除")
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除角色权限关联
	if err := tx.Where("role_id = ?", id).Delete(&models.AdminRolePermission{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除角色权限关联失败: %w", err)
	}

	// 删除角色
	if err := tx.Delete(&role).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除角色失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// GetRoleByID 根据ID获取角色
func (s *AdminRoleService) GetRoleByID(id uint) (*types.RoleResponse, error) {
	var role models.AdminRole
	if err := s.db.Preload("Permissions").First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		return nil, fmt.Errorf("查询角色失败: %w", err)
	}

	response := s.convertToRoleResponse(role)
	return &response, nil
}

// GetPermissionTree 获取权限树
func (s *AdminRoleService) GetPermissionTree(roleID *uint) ([]types.PermissionTreeNode, error) {
	// 获取所有权限定义
	permissionGroups := definitions.GetAllPermissions()

	// 如果指定了角色ID，获取该角色的权限
	var rolePermissions map[string]bool
	if roleID != nil {
		var role models.AdminRole
		if err := s.db.Preload("Permissions").First(&role, *roleID).Error; err != nil {
			return nil, fmt.Errorf("查询角色权限失败: %w", err)
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
				MenuSigns: permission.MenuSigns,
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
func (s *AdminRoleService) AssignRolesToUser(req *types.RoleAssignRequest) error {
	// 检查管理员用户是否存在
	var user models.AdminUser
	if err := s.db.First(&user, req.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("管理员用户不存在")
		}
		return fmt.Errorf("查询管理员用户失败: %w", err)
	}

	// 检查角色是否都存在
	var roleCount int64
	if err := s.db.Model(&models.AdminRole{}).Where("id IN ? AND status = ?", req.RoleIDs, sharedTypes.StatusActive).Count(&roleCount).Error; err != nil {
		return fmt.Errorf("检查角色失败: %w", err)
	}

	if int(roleCount) != len(req.RoleIDs) {
		return errors.New("部分角色不存在或已禁用")
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除用户现有角色
	if err := tx.Where("user_id = ?", req.UserID).Delete(&models.AdminUserRole{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除用户现有角色失败: %w", err)
	}

	// 分配新角色
	for _, roleID := range req.RoleIDs {
		userRole := models.AdminUserRole{
			UserID: req.UserID,
			RoleID: roleID,
		}
		if err := tx.Create(&userRole).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("分配角色失败: %w", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// GetUserRoles 获取管理员用户角色
func (s *AdminRoleService) GetUserRoles(userID uint) (*types.UserRoleResponse, error) {
	var user models.AdminUser
	if err := s.db.Preload("Roles.Role.Permissions").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("管理员用户不存在")
		}
		return nil, fmt.Errorf("查询管理员用户失败: %w", err)
	}

	var roles []types.RoleResponse
	for _, userRole := range user.Roles {
		roleResponse := s.convertToRoleResponse(userRole.Role)
		roles = append(roles, roleResponse)
	}

	return &types.UserRoleResponse{
		UserID:   user.ID,
		Username: user.Username,
		Name:     user.Name,
		Roles:    roles,
	}, nil
}

// GetUserPermissions 获取用户所有权限
func (s *AdminRoleService) GetUserPermissions(userID uint) ([]string, error) {
	var userRoles []models.AdminUserRole
	if err := s.db.Preload("Role.Permissions").Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, fmt.Errorf("查询用户角色失败: %w", err)
	}

	permissionSet := make(map[string]bool)
	for _, userRole := range userRoles {
		if userRole.Role.IsEnabled() {
			for _, permission := range userRole.Role.Permissions {
				permissionSet[permission.PermissionCode] = true
			}
		}
	}

	var permissions []string
	for permission := range permissionSet {
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// GetRoleStats 获取角色统计
func (s *AdminRoleService) GetRoleStats() (*types.RoleStatsResponse, error) {
	var stats types.RoleStatsResponse

	// 总角色数
	if err := s.db.Model(&models.AdminRole{}).Count(&stats.TotalRoles).Error; err != nil {
		return nil, fmt.Errorf("获取总角色数失败: %w", err)
	}

	// 启用角色数
	if err := s.db.Model(&models.AdminRole{}).Where("status = ?", sharedTypes.StatusActive).Count(&stats.ActiveRoles).Error; err != nil {
		return nil, fmt.Errorf("获取启用角色数失败: %w", err)
	}

	// 禁用角色数
	stats.InactiveRoles = stats.TotalRoles - stats.ActiveRoles

	// 拥有角色的用户数
	if err := s.db.Model(&models.AdminUserRole{}).Distinct("user_id").Count(&stats.TotalUsers).Error; err != nil {
		return nil, fmt.Errorf("获取用户数失败: %w", err)
	}

	return &stats, nil
}

// 私有方法

// validatePermissions 验证权限代码
func (s *AdminRoleService) validatePermissions(permissionCodes []string) error {
	validPermissions := definitions.GetPermissionCodes()
	validPermissionSet := make(map[string]bool)
	for _, code := range validPermissions {
		validPermissionSet[code] = true
	}

	for _, code := range permissionCodes {
		if !validPermissionSet[code] {
			return fmt.Errorf("无效的权限代码: %s", code)
		}
	}

	return nil
}

// createRolePermissions 创建角色权限关联
func (s *AdminRoleService) createRolePermissions(tx *gorm.DB, roleID uint, permissionCodes []string) error {
	for _, code := range permissionCodes {
		rolePermission := models.AdminRolePermission{
			RoleID:         roleID,
			PermissionCode: code,
		}
		if err := tx.Create(&rolePermission).Error; err != nil {
			return fmt.Errorf("创建角色权限关联失败: %w", err)
		}
	}
	return nil
}

// convertToRoleResponse 转换为角色响应
func (s *AdminRoleService) convertToRoleResponse(role models.AdminRole) types.RoleResponse {
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
	var userCount int64
	s.db.Model(&models.AdminUserRole{}).Where("role_id = ?", role.ID).Count(&userCount)

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
func (s *AdminRoleService) getStatusText(status int) string {
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
func (s *AdminRoleService) getLevelText(level int) string {
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
