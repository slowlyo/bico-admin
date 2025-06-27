import React from "react";
import { Card, Col, Row, Statistic, Typography } from "antd";
import { UserOutlined, FileTextOutlined, SettingOutlined, DashboardOutlined } from "@ant-design/icons";

const { Title, Paragraph } = Typography;

export const Dashboard: React.FC = () => {
  return (
    <div style={{ padding: "24px" }}>
      <Title level={2}>
        <DashboardOutlined style={{ marginRight: "8px" }} />
        控制台
      </Title>
      
      <Paragraph type="secondary">
        欢迎使用 Bico Admin 管理系统
      </Paragraph>

      <Row gutter={[16, 16]} style={{ marginTop: "24px" }}>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="用户总数"
              value={1}
              prefix={<UserOutlined />}
              valueStyle={{ color: "#3f8600" }}
            />
          </Card>
        </Col>
        
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="文章数量"
              value={0}
              prefix={<FileTextOutlined />}
              valueStyle={{ color: "#1890ff" }}
            />
          </Card>
        </Col>
        
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="分类数量"
              value={0}
              prefix={<SettingOutlined />}
              valueStyle={{ color: "#722ed1" }}
            />
          </Card>
        </Col>
        
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="系统状态"
              value="正常"
              valueStyle={{ color: "#52c41a" }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: "24px" }}>
        <Col xs={24} lg={12}>
          <Card title="系统信息" bordered={false}>
            <div style={{ lineHeight: "2" }}>
              <p><strong>系统名称：</strong>Bico Admin</p>
              <p><strong>版本：</strong>1.0.0</p>
              <p><strong>技术栈：</strong>React + Refine + Ant Design</p>
              <p><strong>后端：</strong>Go + Fiber + GORM</p>
              <p><strong>数据库：</strong>MySQL</p>
            </div>
          </Card>
        </Col>
        
        <Col xs={24} lg={12}>
          <Card title="快速操作" bordered={false}>
            <div style={{ lineHeight: "2" }}>
              <p>• 用户管理</p>
              <p>• 内容管理</p>
              <p>• 系统设置</p>
              <p>• 数据统计</p>
              <p>• 日志查看</p>
            </div>
          </Card>
        </Col>
      </Row>

      <Row style={{ marginTop: "24px" }}>
        <Col span={24}>
          <Card title="最近活动" bordered={false}>
            <div style={{ textAlign: "center", padding: "40px", color: "#999" }}>
              暂无活动记录
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
};
