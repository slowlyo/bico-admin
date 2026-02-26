# 后台认证 API

本文档基于当前代码实现（`internal/admin/router.go`、`auth_handler.go`、`auth_service.go`）。

## 基本信息

- API 前缀：`/admin-api`
- 认证方式：`Authorization: Bearer {token}`
- 默认管理员：首次迁移自动创建 `admin / admin`

## 响应格式

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

说明：

- 大多数业务错误返回 HTTP 200，`code != 0`
- 参数绑定错误使用 HTTP 400（`response.BadRequest`）

## 接口清单

### 公开接口

1. `GET /admin-api/captcha` 获取验证码
2. `POST /admin-api/auth/login` 登录
3. `GET /admin-api/app-config` 获取应用配置

### 需登录接口

1. `POST /admin-api/auth/logout` 退出登录
2. `GET /admin-api/auth/current-user` 获取当前用户
3. `PUT /admin-api/auth/profile` 更新资料
4. `PUT /admin-api/auth/password` 修改密码
5. `POST /admin-api/auth/avatar` 上传头像

## 关键接口说明

### 1) 获取验证码

`GET /admin-api/captcha`

成功示例：

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": "captcha-id",
    "image": "data:image/png;base64,..."
  }
}
```

### 2) 登录

`POST /admin-api/auth/login`

请求体：

```json
{
  "username": "admin",
  "password": "admin",
  "captchaId": "captcha-id",
  "captchaCode": "abcd"
}
```

成功示例：

```json
{
  "code": 0,
  "msg": "登录成功",
  "data": {
    "token": "eyJhbGciOi..."
  }
}
```

常见失败：

- `code=400, msg=验证码错误`
- `code=400, msg=用户不存在`
- `code=400, msg=密码错误`
- `code=400, msg=用户已被禁用`

### 3) 获取当前用户

`GET /admin-api/auth/current-user`

请求头：

```text
Authorization: Bearer <token>
```

成功示例：

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "name": "系统管理员",
    "avatar": "https://...",
    "enabled": true,
    "permissions": ["system:admin_user:list"]
  }
}
```

未登录：`code=401, msg=未授权`

### 4) 退出登录

`POST /admin-api/auth/logout`

说明：

- 会将 token 写入黑名单缓存，默认保留 7 天。

### 5) 更新资料

`PUT /admin-api/auth/profile`

请求体：

```json
{
  "name": "新名字",
  "avatar": "https://example.com/avatar.png"
}
```

### 6) 修改密码

`PUT /admin-api/auth/password`

请求体：

```json
{
  "oldPassword": "old",
  "newPassword": "new"
}
```

### 7) 上传头像

`POST /admin-api/auth/avatar`

- `Content-Type: multipart/form-data`
- 字段名：`avatar`

## 本地联调示例

```bash
# 1. 执行迁移
make migrate

# 2. 启动服务
make serve

# 3. 获取验证码
curl -s "http://localhost:8080/admin-api/captcha"
```

登录 curl 示例（需替换验证码）：

```bash
curl -X POST "http://localhost:8080/admin-api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin","captchaId":"xxx","captchaCode":"xxxx"}'
```
