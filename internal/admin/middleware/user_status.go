package middleware

import (
	"bico-admin/internal/admin/model"
	"bico-admin/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserStatusMiddleware 用户状态检查中间件
type UserStatusMiddleware struct {
	db *gorm.DB
}

// NewUserStatusMiddleware 创建用户状态中间件
func NewUserStatusMiddleware(db *gorm.DB) *UserStatusMiddleware {
	return &UserStatusMiddleware{db: db}
}

// Check 检查用户是否被禁用
func (m *UserStatusMiddleware) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		var user model.AdminUser
		if err := m.db.Select("enabled").First(&user, userID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response.ErrorWithCode(c, 401, "用户不存在")
				c.Abort()
				return
			}
			response.ErrorWithCode(c, 500, "查询用户状态失败")
			c.Abort()
			return
		}

		if !user.Enabled {
			response.ErrorWithCode(c, 401, "账户已被禁用")
			c.Abort()
			return
		}

		c.Next()
	}
}
