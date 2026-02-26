# 前端 CRUD 开发指南

基于 Vue 3 + `useTable` + `ArtTable` 的标准开发方式。

## 最小可用示例

```vue
<template>
  <div class="art-full-height">
    <RoleSearch v-model="searchForm" @search="handleSearch" @reset="resetSearchParams" />

    <ElCard class="art-table-card" shadow="never">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData" />

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

const searchForm = ref({ name: undefined, code: undefined })

const {
  columns,
  columnChecks,
  data,
  loading,
  pagination,
  searchParams,
  getData,
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
      { prop: 'code', label: '角色代码', minWidth: 120 }
    ]
  }
})

const handleSearch = (params: Record<string, any>) => {
  Object.assign(searchParams, params)
  getData()
}
</script>
```

## 刷新策略

`useTable` 内置以下刷新函数：

1. `refreshCreate`：新增后刷新，回到第一页
2. `refreshUpdate`：更新后刷新，保持当前页
3. `refreshRemove`：删除后刷新，必要时回退页码
4. `refreshData`：全量刷新
5. `refreshSoft`：软刷新（轻量）

## 权限控制

可使用 `v-auth` 指令或 `useAuth()`：

```vue
<ElButton v-auth="'system:admin_role:create'">新增</ElButton>
```

## 实践建议

1. API 调用统一放在 `web/src/api/`。
2. 表格页统一复用 `useTable`，避免重复分页与请求逻辑。
3. 搜索时通过 `searchParams + getData()` 驱动，不要直接改 `data`。
4. 涉及权限按钮时，同时做前端显隐和后端权限校验。
