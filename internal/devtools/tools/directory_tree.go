package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"bico-admin/internal/devtools/types"
)

// DirectoryTreeTool 目录树工具
type DirectoryTreeTool struct {
	// 默认排除的目录
	defaultExcludeDirs []string
}

// NewDirectoryTreeTool 创建目录树工具
func NewDirectoryTreeTool() *DirectoryTreeTool {
	return &DirectoryTreeTool{
		defaultExcludeDirs: []string{
			"vendor",
			"node_modules",
			".git",
			".idea",
			".vscode",
			"dist",
			"build",
			"target",
			"bin",
			"obj",
			"tmp",
			"temp",
			".DS_Store",
			"Thumbs.db",
			"*.log",
			"*.tmp",
		},
	}
}

// GetTool 获取MCP工具定义
func (t *DirectoryTreeTool) GetTool() mcp.Tool {
	return mcp.NewTool("read_directory_tree",
		mcp.WithDescription("读取项目目录结构，返回文本格式的目录树，自动排除常见的无关目录"),
		mcp.WithString("root_path",
			mcp.Description("根目录路径，为空时使用当前工作目录"),
		),
		mcp.WithNumber("max_depth",
			mcp.Description("最大遍历深度，默认为10"),
		),
		mcp.WithString("exclude_dirs",
			mcp.Description("额外排除的目录列表，用逗号分隔"),
		),
		mcp.WithBoolean("show_files",
			mcp.Description("是否显示文件，默认为true"),
		),
		mcp.WithBoolean("show_hidden",
			mcp.Description("是否显示隐藏文件和目录，默认为false"),
		),
		mcp.WithString("format",
			mcp.Description("输出格式"),
			mcp.Enum("tree", "list", "json"),
		),
	)
}

// Handle 处理工具调用
func (t *DirectoryTreeTool) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 解析参数
	rootPath := request.GetString("root_path", ".")

	// 解析最大深度
	maxDepth := 10
	if maxDepthStr := request.GetString("max_depth", ""); maxDepthStr != "" {
		if parsed, err := strconv.Atoi(maxDepthStr); err == nil && parsed > 0 {
			maxDepth = parsed
		}
	}

	// 解析排除目录列表
	var excludeDirs []string
	if excludeDirsStr := request.GetString("exclude_dirs", ""); excludeDirsStr != "" {
		excludeDirs = strings.Split(excludeDirsStr, ",")
		// 去除空白字符
		for i, dir := range excludeDirs {
			excludeDirs[i] = strings.TrimSpace(dir)
		}
	}

	showFiles := request.GetBool("show_files", true)
	showHidden := request.GetBool("show_hidden", false)
	format := request.GetString("format", "tree")

	// 获取绝对路径
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("获取绝对路径失败: %v", err)), nil
	}

	// 检查路径是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return mcp.NewToolResultError(fmt.Sprintf("路径不存在: %s", absPath)), nil
	}

	// 合并排除目录
	allExcludeDirs := append(t.defaultExcludeDirs, excludeDirs...)

	// 构建目录树
	tree, err := t.buildDirectoryTree(absPath, maxDepth, allExcludeDirs, showFiles, showHidden)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("构建目录树失败: %v", err)), nil
	}

	// 根据格式返回结果
	switch format {
	case "tree":
		return t.handleTreeFormat(tree, absPath)
	case "list":
		return t.handleListFormat(tree, absPath)
	case "json":
		return t.handleJSONFormat(tree, absPath)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("不支持的输出格式: %s", format)), nil
	}
}

// DirectoryNode 目录节点
type DirectoryNode struct {
	Name     string           `json:"name"`
	Path     string           `json:"path"`
	IsDir    bool             `json:"is_dir"`
	Size     int64            `json:"size,omitempty"`
	Children []*DirectoryNode `json:"children,omitempty"`
}

// buildDirectoryTree 构建目录树
func (t *DirectoryTreeTool) buildDirectoryTree(rootPath string, maxDepth int, excludeDirs []string, showFiles, showHidden bool) (*DirectoryNode, error) {
	return t.buildNode(rootPath, 0, maxDepth, excludeDirs, showFiles, showHidden)
}

// buildNode 构建单个节点
func (t *DirectoryTreeTool) buildNode(path string, currentDepth, maxDepth int, excludeDirs []string, showFiles, showHidden bool) (*DirectoryNode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	node := &DirectoryNode{
		Name:  filepath.Base(path),
		Path:  path,
		IsDir: info.IsDir(),
	}

	if !info.IsDir() {
		node.Size = info.Size()
		return node, nil
	}

	// 如果达到最大深度，不再递归
	if currentDepth >= maxDepth {
		return node, nil
	}

	// 读取目录内容
	entries, err := os.ReadDir(path)
	if err != nil {
		return node, nil // 忽略无法读取的目录
	}

	// 排序
	sort.Slice(entries, func(i, j int) bool {
		// 目录优先，然后按名称排序
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		name := entry.Name()

		// 跳过隐藏文件（如果不显示）
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}

		// 跳过排除的目录
		if entry.IsDir() && t.shouldExclude(name, excludeDirs) {
			continue
		}

		// 跳过文件（如果不显示）
		if !entry.IsDir() && !showFiles {
			continue
		}

		childPath := filepath.Join(path, name)
		childNode, err := t.buildNode(childPath, currentDepth+1, maxDepth, excludeDirs, showFiles, showHidden)
		if err != nil {
			continue // 忽略错误的子节点
		}

		node.Children = append(node.Children, childNode)
	}

	return node, nil
}

// shouldExclude 检查是否应该排除目录
func (t *DirectoryTreeTool) shouldExclude(name string, excludeDirs []string) bool {
	for _, exclude := range excludeDirs {
		if matched, _ := filepath.Match(exclude, name); matched {
			return true
		}
		if exclude == name {
			return true
		}
	}
	return false
}

// handleTreeFormat 处理树形格式输出
func (t *DirectoryTreeTool) handleTreeFormat(tree *DirectoryNode, rootPath string) (*mcp.CallToolResult, error) {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("目录树: %s\n", rootPath))
	builder.WriteString("├── " + tree.Name + "\n")

	t.renderTreeNode(tree, "", true, &builder)

	response := types.ToolResponse{
		Success: true,
		Data:    builder.String(),
		Message: "目录树生成成功",
	}

	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("序列化响应失败: %v", err)), nil
	}
	return mcp.NewToolResultText(string(responseJSON)), nil
}

// renderTreeNode 渲染树形节点
func (t *DirectoryTreeTool) renderTreeNode(node *DirectoryNode, prefix string, isLast bool, builder *strings.Builder) {
	if len(node.Children) == 0 {
		return
	}

	for i, child := range node.Children {
		isChildLast := i == len(node.Children)-1

		// 构建当前行的前缀
		var currentPrefix string
		if isChildLast {
			currentPrefix = prefix + "└── "
		} else {
			currentPrefix = prefix + "├── "
		}

		// 写入当前节点
		if child.IsDir {
			builder.WriteString(currentPrefix + child.Name + "/\n")
		} else {
			sizeStr := ""
			if child.Size > 0 {
				sizeStr = fmt.Sprintf(" (%s)", formatFileSize(child.Size))
			}
			builder.WriteString(currentPrefix + child.Name + sizeStr + "\n")
		}

		// 递归处理子节点
		var nextPrefix string
		if isChildLast {
			nextPrefix = prefix + "    "
		} else {
			nextPrefix = prefix + "│   "
		}

		t.renderTreeNode(child, nextPrefix, isChildLast, builder)
	}
}

// handleListFormat 处理列表格式输出
func (t *DirectoryTreeTool) handleListFormat(tree *DirectoryNode, rootPath string) (*mcp.CallToolResult, error) {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("目录列表: %s\n\n", rootPath))

	t.renderListNode(tree, 0, &builder)

	response := types.ToolResponse{
		Success: true,
		Data:    builder.String(),
		Message: "目录列表生成成功",
	}

	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("序列化响应失败: %v", err)), nil
	}
	return mcp.NewToolResultText(string(responseJSON)), nil
}

// renderListNode 渲染列表节点
func (t *DirectoryTreeTool) renderListNode(node *DirectoryNode, depth int, builder *strings.Builder) {
	indent := strings.Repeat("  ", depth)

	if node.IsDir {
		builder.WriteString(fmt.Sprintf("%s[DIR]  %s/\n", indent, node.Name))
	} else {
		sizeStr := ""
		if node.Size > 0 {
			sizeStr = fmt.Sprintf(" (%s)", formatFileSize(node.Size))
		}
		builder.WriteString(fmt.Sprintf("%s[FILE] %s%s\n", indent, node.Name, sizeStr))
	}

	for _, child := range node.Children {
		t.renderListNode(child, depth+1, builder)
	}
}

// handleJSONFormat 处理JSON格式输出
func (t *DirectoryTreeTool) handleJSONFormat(tree *DirectoryNode, rootPath string) (*mcp.CallToolResult, error) {
	response := types.ToolResponse{
		Success: true,
		Data: map[string]interface{}{
			"root_path": rootPath,
			"tree":      tree,
		},
		Message: "目录树JSON生成成功",
	}

	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("序列化响应失败: %v", err)), nil
	}
	return mcp.NewToolResultText(string(responseJSON)), nil
}

// formatFileSize 格式化文件大小
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
