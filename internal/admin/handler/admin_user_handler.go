package handler

import (
	"net/http"
	"strconv"

	"bico-admin/internal/admin/service"
	"bico-admin/internal/pkg/pagination"
	"github.com/gin-gonic/gin"
)

// AdminUserHandler 用户管理处理器
type AdminUserHandler struct {
	userService *service.AdminUserService
}

// NewAdminUserHandler 创建用户管理处理器
func NewAdminUserHandler(userService *service.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{
		userService: userService,
	}
}

// List 获取用户列表
func (h *AdminUserHandler) List(c *gin.Context) {
	var req service.ListRequest
	req.Pagination = *pagination.FromContext(c)
	req.Username = c.Query("username")
	req.Name = c.Query("name")
	
	if enabledStr := c.Query("enabled"); enabledStr != "" {
		enabled := enabledStr == "true"
		req.Enabled = &enabled
	}
	
	resp, err := h.userService.List(&req)
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

// Get 获取用户详情
func (h *AdminUserHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "无效的用户ID",
		})
		return
	}
	
	user, err := h.userService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}
	
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": user,
	})
}

// Create 创建用户
func (h *AdminUserHandler) Create(c *gin.Context) {
	var req service.CreateAdminUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}
	
	user, err := h.userService.Create(&req)
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
		"data": user,
	})
}

// Update 更新用户
func (h *AdminUserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "无效的用户ID",
		})
		return
	}
	
	var req service.UpdateAdminUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}
	
	user, err := h.userService.Update(uint(id), &req)
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
		"data": user,
	})
}

// Delete 删除用户
func (h *AdminUserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "无效的用户ID",
		})
		return
	}
	
	if err := h.userService.Delete(uint(id)); err != nil {
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
