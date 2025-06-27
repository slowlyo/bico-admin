import { useLogin } from "@refinedev/core";
import { Form, Input, Button, Card, Typography, Alert } from "antd";
import { UserOutlined, LockOutlined } from "@ant-design/icons";

const { Title, Text } = Typography;

export const Login = () => {
  const { mutate: login, isPending, error } = useLogin();

  const onFinish = (values: { username: string; password: string }) => {
    login(values);
  };

  return (
    <div
      style={{
        height: "100vh",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
      }}
    >
      <Card
        style={{
          width: 400,
          boxShadow: "0 4px 12px rgba(0, 0, 0, 0.15)",
          borderRadius: "8px",
        }}
      >
        <div style={{ textAlign: "center", marginBottom: "24px" }}>
          <Title level={2} style={{ color: "#1890ff", marginBottom: "8px" }}>
            Bico Admin
          </Title>
          <Text type="secondary">欢迎登录管理系统</Text>
        </div>

        {error && (
          <Alert
            message="登录失败"
            description={error.message || "用户名或密码错误"}
            type="error"
            showIcon
            style={{ marginBottom: "16px" }}
          />
        )}

        <Form
          name="login"
          initialValues={{ username: "admin", password: "123456" }}
          onFinish={onFinish}
          layout="vertical"
          requiredMark={false}
        >
          <Form.Item
            name="username"
            label="用户名"
            rules={[
              { required: true, message: "请输入用户名" },
              { min: 3, message: "用户名至少3个字符" },
            ]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="请输入用户名"
              size="large"
            />
          </Form.Item>

          <Form.Item
            name="password"
            label="密码"
            rules={[
              { required: true, message: "请输入密码" },
              { min: 6, message: "密码至少6个字符" },
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="请输入密码"
              size="large"
            />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={isPending}
              size="large"
              style={{ width: "100%" }}
            >
              登录
            </Button>
          </Form.Item>
        </Form>

        <div style={{ textAlign: "center", marginTop: "16px" }}>
          <Text type="secondary" style={{ fontSize: "12px" }}>
            © 2024 Bico Admin. All rights reserved.
          </Text>
        </div>
      </Card>
    </div>
  );
};
