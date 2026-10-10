package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestRateLimitDisabledAllowsAll 验证关闭限流时不为请求计数。
func TestRateLimitDisabledAllowsAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	limiter := NewRateLimiter(func() (bool, int, int) {
		return false, 1, 1
	})
	engine := gin.New()
	engine.Use(limiter.RateLimit())
	engine.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	for i := 0; i < 3; i++ {
		request := httptest.NewRequest(http.MethodGet, "/test", nil)
		response := httptest.NewRecorder()
		engine.ServeHTTP(response, request)
		if response.Code != http.StatusNoContent {
			t.Fatalf("关闭限流后第 %d 次请求被拒绝", i+1)
		}
	}
}

// TestRateLimitUsesLatestSettings 验证修改返回值后下一次请求就按新额度限制。
func TestRateLimitUsesLatestSettings(t *testing.T) {
	gin.SetMode(gin.TestMode)
	enabled := false
	limiter := NewRateLimiter(func() (bool, int, int) {
		return enabled, 1, 1
	})
	engine := gin.New()
	engine.Use(limiter.RateLimit())
	engine.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("启用前请求应放行")
	}

	enabled = true
	first := httptest.NewRecorder()
	engine.ServeHTTP(first, request)
	if first.Code != http.StatusNoContent {
		t.Fatalf("启用后首个请求应放行")
	}
	second := httptest.NewRecorder()
	engine.ServeHTTP(second, request)
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("超过突发额度后应返回 429，实际 %d", second.Code)
	}
}
