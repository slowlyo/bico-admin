package compilation

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bico-admin/internal/devtools/generator"
	"bico-admin/tests/generator/utils"
)

func TestGeneratedCode_SyntaxValidation(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	codeGenerator := generator.NewCodeGenerator()

	tests := []struct {
		name    string
		request *generator.GenerateRequest
	}{
		{
			name: "基础模型语法验证",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "SyntaxTestModel",
				Fields: []generator.FieldDefinition{
					{
						Name:     "Name",
						Type:     "string",
						GormTag:  "size:100;not null",
						JsonTag:  "name",
						Comment:  "名称",
					},
					{
						Name:     "Age",
						Type:     "int",
						GormTag:  "check:age >= 0",
						JsonTag:  "age",
						Comment:  "年龄",
					},
					{
						Name:     "Email",
						Type:     "string",
						GormTag:  "uniqueIndex;size:100",
						JsonTag:  "email",
						Comment:  "邮箱",
					},
					{
						Name:     "IsActive",
						Type:     "bool",
						GormTag:  "default:true",
						JsonTag:  "is_active",
						Comment:  "是否激活",
					},
					{
						Name:     "CreatedTime",
						Type:     "*time.Time",
						JsonTag:  "created_time",
						Comment:  "创建时间",
					},
				},
				TableName:   "syntax_test_models",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
					OptimizeImports:   true,
				},
			},
		},
		{
			name: "复杂Repository语法验证",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentRepository,
				ModelName:     "ComplexRepo",
				Fields: []generator.FieldDefinition{
					{
						Name:     "Title",
						Type:     "string",
						GormTag:  "size:200;not null;index",
						JsonTag:  "title",
						Comment:  "标题",
					},
					{
						Name:     "Content",
						Type:     "string",
						GormTag:  "type:text",
						JsonTag:  "content",
						Comment:  "内容",
					},
					{
						Name:     "AuthorID",
						Type:     "uint",
						GormTag:  "not null;index",
						JsonTag:  "author_id",
						Comment:  "作者ID",
					},
					{
						Name:     "Tags",
						Type:     "string",
						GormTag:  "type:json",
						JsonTag:  "tags",
						Comment:  "标签JSON",
					},
					{
						Name:     "PublishedAt",
						Type:     "*time.Time",
						JsonTag:  "published_at",
						Comment:  "发布时间",
					},
				},
				TableName:   "complex_repos",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
				},
			},
		},
		{
			name: "Service和Handler语法验证",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentService,
				ModelName:     "ServiceTest",
				Fields: []generator.FieldDefinition{
					{
						Name:     "Name",
						Type:     "string",
						GormTag:  "size:100;not null",
						JsonTag:  "name",
						Comment:  "名称",
					},
					{
						Name:     "Status",
						Type:     "int",
						GormTag:  "default:1;index",
						JsonTag:  "status",
						Comment:  "状态",
					},
					{
						Name:     "Price",
						Type:     "float64",
						GormTag:  "type:decimal(10,2)",
						JsonTag:  "price",
						Comment:  "价格",
					},
				},
				TableName:   "service_tests",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 生成代码
			response, err := codeGenerator.Generate(tt.request)
			helper.AssertNoError(err)

			if !response.Success {
				t.Errorf("代码生成失败: %v", response.Errors)
				return
			}

			// 验证每个生成文件的语法
			for _, filePath := range response.GeneratedFiles {
				t.Run("语法检查_"+filepath.Base(filePath), func(t *testing.T) {
					validateGoFileSyntax(t, helper, filePath)
				})
			}
		})
	}
}

func validateGoFileSyntax(t *testing.T, helper *utils.TestHelper, filePath string) {
	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("读取文件失败: %v", err)
	}

	// 解析Go代码
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		t.Fatalf("Go语法解析失败: %v", err)
	}

	// 验证AST结构
	validateASTStructure(t, file, filePath)
}

func validateASTStructure(t *testing.T, file *ast.File, filePath string) {
	// 验证包声明
	if file.Name == nil {
		t.Errorf("文件 %s 缺少包声明", filePath)
		return
	}

	// 验证导入声明
	hasImports := len(file.Imports) > 0
	if !hasImports && needsImports(filePath) {
		t.Errorf("文件 %s 应该有导入声明但没有", filePath)
	}

	// 验证声明
	if len(file.Decls) == 0 {
		t.Errorf("文件 %s 没有任何声明", filePath)
		return
	}

	// 验证具体的声明类型
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			validateGenDecl(t, d, filePath)
		case *ast.FuncDecl:
			validateFuncDecl(t, d, filePath)
		default:
			t.Errorf("文件 %s 包含未知的声明类型: %T", filePath, decl)
		}
	}
}

func validateGenDecl(t *testing.T, decl *ast.GenDecl, filePath string) {
	switch decl.Tok {
	case token.IMPORT:
		// 验证导入声明
		for _, spec := range decl.Specs {
			if importSpec, ok := spec.(*ast.ImportSpec); ok {
				if importSpec.Path == nil || importSpec.Path.Value == "" {
					t.Errorf("文件 %s 包含无效的导入路径", filePath)
				}
			}
		}
	case token.TYPE:
		// 验证类型声明
		for _, spec := range decl.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				if typeSpec.Name == nil || typeSpec.Name.Name == "" {
					t.Errorf("文件 %s 包含无效的类型名称", filePath)
				}
				validateTypeSpec(t, typeSpec, filePath)
			}
		}
	case token.CONST, token.VAR:
		// 验证常量和变量声明
		for _, spec := range decl.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				if len(valueSpec.Names) == 0 {
					t.Errorf("文件 %s 包含无效的值声明", filePath)
				}
			}
		}
	}
}

func validateTypeSpec(t *testing.T, typeSpec *ast.TypeSpec, filePath string) {
	switch typeExpr := typeSpec.Type.(type) {
	case *ast.StructType:
		// 验证结构体类型
		if typeExpr.Fields == nil {
			t.Errorf("文件 %s 中的结构体 %s 没有字段列表", filePath, typeSpec.Name.Name)
			return
		}

		for _, field := range typeExpr.Fields.List {
			if field.Type == nil {
				t.Errorf("文件 %s 中的结构体 %s 包含无效的字段类型", filePath, typeSpec.Name.Name)
			}
			
			// 验证字段标签
			if field.Tag != nil {
				validateStructTag(t, field.Tag.Value, filePath, typeSpec.Name.Name)
			}
		}
	case *ast.InterfaceType:
		// 验证接口类型
		if typeExpr.Methods == nil {
			t.Errorf("文件 %s 中的接口 %s 没有方法列表", filePath, typeSpec.Name.Name)
		}
	}
}

func validateStructTag(t *testing.T, tag, filePath, structName string) {
	// 移除反引号
	tag = strings.Trim(tag, "`")
	
	// 验证常见的标签格式
	if strings.Contains(tag, "json:") {
		if !strings.Contains(tag, `json:"`) {
			t.Errorf("文件 %s 中结构体 %s 的JSON标签格式无效: %s", filePath, structName, tag)
		}
	}
	
	if strings.Contains(tag, "gorm:") {
		if !strings.Contains(tag, `gorm:"`) {
			t.Errorf("文件 %s 中结构体 %s 的GORM标签格式无效: %s", filePath, structName, tag)
		}
	}
}

func validateFuncDecl(t *testing.T, decl *ast.FuncDecl, filePath string) {
	if decl.Name == nil || decl.Name.Name == "" {
		t.Errorf("文件 %s 包含无效的函数名称", filePath)
		return
	}

	// 验证函数类型
	if decl.Type == nil {
		t.Errorf("文件 %s 中的函数 %s 没有类型信息", filePath, decl.Name.Name)
		return
	}

	// 验证参数列表
	if decl.Type.Params != nil {
		for _, param := range decl.Type.Params.List {
			if param.Type == nil {
				t.Errorf("文件 %s 中的函数 %s 包含无效的参数类型", filePath, decl.Name.Name)
			}
		}
	}

	// 验证返回值列表
	if decl.Type.Results != nil {
		for _, result := range decl.Type.Results.List {
			if result.Type == nil {
				t.Errorf("文件 %s 中的函数 %s 包含无效的返回值类型", filePath, decl.Name.Name)
			}
		}
	}

	// 验证函数体（如果存在）
	if decl.Body != nil {
		validateBlockStmt(t, decl.Body, filePath, decl.Name.Name)
	}
}

func validateBlockStmt(t *testing.T, block *ast.BlockStmt, filePath, funcName string) {
	if block.List == nil {
		return // 空函数体是允许的
	}

	// 验证语句列表
	for _, stmt := range block.List {
		if stmt == nil {
			t.Errorf("文件 %s 中的函数 %s 包含空语句", filePath, funcName)
		}
	}
}

func needsImports(filePath string) bool {
	// 判断文件是否应该有导入声明
	// 大多数生成的Go文件都需要导入其他包
	return strings.HasSuffix(filePath, ".go") && 
		   !strings.Contains(filePath, "_test.go")
}

func TestGeneratedCode_SpecificSyntaxPatterns(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	codeGenerator := generator.NewCodeGenerator()

	t.Run("泛型语法验证", func(t *testing.T) {
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentRepository,
			ModelName:     "GenericTest",
			Fields: []generator.FieldDefinition{
				{Name: "Name", Type: "string", JsonTag: "name"},
			},
			TableName:   "generic_tests",
			PackagePath: "internal/admin",
			Options: generator.GenerateOptions{
				OverwriteExisting: true,
				FormatCode:        true,
			},
		}

		response, err := codeGenerator.Generate(req)
		helper.AssertNoError(err)

		if !response.Success {
			t.Errorf("生成失败: %v", response.Errors)
			return
		}

		// 验证泛型语法
		for _, filePath := range response.GeneratedFiles {
			content, err := os.ReadFile(filePath)
			helper.AssertNoError(err)

			contentStr := string(content)
			
			// 检查泛型语法模式
			if strings.Contains(contentStr, "BaseRepositoryInterface[") {
				if !strings.Contains(contentStr, "BaseRepositoryInterface[models.GenericTest]") {
					t.Errorf("文件 %s 中的泛型语法可能不正确", filePath)
				}
			}

			if strings.Contains(contentStr, "BaseRepository[") {
				if !strings.Contains(contentStr, "BaseRepository[models.GenericTest]") {
					t.Errorf("文件 %s 中的泛型语法可能不正确", filePath)
				}
			}
		}
	})

	t.Run("接口嵌入语法验证", func(t *testing.T) {
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentService,
			ModelName:     "InterfaceTest",
			Fields: []generator.FieldDefinition{
				{Name: "Name", Type: "string", JsonTag: "name"},
				{Name: "Status", Type: "int", JsonTag: "status"},
			},
			TableName:   "interface_tests",
			PackagePath: "internal/admin",
			Options: generator.GenerateOptions{
				OverwriteExisting: true,
				FormatCode:        true,
			},
		}

		response, err := codeGenerator.Generate(req)
		helper.AssertNoError(err)

		if !response.Success {
			t.Errorf("生成失败: %v", response.Errors)
			return
		}

		// 验证接口嵌入语法
		for _, filePath := range response.GeneratedFiles {
			if strings.HasSuffix(filePath, "_service.go") {
				content, err := os.ReadFile(filePath)
				helper.AssertNoError(err)

				contentStr := string(content)
				
				// 检查接口嵌入语法
				if strings.Contains(contentStr, "type InterfaceTestService interface") {
					if !strings.Contains(contentStr, "service.BaseServiceInterface[") {
						t.Errorf("文件 %s 中的接口嵌入语法可能不正确", filePath)
					}
				}
			}
		}
	})

	t.Run("方法接收者语法验证", func(t *testing.T) {
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentHandler,
			ModelName:     "ReceiverTest",
			Fields: []generator.FieldDefinition{
				{Name: "Name", Type: "string", JsonTag: "name"},
				{Name: "Status", Type: "int", JsonTag: "status"},
			},
			TableName:   "receiver_tests",
			PackagePath: "internal/admin",
			Options: generator.GenerateOptions{
				OverwriteExisting: true,
				FormatCode:        true,
			},
		}

		response, err := codeGenerator.Generate(req)
		helper.AssertNoError(err)

		if !response.Success {
			t.Errorf("生成失败: %v", response.Errors)
			return
		}

		// 验证方法接收者语法
		for _, filePath := range response.GeneratedFiles {
			if strings.HasSuffix(filePath, "_handler.go") {
				content, err := os.ReadFile(filePath)
				helper.AssertNoError(err)

				contentStr := string(content)
				
				// 检查方法接收者语法
				lines := strings.Split(contentStr, "\n")
				for _, line := range lines {
					if strings.Contains(line, "func (h *ReceiverTestHandler)") {
						// 验证接收者语法正确
						if !strings.Contains(line, ") ") {
							t.Errorf("文件 %s 中的方法接收者语法可能不正确: %s", filePath, line)
						}
					}
				}
			}
		}
	})
}
