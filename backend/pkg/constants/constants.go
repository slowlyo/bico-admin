package constants

// 应用常量
const (
	AppName    = "Bico Admin"
	AppVersion = "1.0.0"
)

// 默认值常量
const (
	DefaultPageSize = 10
	MaxPageSize     = 100
	DefaultPage     = 1
)

// 缓存键前缀
const (
	CacheKeyUser        = "user:"
	CacheKeyRole        = "role:"
	CacheKeyPermission  = "permission:"
	CacheKeyToken       = "token:"
	CacheKeyUserRoles   = "user_roles:"
	CacheKeyRolePerms   = "role_permissions:"
)

// 缓存过期时间（秒）
const (
	CacheExpireShort  = 300   // 5分钟
	CacheExpireMedium = 1800  // 30分钟
	CacheExpireLong   = 3600  // 1小时
	CacheExpireDay    = 86400 // 24小时
)

// 文件上传相关
const (
	MaxFileSize     = 10 << 20 // 10MB
	AllowedImageExt = "jpg,jpeg,png,gif,webp"
	AllowedDocExt   = "pdf,doc,docx,xls,xlsx,ppt,pptx,txt"
	AllowedVideoExt = "mp4,avi,mov,wmv,flv,webm"
)

// 用户状态
const (
	UserStatusInactive = 0
	UserStatusActive   = 1
	UserStatusBlocked  = 2
)

// 角色状态
const (
	RoleStatusInactive = 0
	RoleStatusActive   = 1
)

// 权限状态
const (
	PermissionStatusInactive = 0
	PermissionStatusActive   = 1
)

// 权限类型
const (
	PermissionTypeMenu   = 1
	PermissionTypeButton = 2
	PermissionTypeAPI    = 3
)

// 默认角色代码
const (
	RoleCodeSuperAdmin = "super_admin"
	RoleCodeAdmin      = "admin"
	RoleCodeUser       = "user"
)

// 默认权限代码
const (
	PermissionCodeUserManage       = "user:manage"
	PermissionCodeUserCreate       = "user:create"
	PermissionCodeUserUpdate       = "user:update"
	PermissionCodeUserDelete       = "user:delete"
	PermissionCodeUserView         = "user:view"
	PermissionCodeRoleManage       = "role:manage"
	PermissionCodePermissionManage = "permission:manage"
)

// HTTP状态码
const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusConflict            = 409
	StatusUnprocessableEntity = 422
	StatusInternalServerError = 500
)

// 错误消息
const (
	ErrInvalidCredentials = "Invalid credentials"
	ErrUserNotFound       = "User not found"
	ErrUserExists         = "User already exists"
	ErrRoleNotFound       = "Role not found"
	ErrPermissionDenied   = "Permission denied"
	ErrTokenExpired       = "Token expired"
	ErrTokenInvalid       = "Token invalid"
	ErrValidationFailed   = "Validation failed"
	ErrInternalServer     = "Internal server error"
)

// 成功消息
const (
	MsgSuccess      = "Success"
	MsgCreated      = "Created successfully"
	MsgUpdated      = "Updated successfully"
	MsgDeleted      = "Deleted successfully"
	MsgLoginSuccess = "Login successful"
	MsgLogoutSuccess = "Logout successful"
)
