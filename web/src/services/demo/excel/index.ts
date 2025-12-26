import { request } from '@umijs/max';
import { buildApiUrl } from '@/services/config';

export interface DemoExcelImportResult {
  total: number;
  preview: string[][];
}

export async function downloadDemoExcelTemplate() {
  return request(buildApiUrl('/demo/excel/template'), {
    method: 'GET',
    responseType: 'blob',
    getResponse: true,
  });
}

export async function exportDemoExcel() {
  return request(buildApiUrl('/demo/excel/export'), {
    method: 'GET',
    responseType: 'blob',
    getResponse: true,
  });
}

export async function importDemoExcel(file: File) {
  const formData = new FormData();
  formData.append('file', file);

  return request<API.Response<DemoExcelImportResult>>(buildApiUrl('/demo/excel/import'), {
    method: 'POST',
    data: formData,
  });
}
