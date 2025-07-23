# MCP 开发工具服务

基于 Model Context Protocol (MCP) 的开发工具服务，为 Bico Admin 项目提供便利的开发时工具。

## 功能概述

- 📖 **配置读取**: 读取并解析应用配置文件
- 🗄️ **数据库操作**: 执行SQL语句操作数据库（仅开发环境）
- � **目录树查看**: 读取项目目录结构
- 🔍 **表结构查看**: 查看数据表结构信息
- �🔧 **代码生成**: CRUD 模板生成

## 快速开始

### 第一步：构建 MCP 开发工具

在项目根目录下运行：
```bash
make build-devtools
```

构建完成后，命令会自动输出完整的 MCP 配置，直接复制使用即可。

### 第二步：配置客户端

将构建命令输出的 JSON 配置复制到你的 MCP 客户端配置文件中。

**Claude Desktop 配置文件位置**：
- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%\Claude\claude_desktop_config.json`

## 可用工具

- `read_config` - 读取应用配置文件
- `execute_sql` - 执行SQL语句（仅开发环境）
- `read_directory_tree` - 读取项目目录结构
- `inspect_table_schema` - 查看数据表结构
- `generate_code` - 代码生成工具

## 命令参考

```bash
# 构建MCP开发工具（推荐）
make build-devtools

# 查看帮助
./bin/devtools -help

# 查看版本
./bin/devtools -version

# 手动启动（调试用）
./bin/devtools -log-level debug
```

## 工具说明

### read_config
读取并解析应用配置文件，支持JSON和摘要格式输出。

### execute_sql
执行SQL语句进行数据库操作，支持查询和DML操作。

**安全限制**:
- 仅开发环境可用
- 查询结果最多1000条
- 执行超时30秒

### read_directory_tree
读取项目目录结构，返回文本格式的目录树，自动排除常见的无关目录。

### inspect_table_schema
查看数据表结构，支持查看所有表列表、指定表结构或多个表结构。

### generate_code
代码生成工具，支持生成CRUD模板代码。



## 开发指南

### 添加新工具

参考现有工具实现：
- `internal/devtools/tools/config_reader.go`
- `internal/devtools/tools/database_tool.go`

实现步骤：
1. 在 `internal/devtools/tools/` 创建工具文件
2. 实现 `ToolHandler` 接口
3. 在 `server.go` 中注册工具

## 故障排除

### 常见问题

**1. 构建失败**
**解决方案**：
- 确保在项目根目录下运行 `make build-devtools`
- 确保项目依赖已安装：`go mod tidy`

**2. MCP 客户端连接失败**
**解决方案**：
- 确保使用 `make build-devtools` 输出的完整路径配置
- 检查二进制文件是否存在：`ls -la bin/devtools`
- 手动测试：`./bin/devtools -help`

**3. 权限错误**
**解决方案**：
- 确保二进制文件有执行权限：`chmod +x bin/devtools`

## 注意事项

- 仅开发环境使用
- 使用标准输入输出(stdio)与MCP客户端通信
- 服务独立于主应用运行
- MCP客户端会自动启动和管理服务进程
- 使用 `make build-devtools` 构建后，二进制文件位于 `bin/devtools`
