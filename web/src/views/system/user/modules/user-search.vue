<template>
  <ArtSearchBar
    ref="searchBarRef"
    v-model="formData"
    :items="formItems"
    :rules="rules"
    @reset="handleReset"
    @search="handleSearch"
  >
  </ArtSearchBar>
</template>

<script setup lang="ts">
  interface Props {
    modelValue: Record<string, any>
  }
  interface Emits {
    (e: 'update:modelValue', value: Record<string, any>): void
    (e: 'search', params: Record<string, any>): void
    (e: 'reset'): void
  }
  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  // 表单数据双向绑定
  const searchBarRef = ref()
  const formData = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
  })

  // 校验规则
  const rules = {
    // userName: [{ required: true, message: '请输入用户名', trigger: 'blur' }]
  }

  // 动态 options
  const statusOptions = ref<{ label: string; value: string; disabled?: boolean }[]>([])

  // 模拟接口返回状态数据
  function fetchStatusOptions(): Promise<typeof statusOptions.value> {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve([
          { label: '在线', value: '1' },
          { label: '离线', value: '2' },
          { label: '异常', value: '3' },
          { label: '注销', value: '4' }
        ])
      }, 1000)
    })
  }

  onMounted(async () => {
    statusOptions.value = await fetchStatusOptions()
  })

  // 表单配置
  const formItems = computed(() => [
    {
      label: '用户名',
      key: 'username',
      type: 'input',
      placeholder: '请输入用户名',
      clearable: true
    },
    {
      label: '姓名',
      key: 'name',
      type: 'input',
      placeholder: '请输入姓名',
      clearable: true
    },
    {
      label: '状态',
      key: 'enabled',
      type: 'select',
      props: {
        placeholder: '请选择状态',
        options: [
          { label: '启用', value: true },
          { label: '禁用', value: false }
        ]
      }
    }
  ])

  // 事件
  function handleReset() {
    console.log('重置表单')
    emit('reset')
  }

  async function handleSearch() {
    await searchBarRef.value.validate()
    emit('search', formData.value)
    console.log('表单数据', formData.value)
  }
</script>
