import { LogoutOutlined, UserOutlined } from '@ant-design/icons';
import { Avatar, Dropdown, Menu, Spin } from 'antd';
import { history, useModel } from '@umijs/max';
import React, { useCallback } from 'react';
import { logout } from '@/services/auth';
import { createStyles } from 'antd-style';

const useStyles = createStyles(({ token }) => {
  return {
    action: {
      display: 'flex',
      alignItems: 'center',
      height: '48px',
      marginLeft: 'auto',
      cursor: 'pointer',
      padding: '0 12px',
      borderRadius: token.borderRadius,
      '&:hover': {
        backgroundColor: token.colorBgTextHover,
      },
    },
    avatar: {
      marginRight: 8,
      color: token.colorPrimary,
      verticalAlign: 'top',
      background: 'rgba(255, 255, 255, 0.85)',
    },
    name: {
      color: token.colorTextHeading,
      fontWeight: 500,
      fontSize: '14px',
    },
  };
});

export interface GlobalHeaderRightProps {
  menu?: boolean;
  children?: React.ReactNode;
}

export const AvatarName: React.FC = () => {
  const { initialState } = useModel('@@initialState');
  const { currentUser } = initialState || {};
  return <span>{currentUser?.nickname || currentUser?.username}</span>;
};

export const AvatarDropdown: React.FC<GlobalHeaderRightProps> = ({ children }) => {
  const { styles } = useStyles();
  const { initialState, setInitialState } = useModel('@@initialState');

  const onMenuClick = useCallback(
    async (event: any) => {
      const { key } = event;
      if (key === 'logout') {
        try {
          await logout();
        } catch (error) {
          console.log('退出登录失败:', error);
        } finally {
          localStorage.removeItem('token');
          await setInitialState((s: any) => ({ ...s, currentUser: undefined }));
          history.push('/login');
        }
        return;
      }
      if (key === 'profile') {
        history.push('/profile');
        return;
      }
    },
    [setInitialState],
  );

  const loading = (
    <span className={styles.action}>
      <Spin
        size="small"
        style={{
          marginLeft: 8,
          marginRight: 8,
        }}
      />
    </span>
  );

  if (!initialState) {
    return loading;
  }

  const { currentUser } = initialState;

  if (!currentUser || !currentUser.username) {
    return loading;
  }

  const menuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人资料',
    },
    {
      type: 'divider' as const,
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
    },
  ];

  return (
    <Dropdown
      menu={{
        selectedKeys: [],
        onClick: onMenuClick,
        items: menuItems,
      }}
      placement="bottomRight"
    >
      {children || (
        <span className={styles.action}>
          <Avatar
            size="small"
            className={styles.avatar}
            src={currentUser.avatar}
            alt="avatar"
          />
          <span className={styles.name}>
            {currentUser.nickname || currentUser.username}
          </span>
        </span>
      )}
    </Dropdown>
  );
};
