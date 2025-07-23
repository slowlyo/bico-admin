<template>
  <div class="image-upload">
    <ElUpload
      class="image-uploader"
      :action="uploadUrl"
      :headers="uploadHeaders"
      :show-file-list="false"
      :before-upload="beforeImageUpload"
      :on-success="handleImageSuccess"
      :on-error="handleImageError"
      :disabled="disabled"
      name="files"
    >
      <div class="image-container">
        <img v-if="imageUrl" :src="imageUrl" class="image" />
        <div v-else class="image-placeholder">
          <ElIcon class="image-uploader-icon"><Plus /></ElIcon>
        </div>
        <!-- 悬浮操作层 -->
        <div v-if="imageUrl && !disabled" class="image-overlay">
          <ElIcon class="overlay-icon" @click.stop="handlePreview"><ZoomIn /></ElIcon>
          <ElIcon class="overlay-icon" @click.stop="handleDelete"><Delete /></ElIcon>
        </div>
      </div>
    </ElUpload>
    
    <!-- 提示信息 -->
    <div v-if="showTips" class="image-tips">
      <p>{{ formatText || '支持 JPG、PNG 格式，文件大小不超过 2MB' }}</p>
      <p v-if="tipText" class="custom-tip">{{ tipText }}</p>
    </div>

    <!-- 图片预览 -->
    <ElImageViewer
      v-if="previewVisible"
      :url-list="[imageUrl]"
      @close="previewVisible = false"
    />
  </div>
</template>

<script setup lang="ts">
  import { computed, ref } from 'vue'
  import { ElUpload, ElIcon, ElMessage, ElMessageBox, ElImageViewer } from 'element-plus'
  import { Plus, ZoomIn, Delete } from '@element-plus/icons-vue'
  import { useUserStore } from '@/store/modules/user'

  // Props 定义
  interface Props {
    /** 图片值 */
    modelValue?: string
    /** 是否禁用 */
    disabled?: boolean
    /** 图片尺寸 */
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
    /** 预览对话框标题 */
    previewTitle?: string
  }

  const props = withDefaults(defineProps<Props>(), {
    modelValue: '',
    disabled: false,
    size: 100,
    showTips: true,
    tipText: '',
    formatText: '',
    uploadDir: 'images',
    maxSize: 2,
    allowedTypes: () => ['image/jpeg', 'image/png'],
    previewTitle: '图片预览'
  })

  // Emits 定义
  interface Emits {
    (e: 'update:modelValue', value: string): void
    (e: 'change', value: string): void
    (e: 'success', response: any): void
    (e: 'error', error: any): void
  }

  const emit = defineEmits<Emits>()

  // Store
  const userStore = useUserStore()

  // 响应式数据
  const previewVisible = ref(false)

  // 计算属性
  const uploadUrl = computed(() => {
    const apiUrl = import.meta.env.VITE_API_URL
    // 当 VITE_API_URL 为空字符串或仅为单个斜杠时，使用当前页面域名
    const baseUrl = (!apiUrl || apiUrl === '/') ? window.location.origin : apiUrl
    return `${baseUrl}/admin-api/upload?dir=${props.uploadDir}`
  })
  
  const uploadHeaders = computed(() => ({
    Authorization: `Bearer ${userStore.accessToken}`
  }))

  const imageUrl = computed(() => {
    if (!props.modelValue) return ''
    // 直接使用返回的URL，无需额外处理
    return props.modelValue
  })

  // 图片上传前的检查
  const beforeImageUpload = (file: File): boolean => {
    const isValidType = props.allowedTypes.includes(file.type)
    const isLtMaxSize = file.size / 1024 / 1024 < props.maxSize

    if (!isValidType) {
      const typeText = props.allowedTypes.map(type => type.split('/')[1].toUpperCase()).join('、')
      ElMessage.error(`图片只能是 ${typeText} 格式!`)
      return false
    }
    if (!isLtMaxSize) {
      ElMessage.error(`图片大小不能超过 ${props.maxSize}MB!`)
      return false
    }
    return true
  }

  // 图片上传成功
  const handleImageSuccess = (response: any) => {
    console.log('图片上传响应:', response)
    if (response && response.data && response.data.files && response.data.files.length > 0) {
      const fileData = response.data.files[0]
      // 直接使用返回的URL
      const fileValue = fileData.file_url || fileData.file_path
      emit('update:modelValue', fileValue)
      emit('change', fileValue)
      emit('success', response)
      ElMessage.success('图片上传成功')
    } else {
      ElMessage.error('图片上传失败')
      emit('error', response)
    }
  }

  // 图片上传失败
  const handleImageError = (error: any) => {
    console.error('图片上传失败:', error)
    ElMessage.error('图片上传失败，请重试')
    emit('error', error)
  }

  // 预览图片
  const handlePreview = () => {
    if (imageUrl.value) {
      previewVisible.value = true
    }
  }

  // 删除图片
  const handleDelete = async () => {
    try {
      await ElMessageBox.confirm('确定要删除当前图片吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      
      emit('update:modelValue', '')
      emit('change', '')
      ElMessage.success('图片已删除')
    } catch {
      // 用户取消删除
    }
  }
</script>

<style lang="scss" scoped>
  .image-upload {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;

    .image-uploader {
      :deep(.el-upload) {
        border: 1px dashed var(--el-border-color);
        border-radius: 6px;
        cursor: pointer;
        position: relative;
        overflow: hidden;
        transition: var(--el-transition-duration-fast);
        width: v-bind('props.size + "px"');
        height: v-bind('props.size + "px"');
        display: flex;
        align-items: center;
        justify-content: center;

        &:hover {
          border-color: var(--el-color-primary);
        }
      }

      .image-container {
        position: relative;
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;

        .image {
          width: 100%;
          height: 100%;
          display: block;
          object-fit: cover;
          border-radius: 6px;
        }

        .image-placeholder {
          width: 100%;
          height: 100%;
          display: flex;
          align-items: center;
          justify-content: center;

          .image-uploader-icon {
            font-size: 28px;
            color: #8c939d;
          }
        }

        .image-overlay {
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          background: rgba(0, 0, 0, 0.5);
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 10px;
          opacity: 0;
          transition: opacity 0.3s;
          border-radius: 6px;

          .overlay-icon {
            color: white;
            font-size: 18px;
            cursor: pointer;
            padding: 4px;
            border-radius: 4px;
            transition: background-color 0.3s;

            &:hover {
              background-color: rgba(255, 255, 255, 0.2);
            }
          }
        }

        &:hover .image-overlay {
          opacity: 1;
        }
      }
    }

    .image-tips {
      color: var(--el-text-color-placeholder);
      font-size: 12px;
      line-height: 1.4;

      p {
        margin: 0 0 4px 0;

        &:last-child {
          margin-bottom: 0;
        }

        &.custom-tip {
          color: var(--el-text-color-regular);
          font-size: 13px;
        }
      }
    }
  }


</style>
