package generator

import "time"

// ComponentType 组件类型
type ComponentType string

const (
	ComponentModel        ComponentType = "model"
	ComponentRepository   ComponentType = "repository"
	ComponentService      ComponentType = "service"
	ComponentHandler      ComponentType = "handler"
	ComponentRoutes       ComponentType = "routes"
	ComponentWire         ComponentType = "wire"
	ComponentMigration    ComponentType = "migration"
	ComponentPermission   ComponentType = "permission"
	ComponentFrontendAPI  ComponentType = "frontend_api"
	ComponentFrontendPage ComponentType = "frontend_page"
	ComponentFrontendForm ComponentType = "frontend_form"
	ComponentAll          ComponentType = "all"
)

// FieldDefinition 字段定义
type FieldDefinition struct {
	Name     string `json:"name"`     // 字段名
	Type     string `json:"type"`     // 字段类型
	GormTag  string `json:"gorm_tag"` // GORM标签
	JsonTag  string `json:"json_tag"` // JSON标签
	Validate string `json:"validate"` // 验证规则
	Comment  string `json:"comment"`  // 注释
}

// GenerateRequest 生成请求
type GenerateRequest struct {
	ComponentType ComponentType     `json:"component_type"` // 组件类型
	ModelName     string            `json:"model_name"`     // 模型名称
	Fields        []FieldDefinition `json:"fields"`         // 字段定义
	TableName     string            `json:"table_name"`     // 表名（可选）
	PackagePath   string            `json:"package_path"`   // 包路径
	Options       GenerateOptions   `json:"options"`        // 生成选项
}

// GenerateOptions 生成选项
type GenerateOptions struct {
	OverwriteExisting bool `json:"overwrite_existing"` // 是否覆盖已存在的文件
	FormatCode        bool `json:"format_code"`        // 是否格式化代码
	OptimizeImports   bool `json:"optimize_imports"`   // 是否优化导入
}

// GenerateResponse 生成响应
type GenerateResponse struct {
	Success        bool     `json:"success"`
	GeneratedFiles []string `json:"generated_files"`
	Message        string   `json:"message"`
	HistoryUpdated bool     `json:"history_updated"`
	Errors         []string `json:"errors,omitempty"`
}

// TemplateData 模板数据
type TemplateData struct {
	PackageName    string            // 包名
	PackagePath    string            // 包路径
	ModelName      string            // 模型名（如User）
	ModelNameLower string            // 模型名小写（如user）
	ModelNameSnake string            // 模型名蛇形命名（如user_info）
	TableName      string            // 表名（如users）
	Fields         []FieldDefinition // 字段列表
	Imports        []string          // 导入包列表
	HasTimeField   bool              // 是否包含时间字段
	HasValidation  bool              // 是否包含验证
	Timestamp      time.Time         // 生成时间戳
}

// GenerateHistory 生成历史记录
type GenerateHistory struct {
	ModuleName     string    `json:"module_name"`
	GeneratedAt    time.Time `json:"generated_at"`
	Components     []string  `json:"components"`
	ModelName      string    `json:"model_name"`
	TableName      string    `json:"table_name"`
	PackagePath    string    `json:"package_path"`
	GeneratedBy    string    `json:"generated_by"`
	GeneratedFiles []string  `json:"generated_files"` // 存储生成的文件相对路径
}

// HistoryFile 历史文件结构
type HistoryFile struct {
	Version string            `json:"version"`
	History []GenerateHistory `json:"history"`
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error 实现error接口
func (v ValidationError) Error() string {
	return v.Field + ": " + v.Message
}

// ValidationErrors 验证错误集合
type ValidationErrors []ValidationError

// Error 实现error接口
func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	if len(v) == 1 {
		return v[0].Error()
	}
	return "多个验证错误"
}

// HasErrors 是否有错误
func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}
