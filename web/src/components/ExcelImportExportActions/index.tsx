import React, { useMemo, useState } from 'react';
import { Button, Modal, Space, Typography, Upload } from 'antd';
import type { ModalProps } from 'antd';
import type { UploadProps } from 'antd';
import { DownloadOutlined, UploadOutlined } from '@ant-design/icons';

const { Dragger } = Upload;

export interface ExcelImportExportActionsProps {
  onExportClick?: () => void;
  exporting?: boolean;
  importDisabled?: boolean;

  showImport?: boolean;
  showExport?: boolean;

  onImportFile?: (file: File) => Promise<void>;
  onDownloadTemplate?: () => Promise<void> | void;

  showDownloadTemplate?: boolean;
  downloadTemplateText?: string;
  accept?: string;

  importModalTitle?: React.ReactNode;
  importModalContent?: React.ReactNode;
  importModalProps?: Omit<ModalProps, 'open' | 'onCancel' | 'title' | 'footer'>;
}

function ExcelImportExportActions({
  onExportClick,
  exporting,
  importDisabled,
  showImport = true,
  showExport = true,
  onImportFile,
  onDownloadTemplate,
  showDownloadTemplate = true,
  downloadTemplateText = '下载导入模板',
  accept = '.xlsx,.xlsm,.xltx,.xltm,.csv',
  importModalTitle = '导入',
  importModalContent,
  importModalProps,
}: ExcelImportExportActionsProps) {
  const [importOpen, setImportOpen] = useState(false);
  const [importing, setImporting] = useState(false);

  const canOpenImportModal = useMemo(() => {
    // 允许：
    // - 用户完全覆盖 importModalContent
    // - 或者使用内置默认内容（需要 onImportFile）
    return !!importModalContent || !!onImportFile;
  }, [importModalContent, onImportFile]);

  const handleImportClick = () => {
    if (!canOpenImportModal) {
      return;
    }
    setImportOpen(true);
  };

  const handleCustomRequest: UploadProps['customRequest'] = async (options) => {
    const file = options.file as File;
    if (!onImportFile) {
      options.onError?.(new Error('未配置导入处理函数'));
      return;
    }

    setImporting(true);
    try {
      await onImportFile(file);
      options.onSuccess?.({}, file);
    } catch (e: any) {
      options.onError?.(e);
    } finally {
      setImporting(false);
    }
  };

  const defaultImportModalContent = (
    <>
      <Dragger
        name="file"
        multiple={false}
        maxCount={1}
        accept={accept}
        customRequest={handleCustomRequest}
        showUploadList={false}
        disabled={importDisabled || importing}
      >
        <p style={{ marginBottom: 8 }}>将 Excel 文件拖拽到此处，或点击选择文件上传</p>
        <p style={{ marginBottom: 0, color: '#8c8c8c' }}>支持 .xlsx/.xlsm/.xltx/.xltm/.csv</p>
      </Dragger>

      {showDownloadTemplate && (
        <div style={{ marginTop: 12 }}>
          <Button
            type="link"
            onClick={() => onDownloadTemplate?.()}
            style={{ padding: 0 }}
            disabled={!onDownloadTemplate}
          >
            {downloadTemplateText}
          </Button>
        </div>
      )}

      {importing && (
        <Typography.Paragraph type="secondary" style={{ marginTop: 8, marginBottom: 0 }}>
          正在上传并解析，请稍候...
        </Typography.Paragraph>
      )}
    </>
  );

  return (
    <>
      <Space>
        {showImport && (
          <Button
            icon={<UploadOutlined />}
            onClick={handleImportClick}
            disabled={importDisabled || !canOpenImportModal}
          >
            导入
          </Button>
        )}
        {showExport && (
          <Button
            icon={<DownloadOutlined />}
            loading={exporting}
            onClick={onExportClick}
            type="primary"
            disabled={!onExportClick}
          >
            导出
          </Button>
        )}
      </Space>

      {canOpenImportModal && (
        <Modal
          title={importModalTitle}
          open={importOpen}
          onCancel={() => setImportOpen(false)}
          footer={null}
          destroyOnClose
          {...importModalProps}
        >
          {importModalContent || defaultImportModalContent}
        </Modal>
      )}
    </>
  );
}

export default ExcelImportExportActions;
