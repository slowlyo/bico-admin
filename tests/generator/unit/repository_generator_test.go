package unit

import (
	"path/filepath"
	"strings"
	"testing"

	"bico-admin/internal/devtools/generator"
	"bico-admin/tests/generator/utils"
)

func TestRepositoryGenerator_Generate(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	repositoryGenerator := generator.NewRepositoryGenerator()

	tests := []struct {
		name          string
		request       *generator.GenerateRequest
		expectSuccess bool
		expectedFiles int
	}{
		{
			name: "基础Repository生成",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentRepository,
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
			name: "复杂Repository生成",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentRepository,
				ModelName:     "Product",
				Fields: []generator.FieldDefinition{
					{
						Name:    "Name",
						Type:    "string",
						GormTag: "size:200;not null;index",
						JsonTag: "name",
						Comment: "产品名称",
					},
					{
						Name:    "Price",
						Type:    "float64",
						GormTag: "type:decimal(10,2);not null",
						JsonTag: "price",
						Comment: "价格",
					},
					{
						Name:    "CategoryID",
						Type:    "uint",
						GormTag: "not null;index",
						JsonTag: "category_id",
						Comment: "分类ID",
					},
					{
						Name:    "PublishedAt",
						Type:    "*time.Time",
						JsonTag: "published_at",
						Comment: "发布时间",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := repositoryGenerator.Generate(tt.request)
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

					t.Run("验证Repository文件内容_"+filepath.Base(filePath), func(t *testing.T) {
						validateRepositoryFileContent(t, helper, filePath, tt.request)
					})
				}
			}
		})
	}
}

func validateRepositoryFileContent(t *testing.T, helper *utils.TestHelper, filePath string, req *generator.GenerateRequest) {
	// 验证包声明
	helper.AssertFileContains(filePath, "package repository")

	// 验证导入
	helper.AssertFileContains(filePath, `"bico-admin/internal/shared/models"`)
	helper.AssertFileContains(filePath, `"bico-admin/internal/shared/repository"`)
	helper.AssertFileContains(filePath, `"gorm.io/gorm"`)

	// 验证接口定义
	interfaceName := req.ModelName + "Repository"
	helper.AssertFileContains(filePath, "type "+interfaceName+" interface {")
	helper.AssertFileContains(filePath, "repository.BaseRepositoryInterface[models."+req.ModelName+"]")

	// 验证实现结构体
	implName := req.ModelName + "RepositoryImpl"
	helper.AssertFileContains(filePath, "type "+implName+" struct {")
	helper.AssertFileContains(filePath, "*repository.BaseRepository[models."+req.ModelName+"]")

	// 验证构造函数
	constructorName := "New" + req.ModelName + "Repository"
	helper.AssertFileContains(filePath, "func "+constructorName+"(db *gorm.DB) "+interfaceName+" {")

	// 验证ListWithFilter方法
	helper.AssertFileContains(filePath, "func (r *"+implName+") ListWithFilter")
	helper.AssertFileContains(filePath, "types.BasePageQuery")

	// 验证返回语句
	helper.AssertFileContains(filePath, "return &"+implName+"{")
}

func TestRepositoryGenerator_TemplateData(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	repositoryGenerator := generator.NewRepositoryGenerator()

	req := &generator.GenerateRequest{
		ComponentType: generator.ComponentRepository,
		ModelName:     "OrderItem",
		Fields: []generator.FieldDefinition{
			{
				Name:    "OrderID",
				Type:    "uint",
				GormTag: "not null;index",
				JsonTag: "order_id",
				Comment: "订单ID",
			},
			{
				Name:    "ProductID",
				Type:    "uint",
				GormTag: "not null;index",
				JsonTag: "product_id",
				Comment: "产品ID",
			},
			{
				Name:    "Quantity",
				Type:    "int",
				GormTag: "not null;check:quantity > 0",
				JsonTag: "quantity",
				Comment: "数量",
			},
		},
		TableName:   "order_items",
		PackagePath: "internal/admin",
		Options: generator.GenerateOptions{
			OverwriteExisting: true,
			FormatCode:        true,
		},
	}

	response, err := repositoryGenerator.Generate(req)
	helper.AssertNoError(err)

	if !response.Success {
		t.Errorf("生成失败: %v", response.Errors)
		return
	}

	if len(response.GeneratedFiles) > 0 {
		filePath := response.GeneratedFiles[0]

		// 验证模型名称转换
		helper.AssertFileContains(filePath, "type OrderItemRepository interface")
		helper.AssertFileContains(filePath, "type OrderItemRepositoryImpl struct")
		helper.AssertFileContains(filePath, "func NewOrderItemRepository")

		// 验证泛型使用
		helper.AssertFileContains(filePath, "BaseRepositoryInterface[models.OrderItem]")
		helper.AssertFileContains(filePath, "BaseRepository[models.OrderItem]")

		// 验证方法签名
		helper.AssertFileContains(filePath, "ListWithFilter(query types.BasePageQuery)")
	}
}

func TestRepositoryGenerator_ErrorHandling(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	repositoryGenerator := generator.NewRepositoryGenerator()

	t.Run("无效请求处理", func(t *testing.T) {
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentRepository,
			ModelName:     "", // 空模型名
			Fields: []generator.FieldDefinition{
				{Name: "Name", Type: "string"},
			},
		}

		response, err := repositoryGenerator.Generate(req)
		helper.AssertNoError(err)

		if response.Success {
			t.Error("期望生成失败，但生成成功了")
		}

		if len(response.Errors) == 0 {
			t.Error("期望有错误信息，但没有错误")
		}
	})

	t.Run("空字段列表处理", func(t *testing.T) {
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentRepository,
			ModelName:     "EmptyFieldsModel",
			Fields:        []generator.FieldDefinition{}, // 空字段列表
			TableName:     "empty_fields_models",
			PackagePath:   "internal/admin",
		}

		response, err := repositoryGenerator.Generate(req)
		helper.AssertNoError(err)

		if response.Success {
			t.Error("期望生成失败，但生成成功了")
		}

		// 检查错误信息
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
