# 用户信息编辑和修改密码功能

## 功能概述

本功能实现了用户个人资料的编辑和密码修改功能，包括：

1. **用户信息编辑**：用户可以修改自己的基本信息（用户名、邮箱、昵称、手机号等）
2. **修改密码**：用户可以安全地修改自己的登录密码
3. **头像上传**：用户可以上传和更换个人头像（前端界面已准备，后端接口待实现）

## 后端实现

### API接口

#### 1. 获取用户资料
```
GET /auth/profile
Authorization: Bearer {token}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@bico-admin.com",
    "nickname": "管理员",
    "avatar": "",
    "phone": "",
    "status": 1,
    "last_login_at": "2024-01-01T10:00:00Z",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z",
    "roles": []
  }
}
```

#### 2. 更新用户资料
```
PUT /auth/profile
Authorization: Bearer {token}
Content-Type: application/json
```

**请求体：**
```json
{
  "username": "new_username",
  "email": "new_email@example.com",
  "nickname": "新昵称",
  "phone": "13800138000"
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "Profile updated successfully",
  "data": {
    "id": 1,
    "username": "new_username",
    "email": "new_email@example.com",
    "nickname": "新昵称",
    "phone": "13800138000",
    // ... 其他字段
  }
}
```

#### 3. 修改密码
```
POST /auth/change-password
Authorization: Bearer {token}
Content-Type: application/json
```

**请求体：**
```json
{
  "old_password": "current_password",
  "new_password": "new_password"
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "Password changed successfully",
  "data": null
}
```

### 数据模型

#### UserUpdateRequest
```go
type UserUpdateRequest struct {
    Username string     `json:"username" validate:"min=3,max=50"`
    Email    string     `json:"email" validate:"email"`
    Nickname string     `json:"nickname" validate:"max=50"`
    Phone    string     `json:"phone" validate:"max=20"`
    Status   UserStatus `json:"status" validate:"oneof=0 1 2"`
    RoleIDs  []uint     `json:"role_ids"`
}
```

#### UserChangePasswordRequest
```go
type UserChangePasswordRequest struct {
    OldPassword string `json:"old_password" validate:"required"`
    NewPassword string `json:"new_password" validate:"required,min=6"`
}
```

### 安全特性

1. **身份验证**：所有接口都需要有效的JWT token
2. **权限验证**：用户只能修改自己的信息
3. **密码验证**：修改密码时需要验证当前密码
4. **数据验证**：使用validator进行输入数据验证
5. **唯一性检查**：用户名和邮箱的唯一性验证
6. **密码加密**：使用bcrypt加密存储密码

## 前端实现

### 页面路由
- `/profile` - 个人资料页面

### 主要组件

#### Profile 组件
位置：`frontend/src/pages/profile/index.tsx`

**功能特性：**
1. **标签页设计**：基本信息和修改密码分别在不同标签页
2. **表单验证**：前端表单验证确保数据格式正确
3. **头像上传**：支持头像上传（界面已实现，后端接口待开发）
4. **实时反馈**：操作成功/失败的消息提示
5. **响应式布局**：适配不同屏幕尺寸

#### Header 组件更新
位置：`frontend/src/components/header/index.tsx`

**新增功能：**
- 用户下拉菜单中添加"个人资料"链接
- 点击可跳转到个人资料页面

### API调用

#### authAPI 扩展
位置：`frontend/src/utils/request.ts`

**新增方法：**
```typescript
// 更新用户资料
updateProfile: async (data: Partial<UserInfo>): Promise<ApiResponse<UserInfo>>

// 修改密码  
changePassword: async (data: { old_password: string; new_password: string }): Promise<ApiResponse<null>>
```

## 使用指南

### 用户操作流程

1. **访问个人资料页面**
   - 点击右上角用户头像
   - 选择"个人资料"菜单项

2. **编辑基本信息**
   - 在"基本信息"标签页修改个人信息
   - 点击"保存修改"按钮提交

3. **修改密码**
   - 切换到"修改密码"标签页
   - 输入当前密码和新密码
   - 点击"修改密码"按钮提交

### 开发者测试

#### 后端API测试
```bash
# 进入后端目录
cd backend

# 运行测试脚本
./test_user_profile.sh
```

#### 前端功能测试
1. 启动前端开发服务器
2. 登录系统
3. 访问个人资料页面测试各项功能

## 技术实现细节

### 后端架构
- **Handler层**：直接处理HTTP请求和业务逻辑
- **数据库操作**：使用通用的database.Operations工具
- **错误处理**：统一的错误响应格式
- **中间件**：JWT认证中间件保护接口

### 前端架构
- **Refine框架**：使用useGetIdentity获取用户信息
- **Ant Design**：UI组件库提供表单和布局
- **React Router**：页面路由管理
- **Axios**：HTTP请求处理

### 数据流程
```
用户操作 → 前端表单 → API请求 → 后端验证 → 数据库更新 → 响应返回 → 前端更新
```

## 扩展功能

### 待实现功能
1. **头像上传**：文件上传接口和存储
2. **操作日志**：记录用户信息修改历史
3. **邮箱验证**：修改邮箱时发送验证邮件
4. **手机验证**：修改手机号时发送验证短信

### 优化建议
1. **缓存优化**：用户信息缓存减少数据库查询
2. **批量更新**：支持批量修改多个字段
3. **版本控制**：用户信息变更版本管理
4. **审计日志**：详细的操作审计记录

## 注意事项

1. **安全性**：确保只有用户本人可以修改自己的信息
2. **数据一致性**：用户名和邮箱的唯一性约束
3. **密码强度**：前后端都要验证密码强度
4. **错误处理**：友好的错误提示和异常处理
5. **性能考虑**：避免频繁的数据库查询和更新
