<template>
  <div class="art-full-height">
    <ElCard shadow="never" class="art-card">
      <template #header>
        <div class="flex flex-col gap-1">
          <span class="text-lg font-medium">富文本编辑器示例</span>
          <span class="text-sm text-g-500">演示 RichEditor 的基本用法：受控 value + onChange，并支持图片/视频上传。</span>
        </div>
      </template>

      <div class="mb-6">
        <ArtWangEditor v-model="html" height="400px" />
      </div>

      <div class="mt-6">
        <div class="text-base font-medium mb-3">当前 HTML 内容</div>
        <div class="p-4 bg-g-100 rounded-md border border-g-200 overflow-auto">
          <pre class="text-sm text-g-700 whitespace-pre-wrap break-all">{{ html }}</pre>
        </div>
      </div>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import ArtWangEditor from '@/components/core/forms/art-wang-editor/index.vue'

  defineOptions({ name: 'DemoEditor' })

  const html = ref('<p>这是一个示例内容</p>')
</script>
