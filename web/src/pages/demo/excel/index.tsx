import React, { useMemo, useState } from 'react';
import { Button, Card, Descriptions, Typography, message } from 'antd';

import { PageContainer } from '@/components';
import { ExcelImportExportActions } from '@/components';
import { downloadBlob, getFilenameFromContentDisposition } from '@/utils/download';
import {
  downloadDemoExcelTemplate,
  exportDemoExcel,
  importDemoExcel,
} from '@/services/demo/excel';

const DemoExcelPage: React.FC = () => {
  const [importing, setImporting] = useState(false);
  const [exporting, setExporting] = useState(false);
  const [importResult, setImportResult] = useState<{ total: number; preview: string[][] }>();

  const headers = useMemo(() => ['姓名', '手机号', '年龄', '城市'], []);

  const handleDownloadTemplate = async () => {
    try {
      const resp: any = await downloadDemoExcelTemplate();
      const blob = resp?.data as Blob;
      const cd = resp?.response?.headers?.get?.('content-disposition');
      const filename = getFilenameFromContentDisposition(cd) || '导入模板_示例.xlsx';
      downloadBlob(blob, filename);
    } catch (e: any) {
      message.error(e?.message || e?.data?.msg || '下载模板失败');
    }
  };

  const handleExport = async () => {
    setExporting(true);
    try {
      const resp: any = await exportDemoExcel();
      const blob = resp?.data as Blob;
      const cd = resp?.response?.headers?.get?.('content-disposition');
      const filename = getFilenameFromContentDisposition(cd) || '导出_示例.xlsx';
      downloadBlob(blob, filename);
    } catch (e: any) {
      message.error(e?.message || e?.data?.msg || '导出失败');
    } finally {
      setExporting(false);
    }
  };

  const handleImportFile = async (file: File) => {
    setImporting(true);
    try {
      const res = await importDemoExcel(file);
      if (res.code === 0 && res.data) {
        message.success(`导入解析成功，共 ${res.data.total} 行`);
        setImportResult(res.data);
        return;
      }
      message.error(res.msg || '导入失败');
      throw new Error(res.msg || '导入失败');
    } catch (e: any) {
      message.error(e?.message || e?.data?.msg || '导入失败');
      throw e;
    } finally {
      setImporting(false);
    }
  };

  return (
    <PageContainer>
      <Card style={{ background: '#fff' }}>
        <Typography.Title level={4} style={{ marginTop: 0 }}>
          Excel 导入/导出示例
        </Typography.Title>
        <Typography.Paragraph type="secondary">
          演示：下载模板（后端通过 Excel 库输出表头） / 拖拽上传导入 / 导出（按钮 loading 处理）。
        </Typography.Paragraph>

        <ExcelImportExportActions
          onExportClick={handleExport}
          exporting={exporting}
          importDisabled={importing}
          importModalTitle="导入"
          onImportFile={handleImportFile}
          onDownloadTemplate={handleDownloadTemplate}
        />

        <div style={{ marginTop: 16 }}>
          <Descriptions size="small" bordered column={1} title="模板表头">
            <Descriptions.Item label="表头">{headers.join(' / ')}</Descriptions.Item>
          </Descriptions>
        </div>

        {importResult && (
          <div style={{ marginTop: 16 }}>
            <Typography.Title level={5} style={{ marginBottom: 8 }}>
              导入预览（前 5 行）
            </Typography.Title>
            <div style={{ marginBottom: 8 }}>
              <Button
                size="small"
                onClick={() => setImportResult(undefined)}
              >
                清空导入结果
              </Button>
            </div>
            <pre
              style={{
                background: '#fafafa',
                border: '1px solid #f0f0f0',
                borderRadius: 6,
                padding: 12,
                whiteSpace: 'pre-wrap',
                wordBreak: 'break-word',
                marginBottom: 0,
              }}
            >
              {JSON.stringify(importResult.preview, null, 2)}
            </pre>
          </div>
        )}
      </Card>

    </PageContainer>
  );
};

export default DemoExcelPage;
