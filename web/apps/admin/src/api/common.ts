import { requestClient } from '#/api/request';

export interface AppConfig {
  debug: boolean;
  logo: string;
  name: string;
}

// 应用配置是公开接口。debug 为真时才显示开发者工具菜单。
export function getAppConfigApi() {
  return requestClient.get<AppConfig>('/app-config');
}

// 富文本图片上传。字段名用 file，和通用上传接口一致。
export function uploadEditorFileApi(file: File) {
  return requestClient.upload<{ url: string }>('/upload', { file });
}
