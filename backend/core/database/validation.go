package database

import (
	"errors"
	"gorm.io/gorm"
)

// ValidationHelper 数据验证助手 - 提供通用的数据验证方法
type ValidationHelper[T any] struct {
	ops *Operations[T]
}

// NewValidationHelper 创建验证助手实例
func NewValidationHelper[T any](ops *Operations[T]) *ValidationHelper[T] {
	return &ValidationHelper[T]{
		ops: ops,
	}
}

// CheckUniqueField 检查字段值是否唯一
func (v *ValidationHelper[T]) CheckUniqueField(field, value string, excludeID ...uint) error {
	if value == "" {
		return nil // 空值不检查
	}

	condition := field + " = ?"
	args := []interface{}{value}

	// 如果提供了排除ID，添加到条件中
	if len(excludeID) > 0 && excludeID[0] > 0 {
		condition += " AND id != ?"
		args = append(args, excludeID[0])
	}

	existing, err := v.ops.GetByCondition(condition, args...)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existing != nil {
		return errors.New(field + " already exists")
	}

	return nil
}

// CheckMultipleUniqueFields 检查多个字段的唯一性
func (v *ValidationHelper[T]) CheckMultipleUniqueFields(fields map[string]string, excludeID ...uint) error {
	for field, value := range fields {
		if err := v.CheckUniqueField(field, value, excludeID...); err != nil {
			return err
		}
	}
	return nil
}

// UserValidationHelper 用户验证助手 - 专门用于用户相关验证
type UserValidationHelper struct {
	*ValidationHelper[interface{}] // 使用interface{}作为泛型类型
	db *gorm.DB
}

// NewUserValidationHelper 创建用户验证助手
func NewUserValidationHelper(db *gorm.DB) *UserValidationHelper {
	return &UserValidationHelper{
		db: db,
	}
}

// CheckUsernameUnique 检查用户名是否唯一
func (v *UserValidationHelper) CheckUsernameUnique(username string, excludeID ...uint) error {
	if username == "" {
		return nil
	}

	condition := "username = ?"
	args := []interface{}{username}

	if len(excludeID) > 0 && excludeID[0] > 0 {
		condition += " AND id != ?"
		args = append(args, excludeID[0])
	}

	var count int64
	err := v.db.Table("users").Where(condition, args...).Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("username already exists")
	}

	return nil
}

// CheckEmailUnique 检查邮箱是否唯一
func (v *UserValidationHelper) CheckEmailUnique(email string, excludeID ...uint) error {
	if email == "" {
		return nil
	}

	condition := "email = ?"
	args := []interface{}{email}

	if len(excludeID) > 0 && excludeID[0] > 0 {
		condition += " AND id != ?"
		args = append(args, excludeID[0])
	}

	var count int64
	err := v.db.Table("users").Where(condition, args...).Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("email already exists")
	}

	return nil
}

// CheckUserUniqueFields 检查用户的唯一字段
func (v *UserValidationHelper) CheckUserUniqueFields(username, email string, excludeID ...uint) error {
	if err := v.CheckUsernameUnique(username, excludeID...); err != nil {
		return err
	}

	if err := v.CheckEmailUnique(email, excludeID...); err != nil {
		return err
	}

	return nil
}
