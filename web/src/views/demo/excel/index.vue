<template>
  <div class="art-full-height">
    <ElCard shadow="never" class="art-card">
      <template #header>
        <div class="flex flex-col gap-1">
          <span class="text-lg font-medium">Excel 导入/导出示例</span>
          <span class="text-sm text-g-500">演示：下载模板（后端通过 Excel 库输出表头） / 拖拽上传导入 / 导出（按钮 loading 处理）。</span>
        </div>
      </template>

      <div class="flex flex-wrap gap-4 mb-6">
        <ElButton type="primary" @click="handleDownloadTemplate" v-ripple>
          <template #icon><ArtSvgIcon icon="ri:download-line" /></template>
          下载导入模板
        </ElButton>

        <ArtExcelImport
          v-model:visible="importVisible"
          title="导入 Excel"
          @upload="handleImportFile"
        >
          <template #trigger>
            <ElButton type="success" :loading="importing" v-ripple>
              <template #icon><ArtSvgIcon icon="ri:upload-line" /></template>
              导入数据
            </ElButton>
          </template>
        </ArtExcelImport>

        <ElButton type="warning" :loading="exporting" @click="handleExport" v-ripple>
          <template #icon><ArtSvgIcon icon="ri:file-excel-2-line" /></template>
          导出示例
        </ElButton>
      </div>

      <div class="mb-6">
        <div class="text-base font-medium mb-3">模板表头</div>
        <div class="flex flex-wrap gap-2">
          <ElTag v-for="h in headers" :key="h" effect="plain">{{ h }}</ElTag>
        </div>
      </div>

      <div v-if="importResult" class="mt-6">
        <div class="flex items-center justify-between mb-4">
          <span class="text-base font-medium">导入预览（前 5 行）</span>
          <ElButton size="small" @click="importResult = undefined">清空导入结果</ElButton>
        </div>
        <div class="p-4 bg-g-100 rounded-md overflow-auto">
          <pre class="text-sm text-g-700">{{ JSON.stringify(importResult.preview, null, 2) }}</pre>
        </div>
      </div>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed } from 'vue'
  import { ElMessage } from 'element-plus'
  import {
    fetchDownloadDemoExcelTemplate,
    fetchExportDemoExcel,
    fetchImportDemoExcel
  } from '@/api/demo'
  import { downloadBlob } from '@/utils'

  defineOptions({ name: 'DemoExcel' })

  const importing = ref(false)
  const exporting = ref(false)
  const importVisible = ref(false)
  const importResult = ref<{ total: number; preview: string[][] }>()

  const headers = ['姓名', '手机号', '年龄', '城市']

  /** 下载模板 */
  const handleDownloadTemplate = async () => {
    try {
      const resp: any = await fetchDownloadDemoExcelTemplate()
      const blob = resp
      const filename = '导入模板_示例.xlsx'
      downloadBlob(blob, filename)
    } catch (e: any) {
      ElMessage.error('下载模板失败')
    }
  }

  /** 导出 */
  const handleExport = async () => {
    exporting.value = true
    try {
      const resp: any = await fetchExportDemoExcel()
      const blob = resp
      const filename = '导出_示例.xlsx'
      downloadBlob(blob, filename)
    } catch (e: any) {
      ElMessage.error('导出失败')
    } finally {
      exporting.value = false
    }
  }

  /** 导入 */
  const handleImportFile = async (file: File) => {
    importing.value = true
    try {
      const res: any = await fetchImportDemoExcel(file)
      if (res.code === 0 && res.data) {
        ElMessage.success(`导入解析成功，共 ${res.data.total} 行`)
        importResult.value = res.data
        importVisible.value = false
      } else {
        ElMessage.error(res.msg || '导入失败')
      }
    } catch (e: any) {
      ElMessage.error('导入失败')
    } finally {
      importing.value = false
    }
  }
</script>
