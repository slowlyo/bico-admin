package unit

import (
	"path/filepath"
	"strings"
	"testing"

	"bico-admin/internal/devtools/generator"
	"bico-admin/tests/generator/utils"
)

func TestServiceGenerator_Generate(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	serviceGenerator := generator.NewServiceGenerator()

	tests := []struct {
		name          string
		request       *generator.GenerateRequest
		expectSuccess bool
		expectedFiles int
	}{
		{
			name: "基础Service生成",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentService,
				ModelName:     "User",
				Fields: []generator.FieldDefinition{
					{
						Name:    "Username",
						Type:    "string",
						GormTag: "uniqueIndex;size:50",
						JsonTag: "username",
						Comment: "用户名",
					},
					{
						Name:    "Email",
						Type:    "string",
						GormTag: "uniqueIndex;size:100",
						JsonTag: "email",
						Comment: "邮箱",
					},
					{
						Name:    "Status",
						Type:    "int",
						GormTag: "default:1",
						JsonTag: "status",
						Comment: "状态",
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
			name: "包含状态字段的Service生成",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentService,
				ModelName:     "Product",
				Fields: []generator.FieldDefinition{
					{
						Name:    "Name",
						Type:    "string",
						GormTag: "size:200;not null",
						JsonTag: "name",
						Comment: "产品名称",
					},
					{
						Name:    "Status",
						Type:    "int",
						GormTag: "default:1;index",
						JsonTag: "status",
						Comment: "状态：0-下架，1-上架",
					},
					{
						Name:    "Price",
						Type:    "float64",
						GormTag: "type:decimal(10,2)",
						JsonTag: "price",
						Comment: "价格",
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := serviceGenerator.Generate(tt.request)
			helper.AssertNoError(err)

			if response.Success != tt.expectSuccess {
				t.Errorf("期望成功状态 %v, 得到 %v", tt.expectSuccess, response.Success)
			}

			if len(response.GeneratedFiles) != tt.expectedFiles {
				t.Errorf("期望生成 %d 个文件, 得到 %d 个", tt.expectedFiles, len(response.GeneratedFiles))
			}

			if response.Success && len(response.GeneratedFiles) > 0 {
				for _, filePath := range response.GeneratedFiles {
					helper.AssertFileExists(filePath)
					err := helper.ValidateGoSyntax(filePath)
					helper.AssertNoError(err)

					t.Run("验证Service文件内容_"+filepath.Base(filePath), func(t *testing.T) {
						validateServiceFileContent(t, helper, filePath, tt.request)
					})
				}
			}
		})
	}
}

func validateServiceFileContent(t *testing.T, helper *utils.TestHelper, filePath string, req *generator.GenerateRequest) {
	// 验证包声明
	helper.AssertFileContains(filePath, "package service")

	// 验证导入
	helper.AssertFileContains(filePath, `"context"`)
	helper.AssertFileContains(filePath, `"bico-admin/internal/shared/models"`)
	helper.AssertFileContains(filePath, `"bico-admin/internal/shared/service"`)
	helper.AssertFileContains(filePath, `"bico-admin/internal/shared/types"`)
	helper.AssertFileContains(filePath, `"bico-admin/internal/admin/repository"`)

	// 验证接口定义
	interfaceName := req.ModelName + "Service"
	helper.AssertFileContains(filePath, "type "+interfaceName+" interface {")
	helper.AssertFileContains(filePath, "service.BaseServiceInterface[models."+req.ModelName+", repository."+req.ModelName+"Repository]")

	// 验证实现结构体
	implName := req.ModelName + "ServiceImpl"
	helper.AssertFileContains(filePath, "type "+implName+" struct {")
	helper.AssertFileContains(filePath, "*service.BaseService[models."+req.ModelName+", repository."+req.ModelName+"Repository]")

	// 验证构造函数
	constructorName := "New" + req.ModelName + "Service"
	helper.AssertFileContains(filePath, "func "+constructorName+"(repo repository."+req.ModelName+"Repository) "+interfaceName+" {")

	// 验证业务方法
	helper.AssertFileContains(filePath, "func (s *"+implName+") ListWithFilter")
	helper.AssertFileContains(filePath, "func (s *"+implName+") Create")
	helper.AssertFileContains(filePath, "func (s *"+implName+") Update")
	helper.AssertFileContains(filePath, "func (s *"+implName+") Delete")

	// 检查是否有状态字段，验证UpdateStatus方法
	hasStatusField := false
	for _, field := range req.Fields {
		if field.Name == "Status" {
			hasStatusField = true
			break
		}
	}
	if hasStatusField {
		helper.AssertFileContains(filePath, "func (s *"+implName+") UpdateStatus")
	}

	// 验证验证方法
	helper.AssertFileContains(filePath, "func (s *"+implName+") validate"+req.ModelName)
	helper.AssertFileContains(filePath, "func (s *"+implName+") validateDelete"+req.ModelName)
	if hasStatusField {
		helper.AssertFileContains(filePath, "func (s *"+implName+") validateStatusUpdate"+req.ModelName)
	}

	// 验证返回语句
	helper.AssertFileContains(filePath, "return &"+implName+"{")
}

func TestServiceGenerator_ValidationMethods(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	serviceGenerator := generator.NewServiceGenerator()

	req := &generator.GenerateRequest{
		ComponentType: generator.ComponentService,
		ModelName:     "Category",
		Fields: []generator.FieldDefinition{
			{
				Name:    "Name",
				Type:    "string",
				GormTag: "size:100;not null;uniqueIndex",
				JsonTag: "name",
				Comment: "分类名称",
			},
			{
				Name:    "ParentID",
				Type:    "*uint",
				GormTag: "index",
				JsonTag: "parent_id",
				Comment: "父分类ID",
			},
			{
				Name:    "Status",
				Type:    "int",
				GormTag: "default:1",
				JsonTag: "status",
				Comment: "状态",
			},
		},
		TableName:   "categories",
		PackagePath: "internal/admin",
		Options: generator.GenerateOptions{
			OverwriteExisting: true,
			FormatCode:        true,
		},
	}

	response, err := serviceGenerator.Generate(req)
	helper.AssertNoError(err)

	if !response.Success {
		t.Errorf("生成失败: %v", response.Errors)
		return
	}

	if len(response.GeneratedFiles) > 0 {
		filePath := response.GeneratedFiles[0]

		// 验证基础验证方法
		helper.AssertFileContains(filePath, "func (s *CategoryServiceImpl) validateCategory")
		helper.AssertFileContains(filePath, "// TODO: 实现Category实体验证逻辑")

		// 验证删除验证方法
		helper.AssertFileContains(filePath, "func (s *CategoryServiceImpl) validateDeleteCategory")
		helper.AssertFileContains(filePath, "// TODO: 实现Category删除前验证逻辑")

		// 验证状态更新验证方法（因为有Status字段）
		helper.AssertFileContains(filePath, "func (s *CategoryServiceImpl) validateStatusUpdateCategory")
		helper.AssertFileContains(filePath, "// TODO: 实现Category状态更新验证逻辑")

		// 验证验证方法的调用
		helper.AssertFileContains(filePath, "if err := s.validateCategory(entity); err != nil {")
		helper.AssertFileContains(filePath, "if err := s.validateDeleteCategory(entity); err != nil {")
		helper.AssertFileContains(filePath, "if err := s.validateStatusUpdateCategory(entity); err != nil {")
	}
}

func TestServiceGenerator_BusinessLogic(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	serviceGenerator := generator.NewServiceGenerator()

	req := &generator.GenerateRequest{
		ComponentType: generator.ComponentService,
		ModelName:     "Order",
		Fields: []generator.FieldDefinition{
			{
				Name:    "OrderNo",
				Type:    "string",
				GormTag: "uniqueIndex;size:50;not null",
				JsonTag: "order_no",
				Comment: "订单号",
			},
			{
				Name:    "UserID",
				Type:    "uint",
				GormTag: "not null;index",
				JsonTag: "user_id",
				Comment: "用户ID",
			},
			{
				Name:    "TotalAmount",
				Type:    "float64",
				GormTag: "type:decimal(10,2);not null",
				JsonTag: "total_amount",
				Comment: "总金额",
			},
			{
				Name:    "Status",
				Type:    "int",
				GormTag: "default:0;index",
				JsonTag: "status",
				Comment: "订单状态",
			},
		},
		TableName:   "orders",
		PackagePath: "internal/admin",
		Options: generator.GenerateOptions{
			OverwriteExisting: true,
			FormatCode:        true,
		},
	}

	response, err := serviceGenerator.Generate(req)
	helper.AssertNoError(err)

	if !response.Success {
		t.Errorf("生成失败: %v", response.Errors)
		return
	}

	if len(response.GeneratedFiles) > 0 {
		filePath := response.GeneratedFiles[0]

		// 验证ListWithFilter方法的业务逻辑
		helper.AssertFileContains(filePath, "func (s *OrderServiceImpl) ListWithFilter")
		helper.AssertFileContains(filePath, "if query.Page <= 0 {")
		helper.AssertFileContains(filePath, "if query.PageSize <= 0 {")
		helper.AssertFileContains(filePath, "return s.repository.ListWithFilter(query)")

		// 验证Create方法的业务逻辑
		helper.AssertFileContains(filePath, "func (s *OrderServiceImpl) Create")
		helper.AssertFileContains(filePath, "if err := s.validateOrder(entity); err != nil {")
		helper.AssertFileContains(filePath, "return s.BaseService.Create(ctx, entity)")

		// 验证Update方法的业务逻辑
		helper.AssertFileContains(filePath, "func (s *OrderServiceImpl) Update")
		helper.AssertFileContains(filePath, "if err := s.validateOrder(entity); err != nil {")
		helper.AssertFileContains(filePath, "return s.BaseService.Update(ctx, entity)")

		// 验证Delete方法的业务逻辑
		helper.AssertFileContains(filePath, "func (s *OrderServiceImpl) Delete")
		helper.AssertFileContains(filePath, "if err := s.validateDeleteOrder(existing); err != nil {")

		// 验证UpdateStatus方法（因为有Status字段）
		helper.AssertFileContains(filePath, "func (s *OrderServiceImpl) UpdateStatus")
		helper.AssertFileContains(filePath, "if id <= 0 {")
		helper.AssertFileContains(filePath, "if status < 0 {")
		helper.AssertFileContains(filePath, "if err := s.validateStatusUpdateOrder(existing); err != nil {")
	}
}

func TestServiceGenerator_ErrorHandling(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	serviceGenerator := generator.NewServiceGenerator()

	t.Run("无效模型名称", func(t *testing.T) {
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentService,
			ModelName:     "123Invalid",
			Fields: []generator.FieldDefinition{
				{Name: "Name", Type: "string"},
			},
		}

		response, err := serviceGenerator.Generate(req)
		helper.AssertNoError(err)

		if response.Success {
			t.Error("期望生成失败，但生成成功了")
		}

		found := false
		for _, errMsg := range response.Errors {
			if strings.Contains(errMsg, "模型名称必须是有效的Go标识符") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("期望错误信息包含模型名称验证错误，实际错误: %v", response.Errors)
		}
	})

	t.Run("空字段列表", func(t *testing.T) {
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentService,
			ModelName:     "EmptyService",
			Fields:        []generator.FieldDefinition{},
			TableName:     "empty_services",
			PackagePath:   "internal/admin",
		}

		response, err := serviceGenerator.Generate(req)
		helper.AssertNoError(err)

		if response.Success {
			t.Error("期望生成失败，但生成成功了")
		}

		found := false
		for _, errMsg := range response.Errors {
			if strings.Contains(errMsg, "至少需要定义一个字段") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("期望错误信息包含字段验证错误，实际错误: %v", response.Errors)
		}
	})
}
