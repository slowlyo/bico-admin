package generator

import (
	"regexp"
	"strings"
	"unicode"
)

// ToSnakeCase 转换为蛇形命名
func ToSnakeCase(str string) string {
	// 处理连续的大写字母
	re1 := regexp.MustCompile("([A-Z]+)([A-Z][a-z])")
	str = re1.ReplaceAllString(str, "${1}_${2}")

	// 处理大写字母前插入下划线
	re2 := regexp.MustCompile("([a-z\\d])([A-Z])")
	str = re2.ReplaceAllString(str, "${1}_${2}")

	return strings.ToLower(str)
}

// ToCamelCase 转换为驼峰命名
func ToCamelCase(str string) string {
	if str == "" {
		return ""
	}

	// 分割字符串
	parts := strings.FieldsFunc(str, func(c rune) bool {
		return c == '_' || c == '-' || c == ' '
	})

	if len(parts) == 0 {
		return ""
	}

	// 第一个单词小写，其余单词首字母大写
	result := strings.ToLower(parts[0])
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			result += strings.ToUpper(string(parts[i][0])) + strings.ToLower(parts[i][1:])
		}
	}

	return result
}

// ToPascalCase 转换为帕斯卡命名（首字母大写的驼峰）
func ToPascalCase(str string) string {
	if str == "" {
		return ""
	}

	// 分割字符串
	parts := strings.FieldsFunc(str, func(c rune) bool {
		return c == '_' || c == '-' || c == ' '
	})

	if len(parts) == 0 {
		return ""
	}

	// 所有单词首字母大写
	var result strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(string(part[0])))
			if len(part) > 1 {
				result.WriteString(strings.ToLower(part[1:]))
			}
		}
	}

	return result.String()
}

// ToKebabCase 转换为短横线命名
func ToKebabCase(str string) string {
	if str == "" {
		return ""
	}

	// 先转换为蛇形命名，然后替换下划线为短横线
	snakeCase := ToSnakeCase(str)
	return strings.ReplaceAll(snakeCase, "_", "-")
}

// ToPlural 转换为复数形式（简单实现）
func ToPlural(str string) string {
	if str == "" {
		return ""
	}

	str = strings.ToLower(str)

	// 特殊复数形式
	irregulars := map[string]string{
		"person": "people",
		"child":  "children",
		"foot":   "feet",
		"tooth":  "teeth",
		"goose":  "geese",
		"mouse":  "mice",
		"man":    "men",
		"woman":  "women",
	}

	if plural, exists := irregulars[str]; exists {
		return plural
	}

	// 规则变化
	if strings.HasSuffix(str, "y") && len(str) > 1 {
		// 辅音字母+y结尾，变y为ies
		if !isVowel(rune(str[len(str)-2])) {
			return str[:len(str)-1] + "ies"
		}
	}

	if strings.HasSuffix(str, "s") || strings.HasSuffix(str, "sh") ||
		strings.HasSuffix(str, "ch") || strings.HasSuffix(str, "x") ||
		strings.HasSuffix(str, "z") {
		return str + "es"
	}

	if strings.HasSuffix(str, "f") {
		return str[:len(str)-1] + "ves"
	}

	if strings.HasSuffix(str, "fe") {
		return str[:len(str)-2] + "ves"
	}

	// 默认加s
	return str + "s"
}

// isVowel 判断是否为元音字母
func isVowel(r rune) bool {
	vowels := "aeiouAEIOU"
	return strings.ContainsRune(vowels, r)
}

// IsValidGoIdentifier 检查是否为有效的Go标识符
func IsValidGoIdentifier(name string) bool {
	if name == "" {
		return false
	}

	// 第一个字符必须是字母或下划线
	if !unicode.IsLetter(rune(name[0])) && name[0] != '_' {
		return false
	}

	// 其余字符必须是字母、数字或下划线
	for _, r := range name[1:] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}

	return true
}

// IsGoKeyword 检查是否为Go关键字
func IsGoKeyword(name string) bool {
	keywords := []string{
		"break", "case", "chan", "const", "continue", "default", "defer",
		"else", "fallthrough", "for", "func", "go", "goto", "if", "import",
		"interface", "map", "package", "range", "return", "select", "struct",
		"switch", "type", "var",
	}

	for _, keyword := range keywords {
		if name == keyword {
			return true
		}
	}

	return false
}

// SanitizeGoIdentifier 清理Go标识符
func SanitizeGoIdentifier(name string) string {
	if name == "" {
		return "Field"
	}

	// 移除非法字符
	var result strings.Builder
	for i, r := range name {
		if i == 0 {
			if unicode.IsLetter(r) || r == '_' {
				result.WriteRune(r)
			} else {
				result.WriteRune('F') // 默认前缀
				if unicode.IsDigit(r) {
					result.WriteRune(r)
				}
			}
		} else {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				result.WriteRune(r)
			}
		}
	}

	sanitized := result.String()
	if sanitized == "" {
		return "Field"
	}

	// 如果是关键字，添加后缀
	if IsGoKeyword(sanitized) {
		sanitized += "Field"
	}

	return sanitized
}

// GetGoType 获取Go类型
func GetGoType(fieldType string) string {
	// 如果已经是完整的Go类型（包含指针、包名等），直接返回
	if strings.Contains(fieldType, "*") || strings.Contains(fieldType, ".") {
		return fieldType
	}

	typeMap := map[string]string{
		"string":    "string",
		"int":       "int",
		"int32":     "int32",
		"int64":     "int64",
		"uint":      "uint",
		"uint32":    "uint32",
		"uint64":    "uint64",
		"float32":   "float32",
		"float64":   "float64",
		"bool":      "bool",
		"time":      "*time.Time",
		"date":      "*time.Time",
		"datetime":  "*time.Time",
		"timestamp": "*time.Time",
		"text":      "string",
		"json":      "string",
		"decimal":   "float64",
	}

	if goType, exists := typeMap[strings.ToLower(fieldType)]; exists {
		return goType
	}

	// 默认返回string
	return "string"
}

// NeedsTimeImport 检查是否需要导入time包
func NeedsTimeImport(fields []FieldDefinition) bool {
	for _, field := range fields {
		goType := GetGoType(field.Type)
		if strings.Contains(goType, "time.Time") {
			return true
		}
	}
	return false
}

// NeedsValidationImport 检查是否需要导入验证包
func NeedsValidationImport(fields []FieldDefinition) bool {
	for _, field := range fields {
		if field.Validate != "" {
			return true
		}
	}
	return false
}
