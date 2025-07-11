package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"bico-admin/internal/api/types"
)

// APIAuth API认证中间件
func APIAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 获取API Key
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			// 也可以从查询参数获取
			apiKey = c.Query("api_key")
		}

		if apiKey == "" {
			c.JSON(401, types.APIResponse{
				Code:      401,
				Message:   "缺少API密钥",
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		// TODO: 验证API Key
		keyInfo, err := validateAPIKey(apiKey)
		if err != nil {
			c.JSON(401, types.APIResponse{
				Code:      401,
				Message:   "API密钥无效",
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		// 检查权限
		if !hasAPIPermission(keyInfo, c.Request.Method, c.FullPath()) {
			c.JSON(403, types.APIResponse{
				Code:      403,
				Message:   "权限不足",
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		// 将API Key信息存储到上下文
		c.Set("api_key_id", keyInfo.ID)
		c.Set("api_permissions", keyInfo.Permissions)

		c.Next()
	})
}

// APIKeyInfo API密钥信息
type APIKeyInfo struct {
	ID          uint
	Permissions []string
	Status      int
}

// validateAPIKey 验证API密钥（临时实现）
func validateAPIKey(apiKey string) (*APIKeyInfo, error) {
	// TODO: 实现真正的API密钥验证
	// 这里是临时的mock实现
	if strings.HasPrefix(apiKey, "bico_") {
		return &APIKeyInfo{
			ID:          1,
			Permissions: []string{"user:read", "user:write", "stats:read"},
			Status:      1,
		}, nil
	}
	
	return nil, gin.Error{Err: gin.Error{}.Err, Type: gin.ErrorTypePublic}
}

// hasAPIPermission 检查API权限
func hasAPIPermission(keyInfo *APIKeyInfo, method, path string) bool {
	// TODO: 实现真正的权限检查逻辑
	// 这里是临时的简单实现
	
	// 构建权限字符串
	var permission string
	switch {
	case strings.Contains(path, "/users"):
		if method == "GET" {
			permission = "user:read"
		} else {
			permission = "user:write"
		}
	case strings.Contains(path, "/stats"):
		permission = "stats:read"
	default:
		return true // 默认允许
	}

	// 检查权限
	for _, p := range keyInfo.Permissions {
		if p == permission {
			return true
		}
	}

	return false
}
