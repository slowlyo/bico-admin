import { DrawerForm, ProFormText, ProFormSwitch, ProFormSelect, ProFormTextArea } from '@ant-design/pro-components';
import React, { useEffect, useState } from 'react';
import { message } from 'antd';
import { AdminUser, AdminUserCreateRequest, AdminUserUpdateRequest } from '@/services/adminUser';
import { getActiveRoles, Role } from '@/services/role';

export interface AdminUserFormProps {
  drawerVisible: boolean;
  onCancel: () => void;
  onSubmit: (values: AdminUserCreateRequest | AdminUserUpdateRequest) => Promise<void>;
  values?: AdminUser; // 编辑时传入，新建时为 undefined
  isEdit?: boolean; // 是否为编辑模式
}

const AdminUserForm: React.FC<AdminUserFormProps> = (props) => {
  const { drawerVisible, onCancel, onSubmit, values, isEdit = false } = props;
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(false);

  // 加载角色列表
  const loadRoles = async () => {
    setLoading(true);
    try {
      const response = await getActiveRoles();
      if (response.code === 200) {
        setRoles(response.data);
      } else {
        message.error(response.message || '获取角色列表失败');
      }
    } catch (error) {
      message.error('获取角色列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (drawerVisible) {
      loadRoles();
    }
  }, [drawerVisible]);

  return (
    <DrawerForm
      title={isEdit ? '编辑管理员用户' : '新建管理员用户'}
      width={600}
      open={drawerVisible}
      onOpenChange={(visible) => {
        if (!visible) {
          onCancel();
        }
      }}
      initialValues={
        isEdit && values
          ? {
              username: values.username,
              name: values.name,
              email: values.email,
              phone: values.phone,
              avatar: values.avatar,
              remark: values.remark,
              enabled: values.status === 1,
              role_ids: values.roles?.map(role => role.id) || [],
            }
          : {
              enabled: true, // 新建时默认启用
              role_ids: [],
            }
      }
      onFinish={async (value) => {
        await onSubmit(value as AdminUserCreateRequest | AdminUserUpdateRequest);
      }}
      drawerProps={{
        destroyOnClose: true,
      }}
    >
      <ProFormText
        name="username"
        label="用户名"
        rules={[
          {
            required: true,
            message: '用户名为必填项',
          },
          {
            min: 3,
            max: 50,
            message: '用户名长度为3-50个字符',
          },
        ]}
        placeholder="请输入用户名"
      />
      
      <ProFormText.Password
        name="password"
        label="密码"
        rules={[
          {
            required: !isEdit, // 新建时必填，编辑时可选
            message: '密码为必填项',
          },
          {
            min: 6,
            max: 100,
            message: '密码长度为6-100个字符',
          },
        ]}
        placeholder={isEdit ? '留空则不修改密码' : '请输入密码'}
      />
      
      <ProFormText
        name="name"
        label="姓名"
        rules={[
          {
            required: true,
            message: '姓名为必填项',
          },
          {
            max: 100,
            message: '姓名长度不能超过100个字符',
          },
        ]}
        placeholder="请输入姓名"
      />
      
      <ProFormText
        name="email"
        label="邮箱"
        rules={[
          {
            type: 'email',
            message: '请输入正确的邮箱格式',
          },
        ]}
        placeholder="请输入邮箱（可选）"
      />
      
      <ProFormText
        name="phone"
        label="手机号"
        rules={[
          {
            pattern: /^1[3-9]\d{9}$/,
            message: '请输入正确的手机号格式',
          },
        ]}
        placeholder="请输入手机号（可选）"
      />
      
      <ProFormText
        name="avatar"
        label="头像URL"
        placeholder="请输入头像URL（可选）"
      />
      
      <ProFormTextArea
        name="remark"
        label="备注"
        rules={[
          {
            max: 500,
            message: '备注长度不能超过500个字符',
          },
        ]}
        placeholder="请输入备注（可选）"
        fieldProps={{
          rows: 3,
        }}
      />

      <ProFormSelect
        name="role_ids"
        label="角色"
        mode="multiple"
        options={roles.map(role => ({
          label: role.name,
          value: role.id,
        }))}
        placeholder="请选择角色"
        fieldProps={{
          loading: loading,
          showSearch: true,
          filterOption: (input: string, option: any) =>
            option?.label?.toLowerCase().includes(input.toLowerCase()),
        }}
        rules={[
          {
            required: true,
            message: '请至少选择一个角色',
          },
        ]}
      />

      <ProFormSwitch
        name="enabled"
        label="启用状态"
      />
    </DrawerForm>
  );
};

export default AdminUserForm;
