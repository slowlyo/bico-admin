import React, { useState, useEffect } from 'react';
import { ProDescriptions } from '@ant-design/pro-components';
import { Card, Tag, Spin, Button } from 'antd';
import { getRole, getRolePermissions } from '@/services/role';
import type { RoleItem } from '@/services/role';

export type RoleDetailProps = {
  roleId: number;
  onClose: () => void;
};

const RoleDetail: React.FC<RoleDetailProps> = ({ roleId, onClose }) => {
  const [role, setRole] = useState<RoleItem>();
  const [loading, setLoading] = useState(false);
  const [permissions, setPermissions] = useState<any[]>([]);

  const loadRoleDetail = async () => {
    setLoading(true);
    try {
      const [roleResponse, permissionsResponse] = await Promise.all([
        getRole(roleId),
        getRolePermissions(roleId)
      ]);

      if (roleResponse.code === 200) {
        setRole(roleResponse.data);
      }

      if (permissionsResponse.code === 200) {
        setPermissions(permissionsResponse.data || []);
      }
    } catch (error) {
      console.error('加载角色详情失败:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (roleId) {
      loadRoleDetail();
    }
  }, [roleId]);

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px 0' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!role) {
    return <div>角色信息不存在</div>;
  }

  return (
    <div>
      <ProDescriptions
        title="角色详情"
        column={1}
        bordered
        dataSource={role}
        columns={[
          {
            title: 'ID',
            dataIndex: 'id',
            copyable: true,
          },
          {
            title: '角色名称',
            dataIndex: 'name',
            copyable: true,
          },
          {
            title: '角色代码',
            dataIndex: 'code',
            copyable: true,
          },
          {
            title: '描述',
            dataIndex: 'description',
            render: (text) => text || '-',
          },
          {
            title: '状态',
            dataIndex: 'status',
            render: (status) => (
              <Tag color={status === 1 ? 'green' : 'red'}>
                {status === 1 ? '启用' : '禁用'}
              </Tag>
            ),
          },
          {
            title: '创建时间',
            dataIndex: 'created_at',
            valueType: 'dateTime',
          },
          {
            title: '更新时间',
            dataIndex: 'updated_at',
            valueType: 'dateTime',
          },
        ]}
      />

      <Card title="权限列表" style={{ marginTop: 16 }}>
        {permissions.length > 0 ? (
          <div>
            {permissions.map((permission: any) => (
              <Tag key={permission.code || permission.id} color="blue" style={{ margin: '4px' }}>
                {permission.name}
              </Tag>
            ))}
          </div>
        ) : (
          <div style={{ color: '#999', textAlign: 'center', padding: '20px 0' }}>
            暂无权限
          </div>
        )}
      </Card>
    </div>
  );
};

export default RoleDetail;
