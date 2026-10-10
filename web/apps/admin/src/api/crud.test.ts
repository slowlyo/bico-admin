import { beforeEach, describe, expect, it, vi } from 'vitest';

const client = vi.hoisted(() => ({
  delete: vi.fn(),
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
}));

vi.mock('#/api/request', () => ({ requestClient: client }));

const { createCrudApi } = await import('./crud');

describe('createCrudApi', () => {
  const api = createCrudApi<{ id: number }, any, { name: string }>('/things');

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('list 以 params 传查询参数', async () => {
    client.get.mockResolvedValue({ list: [], total: 0 });
    await api.list({ page: 1, pageSize: 10 });
    expect(client.get).toHaveBeenCalledWith('/things', {
      params: { page: 1, pageSize: 10 },
    });
  });

  it('create 与 update 使用 POST 和 PUT', async () => {
    await api.create({ name: 'a' });
    expect(client.post).toHaveBeenCalledWith('/things', { name: 'a' });
    await api.update(7, { name: 'b' });
    expect(client.put).toHaveBeenCalledWith('/things/7', { name: 'b' });
  });

  it('remove 走 DELETE /:id，removeBatch 把 ids 放在请求体', async () => {
    await api.remove(3);
    expect(client.delete).toHaveBeenCalledWith('/things/3');
    await api.removeBatch([1, 2]);
    expect(client.delete).toHaveBeenCalledWith('/things/batch', {
      data: { ids: [1, 2] },
    });
  });
});
