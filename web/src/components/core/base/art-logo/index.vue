<!-- 系统logo -->
<template>
  <div class="flex-cc">
    <img
      :style="logoStyle"
      :src="logoSrc"
      alt="logo"
      class="w-full h-full object-contain"
    />
  </div>
</template>

<script setup lang="ts">
  import { useSettingStore } from '@/store/modules/setting'
  import defaultLogo from '@imgs/common/logo.webp'

  defineOptions({ name: 'ArtLogo' })

  interface Props {
    /** logo 大小 */
    size?: number | string
    /** 自定义 logo 地址 */
    src?: string
  }

  const props = withDefaults(defineProps<Props>(), {
    size: 36
  })

  const settingStore = useSettingStore()

  const logoSrc = computed(() => {
    if (props.src) return props.src
    return settingStore.appConfig.logo || defaultLogo
  })

  const logoStyle = computed(() => ({ width: `${props.size}px` }))
</script>
