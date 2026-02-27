package middleware

import (
	"bico-admin/internal/pkg/response"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 限流器接口
type RateLimiter struct {
	limiters    map[string]*rate.Limiter
	mu          sync.RWMutex
	rate        rate.Limit
	burst       int
	cleanupOnce sync.Once
}

// NewRateLimiter 创建限流器
// rps: 每秒请求数
// burst: 突发流量桶容量
func NewRateLimiter(rps int, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(rps),
		burst:    burst,
	}
}

// getLimiter 获取或创建限流器
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[key]
	rl.mu.RUnlock()

	if exists {
		return limiter
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// 双重检查
	if limiter, exists := rl.limiters[key]; exists {
		return limiter
	}

	limiter = rate.NewLimiter(rl.rate, rl.burst)
	rl.limiters[key] = limiter

	return limiter
}

// startCleanup 启动限流器清理协程（仅启动一次）
func (rl *RateLimiter) startCleanup(interval time.Duration) {
	rl.cleanupOnce.Do(func() {
		rl.cleanupExpiredLimiters(interval)
	})
}

// cleanupExpiredLimiters 清理过期限流器（定期清理避免内存泄漏）
func (rl *RateLimiter) cleanupExpiredLimiters(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			rl.mu.Lock()
			// 简单清空策略，实际可根据最后访问时间清理
			if len(rl.limiters) > 10000 {
				rl.limiters = make(map[string]*rate.Limiter)
			}
			rl.mu.Unlock()
		}
	}()
}

// RateLimit 限流中间件（基于IP）
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	// 启动清理协程（只会启动一次）
	rl.startCleanup(5 * time.Minute)

	return func(c *gin.Context) {
		// 获取客户端IP
		ip := c.ClientIP()
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			response.TooManyRequests(c, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByUser 限流中间件（基于用户ID）
func (rl *RateLimiter) RateLimitByUser() gin.HandlerFunc {
	rl.startCleanup(5 * time.Minute)

	return func(c *gin.Context) {
		// 从JWT中间件获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			// 未认证用户使用IP限流
			ip := c.ClientIP()
			limiter := rl.getLimiter(ip)
			if !limiter.Allow() {
				response.TooManyRequests(c, "请求过于频繁，请稍后再试")
				c.Abort()
				return
			}
		} else {
			// 已认证用户使用用户ID限流
			key := buildUserRateLimitKey(userID)
			limiter := rl.getLimiter(key)
			if !limiter.Allow() {
				response.TooManyRequests(c, "请求过于频繁，请稍后再试")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// buildUserRateLimitKey 构建用户维度限流 key
func buildUserRateLimitKey(userID interface{}) string {
	switch v := userID.(type) {
	case string:
		return "user:" + v
	case uint:
		return fmt.Sprintf("user:%d", v)
	case uint64:
		return fmt.Sprintf("user:%d", v)
	case int:
		return fmt.Sprintf("user:%d", v)
	case int64:
		return fmt.Sprintf("user:%d", v)
	default:
		// 未知类型兜底转字符串，避免类型断言 panic。
		return fmt.Sprintf("user:%v", v)
	}
}

// RateLimitByKey 限流中间件（基于自定义key）
func (rl *RateLimiter) RateLimitByKey(keyFunc func(*gin.Context) string) gin.HandlerFunc {
	rl.startCleanup(5 * time.Minute)

	return func(c *gin.Context) {
		key := keyFunc(c)
		limiter := rl.getLimiter(key)

		if !limiter.Allow() {
			response.TooManyRequests(c, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
