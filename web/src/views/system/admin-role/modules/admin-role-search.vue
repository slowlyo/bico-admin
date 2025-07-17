<template>
  <ArtSearchBar
    v-model:filter="searchForm"
    :items="searchItems"
    @reset="handleReset"
    @search="handleSearch"
  />
</template>

<script setup lang="ts">
  import type { SearchFormItem } from '@/types'

  interface Emits {
    (e: 'search', params: any): void
    (e: 'reset'): void
  }

  const emit = defineEmits<Emits>()

  // 搜索表单
  const searchForm = ref({
    name: '',
    code: '',
    status: undefined
  })

  // 重置搜索
  const handleReset = () => {
    emit('reset')
  }

  // 执行搜索
  const handleSearch = () => {
    console.log('搜索表单数据:', searchForm.value)
    emit('search', searchForm.value)
  }

  // 搜索配置项
  const searchItems: SearchFormItem[] = [
    {
      label: '角色名称',
      prop: 'name',
      type: 'input',
      config: {
        clearable: true,
        placeholder: '请输入角色名称'
      }
    },
    {
      label: '角色代码',
      prop: 'code',
      type: 'input',
      config: {
        clearable: true,
        placeholder: '请输入角色代码'
      }
    },
    {
      label: '状态',
      prop: 'status',
      type: 'select',
      config: {
        clearable: true,
        placeholder: '请选择状态'
      },
      options: () => [
        { label: '启用', value: 1 },
        { label: '禁用', value: 0 }
      ]
    }
  ]
</script>
