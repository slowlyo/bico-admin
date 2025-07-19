package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Validator 验证器
type Validator struct{}

// NewValidator 创建验证器
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateRequest 验证生成请求
func (v *Validator) ValidateRequest(req *GenerateRequest) ValidationErrors {
	var errors ValidationErrors

	// 验证组件类型
	if err := v.validateComponentType(req.ComponentType); err != nil {
		errors = append(errors, ValidationError{
			Field:   "component_type",
			Message: err.Error(),
		})
	}

	// 验证模型名称
	if err := v.validateModelName(req.ModelName); err != nil {
		errors = append(errors, ValidationError{
			Field:   "model_name",
			Message: err.Error(),
		})
	}

	// 验证字段定义
	if err := v.validateFields(req.Fields); err != nil {
		errors = append(errors, ValidationError{
			Field:   "fields",
			Message: err.Error(),
		})
	}

	// 包路径验证（模型固定放在shared/models，无需验证）

	// 验证表名
	if req.TableName != "" {
		if err := v.validateTableName(req.TableName); err != nil {
			errors = append(errors, ValidationError{
				Field:   "table_name",
				Message: err.Error(),
			})
		}
	}

	return errors
}

// validateComponentType 验证组件类型
func (v *Validator) validateComponentType(componentType ComponentType) error {
	validTypes := []ComponentType{
		ComponentModel, ComponentRepository, ComponentService, ComponentHandler,
		ComponentRoutes, ComponentWire, ComponentMigration, ComponentPermission,
		ComponentFrontendAPI, ComponentFrontendPage, ComponentFrontendForm, ComponentFrontendRoute, ComponentAll,
	}

	for _, validType := range validTypes {
		if componentType == validType {
			return nil
		}
	}

	return fmt.Errorf("无效的组件类型: %s", componentType)
}

// validateModelName 验证模型名称
func (v *Validator) validateModelName(modelName string) error {
	if modelName == "" {
		return fmt.Errorf("模型名称不能为空")
	}

	if len(modelName) > 50 {
		return fmt.Errorf("模型名称长度不能超过50个字符")
	}

	if !IsValidGoIdentifier(modelName) {
		return fmt.Errorf("模型名称必须是有效的Go标识符")
	}

	if IsGoKeyword(modelName) {
		return fmt.Errorf("模型名称不能是Go关键字")
	}

	// 模型名称应该是帕斯卡命名
	if modelName[0] < 'A' || modelName[0] > 'Z' {
		return fmt.Errorf("模型名称应该以大写字母开头")
	}

	return nil
}

// validateFields 验证字段定义
func (v *Validator) validateFields(fields []FieldDefinition) error {
	if len(fields) == 0 {
		return fmt.Errorf("至少需要定义一个字段")
	}

	if len(fields) > 50 {
		return fmt.Errorf("字段数量不能超过50个")
	}

	fieldNames := make(map[string]bool)

	for i, field := range fields {
		// 验证字段名称
		if field.Name == "" {
			return fmt.Errorf("第%d个字段名称不能为空", i+1)
		}

		if !IsValidGoIdentifier(field.Name) {
			return fmt.Errorf("第%d个字段名称'%s'不是有效的Go标识符", i+1, field.Name)
		}

		if IsGoKeyword(field.Name) {
			return fmt.Errorf("第%d个字段名称'%s'不能是Go关键字", i+1, field.Name)
		}

		// 检查字段名称重复
		if fieldNames[field.Name] {
			return fmt.Errorf("字段名称'%s'重复", field.Name)
		}
		fieldNames[field.Name] = true

		// 验证字段类型
		if field.Type == "" {
			return fmt.Errorf("第%d个字段'%s'的类型不能为空", i+1, field.Name)
		}

		// 验证JSON标签
		if field.JsonTag == "" {
			// 如果没有指定JSON标签，使用字段名的蛇形命名
			field.JsonTag = ToSnakeCase(field.Name)
		}
	}

	return nil
}

// validatePackagePath 验证包路径
func (v *Validator) validatePackagePath(packagePath string) error {
	if packagePath == "" {
		return fmt.Errorf("包路径不能为空")
	}

	// 检查路径格式
	if !strings.HasPrefix(packagePath, "internal/") {
		return fmt.Errorf("包路径必须以'internal/'开头")
	}

	// 检查路径是否存在
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		return fmt.Errorf("包路径'%s'不存在", packagePath)
	}

	return nil
}

// validateTableName 验证表名
func (v *Validator) validateTableName(tableName string) error {
	if tableName == "" {
		return fmt.Errorf("表名不能为空")
	}

	if len(tableName) > 64 {
		return fmt.Errorf("表名长度不能超过64个字符")
	}

	// 表名应该是蛇形命名
	if !isValidTableName(tableName) {
		return fmt.Errorf("表名应该使用蛇形命名（小写字母、数字、下划线）")
	}

	return nil
}

// isValidTableName 检查是否为有效的表名
func isValidTableName(name string) bool {
	if name == "" {
		return false
	}

	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}

	return true
}

// CheckFileExists 检查文件是否存在
func (v *Validator) CheckFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// CheckFileConflict 检查文件冲突
func (v *Validator) CheckFileConflict(req *GenerateRequest, targetPath string) error {
	if !req.Options.OverwriteExisting {
		if v.CheckFileExists(targetPath) {
			return fmt.Errorf("文件'%s'已存在，如需覆盖请设置overwrite_existing为true", targetPath)
		}
	}

	return nil
}

// ValidateOutputPath 验证输出路径
func (v *Validator) ValidateOutputPath(outputPath string) error {
	// 检查目录是否存在，如果不存在则创建
	dir := filepath.Dir(outputPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录'%s'失败: %w", dir, err)
		}
	}

	return nil
}
