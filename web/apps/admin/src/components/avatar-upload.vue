<script lang="ts" setup>
import { computed, ref } from 'vue';

import { Camera, Eye, LoaderCircle, User } from '@vben/icons';

import { Avatar, Button, Image, Tooltip, Upload, message } from 'antdv-next';

const props = withDefaults(
  defineProps<{
    disabled?: boolean;
    hint?: false | string;
    label?: string;
    modelValue?: string;
    /** 选中文件后上传，返回要展示的地址。写不写用户资料由调用方决定。 */
    onUpload: (file: File) => Promise<string>;
    size?: number;
    successMessage?: string;
  }>(),
  {
    disabled: false,
    hint: false,
    modelValue: '',
    size: 80,
    successMessage: '头像上传成功',
  },
);

const emit = defineEmits<{
  'update:modelValue': [url: string];
}>();

const previewOpen = ref(false);
const uploading = ref(false);

// 头像比外圈小 2px，边框才露得出来。
const avatarSize = computed(() => Math.max(props.size - 2, 1));
// 底部遮罩按头像高度取一小条，太矮时字会挤出圆边。
const maskHeight = computed(() => Math.max(Math.round(props.size * 0.28), 28));

// 把 antd 的上传回调收成地址。失败时不改当前头像。
async function handleUpload(options: {
  file: Blob | File | string;
  onError?: (error: Error) => void;
  onSuccess?: (body: unknown) => void;
}) {
  const file = options.file;
  // customRequest 的 file 在选图时是 File。其它类型不是这次选择的图片。
  if (!(file instanceof File)) {
    const uploadError = new Error('头像上传失败');
    message.error(uploadError.message);
    options.onError?.(uploadError);
    return;
  }

  uploading.value = true;
  try {
    const url = await props.onUpload(file);
    // 没有地址就不能写进表单，否则会把空头像存下去。
    if (!url) {
      throw new Error('上传结果缺少图片地址');
    }
    emit('update:modelValue', url);
    message.success(props.successMessage);
    options.onSuccess?.({ url });
  } catch (error) {
    const uploadError =
      error instanceof Error ? error : new Error('头像上传失败');
    const notified = error as { data?: unknown; response?: unknown };
    // 接口失败时请求层已经提示过后台 msg，这里只补组件自己抛出的错误。
    if (!notified.response && !notified.data) {
      message.error(uploadError.message || '头像上传失败');
    }
    options.onError?.(uploadError);
  } finally {
    uploading.value = false;
  }
}

// 预览关闭时要写回开关，否则下一次点眼睛不会再打开。
function onPreviewOpenChange(open: boolean) {
  previewOpen.value = open;
}
</script>

<template>
  <div class="avatar-upload">
    <div v-if="label" class="avatar-upload__label">{{ label }}</div>
    <div class="avatar-upload__stack">
      <div class="avatar-upload__wrapper">
        <Upload
          accept="image/*"
          :custom-request="handleUpload"
          :disabled="disabled || uploading"
          :show-upload-list="false"
        >
          <button
            aria-label="更换头像"
            class="avatar-upload__trigger"
            :disabled="disabled || uploading"
            :style="{ width: `${size}px`, height: `${size}px` }"
            type="button"
          >
            <Avatar :size="avatarSize" :src="modelValue || undefined" alt="当前头像">
              <template v-if="!modelValue" #icon>
                <User class="size-5" />
              </template>
            </Avatar>
            <span class="avatar-upload__mask" :style="{ height: `${maskHeight}px` }">
              <LoaderCircle v-if="uploading" class="size-3.5 animate-spin" />
              <Camera v-else class="size-3.5" />
              <span>{{ uploading ? '上传中' : '更换' }}</span>
            </span>
          </button>
        </Upload>
        <Tooltip title="预览头像">
          <Button
            aria-label="预览头像"
            class="avatar-upload__preview"
            :disabled="!modelValue"
            shape="circle"
            @click="previewOpen = true"
          >
            <template #icon>
              <Eye class="size-4" />
            </template>
          </Button>
        </Tooltip>
        <Image
          v-if="modelValue"
          alt="头像预览"
          :preview="{ open: previewOpen, onOpenChange: onPreviewOpenChange }"
          :src="modelValue"
        />
      </div>
      <span v-if="hint" class="avatar-upload__hint">{{ hint }}</span>
    </div>
  </div>
</template>

<style scoped>
.avatar-upload__label {
  margin-bottom: 8px;
  font-weight: 500;
}

.avatar-upload__stack {
  display: inline-flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
}

.avatar-upload__wrapper {
  position: relative;
  display: inline-block;
  /* 预览按钮要压在圆角外。用户表单格子是 overflow:hidden，负偏移会被切掉。 */
  padding: 8px;
}

.avatar-upload__wrapper :deep(.ant-upload) {
  display: block;
  line-height: 0;
}

.avatar-upload__trigger {
  position: relative;
  padding: 0;
  overflow: hidden;
  cursor: pointer;
  background: var(--ant-color-bg-container);
  border: 1px solid var(--ant-color-border);
  border-radius: 50%;
}

.avatar-upload__trigger:focus-visible {
  outline: 2px solid var(--ant-color-primary);
  outline-offset: 2px;
}

.avatar-upload__trigger:disabled {
  cursor: not-allowed;
}

.avatar-upload__mask {
  position: absolute;
  right: 0;
  bottom: 0;
  left: 0;
  display: flex;
  gap: 4px;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  line-height: 1;
  color: #fff;
  background: rgb(0 0 0 / 62%);
}

.avatar-upload__preview {
  position: absolute;
  top: 0;
  right: 0;
  z-index: 1;
  box-shadow: var(--ant-box-shadow-tertiary);
}

.avatar-upload__wrapper :deep(.ant-image) {
  display: none;
}

.avatar-upload__hint {
  font-size: 12px;
  color: var(--ant-color-text-secondary);
}
</style>
