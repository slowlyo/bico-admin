package handler

import "bico-admin/internal/admin/service"

// adminResponse 后台通用响应，仅用于 Swagger 文档。
type adminResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// loginRequest 登录请求参数。
type loginRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaID   string `json:"captchaId" binding:"required"`
	CaptchaCode string `json:"captchaCode" binding:"required"`
}

// captchaResponse 验证码响应数据。
type captchaResponse struct {
	ID    string `json:"id"`
	Image string `json:"image"`
}

// uploadResponse 上传响应数据。
type uploadResponse struct {
	URL string `json:"url"`
}

// rolePermissionsResponse 角色权限响应数据。
type rolePermissionsResponse struct {
	Permissions []string `json:"permissions"`
}

// permissionDocItem 权限树节点，仅用于 Swagger 文档。
type permissionDocItem struct {
	Key      string              `json:"key"`
	Label    string              `json:"label"`
	Children []permissionDocItem `json:"children,omitempty"`
}

// adminRoleDocItem 角色简要信息，仅用于 Swagger 文档。
type adminRoleDocItem struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Code        string   `json:"code"`
	Description string   `json:"description"`
	Enabled     bool     `json:"enabled"`
	Permissions []string `json:"permissions"`
}

// demoExcelImportResponse Excel 导入示例响应数据。
type demoExcelImportResponse struct {
	Total   int        `json:"total"`
	Preview [][]string `json:"preview"`
}

// appConfigDocResponse 应用配置文档响应。
type appConfigDocResponse = service.AppConfigResponse

// loginDocResponse 登录文档响应。
type loginDocResponse = service.LoginResponse

// currentUserDocResponse 当前用户文档响应。
type currentUserDocResponse = service.UserInfo
