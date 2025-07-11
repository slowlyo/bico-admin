package types

import (
	"time"

	"bico-admin/internal/shared/types"
)

// SystemInfoResponse 系统信息响应
type SystemInfoResponse struct {
	AppName     string    `json:"app_name"`
	Version     string    `json:"version"`
	Environment string    `json:"environment"`
	StartTime   time.Time `json:"start_time"`
	Uptime      string    `json:"uptime"`
	GoVersion   string    `json:"go_version"`
	ServerTime  time.Time `json:"server_time"`
}

// SystemStatsResponse 系统统计响应
type SystemStatsResponse struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	Goroutines  int     `json:"goroutines"`
	Connections int     `json:"connections"`
}

// ConfigListRequest 配置列表请求
type ConfigListRequest struct {
	types.BasePageQuery
	Group string `form:"group" json:"group"`
}

// ConfigCreateRequest 创建配置请求
type ConfigCreateRequest struct {
	Key         string `json:"key" binding:"required,max=100"`
	Value       string `json:"value" binding:"required"`
	Group       string `json:"group" binding:"required,max=50"`
	Description string `json:"description" binding:"max=255"`
	IsPublic    bool   `json:"is_public"`
}

// ConfigUpdateRequest 更新配置请求
type ConfigUpdateRequest struct {
	Value       string `json:"value" binding:"required"`
	Description string `json:"description" binding:"max=255"`
	IsPublic    bool   `json:"is_public"`
}

// ConfigResponse 配置响应
type ConfigResponse struct {
	ID          uint      `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Group       string    `json:"group"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// LogListRequest 日志列表请求
type LogListRequest struct {
	types.BasePageQuery
	Level     string     `form:"level" json:"level"`
	Module    string     `form:"module" json:"module"`
	StartTime *time.Time `form:"start_time" json:"start_time"`
	EndTime   *time.Time `form:"end_time" json:"end_time"`
}

// LogResponse 日志响应
type LogResponse struct {
	ID        uint      `json:"id"`
	Level     string    `json:"level"`
	Module    string    `json:"module"`
	Message   string    `json:"message"`
	Context   string    `json:"context"`
	UserID    uint      `json:"user_id"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// CacheStatsResponse 缓存统计响应
type CacheStatsResponse struct {
	RedisInfo map[string]interface{} `json:"redis_info"`
	KeyCount  int64                  `json:"key_count"`
	Memory    string                 `json:"memory"`
	Hits      int64                  `json:"hits"`
	Misses    int64                  `json:"misses"`
	HitRate   float64                `json:"hit_rate"`
}
