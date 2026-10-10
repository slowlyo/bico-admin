import { requestClient } from '#/api/request';

interface BlobResponse {
  data: Blob;
  headers: Record<string, string | undefined>;
}

export interface DemoExcelImportResult {
  preview: string[][];
  total: number;
}

// 模板是文件流，必须拿原始响应才能读 Content-Disposition 里的文件名。
export function downloadDemoExcelTemplateApi() {
  return requestClient.download<BlobResponse>('/demo/excel/template', {
    responseReturn: 'raw',
  });
}

// 导出同样返回文件流。ids 为空时后台导出全部示例行。
export function exportDemoExcelApi(ids: number[]) {
  return requestClient.download<BlobResponse>('/demo/excel/export', {
    params: ids.length > 0 ? { ids: ids.join(',') } : undefined,
    responseReturn: 'raw',
  });
}

// 导入只解析并回预览，不落库。
export function importDemoExcelApi(file: File) {
  return requestClient.upload<DemoExcelImportResult>('/demo/excel/import', {
    file,
  });
}
