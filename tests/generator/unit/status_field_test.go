package unit

import (
	"testing"

	"bico-admin/internal/devtools/generator"
)

func TestStatusFieldDetection(t *testing.T) {
	tests := []struct {
		name     string
		field    generator.FieldDefinition
		expected bool
		fieldType generator.StatusFieldType
	}{
		{
			name: "标准int状态字段",
			field: generator.FieldDefinition{
				Name: "Status",
				Type: "int",
			},
			expected:  true,
			fieldType: generator.StatusFieldTypeInt,
		},
		{
			name: "指针int状态字段",
			field: generator.FieldDefinition{
				Name: "Status",
				Type: "*int",
			},
			expected:  true,
			fieldType: generator.StatusFieldTypeIntPointer,
		},
		{
			name: "bool状态字段",
			field: generator.FieldDefinition{
				Name: "Active",
				Type: "bool",
			},
			expected:  true,
			fieldType: generator.StatusFieldTypeBool,
		},
		{
			name: "指针bool状态字段",
			field: generator.FieldDefinition{
				Name: "Enabled",
				Type: "*bool",
			},
			expected:  true,
			fieldType: generator.StatusFieldTypeBoolPointer,
		},
		{
			name: "state字段",
			field: generator.FieldDefinition{
				Name: "State",
				Type: "int",
			},
			expected:  true,
			fieldType: generator.StatusFieldTypeInt,
		},
		{
			name: "非状态字段-字符串",
			field: generator.FieldDefinition{
				Name: "Name",
				Type: "string",
			},
			expected:  false,
			fieldType: generator.StatusFieldTypeNone,
		},
		{
			name: "非状态字段-状态名但类型不对",
			field: generator.FieldDefinition{
				Name: "Status",
				Type: "string",
			},
			expected:  false,
			fieldType: generator.StatusFieldTypeNone,
		},
		{
			name: "大小写不敏感",
			field: generator.FieldDefinition{
				Name: "status",
				Type: "int",
			},
			expected:  true,
			fieldType: generator.StatusFieldTypeInt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试IsStatusField
			result := generator.IsStatusField(tt.field)
			if result != tt.expected {
				t.Errorf("IsStatusField() = %v, expected %v", result, tt.expected)
			}

			// 测试GetStatusFieldType
			fieldType := generator.GetStatusFieldType(tt.field)
			if fieldType != tt.fieldType {
				t.Errorf("GetStatusFieldType() = %v, expected %v", fieldType, tt.fieldType)
			}
		})
	}
}

func TestHasStatusField(t *testing.T) {
	tests := []struct {
		name     string
		fields   []generator.FieldDefinition
		expected bool
	}{
		{
			name: "包含状态字段",
			fields: []generator.FieldDefinition{
				{Name: "Name", Type: "string"},
				{Name: "Status", Type: "int"},
				{Name: "Description", Type: "string"},
			},
			expected: true,
		},
		{
			name: "不包含状态字段",
			fields: []generator.FieldDefinition{
				{Name: "Name", Type: "string"},
				{Name: "Description", Type: "string"},
				{Name: "Count", Type: "int"},
			},
			expected: false,
		},
		{
			name: "包含多种状态字段",
			fields: []generator.FieldDefinition{
				{Name: "Name", Type: "string"},
				{Name: "Active", Type: "bool"},
				{Name: "State", Type: "*int"},
			},
			expected: true,
		},
		{
			name: "空字段列表",
			fields: []generator.FieldDefinition{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.HasStatusField(tt.fields)
			if result != tt.expected {
				t.Errorf("HasStatusField() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestStatusFieldTypeConstants(t *testing.T) {
	// 验证常量值
	if generator.StatusFieldTypeNone != 0 {
		t.Errorf("StatusFieldTypeNone should be 0, got %d", generator.StatusFieldTypeNone)
	}
	if generator.StatusFieldTypeBool != 1 {
		t.Errorf("StatusFieldTypeBool should be 1, got %d", generator.StatusFieldTypeBool)
	}
	if generator.StatusFieldTypeInt != 2 {
		t.Errorf("StatusFieldTypeInt should be 2, got %d", generator.StatusFieldTypeInt)
	}
	if generator.StatusFieldTypeBoolPointer != 3 {
		t.Errorf("StatusFieldTypeBoolPointer should be 3, got %d", generator.StatusFieldTypeBoolPointer)
	}
	if generator.StatusFieldTypeIntPointer != 4 {
		t.Errorf("StatusFieldTypeIntPointer should be 4, got %d", generator.StatusFieldTypeIntPointer)
	}
}
