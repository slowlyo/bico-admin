import { request } from '@umijs/max';
import { buildApiUrl } from './config';

/**
 * 创建标准 CRUD 服务
 * 只需一行代码即可生成完整的 CRUD API
 */
export function createCrudService<
  T extends { id: number },
  CreateParams = Partial<T>,
  UpdateParams = Partial<T>,
  ListParams = Record<string, any>
>(endpoint: string) {
  const url = (path = '') => buildApiUrl(`${endpoint}${path}`);

  return {
    /** 获取列表 */
    list: (params?: ListParams) =>
      request<API.Response<T[]> & { total?: number }>(url(), {
        method: 'GET',
        params,
      }),

    /** 获取详情 */
    get: (id: number) =>
      request<API.Response<T>>(url(`/${id}`), {
        method: 'GET',
      }),

    /** 创建 */
    create: (data: CreateParams) =>
      request<API.Response<T>>(url(), {
        method: 'POST',
        data,
      }),

    /** 更新 */
    update: (id: number, data: UpdateParams) =>
      request<API.Response<T>>(url(`/${id}`), {
        method: 'PUT',
        data,
      }),

    /** 删除 */
    delete: (id: number) =>
      request<API.Response<null>>(url(`/${id}`), {
        method: 'DELETE',
      }),
  };
}

/**
 * 示例用法:
 * 
 * // services/article.ts
 * export const articleService = createCrudService<Article>('/articles');
 * 
 * // 使用
 * const { data } = await articleService.list({ page: 1 });
 * await articleService.create({ title: 'Hello' });
 * await articleService.update(1, { title: 'Updated' });
 * await articleService.delete(1);
 */
