package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"bico-admin/internal/devtools/types"
	"bico-admin/pkg/config"
	"bico-admin/pkg/database"
)

// DatabaseTool 数据库操作工具
type DatabaseTool struct {
	db *gorm.DB
}

// NewDatabaseTool 创建数据库操作工具
func NewDatabaseTool() *DatabaseTool {
	return &DatabaseTool{}
}

// GetTool 获取MCP工具定义
func (t *DatabaseTool) GetTool() mcp.Tool {
	return mcp.NewTool("execute_sql",
		mcp.WithDescription("执行SQL语句，支持查询和DML操作（仅开发环境可用）"),
		mcp.WithString("sql",
			mcp.Required(),
			mcp.Description("要执行的SQL语句"),
		),
		mcp.WithString("operation_type",
			mcp.Description("操作类型"),
			mcp.Enum("query", "exec"),
		),
		mcp.WithNumber("limit",
			mcp.Description("查询结果限制条数（仅对query操作有效）"),
		),
	)
}

// Handle 处理工具调用
func (t *DatabaseTool) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 解析参数
	sql := strings.TrimSpace(request.GetString("sql", ""))
	operationType := request.GetString("operation_type", "query")

	// 获取limit参数，默认100
	limit := 100
	if limitStr := request.GetString("limit", ""); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// 参数验证
	if sql == "" {
		return mcp.NewToolResultError("SQL语句不能为空"), nil
	}

	if limit <= 0 || limit > 1000 {
		limit = 100 // 默认限制100条，最大1000条
	}

	// 初始化数据库连接
	if err := t.initDatabase(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("数据库连接失败: %v", err)), nil
	}

	// 检查环境限制
	if err := t.checkEnvironment(); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// 执行SQL
	switch operationType {
	case "query":
		return t.executeQuery(ctx, sql, limit)
	case "exec":
		return t.executeExec(ctx, sql)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("不支持的操作类型: %s", operationType)), nil
	}
}

// initDatabase 初始化数据库连接
func (t *DatabaseTool) initDatabase() error {
	if t.db != nil {
		return nil // 已经初始化
	}

	// 加载配置
	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 创建数据库连接
	logLevel := logger.Silent
	if cfg.App.Debug {
		logLevel = logger.Info
	}

	switch cfg.Database.Driver {
	case "mysql":
		dbConfig := database.MySQLConfig{
			Host:            cfg.Database.Host,
			Port:            cfg.Database.Port,
			Username:        cfg.Database.Username,
			Password:        cfg.Database.Password,
			Database:        cfg.Database.Database,
			Charset:         cfg.Database.Charset,
			MaxIdleConns:    cfg.Database.MaxIdleConns,
			MaxOpenConns:    cfg.Database.MaxOpenConns,
			ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
			LogLevel:        logLevel,
		}
		t.db, err = database.NewMySQL(dbConfig)

	case "postgres":
		dbConfig := database.PostgresConfig{
			Host:            cfg.Database.Host,
			Port:            cfg.Database.Port,
			Username:        cfg.Database.Username,
			Password:        cfg.Database.Password,
			Database:        cfg.Database.Database,
			MaxIdleConns:    cfg.Database.MaxIdleConns,
			MaxOpenConns:    cfg.Database.MaxOpenConns,
			ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
			LogLevel:        logLevel,
		}
		t.db, err = database.NewPostgres(dbConfig)

	case "sqlite":
		dbConfig := database.SQLiteConfig{
			Database:        cfg.Database.Database,
			MaxIdleConns:    cfg.Database.MaxIdleConns,
			MaxOpenConns:    cfg.Database.MaxOpenConns,
			ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
			LogLevel:        logLevel,
		}
		t.db, err = database.NewSQLite(dbConfig)

	default:
		return fmt.Errorf("不支持的数据库驱动: %s", cfg.Database.Driver)
	}

	return err
}

// checkEnvironment 检查环境限制
func (t *DatabaseTool) checkEnvironment() error {
	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("无法获取环境配置: %w", err)
	}

	// 仅在开发环境下允许使用
	if cfg.App.Environment != "development" && cfg.App.Environment != "dev" {
		return fmt.Errorf("数据库操作工具仅在开发环境下可用，当前环境: %s", cfg.App.Environment)
	}

	return nil
}

// executeQuery 执行查询操作
func (t *DatabaseTool) executeQuery(ctx context.Context, sql string, limit int) (*mcp.CallToolResult, error) {
	// 设置查询超时
	queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 执行查询
	rows, err := t.db.WithContext(queryCtx).Raw(sql).Limit(limit).Rows()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("查询执行失败: %v", err)), nil
	}
	defer rows.Close()

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("获取列信息失败: %v", err)), nil
	}

	// 读取数据
	var results []map[string]interface{}
	for rows.Next() {
		// 创建扫描目标
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// 扫描行数据
		if err := rows.Scan(valuePtrs...); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("扫描行数据失败: %v", err)), nil
		}

		// 构建结果行
		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if val != nil {
				// 处理字节数组
				if b, ok := val.([]byte); ok {
					row[col] = string(b)
				} else {
					row[col] = val
				}
			} else {
				row[col] = nil
			}
		}
		results = append(results, row)
	}

	// 检查行扫描错误
	if err := rows.Err(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("行扫描错误: %v", err)), nil
	}

	// 构建响应
	response := types.ToolResponse{
		Success: true,
		Data: map[string]interface{}{
			"sql":       sql,
			"operation": "query",
			"columns":   columns,
			"rows":      results,
			"row_count": len(results),
			"limited":   len(results) == limit,
			"limit":     limit,
		},
		Message: fmt.Sprintf("查询执行成功，返回 %d 行数据", len(results)),
	}

	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("序列化响应失败: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}

// executeExec 执行DML操作
func (t *DatabaseTool) executeExec(ctx context.Context, sql string) (*mcp.CallToolResult, error) {
	// 设置执行超时
	execCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 执行SQL
	result := t.db.WithContext(execCtx).Exec(sql)
	if result.Error != nil {
		return mcp.NewToolResultError(fmt.Sprintf("SQL执行失败: %v", result.Error)), nil
	}

	// 构建响应
	response := types.ToolResponse{
		Success: true,
		Data: map[string]interface{}{
			"sql":           sql,
			"operation":     "exec",
			"rows_affected": result.RowsAffected,
		},
		Message: fmt.Sprintf("SQL执行成功，影响 %d 行", result.RowsAffected),
	}

	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("序列化响应失败: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}
