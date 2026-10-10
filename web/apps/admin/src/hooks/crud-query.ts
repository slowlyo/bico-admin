// 把搜索表单的值转成后端查询参数：丢弃空值，enabled 字符串转布尔，数组转逗号分隔。
export function normalizeFormValues(values: Record<string, any>) {
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
export function toSortOrder(order?: null | string) {
  if (order === 'asc') return 'ascend';
  if (order === 'desc') return 'descend';
  return undefined;
}
