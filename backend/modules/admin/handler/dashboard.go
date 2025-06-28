package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"bico-admin/core/config"
	"bico-admin/core/database"
	"bico-admin/core/model"
	"bico-admin/pkg/response"
)

// DashboardHandler 仪表板处理器 - 业务逻辑直接在handler中实现
type DashboardHandler struct {
	db      *gorm.DB
	config  *config.Config
	userOps *database.Operations[model.User]
}

// NewDashboardHandler 创建仪表板处理器实例
func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	cfg := config.New()
	return &DashboardHandler{
		db:      db,
		config:  cfg,
		userOps: database.NewOperations[model.User](db),
	}
}

// GetDashboard 获取仪表板数据
func (h *DashboardHandler) GetDashboard(c *fiber.Ctx) error {
	// 获取用户总数
	userCount, err := h.userOps.Count()
	if err != nil {
		return response.InternalServerError(c, "Failed to get user count")
	}

	// 获取活跃用户数
	activeUserCount, err := h.userOps.CountWithCondition("status = ?", 1)
	if err != nil {
		return response.InternalServerError(c, "Failed to get active user count")
	}

	data := fiber.Map{
		"user_count":        userCount,
		"active_user_count": activeUserCount,
		"total_visits":      0, // TODO: 实现访问统计
		"today_visits":      0, // TODO: 实现今日访问统计
	}

	return response.Success(c, data)
}

// GetStats 获取统计数据
func (h *DashboardHandler) GetStats(c *fiber.Ctx) error {
	// 获取各种统计数据
	userCount, _ := h.userOps.Count()
	activeUserCount, _ := h.userOps.CountWithCondition("status = ?", 1)

	stats := fiber.Map{
		"users": fiber.Map{
			"total":  userCount,
			"active": activeUserCount,
		},
		"system": fiber.Map{
			"uptime": "24h",   // TODO: 实现系统运行时间统计
			"memory": "512MB", // TODO: 实现内存使用统计
		},
	}

	return response.Success(c, stats)
}

// GetUsers 获取用户列表
func (h *DashboardHandler) GetUsers(c *fiber.Ctx) error {
	// 解析查询参数
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	search := c.Query("search", "")

	params := database.PaginationParams{
		Page:         page,
		PageSize:     pageSize,
		Search:       search,
		SearchFields: []string{"username", "email", "nickname"},
	}

	result, err := h.userOps.List(params)
	if err != nil {
		return response.InternalServerError(c, "Failed to get users")
	}

	return response.Success(c, result)
}

// CreateUser 创建用户
func (h *DashboardHandler) CreateUser(c *fiber.Ctx) error {
	var req model.UserCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	user, err := h.createUser(req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, user)
}

// GetUser 获取单个用户
func (h *DashboardHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	user, err := h.getUserWithRoles(uint(id))
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	userResponse := user.ToResponse()
	return response.Success(c, userResponse)
}

// UpdateUser 更新用户
func (h *DashboardHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	var req model.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	user, err := h.updateUserProfile(uint(id), req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, user)
}

// DeleteUser 删除用户
func (h *DashboardHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	if err := h.userOps.Delete(uint(id)); err != nil {
		return response.InternalServerError(c, "Failed to delete user")
	}

	return response.SuccessWithMessage(c, "User deleted successfully", nil)
}

// GetRoles 获取角色列表
func (h *DashboardHandler) GetRoles(c *fiber.Ctx) error {
	// TODO: 实现角色列表获取
	return response.Success(c, fiber.Map{
		"message": "Get roles - to be implemented",
	})
}

// CreateRole 创建角色
func (h *DashboardHandler) CreateRole(c *fiber.Ctx) error {
	// TODO: 实现角色创建
	return response.Success(c, fiber.Map{
		"message": "Create role - to be implemented",
	})
}

// GetRole 获取单个角色
func (h *DashboardHandler) GetRole(c *fiber.Ctx) error {
	// TODO: 实现单个角色获取
	return response.Success(c, fiber.Map{
		"message": "Get role - to be implemented",
	})
}

// UpdateRole 更新角色
func (h *DashboardHandler) UpdateRole(c *fiber.Ctx) error {
	// TODO: 实现角色更新
	return response.Success(c, fiber.Map{
		"message": "Update role - to be implemented",
	})
}

// DeleteRole 删除角色
func (h *DashboardHandler) DeleteRole(c *fiber.Ctx) error {
	// TODO: 实现角色删除
	return response.Success(c, fiber.Map{
		"message": "Delete role - to be implemented",
	})
}

// GetPermissions 获取权限列表
func (h *DashboardHandler) GetPermissions(c *fiber.Ctx) error {
	// TODO: 实现权限列表获取
	return response.Success(c, fiber.Map{
		"message": "Get permissions - to be implemented",
	})
}

// CreatePermission 创建权限
func (h *DashboardHandler) CreatePermission(c *fiber.Ctx) error {
	// TODO: 实现权限创建
	return response.Success(c, fiber.Map{
		"message": "Create permission - to be implemented",
	})
}

// GetPermission 获取单个权限
func (h *DashboardHandler) GetPermission(c *fiber.Ctx) error {
	// TODO: 实现单个权限获取
	return response.Success(c, fiber.Map{
		"message": "Get permission - to be implemented",
	})
}

// UpdatePermission 更新权限
func (h *DashboardHandler) UpdatePermission(c *fiber.Ctx) error {
	// TODO: 实现权限更新
	return response.Success(c, fiber.Map{
		"message": "Update permission - to be implemented",
	})
}

// DeletePermission 删除权限
func (h *DashboardHandler) DeletePermission(c *fiber.Ctx) error {
	// TODO: 实现权限删除
	return response.Success(c, fiber.Map{
		"message": "Delete permission - to be implemented",
	})
}

// 业务逻辑方法 - 直接在handler中实现

// createUser 创建用户
func (h *DashboardHandler) createUser(req model.UserCreateRequest) (*model.UserResponse, error) {
	// 检查用户名是否已存在
	if existingUser, _ := h.userOps.GetByCondition("username = ?", req.Username); existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	if existingUser, _ := h.userOps.GetByCondition("email = ?", req.Email); existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Status:   model.UserStatusActive,
	}

	// 加密密码
	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	// 保存用户
	if err := h.userOps.Create(&user); err != nil {
		return nil, err
	}

	userResponse := user.ToResponse()
	return &userResponse, nil
}

// getUserWithRoles 获取用户及其角色信息
func (h *DashboardHandler) getUserWithRoles(userID uint) (*model.User, error) {
	var user model.User
	if err := h.db.Preload("Roles").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// updateUserProfile 更新用户资料
func (h *DashboardHandler) updateUserProfile(userID uint, req model.UserUpdateRequest) (*model.UserResponse, error) {
	// 检查用户是否存在
	user, err := h.userOps.Get(userID)
	if err != nil {
		return nil, err
	}

	// 如果更新用户名，检查是否重复
	if req.Username != "" && req.Username != user.Username {
		if existingUser, _ := h.userOps.GetByCondition("username = ?", req.Username); existingUser != nil {
			return nil, errors.New("username already exists")
		}
	}

	// 如果更新邮箱，检查是否重复
	if req.Email != "" && req.Email != user.Email {
		if existingUser, _ := h.userOps.GetByCondition("email = ?", req.Email); existingUser != nil {
			return nil, errors.New("email already exists")
		}
	}

	// 构建更新数据
	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	// 更新用户信息
	if err := h.userOps.UpdateFields(userID, updates); err != nil {
		return nil, err
	}

	// 获取更新后的用户信息
	updatedUser, err := h.getUserWithRoles(userID)
	if err != nil {
		return nil, err
	}

	userResponse := updatedUser.ToResponse()
	return &userResponse, nil
}
