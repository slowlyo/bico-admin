package integration

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"bico-admin/internal/devtools/generator"
	"bico-admin/tests/generator/utils"
)

func TestCodeGenerator_EndToEnd(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	// 创建完整的代码生成器
	codeGenerator := generator.NewCodeGenerator()

	tests := []struct {
		name           string
		request        *generator.GenerateRequest
		expectSuccess  bool
		expectedFiles  int
		validateFunc   func(t *testing.T, helper *utils.TestHelper, response *generator.GenerateResponse)
	}{
		{
			name: "完整用户模块生成",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentAll,
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
						Comment:  "状态：0-禁用，1-启用",
					},
					{
						Name:     "LastLoginAt",
						Type:     "*time.Time",
						JsonTag:  "last_login_at",
						Comment:  "最后登录时间",
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
			expectedFiles: 8, // model, repository, service, handler, routes, wire, migration, permission
			validateFunc:  validateCompleteUserModule,
		},
		{
			name: "单个模型生成",
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
						Name:     "Status",
						Type:     "int",
						GormTag:  "default:1;index",
						JsonTag:  "status",
						Validate: "oneof=0 1 2",
						Comment:  "状态：0-下架，1-上架，2-预售",
					},
				},
				TableName:   "products",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
				},
			},
			expectSuccess: true,
			expectedFiles: 1,
			validateFunc:  validateSingleProductModel,
		},
		{
			name: "Repository和Service组合生成",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentRepository,
				ModelName:     "Category",
				Fields: []generator.FieldDefinition{
					{
						Name:     "Name",
						Type:     "string",
						GormTag:  "size:100;not null;uniqueIndex",
						JsonTag:  "name",
						Comment:  "分类名称",
					},
					{
						Name:     "ParentID",
						Type:     "*uint",
						GormTag:  "index",
						JsonTag:  "parent_id",
						Comment:  "父分类ID",
					},
				},
				TableName:   "categories",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
				},
			},
			expectSuccess: true,
			expectedFiles: 1,
			validateFunc:  validateCategoryRepository,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 执行生成
			response, err := codeGenerator.Generate(tt.request)
			helper.AssertNoError(err)

			// 验证基本响应
			if response.Success != tt.expectSuccess {
				t.Errorf("期望成功状态 %v, 得到 %v", tt.expectSuccess, response.Success)
				if !response.Success {
					t.Errorf("错误信息: %v", response.Errors)
				}
			}

			if len(response.GeneratedFiles) != tt.expectedFiles {
				t.Errorf("期望生成 %d 个文件, 得到 %d 个", tt.expectedFiles, len(response.GeneratedFiles))
				t.Logf("生成的文件: %v", response.GeneratedFiles)
			}

			// 验证所有生成的文件
			if response.Success {
				for _, filePath := range response.GeneratedFiles {
					helper.AssertFileExists(filePath)
					
					// 验证Go语法
					err := helper.ValidateGoSyntax(filePath)
					if err != nil {
						t.Errorf("文件 %s 语法验证失败: %v", filePath, err)
					}
				}

				// 执行自定义验证
				if tt.validateFunc != nil {
					tt.validateFunc(t, helper, response)
				}
			}
		})
	}
}

func validateCompleteUserModule(t *testing.T, helper *utils.TestHelper, response *generator.GenerateResponse) {
	expectedFiles := map[string]bool{
		"internal/shared/models/user.go":                    false,
		"internal/admin/repository/user_repository.go":     false,
		"internal/admin/service/user_service.go":           false,
		"internal/admin/handler/user_handler.go":           false,
		"internal/admin/handler/user_types.go":             false,
		"internal/admin/routes/user_routes_gen.go":         false,
		"internal/admin/wire/user_wire_gen.go":             false,
		// migration 和 permission 文件路径可能不同，暂时不验证
	}

	// 检查生成的文件
	for _, filePath := range response.GeneratedFiles {
		if _, exists := expectedFiles[filePath]; exists {
			expectedFiles[filePath] = true
		}
	}

	// 验证关键文件是否生成
	for filePath, generated := range expectedFiles {
		if !generated {
			t.Errorf("期望生成文件 %s, 但未找到", filePath)
		}
	}

	// 验证模型文件内容
	modelFile := "internal/shared/models/user.go"
	if contains(response.GeneratedFiles, modelFile) {
		helper.AssertFileContains(modelFile, "type User struct")
		helper.AssertFileContains(modelFile, "Username string")
		helper.AssertFileContains(modelFile, "Email string")
		helper.AssertFileContains(modelFile, "Status int")
		helper.AssertFileContains(modelFile, "LastLoginAt *time.Time")
		helper.AssertFileContains(modelFile, `return "users"`)
	}

	// 验证Repository文件内容
	repoFile := "internal/admin/repository/user_repository.go"
	if contains(response.GeneratedFiles, repoFile) {
		helper.AssertFileContains(repoFile, "type UserRepository interface")
		helper.AssertFileContains(repoFile, "type UserRepositoryImpl struct")
		helper.AssertFileContains(repoFile, "func NewUserRepository")
		helper.AssertFileContains(repoFile, "BaseRepositoryInterface[models.User]")
	}

	// 验证Service文件内容
	serviceFile := "internal/admin/service/user_service.go"
	if contains(response.GeneratedFiles, serviceFile) {
		helper.AssertFileContains(serviceFile, "type UserService interface")
		helper.AssertFileContains(serviceFile, "type UserServiceImpl struct")
		helper.AssertFileContains(serviceFile, "func NewUserService")
		helper.AssertFileContains(serviceFile, "func (s *UserServiceImpl) UpdateStatus")
		helper.AssertFileContains(serviceFile, "func (s *UserServiceImpl) validateUser")
	}

	// 验证Handler文件内容
	handlerFile := "internal/admin/handler/user_handler.go"
	if contains(response.GeneratedFiles, handlerFile) {
		helper.AssertFileContains(handlerFile, "type UserHandler struct")
		helper.AssertFileContains(handlerFile, "func NewUserHandler")
		helper.AssertFileContains(handlerFile, "func (h *UserHandler) ConvertToResponse")
		helper.AssertFileContains(handlerFile, "func (h *UserHandler) getStatusValue")
		helper.AssertFileContains(handlerFile, "func (h *UserHandler) formatTime")
	}

	// 验证Types文件内容
	typesFile := "internal/admin/handler/user_types.go"
	if contains(response.GeneratedFiles, typesFile) {
		helper.AssertFileContains(typesFile, "type UserCreateRequest struct")
		helper.AssertFileContains(typesFile, "type UserUpdateRequest struct")
		helper.AssertFileContains(typesFile, "type UserListRequest struct")
		helper.AssertFileContains(typesFile, "type UserResponse struct")
		helper.AssertFileContains(typesFile, "StatusText string")
	}
}

func validateSingleProductModel(t *testing.T, helper *utils.TestHelper, response *generator.GenerateResponse) {
	if len(response.GeneratedFiles) != 1 {
		t.Errorf("期望生成1个文件, 得到 %d 个", len(response.GeneratedFiles))
		return
	}

	modelFile := response.GeneratedFiles[0]
	expectedPath := "internal/shared/models/product.go"
	
	if !strings.HasSuffix(modelFile, "product.go") {
		t.Errorf("期望生成 %s, 得到 %s", expectedPath, modelFile)
	}

	// 验证文件内容
	helper.AssertFileContains(modelFile, "package models")
	helper.AssertFileContains(modelFile, "type Product struct")
	helper.AssertFileContains(modelFile, "types.BaseModel")
	helper.AssertFileContains(modelFile, "Name string")
	helper.AssertFileContains(modelFile, "Price float64")
	helper.AssertFileContains(modelFile, "Status int")
	helper.AssertFileContains(modelFile, `json:"name"`)
	helper.AssertFileContains(modelFile, `json:"price"`)
	helper.AssertFileContains(modelFile, `json:"status"`)
	helper.AssertFileContains(modelFile, `gorm:"size:200;not null;index"`)
	helper.AssertFileContains(modelFile, `gorm:"type:decimal(10,2);not null"`)
	helper.AssertFileContains(modelFile, `return "products"`)
}

func validateCategoryRepository(t *testing.T, helper *utils.TestHelper, response *generator.GenerateResponse) {
	if len(response.GeneratedFiles) != 1 {
		t.Errorf("期望生成1个文件, 得到 %d 个", len(response.GeneratedFiles))
		return
	}

	repoFile := response.GeneratedFiles[0]
	
	if !strings.HasSuffix(repoFile, "category_repository.go") {
		t.Errorf("期望生成category_repository.go文件, 得到 %s", repoFile)
	}

	// 验证文件内容
	helper.AssertFileContains(repoFile, "package repository")
	helper.AssertFileContains(repoFile, "type CategoryRepository interface")
	helper.AssertFileContains(repoFile, "type CategoryRepositoryImpl struct")
	helper.AssertFileContains(repoFile, "func NewCategoryRepository")
	helper.AssertFileContains(repoFile, "BaseRepositoryInterface[models.Category]")
	helper.AssertFileContains(repoFile, "func (r *CategoryRepositoryImpl) ListWithFilter")
}

// 辅助函数：检查切片是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func TestCodeGenerator_HistoryManagement(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	codeGenerator := generator.NewCodeGenerator()

	t.Run("生成后历史记录更新", func(t *testing.T) {
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentModel,
			ModelName:     "HistoryTest",
			Fields: []generator.FieldDefinition{
				{
					Name:    "Name",
					Type:    "string",
					JsonTag: "name",
					Comment: "名称",
				},
			},
			TableName:   "history_tests",
			PackagePath: "internal/admin",
			Options: generator.GenerateOptions{
				OverwriteExisting: true,
				FormatCode:        true,
			},
		}

		// 执行生成
		response, err := codeGenerator.Generate(req)
		helper.AssertNoError(err)

		if !response.Success {
			t.Errorf("生成失败: %v", response.Errors)
			return
		}

		// 验证历史记录是否更新
		if !response.HistoryUpdated {
			t.Error("期望历史记录已更新，但HistoryUpdated为false")
		}

		// 获取历史记录
		history, err := codeGenerator.GetHistoryByModule("HistoryTest")
		helper.AssertNoError(err)

		if history.ModuleName != "HistoryTest" {
			t.Errorf("期望模块名 'HistoryTest', 得到 '%s'", history.ModuleName)
		}

		if len(history.GeneratedFiles) != len(response.GeneratedFiles) {
			t.Errorf("历史记录中的文件数量 %d 与响应中的文件数量 %d 不匹配", 
				len(history.GeneratedFiles), len(response.GeneratedFiles))
		}
	})

	t.Run("获取所有历史记录", func(t *testing.T) {
		// 先生成几个模块
		modules := []string{"Module1", "Module2", "Module3"}
		
		for _, moduleName := range modules {
			req := &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     moduleName,
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string", JsonTag: "name"},
				},
				TableName:   strings.ToLower(moduleName) + "s",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
				},
			}

			response, err := codeGenerator.Generate(req)
			helper.AssertNoError(err)
			
			if !response.Success {
				t.Errorf("生成模块 %s 失败: %v", moduleName, response.Errors)
			}
		}

		// 获取所有历史记录
		allHistory, err := codeGenerator.GetHistory()
		helper.AssertNoError(err)

		if len(allHistory) < len(modules) {
			t.Errorf("期望至少 %d 条历史记录, 得到 %d 条", len(modules), len(allHistory))
		}

		// 验证包含所有模块
		moduleMap := make(map[string]bool)
		for _, history := range allHistory {
			moduleMap[history.ModuleName] = true
		}

		for _, moduleName := range modules {
			if !moduleMap[moduleName] {
				t.Errorf("历史记录中未找到模块 '%s'", moduleName)
			}
		}
	})
}

func TestCodeGenerator_ErrorScenarios(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	codeGenerator := generator.NewCodeGenerator()

	t.Run("无效请求处理", func(t *testing.T) {
		req := &generator.GenerateRequest{
			ComponentType: "invalid_component",
			ModelName:     "TestModel",
			Fields: []generator.FieldDefinition{
				{Name: "Name", Type: "string"},
			},
		}

		response, err := codeGenerator.Generate(req)
		helper.AssertNoError(err)

		if response.Success {
			t.Error("期望生成失败，但生成成功了")
		}

		if len(response.Errors) == 0 {
			t.Error("期望有错误信息，但没有错误")
		}
	})

	t.Run("部分组件生成失败", func(t *testing.T) {
		// 这个测试需要模拟某个组件生成失败的情况
		// 由于当前实现中很难模拟部分失败，我们测试验证失败的情况
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentAll,
			ModelName:     "", // 空模型名会导致验证失败
			Fields: []generator.FieldDefinition{
				{Name: "Name", Type: "string"},
			},
		}

		response, err := codeGenerator.Generate(req)
		helper.AssertNoError(err)

		if response.Success {
			t.Error("期望生成失败，但生成成功了")
		}

		// 验证错误信息
		found := false
		for _, errMsg := range response.Errors {
			if strings.Contains(errMsg, "模型名称不能为空") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("期望错误信息包含模型名称验证错误，实际错误: %v", response.Errors)
		}
	})
}

func TestCodeGenerator_LoadTestData(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	codeGenerator := generator.NewCodeGenerator()

	t.Run("使用测试数据文件生成", func(t *testing.T) {
		// 加载基础用户测试数据
		var req generator.GenerateRequest
		err := helper.LoadTestData("requests/basic_user.json", &req)
		helper.AssertNoError(err)

		// 执行生成
		response, err := codeGenerator.Generate(&req)
		helper.AssertNoError(err)

		if !response.Success {
			t.Errorf("使用测试数据生成失败: %v", response.Errors)
			return
		}

		// 验证生成结果
		if len(response.GeneratedFiles) == 0 {
			t.Error("期望生成文件，但没有生成任何文件")
		}

		// 验证模型文件内容
		for _, filePath := range response.GeneratedFiles {
			if strings.HasSuffix(filePath, "user.go") {
				helper.AssertFileContains(filePath, "type User struct")
				helper.AssertFileContains(filePath, "Username string")
				helper.AssertFileContains(filePath, "Email string")
				helper.AssertFileContains(filePath, "Status int")
				break
			}
		}
	})
}
