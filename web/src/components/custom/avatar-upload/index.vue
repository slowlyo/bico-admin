<template>
  <ImageUpload
    :model-value="modelValue"
    :disabled="disabled"
    :size="size"
    :show-tips="showTips"
    :tip-text="tipText"
    :format-text="formatText"
    :upload-dir="uploadDir"
    :max-size="maxSize"
    :allowed-types="allowedTypes"
    preview-title="头像预览"
    @update:model-value="handleUpdateModelValue"
    @change="handleChange"
    @success="handleSuccess"
    @error="handleError"
  />
</template>

<script setup lang="ts">
  import ImageUpload from '../image-upload/index.vue'

  // Props 定义 - 保持与原 AvatarUpload 相同的接口
  interface Props {
    /** 头像值 */
    modelValue?: string
    /** 是否禁用 */
    disabled?: boolean
    /** 头像尺寸 */
    size?: number
    /** 是否显示提示信息 */
    showTips?: boolean
    /** 自定义提示文本（可选） */
    tipText?: string
    /** 格式提示文本 */
    formatText?: string
    /** 上传目录 */
    uploadDir?: string
    /** 最大文件大小(MB) */
    maxSize?: number
    /** 允许的文件类型 */
    allowedTypes?: string[]
  }

  const props = withDefaults(defineProps<Props>(), {
    modelValue: '',
    disabled: false,
    size: 100,
    showTips: true,
    tipText: '',
    formatText: '',
    uploadDir: 'avatars', // 保持原有的默认目录
    maxSize: 2,
    allowedTypes: () => ['image/jpeg', 'image/png']
  })

  // Emits 定义
  interface Emits {
    (e: 'update:modelValue', value: string): void
    (e: 'change', value: string): void
    (e: 'success', response: any): void
    (e: 'error', error: any): void
  }

  const emit = defineEmits<Emits>()

  // 事件处理器
  const handleUpdateModelValue = (value: string) => {
    emit('update:modelValue', value)
  }

  const handleChange = (value: string) => {
    emit('change', value)
  }

  const handleSuccess = (response: any) => {
    emit('success', response)
  }

  const handleError = (error: any) => {
    emit('error', error)
  }
</script>


