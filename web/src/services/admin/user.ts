import { request } from '@umijs/max';

/**
 * 更新用户信息
 * PUT /admin-api/auth/profile
 */
export async function updateProfile(data: API.UpdateProfileParams) {
  return request<API.Response<API.CurrentUser>>('/admin-api/auth/profile', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
  });
}

/**
 * 修改密码
 * PUT /admin-api/auth/password
 */
export async function changePassword(data: API.ChangePasswordParams) {
  return request<API.Response<any>>('/admin-api/auth/password', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data,
  });
}

/**
 * 上传头像
 * POST /admin-api/auth/avatar
 */
export async function uploadAvatar(file: File) {
  const formData = new FormData();
  formData.append('avatar', file);
  
  return request<API.Response<{ url: string }>>('/admin-api/auth/avatar', {
    method: 'POST',
    data: formData,
  });
}
