package unit

import (
	"path/filepath"
	"strings"
	"testing"

	"bico-admin/internal/devtools/generator"
	"bico-admin/tests/generator/utils"
)

func TestHandlerGenerator_Generate(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	handlerGenerator := generator.NewHandlerGenerator()

	tests := []struct {
		name          string
		request       *generator.GenerateRequest
		expectSuccess bool
		expectedFiles int
	}{
		{
			name: "基础Handler生成",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentHandler,
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
			expectedFiles: 2, // handler.go 和 types.go
		},
		{
			name: "包含时间字段的Handler生成",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentHandler,
				ModelName:     "Article",
				Fields: []generator.FieldDefinition{
					{
						Name:    "Title",
						Type:    "string",
						GormTag: "size:200;not null",
						JsonTag: "title",
						Comment: "标题",
					},
					{
						Name:    "Content",
						Type:    "string",
						GormTag: "type:text",
						JsonTag: "content",
						Comment: "内容",
					},
					{
						Name:    "Status",
						Type:    "int",
						GormTag: "default:0",
						JsonTag: "status",
						Comment: "状态",
					},
					{
						Name:    "PublishedAt",
						Type:    "*time.Time",
						JsonTag: "published_at",
						Comment: "发布时间",
					},
				},
				TableName:   "articles",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
				},
			},
			expectSuccess: true,
			expectedFiles: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := handlerGenerator.Generate(tt.request)
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

					t.Run("验证Handler文件内容_"+filepath.Base(filePath), func(t *testing.T) {
						if strings.HasSuffix(filePath, "_handler.go") {
							validateHandlerFileContent(t, helper, filePath, tt.request)
						} else if strings.HasSuffix(filePath, "_types.go") {
							validateHandlerTypesFileContent(t, helper, filePath, tt.request)
						}
					})
				}
			}
		})
	}
}

func validateHandlerFileContent(t *testing.T, helper *utils.TestHelper, filePath string, req *generator.GenerateRequest) {
	// 验证包声明
	helper.AssertFileContains(filePath, "package handler")

	// 验证导入
	helper.AssertFileContains(filePath, `"github.com/gin-gonic/gin"`)
	helper.AssertFileContains(filePath, `"bico-admin/internal/shared/models"`)
	helper.AssertFileContains(filePath, `"bico-admin/internal/shared/service"`)
	helper.AssertFileContains(filePath, `"bico-admin/internal/admin/types"`)
	helper.AssertFileContains(filePath, `"bico-admin/pkg/utils"`)

	// 检查是否有时间字段
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

	// 验证处理器结构体
	handlerName := req.ModelName + "Handler"
	helper.AssertFileContains(filePath, "type "+handlerName+" struct {")
	helper.AssertFileContains(filePath, "*BaseHandler[models."+req.ModelName+", types."+req.ModelName+"CreateRequest, types."+req.ModelName+"UpdateRequest, types."+req.ModelName+"ListRequest, types."+req.ModelName+"Response]")

	// 验证构造函数
	constructorName := "New" + handlerName
	helper.AssertFileContains(filePath, "func "+constructorName+"(")

	// 验证选项设置
	helper.AssertFileContains(filePath, "options := DefaultHandlerOptions()")
	helper.AssertFileContains(filePath, "options.EnableSoftDelete = true")

	// 检查是否有状态字段
	hasStatusField := generator.HasStatusField(req.Fields)
	if hasStatusField {
		helper.AssertFileContains(filePath, "options.EnableStatusManagement = true")
	} else {
		helper.AssertFileContains(filePath, "options.EnableStatusManagement = false")
	}

	// 验证转换方法
	helper.AssertFileContains(filePath, "func (h *"+handlerName+") ConvertToResponse")
	helper.AssertFileContains(filePath, "func (h *"+handlerName+") ConvertCreateRequest")
	helper.AssertFileContains(filePath, "func (h *"+handlerName+") ConvertUpdateRequest")
	helper.AssertFileContains(filePath, "func (h *"+handlerName+") ConvertListRequest")

	// 验证状态处理方法
	if hasStatusField {
		helper.AssertFileContains(filePath, "func (h *"+handlerName+") getStatusValue")
		helper.AssertFileContains(filePath, "func (h *"+handlerName+") getStatusText")
	}

	// 验证时间格式化方法
	if hasTimeField {
		helper.AssertFileContains(filePath, "func (h *"+handlerName+") formatTime")
	}
}

func validateHandlerTypesFileContent(t *testing.T, helper *utils.TestHelper, filePath string, req *generator.GenerateRequest) {
	// 验证包声明
	helper.AssertFileContains(filePath, "package types")

	// 验证导入
	helper.AssertFileContains(filePath, `"bico-admin/internal/shared/types"`)

	// 检查是否有时间字段
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

	// 验证请求/响应结构体
	helper.AssertFileContains(filePath, "type "+req.ModelName+"CreateRequest struct {")
	helper.AssertFileContains(filePath, "type "+req.ModelName+"UpdateRequest struct {")
	helper.AssertFileContains(filePath, "type "+req.ModelName+"ListRequest struct {")
	helper.AssertFileContains(filePath, "type "+req.ModelName+"Response struct {")

	// 验证列表请求继承基础分页
	helper.AssertFileContains(filePath, "types.BasePageQuery")

	// 验证响应包含基础字段
	helper.AssertFileContains(filePath, "ID uint `json:\"id\"`")
	helper.AssertFileContains(filePath, "CreatedAt string `json:\"created_at\"`")
	helper.AssertFileContains(filePath, "UpdatedAt string `json:\"updated_at\"`")

	// 验证字段定义
	for _, field := range req.Fields {
		// 检查创建请求字段
		if field.Name != "Status" { // 状态字段通常不在创建请求中
			helper.AssertFileContains(filePath, field.Name+" "+getTypeForRequest(field.Type))
		}

		// 检查响应字段
		if field.Name == "Status" {
			helper.AssertFileContains(filePath, "Status int `json:\"status\"`")
			helper.AssertFileContains(filePath, "StatusText string `json:\"status_text\"`")
		} else if strings.Contains(field.Type, "time.Time") {
			helper.AssertFileContains(filePath, field.Name+" string `json:\""+field.JsonTag+"\"`")
		} else {
			helper.AssertFileContains(filePath, field.Name+" "+field.Type+" `json:\""+field.JsonTag+"\"`")
		}
	}
}

// 辅助函数：获取请求中的字段类型
func getTypeForRequest(fieldType string) string {
	if strings.Contains(fieldType, "time.Time") {
		return "string"
	}
	return fieldType
}

func TestHandlerGenerator_DataConversionMethods(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	handlerGenerator := generator.NewHandlerGenerator()

	req := &generator.GenerateRequest{
		ComponentType: generator.ComponentHandler,
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
				Name:    "Price",
				Type:    "float64",
				GormTag: "type:decimal(10,2);not null",
				JsonTag: "price",
				Comment: "价格",
			},
			{
				Name:    "Status",
				Type:    "int",
				GormTag: "default:1",
				JsonTag: "status",
				Comment: "状态",
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
		},
	}

	response, err := handlerGenerator.Generate(req)
	helper.AssertNoError(err)

	if !response.Success {
		t.Errorf("生成失败: %v", response.Errors)
		return
	}

	// 查找handler文件
	var handlerFilePath string
	for _, filePath := range response.GeneratedFiles {
		if strings.HasSuffix(filePath, "_handler.go") {
			handlerFilePath = filePath
			break
		}
	}

	if handlerFilePath == "" {
		t.Error("未找到生成的handler文件")
		return
	}

	// 验证ConvertToResponse方法
	helper.AssertFileContains(handlerFilePath, "func (h *ProductHandler) ConvertToResponse")
	helper.AssertFileContains(handlerFilePath, "return types.ProductResponse{")
	helper.AssertFileContains(handlerFilePath, "ID: entity.ID,")
	helper.AssertFileContains(handlerFilePath, "Name: entity.Name,")
	helper.AssertFileContains(handlerFilePath, "Price: entity.Price,")
	helper.AssertFileContains(handlerFilePath, "Status: h.getStatusValue(entity.Status),")
	helper.AssertFileContains(handlerFilePath, "StatusText: h.getStatusText(h.getStatusValue(entity.Status)),")
	helper.AssertFileContains(handlerFilePath, "PublishedAt: h.formatTime(entity.PublishedAt),")
	helper.AssertFileContains(handlerFilePath, "CreatedAt: utils.NewFormattedTime(entity.CreatedAt),")
	helper.AssertFileContains(handlerFilePath, "UpdatedAt: utils.NewFormattedTime(entity.UpdatedAt),")

	// 验证ConvertCreateRequest方法
	helper.AssertFileContains(handlerFilePath, "func (h *ProductHandler) ConvertCreateRequest")
	helper.AssertFileContains(handlerFilePath, "return &models.Product{")
	helper.AssertFileContains(handlerFilePath, "Name: req.Name,")
	helper.AssertFileContains(handlerFilePath, "Price: req.Price,")

	// 验证ConvertUpdateRequest方法
	helper.AssertFileContains(handlerFilePath, "func (h *ProductHandler) ConvertUpdateRequest")
	helper.AssertFileContains(handlerFilePath, "entity.Name = req.Name")
	helper.AssertFileContains(handlerFilePath, "entity.Price = req.Price")

	// 验证ConvertListRequest方法
	helper.AssertFileContains(handlerFilePath, "func (h *ProductHandler) ConvertListRequest")
	helper.AssertFileContains(handlerFilePath, "return types.BasePageQuery{")
	helper.AssertFileContains(handlerFilePath, "Page: req.Page,")
	helper.AssertFileContains(handlerFilePath, "PageSize: req.PageSize,")
	helper.AssertFileContains(handlerFilePath, "Keyword: req.Keyword,")
	helper.AssertFileContains(handlerFilePath, "OrderBy: req.OrderBy,")
	helper.AssertFileContains(handlerFilePath, "OrderDirection: req.OrderDirection,")

	// 验证辅助方法
	helper.AssertFileContains(handlerFilePath, "func (h *ProductHandler) getStatusValue")
	helper.AssertFileContains(handlerFilePath, "func (h *ProductHandler) getStatusText")
	helper.AssertFileContains(handlerFilePath, "func (h *ProductHandler) formatTime")
}

func TestHandlerGenerator_ErrorHandling(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	handlerGenerator := generator.NewHandlerGenerator()

	t.Run("无效模型名称", func(t *testing.T) {
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentHandler,
			ModelName:     "123Invalid",
			Fields: []generator.FieldDefinition{
				{Name: "Name", Type: "string"},
			},
		}

		response, err := handlerGenerator.Generate(req)
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
			ComponentType: generator.ComponentHandler,
			ModelName:     "EmptyHandler",
			Fields:        []generator.FieldDefinition{},
			TableName:     "empty_handlers",
			PackagePath:   "internal/admin",
		}

		response, err := handlerGenerator.Generate(req)
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
