import type { ParamsType } from '@ant-design/pro-components';
import type { SortOrder } from 'antd/es/table/interface';

/**
 * ProTable request 参数转换
 * 将 ProTable 的 params 和 sort 转换为后端 API 需要的格式
 */
export function transformTableParams<T = any>(
  params: ParamsType & {
    pageSize?: number;
    current?: number;
  },
  sort?: Record<string, SortOrder>
): T {
  const sortField = Object.keys(sort || {})[0];
  const sortOrder = sort?.[sortField];

  return {
    ...params,
    page: params.current,
    pageSize: params.pageSize,
    sortField: sortField,
    sortOrder: sortOrder || undefined,
  } as T;
}

/**
 * ProTable response 转换
 * 将后端 API 响应转换为 ProTable 需要的格式
 */
export function transformTableResponse<T>(response: {
  code: number;
  data?: T[];
  total?: number;
}) {
  return {
    data: (response.data || []) as T[],
    total: response.total || 0,
    success: response.code === 0,
  };
}
