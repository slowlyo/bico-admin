package unit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bico-admin/internal/devtools/generator"
	"bico-admin/tests/generator/utils"
)

func TestModelGenerator_Generate(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	// 创建模型生成器
	modelGenerator := generator.NewModelGenerator()

	tests := []struct {
		name           string
		request        *generator.GenerateRequest
		expectSuccess  bool
		expectedFiles  int
		expectedErrors []string
	}{
		{
			name: "基础用户模型生成",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "User",
				Fields: []generator.FieldDefinition{
					{
						Name:     "Username",
						Type:     "string",
						GormTag:  "uniqueIndex;size:50;not null",
						JsonTag:  "username",
						Validate: "required,min=3,max=50",
						Comment:  "用户名",
					},
					{
						Name:     "Email",
						Type:     "string",
						GormTag:  "uniqueIndex;size:100",
						JsonTag:  "email",
						Validate: "required,email",
						Comment:  "邮箱",
					},
					{
						Name:     "Status",
						Type:     "int",
						GormTag:  "default:1",
						JsonTag:  "status",
						Validate: "oneof=0 1",
						Comment:  "状态",
					},
				},
				TableName:   "users",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
					OptimizeImports:   true,
				},
			},
			expectSuccess: true,
			expectedFiles: 1,
		},
		{
			name: "复杂产品模型生成",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "Product",
				Fields: []generator.FieldDefinition{
					{
						Name:     "Name",
						Type:     "string",
						GormTag:  "size:200;not null;index",
						JsonTag:  "name",
						Validate: "required,max=200",
						Comment:  "产品名称",
					},
					{
						Name:     "Price",
						Type:     "float64",
						GormTag:  "type:decimal(10,2);not null",
						JsonTag:  "price",
						Validate: "required,min=0",
						Comment:  "价格",
					},
					{
						Name:     "PublishedAt",
						Type:     "*time.Time",
						GormTag:  "",
						JsonTag:  "published_at",
						Validate: "",
						Comment:  "发布时间",
					},
					{
						Name:     "IsActive",
						Type:     "bool",
						GormTag:  "default:true",
						JsonTag:  "is_active",
						Validate: "",
						Comment:  "是否激活",
					},
				},
				TableName:   "products",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
					OptimizeImports:   true,
				},
			},
			expectSuccess: true,
			expectedFiles: 1,
		},
		{
			name: "无效请求-空字段",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "EmptyModel",
				Fields:        []generator.FieldDefinition{},
				TableName:     "empty_models",
				PackagePath:   "internal/admin",
			},
			expectSuccess:  false,
			expectedFiles:  0,
			expectedErrors: []string{"至少需要定义一个字段"},
		},
		{
			name: "无效请求-无效模型名",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "123Invalid",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string"},
				},
				TableName:   "invalid_models",
				PackagePath: "internal/admin",
			},
			expectSuccess:  false,
			expectedFiles:  0,
			expectedErrors: []string{"模型名称必须是有效的Go标识符"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 执行生成
			response, err := modelGenerator.Generate(tt.request)

			// 检查基本错误
			helper.AssertNoError(err)

			// 检查成功状态
			if response.Success != tt.expectSuccess {
				t.Errorf("期望成功状态 %v, 得到 %v", tt.expectSuccess, response.Success)
			}

			// 检查生成文件数量
			if len(response.GeneratedFiles) != tt.expectedFiles {
				t.Errorf("期望生成 %d 个文件, 得到 %d 个", tt.expectedFiles, len(response.GeneratedFiles))
			}

			// 检查错误信息
			if len(tt.expectedErrors) > 0 {
				if len(response.Errors) == 0 {
					t.Errorf("期望有错误信息，但没有错误")
				} else {
					for _, expectedError := range tt.expectedErrors {
						found := false
						for _, actualError := range response.Errors {
							if strings.Contains(actualError, expectedError) {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("期望错误信息包含 '%s', 实际错误: %v", expectedError, response.Errors)
						}
					}
				}
			}

			// 如果生成成功，验证文件内容
			if response.Success && len(response.GeneratedFiles) > 0 {
				for _, filePath := range response.GeneratedFiles {
					// 检查文件是否存在
					helper.AssertFileExists(filePath)

					// 验证Go语法
					err := helper.ValidateGoSyntax(filePath)
					helper.AssertNoError(err)

					// 验证文件内容
					t.Run("验证文件内容_"+filepath.Base(filePath), func(t *testing.T) {
						validateModelFileContent(t, helper, filePath, tt.request)
					})
				}
			}
		})
	}
}

func validateModelFileContent(t *testing.T, helper *utils.TestHelper, filePath string, req *generator.GenerateRequest) {
	// 验证包声明
	helper.AssertFileContains(filePath, "package models")

	// 验证导入
	helper.AssertFileContains(filePath, `"bico-admin/internal/shared/types"`)

	// 检查是否有时间字段，验证time包导入
	hasTimeField := false
	for _, field := range req.Fields {
		if strings.Contains(field.Type, "time.Time") {
			hasTimeField = true
			break
		}
	}
	if hasTimeField {
		helper.AssertFileContains(filePath, `"time"`)
	}

	// 验证结构体定义
	helper.AssertFileContains(filePath, "type "+req.ModelName+" struct {")
	helper.AssertFileContains(filePath, "types.BaseModel")

	// 验证字段定义
	for _, field := range req.Fields {
		// 检查字段声明
		helper.AssertFileContains(filePath, field.Name+" "+field.Type)

		// 检查JSON标签
		if field.JsonTag != "" {
			helper.AssertFileContains(filePath, `json:"`+field.JsonTag+`"`)
		}

		// 检查GORM标签
		if field.GormTag != "" {
			helper.AssertFileContains(filePath, `gorm:"`+field.GormTag+`"`)
		}
	}

	// 验证TableName方法
	helper.AssertFileContains(filePath, "func ("+req.ModelName+") TableName() string {")
	helper.AssertFileContains(filePath, `return "`+req.TableName+`"`)

	// 验证注释
	helper.AssertFileContains(filePath, "// "+req.ModelName+" "+req.ModelName+"模型")
}

func TestModelGenerator_PrepareTemplateData(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	modelGenerator := generator.NewModelGenerator()

	req := &generator.GenerateRequest{
		ComponentType: generator.ComponentModel,
		ModelName:     "UserProfile",
		Fields: []generator.FieldDefinition{
			{
				Name:    "FirstName",
				Type:    "string",
				GormTag: "size:50",
				JsonTag: "first_name",
				Comment: "名字",
			},
			{
				Name:    "LastName",
				Type:    "string",
				GormTag: "size:50",
				JsonTag: "last_name",
				Comment: "姓氏",
			},
			{
				Name:    "BirthDate",
				Type:    "*time.Time",
				JsonTag: "birth_date",
				Comment: "出生日期",
			},
		},
		TableName:   "user_profiles",
		PackagePath: "internal/admin",
	}

	// 这里我们需要通过反射或其他方式测试私有方法
	// 由于prepareTemplateData是私有方法，我们通过Generate方法间接测试
	response, err := modelGenerator.Generate(req)
	helper.AssertNoError(err)

	if !response.Success {
		t.Errorf("生成失败: %v", response.Errors)
		return
	}

	// 验证生成的文件包含正确的模板数据处理结果
	if len(response.GeneratedFiles) > 0 {
		filePath := response.GeneratedFiles[0]

		// 验证模型名称转换
		helper.AssertFileContains(filePath, "type UserProfile struct")

		// 验证表名转换
		helper.AssertFileContains(filePath, `return "user_profiles"`)

		// 验证时间字段导入
		helper.AssertFileContains(filePath, `"time"`)

		// 验证字段定义
		helper.AssertFileContains(filePath, "FirstName string")
		helper.AssertFileContains(filePath, "LastName string")
		helper.AssertFileContains(filePath, "BirthDate *time.Time")

		// 验证JSON标签
		helper.AssertFileContains(filePath, `json:"first_name"`)
		helper.AssertFileContains(filePath, `json:"last_name"`)
		helper.AssertFileContains(filePath, `json:"birth_date"`)
	}
}

func TestModelGenerator_FileConflict(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	modelGenerator := generator.NewModelGenerator()

	req := &generator.GenerateRequest{
		ComponentType: generator.ComponentModel,
		ModelName:     "TestModel",
		Fields: []generator.FieldDefinition{
			{Name: "Name", Type: "string", JsonTag: "name"},
		},
		TableName:   "test_models",
		PackagePath: "internal/admin",
		Options: generator.GenerateOptions{
			OverwriteExisting: false, // 不覆盖现有文件
		},
	}

	// 创建一个已存在的文件
	existingFilePath := "internal/shared/models/test_model.go"
	err := os.MkdirAll(filepath.Dir(existingFilePath), 0755)
	helper.AssertNoError(err)

	err = os.WriteFile(existingFilePath, []byte("existing content"), 0644)
	helper.AssertNoError(err)

	// 尝试生成，应该失败
	response, err := modelGenerator.Generate(req)
	helper.AssertNoError(err)

	if response.Success {
		t.Error("期望生成失败由于文件冲突，但生成成功了")
	}

	// 检查错误信息
	found := false
	for _, errMsg := range response.Errors {
		if strings.Contains(errMsg, "已存在") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("期望错误信息包含文件冲突提示，实际错误: %v", response.Errors)
	}

	// 清理
	os.Remove(existingFilePath)
}
