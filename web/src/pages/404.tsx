import { Button, Result } from 'antd';
import { history } from '@umijs/max';

const NotFoundPage: React.FC = () => {
  const handleBackHome = () => {
    history.push('/home');
  };

  const handleGoBack = () => {
    history.back();
  };

  return (
    <Result
      status="404"
      title="404"
      subTitle="抱歉，您访问的页面不存在。"
      extra={
        <div style={{ display: 'flex', gap: '8px', justifyContent: 'center' }}>
          <Button type="primary" onClick={handleBackHome}>
            返回首页
          </Button>
          <Button onClick={handleGoBack}>
            返回上一页
          </Button>
        </div>
      }
    />
  );
};

export default NotFoundPage;
