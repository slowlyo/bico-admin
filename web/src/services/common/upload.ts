import { request } from '@umijs/max';
import { buildApiUrl } from '../config';

export type UploadType = 'image' | 'video';

export interface WangEditorUploadResponse {
  errno: number;
  message?: string;
  data?: {
    url: string;
  };
}

/**
 * 通用上传（用于富文本编辑器图片/视频上传）
 */
export async function uploadForEditor(file: File, type: UploadType) {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('type', type);

  return request<WangEditorUploadResponse>(buildApiUrl('/upload'), {
    method: 'POST',
    data: formData,
  });
}
