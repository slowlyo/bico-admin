package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"bico-admin/internal/devtools/types"
	"bico-admin/pkg/config"
	"bico-admin/pkg/database"
)

// TableSchemaTool 数据表结构查看工具
type TableSchemaTool struct {
	db *gorm.DB
}

// NewTableSchemaTool 创建数据表结构查看工具
func NewTableSchemaTool() *TableSchemaTool {
	return &TableSchemaTool{}
}

// TableInfo 表信息结构
type TableInfo struct {
	TableName    string       `json:"table_name"`
	TableComment string       `json:"table_comment,omitempty"`
	Columns      []ColumnInfo `json:"columns"`
	Indexes      []IndexInfo  `json:"indexes,omitempty"`
	RowCount     int64        `json:"row_count,omitempty"`
}

// ColumnInfo 列信息结构
type ColumnInfo struct {
	ColumnName    string `json:"column_name"`
	DataType      string `json:"data_type"`
	IsNullable    string `json:"is_nullable"`
	ColumnDefault string `json:"column_default,omitempty"`
	ColumnComment string `json:"column_comment,omitempty"`
	Extra         string `json:"extra,omitempty"`
	ColumnKey     string `json:"column_key,omitempty"`
}

// IndexInfo 索引信息结构
type IndexInfo struct {
	IndexName    string `json:"index_name"`
	ColumnName   string `json:"column_name"`
	NonUnique    int    `json:"non_unique"`
	IndexType    string `json:"index_type,omitempty"`
	IndexComment string `json:"index_comment,omitempty"`
}

// GetTool 获取MCP工具定义
func (t *TableSchemaTool) GetTool() mcp.Tool {
	return mcp.NewTool("inspect_table_schema",
		mcp.WithDescription("查看数据表结构，支持查看所有表列表、指定表结构或多个表结构"),
		mcp.WithString("action",
			mcp.Required(),
			mcp.Description("操作类型"),
			mcp.Enum("list_tables", "describe_table", "describe_tables"),
		),
		mcp.WithString("table_names",
			mcp.Description("表名，多个表名用逗号分隔（action为describe_table或describe_tables时必需）"),
		),
		mcp.WithBoolean("include_indexes",
			mcp.Description("是否包含索引信息（默认true）"),
		),
		mcp.WithBoolean("include_row_count",
			mcp.Description("是否包含行数统计（默认false，大表可能较慢）"),
		),
	)
}

// Handle 处理工具调用
func (t *TableSchemaTool) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 解析参数
	action := strings.TrimSpace(request.GetString("action", ""))
	tableNames := strings.TrimSpace(request.GetString("table_names", ""))
	includeIndexes := request.GetBool("include_indexes", true)
	includeRowCount := request.GetBool("include_row_count", false)

	// 参数验证
	if action == "" {
		return mcp.NewToolResultError("action参数不能为空"), nil
	}

	if (action == "describe_table" || action == "describe_tables") && tableNames == "" {
		return mcp.NewToolResultError("describe_table和describe_tables操作需要指定table_names参数"), nil
	}

	// 初始化数据库连接
	if err := t.initDatabase(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("数据库连接失败: %v", err)), nil
	}

	// 检查环境限制
	if err := t.checkEnvironment(); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// 执行相应操作
	switch action {
	case "list_tables":
		return t.listTables(ctx)
	case "describe_table":
		tables := []string{tableNames}
		return t.describeTables(ctx, tables, includeIndexes, includeRowCount)
	case "describe_tables":
		tables := strings.Split(tableNames, ",")
		for i := range tables {
			tables[i] = strings.TrimSpace(tables[i])
		}
		return t.describeTables(ctx, tables, includeIndexes, includeRowCount)
	default:
		return mcp.NewToolResultError(fmt.Sprintf("不支持的操作类型: %s", action)), nil
	}
}

// initDatabase 初始化数据库连接
func (t *TableSchemaTool) initDatabase() error {
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
func (t *TableSchemaTool) checkEnvironment() error {
	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("无法获取环境配置: %w", err)
	}

	// 仅在开发环境下允许使用
	if cfg.App.Environment != "development" && cfg.App.Environment != "dev" {
		return fmt.Errorf("表结构查看工具仅在开发环境下可用，当前环境: %s", cfg.App.Environment)
	}

	return nil
}

// listTables 列出所有表
func (t *TableSchemaTool) listTables(ctx context.Context) (*mcp.CallToolResult, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 获取数据库驱动类型
	cfg, _ := config.Load("")
	var sql string

	switch cfg.Database.Driver {
	case "mysql":
		sql = "SELECT table_name, table_comment FROM information_schema.tables WHERE table_schema = DATABASE() ORDER BY table_name"
	case "postgres":
		sql = `SELECT table_name, '' as table_comment 
			   FROM information_schema.tables 
			   WHERE table_schema = 'public' AND table_type = 'BASE TABLE' 
			   ORDER BY table_name`
	case "sqlite":
		sql = "SELECT name as table_name, '' as table_comment FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name"
	default:
		return mcp.NewToolResultError(fmt.Sprintf("不支持的数据库驱动: %s", cfg.Database.Driver)), nil
	}

	rows, err := t.db.WithContext(queryCtx).Raw(sql).Rows()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("查询表列表失败: %v", err)), nil
	}
	defer rows.Close()

	var tables []map[string]string
	for rows.Next() {
		var tableName, tableComment string
		if err := rows.Scan(&tableName, &tableComment); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("扫描表信息失败: %v", err)), nil
		}
		tables = append(tables, map[string]string{
			"table_name":    tableName,
			"table_comment": tableComment,
		})
	}

	if err := rows.Err(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("行扫描错误: %v", err)), nil
	}

	// 构建响应
	response := types.ToolResponse{
		Success: true,
		Data: map[string]any{
			"action":      "list_tables",
			"tables":      tables,
			"table_count": len(tables),
		},
		Message: fmt.Sprintf("成功获取 %d 个表的信息", len(tables)),
	}

	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("序列化响应失败: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}

// describeTables 描述表结构
func (t *TableSchemaTool) describeTables(ctx context.Context, tableNames []string, includeIndexes, includeRowCount bool) (*mcp.CallToolResult, error) {
	var tableInfos []TableInfo

	for _, tableName := range tableNames {
		if tableName == "" {
			continue
		}

		tableInfo, err := t.getTableInfo(ctx, tableName, includeIndexes, includeRowCount)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("获取表 %s 信息失败: %v", tableName, err)), nil
		}
		tableInfos = append(tableInfos, *tableInfo)
	}

	// 构建响应
	action := "describe_table"
	if len(tableNames) > 1 {
		action = "describe_tables"
	}

	response := types.ToolResponse{
		Success: true,
		Data: map[string]any{
			"action":      action,
			"tables":      tableInfos,
			"table_count": len(tableInfos),
		},
		Message: fmt.Sprintf("成功获取 %d 个表的结构信息", len(tableInfos)),
	}

	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("序列化响应失败: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}

// getTableInfo 获取单个表的详细信息
func (t *TableSchemaTool) getTableInfo(ctx context.Context, tableName string, includeIndexes, includeRowCount bool) (*TableInfo, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cfg, _ := config.Load("")
	tableInfo := &TableInfo{
		TableName: tableName,
	}

	// 获取列信息
	columns, err := t.getTableColumns(queryCtx, tableName, cfg.Database.Driver)
	if err != nil {
		return nil, fmt.Errorf("获取列信息失败: %w", err)
	}
	tableInfo.Columns = columns

	// 获取表注释（MySQL支持）
	if cfg.Database.Driver == "mysql" {
		comment, err := t.getTableComment(queryCtx, tableName)
		if err == nil {
			tableInfo.TableComment = comment
		}
	}

	// 获取索引信息
	if includeIndexes {
		indexes, err := t.getTableIndexes(queryCtx, tableName, cfg.Database.Driver)
		if err == nil {
			tableInfo.Indexes = indexes
		}
	}

	// 获取行数统计
	if includeRowCount {
		count, err := t.getTableRowCount(queryCtx, tableName)
		if err == nil {
			tableInfo.RowCount = count
		}
	}

	return tableInfo, nil
}

// getTableColumns 获取表的列信息
func (t *TableSchemaTool) getTableColumns(ctx context.Context, tableName, driver string) ([]ColumnInfo, error) {
	var sql string

	switch driver {
	case "mysql":
		sql = `SELECT
			column_name,
			data_type,
			is_nullable,
			COALESCE(column_default, '') as column_default,
			COALESCE(column_comment, '') as column_comment,
			COALESCE(extra, '') as extra,
			COALESCE(column_key, '') as column_key
		FROM information_schema.columns
		WHERE table_schema = DATABASE() AND table_name = ?
		ORDER BY ordinal_position`

	case "postgres":
		sql = `SELECT
			column_name,
			data_type,
			is_nullable,
			COALESCE(column_default, '') as column_default,
			'' as column_comment,
			'' as extra,
			'' as column_key
		FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = $1
		ORDER BY ordinal_position`

	case "sqlite":
		// SQLite使用PRAGMA table_info
		sql = fmt.Sprintf("PRAGMA table_info(%s)", tableName)

	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", driver)
	}

	var columns []ColumnInfo

	if driver == "sqlite" {
		// SQLite特殊处理
		rows, err := t.db.WithContext(ctx).Raw(sql).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var cid int
			var name, dataType string
			var defaultValue *string // 使用指针类型处理NULL值
			var notNull, pk int

			if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk); err != nil {
				return nil, err
			}

			isNullable := "YES"
			if notNull == 1 {
				isNullable = "NO"
			}

			columnKey := ""
			if pk == 1 {
				columnKey = "PRI"
			}

			// 处理默认值
			defaultVal := ""
			if defaultValue != nil {
				defaultVal = *defaultValue
			}

			columns = append(columns, ColumnInfo{
				ColumnName:    name,
				DataType:      dataType,
				IsNullable:    isNullable,
				ColumnDefault: defaultVal,
				ColumnKey:     columnKey,
			})
		}
	} else {
		// MySQL和PostgreSQL
		rows, err := t.db.WithContext(ctx).Raw(sql, tableName).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var col ColumnInfo
			if err := rows.Scan(
				&col.ColumnName,
				&col.DataType,
				&col.IsNullable,
				&col.ColumnDefault,
				&col.ColumnComment,
				&col.Extra,
				&col.ColumnKey,
			); err != nil {
				return nil, err
			}
			columns = append(columns, col)
		}
	}

	return columns, nil
}

// getTableComment 获取表注释（MySQL）
func (t *TableSchemaTool) getTableComment(ctx context.Context, tableName string) (string, error) {
	sql := `SELECT table_comment
			FROM information_schema.tables
			WHERE table_schema = DATABASE() AND table_name = ?`

	var comment string
	err := t.db.WithContext(ctx).Raw(sql, tableName).Scan(&comment).Error
	return comment, err
}

// getTableIndexes 获取表索引信息
func (t *TableSchemaTool) getTableIndexes(ctx context.Context, tableName, driver string) ([]IndexInfo, error) {
	var sql string

	switch driver {
	case "mysql":
		sql = `SELECT
			index_name,
			column_name,
			non_unique,
			index_type,
			COALESCE(index_comment, '') as index_comment
		FROM information_schema.statistics
		WHERE table_schema = DATABASE() AND table_name = ?
		ORDER BY index_name, seq_in_index`

	case "postgres":
		sql = `SELECT
			i.relname as index_name,
			a.attname as column_name,
			CASE WHEN ix.indisunique THEN 0 ELSE 1 END as non_unique,
			am.amname as index_type,
			'' as index_comment
		FROM pg_class t
		JOIN pg_index ix ON t.oid = ix.indrelid
		JOIN pg_class i ON i.oid = ix.indexrelid
		JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = ANY(ix.indkey)
		JOIN pg_am am ON i.relam = am.oid
		WHERE t.relname = $1 AND t.relkind = 'r'
		ORDER BY i.relname, a.attnum`

	case "sqlite":
		// SQLite索引查询比较复杂，这里简化处理
		sql = fmt.Sprintf("PRAGMA index_list(%s)", tableName)

	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", driver)
	}

	var indexes []IndexInfo

	if driver == "sqlite" {
		// SQLite特殊处理
		rows, err := t.db.WithContext(ctx).Raw(sql).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var seq int
			var name, origin string
			var unique int
			var partial string

			if err := rows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
				return nil, err
			}

			// 获取索引的列信息
			indexInfoSQL := fmt.Sprintf("PRAGMA index_info(%s)", name)
			infoRows, err := t.db.WithContext(ctx).Raw(indexInfoSQL).Rows()
			if err != nil {
				continue
			}

			for infoRows.Next() {
				var seqno, cid int
				var columnName string
				if err := infoRows.Scan(&seqno, &cid, &columnName); err != nil {
					continue
				}

				indexes = append(indexes, IndexInfo{
					IndexName:  name,
					ColumnName: columnName,
					NonUnique:  1 - unique, // SQLite的unique字段含义相反
					IndexType:  "BTREE",
				})
			}
			infoRows.Close()
		}
	} else {
		// MySQL和PostgreSQL
		rows, err := t.db.WithContext(ctx).Raw(sql, tableName).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var idx IndexInfo
			if err := rows.Scan(
				&idx.IndexName,
				&idx.ColumnName,
				&idx.NonUnique,
				&idx.IndexType,
				&idx.IndexComment,
			); err != nil {
				return nil, err
			}
			indexes = append(indexes, idx)
		}
	}

	// 按索引名排序
	sort.Slice(indexes, func(i, j int) bool {
		if indexes[i].IndexName == indexes[j].IndexName {
			return indexes[i].ColumnName < indexes[j].ColumnName
		}
		return indexes[i].IndexName < indexes[j].IndexName
	})

	return indexes, nil
}

// getTableRowCount 获取表行数
func (t *TableSchemaTool) getTableRowCount(ctx context.Context, tableName string) (int64, error) {
	sql := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)

	var count int64
	err := t.db.WithContext(ctx).Raw(sql).Scan(&count).Error
	return count, err
}
