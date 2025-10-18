package handler

import (
	"net/http"

	"bico-admin/internal/admin/service"
	"bico-admin/internal/core/upload"
	"bico-admin/internal/pkg/captcha"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.IAuthService
	uploader    upload.Uploader
	captcha     *captcha.Captcha
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService service.IAuthService, uploader upload.Uploader, cap *captcha.Captcha) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		uploader:    uploader,
		captcha:     cap,
	}
}

// Login 登录接口
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		CaptchaID   string `json:"captchaId" binding:"required"`
		CaptchaCode string `json:"captchaCode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 验证验证码
	if !h.captcha.Verify(req.CaptchaID, req.CaptchaCode) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 400,
			"msg":  "验证码错误",
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

// UpdateProfile 更新用户资料
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists || userID == nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	user, err := h.authService.UpdateProfile(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "更新成功",
		"data": user,
	})
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists || userID == nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}

	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	err := h.authService.ChangePassword(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "密码修改成功",
	})
}

// UploadAvatar 上传头像
func (h *AuthHandler) UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists || userID == nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": 400,
			"msg":  "请上传头像文件",
		})
		return
	}

	url, err := h.uploader.Upload(file, "avatars")
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "上传成功",
		"data": map[string]interface{}{
			"url": url,
		},
	})
}

// GetCaptcha 获取验证码
func (h *AuthHandler) GetCaptcha(c *gin.Context) {
	id, b64s, err := h.captcha.Generate()
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 500,
			"msg":  "生成验证码失败",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": map[string]interface{}{
			"id":    id,
			"image": b64s,
		},
	})
}
