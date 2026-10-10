import type { UserInfo } from '@vben/types';

import { preferences } from '@vben/preferences';

import { requestClient } from '#/api/request';

export interface CurrentUser {
  avatar: string;
  enabled: boolean;
  id: number;
  name: string;
  permissions: string[];
  username: string;
}

// 把后台用户转成布局层需要的字段。权限码放到 roles，供路由守卫读取。
function toUserInfo(user: CurrentUser): UserInfo {
  return {
    avatar: user.avatar || '',
    desc: '',
    homePath: preferences.app.defaultHomePath,
    realName: user.name || user.username,
    roles: user.permissions ?? [],
    token: '',
    userId: String(user.id),
    username: user.username,
  };
}

export interface UpdateProfileParams {
  avatar?: string;
  name?: string;
  username?: string;
}

export interface ChangePasswordParams {
  newPassword: string;
  oldPassword: string;
}

// 获取当前登录用户。后台把资料和权限放在同一个接口里。
export async function getUserInfoApi() {
  const user = await requestClient.get<CurrentUser>('/auth/current-user');
  return toUserInfo(user);
}

// 更新姓名、用户名或头像。没传的字段后台会保持原值。
export async function updateProfileApi(data: UpdateProfileParams) {
  const user = await requestClient.put<CurrentUser>('/auth/profile', data);
  return toUserInfo(user);
}

// 先上传文件拿到地址，再写入用户资料。上传接口本身不改头像字段。
export async function uploadAvatarApi(file: File) {
  const uploaded = await requestClient.upload<{ url: string }>('/auth/avatar', {
    avatar: file,
    file,
  });
  return updateProfileApi({ avatar: uploaded.url });
}

// 修改当前用户密码。成功后服务端会使当前 token 失效。
export async function changePasswordApi(data: ChangePasswordParams) {
  return requestClient.put('/auth/password', data);
}
