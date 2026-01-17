<template>
  <ElConfigProvider size="default" :locale="locales[language]" :z-index="3000">
    <RouterView></RouterView>
  </ElConfigProvider>
</template>

<script setup lang="ts">
  import { useUserStore } from './store/modules/user'
  import { useSettingStore } from './store/modules/setting'
  import { fetchAppConfig } from './api/auth'
  import zh from 'element-plus/es/locale/lang/zh-cn'
  import en from 'element-plus/es/locale/lang/en'
  import { systemUpgrade } from './utils/sys'
  import { toggleTransition } from './utils/ui/animation'
  import { checkStorageCompatibility } from './utils/storage'
  import { initializeTheme } from './hooks/core/useTheme'

  const userStore = useUserStore()
  const settingStore = useSettingStore()
  const { language } = storeToRefs(userStore)

  const locales = {
    zh: zh,
    en: en
  }

  // 初始化应用配置
  const initAppConfig = async () => {
    try {
      const config = await fetchAppConfig()
      if (config) {
        settingStore.setAppConfig(config)
      }
    } catch (error) {
      console.error('初始化应用配置失败:', error)
    }
  }

  onBeforeMount(() => {
    toggleTransition(true)
    initializeTheme()
    initAppConfig()
  })

  onMounted(() => {
    checkStorageCompatibility()
    toggleTransition(false)
    systemUpgrade()
  })
</script>
