import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { ProTable } from '@ant-design/pro-components';
import { Button, Grid, message, Popconfirm, Space } from 'antd';
import React, { useRef, useState, useMemo, useCallback } from 'react';
import { useAccess } from '@umijs/max';
import { PageContainer } from '@/components';
import { DEFAULT_PAGINATION } from '@/constants';
import { transformTableParams, transformTableResponse } from '@/utils/table';
import CrudModal from '../CrudModal';

type CrudService<T> = {
  list: (params?: any) => Promise<API.Response<T[]> & { total?: number }>;
  create?: (data: any) => Promise<API.Response<T>>;
  update?: (id: number, data: any) => Promise<API.Response<T>>;
  delete?: (id: number) => Promise<API.Response<null>>;
};

export interface CrudTableProps<T extends { id: number }> {
  /** 模块名称，如 "用户" */
  title: string;
  /** 权限前缀，如 "system:admin_user"，会自动生成 :list/:create/:edit/:delete */
  permissionPrefix: string;
  /** CRUD 服务 */
  service: CrudService<T>;
  /** 列配置（不含操作列） */
  columns: ProColumns<T>[];
  /** 弹窗表单内容 */
  formContent: React.ReactNode;
  /** 将记录转换为表单初始值 */
  recordToValues?: (record: T) => any;
  /** 自定义请求参数转换 */
  transformParams?: (params: any) => any;
  /** 表格 rowKey，默认 "id" */
  rowKey?: string;
  /** 表格横向滚动宽度 */
  scrollX?: number;
  /** 自定义操作列渲染 */
  renderActions?: (record: T, defaultActions: React.ReactNode) => React.ReactNode;
  /** 额外的工具栏按钮 */
  toolBarExtra?: React.ReactNode[];
  /** 是否显示新建按钮，默认 true */
  showCreate?: boolean;
  /** 是否显示删除确认，默认 true */
  showDeleteConfirm?: boolean;
  /** 表格 actionRef，用于外部控制刷新等 */
  actionRef?: React.MutableRefObject<ActionType | null>;
}

function CrudTable<T extends { id: number }>({
  title,
  permissionPrefix,
  service,
  columns,
  formContent,
  recordToValues,
  transformParams,
  rowKey = 'id',
  scrollX = 1200,
  renderActions,
  toolBarExtra,
  showCreate = true,
  showDeleteConfirm = true,
  actionRef: externalActionRef,
}: CrudTableProps<T>) {
  const internalActionRef = useRef<ActionType>(null);
  const actionRef = externalActionRef || internalActionRef;
  const [modalOpen, setModalOpen] = useState(false);
  const [currentRow, setCurrentRow] = useState<T>();
  const access = useAccess() as Record<string, boolean>;
  const screens = Grid.useBreakpoint();
  const isMobile = !screens.md;

  // 权限 keys
  const perms = useMemo(
    () => ({
      create: `${permissionPrefix}:create`,
      edit: `${permissionPrefix}:edit`,
      delete: `${permissionPrefix}:delete`,
    }),
    [permissionPrefix]
  );

  const handleSuccess = useCallback(() => {
    setModalOpen(false);
    setCurrentRow(undefined);
    actionRef.current?.reload();
  }, []);

  const handleDelete = useCallback(
    async (id: number) => {
      if (!service.delete) return;
      try {
        const res = await service.delete(id);
        if (res.code === 0) {
          message.success('删除成功');
          actionRef.current?.reload();
        } else {
          message.error(res.msg || '删除失败');
        }
      } catch (error: any) {
        message.error(error.message || '删除失败');
      }
    },
    [service]
  );

  const handleEdit = useCallback((record: T) => {
    setCurrentRow(record);
    setModalOpen(true);
  }, []);

  const handleCreate = useCallback(() => {
    setCurrentRow(undefined);
    setModalOpen(true);
  }, []);

  // 默认操作列
  const defaultActions = useCallback(
    (record: T) => (
      <Space>
        {access[perms.edit] && service.update && (
          <a onClick={() => handleEdit(record)}>编辑</a>
        )}
        {access[perms.delete] && service.delete && (
          showDeleteConfirm ? (
            <Popconfirm
              title={`确定删除该${title}吗？`}
              onConfirm={() => handleDelete(record.id)}
              okText="确定"
              cancelText="取消"
            >
              <a style={{ color: '#ff4d4f' }}>删除</a>
            </Popconfirm>
          ) : (
            <a style={{ color: '#ff4d4f' }} onClick={() => handleDelete(record.id)}>
              删除
            </a>
          )
        )}
      </Space>
    ),
    [access, perms, service, title, handleDelete, handleEdit, showDeleteConfirm]
  );

  // 合并操作列
  const finalColumns: ProColumns<T>[] = useMemo(
    () => [
      ...columns,
      {
        title: '操作',
        valueType: 'option' as const,
        width: 150,
        fixed: isMobile ? false : ('right' as const),
        render: (_: any, record: T) =>
          renderActions ? renderActions(record, defaultActions(record)) : defaultActions(record),
      },
    ],
    [columns, renderActions, defaultActions, isMobile]
  );

  return (
    <PageContainer>
      <ProTable<T>
        actionRef={actionRef}
        rowKey={rowKey}
        search={{ labelWidth: 120 }}
        pagination={DEFAULT_PAGINATION}
        toolBarRender={() => [
          ...(toolBarExtra || []),
          showCreate && access[perms.create] && service.create && (
            <Button type="primary" key="create" icon={<PlusOutlined />} onClick={handleCreate}>
              新建
            </Button>
          ),
        ].filter(Boolean)}
        request={async (params, sort) => {
          const apiParams = transformParams
            ? transformParams(transformTableParams(params, sort))
            : transformTableParams(params, sort);
          const res = await service.list(apiParams);
          return transformTableResponse<T>(res);
        }}
        columns={finalColumns}
        scroll={{ x: scrollX }}
      />

      <CrudModal<T>
        title={title}
        open={modalOpen}
        onOpenChange={(visible) => {
          setModalOpen(visible);
          if (!visible) setCurrentRow(undefined);
        }}
        record={currentRow}
        onCreate={service.create}
        onUpdate={service.update}
        onSuccess={handleSuccess}
        recordToValues={recordToValues}
      >
        {formContent}
      </CrudModal>
    </PageContainer>
  );
}

export default CrudTable;
