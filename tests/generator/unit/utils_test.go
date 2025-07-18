package unit

import (
	"testing"

	"bico-admin/internal/devtools/generator"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"User", "user"},
		{"UserProfile", "user_profile"},
		{"XMLHttpRequest", "xml_http_request"},
		{"HTTPSConnection", "https_connection"},
		{"APIKey", "api_key"},
		{"ID", "id"},
		{"UserID", "user_id"},
		{"HTMLParser", "html_parser"},
		{"JSONData", "json_data"},
		{"URLPath", "url_path"},
		{"", ""},
		{"A", "a"},
		{"AB", "ab"},
		{"ABC", "abc"},
		{"AbC", "ab_c"},
		{"AbCd", "ab_cd"},
		{"AbCdE", "ab_cd_e"},
		{"ProductCategory", "product_category"},
		{"OrderItem", "order_item"},
		{"PaymentMethod", "payment_method"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := generator.ToSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToSnakeCase(%q) = %q, 期望 %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user", "user"},
		{"user_profile", "userProfile"},
		{"user-profile", "userProfile"},
		{"user profile", "userProfile"},
		{"xml_http_request", "xmlHttpRequest"},
		{"api_key", "apiKey"},
		{"", ""},
		{"a", "a"},
		{"a_b", "aB"},
		{"a_b_c", "aBC"},
		{"product_category", "productCategory"},
		{"order_item", "orderItem"},
		{"payment_method", "paymentMethod"},
		{"first_name", "firstName"},
		{"last_name", "lastName"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := generator.ToCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToCamelCase(%q) = %q, 期望 %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user", "User"},
		{"user_profile", "UserProfile"},
		{"user-profile", "UserProfile"},
		{"user profile", "UserProfile"},
		{"xml_http_request", "XmlHttpRequest"},
		{"api_key", "ApiKey"},
		{"", ""},
		{"a", "A"},
		{"a_b", "AB"},
		{"a_b_c", "ABC"},
		{"product_category", "ProductCategory"},
		{"order_item", "OrderItem"},
		{"payment_method", "PaymentMethod"},
		{"first_name", "FirstName"},
		{"last_name", "LastName"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := generator.ToPascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToPascalCase(%q) = %q, 期望 %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"User", "user"},
		{"UserProfile", "user-profile"},
		{"XMLHttpRequest", "xml-http-request"},
		{"APIKey", "api-key"},
		{"", ""},
		{"A", "a"},
		{"ProductCategory", "product-category"},
		{"OrderItem", "order-item"},
		{"PaymentMethod", "payment-method"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := generator.ToKebabCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToKebabCase(%q) = %q, 期望 %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToPlural(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// 规则变化
		{"user", "users"},
		{"product", "products"},
		{"category", "categories"},
		{"company", "companies"},
		{"box", "boxes"},
		{"class", "classes"},
		{"dish", "dishes"},
		{"church", "churches"},
		{"quiz", "quizes"},
		{"leaf", "leaves"},
		{"knife", "knives"},
		{"life", "lives"},
		{"wife", "wives"},

		// 不规则变化
		{"person", "people"},
		{"child", "children"},
		{"foot", "feet"},
		{"tooth", "teeth"},
		{"goose", "geese"},
		{"mouse", "mice"},
		{"man", "men"},
		{"woman", "women"},

		// 边界情况
		{"", ""},
		{"a", "as"},
		{"order", "orders"},
		{"item", "items"},
		{"status", "statuses"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := generator.ToPlural(tt.input)
			if result != tt.expected {
				t.Errorf("ToPlural(%q) = %q, 期望 %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidGoIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// 有效标识符
		{"User", true},
		{"user", true},
		{"_user", true},
		{"User123", true},
		{"user_name", true},
		{"userName", true},
		{"_", true},
		{"a", true},
		{"A", true},
		{"ID", true},
		{"userID", true},

		// 无效标识符
		{"", false},
		{"123user", false},
		{"user-name", false},
		{"user name", false},
		{"user.name", false},
		{"user@name", false},
		{"user#name", false},
		{"user$name", false},
		{"user%name", false},
		{"user&name", false},
		{"user*name", false},
		{"user+name", false},
		{"user=name", false},
		{"user!name", false},
		{"user?name", false},
		{"user<name", false},
		{"user>name", false},
		{"user,name", false},
		{"user;name", false},
		{"user:name", false},
		{"user'name", false},
		{"user\"name", false},
		{"user[name", false},
		{"user]name", false},
		{"user{name", false},
		{"user}name", false},
		{"user|name", false},
		{"user\\name", false},
		{"user/name", false},
		{"user`name", false},
		{"user~name", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := generator.IsValidGoIdentifier(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidGoIdentifier(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsGoKeyword(t *testing.T) {
	// Go关键字
	keywords := []string{
		"break", "case", "chan", "const", "continue",
		"default", "defer", "else", "fallthrough", "for",
		"func", "go", "goto", "if", "import",
		"interface", "map", "package", "range", "return",
		"select", "struct", "switch", "type", "var",
	}

	// 测试关键字
	for _, keyword := range keywords {
		t.Run(keyword, func(t *testing.T) {
			if !generator.IsGoKeyword(keyword) {
				t.Errorf("IsGoKeyword(%q) = false, 期望 true", keyword)
			}
		})
	}

	// 测试非关键字
	nonKeywords := []string{
		"User", "user", "name", "value", "data",
		"string", "int", "bool", "float64",
		"ID", "userID", "userName", "firstName",
	}

	for _, nonKeyword := range nonKeywords {
		t.Run(nonKeyword, func(t *testing.T) {
			if generator.IsGoKeyword(nonKeyword) {
				t.Errorf("IsGoKeyword(%q) = true, 期望 false", nonKeyword)
			}
		})
	}
}

func TestSanitizeGoIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// 有效标识符保持不变
		{"User", "User"},
		{"userName", "userName"},
		{"_private", "_private"},

		// 无效字符被移除
		{"user-name", "username"},
		{"user name", "username"},
		{"user.name", "username"},
		{"user@name", "username"},

		// 数字开头添加前缀
		{"123user", "F123user"},
		{"1", "F1"},

		// 空字符串返回默认值
		{"", "Field"},

		// 特殊字符组合
		{"user123name", "user123name"},
		{"_123", "_123"},

		// Go关键字添加后缀
		{"type", "typeField"},
		{"var", "varField"},
		{"func", "funcField"},
		{"if", "ifField"},
		{"for", "forField"},

		// 复杂情况
		{"123-user@name", "F123username"},
		{"@#$%", "F"},
		{"user!@#name", "username"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := generator.SanitizeGoIdentifier(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeGoIdentifier(%q) = %q, 期望 %q", tt.input, result, tt.expected)
			}

			// 验证结果是有效的Go标识符
			if !generator.IsValidGoIdentifier(result) {
				t.Errorf("SanitizeGoIdentifier(%q) = %q, 结果不是有效的Go标识符", tt.input, result)
			}

			// 验证结果不是Go关键字
			if generator.IsGoKeyword(result) {
				t.Errorf("SanitizeGoIdentifier(%q) = %q, 结果是Go关键字", tt.input, result)
			}
		})
	}
}
