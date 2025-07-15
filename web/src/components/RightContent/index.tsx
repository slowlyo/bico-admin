import React from 'react';
import { Dropdown, Avatar, Space } from 'antd';
import { UserOutlined, LogoutOutlined } from '@ant-design/icons';
import { history, useModel } from '@umijs/max';
import type { MenuProps } from 'antd';
import styles from './index.less';

const RightContent: React.FC = () => {
  const { initialState } = useModel('@@initialState');

  // 登出处理
  const handleLogout = async () => {
    const { logout } = await import('@/services/auth');
    try {
      await logout();
    } catch (error) {
      console.error('登出失败:', error);
    } finally {
      // 清除本地存储
      localStorage.removeItem('token');
      localStorage.removeItem('userInfo');
      localStorage.removeItem('permissions');

      // 跳转到登录页
      history.replace('/login');

      // 刷新页面以清除状态
      window.location.reload();
    }
  };

  // 用户下拉菜单项
  const menuItems: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人信息',
      onClick: () => {
        history.push('/profile');
      },
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ];

  if (!initialState?.currentUser) {
    return null;
  }

  return (
    <Space className={styles.right}>
      <Dropdown
        menu={{ items: menuItems }}
        placement="bottomRight"
        arrow
      >
        <span className={styles.action}>
          <Avatar
            size="small"
            src={initialState.currentUser.avatar}
            icon={<UserOutlined />}
            className={styles.avatar}
          />
          <span className={styles.name}>{initialState.currentUser.nickname}</span>
        </span>
      </Dropdown>
    </Space>
  );
};

export default RightContent;
