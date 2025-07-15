import { request } from '@umijs/max';

// 上传响应
export interface UploadResponse {
  file_name: string; // 原始文件名
  file_path: string; // 文件访问路径
  file_size: number; // 文件大小（字节）
  file_type: string; // 文件类型
}

// 多文件上传响应
export interface MultiUploadResponse {
  files: UploadResponse[];
  total: number;
}

// API响应格式
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

/**
 * 上传文件
 * @param files 文件列表
 * @param dir 上传目录（可选）
 */
export async function uploadFiles(
  files: File[],
  dir?: string,
  options?: { [key: string]: any }
): Promise<ApiResponse<MultiUploadResponse>> {
  const formData = new FormData();
  
  // 添加文件
  files.forEach(file => {
    formData.append('files', file);
  });
  
  // 添加目录参数
  if (dir) {
    formData.append('dir', dir);
  }

  return request<ApiResponse<MultiUploadResponse>>('/admin/upload', {
    method: 'POST',
    data: formData,
    ...(options || {}),
  });
}

/**
 * 上传单个文件
 * @param file 文件
 * @param dir 上传目录（可选）
 */
export async function uploadFile(
  file: File,
  dir?: string,
  options?: { [key: string]: any }
): Promise<ApiResponse<MultiUploadResponse>> {
  return uploadFiles([file], dir, options);
}

/**
 * 上传头像
 * @param file 头像文件
 */
export async function uploadAvatar(
  file: File,
  options?: { [key: string]: any }
): Promise<ApiResponse<MultiUploadResponse>> {
  return uploadFile(file, 'avatars', options);
}
