import React, { useState } from 'react';
import { Tabs } from 'antd';

import { PageContainer } from '@/components';

type DocType = 'admin' | 'api';

const docOptions = [
  { label: 'Admin API', key: 'admin' },
  { label: 'Open API', key: 'api' },
];

const docPathMap: Record<DocType, string> = {
  admin: '/swagger/admin/index.html',
  api: '/swagger/api/index.html',
};

/**
 * 接口文档页面。
 *
 * 说明：后端已经提供 Swagger UI，这里只做轻量嵌入，避免在前端重复打包 Swagger 运行时代码。
 */
const ApiDocsPage: React.FC = () => {
  const [docType, setDocType] = useState<DocType>('admin');

  return (
    <PageContainer>
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          gap: 0,
          height: 'calc(100vh - 152px)',
          minHeight: 640,
          overflow: 'hidden',
          border: '1px solid #e5e7eb',
          borderRadius: 8,
          background: '#fff',
        }}
      >
        <Tabs
          activeKey={docType}
          items={docOptions}
          onChange={(key) => setDocType(key as DocType)}
          style={{
            padding: '0 16px',
            marginBottom: 0,
          }}
        />
        <div
          style={{
            flex: 1,
            position: 'relative',
            minHeight: 0,
          }}
        >
          {(Object.keys(docPathMap) as DocType[]).map((type) => (
            <iframe
              key={type}
              title={type === 'admin' ? 'Admin API 文档' : 'Open API 文档'}
              src={docPathMap[type]}
              loading="eager"
              style={{
                display: docType === type ? 'block' : 'none',
                width: '100%',
                height: '100%',
                border: 0,
                background: '#fff',
              }}
            />
          ))}
        </div>
      </div>
    </PageContainer>
  );
};

export default ApiDocsPage;
