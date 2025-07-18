package utils

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"bico-admin/internal/devtools/generator"
)

// TestHelper 测试辅助工具
type TestHelper struct {
	t       *testing.T
	tempDir string
}

// NewTestHelper 创建测试辅助工具
func NewTestHelper(t *testing.T) *TestHelper {
	tempDir, err := os.MkdirTemp("", "generator_test_*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}

	return &TestHelper{
		t:       t,
		tempDir: tempDir,
	}
}

// Cleanup 清理测试资源
func (h *TestHelper) Cleanup() {
	if err := os.RemoveAll(h.tempDir); err != nil {
		h.t.Logf("清理临时目录失败: %v", err)
	}
}

// GetTempDir 获取临时目录
func (h *TestHelper) GetTempDir() string {
	return h.tempDir
}

// CreateTestRequest 创建测试请求
func (h *TestHelper) CreateTestRequest(componentType generator.ComponentType, modelName string, fields []generator.FieldDefinition) *generator.GenerateRequest {
	return &generator.GenerateRequest{
		ComponentType: componentType,
		ModelName:     modelName,
		Fields:        fields,
		TableName:     generator.ToSnakeCase(modelName) + "s",
		PackagePath:   "internal/admin",
		Options: generator.GenerateOptions{
			OverwriteExisting: true,
			FormatCode:        true,
			OptimizeImports:   true,
		},
	}
}

// CreateBasicFields 创建基础字段定义
func (h *TestHelper) CreateBasicFields() []generator.FieldDefinition {
	return []generator.FieldDefinition{
		{
			Name:     "Name",
			Type:     "string",
			GormTag:  "size:100;not null",
			JsonTag:  "name",
			Validate: "required,max=100",
			Comment:  "名称",
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
			Name:     "Description",
			Type:     "string",
			GormTag:  "type:text",
			JsonTag:  "description",
			Validate: "",
			Comment:  "描述",
		},
	}
}

// CreateComplexFields 创建复杂字段定义
func (h *TestHelper) CreateComplexFields() []generator.FieldDefinition {
	return []generator.FieldDefinition{
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
			Name:     "Age",
			Type:     "int",
			GormTag:  "check:age >= 0 AND age <= 150",
			JsonTag:  "age",
			Validate: "min=0,max=150",
			Comment:  "年龄",
		},
		{
			Name:     "Balance",
			Type:     "float64",
			GormTag:  "type:decimal(10,2);default:0.00",
			JsonTag:  "balance",
			Validate: "min=0",
			Comment:  "余额",
		},
		{
			Name:     "IsActive",
			Type:     "bool",
			GormTag:  "default:true",
			JsonTag:  "is_active",
			Validate: "",
			Comment:  "是否激活",
		},
		{
			Name:     "LastLoginAt",
			Type:     "*time.Time",
			GormTag:  "",
			JsonTag:  "last_login_at",
			Validate: "",
			Comment:  "最后登录时间",
		},
		{
			Name:     "Profile",
			Type:     "string",
			GormTag:  "type:json",
			JsonTag:  "profile",
			Validate: "",
			Comment:  "用户资料JSON",
		},
	}
}

// LoadTestData 加载测试数据
func (h *TestHelper) LoadTestData(filename string, target interface{}) error {
	data, err := os.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		return fmt.Errorf("读取测试数据文件失败: %w", err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("解析测试数据失败: %w", err)
	}

	return nil
}

// SaveTestData 保存测试数据
func (h *TestHelper) SaveTestData(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化测试数据失败: %w", err)
	}

	testDataDir := "testdata"
	if err := os.MkdirAll(testDataDir, 0755); err != nil {
		return fmt.Errorf("创建测试数据目录失败: %w", err)
	}

	filePath := filepath.Join(testDataDir, filename)
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("写入测试数据文件失败: %w", err)
	}

	return nil
}

// ValidateGoSyntax 验证Go代码语法
func (h *TestHelper) ValidateGoSyntax(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	fset := token.NewFileSet()
	_, err = parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("Go语法解析失败: %w", err)
	}

	return nil
}

// ValidateGoFile 验证Go文件结构
func (h *TestHelper) ValidateGoFile(filePath string) (*ast.File, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("Go语法解析失败: %w", err)
	}

	return file, nil
}

// AssertFileExists 断言文件存在
func (h *TestHelper) AssertFileExists(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		h.t.Errorf("期望文件存在但不存在: %s", filePath)
	}
}

// AssertFileNotExists 断言文件不存在
func (h *TestHelper) AssertFileNotExists(filePath string) {
	if _, err := os.Stat(filePath); err == nil {
		h.t.Errorf("期望文件不存在但存在: %s", filePath)
	}
}

// AssertFileContains 断言文件包含指定内容
func (h *TestHelper) AssertFileContains(filePath, expectedContent string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		h.t.Errorf("读取文件失败: %v", err)
		return
	}

	if !strings.Contains(string(content), expectedContent) {
		h.t.Errorf("文件 %s 不包含期望内容: %s", filePath, expectedContent)
	}
}

// AssertFileNotContains 断言文件不包含指定内容
func (h *TestHelper) AssertFileNotContains(filePath, unexpectedContent string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		h.t.Errorf("读取文件失败: %v", err)
		return
	}

	if strings.Contains(string(content), unexpectedContent) {
		h.t.Errorf("文件 %s 包含不期望的内容: %s", filePath, unexpectedContent)
	}
}

// CreateTempFile 创建临时文件
func (h *TestHelper) CreateTempFile(name, content string) string {
	filePath := filepath.Join(h.tempDir, name)
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		h.t.Fatalf("创建目录失败: %v", err)
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		h.t.Fatalf("创建临时文件失败: %v", err)
	}

	return filePath
}

// CompareFiles 比较两个文件内容
func (h *TestHelper) CompareFiles(file1, file2 string) bool {
	content1, err1 := os.ReadFile(file1)
	content2, err2 := os.ReadFile(file2)

	if err1 != nil || err2 != nil {
		return false
	}

	return string(content1) == string(content2)
}

// WalkGeneratedFiles 遍历生成的文件
func (h *TestHelper) WalkGeneratedFiles(rootDir string, callback func(path string, info fs.FileInfo) error) error {
	return filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".go") {
			info, err := d.Info()
			if err != nil {
				return err
			}
			return callback(path, info)
		}

		return nil
	})
}

// MockTime 模拟时间
func (h *TestHelper) MockTime() time.Time {
	return time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
}

// AssertNoError 断言无错误
func (h *TestHelper) AssertNoError(err error) {
	if err != nil {
		h.t.Errorf("期望无错误但得到错误: %v", err)
	}
}

// AssertError 断言有错误
func (h *TestHelper) AssertError(err error) {
	if err == nil {
		h.t.Error("期望有错误但没有错误")
	}
}

// AssertErrorContains 断言错误包含指定消息
func (h *TestHelper) AssertErrorContains(err error, expectedMsg string) {
	if err == nil {
		h.t.Error("期望有错误但没有错误")
		return
	}

	if !strings.Contains(err.Error(), expectedMsg) {
		h.t.Errorf("错误消息不包含期望内容。期望包含: %s, 实际错误: %v", expectedMsg, err)
	}
}
