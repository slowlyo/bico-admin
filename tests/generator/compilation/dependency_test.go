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

func TestGeneratedCode_DependencyValidation(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	codeGenerator := generator.NewCodeGenerator()

	tests := []struct {
		name             string
		request          *generator.GenerateRequest
		expectedImports  map[string][]string // 文件类型 -> 期望的导入包
		forbiddenImports map[string][]string // 文件类型 -> 禁止的导入包
	}{
		{
			name: "模型文件依赖验证",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "DependencyTestModel",
				Fields: []generator.FieldDefinition{
					{
						Name:     "Name",
						Type:     "string",
						JsonTag:  "name",
						Comment:  "名称",
					},
					{
						Name:     "CreatedAt",
						Type:     "*time.Time",
						JsonTag:  "created_at",
						Comment:  "创建时间",
					},
				},
				TableName:   "dependency_test_models",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
				},
			},
			expectedImports: map[string][]string{
				"model": {
					"time",
					"bico-admin/internal/shared/types",
				},
			},
			forbiddenImports: map[string][]string{
				"model": {
					"github.com/gin-gonic/gin",
					"gorm.io/gorm",
				},
			},
		},
		{
			name: "Repository文件依赖验证",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentRepository,
				ModelName:     "RepoDepTest",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string", JsonTag: "name"},
					{Name: "Status", Type: "int", JsonTag: "status"},
				},
				TableName:   "repo_dep_tests",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
				},
			},
			expectedImports: map[string][]string{
				"repository": {
					"bico-admin/internal/shared/models",
					"bico-admin/internal/shared/repository",
					"bico-admin/internal/shared/types",
					"gorm.io/gorm",
				},
			},
			forbiddenImports: map[string][]string{
				"repository": {
					"github.com/gin-gonic/gin",
					"time", // 如果没有时间字段，不应该导入time包
				},
			},
		},
		{
			name: "Service文件依赖验证",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentService,
				ModelName:     "ServiceDepTest",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string", JsonTag: "name"},
					{Name: "Status", Type: "int", JsonTag: "status"},
				},
				TableName:   "service_dep_tests",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
				},
			},
			expectedImports: map[string][]string{
				"service": {
					"context",
					"bico-admin/internal/shared/models",
					"bico-admin/internal/shared/service",
					"bico-admin/internal/shared/types",
					"bico-admin/internal/admin/repository",
				},
			},
			forbiddenImports: map[string][]string{
				"service": {
					"github.com/gin-gonic/gin",
				},
			},
		},
		{
			name: "Handler文件依赖验证",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentHandler,
				ModelName:     "HandlerDepTest",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string", JsonTag: "name"},
					{Name: "Status", Type: "int", JsonTag: "status"},
					{Name: "UpdatedAt", Type: "*time.Time", JsonTag: "updated_at"},
				},
				TableName:   "handler_dep_tests",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
				},
			},
			expectedImports: map[string][]string{
				"handler": {
					"time",
					"github.com/gin-gonic/gin",
					"bico-admin/internal/shared/models",
					"bico-admin/internal/shared/service",
					"bico-admin/internal/admin/types",
					"bico-admin/pkg/utils",
				},
				"types": {
					"bico-admin/internal/shared/types",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := codeGenerator.Generate(tt.request)
			helper.AssertNoError(err)

			if !response.Success {
				t.Errorf("代码生成失败: %v", response.Errors)
				return
			}

			// 验证每个生成文件的依赖
			for _, filePath := range response.GeneratedFiles {
				fileType := getFileType(filePath)
				
				t.Run("依赖检查_"+filepath.Base(filePath), func(t *testing.T) {
					validateFileDependencies(t, helper, filePath, fileType, tt.expectedImports, tt.forbiddenImports)
				})
			}
		})
	}
}

func getFileType(filePath string) string {
	fileName := filepath.Base(filePath)
	
	if strings.Contains(fileName, "_repository.go") {
		return "repository"
	} else if strings.Contains(fileName, "_service.go") {
		return "service"
	} else if strings.Contains(fileName, "_handler.go") {
		return "handler"
	} else if strings.Contains(fileName, "_types.go") {
		return "types"
	} else if strings.Contains(fileName, "_routes") {
		return "routes"
	} else if strings.Contains(fileName, "_wire") {
		return "wire"
	} else if strings.HasSuffix(fileName, ".go") && strings.Contains(filePath, "/models/") {
		return "model"
	}
	
	return "unknown"
}

func validateFileDependencies(t *testing.T, helper *utils.TestHelper, filePath, fileType string, expectedImports, forbiddenImports map[string][]string) {
	// 解析文件获取导入信息
	imports, err := extractImports(filePath)
	if err != nil {
		t.Errorf("提取导入信息失败: %v", err)
		return
	}

	// 检查期望的导入
	if expectedList, exists := expectedImports[fileType]; exists {
		for _, expectedImport := range expectedList {
			if !containsImport(imports, expectedImport) {
				t.Errorf("文件 %s 缺少期望的导入: %s", filePath, expectedImport)
			}
		}
	}

	// 检查禁止的导入
	if forbiddenList, exists := forbiddenImports[fileType]; exists {
		for _, forbiddenImport := range forbiddenList {
			if containsImport(imports, forbiddenImport) {
				t.Errorf("文件 %s 包含禁止的导入: %s", filePath, forbiddenImport)
			}
		}
	}

	// 验证导入路径的有效性
	for _, importPath := range imports {
		validateImportPath(t, importPath, filePath)
	}
}

func extractImports(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var imports []string
	for _, importSpec := range file.Imports {
		if importSpec.Path != nil {
			// 移除引号
			importPath := strings.Trim(importSpec.Path.Value, `"`)
			imports = append(imports, importPath)
		}
	}

	return imports, nil
}

func containsImport(imports []string, target string) bool {
	for _, imp := range imports {
		if imp == target {
			return true
		}
	}
	return false
}

func validateImportPath(t *testing.T, importPath, filePath string) {
	// 验证导入路径格式
	if importPath == "" {
		t.Errorf("文件 %s 包含空的导入路径", filePath)
		return
	}

	// 验证标准库导入
	if isStandardLibrary(importPath) {
		return // 标准库导入总是有效的
	}

	// 验证项目内部导入
	if strings.HasPrefix(importPath, "bico-admin/") {
		validateInternalImport(t, importPath, filePath)
		return
	}

	// 验证第三方导入
	if strings.Contains(importPath, ".") {
		validateThirdPartyImport(t, importPath, filePath)
		return
	}

	t.Errorf("文件 %s 包含无效的导入路径格式: %s", filePath, importPath)
}

func isStandardLibrary(importPath string) bool {
	standardLibs := []string{
		"context", "fmt", "time", "strings", "strconv", "errors",
		"encoding/json", "net/http", "os", "path/filepath",
		"go/ast", "go/parser", "go/token", "testing",
	}

	for _, lib := range standardLibs {
		if importPath == lib {
			return true
		}
	}

	// 检查标准库包前缀
	standardPrefixes := []string{
		"encoding/", "net/", "crypto/", "database/", "go/",
		"html/", "image/", "mime/", "text/", "unicode/",
	}

	for _, prefix := range standardPrefixes {
		if strings.HasPrefix(importPath, prefix) {
			return true
		}
	}

	return false
}

func validateInternalImport(t *testing.T, importPath, filePath string) {
	// 验证内部导入路径的合理性
	validInternalPaths := []string{
		"bico-admin/internal/shared/models",
		"bico-admin/internal/shared/types",
		"bico-admin/internal/shared/service",
		"bico-admin/internal/shared/repository",
		"bico-admin/internal/admin/repository",
		"bico-admin/internal/admin/service",
		"bico-admin/internal/admin/handler",
		"bico-admin/internal/admin/types",
		"bico-admin/pkg/utils",
		"bico-admin/pkg/response",
		"bico-admin/pkg/logger",
	}

	for _, validPath := range validInternalPaths {
		if importPath == validPath {
			return
		}
	}

	// 允许一些动态的内部路径
	if strings.HasPrefix(importPath, "bico-admin/internal/") ||
		strings.HasPrefix(importPath, "bico-admin/pkg/") {
		return
	}

	t.Errorf("文件 %s 包含可疑的内部导入路径: %s", filePath, importPath)
}

func validateThirdPartyImport(t *testing.T, importPath, filePath string) {
	// 验证第三方导入路径
	validThirdPartyPaths := []string{
		"github.com/gin-gonic/gin",
		"gorm.io/gorm",
		"gorm.io/driver/sqlite",
		"gorm.io/driver/mysql",
		"gorm.io/driver/postgres",
		"github.com/go-playground/validator/v10",
		"github.com/google/wire",
	}

	for _, validPath := range validThirdPartyPaths {
		if importPath == validPath {
			return
		}
	}

	// 允许一些常见的第三方包前缀
	validPrefixes := []string{
		"github.com/",
		"golang.org/x/",
		"google.golang.org/",
		"gorm.io/",
	}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(importPath, prefix) {
			return
		}
	}

	t.Errorf("文件 %s 包含未知的第三方导入路径: %s", filePath, importPath)
}

func TestGeneratedCode_CircularDependency(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	codeGenerator := generator.NewCodeGenerator()

	t.Run("检查循环依赖", func(t *testing.T) {
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentAll,
			ModelName:     "CircularTest",
			Fields: []generator.FieldDefinition{
				{Name: "Name", Type: "string", JsonTag: "name"},
				{Name: "Status", Type: "int", JsonTag: "status"},
			},
			TableName:   "circular_tests",
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

		// 构建依赖图
		dependencyGraph := make(map[string][]string)
		for _, filePath := range response.GeneratedFiles {
			imports, err := extractImports(filePath)
			if err != nil {
				t.Errorf("提取导入失败: %v", err)
				continue
			}

			// 只关注项目内部的依赖
			var internalImports []string
			for _, imp := range imports {
				if strings.HasPrefix(imp, "bico-admin/internal/") {
					internalImports = append(internalImports, imp)
				}
			}

			dependencyGraph[filePath] = internalImports
		}

		// 检查循环依赖
		visited := make(map[string]bool)
		recStack := make(map[string]bool)

		for filePath := range dependencyGraph {
			if !visited[filePath] {
				if hasCircularDependency(filePath, dependencyGraph, visited, recStack) {
					t.Errorf("检测到循环依赖，涉及文件: %s", filePath)
				}
			}
		}
	})
}

func hasCircularDependency(filePath string, graph map[string][]string, visited, recStack map[string]bool) bool {
	visited[filePath] = true
	recStack[filePath] = true

	for _, dep := range graph[filePath] {
		// 将导入路径转换为文件路径进行比较
		if !visited[dep] {
			if hasCircularDependency(dep, graph, visited, recStack) {
				return true
			}
		} else if recStack[dep] {
			return true
		}
	}

	recStack[filePath] = false
	return false
}

func TestGeneratedCode_UnusedImports(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	codeGenerator := generator.NewCodeGenerator()

	t.Run("检查未使用的导入", func(t *testing.T) {
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentModel,
			ModelName:     "UnusedImportTest",
			Fields: []generator.FieldDefinition{
				{Name: "Name", Type: "string", JsonTag: "name"},
				// 没有时间字段，不应该导入time包
			},
			TableName:   "unused_import_tests",
			PackagePath: "internal/admin",
			Options: generator.GenerateOptions{
				OverwriteExisting: true,
				FormatCode:        true,
				OptimizeImports:   true, // 启用导入优化
			},
		}

		response, err := codeGenerator.Generate(req)
		helper.AssertNoError(err)

		if !response.Success {
			t.Errorf("生成失败: %v", response.Errors)
			return
		}

		// 检查生成的文件是否有未使用的导入
		for _, filePath := range response.GeneratedFiles {
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Errorf("读取文件失败: %v", err)
				continue
			}

			contentStr := string(content)
			
			// 如果没有时间字段，不应该导入time包
			if !strings.Contains(contentStr, "time.Time") && strings.Contains(contentStr, `"time"`) {
				t.Errorf("文件 %s 可能包含未使用的time包导入", filePath)
			}

			// 检查其他可能未使用的导入
			checkUnusedImport(t, contentStr, filePath, "context", "context.")
			checkUnusedImport(t, contentStr, filePath, "fmt", "fmt.")
			checkUnusedImport(t, contentStr, filePath, "strings", "strings.")
		}
	})
}

func checkUnusedImport(t *testing.T, content, filePath, importName, usage string) {
	hasImport := strings.Contains(content, `"`+importName+`"`)
	hasUsage := strings.Contains(content, usage)
	
	if hasImport && !hasUsage {
		t.Errorf("文件 %s 可能包含未使用的 %s 包导入", filePath, importName)
	}
}
