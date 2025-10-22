import { request } from '@umijs/max';
import { buildApiUrl } from '../config';
import type {
  UpdateProfileParams,
  ChangePasswordParams,
  CurrentUser,
} from './types';

/**
 * 更新用户资料
 */
export async function updateProfile(data: UpdateProfileParams) {
  return request<API.Response<CurrentUser>>(buildApiUrl('/auth/profile'), {
    method: 'PUT',
    data,
  });
}

/**
 * 修改密码
 */
export async function changePassword(data: ChangePasswordParams) {
  return request<API.Response<any>>(buildApiUrl('/auth/password'), {
    method: 'PUT',
    data,
  });
}

/**
 * 上传头像
 */
export async function uploadAvatar(file: File) {
  const formData = new FormData();
  formData.append('avatar', file);

  return request<API.Response<{ url: string }>>(buildApiUrl('/auth/avatar'), {
    method: 'POST',
    data: formData,
  });
}
