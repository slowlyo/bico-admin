import React, { useState } from 'react';
import { PageContainer } from '@ant-design/pro-components';
import { Card, Tabs, message } from 'antd';
import { useModel } from '@umijs/max';
import BasicInfo from './components/BasicInfo';
import ChangePassword from './components/ChangePassword';

const Profile: React.FC = () => {
  const { initialState, setInitialState } = useModel('@@initialState');
  const [activeTab, setActiveTab] = useState('basic');

  // 更新用户信息后刷新全局状态
  const handleUserInfoUpdate = async () => {
    if (initialState?.fetchUserInfo) {
      try {
        const { userInfo, permissions } = await initialState.fetchUserInfo();
        if (userInfo) {
          setInitialState({
            ...initialState,
            currentUser: userInfo,
            permissions: permissions || [],
            name: userInfo.nickname,
            avatar: userInfo.avatar,
          });
          message.success('个人信息更新成功');
        }
      } catch (error) {
        console.error('刷新用户信息失败:', error);
      }
    }
  };

  const tabItems = [
    {
      key: 'basic',
      label: '基本信息',
      children: (
        <BasicInfo 
          userInfo={initialState?.currentUser}
          onUpdate={handleUserInfoUpdate}
        />
      ),
    },
    {
      key: 'password',
      label: '修改密码',
      children: <ChangePassword />,
    },
  ];

  return (
    <PageContainer
      title="个人信息"
      content="管理您的个人信息和账户设置"
    >
      <Card title={false}>
        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          items={tabItems}
        />
      </Card>
    </PageContainer>
  );
};

export default Profile;
