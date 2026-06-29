import { PageContainer as ProPageContainer } from '@ant-design/pro-components';
import type { PageContainerProps } from '@ant-design/pro-components';
import React from 'react';

/**
 * 自定义 PageContainer 组件，默认隐藏 title
 */
const PageContainer: React.FC<PageContainerProps> = (props) => {
  return <ProPageContainer title={false} {...props} />;
};

export default PageContainer;
