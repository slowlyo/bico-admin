import React, { useState } from 'react';
import { Card, Typography } from 'antd';

import { PageContainer } from '@/components';
import RichEditor from '@/components/RichEditor';

/**
 * 编辑器示例页面
 */
const EditorDemoPage: React.FC = () => {
  const [html, setHtml] = useState<string>('<p>这是一个示例内容</p>');

  return (
    <PageContainer>
      <Card style={{ background: '#fff' }}>
        <Typography.Title level={4} style={{ marginTop: 0 }}>
          富文本编辑器示例
        </Typography.Title>
        <Typography.Paragraph type="secondary">
          演示 RichEditor 的基本用法：受控 value + onChange，并支持图片/视频上传。
        </Typography.Paragraph>

        <RichEditor value={html} onChange={setHtml} />

        <div style={{ marginTop: 16 }}>
          <Typography.Title level={5}>当前 HTML</Typography.Title>
          <pre
            style={{
              background: '#fafafa',
              border: '1px solid #f0f0f0',
              borderRadius: 6,
              padding: 12,
              whiteSpace: 'pre-wrap',
              wordBreak: 'break-word',
            }}
          >
            {html}
          </pre>
        </div>
      </Card>
    </PageContainer>
  );
};

export default EditorDemoPage;
