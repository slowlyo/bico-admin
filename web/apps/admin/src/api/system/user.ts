import type { definitions, paths } from '#/api/generated/schema';

import type { Defined, PageQuery } from '#/api/crud';

import { createCrudApi } from '#/api/crud';
import { requestClient } from '#/api/request';

// 以下类型来自后端 Swagger（api/generated/schema.d.ts），后端改字段后执行 make web-types 同步。
export type AdminUserRole = Defined<definitions['model.AdminRole']>;

export type AdminUser = Omit<Defined<definitions['model.AdminUser']>, 'roles'> & {
  roles?: AdminUserRole[];
};

export type AdminUserQuery = NonNullable<
  paths['/admin-users']['get']['parameters']['query']
> &
  PageQuery;

export type AdminUserCreatePayload = definitions['handler.createUserReq'];

export type AdminUserUpdatePayload = definitions['handler.updateUserReq'];

// 用户走标准 CRUD 接口。角色筛选是逗号分隔的 id；密码只在创建时必填，更新留空时不传。
export const adminUserApi = createCrudApi<
  AdminUser,
  AdminUserQuery,
  AdminUserCreatePayload,
  AdminUserUpdatePayload
>('/admin-users');

// 只拿上传后的地址，不改当前登录用户的头像。
export function uploadUserAvatarApi(file: File) {
  return requestClient.upload<{ url: string }>('/auth/avatar', {
    avatar: file,
    file,
  });
}
