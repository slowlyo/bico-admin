# 前端服务层架构

## 目录结构

服务层按照**菜单层级**组织，结构清晰，便于维护：

```
src/services/
├── config.ts              # API配置（前缀、工具函数）
├── auth/                  # 认证相关服务
│   ├── index.ts          # 登录、登出、获取当前用户、验证码
│   ├── profile.ts        # 个人资料、修改密码、上传头像
│   └── types.ts          # 类型定义
├── common/                # 公共服务
│   ├── index.ts          # 应用配置
│   └── types.ts          # 类型定义
└── system/                # 系统管理模块
    ├── admin-user/        # 管理员管理
    │   ├── index.ts      # CRUD接口
    │   └── types.ts      # 类型定义
    └── admin-role/        # 角色管理
        ├── index.ts      # CRUD接口、权限管理
        └── types.ts      # 类型定义
```

## API前缀配置

所有API请求前缀统一在 `config.ts` 中配置：

```typescript
// src/services/config.ts
export const API_PREFIX = '/admin-api';

export function buildApiUrl(path: string): string {
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;
  return `${API_PREFIX}${normalizedPath}`;
}
```

### 使用方法

```typescript
import { buildApiUrl } from '@/services/config';

// 构建完整的API路径
const url = buildApiUrl('/auth/login');  // => '/admin-api/auth/login'
```

## 模块说明

### auth - 认证服务

**index.ts** - 基础认证功能
- `login()` - 用户登录
- `logout()` - 退出登录
- `getCurrentUser()` - 获取当前用户信息
- `getCaptcha()` - 获取验证码

**profile.ts** - 个人资料管理
- `updateProfile()` - 更新个人资料
- `changePassword()` - 修改密码
- `uploadAvatar()` - 上传头像

**types.ts** - 类型定义
- `LoginParams`, `LoginResult`
- `CurrentUser`
- `UpdateProfileParams`, `ChangePasswordParams`
- `CaptchaResult`

### common - 公共服务

**index.ts** - 应用级配置
- `getAppConfig()` - 获取应用配置

**types.ts** - 类型定义
- `AppConfig`

### system - 系统管理

#### admin-user - 管理员管理

**index.ts** - 管理员CRUD
- `getAdminUserList()` - 获取管理员列表
- `getAdminUser()` - 获取管理员详情
- `createAdminUser()` - 创建管理员
- `updateAdminUser()` - 更新管理员
- `deleteAdminUser()` - 删除管理员

**types.ts** - 类型定义
- `AdminUser`, `AdminRole`
- `AdminUserListParams`, `AdminUserCreateParams`, `AdminUserUpdateParams`

#### admin-role - 角色管理

**index.ts** - 角色CRUD及权限管理
- `getAdminRoleList()` - 获取角色列表
- `getAllAdminRoles()` - 获取所有角色（不分页）
- `getAdminRole()` - 获取角色详情
- `createAdminRole()` - 创建角色
- `updateAdminRole()` - 更新角色
- `deleteAdminRole()` - 删除角色
- `updateRolePermissions()` - 更新角色权限
- `getAllPermissions()` - 获取所有权限

**types.ts** - 类型定义
- `AdminRole`, `Permission`
- `AdminRoleListParams`, `AdminRoleCreateParams`, `AdminRoleUpdateParams`
- `UpdatePermissionsParams`

## 最佳实践

### 1. 导入规范

推荐分别导入函数和类型：

```typescript
// ✅ 推荐
import { getAdminUserList, createAdminUser } from '@/services/system/admin-user';
import type { AdminUser, AdminUserCreateParams } from '@/services/system/admin-user/types';

// ❌ 不推荐（类型和函数混在一起）
import { getAdminUserList, type AdminUser } from '@/services/system/admin-user';
```

### 2. API路径使用

禁止硬编码API路径，统一使用 `buildApiUrl` 函数：

```typescript
// ✅ 推荐
import { buildApiUrl } from '@/services/config';
const url = buildApiUrl('/auth/avatar');

// ❌ 禁止（硬编码前缀）
const url = '/admin-api/auth/avatar';
```

### 3. 类型安全

所有接口调用必须使用TypeScript类型：

```typescript
import { createAdminUser } from '@/services/system/admin-user';
import type { AdminUserCreateParams } from '@/services/system/admin-user/types';

const handleCreate = async (values: AdminUserCreateParams) => {
  const res = await createAdminUser(values);
  // ...
};
```

### 4. 模块组织原则

- **按菜单层级组织**：services目录结构与路由菜单保持一致
- **职责单一**：每个模块只负责对应功能的API调用
- **类型分离**：types.ts 单独管理类型定义
- **配置统一**：所有API配置在 config.ts 中集中管理

## 新增模块步骤

1. 在对应的菜单层级下创建目录（如 `services/system/xxx/`）
2. 创建 `index.ts` - 编写API调用函数
3. 创建 `types.ts` - 定义相关类型
4. 使用 `buildApiUrl()` 构建API路径
5. 在页面/组件中导入使用

示例：

```typescript
// services/system/xxx/types.ts
export interface XxxItem {
  id: number;
  name: string;
}

export interface XxxListParams {
  page?: number;
  pageSize?: number;
}

// services/system/xxx/index.ts
import { request } from '@umijs/max';
import { buildApiUrl } from '../../config';
import type { XxxItem, XxxListParams } from './types';

export async function getXxxList(params: XxxListParams) {
  return request<API.Response<XxxItem[]>>(buildApiUrl('/xxx'), {
    method: 'GET',
    params,
  });
}
```
