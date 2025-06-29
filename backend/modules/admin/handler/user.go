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

// UserHandler 用户管理处理器 - 业务逻辑直接在handler中实现
type UserHandler struct {
	db            *gorm.DB
	config        *config.Config
	userOps       *database.Operations[model.User]
	userValidator *database.UserValidationHelper
}

// NewUserHandler 创建用户管理处理器实例
func NewUserHandler(db *gorm.DB) *UserHandler {
	cfg := config.New()
	return &UserHandler{
		db:            db,
		config:        cfg,
		userOps:       database.NewOperations[model.User](db),
		userValidator: database.NewUserValidationHelper(db),
	}
}

// GetUsers 获取用户列表
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	// 解析查询参数
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	search := c.Query("search", "")

	params := database.PaginationParams{
		Page:         page,
		PageSize:     pageSize,
		Search:       search,
		SearchFields: []string{"username", "email", "nickname"},
		Preloads:     []string{"Roles"},
	}

	result, err := h.userOps.List(params)
	if err != nil {
		return response.InternalServerError(c, "Failed to get users")
	}

	// 使用Ant Design Pro标准的分页响应格式
	return response.Pagination(c, result.Data, result.Total, result.Page, result.PageSize)
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
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
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	// 获取用户及其角色信息
	var user model.User
	if err := h.db.Preload("Roles").Where("id = ?", uint(id)).First(&user).Error; err != nil {
		return response.NotFound(c, "User not found")
	}

	userResponse := user.ToResponse()
	return response.Success(c, userResponse)
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
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
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	if err := h.userOps.Delete(uint(id)); err != nil {
		return response.InternalServerError(c, "Failed to delete user")
	}

	return response.SuccessWithMessage(c, "User deleted successfully", nil)
}

// BatchDeleteUsers 批量删除用户
func (h *UserHandler) BatchDeleteUsers(c *fiber.Ctx) error {
	var req struct {
		IDs []uint `json:"ids" validate:"required,min=1"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if len(req.IDs) == 0 {
		return response.BadRequest(c, "No user IDs provided")
	}

	// 批量删除用户
	for _, id := range req.IDs {
		if err := h.userOps.Delete(id); err != nil {
			return response.InternalServerError(c, "Failed to delete some users")
		}
	}

	return response.SuccessWithMessage(c, "Users deleted successfully", nil)
}

// UpdateUserStatus 更新用户状态
func (h *UserHandler) UpdateUserStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	var req struct {
		Status model.UserStatus `json:"status" validate:"required,oneof=0 1 2"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 检查用户是否存在
	_, err = h.userOps.Get(uint(id))
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	// 更新用户状态
	updates := map[string]interface{}{
		"status": req.Status,
	}

	if err := h.userOps.UpdateFields(uint(id), updates); err != nil {
		return response.InternalServerError(c, "Failed to update user status")
	}

	// 获取更新后的用户信息（包含角色）
	var updatedUser model.User
	if err := h.db.Preload("Roles").First(&updatedUser, uint(id)).Error; err != nil {
		return response.InternalServerError(c, "Failed to get updated user")
	}

	userResponse := updatedUser.ToResponse()
	return response.Success(c, userResponse)
}

// ResetUserPassword 重置用户密码
func (h *UserHandler) ResetUserPassword(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	var req struct {
		NewPassword string `json:"new_password" validate:"required,min=6"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 检查用户是否存在
	user, err := h.userOps.Get(uint(id))
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	// 加密新密码
	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		return response.InternalServerError(c, "Failed to hash password")
	}

	// 更新密码
	updates := map[string]interface{}{
		"password": user.Password,
	}

	if err := h.userOps.UpdateFields(uint(id), updates); err != nil {
		return response.InternalServerError(c, "Failed to reset password")
	}

	return response.SuccessWithMessage(c, "Password reset successfully", nil)
}

// ChangeUserPassword 修改用户密码（需要旧密码验证）
func (h *UserHandler) ChangeUserPassword(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	var req model.UserChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 获取用户信息
	user, err := h.userOps.Get(uint(id))
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	// 验证旧密码
	if !user.CheckPassword(req.OldPassword) {
		return response.BadRequest(c, "Old password is incorrect")
	}

	// 加密新密码
	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		return response.InternalServerError(c, "Failed to hash password")
	}

	// 更新密码
	updates := map[string]interface{}{
		"password": user.Password,
	}

	if err := h.userOps.UpdateFields(uint(id), updates); err != nil {
		return response.InternalServerError(c, "Failed to change password")
	}

	return response.SuccessWithMessage(c, "Password changed successfully", nil)
}

// 业务逻辑方法 - 直接在handler中实现

// createUser 创建用户
func (h *UserHandler) createUser(req model.UserCreateRequest) (*model.UserResponse, error) {
	// 检查用户名和邮箱是否已存在
	if err := h.userValidator.CheckUserUniqueFields(req.Username, req.Email); err != nil {
		return nil, err
	}

	// 设置默认角色
	role := req.Role
	if role == "" {
		role = "user" // 默认角色
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Role:     role,
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

	// 如果有角色ID，关联角色
	if len(req.RoleIDs) > 0 {
		// TODO: 实现角色分配逻辑，当角色管理功能完成后
		// 暂时跳过角色关联，因为角色管理还未完全实现
	}

	// 获取完整的用户信息（包含角色）
	var fullUser model.User
	if err := h.db.Preload("Roles").First(&fullUser, user.ID).Error; err != nil {
		// 如果获取失败，返回基本信息
		userResponse := user.ToResponse()
		return &userResponse, nil
	}

	userResponse := fullUser.ToResponse()
	return &userResponse, nil
}

// updateUserProfile 更新用户资料
func (h *UserHandler) updateUserProfile(userID uint, req model.UserUpdateRequest) (*model.UserResponse, error) {
	// 检查用户是否存在
	user, err := h.userOps.Get(userID)
	if err != nil {
		return nil, err
	}

	// 检查用户名和邮箱是否重复（排除当前用户）
	fieldsToCheck := make(map[string]string)
	if req.Username != "" && req.Username != user.Username {
		fieldsToCheck["username"] = req.Username
	}
	if req.Email != "" && req.Email != user.Email {
		fieldsToCheck["email"] = req.Email
	}

	// 如果有需要检查的字段，进行唯一性验证
	if len(fieldsToCheck) > 0 {
		if err := h.userValidator.CheckUserUniqueFields(fieldsToCheck["username"], fieldsToCheck["email"], userID); err != nil {
			return nil, err
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
	if req.Status != 0 {
		updates["status"] = req.Status
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}

	// 更新用户信息
	if len(updates) > 0 {
		if err := h.userOps.UpdateFields(userID, updates); err != nil {
			return nil, err
		}
	}

	// 如果有角色ID，更新角色关联
	if len(req.RoleIDs) > 0 {
		// TODO: 实现角色更新逻辑，当角色管理功能完成后
		// 暂时跳过角色更新，因为角色管理还未完全实现
	}

	// 获取更新后的用户信息（包含角色）
	var updatedUser model.User
	if err := h.db.Preload("Roles").First(&updatedUser, userID).Error; err != nil {
		return nil, err
	}

	userResponse := updatedUser.ToResponse()
	return &userResponse, nil
}
