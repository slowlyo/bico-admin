import request from '@/utils/http'

/** 下载 Excel 模板 */
export function fetchDownloadDemoExcelTemplate() {
  return request.get({
    url: '/admin-api/demo/excel/template',
    responseType: 'blob'
  })
}

/** 导出 Excel */
export function fetchExportDemoExcel() {
  return request.get({
    url: '/admin-api/demo/excel/export',
    responseType: 'blob'
  })
}

/** 导入 Excel */
export function fetchImportDemoExcel(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post({
    url: '/admin-api/demo/excel/import',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
