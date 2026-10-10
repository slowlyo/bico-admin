<script lang="ts" setup>
import type { DashboardOverview } from '#/api/dashboard';

import { computed, onMounted, ref } from 'vue';

import { useUserStore } from '@vben/stores';

import {
  Card,
  Descriptions,
  DescriptionsItem,
  Statistic,
  Tag,
} from 'antdv-next';

import { getDashboardOverviewApi } from '#/api/dashboard';

const userStore = useUserStore();
const overview = ref<DashboardOverview>();

const displayName = computed(
  () => userStore.userInfo?.realName || userStore.userInfo?.username || '',
);

// 把秒数收成小时和分钟，工作台不展示原始秒数。
function formatUptime(seconds: number) {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  return `${hours}小时${minutes}分`;
}

// warning 用红色，其余保持中性。状态只来自后台给的 status。
function statusColor(status: string) {
  return status === 'warning' ? '#cf1322' : '#3f8600';
}

// 进入工作台就拉概览。失败时拦截器会提示，页面留在欢迎卡片。
async function loadOverview() {
  overview.value = await getDashboardOverviewApi();
}

onMounted(() => {
  loadOverview();
});
</script>

<template>
  <div class="flex flex-col gap-4 p-4">
    <Card>
      <h2 class="text-lg font-medium">欢迎回来，{{ displayName }}</h2>
      <p class="text-muted-foreground mt-2">这是你的工作台。</p>
    </Card>

    <div
      v-if="overview"
      class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4"
    >
      <Card
        v-for="metric in overview.monitor.metrics"
        :key="metric.key"
        class="metric-card h-full"
      >
        <Statistic
          :title="metric.label"
          :value="metric.value"
          :suffix="metric.unit"
          :value-style="{ color: statusColor(metric.status) }"
        />
      </Card>
    </div>

    <div v-if="overview" class="grid grid-cols-1 gap-4 lg:grid-cols-3">
      <Card title="服务器" class="h-full">
          <Descriptions :column="1" size="small">
            <DescriptionsItem label="主机">
              {{ overview.server.hostname }}
            </DescriptionsItem>
            <DescriptionsItem label="系统">
              {{ overview.server.os }} / {{ overview.server.arch }}
            </DescriptionsItem>
            <DescriptionsItem label="Go">
              {{ overview.server.goVersion }}
            </DescriptionsItem>
            <DescriptionsItem label="模式">
              <Tag>{{ overview.server.mode }}</Tag>
            </DescriptionsItem>
            <DescriptionsItem label="运行">
              {{ formatUptime(overview.server.uptimeSeconds) }}
            </DescriptionsItem>
          </Descriptions>
      </Card>
      <Card title="运行时" class="h-full">
          <Descriptions :column="1" size="small">
            <DescriptionsItem label="CPU">
              {{ overview.runtime.cpuCores }}
            </DescriptionsItem>
            <DescriptionsItem label="协程">
              {{ overview.runtime.goroutines }}
            </DescriptionsItem>
            <DescriptionsItem label="内存">
              {{ overview.runtime.allocMb }} / {{ overview.runtime.sysMb }} MB
            </DescriptionsItem>
            <DescriptionsItem label="GC">
              {{ overview.runtime.gcCycles }}
            </DescriptionsItem>
          </Descriptions>
      </Card>
      <Card title="数据库" class="h-full">
          <Descriptions :column="1" size="small">
            <DescriptionsItem label="驱动">
              {{ overview.database.driver }}
            </DescriptionsItem>
            <DescriptionsItem label="连接">
              {{ overview.database.openConnections }} /
              {{ overview.database.maxOpenConnections }}
            </DescriptionsItem>
            <DescriptionsItem label="使用中">
              {{ overview.database.inUse }}
            </DescriptionsItem>
            <DescriptionsItem label="空闲">
              {{ overview.database.idle }}
            </DescriptionsItem>
            <DescriptionsItem label="等待">
              {{ overview.database.waitCount }}
            </DescriptionsItem>
          </Descriptions>
      </Card>
    </div>
  </div>
</template>

<style scoped>
/* 单位和数字保持同一行，避免第一张卡片被挤高后底边对不齐。 */
.metric-card :deep(.ant-statistic-content) {
  white-space: nowrap;
}
</style>
