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
  data?: any;
  total?: number;
}) {
  const responseData = response.data;
  
  // 处理后端分页响应数据（包裹在 data.list 中）
  if (responseData && typeof responseData === 'object' && 'list' in responseData) {
    return {
      data: (responseData.list || []) as T[],
      total: responseData.total || 0,
      success: response.code === 0,
    };
  }

  // 处理无分页或旧接口响应数据（data 直接为数组）
  return {
    data: (responseData || []) as T[],
    total: response.total || 0,
    success: response.code === 0,
  };
}
