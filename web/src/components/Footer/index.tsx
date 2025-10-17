import { DefaultFooter } from '@ant-design/pro-components';
import { useModel } from '@umijs/max';
import React from 'react';

const Footer: React.FC = () => {
  const { initialState } = useModel('@@initialState');
  const appName = initialState?.appConfig?.name || 'Bico Admin';
  
  return (
    <DefaultFooter
      style={{
        background: 'none',
      }}
      copyright={`Powered by ${appName}`}
      links={[]}
    />
  );
};

export default Footer;
