# MCP 开发工具服务

基于 Model Context Protocol (MCP) 的开发工具服务，为 Bico Admin 项目提供便利的开发时工具。

## 功能概述

- 📖 **配置读取**: 读取并解析应用配置文件
- 🗄️ **数据库操作**: 执行SQL语句操作数据库（仅开发环境）
- 🔧 **代码生成**: CRUD 模板生成（规划中）

## 快速开始

### 启动服务

```bash
# 启动 MCP 服务
make devtools
```

服务将在 `http://localhost:18901/mcp` 启动。

### 配置客户端

在 MCP 客户端配置文件中添加：

```json
{
  "mcpServers": {
    "bico-admin-devtools": {
      "url": "http://localhost:18901/mcp"
    }
  }
}
```

## 可用工具

- `read_config` - 读取应用配置文件
- `execute_sql` - 执行SQL语句（仅开发环境）

## 命令参考

```bash
# 启动服务
make devtools

# 查看帮助
make devtools-help

# 查看版本
make devtools-version

# 手动启动（调试用）
go run ./cmd/devtools -host 0.0.0.0 -port 18901 -log-level debug
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

## 配置示例

**Claude Desktop** (`~/Library/Application Support/Claude/claude_desktop_config.json`):
```json
{
  "mcpServers": {
    "bico-admin-devtools": {
      "url": "http://localhost:18901/mcp"
    }
  }
}
```

**使用前需要先启动服务**:
```bash
cd /path/to/bico-admin
make devtools
```

## 开发指南

### 添加新工具

参考现有工具实现：
- `internal/devtools/tools/config_reader.go`
- `internal/devtools/tools/database_tool.go`

实现步骤：
1. 在 `internal/devtools/tools/` 创建工具文件
2. 实现 `ToolHandler` 接口
3. 在 `server.go` 中注册工具

## 注意事项

- 仅开发环境使用
- 默认端口 18901
- 服务独立于主应用运行
