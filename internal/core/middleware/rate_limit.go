package middleware

import (
	"bico-admin/internal/pkg/response"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 按 key 保存令牌桶，参数在每次请求时读取。
type RateLimiter struct {
	limiters    map[string]*rate.Limiter
	mu          sync.RWMutex
	settings    func() (enabled bool, rps int, burst int)
	cleanupOnce sync.Once
}

// NewRateLimiter 创建限流器。
//
// 说明：settings 返回当前是否启用以及速率，供配置热更新后立即生效。
func NewRateLimiter(settings func() (enabled bool, rps int, burst int)) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		settings: settings,
	}
}

// getLimiter 获取或创建限流器，并把已有桶调整到当前速率。
func (rl *RateLimiter) getLimiter(key string, rps int, burst int) *rate.Limiter {
	limit := rate.Limit(rps)

	rl.mu.RLock()
	limiter, exists := rl.limiters[key]
	rl.mu.RUnlock()
	if exists {
		limiter.SetLimit(limit)
		limiter.SetBurst(burst)
		return limiter
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// 拿到写锁后再看一次，避免并发下为同一个 key 建两个桶。
	if limiter, exists = rl.limiters[key]; exists {
		limiter.SetLimit(limit)
		limiter.SetBurst(burst)
		return limiter
	}

	limiter = rate.NewLimiter(limit, burst)
	rl.limiters[key] = limiter
	return limiter
}

// allow 按当前配置决定这次请求是否放行。
//
// 说明：未启用或参数无效时直接放行，不为每个 IP 建桶。
func (rl *RateLimiter) allow(key string) bool {
	if rl.settings == nil {
		return true
	}
	enabled, rps, burst := rl.settings()
	if !enabled || rps <= 0 || burst <= 0 {
		return true
	}

	rl.startCleanup(5 * time.Minute)
	return rl.getLimiter(key, rps, burst).Allow()
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
			// 只删满桶的空闲项。活跃桶继续保留，避免整表清空后把正常用户打回突发额度。
			if len(rl.limiters) > 10000 {
				for key, limiter := range rl.limiters {
					if limiter.Tokens() >= float64(limiter.Burst()) {
						delete(rl.limiters, key)
					}
				}
			}
			rl.mu.Unlock()
		}
	}()
}

// RateLimit 限流中间件（基于IP）
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.allow(c.ClientIP()) {
			response.TooManyRequests(c, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByUser 限流中间件（基于用户ID）
func (rl *RateLimiter) RateLimitByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 未登录时没有用户 ID，退回 IP，避免匿名流量共用一个桶。
		key := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			key = buildUserRateLimitKey(userID)
		}
		if !rl.allow(key) {
			response.TooManyRequests(c, "请求过于频繁，请稍后再试")
			c.Abort()
			return
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
	return func(c *gin.Context) {
		if !rl.allow(keyFunc(c)) {
			response.TooManyRequests(c, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
