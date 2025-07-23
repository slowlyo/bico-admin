# FileUpload 文件上传组件

## 概述

`FileUpload` 是一个功能强大的文件上传组件，支持拖拽上传、多文件上传、文件类型限制等功能。基于 Element Plus 的 Upload 组件封装，提供了更好的用户体验和更丰富的功能。

## 主要特性

- ✅ **拖拽上传**: 支持将文件拖拽到上传区域
- ✅ **多文件上传**: 支持同时上传多个文件
- ✅ **文件类型限制**: 可配置允许的文件类型
- ✅ **文件大小限制**: 可配置单个文件最大大小
- ✅ **自定义文件列表**: 提供美观的文件列表展示
- ✅ **文件预览**: 支持文件预览功能
- ✅ **智能图标**: 根据文件类型显示对应图标
- ✅ **灵活配置**: 支持拖拽模式和按钮模式

## 使用方法

### 基础拖拽上传

```vue
<template>
  <FileUpload
    v-model="fileList"
    :multiple="true"
    :limit="5"
    upload-dir="documents"
    :max-size="10"
    :allowed-types="['.pdf', '.doc', '.docx']"
    @change="handleFileChange"
    @success="handleUploadSuccess"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import FileUpload from '@/components/custom/file-upload/index.vue'
import type { UploadFile } from 'element-plus'

const fileList = ref<UploadFile[]>([])

const handleFileChange = (files: UploadFile[]) => {
  console.log('文件列表变化:', files)
}

const handleUploadSuccess = (response: any, file: UploadFile) => {
  console.log('文件上传成功:', response, file)
}
</script>
```

### 按钮模式上传

```vue
<template>
  <FileUpload
    v-model="fileList"
    :drag="false"
    :multiple="false"
    button-text="选择文档"
    upload-dir="attachments"
    :allowed-types="['image/*', '.pdf']"
  />
</template>
```

### 图片上传专用

```vue
<template>
  <FileUpload
    v-model="imageFiles"
    drag-text="拖拽图片到此处上传"
    :allowed-types="['image/jpeg', 'image/png', 'image/gif']"
    :max-size="5"
    upload-dir="images"
  />
</template>
```

## Props 参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `modelValue` | `UploadFile[]` | `[]` | 文件列表，支持 v-model |
| `disabled` | `boolean` | `false` | 是否禁用上传 |
| `drag` | `boolean` | `true` | 是否支持拖拽上传 |
| `multiple` | `boolean` | `true` | 是否支持多选 |
| `limit` | `number` | `10` | 最大上传数量 |
| `showFileList` | `boolean` | `false` | 是否显示 Element 默认文件列表 |
| `showTips` | `boolean` | `true` | 是否显示提示信息 |
| `tipText` | `string` | `''` | 自定义提示文本 |
| `formatText` | `string` | `''` | 格式提示文本（自动生成） |
| `dragText` | `string` | `''` | 拖拽区域文本 |
| `buttonText` | `string` | `''` | 按钮文本 |
| `uploadDir` | `string` | `'files'` | 上传目录 |
| `maxSize` | `number` | `10` | 最大文件大小（MB） |
| `allowedTypes` | `string[]` | `[]` | 允许的文件类型 |

## Events 事件

| 事件名 | 参数 | 说明 |
|--------|------|------|
| `update:modelValue` | `(files: UploadFile[])` | 文件列表更新时触发 |
| `change` | `(files: UploadFile[])` | 文件列表变化时触发 |
| `success` | `(response: any, file: UploadFile)` | 单个文件上传成功时触发 |
| `error` | `(error: any, file: UploadFile)` | 单个文件上传失败时触发 |
| `remove` | `(file: UploadFile)` | 移除文件时触发 |
| `preview` | `(file: UploadFile)` | 预览文件时触发 |

## 文件类型配置

### 支持的文件类型格式

```typescript
// MIME 类型
allowedTypes: ['image/jpeg', 'image/png', 'application/pdf']

// 文件扩展名
allowedTypes: ['.jpg', '.png', '.pdf', '.doc', '.docx']

// 通配符
allowedTypes: ['image/*', 'video/*']

// 混合使用
allowedTypes: ['image/*', '.pdf', '.doc', '.docx']
```

### 常用文件类型预设

```typescript
// 图片文件
const imageTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']

// 文档文件
const documentTypes = ['.pdf', '.doc', '.docx', '.xls', '.xlsx', '.ppt', '.pptx']

// 压缩文件
const archiveTypes = ['.zip', '.rar', '.7z', '.tar', '.gz']

// 媒体文件
const mediaTypes = ['video/*', 'audio/*']
```

## 样式定制

### CSS 变量

```scss
.file-upload {
  // 拖拽区域
  --upload-drag-height: 180px;
  --upload-drag-border: 2px dashed var(--el-border-color);
  --upload-drag-radius: 8px;
  
  // 文件列表
  --file-item-padding: 12px 16px;
  --file-icon-size: 20px;
  --file-name-size: 14px;
}
```

### 自定义样式

```vue
<style>
.custom-file-upload .file-uploader :deep(.el-upload-dragger) {
  height: 120px;
  border-style: solid;
}

.custom-file-upload .drag-upload-area .upload-icon {
  color: #409eff;
}
</style>
```

## 高级用法

### 自定义文件预览

```vue
<template>
  <FileUpload
    v-model="fileList"
    @preview="handleCustomPreview"
  />
</template>

<script setup lang="ts">
const handleCustomPreview = (file: UploadFile) => {
  if (file.name.endsWith('.pdf')) {
    // 自定义 PDF 预览逻辑
    openPdfViewer(file.url)
  } else {
    // 默认预览
    window.open(file.url, '_blank')
  }
}
</script>
```

### 文件上传进度

```vue
<template>
  <FileUpload
    v-model="fileList"
    :show-file-list="true"
    @success="handleSuccess"
    @error="handleError"
  />
</template>
```

## 注意事项

- 组件依赖 `useUserStore` 获取认证 token
- 上传接口路径为 `/admin-api/upload`
- 拖拽功能需要现代浏览器支持
- 文件类型检查同时支持 MIME 类型和文件扩展名
- 自定义文件列表提供了更好的视觉效果和交互体验
