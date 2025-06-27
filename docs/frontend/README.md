# 前端开发指南

## 🎯 技术栈

- **框架**: React 18 + TypeScript
- **UI库**: [Refine](https://refine.dev/) + [Ant Design](https://ant.design/)
- **状态管理**: React Query (内置于Refine)
- **路由**: React Router
- **构建工具**: Vite
- **包管理**: pnpm

## 🏗️ 项目结构

```
frontend/
├── public/              # 静态资源
├── src/                 # 源代码
│   ├── components/      # 通用组件
│   ├── pages/          # 页面组件
│   ├── contexts/       # React上下文
│   ├── config/         # 配置文件
│   ├── App.tsx         # 主应用组件
│   └── main.tsx        # 入口文件
├── package.json        # 项目配置
├── vite.config.ts      # Vite配置
└── tsconfig.json       # TypeScript配置
```

## 🚀 快速开始

### 1. 环境准备
```bash
# 安装Node.js 18+
node --version

# 安装pnpm
npm install -g pnpm

# 安装依赖
cd frontend
pnpm install
```

### 2. 配置环境变量
```bash
# 复制环境变量文件
cp .env.example .env

# 编辑配置
# VITE_ADMIN_API_URL=http://localhost:8080/admin/api
# VITE_API_URL=http://localhost:8080/api
```

### 3. 启动开发服务器
```bash
# 开发模式
pnpm dev

# 或使用Makefile
make dev-frontend
```

## 📋 开发规范

### 组件开发规范

#### 1. 组件命名
- 使用大驼峰命名法 (PascalCase)
- 文件名与组件名保持一致
- 使用描述性的名称

```typescript
// ✅ 正确
const UserList = () => { ... }
// 文件名: UserList.tsx

// ❌ 错误
const userlist = () => { ... }
const List = () => { ... }
```

#### 2. 组件结构
```typescript
// UserList.tsx
import React from 'react';
import { List, useTable } from '@refinedev/antd';
import { Table } from 'antd';

interface UserListProps {
  // 定义props类型
}

export const UserList: React.FC<UserListProps> = () => {
  const { tableProps } = useTable();

  return (
    <List>
      <Table {...tableProps}>
        {/* 表格列定义 */}
      </Table>
    </List>
  );
};
```

#### 3. 类型定义
```typescript
// types/user.ts
export interface User {
  id: number;
  username: string;
  email: string;
  status: UserStatus;
  created_at: string;
  updated_at: string;
}

export enum UserStatus {
  INACTIVE = 0,
  ACTIVE = 1,
  BLOCKED = 2,
}

export interface UserCreateRequest {
  username: string;
  email: string;
  password: string;
  role_ids?: number[];
}
```

## 🔧 Refine 使用指南

### 1. 数据提供者配置
```typescript
// App.tsx
import { Refine } from '@refinedev/core';
import { dataProvider } from '@refinedev/simple-rest';
import config from './config';

function App() {
  return (
    <Refine
      dataProvider={dataProvider(config.adminApiUrl)}
      resources={[
        {
          name: 'users',
          list: '/users',
          create: '/users/create',
          edit: '/users/edit/:id',
          show: '/users/show/:id',
          meta: {
            canDelete: true,
          },
        },
      ]}
    >
      {/* 应用内容 */}
    </Refine>
  );
}
```

### 2. 列表页面开发
```typescript
// pages/users/list.tsx
import React from 'react';
import { List, useTable, EditButton, ShowButton, DeleteButton } from '@refinedev/antd';
import { Table, Space } from 'antd';

export const UserList = () => {
  const { tableProps } = useTable();

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="id" title="ID" />
        <Table.Column dataIndex="username" title="用户名" />
        <Table.Column dataIndex="email" title="邮箱" />
        <Table.Column
          dataIndex="status"
          title="状态"
          render={(value) => (
            <span>{value === 1 ? '激活' : '未激活'}</span>
          )}
        />
        <Table.Column
          title="操作"
          dataIndex="actions"
          render={(_, record) => (
            <Space>
              <EditButton hideText size="small" recordItemId={record.id} />
              <ShowButton hideText size="small" recordItemId={record.id} />
              <DeleteButton hideText size="small" recordItemId={record.id} />
            </Space>
          )}
        />
      </Table>
    </List>
  );
};
```

### 3. 表单页面开发
```typescript
// pages/users/create.tsx
import React from 'react';
import { Create, useForm } from '@refinedev/antd';
import { Form, Input, Select } from 'antd';

export const UserCreate = () => {
  const { formProps, saveButtonProps } = useForm();

  return (
    <Create saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="用户名"
          name="username"
          rules={[
            {
              required: true,
              message: '请输入用户名',
            },
          ]}
        >
          <Input />
        </Form.Item>
        
        <Form.Item
          label="邮箱"
          name="email"
          rules={[
            {
              required: true,
              type: 'email',
              message: '请输入有效的邮箱地址',
            },
          ]}
        >
          <Input />
        </Form.Item>
        
        <Form.Item
          label="密码"
          name="password"
          rules={[
            {
              required: true,
              min: 6,
              message: '密码至少6位',
            },
          ]}
        >
          <Input.Password />
        </Form.Item>
        
        <Form.Item
          label="状态"
          name="status"
          initialValue={1}
        >
          <Select>
            <Select.Option value={0}>未激活</Select.Option>
            <Select.Option value={1}>激活</Select.Option>
          </Select>
        </Form.Item>
      </Form>
    </Create>
  );
};
```

## 🎨 样式和主题

### 1. Ant Design 主题定制
```typescript
// App.tsx
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';

const theme = {
  token: {
    colorPrimary: '#1890ff',
    borderRadius: 6,
  },
};

function App() {
  return (
    <ConfigProvider theme={theme} locale={zhCN}>
      {/* 应用内容 */}
    </ConfigProvider>
  );
}
```

### 2. 自定义样式
```css
/* styles/global.css */
.custom-table {
  .ant-table-thead > tr > th {
    background-color: #fafafa;
    font-weight: 600;
  }
}

.custom-form {
  .ant-form-item-label > label {
    font-weight: 500;
  }
}
```

## 🔐 认证和权限

### 1. 认证提供者
```typescript
// authProvider.ts
import { AuthProvider } from '@refinedev/core';
import config from './config';

export const authProvider: AuthProvider = {
  login: async ({ username, password }) => {
    const response = await fetch(`${config.adminApiUrl}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    if (response.ok) {
      const data = await response.json();
      localStorage.setItem('token', data.data.token);
      localStorage.setItem('user', JSON.stringify(data.data.user));
      return {
        success: true,
        redirectTo: '/',
      };
    }

    return {
      success: false,
      error: {
        name: 'LoginError',
        message: '用户名或密码错误',
      },
    };
  },

  logout: async () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    return {
      success: true,
      redirectTo: '/login',
    };
  },

  check: async () => {
    const token = localStorage.getItem('token');
    if (token) {
      return {
        authenticated: true,
      };
    }

    return {
      authenticated: false,
      redirectTo: '/login',
    };
  },

  getPermissions: async () => {
    const user = localStorage.getItem('user');
    if (user) {
      const userData = JSON.parse(user);
      return userData.roles;
    }
    return null;
  },

  getIdentity: async () => {
    const user = localStorage.getItem('user');
    if (user) {
      return JSON.parse(user);
    }
    return null;
  },

  onError: async (error) => {
    if (error.response?.status === 401) {
      return {
        logout: true,
      };
    }

    return { error };
  },
};
```

### 2. 权限控制
```typescript
// components/ProtectedComponent.tsx
import { useCan } from '@refinedev/core';

interface ProtectedComponentProps {
  resource: string;
  action: string;
  children: React.ReactNode;
}

export const ProtectedComponent: React.FC<ProtectedComponentProps> = ({
  resource,
  action,
  children,
}) => {
  const { data } = useCan({
    resource,
    action,
  });

  if (data?.can) {
    return <>{children}</>;
  }

  return null;
};

// 使用示例
<ProtectedComponent resource="users" action="create">
  <CreateButton />
</ProtectedComponent>
```

## 🧪 测试

### 1. 组件测试
```typescript
// __tests__/UserList.test.tsx
import { render, screen } from '@testing-library/react';
import { TestWrapper } from './test-utils';
import { UserList } from '../pages/users/list';

describe('UserList', () => {
  it('renders user list correctly', () => {
    render(
      <TestWrapper>
        <UserList />
      </TestWrapper>
    );

    expect(screen.getByText('用户列表')).toBeInTheDocument();
  });
});
```

### 2. 测试工具配置
```typescript
// test-utils.tsx
import React from 'react';
import { Refine } from '@refinedev/core';
import { MockJSONServer } from '@refinedev/simple-rest';

export const TestWrapper: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  return (
    <Refine
      dataProvider={MockJSONServer}
      resources={[
        {
          name: 'users',
        },
      ]}
    >
      {children}
    </Refine>
  );
};
```

## 📦 构建和部署

### 1. 构建配置
```typescript
// vite.config.ts
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: 'dist',
    sourcemap: false,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom'],
          antd: ['antd'],
          refine: ['@refinedev/core', '@refinedev/antd'],
        },
      },
    },
  },
});
```

### 2. 环境变量
```bash
# .env.production
VITE_ADMIN_API_URL=https://api.yourdomain.com/admin/api
VITE_API_URL=https://api.yourdomain.com/api
```

## 📝 最佳实践

1. **组件设计**
   - 保持组件单一职责
   - 使用TypeScript类型检查
   - 编写可复用的组件

2. **状态管理**
   - 优先使用React Query缓存
   - 避免不必要的全局状态
   - 合理使用Context

3. **性能优化**
   - 使用React.memo优化渲染
   - 实现虚拟滚动处理大列表
   - 代码分割和懒加载

4. **用户体验**
   - 提供加载状态反馈
   - 实现错误边界处理
   - 支持键盘导航

5. **代码质量**
   - 使用ESLint和Prettier
   - 编写单元测试
   - 定期重构代码
