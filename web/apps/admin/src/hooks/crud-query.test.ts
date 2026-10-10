import { describe, expect, it } from 'vitest';

import { normalizeFormValues, toSortOrder } from './crud-query';

describe('normalizeFormValues', () => {
  it('丢弃空字符串、null 和 undefined', () => {
    expect(
      normalizeFormValues({ a: '', b: null, c: undefined, name: 'x' }),
    ).toEqual({ name: 'x' });
  });

  it('enabled 的字符串转布尔，布尔原样保留', () => {
    expect(normalizeFormValues({ enabled: 'true' })).toEqual({ enabled: true });
    expect(normalizeFormValues({ enabled: 'false' })).toEqual({
      enabled: false,
    });
    expect(normalizeFormValues({ enabled: false })).toEqual({ enabled: false });
  });

  it('数组转逗号分隔，空数组被丢弃', () => {
    expect(normalizeFormValues({ role_ids: [1, 2, 3] })).toEqual({
      role_ids: '1,2,3',
    });
    expect(normalizeFormValues({ role_ids: [] })).toEqual({});
  });

  it('数字 0 是有效值，不能当成空值丢弃', () => {
    expect(normalizeFormValues({ sort: 0 })).toEqual({ sort: 0 });
  });
});

describe('toSortOrder', () => {
  it('把 vxe 的方向映射为后端约定', () => {
    expect(toSortOrder('asc')).toBe('ascend');
    expect(toSortOrder('desc')).toBe('descend');
  });

  it('无排序或未知值返回 undefined', () => {
    expect(toSortOrder(null)).toBeUndefined();
    expect(toSortOrder(undefined)).toBeUndefined();
    expect(toSortOrder('')).toBeUndefined();
  });
});
