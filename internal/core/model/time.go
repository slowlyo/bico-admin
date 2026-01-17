package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// JSONTime 自定义时间类型，用于 JSON 序列化格式化
type JSONTime time.Time

const timeFormat = "2006-01-02 15:04:05"

// MarshalJSON 实现 JSON 序列化接口
func (t JSONTime) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", time.Time(t).Format(timeFormat))
	return []byte(formatted), nil
}

// UnmarshalJSON 实现 JSON 反序列化接口
func (t *JSONTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	*t = JSONTime(now)
	return err
}

// Scan 实现 sql.Scanner 接口
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// Value 实现 driver.Valuer 接口
func (t JSONTime) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return nil, nil
	}
	return time.Time(t), nil
}

// String 重写 String 方法
func (t JSONTime) String() string {
	return time.Time(t).Format(timeFormat)
}
