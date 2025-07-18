package mocks

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// MockFileSystem 模拟文件系统
type MockFileSystem struct {
	mu    sync.RWMutex
	files map[string][]byte
	dirs  map[string]bool
}

// NewMockFileSystem 创建模拟文件系统
func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		files: make(map[string][]byte),
		dirs:  make(map[string]bool),
	}
}

// WriteFile 写入文件
func (fs *MockFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	// 确保目录存在
	dir := filepath.Dir(filename)
	fs.dirs[dir] = true

	fs.files[filename] = data
	return nil
}

// ReadFile 读取文件
func (fs *MockFileSystem) ReadFile(filename string) ([]byte, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	data, exists := fs.files[filename]
	if !exists {
		return nil, os.ErrNotExist
	}

	return data, nil
}

// Stat 获取文件信息
func (fs *MockFileSystem) Stat(filename string) (os.FileInfo, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	if _, exists := fs.files[filename]; exists {
		return &MockFileInfo{
			name: filepath.Base(filename),
			size: int64(len(fs.files[filename])),
		}, nil
	}

	if _, exists := fs.dirs[filename]; exists {
		return &MockFileInfo{
			name:  filepath.Base(filename),
			isDir: true,
		}, nil
	}

	return nil, os.ErrNotExist
}

// MkdirAll 创建目录
func (fs *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.dirs[path] = true
	return nil
}

// Remove 删除文件
func (fs *MockFileSystem) Remove(filename string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if _, exists := fs.files[filename]; exists {
		delete(fs.files, filename)
		return nil
	}

	if _, exists := fs.dirs[filename]; exists {
		delete(fs.dirs, filename)
		return nil
	}

	return os.ErrNotExist
}

// Exists 检查文件是否存在
func (fs *MockFileSystem) Exists(filename string) bool {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	_, fileExists := fs.files[filename]
	_, dirExists := fs.dirs[filename]
	return fileExists || dirExists
}

// ListFiles 列出所有文件
func (fs *MockFileSystem) ListFiles() []string {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	var files []string
	for filename := range fs.files {
		files = append(files, filename)
	}
	return files
}

// ListDirs 列出所有目录
func (fs *MockFileSystem) ListDirs() []string {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	var dirs []string
	for dirname := range fs.dirs {
		dirs = append(dirs, dirname)
	}
	return dirs
}

// Clear 清空文件系统
func (fs *MockFileSystem) Clear() {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.files = make(map[string][]byte)
	fs.dirs = make(map[string]bool)
}

// MockFileInfo 模拟文件信息
type MockFileInfo struct {
	name  string
	size  int64
	isDir bool
}

func (fi *MockFileInfo) Name() string       { return fi.name }
func (fi *MockFileInfo) Size() int64        { return fi.size }
func (fi *MockFileInfo) Mode() os.FileMode  { return 0644 }
func (fi *MockFileInfo) ModTime() time.Time { return time.Now() }
func (fi *MockFileInfo) IsDir() bool        { return fi.isDir }
func (fi *MockFileInfo) Sys() interface{}   { return nil }

// MockTemplate 模拟模板
type MockTemplate struct {
	content string
	err     error
}

// NewMockTemplate 创建模拟模板
func NewMockTemplate(content string, err error) *MockTemplate {
	return &MockTemplate{
		content: content,
		err:     err,
	}
}

// Execute 执行模板
func (t *MockTemplate) Execute(data interface{}) (string, error) {
	if t.err != nil {
		return "", t.err
	}
	return t.content, nil
}

// MockValidator 模拟验证器
type MockValidator struct {
	validateFunc func(interface{}) error
}

// NewMockValidator 创建模拟验证器
func NewMockValidator(validateFunc func(interface{}) error) *MockValidator {
	return &MockValidator{
		validateFunc: validateFunc,
	}
}

// Validate 验证
func (v *MockValidator) Validate(data interface{}) error {
	if v.validateFunc != nil {
		return v.validateFunc(data)
	}
	return nil
}

// MockHistoryManager 模拟历史记录管理器
type MockHistoryManager struct {
	histories map[string]interface{}
	err       error
}

// NewMockHistoryManager 创建模拟历史记录管理器
func NewMockHistoryManager() *MockHistoryManager {
	return &MockHistoryManager{
		histories: make(map[string]interface{}),
	}
}

// SetError 设置错误
func (h *MockHistoryManager) SetError(err error) {
	h.err = err
}

// AddHistory 添加历史记录
func (h *MockHistoryManager) AddHistory(key string, data interface{}) error {
	if h.err != nil {
		return h.err
	}
	h.histories[key] = data
	return nil
}

// GetHistory 获取历史记录
func (h *MockHistoryManager) GetHistory(key string) (interface{}, error) {
	if h.err != nil {
		return nil, h.err
	}
	data, exists := h.histories[key]
	if !exists {
		return nil, fmt.Errorf("历史记录不存在: %s", key)
	}
	return data, nil
}

// DeleteHistory 删除历史记录
func (h *MockHistoryManager) DeleteHistory(key string) error {
	if h.err != nil {
		return h.err
	}
	delete(h.histories, key)
	return nil
}

// ClearHistory 清空历史记录
func (h *MockHistoryManager) ClearHistory() error {
	if h.err != nil {
		return h.err
	}
	h.histories = make(map[string]interface{})
	return nil
}

// ListHistories 列出所有历史记录
func (h *MockHistoryManager) ListHistories() []string {
	var keys []string
	for key := range h.histories {
		keys = append(keys, key)
	}
	return keys
}
