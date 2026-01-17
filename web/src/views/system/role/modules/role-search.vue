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

  const searchBarRef = ref()

  /**
   * 表单数据双向绑定
   */
  const formData = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
  })

  /**
   * 表单校验规则
   */
  const rules = {}

  /**
   * 角色状态选项
   */
  const statusOptions = ref([
    { label: '启用', value: true },
    { label: '禁用', value: false }
  ])

  /**
   * 搜索表单配置项
   */
  const formItems = computed(() => [
    {
      label: '角色名称',
      key: 'name',
      type: 'input',
      placeholder: '请输入角色名称',
      clearable: true
    },
    {
      label: '角色代码',
      key: 'code',
      type: 'input',
      placeholder: '请输入角色代码',
      clearable: true
    },
    {
      label: '描述',
      key: 'description',
      type: 'input',
      placeholder: '请输入描述',
      clearable: true
    },
    {
      label: '状态',
      key: 'enabled',
      type: 'select',
      props: {
        placeholder: '请选择状态',
        options: statusOptions.value,
        clearable: true
      }
    }
  ])

  /**
   * 处理重置事件
   */
  const handleReset = () => {
    emit('reset')
  }

  /**
   * 处理搜索事件
   * 验证表单后触发搜索
   */
  const handleSearch = async () => {
    await searchBarRef.value.validate()
    emit('search', formData.value)
  }
</script>
