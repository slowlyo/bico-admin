package handler

import (
	"context"

	"bico-admin/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type ctxKey int

const ctxKeyOperatorUserID ctxKey = 1

// withOperatorUserID 把当前操作者 ID 写入请求上下文，供事务内 hook 读取。
func withOperatorUserID(ctx context.Context, operatorID uint) context.Context {
	return context.WithValue(ctx, ctxKeyOperatorUserID, operatorID)
}

// operatorUserIDFromContext 从上下文读取操作者 ID；缺失时返回 0。
func operatorUserIDFromContext(ctx context.Context) uint {
	if ctx == nil {
		return 0
	}
	v, ok := ctx.Value(ctxKeyOperatorUserID).(uint)
	if !ok {
		return 0
	}
	return v
}

// mustOperatorUserID 从 Gin 上下文读取当前登录用户 ID，失败时直接写 401 响应。
func mustOperatorUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists || userID == nil {
		response.ErrorWithCode(c, 401, "未授权")
		return 0, false
	}
	uid, ok := userID.(uint)
	if !ok || uid == 0 {
		response.ErrorWithCode(c, 401, "未授权")
		return 0, false
	}
	return uid, true
}
