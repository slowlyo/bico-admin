import { Skeleton } from 'antd';

const Loading: React.FC = () => (
  <div
    style={{
      display: 'flex',
      alignItems: 'start',
      justifyContent: 'center',
      width: '100%',
      height: '100vh',
      overflow: 'hidden',
    }}
  >
    <Skeleton
      style={{
        maxWidth: '1200px',
        padding: '24px',
      }}
      active
      paragraph={{ rows: 8 }}
    />
  </div>
);

export default Loading;
