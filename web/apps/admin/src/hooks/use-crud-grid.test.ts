import { beforeEach, describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
  captured: {} as any,
  gridApi: {
    grid: { getCheckboxRecords: vi.fn() },
    query: vi.fn(),
  },
  message: { success: vi.fn(), warning: vi.fn() },
}));

vi.mock('antdv-next', () => ({ message: mocks.message }));
vi.mock('#/adapter/vxe-table', () => ({
  useVbenVxeGrid: (options: any) => {
    mocks.captured = options;
    return [{}, mocks.gridApi];
  },
}));

const { useCrudGrid } = await import('./use-crud-grid');

interface Row {
  id: number;
  system?: boolean;
}

function setup(canDelete?: (row: Row) => boolean) {
  const api = {
    list: vi.fn().mockResolvedValue({ list: [{ id: 1 }], total: 1 }),
    remove: vi.fn().mockResolvedValue(undefined),
    removeBatch: vi.fn().mockResolvedValue(undefined),
  };
  const result = useCrudGrid<Row, any>({
    api: api as any,
    canDelete,
    entityName: '角色',
    formOptions: { schema: [] },
    gridOptions: { columns: [] },
  });
  return { api, ...result };
}

describe('useCrudGrid 配置', () => {
  beforeEach(() => vi.clearAllMocks());

  it('开启服务端排序与排序代理，并带上默认配置', () => {
    setup();
    const options = mocks.captured.gridOptions;
    expect(options.sortConfig).toEqual({ remote: true });
    expect(options.proxyConfig.sort).toBe(true);
    expect(options.rowConfig).toEqual({ keyField: 'id' });
    expect(options.height).toBe('auto');
  });

  it('页面传入的 gridOptions 优先于默认值，但不能覆盖排序代理', () => {
    const api = { list: vi.fn(), remove: vi.fn(), removeBatch: vi.fn() };
    useCrudGrid<Row, any>({
      api: api as any,
      entityName: 'x',
      formOptions: { schema: [] },
      gridOptions: { checkboxConfig: { highlight: false }, height: 500 },
    });
    const options = mocks.captured.gridOptions;
    expect(options.checkboxConfig).toEqual({ highlight: false });
    expect(options.height).toBe(500);
    expect(options.proxyConfig.sort).toBe(true);
  });

  it('query 转换分页、排序与表单值，并把 list 映射为 items', async () => {
    const { api } = setup();
    const query = mocks.captured.gridOptions.proxyConfig.ajax.query;
    const result = await query(
      {
        page: { currentPage: 2, pageSize: 20 },
        sort: { field: 'name', order: 'desc' },
      },
      { enabled: 'false', name: '', role_ids: [1, 2] },
    );
    expect(api.list).toHaveBeenCalledWith({
      enabled: false,
      page: 2,
      pageSize: 20,
      role_ids: '1,2',
      sortField: 'name',
      sortOrder: 'descend',
    });
    expect(result).toEqual({ items: [{ id: 1 }], total: 1 });
  });

  it('没有排序时 sortField 与 sortOrder 为 undefined，空结果兜底为空数组', async () => {
    const { api } = setup();
    api.list.mockResolvedValue({});
    const query = mocks.captured.gridOptions.proxyConfig.ajax.query;
    const result = await query({ page: { currentPage: 1, pageSize: 10 } }, {});
    expect(api.list).toHaveBeenCalledWith({
      page: 1,
      pageSize: 10,
      sortField: undefined,
      sortOrder: undefined,
    });
    expect(result).toEqual({ items: [], total: 0 });
  });
});

describe('useCrudGrid 删除', () => {
  beforeEach(() => vi.clearAllMocks());

  it('remove 删除后提示并刷新', async () => {
    const { api, remove } = setup();
    await remove({ id: 5 });
    expect(api.remove).toHaveBeenCalledWith(5);
    expect(mocks.message.success).toHaveBeenCalledWith('删除成功');
    expect(mocks.gridApi.query).toHaveBeenCalled();
  });

  it('canDelete 为 false 的行直接忽略', async () => {
    const { api, remove } = setup((row) => !row.system);
    await remove({ id: 1, system: true });
    expect(api.remove).not.toHaveBeenCalled();
    expect(mocks.gridApi.query).not.toHaveBeenCalled();
  });

  it('removeSelected 只提交允许删除的勾选行', async () => {
    const { api, removeSelected } = setup((row) => !row.system);
    mocks.gridApi.grid.getCheckboxRecords.mockReturnValue([
      { id: 1, system: true },
      { id: 2 },
      { id: 3 },
    ]);
    await removeSelected();
    expect(api.removeBatch).toHaveBeenCalledWith([2, 3]);
    expect(mocks.gridApi.query).toHaveBeenCalled();
  });

  it('没有可删除的勾选行时提示并不请求', async () => {
    const { api, removeSelected } = setup((row) => !row.system);
    mocks.gridApi.grid.getCheckboxRecords.mockReturnValue([
      { id: 1, system: true },
    ]);
    await removeSelected();
    expect(mocks.message.warning).toHaveBeenCalledWith('请选择角色');
    expect(api.removeBatch).not.toHaveBeenCalled();
  });
});
