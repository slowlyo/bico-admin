import { request } from '@umijs/max';

// 更新个人信息请求参数
export interface UpdateProfileRequest {
  name: string;
  email?: string;
  phone?: string;
  avatar?: string;
}

// 修改密码请求参数
export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

// 更新个人信息
export async function updateProfile(data: UpdateProfileRequest) {
  return request('/admin/profile', {
    method: 'PUT',
    data,
  });
}

// 修改密码
export async function changePassword(data: ChangePasswordRequest) {
  return request('/admin/profile/password', {
    method: 'PUT',
    data,
  });
}
