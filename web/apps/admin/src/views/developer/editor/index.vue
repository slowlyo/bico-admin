<script lang="ts" setup>
import { ref } from 'vue';

import { Page } from '@vben/common-ui';
import { VbenTiptap } from '@vben/plugins/tiptap';

import { uploadEditorFileApi } from '#/api/common';

const html = ref('<p>这是一个示例内容</p>');

// 编辑器图片走通用上传，返回可插入的地址。
async function uploadImage(file: File) {
  const uploaded = await uploadEditorFileApi(file);
  return uploaded.url;
}
</script>

<template>
  <Page
    description="基本排版和图片上传。下方是当前 HTML。"
    title="富文本编辑器示例"
  >
    <VbenTiptap
      v-model="html"
      :image-upload="{ upload: uploadImage }"
      :min-height="240"
    />
    <pre class="bg-muted mt-4 overflow-auto rounded-md border p-3 text-sm">{{
      html
    }}</pre>
  </Page>
</template>
