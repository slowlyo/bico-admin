# 前端 CRUD 开发指南

> 基于 Vue 3 + Artd Pro，通过 `useTable` Hook 和 `ArtTable` 组件实现高效的 CRUD 开发。

## 核心 Hook: `useTable`

`useTable` 是整个框架最核心的表格管理 Hook，封装了数据获取、分页、搜索、缓存、刷新策略以及列配置。

### 最小示例

```vue
<template>
  <div class="art-full-height">
    <!-- 搜索栏 -->
    <RoleSearch v-model="searchForm" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card" shadow="never">
      <!-- 表格工具栏 -->
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton type="primary" @click="showDialog('add')">新增角色</ElButton>
        </template>
      </ArtTableHeader>

      <!-- 数据表格 -->
      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
import { useTable } from '@/hooks/core/useTable'
import { fetchGetRoleList } from '@/api/system-manage'

const {
  columns,
  columnChecks,
  data,
  loading,
  pagination,
  searchParams,
  handleSearch,
  resetSearchParams,
  handleSizeChange,
  handleCurrentChange,
  refreshData
} = useTable({
  core: {
    apiFn: fetchGetRoleList,
    columnsFactory: () => [
      { prop: 'id', label: 'ID', width: 80 },
      { prop: 'name', label: '角色名称', minWidth: 120 },
      { prop: 'code', label: '角色代码', minWidth: 120 },
      { prop: 'created_at', label: '创建时间', width: 180 }
    ]
  }
})
</script>
```

---

## 刷新策略

`useTable` 提供了 5 种刷新方法，适配不同业务场景，能智能处理缓存清理和页码跳转：

1. **`refreshCreate`**: 新增后刷新。回到第一页并清空分页相关缓存。
2. **`refreshUpdate`**: 更新后刷新。保持当前页，仅清空当前查询条件的缓存。
3. **`refreshRemove`**: 删除后刷新。保持当前页并清空缓存，若当前页数据删完则自动跳到前一页。
4. **`refreshData`**: 手动全量刷新。清空所有缓存并重新获取。
5. **`refreshSoft`**: 软刷新。仅清理缓存并获取，不改变任何状态，适用于定时任务。

---

## 权限控制

权限控制通过 `v-auth` 指令或 `useAuth` Hook 实现，后端权限 key 对应 `meta.permissions`。

```vue
<!-- 指令方式 -->
<ElButton v-auth="'system:role:create'">新增</ElButton>

<!-- Hook 方式 -->
<script setup>
const { hasAuth } = useAuth()
if (hasAuth('system:role:edit')) {
  // ...
}
</script>
```

---

## 搜索处理

搜索参数统一维护在 `searchParams` 中。

```ts
const handleSearch = (params: Record<string, any>) => {
  Object.assign(searchParams, params)
  getData() // useTable 内部方法，重置到第一页并搜索
}
```

---

## 最佳实践

1. **Columns 抽离**: 建议将 `columnsFactory` 抽离到单独的文件或在 `setup` 顶部定义，保持逻辑清晰。
2. **Api 函数规范**: 统一在 `src/api/` 下定义，使用 `request` 工具，并定义好 `Api` 命名空间下的类型。
3. **缓存使用**: 对于变动不频繁的数据（如配置列表），开启 `performance.enableCache` 可显著提升用户体验。
4. **响应式参数**: `searchParams` 是响应式的，可以直接绑定到搜索表单，但建议通过事件触发搜索以保证性能（防抖）。
