<script lang="ts" setup>
import type { DemoExcelImportResult } from '#/api/demo/excel';

import { ref } from 'vue';

import { JsonViewer, Page, useVbenModal } from '@vben/common-ui';

import { Button, Upload, message } from 'antdv-next';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  downloadDemoExcelTemplateApi,
  exportDemoExcelApi,
  importDemoExcelApi,
} from '#/api/demo/excel';

const headers = ['姓名', '手机号', '年龄', '城市'];
const rows = [
  { age: 32, city: '深圳', id: 1, name: '王五', phone: '13600000000' },
  { age: 25, city: '杭州', id: 2, name: '赵六', phone: '13700000000' },
  { age: 41, city: '成都', id: 3, name: '钱七', phone: '13500000000' },
];

const importing = ref(false);
const exporting = ref(false);
const importResult = ref<DemoExcelImportResult>();

const [Grid, gridApi] = useVbenVxeGrid({
  gridOptions: {
    columns: [
      { type: 'checkbox', width: 48 },
      { field: 'name', title: '姓名' },
      { field: 'phone', title: '手机号' },
      { field: 'age', title: '年龄' },
      { field: 'city', title: '城市' },
    ],
    data: rows,
    height: 280,
    pagerConfig: { enabled: false },
    rowConfig: { keyField: 'id' },
    toolbarConfig: { custom: false, refresh: false, zoom: false },
  },
});

const [ImportModal, importModalApi] = useVbenModal({
  footer: false,
  title: '导入',
});

// 从响应头取文件名。没有时用调用方给的默认名。
function filenameFromHeader(header: string | undefined, fallback: string) {
  if (!header) {
    return fallback;
  }
  const encoded = /filename\*=UTF-8''([^;]+)/i.exec(header);
  if (encoded?.[1]) {
    return decodeURIComponent(encoded[1]);
  }
  const plain = /filename="?([^";]+)"?/i.exec(header);
  return plain?.[1] || fallback;
}

// 触发浏览器下载。Blob 来自原始响应体。
function saveBlob(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  link.click();
  URL.revokeObjectURL(url);
}

// 下载后端生成的表头模板。
async function downloadTemplate() {
  const response = await downloadDemoExcelTemplateApi();
  saveBlob(
    response.data,
    filenameFromHeader(response.headers['content-disposition'], '导入模板_示例.xlsx'),
  );
}

// 按勾选行导出。没勾选时后台导出全部示例行。
async function exportRows() {
  exporting.value = true;
  try {
    const selected = (gridApi.grid?.getCheckboxRecords?.() ?? []) as Array<{
      id: number;
    }>;
    const response = await exportDemoExcelApi(selected.map((row) => row.id));
    saveBlob(
      response.data,
      filenameFromHeader(response.headers['content-disposition'], '导出_示例.xlsx'),
    );
  } finally {
    exporting.value = false;
  }
}

// 拦截上传组件自己的提交，改走示例导入接口。
async function beforeImport(file: File) {
  importing.value = true;
  try {
    importResult.value = await importDemoExcelApi(file);
    importModalApi.close();
    message.success(`导入解析成功，共 ${importResult.value.total} 行`);
  } finally {
    importing.value = false;
  }
  return false;
}
</script>

<template>
  <Page
    description="下载模板、上传解析、按勾选行导出。导入只返回预览，不写入数据库。"
    title="Excel 导入/导出示例"
  >
    <Grid>
      <template #toolbar-tools>
        <Button :loading="exporting" type="primary" @click="exportRows">
          导出
        </Button>
        <Button class="ml-2" @click="importModalApi.open()">导入</Button>
      </template>
    </Grid>

    <p class="text-muted-foreground mt-4 text-sm">
      模板表头：{{ headers.join(' / ') }}
    </p>

    <div v-if="importResult" class="mt-4">
      <div class="mb-2 flex items-center justify-between">
        <span class="font-medium">导入预览（前 5 行）</span>
        <Button size="small" @click="importResult = undefined">清空</Button>
      </div>
      <JsonViewer boxed :expand-depth="2" :value="importResult.preview" />
    </div>

    <ImportModal>
      <Upload.Dragger
        accept=".xlsx,.xls,.csv"
        :before-upload="beforeImport"
        :disabled="importing"
        :show-upload-list="false"
      >
        <p>点击或拖拽文件到这里</p>
      </Upload.Dragger>
      <Button class="mt-3" type="link" @click="downloadTemplate">下载模板</Button>
    </ImportModal>
  </Page>
</template>
