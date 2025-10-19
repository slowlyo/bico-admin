import type { TablePaginationConfig } from 'antd';

/**
 * ProTable 默认分页配置
 */
export const DEFAULT_PAGINATION: TablePaginationConfig = {
  showSizeChanger: true,
  showQuickJumper: true,
  pageSizeOptions: ['10', '20', '50', '100'],
  defaultPageSize: 10,
};
