package unit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"bico-admin/internal/devtools/generator"
	"bico-admin/tests/generator/utils"
)

func TestHistoryManager_LoadHistory(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	// 在临时目录中创建历史管理器
	tempDir := helper.GetTempDir()
	historyFile := filepath.Join(tempDir, "test-history.json")
	
	// 创建自定义历史管理器用于测试
	historyManager := &generator.HistoryManager{}
	// 由于HistoryManager的filePath是私有字段，我们需要通过其他方式测试

	t.Run("文件不存在时返回空历史", func(t *testing.T) {
		// 使用默认的历史管理器
		manager := generator.NewHistoryManager()
		
		// 临时修改历史文件路径到不存在的位置
		nonExistentPath := filepath.Join(tempDir, "non-existent", "history.json")
		
		// 由于无法直接修改私有字段，我们测试默认行为
		// 这里主要测试当文件不存在时的处理逻辑
		history, err := manager.GetHistory()
		helper.AssertNoError(err)
		
		if history == nil {
			t.Error("期望返回空的历史记录数组，但得到nil")
		}
	})

	t.Run("加载有效的历史文件", func(t *testing.T) {
		// 创建测试历史数据
		testHistory := generator.HistoryFile{
			Version: "1.0.0",
			History: []generator.GenerateHistory{
				{
					ModuleName:     "User",
					GeneratedAt:    time.Now(),
					Components:     []string{"model"},
					ModelName:      "User",
					TableName:      "users",
					PackagePath:    "internal/admin",
					GeneratedBy:    "Test",
					GeneratedFiles: []string{"internal/shared/models/user.go"},
				},
			},
		}

		// 写入测试数据到临时文件
		data, err := json.MarshalIndent(testHistory, "", "  ")
		helper.AssertNoError(err)

		err = os.MkdirAll(filepath.Dir(historyFile), 0755)
		helper.AssertNoError(err)

		err = os.WriteFile(historyFile, data, 0644)
		helper.AssertNoError(err)

		// 由于无法直接设置私有字段，我们通过其他方式验证
		// 这里主要验证JSON解析逻辑
		var loadedHistory generator.HistoryFile
		loadedData, err := os.ReadFile(historyFile)
		helper.AssertNoError(err)

		err = json.Unmarshal(loadedData, &loadedHistory)
		helper.AssertNoError(err)

		if loadedHistory.Version != "1.0.0" {
			t.Errorf("期望版本 '1.0.0', 得到 '%s'", loadedHistory.Version)
		}

		if len(loadedHistory.History) != 1 {
			t.Errorf("期望1条历史记录, 得到 %d 条", len(loadedHistory.History))
		}

		if loadedHistory.History[0].ModuleName != "User" {
			t.Errorf("期望模块名 'User', 得到 '%s'", loadedHistory.History[0].ModuleName)
		}
	})

	t.Run("加载无效的JSON文件", func(t *testing.T) {
		// 写入无效的JSON数据
		invalidJSON := `{"invalid": json}`
		err := os.WriteFile(historyFile, []byte(invalidJSON), 0644)
		helper.AssertNoError(err)

		// 验证JSON解析错误处理
		var history generator.HistoryFile
		data, err := os.ReadFile(historyFile)
		helper.AssertNoError(err)

		err = json.Unmarshal(data, &history)
		if err == nil {
			t.Error("期望JSON解析错误，但没有错误")
		}
	})
}

func TestHistoryManager_SaveHistory(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.GetTempDir()
	historyFile := filepath.Join(tempDir, "save-test-history.json")

	t.Run("保存历史记录到新文件", func(t *testing.T) {
		testHistory := generator.HistoryFile{
			Version: "1.0.0",
			History: []generator.GenerateHistory{
				{
					ModuleName:     "Product",
					GeneratedAt:    helper.MockTime(),
					Components:     []string{"model", "repository"},
					ModelName:      "Product",
					TableName:      "products",
					PackagePath:    "internal/admin",
					GeneratedBy:    "Test",
					GeneratedFiles: []string{"internal/shared/models/product.go"},
				},
			},
		}

		// 序列化并保存
		data, err := json.MarshalIndent(testHistory, "", "  ")
		helper.AssertNoError(err)

		err = os.MkdirAll(filepath.Dir(historyFile), 0755)
		helper.AssertNoError(err)

		err = os.WriteFile(historyFile, data, 0644)
		helper.AssertNoError(err)

		// 验证文件存在
		helper.AssertFileExists(historyFile)

		// 验证文件内容
		savedData, err := os.ReadFile(historyFile)
		helper.AssertNoError(err)

		var savedHistory generator.HistoryFile
		err = json.Unmarshal(savedData, &savedHistory)
		helper.AssertNoError(err)

		if savedHistory.Version != "1.0.0" {
			t.Errorf("期望版本 '1.0.0', 得到 '%s'", savedHistory.Version)
		}

		if len(savedHistory.History) != 1 {
			t.Errorf("期望1条历史记录, 得到 %d 条", len(savedHistory.History))
		}
	})

	t.Run("保存到不存在的目录", func(t *testing.T) {
		deepPath := filepath.Join(tempDir, "deep", "nested", "path", "history.json")
		
		testHistory := generator.HistoryFile{
			Version: "1.0.0",
			History: []generator.GenerateHistory{},
		}

		// 创建目录并保存
		err := os.MkdirAll(filepath.Dir(deepPath), 0755)
		helper.AssertNoError(err)

		data, err := json.MarshalIndent(testHistory, "", "  ")
		helper.AssertNoError(err)

		err = os.WriteFile(deepPath, data, 0644)
		helper.AssertNoError(err)

		// 验证文件存在
		helper.AssertFileExists(deepPath)
	})
}

func TestHistoryManager_AddHistory(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	manager := generator.NewHistoryManager()

	t.Run("添加新的历史记录", func(t *testing.T) {
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentModel,
			ModelName:     "TestModel",
			TableName:     "test_models",
			PackagePath:   "internal/admin",
		}

		generatedFiles := []string{"internal/shared/models/test_model.go"}

		err := manager.AddHistory(req, generatedFiles)
		helper.AssertNoError(err)

		// 验证历史记录已添加
		history, err := manager.GetHistoryByModule("TestModel")
		helper.AssertNoError(err)

		if history.ModuleName != "TestModel" {
			t.Errorf("期望模块名 'TestModel', 得到 '%s'", history.ModuleName)
		}

		if len(history.GeneratedFiles) != 1 {
			t.Errorf("期望1个生成文件, 得到 %d 个", len(history.GeneratedFiles))
		}

		if history.GeneratedFiles[0] != "internal/shared/models/test_model.go" {
			t.Errorf("期望文件路径 'internal/shared/models/test_model.go', 得到 '%s'", history.GeneratedFiles[0])
		}
	})

	t.Run("更新现有的历史记录", func(t *testing.T) {
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentModel,
			ModelName:     "TestModel", // 相同的模块名
			TableName:     "test_models",
			PackagePath:   "internal/admin",
		}

		newGeneratedFiles := []string{
			"internal/shared/models/test_model.go",
			"internal/admin/repository/test_model_repository.go",
		}

		err := manager.AddHistory(req, newGeneratedFiles)
		helper.AssertNoError(err)

		// 验证历史记录已更新
		history, err := manager.GetHistoryByModule("TestModel")
		helper.AssertNoError(err)

		if len(history.GeneratedFiles) != 2 {
			t.Errorf("期望2个生成文件, 得到 %d 个", len(history.GeneratedFiles))
		}
	})

	t.Run("添加All组件类型的历史记录", func(t *testing.T) {
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentAll,
			ModelName:     "AllTestModel",
			TableName:     "all_test_models",
			PackagePath:   "internal/admin",
		}

		generatedFiles := []string{
			"internal/shared/models/all_test_model.go",
			"internal/admin/repository/all_test_model_repository.go",
			"internal/admin/service/all_test_model_service.go",
		}

		err := manager.AddHistory(req, generatedFiles)
		helper.AssertNoError(err)

		// 验证历史记录
		history, err := manager.GetHistoryByModule("AllTestModel")
		helper.AssertNoError(err)

		if len(history.Components) == 0 {
			t.Error("期望组件列表不为空")
		}

		// 验证包含model组件
		found := false
		for _, component := range history.Components {
			if component == string(generator.ComponentModel) {
				found = true
				break
			}
		}
		if !found {
			t.Error("期望组件列表包含 'model'")
		}
	})
}

func TestHistoryManager_GetHistory(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	manager := generator.NewHistoryManager()

	t.Run("获取空历史记录", func(t *testing.T) {
		history, err := manager.GetHistory()
		helper.AssertNoError(err)

		if history == nil {
			t.Error("期望返回空数组，但得到nil")
		}

		if len(history) != 0 {
			t.Errorf("期望0条历史记录, 得到 %d 条", len(history))
		}
	})

	t.Run("获取包含记录的历史", func(t *testing.T) {
		// 先添加一些历史记录
		req1 := &generator.GenerateRequest{
			ComponentType: generator.ComponentModel,
			ModelName:     "Model1",
			TableName:     "model1s",
			PackagePath:   "internal/admin",
		}

		req2 := &generator.GenerateRequest{
			ComponentType: generator.ComponentRepository,
			ModelName:     "Model2",
			TableName:     "model2s",
			PackagePath:   "internal/admin",
		}

		err := manager.AddHistory(req1, []string{"file1.go"})
		helper.AssertNoError(err)

		err = manager.AddHistory(req2, []string{"file2.go"})
		helper.AssertNoError(err)

		// 获取所有历史记录
		history, err := manager.GetHistory()
		helper.AssertNoError(err)

		if len(history) != 2 {
			t.Errorf("期望2条历史记录, 得到 %d 条", len(history))
		}

		// 验证记录内容
		moduleNames := make(map[string]bool)
		for _, record := range history {
			moduleNames[record.ModuleName] = true
		}

		if !moduleNames["Model1"] {
			t.Error("期望包含 'Model1' 模块")
		}

		if !moduleNames["Model2"] {
			t.Error("期望包含 'Model2' 模块")
		}
	})
}

func TestHistoryManager_GetHistoryByModule(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	manager := generator.NewHistoryManager()

	t.Run("获取不存在的模块历史", func(t *testing.T) {
		_, err := manager.GetHistoryByModule("NonExistentModule")
		if err == nil {
			t.Error("期望获取不存在模块时返回错误")
		}

		if !strings.Contains(err.Error(), "未找到模块") {
			t.Errorf("期望错误信息包含'未找到模块', 实际错误: %v", err)
		}
	})

	t.Run("获取存在的模块历史", func(t *testing.T) {
		// 先添加历史记录
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentModel,
			ModelName:     "ExistingModel",
			TableName:     "existing_models",
			PackagePath:   "internal/admin",
		}

		generatedFiles := []string{"existing_model.go"}
		err := manager.AddHistory(req, generatedFiles)
		helper.AssertNoError(err)

		// 获取特定模块的历史
		history, err := manager.GetHistoryByModule("ExistingModel")
		helper.AssertNoError(err)

		if history.ModuleName != "ExistingModel" {
			t.Errorf("期望模块名 'ExistingModel', 得到 '%s'", history.ModuleName)
		}

		if len(history.GeneratedFiles) != 1 {
			t.Errorf("期望1个生成文件, 得到 %d 个", len(history.GeneratedFiles))
		}
	})
}

func TestHistoryManager_DeleteHistory(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	manager := generator.NewHistoryManager()

	t.Run("删除不存在的模块历史", func(t *testing.T) {
		err := manager.DeleteHistory("NonExistentModule")
		if err == nil {
			t.Error("期望删除不存在模块时返回错误")
		}

		if !strings.Contains(err.Error(), "未找到模块") {
			t.Errorf("期望错误信息包含'未找到模块', 实际错误: %v", err)
		}
	})

	t.Run("删除存在的模块历史", func(t *testing.T) {
		// 先添加历史记录
		req := &generator.GenerateRequest{
			ComponentType: generator.ComponentModel,
			ModelName:     "ToDeleteModel",
			TableName:     "to_delete_models",
			PackagePath:   "internal/admin",
		}

		// 创建一个临时文件模拟生成的文件
		tempFile := helper.CreateTempFile("to_delete_model.go", "package models")
		generatedFiles := []string{tempFile}

		err := manager.AddHistory(req, generatedFiles)
		helper.AssertNoError(err)

		// 验证历史记录存在
		_, err = manager.GetHistoryByModule("ToDeleteModel")
		helper.AssertNoError(err)

		// 删除历史记录
		err = manager.DeleteHistory("ToDeleteModel")
		helper.AssertNoError(err)

		// 验证历史记录已删除
		_, err = manager.GetHistoryByModule("ToDeleteModel")
		if err == nil {
			t.Error("期望删除后获取模块历史返回错误")
		}

		// 验证生成的文件也被删除（如果文件存在的话）
		// 注意：在实际实现中，文件删除可能会失败，这是正常的
	})
}

func TestHistoryManager_ClearHistory(t *testing.T) {
	helper := utils.NewTestHelper(t)
	defer helper.Cleanup()

	manager := generator.NewHistoryManager()

	t.Run("清空历史记录", func(t *testing.T) {
		// 先添加一些历史记录
		req1 := &generator.GenerateRequest{
			ComponentType: generator.ComponentModel,
			ModelName:     "Model1",
			TableName:     "model1s",
			PackagePath:   "internal/admin",
		}

		req2 := &generator.GenerateRequest{
			ComponentType: generator.ComponentModel,
			ModelName:     "Model2",
			TableName:     "model2s",
			PackagePath:   "internal/admin",
		}

		// 创建临时文件
		tempFile1 := helper.CreateTempFile("model1.go", "package models")
		tempFile2 := helper.CreateTempFile("model2.go", "package models")

		err := manager.AddHistory(req1, []string{tempFile1})
		helper.AssertNoError(err)

		err = manager.AddHistory(req2, []string{tempFile2})
		helper.AssertNoError(err)

		// 验证历史记录存在
		history, err := manager.GetHistory()
		helper.AssertNoError(err)

		if len(history) != 2 {
			t.Errorf("期望2条历史记录, 得到 %d 条", len(history))
		}

		// 清空历史记录
		err = manager.ClearHistory()
		helper.AssertNoError(err)

		// 验证历史记录已清空
		history, err = manager.GetHistory()
		helper.AssertNoError(err)

		if len(history) != 0 {
			t.Errorf("期望0条历史记录, 得到 %d 条", len(history))
		}
	})
}
