package utils

import (
	"encoding/json"
	"time"
)

// 常用时间格式常量
const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
)

// FormattedTime 格式化时间类型，支持自定义JSON序列化
type FormattedTime struct {
	time.Time
}

// NewFormattedTime 创建格式化时间
func NewFormattedTime(t time.Time) FormattedTime {
	return FormattedTime{Time: t}
}

// MarshalJSON 自定义JSON序列化，输出格式化的时间字符串
func (ft FormattedTime) MarshalJSON() ([]byte, error) {
	if ft.Time.IsZero() {
		return []byte("null"), nil
	}
	formatted := ft.Time.Format(DateTimeFormat)
	return json.Marshal(formatted)
}

// UnmarshalJSON 自定义JSON反序列化
func (ft *FormattedTime) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}
	
	if timeStr == "" {
		ft.Time = time.Time{}
		return nil
	}
	
	// 尝试多种时间格式
	formats := []string{
		DateTimeFormat,
		time.RFC3339,
		time.RFC3339Nano,
		DateFormat,
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			ft.Time = t
			return nil
		}
	}
	
	return &time.ParseError{
		Layout: DateTimeFormat,
		Value:  timeStr,
	}
}

// FormatTime 格式化时间为字符串
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DateTimeFormat)
}

// FormatTimePtr 格式化时间指针为字符串
func FormatTimePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(DateTimeFormat)
}

// FormatDate 格式化日期为字符串
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DateFormat)
}

// FormatDatePtr 格式化日期指针为字符串
func FormatDatePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(DateFormat)
}

// TimeResponse 时间响应结构，包含原始时间和格式化字符串
type TimeResponse struct {
	Raw       time.Time `json:"raw"`       // 原始时间（ISO格式）
	Formatted string    `json:"formatted"` // 格式化时间字符串
}

// NewTimeResponse 创建时间响应
func NewTimeResponse(t time.Time) TimeResponse {
	return TimeResponse{
		Raw:       t,
		Formatted: FormatTime(t),
	}
}

// NewTimeResponsePtr 创建时间指针响应
func NewTimeResponsePtr(t *time.Time) *TimeResponse {
	if t == nil || t.IsZero() {
		return nil
	}
	resp := NewTimeResponse(*t)
	return &resp
}
