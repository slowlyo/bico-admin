<script lang="ts" setup>
import { ref } from 'vue';

import { Tabs } from 'antdv-next';

const docType = ref<'admin' | 'api'>('admin');

const docs = {
  admin: '/swagger/admin/index.html',
  api: '/swagger/api/index.html',
};
</script>

<template>
  <div class="flex h-[calc(100vh-120px)] min-h-[640px] flex-col p-4">
    <div class="bg-background flex min-h-0 flex-1 flex-col overflow-hidden rounded-lg border">
      <Tabs
        v-model:active-key="docType"
        style="padding-inline: 16px"
        :items="[
          { key: 'admin', label: 'Admin API' },
          { key: 'api', label: 'Open API' },
        ]"
      />
      <iframe
        class="min-h-0 w-full flex-1 border-0"
        :src="docs[docType]"
        :title="docType === 'admin' ? 'Admin API 文档' : 'Open API 文档'"
      />
    </div>
  </div>
</template>
