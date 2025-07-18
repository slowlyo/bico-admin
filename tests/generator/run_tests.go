package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"bico-admin/internal/devtools/generator"
)

// TestRunner 测试运行器
type TestRunner struct {
	codeGenerator *generator.CodeGenerator
	testResults   []TestResult
}

// TestResult 测试结果
type TestResult struct {
	Name           string        `json:"name"`
	Success        bool          `json:"success"`
	Duration       time.Duration `json:"duration"`
	Errors         []string      `json:"errors,omitempty"`
	GeneratedFiles int           `json:"generated_files"`
	Description    string        `json:"description,omitempty"`
}

// TestSuite 测试套件
type TestSuite struct {
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Tests       []generator.GenerateRequest `json:"tests"`
}

func main() {
	fmt.Println("🚀 开始运行代码生成器测试套件...")

	runner := &TestRunner{
		codeGenerator: generator.NewCodeGenerator(),
		testResults:   make([]TestResult, 0),
	}

	// 运行基础功能测试
	runner.runBasicTests()

	// 运行业务场景测试
	runner.runBusinessScenarioTests()

	// 运行边界情况测试
	runner.runEdgeCaseTests()

	// 运行字段类型测试
	runner.runFieldTypeTests()

	// 运行编译测试
	runner.runCompilationTests()

	// 生成测试报告
	runner.generateReport()

	fmt.Println("✅ 测试套件运行完成！")
}

func (r *TestRunner) runBasicTests() {
	fmt.Println("\n📋 运行基础功能测试...")

	tests := []struct {
		name        string
		description string
		request     *generator.GenerateRequest
	}{
		{
			name:        "基础模型生成",
			description: "测试基础模型的生成功能",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentModel,
				ModelName:     "BasicTestModel",
				Fields: []generator.FieldDefinition{
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
						Comment:  "状态",
					},
				},
				TableName:   "basic_test_models",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
				},
			},
		},
		{
			name:        "Repository生成",
			description: "测试Repository层的生成功能",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentRepository,
				ModelName:     "RepoTestModel",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string", JsonTag: "name"},
					{Name: "Status", Type: "int", JsonTag: "status"},
				},
				TableName:   "repo_test_models",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
				},
			},
		},
		{
			name:        "Service生成",
			description: "测试Service层的生成功能",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentService,
				ModelName:     "ServiceTestModel",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string", JsonTag: "name"},
					{Name: "Status", Type: "int", JsonTag: "status"},
				},
				TableName:   "service_test_models",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
				},
			},
		},
		{
			name:        "Handler生成",
			description: "测试Handler层的生成功能",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentHandler,
				ModelName:     "HandlerTestModel",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string", JsonTag: "name"},
					{Name: "Status", Type: "int", JsonTag: "status"},
				},
				TableName:   "handler_test_models",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
				},
			},
		},
		{
			name:        "完整模块生成",
			description: "测试完整模块的生成功能",
			request: &generator.GenerateRequest{
				ComponentType: generator.ComponentAll,
				ModelName:     "CompleteTestModel",
				Fields: []generator.FieldDefinition{
					{Name: "Name", Type: "string", JsonTag: "name"},
					{Name: "Status", Type: "int", JsonTag: "status"},
					{Name: "CreatedTime", Type: "*time.Time", JsonTag: "created_time"},
				},
				TableName:   "complete_test_models",
				PackagePath: "internal/admin",
				Options: generator.GenerateOptions{
					OverwriteExisting: true,
					FormatCode:        true,
				},
			},
		},
	}

	for _, test := range tests {
		r.runSingleTest(test.name, test.description, test.request)
	}
}

func (r *TestRunner) runBusinessScenarioTests() {
	fmt.Println("\n🏢 运行业务场景测试...")

	// 加载业务场景测试数据
	scenariosFile := "tests/generator/testdata/requests/business_scenarios.json"
	data, err := os.ReadFile(scenariosFile)
	if err != nil {
		log.Printf("⚠️  无法加载业务场景测试数据: %v", err)
		return
	}

	var scenarios []struct {
		Name        string                    `json:"name"`
		Description string                    `json:"description"`
		Request     generator.GenerateRequest `json:"request"`
	}

	if err := json.Unmarshal(data, &scenarios); err != nil {
		log.Printf("⚠️  解析业务场景测试数据失败: %v", err)
		return
	}

	for _, scenario := range scenarios {
		r.runSingleTest(scenario.Name, scenario.Description, &scenario.Request)
	}
}

func (r *TestRunner) runEdgeCaseTests() {
	fmt.Println("\n⚠️  运行边界情况测试...")

	// 加载边界情况测试数据
	edgeCasesFile := "tests/generator/testdata/requests/edge_cases.json"
	data, err := os.ReadFile(edgeCasesFile)
	if err != nil {
		log.Printf("⚠️  无法加载边界情况测试数据: %v", err)
		return
	}

	var edgeCases []struct {
		Name          string                    `json:"name"`
		Description   string                    `json:"description"`
		Request       generator.GenerateRequest `json:"request"`
		ExpectedError string                    `json:"expected_error"`
	}

	if err := json.Unmarshal(data, &edgeCases); err != nil {
		log.Printf("⚠️  解析边界情况测试数据失败: %v", err)
		return
	}

	for _, edgeCase := range edgeCases {
		// 特殊处理字段数量超限测试
		if edgeCase.Name == "too_many_fields" {
			// 动态生成51个字段
			fields := make([]generator.FieldDefinition, 51)
			for i := 0; i < 51; i++ {
				fields[i] = generator.FieldDefinition{
					Name:    fmt.Sprintf("Field%d", i+1),
					Type:    "string",
					JsonTag: fmt.Sprintf("field_%d", i+1),
					Comment: fmt.Sprintf("字段%d", i+1),
				}
			}
			edgeCase.Request.Fields = fields
		}

		// 对于边界情况，我们期望生成失败
		result := r.runSingleTestWithExpectedFailure(edgeCase.Name, edgeCase.Description, &edgeCase.Request, edgeCase.ExpectedError)
		r.testResults = append(r.testResults, result)
	}
}

func (r *TestRunner) runFieldTypeTests() {
	fmt.Println("\n🔧 运行字段类型测试...")

	// 加载字段类型测试数据
	fieldTypesFile := "tests/generator/testdata/requests/field_types_comprehensive.json"
	data, err := os.ReadFile(fieldTypesFile)
	if err != nil {
		log.Printf("⚠️  无法加载字段类型测试数据: %v", err)
		return
	}

	var request generator.GenerateRequest
	if err := json.Unmarshal(data, &request); err != nil {
		log.Printf("⚠️  解析字段类型测试数据失败: %v", err)
		return
	}

	r.runSingleTest("综合字段类型测试", "测试所有支持的字段类型", &request)
}

func (r *TestRunner) runSingleTest(name, description string, request *generator.GenerateRequest) {
	fmt.Printf("  🧪 运行测试: %s\n", name)

	startTime := time.Now()
	response, err := r.codeGenerator.Generate(request)
	duration := time.Since(startTime)

	result := TestResult{
		Name:        name,
		Duration:    duration,
		Description: description,
	}

	if err != nil {
		result.Success = false
		result.Errors = []string{err.Error()}
		fmt.Printf("    ❌ 测试失败: %v\n", err)
	} else if !response.Success {
		result.Success = false
		result.Errors = response.Errors
		fmt.Printf("    ❌ 生成失败: %v\n", response.Errors)
	} else {
		result.Success = true
		result.GeneratedFiles = len(response.GeneratedFiles)
		fmt.Printf("    ✅ 测试通过 (生成了 %d 个文件)\n", len(response.GeneratedFiles))

		// 验证生成的文件
		for _, filePath := range response.GeneratedFiles {
			if err := r.validateGeneratedFile(filePath); err != nil {
				result.Success = false
				result.Errors = append(result.Errors, fmt.Sprintf("文件验证失败 %s: %v", filePath, err))
			}
		}
	}

	r.testResults = append(r.testResults, result)
}

func (r *TestRunner) runSingleTestWithExpectedFailure(name, description string, request *generator.GenerateRequest, expectedError string) TestResult {
	fmt.Printf("  🧪 运行边界测试: %s\n", name)

	startTime := time.Now()
	response, err := r.codeGenerator.Generate(request)
	duration := time.Since(startTime)

	result := TestResult{
		Name:        name,
		Duration:    duration,
		Description: description,
	}

	// 对于边界情况，我们期望生成失败
	if err != nil || !response.Success {
		// 检查错误信息是否包含期望的错误
		errorFound := false
		if err != nil && strings.Contains(err.Error(), expectedError) {
			errorFound = true
		}
		if !errorFound && len(response.Errors) > 0 {
			for _, errMsg := range response.Errors {
				if strings.Contains(errMsg, expectedError) {
					errorFound = true
					break
				}
			}
		}

		if errorFound {
			result.Success = true
			fmt.Printf("    ✅ 边界测试通过 (正确捕获了期望的错误)\n")
		} else {
			result.Success = false
			result.Errors = []string{fmt.Sprintf("未找到期望的错误信息: %s", expectedError)}
			fmt.Printf("    ❌ 边界测试失败 (未找到期望的错误信息)\n")
		}
	} else {
		// 期望失败但成功了
		result.Success = false
		result.Errors = []string{"期望生成失败但生成成功了"}
		fmt.Printf("    ❌ 边界测试失败 (期望失败但生成成功了)\n")
	}

	return result
}

func (r *TestRunner) validateGeneratedFile(filePath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在")
	}

	// 检查文件是否为空
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		return fmt.Errorf("文件为空")
	}

	// 检查Go文件的基本语法
	if strings.HasSuffix(filePath, ".go") {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		contentStr := string(content)

		// 基本的语法检查
		if !strings.Contains(contentStr, "package ") {
			return fmt.Errorf("缺少package声明")
		}

		// 检查是否有明显的语法错误
		if strings.Contains(contentStr, "{{") || strings.Contains(contentStr, "}}") {
			return fmt.Errorf("包含未处理的模板标记")
		}
	}

	return nil
}

func (r *TestRunner) generateReport() {
	fmt.Println("\n📊 生成测试报告...")

	// 统计结果
	totalTests := len(r.testResults)
	successTests := 0
	totalDuration := time.Duration(0)

	for _, result := range r.testResults {
		if result.Success {
			successTests++
		}
		totalDuration += result.Duration
	}

	// 控制台报告
	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Printf("📈 测试报告摘要\n")
	fmt.Printf(strings.Repeat("=", 60) + "\n")
	fmt.Printf("总测试数: %d\n", totalTests)
	fmt.Printf("成功测试: %d\n", successTests)
	fmt.Printf("失败测试: %d\n", totalTests-successTests)
	fmt.Printf("成功率: %.2f%%\n", float64(successTests)/float64(totalTests)*100)
	fmt.Printf("总耗时: %v\n", totalDuration)
	fmt.Printf("平均耗时: %v\n", totalDuration/time.Duration(totalTests))
	fmt.Printf(strings.Repeat("=", 60) + "\n")

	// 详细结果
	fmt.Println("\n📋 详细测试结果:")
	for _, result := range r.testResults {
		status := "✅"
		if !result.Success {
			status = "❌"
		}
		fmt.Printf("  %s %s (耗时: %v)\n", status, result.Name, result.Duration)
		if !result.Success && len(result.Errors) > 0 {
			for _, err := range result.Errors {
				fmt.Printf("      错误: %s\n", err)
			}
		}
	}

	// 生成JSON报告
	reportFile := "test_report.json"
	report := map[string]interface{}{
		"summary": map[string]interface{}{
			"total_tests":      totalTests,
			"success_tests":    successTests,
			"failed_tests":     totalTests - successTests,
			"success_rate":     float64(successTests) / float64(totalTests) * 100,
			"total_duration":   totalDuration.String(),
			"average_duration": (totalDuration / time.Duration(totalTests)).String(),
			"generated_at":     time.Now().Format(time.RFC3339),
		},
		"results": r.testResults,
	}

	reportData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Printf("⚠️  生成JSON报告失败: %v", err)
		return
	}

	if err := os.WriteFile(reportFile, reportData, 0644); err != nil {
		log.Printf("⚠️  保存测试报告失败: %v", err)
		return
	}

	fmt.Printf("\n💾 测试报告已保存到: %s\n", reportFile)
}

func (r *TestRunner) runCompilationTests() {
	fmt.Println("\n🔨 运行编译测试...")

	// 测试生成的代码是否能正常编译
	compilationTests := []struct {
		name        string
		description string
		buildPath   string
	}{
		{
			name:        "模型编译测试",
			description: "验证生成的模型文件能否正常编译",
			buildPath:   "./internal/shared/models/...",
		},
		{
			name:        "Repository编译测试",
			description: "验证生成的Repository文件能否正常编译",
			buildPath:   "./internal/admin/repository/...",
		},
		{
			name:        "Service编译测试",
			description: "验证生成的Service文件能否正常编译",
			buildPath:   "./internal/admin/service/...",
		},
		{
			name:        "Handler编译测试",
			description: "验证生成的Handler文件能否正常编译",
			buildPath:   "./internal/admin/handler/...",
		},
		{
			name:        "Routes编译测试",
			description: "验证生成的Routes文件能否正常编译",
			buildPath:   "./internal/admin/routes/...",
		},
		{
			name:        "Wire编译测试",
			description: "验证生成的Wire文件能否正常编译",
			buildPath:   "./internal/admin/wire/...",
		},
		{
			name:        "Permission编译测试",
			description: "验证生成的Permission文件能否正常编译",
			buildPath:   "./internal/admin/definitions/...",
		},
	}

	for _, test := range compilationTests {
		r.runCompilationTest(test.name, test.description, test.buildPath)
	}
}

func (r *TestRunner) runCompilationTest(name, description, buildPath string) {
	fmt.Printf("  🔨 编译测试: %s\n", name)

	startTime := time.Now()

	// 执行 go build 命令
	cmd := exec.Command("go", "build", buildPath)
	cmd.Dir = "." // 设置工作目录

	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime)

	result := TestResult{
		Name:        name,
		Duration:    duration,
		Description: description,
	}

	if err != nil {
		result.Success = false
		result.Errors = []string{
			fmt.Sprintf("编译失败: %v", err),
			fmt.Sprintf("编译输出: %s", string(output)),
		}
		fmt.Printf("    ❌ 编译失败: %v\n", err)
		if len(output) > 0 {
			fmt.Printf("    编译输出: %s\n", string(output))
		}
	} else {
		result.Success = true
		fmt.Printf("    ✅ 编译通过\n")
	}

	r.testResults = append(r.testResults, result)
}
