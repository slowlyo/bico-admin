<template>
  <div class="avatar-upload">
    <ElUpload
      class="avatar-uploader"
      :action="uploadUrl"
      :headers="uploadHeaders"
      :show-file-list="false"
      :before-upload="beforeAvatarUpload"
      :on-success="handleAvatarSuccess"
      :on-error="handleAvatarError"
      :disabled="disabled"
      name="files"
    >
      <div class="avatar-container">
        <img v-if="avatarUrl" :src="avatarUrl" class="avatar" />
        <div v-else class="avatar-placeholder">
          <ElIcon class="avatar-uploader-icon"><Plus /></ElIcon>
        </div>
        <!-- 悬浮操作层 -->
        <div v-if="avatarUrl && !disabled" class="avatar-overlay">
          <ElIcon class="overlay-icon" @click.stop="handlePreview"><ZoomIn /></ElIcon>
          <ElIcon class="overlay-icon" @click.stop="handleDelete"><Delete /></ElIcon>
        </div>
      </div>
    </ElUpload>
    
    <!-- 提示信息 -->
    <div v-if="showTips" class="avatar-tips">
      <p>{{ formatText || '支持 JPG、PNG 格式，文件大小不超过 2MB' }}</p>
      <p v-if="tipText" class="custom-tip">{{ tipText }}</p>
    </div>

    <!-- 图片预览对话框 -->
    <ElDialog v-model="previewVisible" title="头像预览" width="400px" align-center>
      <div class="preview-container">
        <img :src="avatarUrl" class="preview-image" />
      </div>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { computed, ref } from 'vue'
  import { ElUpload, ElIcon, ElMessage, ElDialog, ElMessageBox } from 'element-plus'
  import { Plus, ZoomIn, Delete } from '@element-plus/icons-vue'
  import { useUserStore } from '@/store/modules/user'

  // Props 定义
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
    uploadDir: 'avatars',
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

  // Store
  const userStore = useUserStore()

  // 响应式数据
  const previewVisible = ref(false)

  // 计算属性
  const uploadUrl = computed(() => `${import.meta.env.VITE_API_URL}/admin-api/upload?dir=${props.uploadDir}`)
  
  const uploadHeaders = computed(() => ({
    Authorization: `Bearer ${userStore.accessToken}`
  }))

  const avatarUrl = computed(() => {
    if (!props.modelValue) return ''
    // 如果已经是完整URL，直接返回
    if (props.modelValue.startsWith('http')) return props.modelValue
    // 如果是以 /uploads 开头的路径，拼接当前域名
    if (props.modelValue.startsWith('/uploads/')) {
      return `${window.location.origin}${props.modelValue}`
    }
    // 兼容旧的相对路径格式
    return `${window.location.origin}/uploads/${props.modelValue}`
  })



  // 头像上传前的检查
  const beforeAvatarUpload = (file: File): boolean => {
    const isValidType = props.allowedTypes.includes(file.type)
    const isLtMaxSize = file.size / 1024 / 1024 < props.maxSize

    if (!isValidType) {
      const typeText = props.allowedTypes.map(type => type.split('/')[1].toUpperCase()).join('、')
      ElMessage.error(`头像只能是 ${typeText} 格式!`)
      return false
    }
    if (!isLtMaxSize) {
      ElMessage.error(`头像大小不能超过 ${props.maxSize}MB!`)
      return false
    }
    return true
  }

  // 头像上传成功
  const handleAvatarSuccess = (response: any) => {
    console.log('头像上传响应:', response)
    if (response && response.data && response.data.files && response.data.files.length > 0) {
      const fileData = response.data.files[0]
      // 优先使用完整URL，如果没有则使用相对路径
      const fileValue = fileData.file_url || fileData.file_path
      emit('update:modelValue', fileValue)
      emit('change', fileValue)
      emit('success', response)
      ElMessage.success('头像上传成功')
    } else {
      ElMessage.error('头像上传失败')
      emit('error', response)
    }
  }

  // 头像上传失败
  const handleAvatarError = (error: any) => {
    console.error('头像上传失败:', error)
    ElMessage.error('头像上传失败，请重试')
    emit('error', error)
  }

  // 预览头像
  const handlePreview = () => {
    if (avatarUrl.value) {
      previewVisible.value = true
    }
  }

  // 删除头像
  const handleDelete = async () => {
    try {
      await ElMessageBox.confirm('确定要删除当前头像吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      
      emit('update:modelValue', '')
      emit('change', '')
      ElMessage.success('头像已删除')
    } catch {
      // 用户取消删除
    }
  }
</script>

<style lang="scss" scoped>
  .avatar-upload {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;

    .avatar-uploader {
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

      .avatar-container {
        position: relative;
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;

        .avatar {
          width: 100%;
          height: 100%;
          display: block;
          object-fit: cover;
          border-radius: 6px;
        }

        .avatar-placeholder {
          width: 100%;
          height: 100%;
          display: flex;
          align-items: center;
          justify-content: center;

          .avatar-uploader-icon {
            font-size: 28px;
            color: #8c939d;
          }
        }

        .avatar-overlay {
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

        &:hover .avatar-overlay {
          opacity: 1;
        }
      }
    }

    .avatar-tips {
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

  .preview-container {
    text-align: center;

    .preview-image {
      max-width: 100%;
      max-height: 400px;
      border-radius: 6px;
    }
  }
</style>
