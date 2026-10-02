package middleware

import (
	"bico-admin/internal/pkg/jwt"
	"bico-admin/internal/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// Session 一次读取到的账号状态，供同一次请求里的后续中间件复用。
type Session struct {
	Found     bool
	Enabled   bool
	VersionOK bool
}

// CtxUserEnabled 标记本次请求已经确认账号启用，避免状态中间件再查一次。
const CtxUserEnabled = "user_enabled"

// CtxUserPermissions 保存本次请求的权限集合。
const CtxUserPermissions = "user_permissions"

// JWTAuth JWT认证中间件
func JWTAuth(jwtManager *jwt.JWTManager, authService interface {
	IsTokenBlacklisted(token string) bool
	LoadSession(userID uint, version uint) (Session, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorWithCode(c, 401, "请先登录")
			c.Abort()
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorWithCode(c, 401, "token 格式错误")
			c.Abort()
			return
		}

		token := parts[1]

		// 验证 token
		claims, err := jwtManager.ParseToken(token)

		if err != nil {
			response.ErrorWithCode(c, 401, "token 无效或已过期")
			c.Abort()
			return
		}
		if claims == nil || claims.UserID == 0 || claims.Username == "" {
			response.ErrorWithCode(c, 401, "token 无效")
			c.Abort()
			return
		}
		if authService.IsTokenBlacklisted(token) {
			response.ErrorWithCode(c, 401, "token 已失效")
			c.Abort()
			return
		}

		session, err := authService.LoadSession(claims.UserID, claims.Version)
		if err != nil {
			response.ErrorWithCode(c, 500, "鉴权失败")
			c.Abort()
			return
		}
		// 用户已删除与令牌被主动作废要分开，避免把不存在的账号说成令牌过期。
		if !session.Found {
			response.ErrorWithCode(c, 401, "用户不存在")
			c.Abort()
			return
		}
		if !session.VersionOK {
			response.ErrorWithCode(c, 401, "token 已失效")
			c.Abort()
			return
		}
		if !session.Enabled {
			response.ErrorWithCode(c, 401, "账户已被禁用")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set(CtxUserEnabled, true)

		c.Next()
	}
}
