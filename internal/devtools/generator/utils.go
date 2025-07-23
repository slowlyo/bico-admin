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

// ToLowerCamelCase 将PascalCase转换为camelCase（首字母小写的驼峰）
func ToLowerCamelCase(str string) string {
	if str == "" {
		return ""
	}

	// 如果字符串只有一个字符，直接转小写
	if len(str) == 1 {
		return strings.ToLower(str)
	}

	// 首字母小写，其余保持不变
	return strings.ToLower(string(str[0])) + str[1:]
}

// ToPascalCase 转换为帕斯卡命名（首字母大写的驼峰）
func ToPascalCase(str string) string {
	if str == "" {
		return ""
	}

	// 如果字符串已经是有效的PascalCase，直接返回
	if isPascalCase(str) {
		return str
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

// isPascalCase 检查字符串是否已经是PascalCase格式
func isPascalCase(str string) bool {
	if str == "" {
		return false
	}

	// 首字母必须是大写
	if !unicode.IsUpper(rune(str[0])) {
		return false
	}

	// 不能包含下划线、短横线或空格
	if strings.ContainsAny(str, "_- ") {
		return false
	}

	// 检查是否只包含字母和数字
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}

	return true
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

	// 确保首字母大写（导出字段）
	if len(sanitized) > 0 {
		sanitized = strings.ToUpper(string(sanitized[0])) + sanitized[1:]
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
		"string":          "string",
		"int":             "int",
		"int32":           "int32",
		"int64":           "int64",
		"uint":            "uint",
		"uint32":          "uint32",
		"uint64":          "uint64",
		"float32":         "float32",
		"float64":         "float64",
		"bool":            "bool",
		"time":            "*time.Time",
		"date":            "*time.Time",
		"datetime":        "*time.Time",
		"timestamp":       "*time.Time",
		"text":            "string",
		"json":            "string",
		"decimal":         "float64",
		"decimal.decimal": "float64",
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

// NeedsDecimalImport 检查是否需要导入decimal包
func NeedsDecimalImport(fields []FieldDefinition) bool {
	for _, field := range fields {
		goType := GetGoType(field.Type)
		if strings.Contains(goType, "decimal.Decimal") {
			return true
		}
	}
	return false
}

// removeDuplicateStrings 去除字符串切片中的重复项
func removeDuplicateStrings(slice []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
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

// HasStatusField 检查是否有状态字段
// 支持多种状态字段名称和类型检测
func HasStatusField(fields []FieldDefinition) bool {
	for _, field := range fields {
		if IsStatusField(field) {
			return true
		}
	}
	return false
}

// IsStatusField 判断字段是否为状态字段
// 检查字段名和类型是否符合状态字段的特征
func IsStatusField(field FieldDefinition) bool {
	// 检查字段名（忽略大小写）
	fieldNameLower := strings.ToLower(field.Name)
	statusNames := []string{"status", "state", "enabled", "active"}

	isStatusName := false
	for _, name := range statusNames {
		if fieldNameLower == name {
			isStatusName = true
			break
		}
	}

	if !isStatusName {
		return false
	}

	// 检查字段类型是否为状态类型
	return IsStatusType(field.Type)
}

// IsStatusType 判断类型是否为状态类型
func IsStatusType(fieldType string) bool {
	// 获取Go类型
	goType := GetGoType(fieldType)

	// 状态字段通常是整数类型或布尔类型
	statusTypes := []string{
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"bool",
		"*int", "*int8", "*int16", "*int32", "*int64",
		"*uint", "*uint8", "*uint16", "*uint32", "*uint64",
		"*bool",
	}

	for _, statusType := range statusTypes {
		if goType == statusType {
			return true
		}
	}

	return false
}

// GetStatusFieldType 获取状态字段的类型信息
func GetStatusFieldType(field FieldDefinition) StatusFieldType {
	if !IsStatusField(field) {
		return StatusFieldTypeNone
	}

	goType := GetGoType(field.Type)

	// 判断是否为指针类型
	isPointer := strings.HasPrefix(goType, "*")
	baseType := strings.TrimPrefix(goType, "*")

	switch baseType {
	case "bool":
		if isPointer {
			return StatusFieldTypeBoolPointer
		}
		return StatusFieldTypeBool
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		if isPointer {
			return StatusFieldTypeIntPointer
		}
		return StatusFieldTypeInt
	default:
		return StatusFieldTypeNone
	}
}

// StatusFieldType 状态字段类型枚举
type StatusFieldType int

const (
	StatusFieldTypeNone        StatusFieldType = iota // 非状态字段 (0)
	StatusFieldTypeBool                               // bool类型 (1)
	StatusFieldTypeInt                                // int类型 (2)
	StatusFieldTypeBoolPointer                        // *bool类型 (3)
	StatusFieldTypeIntPointer                         // *int类型 (4)
)

// CleanComment 清理字段注释，移除括号及其内容
// 支持多种括号格式：()、（）、[]、【】等
func CleanComment(comment string) string {
	if comment == "" {
		return ""
	}

	// 定义需要移除的括号对
	bracketPairs := []struct {
		open  string
		close string
	}{
		{"(", ")"}, // 英文圆括号
		{"（", "）"}, // 中文圆括号
		{"[", "]"}, // 英文方括号
		{"【", "】"}, // 中文方括号
		{"{", "}"}, // 英文花括号
		{"｛", "｝"}, // 中文花括号
	}

	result := comment

	// 移除所有类型的括号及其内容
	for _, pair := range bracketPairs {
		result = removeBracketContent(result, pair.open, pair.close)
	}

	// 清理多余的空格和标点
	result = strings.TrimSpace(result)
	result = strings.TrimSuffix(result, "，")
	result = strings.TrimSuffix(result, ",")
	result = strings.TrimSuffix(result, "；")
	result = strings.TrimSuffix(result, ";")
	result = strings.TrimSuffix(result, "：")
	result = strings.TrimSuffix(result, ":")

	return strings.TrimSpace(result)
}

// removeBracketContent 移除指定括号及其内容，支持嵌套括号
func removeBracketContent(text, openBracket, closeBracket string) string {
	for {
		openIndex := strings.Index(text, openBracket)
		if openIndex == -1 {
			break
		}

		// 查找匹配的闭合括号，考虑嵌套情况
		closeIndex := findMatchingCloseBracket(text, openIndex, openBracket, closeBracket)
		if closeIndex == -1 {
			// 没有找到闭合括号，移除从开括号到结尾的所有内容
			text = text[:openIndex]
			break
		}

		// 移除括号及其内容
		text = text[:openIndex] + text[closeIndex+len(closeBracket):]
	}

	return text
}

// findMatchingCloseBracket 查找匹配的闭合括号，处理嵌套情况
func findMatchingCloseBracket(text string, startIndex int, openBracket, closeBracket string) int {
	openCount := 1
	searchStart := startIndex + len(openBracket)

	for i := searchStart; i < len(text); {
		if strings.HasPrefix(text[i:], openBracket) {
			openCount++
			i += len(openBracket)
		} else if strings.HasPrefix(text[i:], closeBracket) {
			openCount--
			if openCount == 0 {
				return i
			}
			i += len(closeBracket)
		} else {
			i++
		}
	}

	return -1 // 没有找到匹配的闭合括号
}

// GetDisplayComment 获取用于显示的注释（清理后的）
func GetDisplayComment(comment string) string {
	return CleanComment(comment)
}

// GetFullComment 获取完整的注释（保留原始内容）
func GetFullComment(comment string) string {
	return comment
}
