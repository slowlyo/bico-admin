import type { CrudApi, PageQuery } from '#/api/crud';

import { message } from 'antdv-next';

import { useVbenVxeGrid } from '#/adapter/vxe-table';

// 取适配器 useVbenVxeGrid 的入参类型，保证表单组件类型与全局适配一致。
type GridProps<T extends Record<string, any>> = NonNullable<
  Parameters<typeof useVbenVxeGrid<T>>[0]
>;

interface CrudGridOptions<T extends { id: number }, Q extends PageQuery> {
  // 只用到列表与删除接口，创建和更新由弹窗表单自己调用。
  api: Pick<CrudApi<T, Q, unknown>, 'list' | 'remove' | 'removeBatch'>;
  // 行是否允许删除，同时用于批量删除时过滤勾选行（如系统保留角色）。
  canDelete?: (row: T) => boolean;
  // 用于“请选择xx”提示文案，如“用户”。
  entityName: string;
  formOptions: GridProps<T>['formOptions'];
  gridOptions: GridProps<T>['gridOptions'];
}

// 把搜索表单的值转成后端查询参数：丢弃空值，enabled 字符串转布尔，数组转逗号分隔。
function normalizeFormValues(values: Record<string, any>) {
  const result: Record<string, any> = {};
  for (const [key, value] of Object.entries(values)) {
    if (value === undefined || value === null || value === '') {
      // 空值不传，避免后端把空字符串当成有效筛选条件。
      continue;
    }
    if (key === 'enabled') {
      // Select 的选项值是字符串，后端 enabled 是布尔。
      result[key] = value === 'true' || value === true;
    } else if (Array.isArray(value)) {
      // 多选筛选后端按逗号分隔的 id 解析。
      if (value.length > 0) result[key] = value.join(',');
    } else {
      result[key] = value;
    }
  }
  return result;
}

// vxe 的排序方向 asc/desc 对应后端的 ascend/descend。
function toSortOrder(order?: null | string) {
  if (order === 'asc') return 'ascend';
  if (order === 'desc') return 'descend';
  return undefined;
}

// 标准 CRUD 列表页的通用封装：分页、服务端排序、搜索表单映射、单删与批量删除。
// 新建与编辑弹窗仍由页面自己决定，成功后调用 gridApi.query() 刷新。
export function useCrudGrid<T extends { id: number }, Q extends PageQuery>(
  options: CrudGridOptions<T, Q>,
) {
  const { api, canDelete, entityName, formOptions, gridOptions } = options;

  const [Grid, gridApi] = useVbenVxeGrid<T>({
    formOptions,
    gridOptions: {
      checkboxConfig: { highlight: true },
      height: 'auto',
      rowConfig: { keyField: 'id' },
      toolbarConfig: { custom: true, refresh: true, search: true, zoom: true },
      ...gridOptions,
      proxyConfig: {
        ajax: {
          // vxe 读取 items/total，后台分页字段是 list/total。
          query: async (
            {
              page,
              sort,
            }: {
              page: { currentPage: number; pageSize: number };
              sort?: { field?: string; order?: null | string };
            },
            formValues: Record<string, any>,
          ) => {
            const result = await api.list({
              ...normalizeFormValues(formValues),
              page: page.currentPage,
              pageSize: page.pageSize,
              sortField: sort?.field || undefined,
              sortOrder: toSortOrder(sort?.order),
            } as Q);
            return { items: result.list ?? [], total: result.total ?? 0 };
          },
        },
        // 开启排序代理，点击表头排序时会重新调用 query 并携带 sort 参数。
        sort: true,
      },
      // 服务端排序：不在浏览器里对当前页数据排序，排序字段需在后端白名单内。
      sortConfig: { remote: true },
    },
  });

  // 删除单行后刷新当前查询条件；不可删除的行直接忽略，保持与按钮禁用一致。
  async function remove(row: T) {
    if (canDelete && !canDelete(row)) return;
    await api.remove(row.id);
    message.success('删除成功');
    await gridApi.query();
  }

  // 批量删除只处理当前勾选且允许删除的行，避免整批被后端因保留记录拒绝。
  async function removeSelected() {
    const checked = (gridApi.grid?.getCheckboxRecords?.() ?? []) as T[];
    const rows = canDelete ? checked.filter((row) => canDelete(row)) : checked;
    if (rows.length === 0) {
      message.warning(`请选择${entityName}`);
      return;
    }
    await api.removeBatch(rows.map((row) => row.id));
    message.success('删除成功');
    await gridApi.query();
  }

  return { Grid, gridApi, remove, removeSelected };
}
