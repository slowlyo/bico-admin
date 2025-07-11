package types

import (
	"time"

	"bico-admin/internal/shared/types"
)

// APIResponse API统一响应格式
type APIResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// HealthCheckResponse 健康检查响应
type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Services  map[string]ServiceStatus `json:"services"`
}

// ServiceStatus 服务状态
type ServiceStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// APIKeyRequest API密钥请求
type APIKeyRequest struct {
	Name        string   `json:"name" binding:"required,max=100"`
	Description string   `json:"description" binding:"max=255"`
	Permissions []string `json:"permissions" binding:"required"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

// APIKeyResponse API密钥响应
type APIKeyResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Key         string     `json:"key,omitempty"` // 只在创建时返回
	Description string     `json:"description"`
	Permissions []string   `json:"permissions"`
	Status      int        `json:"status"`
	ExpiresAt   *time.Time `json:"expires_at"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// APIKeyListRequest API密钥列表请求
type APIKeyListRequest struct {
	types.BasePageQuery
	Status *int `form:"status" json:"status"`
}

// RateLimitInfo 限流信息
type RateLimitInfo struct {
	Limit     int   `json:"limit"`
	Remaining int   `json:"remaining"`
	Reset     int64 `json:"reset"`
}

// APIUsageStats API使用统计
type APIUsageStats struct {
	TotalRequests   int64 `json:"total_requests"`
	SuccessRequests int64 `json:"success_requests"`
	ErrorRequests   int64 `json:"error_requests"`
	AvgResponseTime int64 `json:"avg_response_time"` // 毫秒
	QPS             int64 `json:"qps"`
}

// EndpointStats 端点统计
type EndpointStats struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Count       int64  `json:"count"`
	AvgTime     int64  `json:"avg_time"`
	ErrorRate   float64 `json:"error_rate"`
	LastAccess  time.Time `json:"last_access"`
}
