import { PageContainer } from '@ant-design/pro-components';
import { Access, useAccess, useModel } from '@umijs/max';
import { Button, Card, Space, Tag, Divider } from 'antd';
import { PERMISSIONS } from '@/utils/permission';

const AccessPage: React.FC = () => {
  const access = useAccess();
  const { initialState } = useModel('@@initialState');
  const userPermissions = initialState?.permissions || [];

  return (
    <PageContainer
      ghost
      header={{
        title: '权限控制示例',
      }}
    >
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        {/* 用户权限信息 */}
        <Card title="当前用户权限">
          <div>
            <strong>权限列表：</strong>
            <div style={{ marginTop: 8 }}>
              {userPermissions.length > 0 ? (
                userPermissions.map(permission => (
                  <Tag key={permission} color="blue" style={{ margin: 4 }}>
                    {permission}
                  </Tag>
                ))
              ) : (
                <Tag color="red">无权限</Tag>
              )}
            </div>
          </div>
        </Card>

        {/* 基础权限控制 */}
        <Card title="基础权限控制">
          <Space wrap>
            <Access accessible={access.isLogin}>
              <Button type="primary">已登录用户可见</Button>
            </Access>

            <Access accessible={access.canSeeAdmin}>
              <Button>管理员可见</Button>
            </Access>
          </Space>
        </Card>

        {/* 管理员管理权限 */}
        <Card title="管理员管理权限">
          <Space wrap>
            <Access accessible={access.canViewAdminUsers}>
              <Button>查看管理员列表</Button>
            </Access>

            <Access accessible={access.canCreateAdminUser}>
              <Button type="primary">创建管理员</Button>
            </Access>

            <Access accessible={access.canEditAdminUser}>
              <Button>编辑管理员</Button>
            </Access>

            <Access accessible={access.canDeleteAdminUser}>
              <Button danger>删除管理员</Button>
            </Access>

            <Access accessible={access.canResetAdminPassword}>
              <Button>重置密码</Button>
            </Access>
          </Space>
        </Card>

        {/* 角色管理权限 */}
        <Card title="角色管理权限">
          <Space wrap>
            <Access accessible={access.canViewRoles}>
              <Button>查看角色列表</Button>
            </Access>

            <Access accessible={access.canCreateRole}>
              <Button type="primary">创建角色</Button>
            </Access>

            <Access accessible={access.canEditRole}>
              <Button>编辑角色</Button>
            </Access>

            <Access accessible={access.canDeleteRole}>
              <Button danger>删除角色</Button>
            </Access>
          </Space>
        </Card>

        {/* 使用权限工具函数 */}
        <Card title="使用权限工具函数">
          <Space wrap>
            {access.hasPermission(PERMISSIONS.ADMIN_USER_LIST) && (
              <Button>通过工具函数检查权限</Button>
            )}

            {access.hasAnyPermission([PERMISSIONS.ADMIN_USER_CREATE, PERMISSIONS.ROLE_CREATE]) && (
              <Button type="primary">有任意创建权限</Button>
            )}
          </Space>
        </Card>
      </Space>
    </PageContainer>
  );
};

export default AccessPage;
