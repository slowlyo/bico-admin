package handler

import (
	"net/http"
	"strconv"

	"bico-admin/internal/admin/consts"
	"bico-admin/internal/admin/service"
	"bico-admin/internal/pkg/pagination"
	"github.com/gin-gonic/gin"
)

// AdminRoleHandler 角色管理处理器
type AdminRoleHandler struct {
	roleService *service.AdminRoleService
}

// NewAdminRoleHandler 创建角色管理处理器
func NewAdminRoleHandler(roleService *service.AdminRoleService) *AdminRoleHandler {
	return &AdminRoleHandler{
		roleService: roleService,
	}
}

// List 获取角色列表
func (h *AdminRoleHandler) List(c *gin.Context) {
	var req service.RoleListRequest
	req.Pagination = *pagination.FromContext(c)
	req.Name = c.Query("name")
	req.Code = c.Query("code")
	
	if enabledStr := c.Query("enabled"); enabledStr != "" {
		enabled := enabledStr == "true"
		req.Enabled = &enabled
	}
	
	resp, err := h.roleService.ListRoles(&req)
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": resp.Data,
		"total": resp.Total,
	})
}

// Get 获取角色详情
func (h *AdminRoleHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "无效的角色ID",
		})
		return
	}
	
	role, err := h.roleService.GetRole(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 404,
			"msg":  "角色不存在",
		})
		return
	}
	
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": role,
	})
}

// Create 创建角色
func (h *AdminRoleHandler) Create(c *gin.Context) {
	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}
	
	role, err := h.roleService.CreateRole(&req)
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "创建成功",
		"data": role,
	})
}

// Update 更新角色
func (h *AdminRoleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "无效的角色ID",
		})
		return
	}
	
	var req service.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}
	
	role, err := h.roleService.UpdateRole(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "更新成功",
		"data": role,
	})
}

// Delete 删除角色
func (h *AdminRoleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "无效的角色ID",
		})
		return
	}
	
	if err := h.roleService.DeleteRole(uint(id)); err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "删除成功",
	})
}

// UpdatePermissions 更新角色权限
func (h *AdminRoleHandler) UpdatePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "无效的角色ID",
		})
		return
	}
	
	var req service.UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}
	
	if err := h.roleService.UpdateRolePermissions(uint(id), &req); err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "权限配置成功",
	})
}

// GetAllPermissions 获取所有权限树
func (h *AdminRoleHandler) GetAllPermissions(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": consts.AllPermissions,
	})
}

// GetAll 获取所有角色（用于下拉选择）
func (h *AdminRoleHandler) GetAll(c *gin.Context) {
	roles, err := h.roleService.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": roles,
	})
}
