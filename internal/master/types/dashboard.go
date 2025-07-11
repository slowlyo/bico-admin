package types

import "time"

// DashboardStatsResponse 主控端仪表板统计响应
type DashboardStatsResponse struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	TotalSessions int64 `json:"total_sessions"`
	OnlineUsers   int64 `json:"online_users"`
	SystemLoad    SystemLoadStats `json:"system_load"`
	RecentActivity []ActivityItem `json:"recent_activity"`
}

// SystemLoadStats 系统负载统计
type SystemLoadStats struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	NetworkIO   NetworkIOStats `json:"network_io"`
}

// NetworkIOStats 网络IO统计
type NetworkIOStats struct {
	BytesIn  int64 `json:"bytes_in"`
	BytesOut int64 `json:"bytes_out"`
}

// ActivityItem 活动项
type ActivityItem struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

// MonitoringDataResponse 监控数据响应
type MonitoringDataResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Metrics   MetricsData `json:"metrics"`
}

// MetricsData 指标数据
type MetricsData struct {
	CPU     float64 `json:"cpu"`
	Memory  float64 `json:"memory"`
	Disk    float64 `json:"disk"`
	Network NetworkMetrics `json:"network"`
	Database DatabaseMetrics `json:"database"`
	Cache   CacheMetrics `json:"cache"`
}

// NetworkMetrics 网络指标
type NetworkMetrics struct {
	Connections int64 `json:"connections"`
	Throughput  int64 `json:"throughput"`
}

// DatabaseMetrics 数据库指标
type DatabaseMetrics struct {
	Connections int64 `json:"connections"`
	QPS         int64 `json:"qps"`
	SlowQueries int64 `json:"slow_queries"`
}

// CacheMetrics 缓存指标
type CacheMetrics struct {
	HitRate     float64 `json:"hit_rate"`
	Memory      int64   `json:"memory"`
	Connections int64   `json:"connections"`
}
