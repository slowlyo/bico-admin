# 前端服务层架构

## 目录结构

服务层统一放在 `src/api/` 目录下，按照业务模块组织：

```
src/api/
├── auth.ts              # 认证相关（登录、登出、用户信息、验证码）
├── system-manage.ts     # 系统管理（用户、角色、菜单、部门）
├── demo.ts              # 示例模块
└── ...                  # 其他业务模块
```

## 请求工具

项目使用统一的请求工具 `src/utils/http`，它是基于 Axios 的二次封装，支持自动注入 Token、响应拦截、错误处理等功能。

### 基础用法

```typescript
import request from '@/utils/http'

/** 获取列表 */
export function fetchGetList(params: any) {
  return request.get<Api.Common.PaginatedResponse<any>>({
    url: '/admin-api/xxx',
    params
  })
}

/** 创建数据 */
export function fetchCreate(data: any) {
  return request.post({
    url: '/admin-api/xxx',
    data
  })
}
```

## API 命名规范

- **前缀**: 统一使用 `fetch` 开头（如 `fetchLogin`, `fetchGetUserList`）。
- **参数**: 复杂参数建议定义在 `src/types/api/` 下的命名空间中（如 `Api.Auth.LoginParams`）。
- **类型**: 尽量明确接口返回值类型，便于 TypeScript 类型推导。

## 常用模块说明

### 1. auth - 认证服务 (`src/api/auth.ts`)
- `fetchCaptcha()`: 获取验证码
- `fetchLogin()`: 登录
- `fetchLogout()`: 退出登录
- `fetchGetUserInfo()`: 获取当前用户信息
- `fetchUpdateProfile()`: 更新个人资料
- `fetchChangePassword()`: 修改密码

### 2. system-manage - 系统管理 (`src/api/system-manage.ts`)
- `fetchGetUserList()`: 用户管理相关
- `fetchGetRoleList()`: 角色管理相关
- `fetchGetAllRoles()`: 获取所有角色（下拉框使用）
- `fetchGetAllPermissions()`: 获取权限树

## 最佳实践

1. **统一管理**: 禁止在组件内直接写 Axios 请求，必须在 `api/` 目录下定义函数。
2. **类型安全**: 充分利用 `src/types/api/` 下定义的类型声明，确保前后端字段一致。
3. **错误处理**: `http` 工具已处理通用错误（如 401 自动跳转、业务错误弹窗），组件内只需关注特殊逻辑。
