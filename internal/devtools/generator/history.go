package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	HistoryFileName = "data/code-generate-history.json"
	HistoryVersion  = "1.0.0"
)

// HistoryManager 历史记录管理器
type HistoryManager struct {
	filePath string
}

// NewHistoryManager 创建历史记录管理器
func NewHistoryManager() *HistoryManager {
	return &HistoryManager{
		filePath: HistoryFileName,
	}
}

// LoadHistory 加载历史记录
func (h *HistoryManager) LoadHistory() (*HistoryFile, error) {
	// 确保目录存在
	dir := filepath.Dir(h.filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("创建目录失败: %w", err)
		}
	}

	// 如果文件不存在，返回空的历史记录
	if _, err := os.Stat(h.filePath); os.IsNotExist(err) {
		return &HistoryFile{
			Version: HistoryVersion,
			History: []GenerateHistory{},
		}, nil
	}

	// 读取文件
	data, err := os.ReadFile(h.filePath)
	if err != nil {
		return nil, fmt.Errorf("读取历史文件失败: %w", err)
	}

	// 解析JSON
	var historyFile HistoryFile
	if err := json.Unmarshal(data, &historyFile); err != nil {
		return nil, fmt.Errorf("解析历史文件失败: %w", err)
	}

	return &historyFile, nil
}

// SaveHistory 保存历史记录
func (h *HistoryManager) SaveHistory(historyFile *HistoryFile) error {
	// 确保目录存在
	dir := filepath.Dir(h.filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
	}

	// 序列化为JSON
	data, err := json.MarshalIndent(historyFile, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化历史文件失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(h.filePath, data, 0644); err != nil {
		return fmt.Errorf("写入历史文件失败: %w", err)
	}

	return nil
}

// AddHistory 添加历史记录
func (h *HistoryManager) AddHistory(req *GenerateRequest, generatedFiles []string) error {
	// 加载现有历史
	historyFile, err := h.LoadHistory()
	if err != nil {
		return err
	}

	// 创建新的历史记录
	history := GenerateHistory{
		ModuleName:     req.ModelName,
		GeneratedAt:    time.Now(),
		Components:     []string{string(req.ComponentType)},
		ModelName:      req.ModelName,
		TableName:      req.TableName,
		PackagePath:    req.PackagePath,
		GeneratedBy:    "MCP Code Generator",
		GeneratedFiles: generatedFiles, // 存储生成的文件路径
	}

	// 如果是生成所有组件，记录具体组件列表
	if req.ComponentType == ComponentAll {
		history.Components = []string{
			string(ComponentModel),
			// TODO: 后续添加其他组件
		}
	}

	// 检查是否已存在相同模块的记录
	found := false
	for i, existing := range historyFile.History {
		if existing.ModuleName == req.ModelName && existing.PackagePath == req.PackagePath {
			// 更新现有记录
			historyFile.History[i] = history
			found = true
			break
		}
	}

	// 如果没有找到，添加新记录
	if !found {
		historyFile.History = append(historyFile.History, history)
	}

	// 保存历史记录
	return h.SaveHistory(historyFile)
}

// GetHistory 获取历史记录
func (h *HistoryManager) GetHistory() ([]GenerateHistory, error) {
	historyFile, err := h.LoadHistory()
	if err != nil {
		return nil, err
	}

	return historyFile.History, nil
}

// GetHistoryByModule 根据模块名获取历史记录
func (h *HistoryManager) GetHistoryByModule(moduleName string) (*GenerateHistory, error) {
	historyFile, err := h.LoadHistory()
	if err != nil {
		return nil, err
	}

	for _, history := range historyFile.History {
		if history.ModuleName == moduleName {
			return &history, nil
		}
	}

	return nil, fmt.Errorf("未找到模块'%s'的历史记录", moduleName)
}

// DeleteHistory 删除历史记录
func (h *HistoryManager) DeleteHistory(moduleName string) error {
	historyFile, err := h.LoadHistory()
	if err != nil {
		return err
	}

	// 查找并删除记录
	for i, history := range historyFile.History {
		if history.ModuleName == moduleName {
			// 删除生成的文件
			for _, filePath := range history.GeneratedFiles {
				if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
					fmt.Printf("警告: 删除文件'%s'失败: %v\n", filePath, err)
				}
			}

			// 删除记录
			historyFile.History = append(historyFile.History[:i], historyFile.History[i+1:]...)
			return h.SaveHistory(historyFile)
		}
	}

	return fmt.Errorf("未找到模块'%s'的历史记录", moduleName)
}

// ClearHistory 清空历史记录
func (h *HistoryManager) ClearHistory() error {
	// 先加载现有历史，删除所有文件
	historyFile, err := h.LoadHistory()
	if err != nil {
		return err
	}

	// 删除所有生成的文件
	for _, history := range historyFile.History {
		for _, filePath := range history.GeneratedFiles {
			if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
				fmt.Printf("警告: 删除文件'%s'失败: %v\n", filePath, err)
			}
		}
	}

	// 清空历史记录
	historyFile = &HistoryFile{
		Version: HistoryVersion,
		History: []GenerateHistory{},
	}

	return h.SaveHistory(historyFile)
}
