import { PageContainer } from '@ant-design/pro-components';
import { Card, Typography } from 'antd';
import { useModel } from '@umijs/max';
import styles from './index.less';

const { Title, Paragraph } = Typography;

const HomePage: React.FC = () => {
  const { initialState } = useModel('@@initialState');
  const currentUser = initialState?.currentUser;

  return (
    <PageContainer ghost>
      <div className={styles.container}>
        <Card>
          <Typography>
            <Title level={2}>欢迎使用 Bico Admin</Title>
            <Paragraph>
              这是一个基于 Go Fiber + UmiJS + Ant Design Pro 构建的现代化管理后台系统。
            </Paragraph>
            <Paragraph>
              当前用户：{currentUser?.nickname || currentUser?.username || '游客'}
            </Paragraph>
          </Typography>
        </Card>
      </div>
    </PageContainer>
  );
};

export default HomePage;
