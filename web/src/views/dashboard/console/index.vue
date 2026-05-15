<!-- 工作台页面 -->
<template>
  <div class="dashboard-console">
    <ElSkeleton :loading="loading" animated>
      <template #template>
        <div class="art-card hero-card mb-5">
          <ElSkeletonItem variant="text" class="!w-24" />
          <ElSkeletonItem variant="h1" class="!w-64 mt-4" />
          <ElSkeletonItem variant="text" class="!w-96 mt-4 max-sm:!w-full" />
        </div>
        <ElRow :gutter="20">
          <ElCol v-for="item in 4" :key="item" :xs="24" :sm="12" :lg="6">
            <div class="art-card h-34 px-5 py-4 mb-5 max-sm:mb-4">
              <ElSkeletonItem variant="text" class="!w-20" />
              <ElSkeletonItem variant="h1" class="!w-24 mt-4" />
              <ElSkeletonItem variant="text" class="!w-32 mt-4" />
            </div>
          </ElCol>
        </ElRow>
      </template>

      <template #default>
        <section class="art-card hero-card mb-5">
          <div class="hero-main">
            <div class="status-icon" :class="healthStatus.className">
              <ArtSvgIcon :icon="healthStatus.icon" />
            </div>
            <div class="min-w-0">
              <div class="flex-c flex-wrap gap-2">
                <h2>{{ healthStatus.title }}</h2>
                <ElTag :type="healthStatus.tagType" effect="plain">{{ healthStatus.tagText }}</ElTag>
              </div>
              <p>{{ healthStatus.description }}</p>
            </div>
          </div>
          <div class="hero-side">
            <div>
              <span>已稳定运行</span>
              <strong>{{ uptimeText }}</strong>
            </div>
            <ElButton :loading="loading" @click="loadOverview" v-ripple>
              <ArtSvgIcon icon="ri:refresh-line" class="mr-1" />
              刷新
            </ElButton>
          </div>
        </section>

        <ElRow :gutter="20">
          <ElCol v-for="card in summaryCards" :key="card.key" :xs="24" :sm="12" :lg="6">
            <div class="art-card summary-card mb-5 max-sm:mb-4">
              <div class="summary-head">
                <div class="summary-icon">
                  <ArtSvgIcon :icon="card.icon" />
                </div>
                <ElTag :type="card.tagType" effect="plain" size="small">{{ card.tagText }}</ElTag>
              </div>
              <p>{{ card.label }}</p>
              <div class="summary-value">
                <ArtCountTo
                  :target="card.value"
                  :decimals="card.decimals"
                  :duration="900"
                  separator=","
                />
                <span>{{ card.unit }}</span>
              </div>
              <small>{{ card.description }}</small>
            </div>
          </ElCol>
        </ElRow>

        <ElRow :gutter="20">
          <ElCol :xs="24" :lg="15">
            <div class="art-card panel-card mb-5 max-sm:mb-4">
              <div class="panel-header">
                <div>
                  <h3>服务状态</h3>
                  <p>最近一次检查：{{ collectedTime }}</p>
                </div>
                <ElTag :type="healthStatus.tagType" effect="plain">{{ healthStatus.tagText }}</ElTag>
              </div>
              <ArtLineChart
                height="18rem"
                :data="statusChartData"
                :xAxisData="statusChartLabels"
                :showAreaColor="true"
                :showAxisLine="false"
                :isEmpty="!statusChartData.length"
              />
            </div>
          </ElCol>

          <ElCol :xs="24" :lg="9">
            <div class="art-card panel-card mb-5 max-sm:mb-4">
              <div class="panel-header">
                <div>
                  <h3>服务入口</h3>
                  <p>{{ serverAddress }}</p>
                </div>
              </div>
              <div class="server-card">
                <div class="server-name">{{ serverName }}</div>
                <div class="server-desc">当前后台服务正在响应管理端请求</div>
              </div>
              <div class="detail-list">
                <div v-for="item in serverItems" :key="item.label">
                  <span>{{ item.label }}</span>
                  <strong>{{ item.value }}</strong>
                </div>
              </div>
            </div>
          </ElCol>
        </ElRow>

        <ElRow :gutter="20">
          <ElCol :xs="24" :lg="12">
            <div class="art-card panel-card mb-5 max-sm:mb-4">
              <div class="panel-header">
                <div>
                  <h3>资源使用</h3>
                  <p>用于判断服务是否轻松运行</p>
                </div>
              </div>
              <div class="progress-list">
                <div v-for="item in resourceItems" :key="item.label">
                  <div class="progress-title">
                    <span>{{ item.label }}</span>
                    <strong>{{ item.value }}</strong>
                  </div>
                  <ElProgress
                    :percentage="item.percentage"
                    :stroke-width="7"
                    :show-text="false"
                    :color="progressColor(item.percentage)"
                    class="[&_.el-progress-bar__outer]:bg-g-200"
                  />
                  <p>{{ item.description }}</p>
                </div>
              </div>
            </div>
          </ElCol>

          <ElCol :xs="24" :lg="12">
            <div class="art-card panel-card mb-5 max-sm:mb-4">
              <div class="panel-header">
                <div>
                  <h3>数据服务</h3>
                  <p>{{ databaseDriver }}</p>
                </div>
              </div>
              <div class="detail-grid">
                <div v-for="item in databaseItems" :key="item.label" class="detail-block">
                  <span>{{ item.label }}</span>
                  <strong>{{ item.value }}</strong>
                </div>
              </div>
            </div>
          </ElCol>
        </ElRow>
      </template>
    </ElSkeleton>
  </div>
</template>

<script setup lang="ts">
  import { fetchDashboardOverview } from '@/api/system-manage'
  import { ElMessage } from 'element-plus'

  defineOptions({ name: 'Console' })

  interface SummaryCard {
    key: string
    label: string
    value: number
    unit: string
    icon: string
    decimals: number
    description: string
    tagType: 'success' | 'warning'
    tagText: string
  }

  interface InfoItem {
    label: string
    value: string
  }

  interface ResourceItem {
    label: string
    value: string
    percentage: number
    description: string
  }

  const loading = ref(true)
  const overview = ref<Api.Dashboard.Overview>()

  const metricIconMap: Record<string, string> = {
    allocMemory: 'ri:hard-drive-2-line',
    goroutines: 'ri:flow-chart',
    openConnections: 'ri:link-m',
    gcCycles: 'ri:loop-left-line'
  }

  // 页面挂载后读取一次概览数据，让首屏状态来自后端实时快照。
  onMounted(() => {
    loadOverview()
  })

  // loadOverview 获取工作台概览，失败时保留页面结构并提示用户刷新。
  async function loadOverview() {
    loading.value = true
    try {
      overview.value = await fetchDashboardOverview()
    } catch {
      // 工作台是普通用户判断系统状态的入口，接口失败必须给出明确反馈。
      ElMessage.error('工作台数据加载失败')
    } finally {
      // 请求完成后关闭加载态，避免失败场景一直停留在骨架屏。
      loading.value = false
    }
  }

  const hasWarning = computed(() => {
    return (overview.value?.monitor.metrics ?? []).some((item) => item.status === 'warning')
  })

  const healthStatus = computed(() => {
    if (!overview.value) {
      // 尚未拿到数据时用中性状态，避免把加载过程展示成异常。
      return {
        title: '正在检查系统',
        description: '正在读取服务器和数据服务状态。',
        tagText: '检查中',
        tagType: 'warning' as const,
        icon: 'ri:time-line',
        className: 'is-waiting'
      }
    }
    if (hasWarning.value) {
      // 后端返回需关注指标时，用普通用户能理解的语言提示检查。
      return {
        title: '系统需要关注',
        description: '部分服务指标偏高，建议稍后复查或联系管理员。',
        tagText: '需关注',
        tagType: 'warning' as const,
        icon: 'ri:alert-line',
        className: 'is-warning'
      }
    }
    // 所有监控指标正常时，把结论前置，减少用户理解成本。
    return {
      title: '系统运行正常',
      description: '后台服务、资源使用和数据连接都处于正常状态。',
      tagText: '正常',
      tagType: 'success' as const,
      icon: 'ri:check-line',
      className: 'is-normal'
    }
  })

  const summaryCards = computed<SummaryCard[]>(() => {
    const metrics = overview.value?.monitor.metrics ?? []
    return metrics.map((item) => ({
      key: item.key,
      label: metricLabel(item.key, item.label),
      value: item.value,
      unit: item.unit,
      icon: metricIconMap[item.key] ?? 'ri:pulse-line',
      decimals: item.unit === 'MB' ? 1 : 0,
      description: metricDescription(item.key),
      tagType: item.status === 'warning' ? 'warning' : 'success',
      tagText: item.status === 'warning' ? '需关注' : '正常'
    }))
  })

  const statusChartLabels = computed(() => resourceItems.value.map((item) => item.label))

  const statusChartData = computed(() => resourceItems.value.map((item) => item.percentage))

  const collectedTime = computed(() => formatDateTime(overview.value?.monitor.collectedAt))

  const serverName = computed(() => overview.value?.server.hostname || '后台服务')

  const serverAddress = computed(() => {
    const server = overview.value?.server
    if (!server) {
      // 数据未返回前不拼接地址，避免出现 undefined。
      return '等待数据'
    }
    return `${server.hostname}:${server.port}`
  })

  const uptimeText = computed(() => {
    const seconds = overview.value?.server.uptimeSeconds
    if (seconds === undefined) {
      // 没有启动时长时显示等待状态，避免误导用户认为刚刚启动。
      return '等待数据'
    }
    return formatDuration(seconds)
  })

  const serverItems = computed<InfoItem[]>(() => {
    const server = overview.value?.server
    if (!server) {
      // 服务信息缺失时不渲染详情，保持页面信息真实。
      return []
    }
    return [
      { label: '运行环境', value: server.mode },
      { label: '系统类型', value: `${server.os}/${server.arch}` },
      { label: '服务版本', value: server.goVersion },
      { label: '启动时间', value: formatDateTime(server.startedAt) }
    ]
  })

  const resourceItems = computed<ResourceItem[]>(() => {
    const runtime = overview.value?.runtime
    const database = overview.value?.database
    if (!runtime || !database) {
      // 资源数据不完整时不渲染进度条，避免展示无意义百分比。
      return []
    }
    return [
      {
        label: '内存使用',
        value: `${runtime.allocMb} MB`,
        percentage: percent(runtime.allocMb, runtime.nextGcMb),
        description: '数值越低，服务越轻松。'
      },
      {
        label: '处理能力',
        value: `${runtime.goMaxProcs}/${runtime.cpuCores} 核`,
        percentage: percent(runtime.goMaxProcs, runtime.cpuCores),
        description: '当前可用于处理请求的 CPU 能力。'
      },
      {
        label: '数据连接',
        value: `${database.openConnections}/${database.maxOpenConnections}`,
        percentage: percent(database.openConnections, database.maxOpenConnections),
        description: '连接越接近上限，越需要关注。'
      },
      {
        label: '后台任务',
        value: `${runtime.goroutines} 个`,
        percentage: percent(runtime.goroutines, 10000),
        description: '后台任务数量用于观察服务压力。'
      }
    ]
  })

  const databaseDriver = computed(() => {
    const driver = overview.value?.database.driver
    return driver ? `${driver} 数据连接` : '等待数据'
  })

  const databaseItems = computed<InfoItem[]>(() => {
    const database = overview.value?.database
    if (!database) {
      // 数据库信息缺失时不展示详情，避免用户误解为连接数为 0。
      return []
    }
    return [
      { label: '当前连接', value: String(database.openConnections) },
      { label: '正在使用', value: String(database.inUse) },
      { label: '空闲可用', value: String(database.idle) },
      { label: '连接上限', value: String(database.maxOpenConnections) },
      { label: '等待次数', value: String(database.waitCount) },
      { label: '等待耗时', value: `${database.waitDurationSeconds.toFixed(2)} 秒` }
    ]
  })

  // metricLabel 将后端技术指标转换成普通用户更容易理解的名称。
  function metricLabel(key: string, fallback: string) {
    const labels: Record<string, string> = {
      allocMemory: '内存使用',
      goroutines: '后台任务',
      openConnections: '数据连接',
      gcCycles: '自动清理'
    }
    return labels[key] ?? fallback
  }

  // metricDescription 为指标卡片生成短说明，减少页面技术黑话。
  function metricDescription(key: string) {
    const descriptions: Record<string, string> = {
      allocMemory: '服务当前占用的内存',
      goroutines: '后台正在处理的任务',
      openConnections: '数据库当前连接数',
      gcCycles: '服务自动整理内存次数'
    }
    return descriptions[key] ?? '服务运行指标'
  }

  // formatDateTime 格式化接口时间，空值返回占位文案以保持页面可读。
  function formatDateTime(value?: string) {
    if (!value) {
      // 接口尚未返回或字段为空时，不展示无效日期。
      return '暂无数据'
    }
    return new Date(value).toLocaleString()
  }

  // formatDuration 将秒数转成天时分格式，优先展示较大的时间单位。
  function formatDuration(seconds: number) {
    const days = Math.floor(seconds / 86400)
    const hours = Math.floor((seconds % 86400) / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)

    if (days > 0) {
      // 超过一天时保留天和小时，避免文本过长挤压状态区。
      return `${days} 天 ${hours} 小时`
    }
    if (hours > 0) {
      // 不足一天但超过一小时，展示小时和分钟便于判断近期重启。
      return `${hours} 小时 ${minutes} 分钟`
    }
    // 小于一小时的服务展示分钟，符合刚启动场景的阅读习惯。
    return `${minutes} 分钟`
  }

  // percent 计算进度百分比，边界值统一限制在 0-100。
  function percent(value: number, total: number) {
    if (!Number.isFinite(value) || !Number.isFinite(total) || total <= 0) {
      // 分母无效时返回 0，避免进度条出现 NaN。
      return 0
    }
    return Math.min(100, Math.round((value / total) * 100))
  }

  // progressColor 根据资源占用选择语义颜色，避免只靠数字判断状态。
  function progressColor(value: number) {
    if (value >= 85) {
      // 高占用需要醒目标识，提示用户关注资源压力。
      return '#e6a23c'
    }
    // 正常占用使用成功色，和页面整体健康状态保持一致。
    return '#67c23a'
  }
</script>

<style scoped lang="scss">
  .dashboard-console {
    .hero-card {
      display: flex;
      align-items: center;
      justify-content: space-between;
      min-height: 148px;
      padding: 24px;
      gap: 20px;

      @media (max-width: 768px) {
        align-items: flex-start;
        flex-direction: column;
      }
    }

    .hero-main {
      display: flex;
      align-items: center;
      min-width: 0;
      gap: 16px;

      h2 {
        margin: 0;
        font-size: 26px;
        font-weight: 650;
        color: var(--art-gray-900);
      }

      p {
        margin-top: 10px;
        font-size: 14px;
        line-height: 1.6;
        color: var(--art-gray-600);
      }
    }

    .status-icon {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 56px;
      height: 56px;
      flex: 0 0 auto;
      border-radius: 16px;
      font-size: 28px;

      &.is-normal {
        color: #16794c;
        background: rgba(22, 121, 76, 0.12);
      }

      &.is-warning,
      &.is-waiting {
        color: #9a5b00;
        background: rgba(154, 91, 0, 0.12);
      }
    }

    .hero-side {
      display: flex;
      align-items: center;
      gap: 16px;

      span {
        display: block;
        font-size: 12px;
        color: var(--art-gray-500);
      }

      strong {
        display: block;
        margin-top: 4px;
        font-size: 18px;
        font-weight: 650;
        white-space: nowrap;
        color: var(--art-gray-900);
      }
    }

    .summary-card {
      position: relative;
      min-height: 150px;
      padding: 18px 20px;

      p {
        margin-top: 18px;
        font-size: 13px;
        color: var(--art-gray-600);
      }

      small {
        display: block;
        margin-top: 8px;
        font-size: 12px;
        color: var(--art-gray-500);
      }
    }

    .summary-head {
      display: flex;
      align-items: center;
      justify-content: space-between;
    }

    .summary-icon {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 42px;
      height: 42px;
      border-radius: 12px;
      font-size: 22px;
      color: var(--art-gray-800);
      background: var(--art-gray-200);
    }

    .summary-value {
      display: flex;
      align-items: flex-end;
      gap: 6px;
      margin-top: 6px;

      :deep(.text-g-900) {
        font-size: 27px;
        font-weight: 650;
        color: var(--art-gray-900);
      }

      span {
        margin-bottom: 5px;
        font-size: 13px;
        color: var(--art-gray-500);
      }
    }

    .panel-card {
      min-height: 350px;
      padding: 20px;
    }

    .panel-header {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 12px;
      margin-bottom: 18px;

      h3 {
        margin: 0;
        font-size: 18px;
        font-weight: 600;
        color: var(--art-gray-900);
      }

      p {
        margin-top: 6px;
        font-size: 13px;
        color: var(--art-gray-500);
      }
    }

    .server-card {
      padding: 16px;
      margin-bottom: 14px;
      background: var(--art-gray-100);
      border: 1px solid var(--default-border);
      border-radius: calc(var(--custom-radius) / 2 + 2px);
    }

    .server-name {
      overflow: hidden;
      font-size: 18px;
      font-weight: 650;
      text-overflow: ellipsis;
      white-space: nowrap;
      color: var(--art-gray-900);
    }

    .server-desc {
      margin-top: 8px;
      font-size: 13px;
      line-height: 1.6;
      color: var(--art-gray-600);
    }

    .detail-list {
      display: grid;
      gap: 12px;

      div {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 12px;
        padding-bottom: 12px;
        border-bottom: 1px solid var(--default-border);

        &:last-child {
          padding-bottom: 0;
          border-bottom: 0;
        }
      }

      span {
        flex: 0 0 auto;
        font-size: 13px;
        color: var(--art-gray-500);
      }

      strong {
        min-width: 0;
        overflow-wrap: anywhere;
        font-size: 13px;
        font-weight: 500;
        text-align: right;
        color: var(--art-gray-900);
      }
    }

    .progress-list {
      display: grid;
      gap: 18px;

      p {
        margin-top: 7px;
        font-size: 12px;
        color: var(--art-gray-500);
      }
    }

    .progress-title {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
      margin-bottom: 8px;

      span {
        font-size: 13px;
        color: var(--art-gray-600);
      }

      strong {
        font-size: 13px;
        font-weight: 600;
        color: var(--art-gray-900);
      }
    }

    .detail-grid {
      display: grid;
      grid-template-columns: repeat(3, minmax(0, 1fr));
      gap: 12px;

      @media (max-width: 640px) {
        grid-template-columns: 1fr;
      }
    }

    .detail-block {
      padding: 14px;
      background: var(--art-gray-100);
      border: 1px solid var(--default-border);
      border-radius: calc(var(--custom-radius) / 2 + 2px);

      span {
        display: block;
        font-size: 12px;
        color: var(--art-gray-500);
      }

      strong {
        display: block;
        margin-top: 9px;
        overflow-wrap: anywhere;
        font-size: 18px;
        font-weight: 650;
        color: var(--art-gray-900);
      }
    }
  }
</style>
