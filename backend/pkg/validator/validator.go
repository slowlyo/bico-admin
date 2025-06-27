package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator 验证器实例
var Validator *validator.Validate

// ValidationError 验证错误结构
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// init 初始化验证器
func init() {
	Validator = validator.New()
	
	// 注册自定义验证器
	registerCustomValidators()
	
	// 使用JSON标签作为字段名
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// registerCustomValidators 注册自定义验证器
func registerCustomValidators() {
	// 注册手机号验证器
	Validator.RegisterValidation("phone", validatePhone)
	
	// 注册用户名验证器
	Validator.RegisterValidation("username", validateUsername)
}

// validatePhone 验证手机号
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true // 空值由required验证
	}
	
	// 简单的手机号验证（中国大陆）
	if len(phone) != 11 {
		return false
	}
	
	if !strings.HasPrefix(phone, "1") {
		return false
	}
	
	return true
}

// validateUsername 验证用户名
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if username == "" {
		return true // 空值由required验证
	}
	
	// 用户名只能包含字母、数字、下划线
	for _, char := range username {
		if !((char >= 'a' && char <= 'z') || 
			 (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') || 
			 char == '_') {
			return false
		}
	}
	
	return true
}

// Validate 验证结构体
func Validate(data interface{}) []ValidationError {
	var errors []ValidationError
	
	err := Validator.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   fmt.Sprintf("%v", err.Value()),
				Message: getErrorMessage(err),
			})
		}
	}
	
	return errors
}

// getErrorMessage 获取错误消息
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", err.Field(), err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param())
	case "phone":
		return fmt.Sprintf("%s must be a valid phone number", err.Field())
	case "username":
		return fmt.Sprintf("%s can only contain letters, numbers and underscores", err.Field())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}

// ValidateVar 验证单个变量
func ValidateVar(field interface{}, tag string) error {
	return Validator.Var(field, tag)
}
