package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"bico-admin/core/config"
	"bico-admin/core/database"
	"bico-admin/core/model"
	"bico-admin/pkg/response"
)

// PermissionHandler 权限管理处理器 - 业务逻辑直接在handler中实现
type PermissionHandler struct {
	db      *gorm.DB
	config  *config.Config
	permOps *database.Operations[model.Permission]
}

// NewPermissionHandler 创建权限管理处理器实例
func NewPermissionHandler(db *gorm.DB) *PermissionHandler {
	cfg := config.New()
	return &PermissionHandler{
		db:      db,
		config:  cfg,
		permOps: database.NewOperations[model.Permission](db),
	}
}

// GetPermissions 获取权限列表
func (h *PermissionHandler) GetPermissions(c *fiber.Ctx) error {
	// 解析查询参数
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	search := c.Query("search", "")
	permType := c.Query("type", "")

	params := database.PaginationParams{
		Page:         page,
		PageSize:     pageSize,
		Search:       search,
		SearchFields: []string{"name", "code", "resource", "action"},
		Preloads:     []string{"Parent", "Children"},
	}

	// 如果指定了权限类型，添加过滤条件
	if permType != "" {
		params.Filters = map[string]interface{}{
			"type": permType,
		}
	}

	result, err := h.permOps.List(params)
	if err != nil {
		return response.InternalServerError(c, "Failed to get permissions")
	}

	// 使用Ant Design Pro标准的分页响应格式
	return response.Pagination(c, result.Data, result.Total, result.Page, result.PageSize)
}

// GetPermissionTree 获取权限树结构
func (h *PermissionHandler) GetPermissionTree(c *fiber.Ctx) error {
	// 获取所有权限
	var permissions []model.Permission
	if err := h.db.Preload("Children").Where("parent_id IS NULL").Order("sort ASC").Find(&permissions).Error; err != nil {
		return response.InternalServerError(c, "Failed to get permission tree")
	}

	return response.Success(c, permissions)
}

// CreatePermission 创建权限
func (h *PermissionHandler) CreatePermission(c *fiber.Ctx) error {
	var req model.PermissionCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 检查权限名称和代码是否已存在
	if existingPerm, _ := h.permOps.GetByCondition("name = ?", req.Name); existingPerm != nil {
		return response.BadRequest(c, "Permission name already exists")
	}
	if existingPerm, _ := h.permOps.GetByCondition("code = ?", req.Code); existingPerm != nil {
		return response.BadRequest(c, "Permission code already exists")
	}

	// 创建权限
	permission := model.Permission{
		Name:        req.Name,
		Code:        req.Code,
		Type:        req.Type,
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
		ParentID:    req.ParentID,
		Sort:        req.Sort,
		Status:      model.PermissionStatusActive,
	}

	// 保存权限
	if err := h.permOps.Create(&permission); err != nil {
		return response.InternalServerError(c, "Failed to create permission")
	}

	// 获取完整的权限信息（包含父级和子级）
	var fullPermission model.Permission
	if err := h.db.Preload("Parent").Preload("Children").First(&fullPermission, permission.ID).Error; err != nil {
		// 如果获取失败，返回基本信息
		return response.Success(c, permission)
	}

	return response.Success(c, fullPermission)
}

// GetPermission 获取单个权限
func (h *PermissionHandler) GetPermission(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid permission ID")
	}

	// 获取权限及其父级和子级信息
	var permission model.Permission
	if err := h.db.Preload("Parent").Preload("Children").Where("id = ?", uint(id)).First(&permission).Error; err != nil {
		return response.NotFound(c, "Permission not found")
	}

	return response.Success(c, permission)
}

// UpdatePermission 更新权限
func (h *PermissionHandler) UpdatePermission(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid permission ID")
	}

	var req model.PermissionUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 检查权限是否存在
	permission, err := h.permOps.Get(uint(id))
	if err != nil {
		return response.NotFound(c, "Permission not found")
	}

	// 检查名称和代码是否重复（排除当前权限）
	if req.Name != "" && req.Name != permission.Name {
		if existingPerm, _ := h.permOps.GetByCondition("name = ? AND id != ?", req.Name, uint(id)); existingPerm != nil {
			return response.BadRequest(c, "Permission name already exists")
		}
	}
	if req.Code != "" && req.Code != permission.Code {
		if existingPerm, _ := h.permOps.GetByCondition("code = ? AND id != ?", req.Code, uint(id)); existingPerm != nil {
			return response.BadRequest(c, "Permission code already exists")
		}
	}

	// 构建更新数据
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}
	if req.Type != 0 {
		updates["type"] = req.Type
	}
	if req.Resource != "" {
		updates["resource"] = req.Resource
	}
	if req.Action != "" {
		updates["action"] = req.Action
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.ParentID != nil {
		updates["parent_id"] = req.ParentID
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}

	// 更新权限信息
	if len(updates) > 0 {
		if err := h.permOps.UpdateFields(uint(id), updates); err != nil {
			return response.InternalServerError(c, "Failed to update permission")
		}
	}

	// 获取更新后的权限信息（包含父级和子级）
	var updatedPermission model.Permission
	if err := h.db.Preload("Parent").Preload("Children").First(&updatedPermission, uint(id)).Error; err != nil {
		return response.InternalServerError(c, "Failed to get updated permission")
	}

	return response.Success(c, updatedPermission)
}

// DeletePermission 删除权限
func (h *PermissionHandler) DeletePermission(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid permission ID")
	}

	// 检查权限是否有子权限
	var childCount int64
	if err := h.db.Model(&model.Permission{}).Where("parent_id = ?", uint(id)).Count(&childCount).Error; err != nil {
		return response.InternalServerError(c, "Failed to check child permissions")
	}

	if childCount > 0 {
		return response.BadRequest(c, "Cannot delete permission that has child permissions")
	}

	// 检查权限是否被角色使用
	var roleCount int64
	if err := h.db.Table("role_permissions").Where("permission_id = ?", uint(id)).Count(&roleCount).Error; err != nil {
		return response.InternalServerError(c, "Failed to check permission usage")
	}

	if roleCount > 0 {
		return response.BadRequest(c, "Cannot delete permission that is assigned to roles")
	}

	// 删除权限
	if err := h.permOps.Delete(uint(id)); err != nil {
		return response.InternalServerError(c, "Failed to delete permission")
	}

	return response.SuccessWithMessage(c, "Permission deleted successfully", nil)
}

// BatchDeletePermissions 批量删除权限
func (h *PermissionHandler) BatchDeletePermissions(c *fiber.Ctx) error {
	var req struct {
		IDs []uint `json:"ids" validate:"required,min=1"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if len(req.IDs) == 0 {
		return response.BadRequest(c, "No permission IDs provided")
	}

	// 检查权限是否有子权限或被角色使用
	for _, id := range req.IDs {
		var childCount int64
		if err := h.db.Model(&model.Permission{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
			return response.InternalServerError(c, "Failed to check child permissions")
		}
		if childCount > 0 {
			return response.BadRequest(c, "Cannot delete permissions that have child permissions")
		}

		var roleCount int64
		if err := h.db.Table("role_permissions").Where("permission_id = ?", id).Count(&roleCount).Error; err != nil {
			return response.InternalServerError(c, "Failed to check permission usage")
		}
		if roleCount > 0 {
			return response.BadRequest(c, "Cannot delete permissions that are assigned to roles")
		}
	}

	// 批量删除权限
	for _, id := range req.IDs {
		if err := h.permOps.Delete(id); err != nil {
			return response.InternalServerError(c, "Failed to delete some permissions")
		}
	}

	return response.SuccessWithMessage(c, "Permissions deleted successfully", nil)
}

// UpdatePermissionStatus 更新权限状态
func (h *PermissionHandler) UpdatePermissionStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid permission ID")
	}

	var req struct {
		Status model.PermissionStatus `json:"status" validate:"required,oneof=0 1"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 检查权限是否存在
	_, err = h.permOps.Get(uint(id))
	if err != nil {
		return response.NotFound(c, "Permission not found")
	}

	// 更新权限状态
	updates := map[string]interface{}{
		"status": req.Status,
	}

	if err := h.permOps.UpdateFields(uint(id), updates); err != nil {
		return response.InternalServerError(c, "Failed to update permission status")
	}

	// 获取更新后的权限信息（包含父级和子级）
	var updatedPermission model.Permission
	if err := h.db.Preload("Parent").Preload("Children").First(&updatedPermission, uint(id)).Error; err != nil {
		return response.InternalServerError(c, "Failed to get updated permission")
	}

	return response.Success(c, updatedPermission)
}
