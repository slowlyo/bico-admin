import { requestClient } from '#/api/request';

export interface PageResult<T> {
  list: T[];
  total: number;
}

// 后端分页与排序的公共查询参数，对应 pagination 包。
export interface PageQuery {
  page: number;
  pageSize: number;
  sortField?: string;
  sortOrder?: string;
}

// 后端 swagger 里字段默认可选，响应数据实际总会返回这些字段，因此展开为必填。
export type Defined<T> = Required<T>;

export type CrudApi<T, Q extends PageQuery, C, U = C> = ReturnType<
  typeof createCrudApi<T, Q, C, U>
>;

// 按后端 CRUDHandler 的标准路由生成接口：列表、创建、更新、单删、批量删。
// 创建与更新载荷类型可以不同（如更新时密码可选）。base 是模块分组路径，如 '/admin-users'；批量删除固定走 DELETE {base}/batch。
export function createCrudApi<T, Q extends PageQuery, C, U = C>(
  base: string,
) {
  return {
    create: (data: C) => requestClient.post(base, data),
    list: (params: Q) => requestClient.get<PageResult<T>>(base, { params }),
    remove: (id: number) => requestClient.delete(`${base}/${id}`),
    removeBatch: (ids: number[]) =>
      requestClient.delete(`${base}/batch`, { data: { ids } }),
    update: (id: number, data: U) => requestClient.put(`${base}/${id}`, data),
  };
}
