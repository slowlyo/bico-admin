<template>
  <ElConfigProvider size="default" :locale="locales[language]" :z-index="3000">
    <RouterView></RouterView>
  </ElConfigProvider>
</template>

<script setup lang="ts">
  import { useUserStore } from './store/modules/user'
  import zh from 'element-plus/es/locale/lang/zh-cn'
  import en from 'element-plus/es/locale/lang/en'
  import { systemUpgrade } from './utils/sys'
  import { UserService } from './api/usersApi'
  import { setThemeTransitionClass } from './utils/theme/animation'
  import { checkStorageCompatibility } from './utils/storage'
  import { startTokenRefreshTimer } from './utils/tokenManager'

  const userStore = useUserStore()
  const { language } = storeToRefs(userStore)

  // token刷新定时器清理函数
  let stopTokenRefreshTimer: (() => void) | null = null

  const locales = {
    zh: zh,
    en: en
  }

  onBeforeMount(() => {
    setThemeTransitionClass(true)
  })

  onMounted(() => {
    // 检查存储兼容性
    checkStorageCompatibility()
    // 提升暗黑主题下页面刷新视觉体验
    setThemeTransitionClass(false)
    // 系统升级
    systemUpgrade()
    // 获取用户信息
    getUserInfo()
    // 启动token自动刷新定时器
    if (userStore.isLogin) {
      stopTokenRefreshTimer = startTokenRefreshTimer()
    }
  })

  onUnmounted(() => {
    // 清理token刷新定时器
    if (stopTokenRefreshTimer) {
      stopTokenRefreshTimer()
    }
  })

  // 监听登录状态变化，自动获取用户信息和管理token刷新定时器
  watch(() => userStore.isLogin, (newValue) => {
    if (newValue) {
      getUserInfo()
      // 启动token自动刷新定时器
      if (!stopTokenRefreshTimer) {
        stopTokenRefreshTimer = startTokenRefreshTimer()
      }
    } else {
      // 停止token刷新定时器
      if (stopTokenRefreshTimer) {
        stopTokenRefreshTimer()
        stopTokenRefreshTimer = null
      }
    }
  })

  // 获取用户信息
  const getUserInfo = async () => {
    if (userStore.isLogin) {
      try {
        const data = await UserService.getUserInfo()
        userStore.setUserInfo(data.user_info)
        userStore.setPermissions(data.permissions)
      } catch (error) {
        console.error('获取用户信息失败', error)
      }
    }
  }
</script>
