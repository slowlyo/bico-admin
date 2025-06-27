package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"bico-admin/core/middleware"
	"bico-admin/core/model"
	"bico-admin/core/repository"
)

// AuthService 认证服务接口
type AuthService interface {
	Login(req model.UserLoginRequest) (*model.UserResponse, string, error)
	Register(req model.UserCreateRequest) (*model.UserResponse, error)
	GetProfile(userID uint) (*model.UserResponse, error)
	UpdateProfile(userID uint, req model.UserUpdateRequest) (*model.UserResponse, error)
	ChangePassword(userID uint, oldPassword, newPassword string) error
}

// AuthServiceImpl 认证服务实现
type AuthServiceImpl struct {
	userRepo   repository.UserRepository
	jwtSecret  string
	jwtExpire  time.Duration
}

// NewAuthService 创建认证服务实例
func NewAuthService(db *gorm.DB, jwtSecret string, jwtExpire time.Duration) AuthService {
	return &AuthServiceImpl{
		userRepo:  repository.NewUserRepository(db),
		jwtSecret: jwtSecret,
		jwtExpire: jwtExpire,
	}
}

// Login 用户登录
func (s *AuthServiceImpl) Login(req model.UserLoginRequest) (*model.UserResponse, string, error) {
	// 根据用户名或邮箱查找用户
	var user *model.User
	var err error
	
	if user, err = s.userRepo.GetByUsername(req.Username); err != nil {
		if user, err = s.userRepo.GetByEmail(req.Username); err != nil {
			return nil, "", errors.New("invalid credentials")
		}
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, "", errors.New("invalid credentials")
	}

	// 检查用户状态
	if user.Status != model.UserStatusActive {
		return nil, "", errors.New("user account is not active")
	}

	// 生成JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, "", err
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user.ID, user)

	// 获取用户角色信息
	userWithRoles, _ := s.userRepo.GetWithRoles(user.ID)
	if userWithRoles != nil {
		user = userWithRoles
	}

	response := user.ToResponse()
	return &response, token, nil
}

// Register 用户注册
func (s *AuthServiceImpl) Register(req model.UserCreateRequest) (*model.UserResponse, error) {
	// 检查用户名是否已存在
	if existingUser, _ := s.userRepo.GetByUsername(req.Username); existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	if existingUser, _ := s.userRepo.GetByEmail(req.Email); existingUser != nil {
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
	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// GetProfile 获取用户资料
func (s *AuthServiceImpl) GetProfile(userID uint) (*model.UserResponse, error) {
	user, err := s.userRepo.GetWithRoles(userID)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// UpdateProfile 更新用户资料
func (s *AuthServiceImpl) UpdateProfile(userID uint, req model.UserUpdateRequest) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	// 保存更新
	if err := s.userRepo.Update(userID, user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// ChangePassword 修改密码
func (s *AuthServiceImpl) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !user.CheckPassword(oldPassword) {
		return errors.New("old password is incorrect")
	}

	// 设置新密码
	user.Password = newPassword
	if err := user.HashPassword(); err != nil {
		return err
	}

	// 保存更新
	return s.userRepo.Update(userID, user)
}
