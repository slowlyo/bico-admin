package repository

import (
	"gorm.io/gorm"

	"bico-admin/core/model"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	BaseRepository[model.User]
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetWithRoles(id uint) (*model.User, error)
}

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	*BaseRepositoryImpl[model.User]
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		BaseRepositoryImpl: &BaseRepositoryImpl[model.User]{DB: db},
	}
}

// GetByUsername 根据用户名获取用户
func (r *UserRepositoryImpl) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepositoryImpl) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetWithRoles 获取用户及其角色信息
func (r *UserRepositoryImpl) GetWithRoles(id uint) (*model.User, error) {
	var user model.User
	err := r.DB.Preload("Roles").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
