package middleware

import (
	"bico-admin/internal/admin/service"
	"bico-admin/internal/pkg/response"
	"errors"

	"github.com/gin-gonic/gin"
)

// UserStatusMiddleware 用户状态检查中间件
type UserStatusMiddleware struct {
	userService interface {
		IsUserEnabled(userID uint) (bool, error)
	}
}

// NewUserStatusMiddleware 创建用户状态中间件
func NewUserStatusMiddleware(userService interface {
	IsUserEnabled(userID uint) (bool, error)
}) *UserStatusMiddleware {
	return &UserStatusMiddleware{userService: userService}
}

// Check 检查用户是否被禁用
func (m *UserStatusMiddleware) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		// 未登录请求直接放行，由后续鉴权中间件处理。
		if !exists {
			c.Next()
			return
		}

		uid, ok := userID.(uint)
		// user_id 类型异常时按未授权处理，避免 panic。
		if !ok {
			response.ErrorWithCode(c, 401, "未授权")
			c.Abort()
			return
		}

		enabled, err := m.userService.IsUserEnabled(uid)
		if err != nil {
			// 用户不存在时返回 401，其他错误按服务错误处理。
			if errors.Is(err, service.ErrUserNotFound) {
				response.ErrorWithCode(c, 401, "用户不存在")
				c.Abort()
				return
			}
			response.ErrorWithCode(c, 500, "查询用户状态失败")
			c.Abort()
			return
		}

		// 用户被禁用时拒绝访问。
		if !enabled {
			response.ErrorWithCode(c, 401, "账户已被禁用")
			c.Abort()
			return
		}

		c.Next()
	}
}
