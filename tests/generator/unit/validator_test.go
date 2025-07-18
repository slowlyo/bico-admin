package unit

import (
	"fmt"
	"strings"
	"testing"

	"bico-admin/internal/devtools/generator"
	"bico-admin/tests/generator/utils"
)

func TestValidator_ValidateRequest(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	validator := generator.NewValidator()

	tests := []struct {
		name        string
		request     *generator.GenerateRequest
		expectError bool
		errorField  string
	}{
		{
			name: "有效的基础请求",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "User",
				Fields: []generator.FieldDefinition{
					{
						Name:     "Name",
						Type:     "string",
						GormTag:  "size:100",
						JsonTag:  "name",
						Validate: "required",
						Comment:  "名称",
					},
				},
				TableName:   "users",
				PackagePath: "internal/admin",
			},
			expectError: false,
		},
		{
			name: "无效的组件类型",
			request: &generator.GenerateRequest{
				ComponentType: "invalid_type",
				ModelName:     "User",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string"},
				},
			},
			expectError: true,
			errorField:  "component_type",
		},
		{
			name: "空模型名称",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string"},
				},
			},
			expectError: true,
			errorField:  "model_name",
		},
		{
			name: "无效的模型名称-数字开头",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "123User",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string"},
				},
			},
			expectError: true,
			errorField:  "model_name",
		},
		{
			name: "Go关键字模型名称",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "type",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string"},
				},
			},
			expectError: true,
			errorField:  "model_name",
		},
		{
			name: "小写开头的模型名称",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "user",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string"},
				},
			},
			expectError: true,
			errorField:  "model_name",
		},
		{
			name: "空字段列表",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "User",
				Fields:        []generator.FieldDefinition{},
			},
			expectError: true,
			errorField:  "fields",
		},
		{
			name: "重复字段名称",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "User",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string"},
					{Name: "Name", Type: "int"},
				},
			},
			expectError: true,
			errorField:  "fields",
		},
		{
			name: "无效的字段名称",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "User",
				Fields: []generator.FieldDefinition{
					{Name: "123invalid", Type: "string"},
				},
			},
			expectError: true,
			errorField:  "fields",
		},
		{
			name: "Go关键字字段名称",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "User",
				Fields: []generator.FieldDefinition{
					{Name: "type", Type: "string"},
				},
			},
			expectError: true,
			errorField:  "fields",
		},
		{
			name: "空字段类型",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "User",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: ""},
				},
			},
			expectError: true,
			errorField:  "fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.ValidateRequest(tt.request)

			if tt.expectError {
				if !errors.HasErrors() {
					t.Errorf("期望有验证错误，但没有错误")
					return
				}

				// 检查是否包含期望的错误字段
				found := false
				for _, err := range errors {
					if err.Field == tt.errorField {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("期望字段 '%s' 有错误，但没有找到。实际错误: %v", tt.errorField, errors)
				}
			} else {
				if errors.HasErrors() {
					t.Errorf("期望无验证错误，但得到错误: %v", errors)
				}
			}
		})
	}
}

func TestValidator_ValidateComponentType(t *testing.T) {
	validator := generator.NewValidator()

	validTypes := []generator.ComponentType{
		generator.ComponentModel,
		generator.ComponentRepository,
		generator.ComponentService,
		generator.ComponentHandler,
		generator.ComponentRoutes,
		generator.ComponentWire,
		generator.ComponentMigration,
		generator.ComponentPermission,
		generator.ComponentAll,
	}

	// 测试有效的组件类型
	for _, componentType := range validTypes {
		t.Run(string(componentType), func(t *testing.T) {
			req := &generator.GenerateRequest{
				ComponentType: componentType,
				ModelName:     "TestModel",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string"},
				},
			}

			errors := validator.ValidateRequest(req)

			// 检查组件类型是否有错误
			for _, err := range errors {
				if err.Field == "component_type" {
					t.Errorf("有效的组件类型 '%s' 不应该有验证错误: %s", componentType, err.Message)
				}
			}
		})
	}

	// 测试无效的组件类型
	invalidTypes := []string{
		"invalid",
		"",
		"MODEL",
		"model_invalid",
	}

	for _, invalidType := range invalidTypes {
		t.Run("invalid_"+invalidType, func(t *testing.T) {
			req := &generator.GenerateRequest{
				ComponentType: generator.ComponentType(invalidType),
				ModelName:     "TestModel",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string"},
				},
			}

			errors := validator.ValidateRequest(req)

			// 应该有组件类型错误
			found := false
			for _, err := range errors {
				if err.Field == "component_type" {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("无效的组件类型 '%s' 应该有验证错误", invalidType)
			}
		})
	}
}

func TestValidator_ValidateFields_EdgeCases(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	validator := generator.NewValidator()

	t.Run("字段数量超限", func(t *testing.T) {
		// 创建51个字段（超过50个限制）
		fields := make([]generator.FieldDefinition, 51)
		for i := 0; i < 51; i++ {
			fields[i] = generator.FieldDefinition{
				Name: fmt.Sprintf("Field%d", i+1),
				Type: "string",
			}
		}

		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentModel,
			ModelName:     "TestModel",
			Fields:        fields,
		}

		errors := validator.ValidateRequest(req)

		found := false
		for _, err := range errors {
			if err.Field == "fields" && strings.Contains(err.Message, "字段数量不能超过50个") {
				found = true
				break
			}
		}
		if !found {
			t.Error("应该有字段数量超限的验证错误")
		}
	})

	t.Run("字段名称为空", func(t *testing.T) {
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentModel,
			ModelName:     "TestModel",
			Fields: []generator.FieldDefinition{
				{Name: "", Type: "string"},
			},
		}

		errors := validator.ValidateRequest(req)

		found := false
		for _, err := range errors {
			if err.Field == "fields" && strings.Contains(err.Message, "字段名称不能为空") {
				found = true
				break
			}
		}
		if !found {
			t.Error("应该有字段名称为空的验证错误")
		}
	})
}
