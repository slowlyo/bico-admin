package handler

import (
	"net/http"
	
	"bico-admin/internal/admin/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.IAuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService service.IAuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 登录接口
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}
	
	resp, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "登录成功",
		"data": resp,
	})
}

// Logout 退出登录接口
func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "" && len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	
	h.authService.Logout(token)
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "退出成功",
	})
}

// CurrentUser 获取当前用户信息
func (h *AuthHandler) CurrentUser(c *gin.Context) {
	// 从上下文中获取用户信息（由 JWT 中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		userID = nil
	}
	if userID == nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}
	
	user, err := h.authService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}
	
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": user,
	})
}
