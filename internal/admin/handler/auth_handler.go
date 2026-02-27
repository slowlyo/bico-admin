package handler

import (
	"bico-admin/internal/admin/service"
	"bico-admin/internal/core/upload"
	"bico-admin/internal/pkg/captcha"
	"bico-admin/internal/pkg/crud"
	"bico-admin/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	crud.BaseHandler
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

	if err := h.BindJSON(c, &req); err != nil {
		return
	}

	// 验证验证码
	if !h.captcha.Verify(req.CaptchaID, req.CaptchaCode) {
		h.Error(c, "验证码错误")
		return
	}

	// 验证码校验通过后，仅将账号密码传入服务层。
	loginReq := &service.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	resp, err := h.authService.Login(loginReq)
	if err != nil {
		h.Error(c, err.Error())
		return
	}

	h.SuccessWithMessage(c, "登录成功", resp)
}

// Logout 退出登录接口
func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	// Authorization 为 Bearer 时提取 token 正文。
	if token != "" && len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	h.authService.Logout(token)
	h.SuccessWithMessage(c, "退出成功", nil)
}

// CurrentUser 获取当前用户信息
func (h *AuthHandler) CurrentUser(c *gin.Context) {
	userID, ok := h.mustUserID(c)
	if !ok {
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		h.NotFound(c, "用户不存在")
		return
	}

	h.Success(c, user)
}

// UpdateProfile 更新用户资料
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, ok := h.mustUserID(c)
	if !ok {
		return
	}

	var req service.UpdateProfileRequest
	if err := h.BindJSON(c, &req); err != nil {
		return
	}

	user, err := h.authService.UpdateProfile(userID, &req)
	if err != nil {
		h.Error(c, err.Error())
		return
	}

	h.SuccessWithMessage(c, "更新成功", user)
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, ok := h.mustUserID(c)
	if !ok {
		return
	}

	var req service.ChangePasswordRequest
	if err := h.BindJSON(c, &req); err != nil {
		return
	}

	err := h.authService.ChangePassword(userID, &req)
	if err != nil {
		h.Error(c, err.Error())
		return
	}

	h.SuccessWithMessage(c, "密码修改成功", nil)
}

// UploadAvatar 上传头像
func (h *AuthHandler) UploadAvatar(c *gin.Context) {
	_, ok := h.mustUserID(c)
	if !ok {
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		response.BadRequest(c, "请上传头像文件")
		return
	}

	url, err := h.uploader.Upload(file, "avatars")
	if err != nil {
		h.Error(c, err.Error())
		return
	}

	h.SuccessWithMessage(c, "上传成功", gin.H{"url": url})
}

// GetCaptcha 获取验证码
func (h *AuthHandler) GetCaptcha(c *gin.Context) {
	id, b64s, err := h.captcha.Generate()
	if err != nil {
		response.ErrorWithCode(c, 500, "生成验证码失败")
		return
	}

	h.Success(c, gin.H{
		"id":    id,
		"image": b64s,
	})
}

// mustUserID 获取当前登录用户 ID。
//
// 说明：该方法统一处理 user_id 读取与类型校验，避免每个接口重复写 401 分支。
func (h *AuthHandler) mustUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	// 未登录或上下文缺失时统一返回 401。
	if !exists || userID == nil {
		response.ErrorWithCode(c, 401, "未授权")
		return 0, false
	}

	uid, ok := userID.(uint)
	// user_id 类型异常时按未授权处理，避免 panic。
	if !ok {
		response.ErrorWithCode(c, 401, "未授权")
		return 0, false
	}

	return uid, true
}
