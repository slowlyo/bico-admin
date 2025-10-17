package handler

import (
	"net/http"
)

type AuthServiceInterface interface {
	Login(req interface{}) (interface{}, error)
	Logout(token string) error
}

// AuthHandler 认证处理器
type AuthHandler struct {
	authService AuthServiceInterface
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 登录接口
func (h *AuthHandler) Login(c interface {
	ShouldBindJSON(obj interface{}) error
	JSON(code int, obj interface{})
}) {
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
			"code": 401,
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
func (h *AuthHandler) Logout(c interface {
	GetHeader(key string) string
	JSON(code int, obj interface{})
}) {
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
