<template>
  <div class="file-upload">
    <ElUpload
      class="file-uploader"
      :action="uploadUrl"
      :headers="uploadHeaders"
      :show-file-list="showFileList"
      :file-list="fileList"
      :before-upload="beforeFileUpload"
      :on-success="handleFileSuccess"
      :on-error="handleFileError"
      :on-remove="handleFileRemove"
      :disabled="disabled"
      :multiple="multiple"
      :limit="limit"
      :drag="drag"
      name="files"
    >
      <div v-if="drag" class="drag-upload-area">
        <ElIcon class="upload-icon"><UploadFilled /></ElIcon>
        <div class="upload-text">
          <p class="primary-text">{{ dragText || '将文件拖到此处，或点击上传' }}</p>
          <p class="secondary-text">{{ formatText || getDefaultFormatText() }}</p>
        </div>
      </div>
      <div v-else class="button-upload-area">
        <ElButton type="primary" :disabled="disabled">
          <ElIcon><Upload /></ElIcon>
          {{ buttonText || '选择文件' }}
        </ElButton>
      </div>
    </ElUpload>
    
    <!-- 提示信息 -->
    <div v-if="showTips && !drag" class="file-tips">
      <p>{{ formatText || getDefaultFormatText() }}</p>
      <p v-if="tipText" class="custom-tip">{{ tipText }}</p>
    </div>

    <!-- 自定义文件列表 -->
    <div v-if="!showFileList && fileList.length > 0" class="custom-file-list">
      <div
        v-for="(file, index) in fileList"
        :key="file.uid || index"
        class="file-item"
      >
        <div class="file-info">
          <ElIcon class="file-icon">
            <Document v-if="isDocumentFile(file.name)" />
            <Picture v-else-if="isImageFile(file.name)" />
            <VideoPlay v-else-if="isVideoFile(file.name)" />
            <Headset v-else-if="isAudioFile(file.name)" />
            <Files v-else />
          </ElIcon>
          <div class="file-details">
            <span class="file-name" :title="file.name">{{ file.name }}</span>
            <span class="file-size">{{ formatFileSize(file.size) }}</span>
          </div>
        </div>
        <div class="file-actions">
          <ElButton
            v-if="file.url"
            type="primary"
            link
            size="small"
            @click="handleFilePreview(file)"
          >
            <ElIcon><View /></ElIcon>
          </ElButton>
          <ElButton
            type="danger"
            link
            size="small"
            @click="handleFileRemove(file)"
          >
            <ElIcon><Delete /></ElIcon>
          </ElButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { computed, ref, watch } from 'vue'
  import { ElUpload, ElIcon, ElMessage, ElButton } from 'element-plus'
  import { 
    UploadFilled, 
    Upload, 
    Document, 
    Picture, 
    VideoPlay, 
    Headset, 
    Files,
    View,
    Delete
  } from '@element-plus/icons-vue'
  import { useUserStore } from '@/store/modules/user'
  import type { UploadFile, UploadFiles } from 'element-plus'

  // Props 定义
  interface Props {
    /** 文件列表值 */
    modelValue?: UploadFile[]
    /** 是否禁用 */
    disabled?: boolean
    /** 是否支持拖拽上传 */
    drag?: boolean
    /** 是否支持多选 */
    multiple?: boolean
    /** 最大上传数量 */
    limit?: number
    /** 是否显示文件列表 */
    showFileList?: boolean
    /** 是否显示提示信息 */
    showTips?: boolean
    /** 自定义提示文本（可选） */
    tipText?: string
    /** 格式提示文本 */
    formatText?: string
    /** 拖拽区域文本 */
    dragText?: string
    /** 按钮文本 */
    buttonText?: string
    /** 上传目录 */
    uploadDir?: string
    /** 最大文件大小(MB) */
    maxSize?: number
    /** 允许的文件类型 */
    allowedTypes?: string[]
  }

  const props = withDefaults(defineProps<Props>(), {
    modelValue: () => [],
    disabled: false,
    drag: true,
    multiple: true,
    limit: 10,
    showFileList: false,
    showTips: true,
    tipText: '',
    formatText: '',
    dragText: '',
    buttonText: '',
    uploadDir: 'files',
    maxSize: 10,
    allowedTypes: () => []
  })

  // Emits 定义
  interface Emits {
    (e: 'update:modelValue', value: UploadFile[]): void
    (e: 'change', value: UploadFile[]): void
    (e: 'success', response: any, file: UploadFile): void
    (e: 'error', error: any, file: UploadFile): void
    (e: 'remove', file: UploadFile): void
    (e: 'preview', file: UploadFile): void
  }

  const emit = defineEmits<Emits>()

  // Store
  const userStore = useUserStore()

  // 响应式数据
  const fileList = ref<UploadFile[]>([...props.modelValue])

  // 监听 modelValue 变化
  watch(() => props.modelValue, (newVal) => {
    fileList.value = [...newVal]
  }, { deep: true })

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

  // 获取默认格式提示文本
  const getDefaultFormatText = () => {
    if (props.allowedTypes.length > 0) {
      const types = props.allowedTypes.map(type => {
        if (type.startsWith('.')) return type.toUpperCase()
        return type.split('/')[1]?.toUpperCase() || type
      }).join('、')
      return `支持 ${types} 格式，单个文件不超过 ${props.maxSize}MB`
    }
    return `单个文件不超过 ${props.maxSize}MB`
  }

  // 文件上传前的检查
  const beforeFileUpload = (file: File): boolean => {
    // 检查文件类型
    if (props.allowedTypes.length > 0) {
      const isValidType = props.allowedTypes.some(type => {
        if (type.startsWith('.')) {
          return file.name.toLowerCase().endsWith(type.toLowerCase())
        }
        return file.type === type
      })
      
      if (!isValidType) {
        const typeText = props.allowedTypes.join('、')
        ElMessage.error(`文件格式不支持，请选择 ${typeText} 格式的文件`)
        return false
      }
    }

    // 检查文件大小
    const isLtMaxSize = file.size / 1024 / 1024 < props.maxSize
    if (!isLtMaxSize) {
      ElMessage.error(`文件大小不能超过 ${props.maxSize}MB`)
      return false
    }

    return true
  }

  // 文件上传成功
  const handleFileSuccess = (response: any, file: UploadFile, files: UploadFiles) => {
    console.log('文件上传响应:', response)
    if (response && response.data && response.data.files && response.data.files.length > 0) {
      const fileData = response.data.files[0]
      // 更新文件信息
      file.url = fileData.file_url || fileData.file_path
      file.response = response
      
      fileList.value = [...files]
      emit('update:modelValue', fileList.value)
      emit('change', fileList.value)
      emit('success', response, file)
      ElMessage.success('文件上传成功')
    } else {
      ElMessage.error('文件上传失败')
      emit('error', response, file)
    }
  }

  // 文件上传失败
  const handleFileError = (error: any, file: UploadFile, files: UploadFiles) => {
    console.error('文件上传失败:', error)
    ElMessage.error('文件上传失败，请重试')
    emit('error', error, file)
  }

  // 移除文件
  const handleFileRemove = (file: UploadFile) => {
    const index = fileList.value.findIndex(item => item.uid === file.uid)
    if (index > -1) {
      fileList.value.splice(index, 1)
      emit('update:modelValue', fileList.value)
      emit('change', fileList.value)
      emit('remove', file)
    }
  }

  // 预览文件
  const handleFilePreview = (file: UploadFile) => {
    if (file.url) {
      window.open(file.url, '_blank')
    }
    emit('preview', file)
  }

  // 文件类型判断
  const isImageFile = (fileName: string) => {
    const imageExts = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp', '.svg']
    return imageExts.some(ext => fileName.toLowerCase().endsWith(ext))
  }

  const isDocumentFile = (fileName: string) => {
    const docExts = ['.pdf', '.doc', '.docx', '.xls', '.xlsx', '.ppt', '.pptx', '.txt']
    return docExts.some(ext => fileName.toLowerCase().endsWith(ext))
  }

  const isVideoFile = (fileName: string) => {
    const videoExts = ['.mp4', '.avi', '.mov', '.wmv', '.flv', '.mkv']
    return videoExts.some(ext => fileName.toLowerCase().endsWith(ext))
  }

  const isAudioFile = (fileName: string) => {
    const audioExts = ['.mp3', '.wav', '.flac', '.aac', '.ogg']
    return audioExts.some(ext => fileName.toLowerCase().endsWith(ext))
  }

  // 格式化文件大小
  const formatFileSize = (size?: number) => {
    if (!size) return '0 B'
    const units = ['B', 'KB', 'MB', 'GB']
    let index = 0
    let fileSize = size
    
    while (fileSize >= 1024 && index < units.length - 1) {
      fileSize /= 1024
      index++
    }
    
    return `${fileSize.toFixed(1)} ${units[index]}`
  }
</script>

<style lang="scss" scoped>
  .file-upload {
    .file-uploader {
      width: 100%;

      // 拖拽上传区域样式
      :deep(.el-upload-dragger) {
        width: 100%;
        height: auto;
        min-height: 180px;
        border: 2px dashed var(--el-border-color);
        border-radius: 8px;
        background: var(--el-fill-color-blank);
        transition: all 0.3s;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 40px 20px;

        &:hover {
          border-color: var(--el-color-primary);
          background: var(--el-color-primary-light-9);
        }

        &.is-dragover {
          border-color: var(--el-color-primary);
          background: var(--el-color-primary-light-8);
        }
      }

      // 按钮上传区域样式
      :deep(.el-upload) {
        display: block;
        width: 100%;
      }

      .drag-upload-area {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        width: 100%;

        .upload-icon {
          font-size: 48px;
          color: var(--el-color-info);
          margin-bottom: 16px;
        }

        .upload-text {
          text-align: center;

          .primary-text {
            font-size: 16px;
            color: var(--el-text-color-primary);
            margin: 0 0 8px 0;
            font-weight: 500;
          }

          .secondary-text {
            font-size: 14px;
            color: var(--el-text-color-regular);
            margin: 0;
          }
        }
      }

      .button-upload-area {
        text-align: center;
      }
    }

    .file-tips {
      margin-top: 12px;
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

    .custom-file-list {
      margin-top: 16px;
      border: 1px solid var(--el-border-color);
      border-radius: 6px;
      overflow: hidden;

      .file-item {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 12px 16px;
        border-bottom: 1px solid var(--el-border-color-lighter);
        transition: background-color 0.3s;

        &:last-child {
          border-bottom: none;
        }

        &:hover {
          background: var(--el-fill-color-light);
        }

        .file-info {
          display: flex;
          align-items: center;
          flex: 1;
          min-width: 0;

          .file-icon {
            font-size: 20px;
            color: var(--el-color-primary);
            margin-right: 12px;
            flex-shrink: 0;
          }

          .file-details {
            display: flex;
            flex-direction: column;
            min-width: 0;
            flex: 1;

            .file-name {
              font-size: 14px;
              color: var(--el-text-color-primary);
              font-weight: 500;
              white-space: nowrap;
              overflow: hidden;
              text-overflow: ellipsis;
              margin-bottom: 2px;
            }

            .file-size {
              font-size: 12px;
              color: var(--el-text-color-regular);
            }
          }
        }

        .file-actions {
          display: flex;
          align-items: center;
          gap: 8px;
          flex-shrink: 0;
        }
      }
    }
  }
</style>
