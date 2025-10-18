import { PageContainer } from '@ant-design/pro-components';
import { Card, Col, Row, Statistic } from 'antd';
import { UserOutlined, TeamOutlined, SafetyOutlined } from '@ant-design/icons';
import React from 'react';
import { useModel } from '@umijs/max';

const Dashboard: React.FC = () => {
  const { initialState } = useModel('@@initialState');
  const { currentUser } = initialState || {};

  return (
    <PageContainer>
      <Card style={{ marginBottom: 24 }}>
        <h2>欢迎回来，{currentUser?.name || currentUser?.username}！</h2>
        <p>这是你的工作台，祝你工作愉快！</p>
      </Card>

      <Row gutter={16}>
        <Col span={8}>
          <Card>
            <Statistic
              title="用户总数"
              value={0}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="角色总数"
              value={0}
              prefix={<TeamOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="权限总数"
              value={0}
              prefix={<SafetyOutlined />}
              valueStyle={{ color: '#cf1322' }}
            />
          </Card>
        </Col>
      </Row>
    </PageContainer>
  );
};

export default Dashboard;
