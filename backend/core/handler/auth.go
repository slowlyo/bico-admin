package handler

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"bico-admin/core/cache"
	"bico-admin/core/config"
	"bico-admin/core/database"
	"bico-admin/core/middleware"
	"bico-admin/core/model"
	"bico-admin/pkg/response"
	"bico-admin/pkg/validator"
)

// AuthHandler 认证处理器 - 业务逻辑直接在handler中实现
type AuthHandler struct {
	db        *gorm.DB
	config    *config.Config
	userOps   *database.Operations[model.User]
	cache     *cache.Cache
	jwtSecret string
	jwtExpire time.Duration
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	jwtExpire, err := time.ParseDuration(cfg.JWT.Expire)
	if err != nil {
		jwtExpire = 24 * time.Hour // 默认24小时
	}

	return &AuthHandler{
		db:        db,
		config:    cfg,
		userOps:   database.NewOperations[model.User](db),
		cache:     cache.GetCache(),
		jwtSecret: cfg.JWT.Secret,
		jwtExpire: jwtExpire,
	}
}

// Login 用户登录
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 验证请求参数
	if errors := validator.Validate(req); len(errors) > 0 {
		return response.ValidationError(c, errors)
	}

	// 检查登录尝试次数限制
	clientIP := c.IP()
	loginIdentifier := req.Username + ":" + clientIP
	attemptCount, err := h.cache.GetLoginAttemptCount(loginIdentifier)
	if err == nil && attemptCount >= 5 {
		return response.Unauthorized(c, "Too many login attempts, please try again later")
	}

	// 根据用户名或邮箱查找用户
	user, err := h.userOps.GetByCondition("username = ? OR email = ?", req.Username, req.Username)
	if err != nil {
		// 增加登录失败次数
		h.cache.IncrementLoginAttempt(loginIdentifier, 15*time.Minute)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Unauthorized(c, "invalid credentials")
		}
		return response.Unauthorized(c, "login failed")
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		// 增加登录失败次数
		h.cache.IncrementLoginAttempt(loginIdentifier, 15*time.Minute)
		return response.Unauthorized(c, "invalid credentials")
	}

	// 检查用户状态
	if user.Status != model.UserStatusActive {
		// 增加登录失败次数
		h.cache.IncrementLoginAttempt(loginIdentifier, 15*time.Minute)
		return response.Unauthorized(c, "user account is not active")
	}

	// 登录成功，重置登录尝试次数
	h.cache.ResetLoginAttempt(loginIdentifier)

	// 更新最后登录时间和IP
	now := time.Now()
	loginIP := c.IP()
	err = h.userOps.UpdateFields(user.ID, map[string]interface{}{
		"last_login_at": &now,
		"last_login_ip": loginIP,
	})
	if err != nil {
		// 记录日志但不影响登录流程
		// log.Printf("Failed to update login info: %v", err)
	}

	// 生成JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, h.jwtSecret, h.jwtExpire)
	if err != nil {
		return response.Unauthorized(c, "failed to generate token")
	}

	// 获取用户完整信息（包含角色）
	var userWithRoles model.User
	err = h.db.Preload("Roles").First(&userWithRoles, user.ID).Error
	if err != nil {
		// 如果获取角色失败，使用基本用户信息
		userWithRoles = *user
	}

	userResponse := userWithRoles.ToResponse()

	// 设置用户会话缓存
	sessionData := map[string]interface{}{
		"user_id":    user.ID,
		"username":   user.Username,
		"login_time": time.Now(),
		"login_ip":   loginIP,
	}
	h.cache.SetUserSession(user.ID, sessionData, h.jwtExpire)

	return response.Success(c, fiber.Map{
		"user":  userResponse,
		"token": token,
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	// 获取当前token并添加到黑名单
	authHeader := c.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 将token添加到黑名单，设置过期时间为JWT的剩余有效期
		err := h.cache.AddTokenToBlacklist(tokenString, h.jwtExpire)
		if err != nil {
			// 记录错误但不影响登出流程
			// log.Printf("Failed to add token to blacklist: %v", err)
		}
	}

	// 删除用户会话缓存
	err := h.cache.DeleteUserSession(userID)
	if err != nil {
		// 记录错误但不影响登出流程
		// log.Printf("Failed to delete user session: %v", err)
	}

	// 记录登出日志
	// log.Printf("User %d logged out from IP: %s", userID, c.IP())

	return response.SuccessWithMessage(c, "Logout successful", nil)
}

// GetProfile 获取用户资料
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	// 获取用户信息（包含角色）
	var user model.User
	err := h.db.Preload("Roles").First(&user, userID).Error
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	userResponse := user.ToResponse()
	return response.Success(c, userResponse)
}

// GetUserPermissions 获取当前用户权限
func (h *AuthHandler) GetUserPermissions(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	// 检查用户是否为超级管理员
	isSuperAdmin, err := h.isSuperAdmin(userID)
	if err != nil {
		return response.InternalServerError(c, "Failed to check user role")
	}

	// 超级管理员拥有所有权限
	if isSuperAdmin {
		// 从代码中获取所有权限常量
		allPermissionCodes := []string{
			// 系统管理权限
			"system:view", "system:manage",
			// 用户管理权限
			"user:view", "user:create", "user:update", "user:delete", "user:manage_status", "user:reset_password",
			// 角色管理权限
			"role:view", "role:create", "role:update", "role:delete", "role:assign_permissions",
		}
		return response.Success(c, allPermissionCodes)
	}

	// 普通用户通过角色获取权限
	var permissionCodes []string
	err = h.db.Table("role_permissions rp").
		Select("DISTINCT rp.permission_code").
		Joins("JOIN user_roles ur ON rp.role_id = ur.role_id").
		Where("ur.user_id = ?", userID).
		Pluck("permission_code", &permissionCodes).Error

	if err != nil {
		return response.InternalServerError(c, "Failed to get user permissions")
	}

	return response.Success(c, permissionCodes)
}

// isSuperAdmin 检查用户是否为超级管理员
func (h *AuthHandler) isSuperAdmin(userID uint) (bool, error) {
	var count int64
	err := h.db.Table("users u").
		Joins("JOIN user_roles ur ON u.id = ur.user_id").
		Joins("JOIN roles r ON ur.role_id = r.id").
		Where("u.id = ? AND r.code = ? AND r.status = ?",
			userID, "super_admin", model.RoleStatusActive).
		Count(&count).Error

	if err != nil {
		// 如果查询失败，尝试通过用户表的role字段检查（兼容旧版本）
		var user model.User
		if err := h.db.First(&user, userID).Error; err != nil {
			return false, err
		}
		return user.Role == "super_admin", nil
	}

	return count > 0, nil
}

// UpdateProfile 更新用户资料
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	var req model.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 验证请求参数
	if errors := validator.Validate(req); len(errors) > 0 {
		return response.ValidationError(c, errors)
	}

	// 获取当前用户
	user, err := h.userOps.GetById(userID)
	if err != nil {
		return response.BadRequest(c, "user not found")
	}

	// 检查用户名是否已存在（如果要更新用户名）
	if req.Username != "" && req.Username != user.Username {
		existingUser, err := h.userOps.GetByCondition("username = ? AND id != ?", req.Username, userID)
		if err == nil && existingUser.ID != 0 {
			return response.BadRequest(c, "username already exists")
		}
	}

	// 检查邮箱是否已存在（如果要更新邮箱）
	if req.Email != "" && req.Email != user.Email {
		existingUser, err := h.userOps.GetByCondition("email = ? AND id != ?", req.Email, userID)
		if err == nil && existingUser.ID != 0 {
			return response.BadRequest(c, "email already exists")
		}
	}

	// 准备更新字段
	updateFields := make(map[string]interface{})
	if req.Username != "" {
		updateFields["username"] = req.Username
	}
	if req.Email != "" {
		updateFields["email"] = req.Email
	}
	if req.Nickname != "" {
		updateFields["nickname"] = req.Nickname
	}
	if req.Phone != "" {
		updateFields["phone"] = req.Phone
	}

	// 更新用户信息
	if len(updateFields) > 0 {
		err = h.userOps.UpdateFields(userID, updateFields)
		if err != nil {
			return response.BadRequest(c, "failed to update user profile")
		}
	}

	// 处理角色更新（如果提供了角色ID）
	if len(req.RoleIDs) > 0 {
		// 获取角色
		var roles []model.Role
		err = h.db.Where("id IN ?", req.RoleIDs).Find(&roles).Error
		if err != nil {
			return response.BadRequest(c, "invalid role IDs")
		}

		// 更新用户角色关联
		err = h.db.Model(user).Association("Roles").Replace(roles)
		if err != nil {
			return response.BadRequest(c, "failed to update user roles")
		}
	}

	// 获取更新后的用户信息（包含角色）
	var updatedUser model.User
	err = h.db.Preload("Roles").First(&updatedUser, userID).Error
	if err != nil {
		return response.BadRequest(c, "failed to get updated user info")
	}

	userResponse := updatedUser.ToResponse()
	return response.SuccessWithMessage(c, "Profile updated successfully", userResponse)
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "User not authenticated")
	}

	var req model.UserChangePasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// 验证请求参数
	if errors := validator.Validate(req); len(errors) > 0 {
		return response.ValidationError(c, errors)
	}

	// 获取当前用户
	user, err := h.userOps.GetById(userID)
	if err != nil {
		return response.BadRequest(c, "user not found")
	}

	// 验证旧密码
	if !user.CheckPassword(req.OldPassword) {
		return response.BadRequest(c, "old password is incorrect")
	}

	// 创建临时用户对象来加密新密码
	tempUser := &model.User{Password: req.NewPassword}
	if err := tempUser.HashPassword(); err != nil {
		return response.BadRequest(c, "failed to encrypt password")
	}

	// 更新密码
	err = h.userOps.UpdateFields(userID, map[string]interface{}{
		"password": tempUser.Password,
	})
	if err != nil {
		return response.BadRequest(c, "failed to update password")
	}

	return response.SuccessWithMessage(c, "Password changed successfully", nil)
}
