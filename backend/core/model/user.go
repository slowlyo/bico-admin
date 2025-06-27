package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
type User struct {
	BaseModel
	Username    string     `json:"username" gorm:"uniqueIndex;size:50;not null" validate:"required,min=3,max=50"`
	Email       string     `json:"email" gorm:"uniqueIndex;size:100;not null" validate:"required,email"`
	Password    string     `json:"-" gorm:"size:255;not null" validate:"required,min=6"`
	Nickname    string     `json:"nickname" gorm:"size:50"`
	Avatar      string     `json:"avatar" gorm:"size:255"`
	Phone       string     `json:"phone" gorm:"size:20"`
	Status      UserStatus `json:"status" gorm:"default:1"`
	LastLoginAt *time.Time `json:"last_login_at"`
	LastLoginIP string     `json:"last_login_ip" gorm:"size:45"`
	
	// 关联关系
	Roles []Role `json:"roles" gorm:"many2many:user_roles;"`
}

// UserStatus 用户状态
type UserStatus int

const (
	UserStatusInactive UserStatus = 0 // 未激活
	UserStatusActive   UserStatus = 1 // 激活
	UserStatusBlocked  UserStatus = 2 // 被封禁
)

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Nickname string `json:"nickname" validate:"max=50"`
	Phone    string `json:"phone" validate:"max=20"`
	RoleIDs  []uint `json:"role_ids"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	Username string     `json:"username" validate:"min=3,max=50"`
	Email    string     `json:"email" validate:"email"`
	Nickname string     `json:"nickname" validate:"max=50"`
	Phone    string     `json:"phone" validate:"max=20"`
	Status   UserStatus `json:"status" validate:"oneof=0 1 2"`
	RoleIDs  []uint     `json:"role_ids"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID          uint       `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Nickname    string     `json:"nickname"`
	Avatar      string     `json:"avatar"`
	Phone       string     `json:"phone"`
	Status      UserStatus `json:"status"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Roles       []Role     `json:"roles"`
}

// HashPassword 加密密码
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		Nickname:    u.Nickname,
		Avatar:      u.Avatar,
		Phone:       u.Phone,
		Status:      u.Status,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Roles:       u.Roles,
	}
}
