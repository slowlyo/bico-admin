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
  import { RoleService } from '@/api/adminUserApi'

  interface Emits {
    (e: 'search', params: any): void
    (e: 'reset'): void
  }

  const emit = defineEmits<Emits>()

  // 角色选项
  const roleOptions = ref<Api.Role.RoleOption[]>([])

  // 计算角色选项
  const roleSelectOptions = computed(() =>
    roleOptions.value.map(role => ({
      label: role.name,
      value: role.id
    }))
  )

  // 搜索表单
  const searchForm = ref({
    username: '',
    name: '',
    status: undefined,
    role_id: undefined
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

  // 加载角色选项
  const loadRoleOptions = async () => {
    try {
      const response = await RoleService.getActiveRoles()
      roleOptions.value = response || []
    } catch (error) {
      console.error('获取角色选项失败:', error)
    }
  }

  // 搜索配置项
  const searchItems: SearchFormItem[] = [
    {
      label: '用户名',
      prop: 'username',
      type: 'input',
      config: {
        clearable: true,
        placeholder: '请输入用户名'
      }
    },
    {
      label: '姓名',
      prop: 'name',
      type: 'input',
      config: {
        clearable: true,
        placeholder: '请输入姓名'
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
    },
    {
      label: '角色',
      prop: 'role_id',
      type: 'select',
      config: {
        clearable: true,
        placeholder: '请选择角色'
      },
      options: () => roleSelectOptions.value
    }
  ]

  // 组件挂载时加载角色选项
  onMounted(() => {
    loadRoleOptions()
  })
</script>
