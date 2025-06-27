package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// StringToUint 字符串转uint
func StringToUint(s string) (uint, error) {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}

// StringToInt 字符串转int
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// UintToString uint转字符串
func UintToString(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}

// IntToString int转字符串
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// IsEmail 验证邮箱格式
func IsEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// IsPhone 验证手机号格式（中国大陆）
func IsPhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// TrimSpaces 去除字符串首尾空格
func TrimSpaces(s string) string {
	return strings.TrimSpace(s)
}

// IsEmpty 检查字符串是否为空
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Contains 检查切片是否包含指定元素
func Contains(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

// RemoveDuplicates 去除切片中的重复元素
func RemoveDuplicates(slice interface{}) interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return slice
	}

	seen := make(map[interface{}]bool)
	result := reflect.MakeSlice(s.Type(), 0, s.Len())

	for i := 0; i < s.Len(); i++ {
		item := s.Index(i).Interface()
		if !seen[item] {
			seen[item] = true
			result = reflect.Append(result, s.Index(i))
		}
	}

	return result.Interface()
}

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// ParseTime 解析时间字符串
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", timeStr)
}

// GetCurrentTime 获取当前时间
func GetCurrentTime() time.Time {
	return time.Now()
}

// GetCurrentTimeString 获取当前时间字符串
func GetCurrentTimeString() string {
	return FormatTime(GetCurrentTime())
}

// CalculateAge 计算年龄
func CalculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	
	return age
}

// GenerateOrderNumber 生成订单号
func GenerateOrderNumber(prefix string) string {
	timestamp := time.Now().Unix()
	random, _ := GenerateRandomString(8)
	return fmt.Sprintf("%s%d%s", prefix, timestamp, random)
}
