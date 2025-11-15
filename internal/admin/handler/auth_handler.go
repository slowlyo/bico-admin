package handler

import (
	"bico-admin/internal/admin/service"
	"bico-admin/internal/core/upload"
	"bico-admin/internal/pkg/captcha"
	"bico-admin/internal/pkg/response"

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
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证验证码
	if !h.captcha.Verify(req.CaptchaID, req.CaptchaCode) {
		response.ErrorWithCode(c, 400, "验证码错误")
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		response.ErrorWithCode(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, "登录成功", resp)
}

// Logout 退出登录接口
func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "" && len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	h.authService.Logout(token)
	response.SuccessWithMessage(c, "退出成功", nil)
}

// CurrentUser 获取当前用户信息
func (h *AuthHandler) CurrentUser(c *gin.Context) {
	// 从上下文中获取用户信息（由 JWT 中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		userID = nil
	}
	if userID == nil {
		response.ErrorWithCode(c, 401, "未授权")
		return
	}

	user, err := h.authService.GetUserByID(userID.(uint))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.SuccessWithData(c, user)
}

// UpdateProfile 更新用户资料
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists || userID == nil {
		response.ErrorWithCode(c, 401, "未授权")
		return
	}

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.authService.UpdateProfile(userID.(uint), &req)
	if err != nil {
		response.ErrorWithCode(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", user)
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists || userID == nil {
		response.ErrorWithCode(c, 401, "未授权")
		return
	}

	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	err := h.authService.ChangePassword(userID.(uint), &req)
	if err != nil {
		response.ErrorWithCode(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}

// UploadAvatar 上传头像
func (h *AuthHandler) UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists || userID == nil {
		response.ErrorWithCode(c, 401, "未授权")
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		response.BadRequest(c, "请上传头像文件")
		return
	}

	url, err := h.uploader.Upload(file, "avatars")
	if err != nil {
		response.ErrorWithCode(c, 400, err.Error())
		return
	}

	response.SuccessWithMessage(c, "上传成功", gin.H{"url": url})
}

// GetCaptcha 获取验证码
func (h *AuthHandler) GetCaptcha(c *gin.Context) {
	id, b64s, err := h.captcha.Generate()
	if err != nil {
		response.ErrorWithCode(c, 500, "生成验证码失败")
		return
	}

	response.SuccessWithData(c, gin.H{
		"id":    id,
		"image": b64s,
	})
}
