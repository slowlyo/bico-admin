<template>
  <div class="login">
    <LoginLeftView></LoginLeftView>

    <div class="right-wrap">
      <div class="top-right-wrap">
        <div class="btn theme-btn" @click="toggleTheme">
          <i class="iconfont-sys">
            {{ isDark ? '&#xe6b5;' : '&#xe725;' }}
          </i>
        </div>
      </div>
      <div class="header">
        <ArtLogo class="icon" />
        <h1>{{ systemName }}</h1>
      </div>
      <div class="login-wrap">
        <div class="form">
          <h3 class="title">{{ $t('login.title') }}</h3>
          <p class="sub-title">{{ $t('login.subTitle') }}</p>
          <ElForm
            ref="formRef"
            :model="formData"
            :rules="rules"
            @keyup.enter="handleSubmit"
            style="margin-top: 25px"
          >

            <ElFormItem prop="username">
              <ElInput :placeholder="$t('login.placeholder[0]')" v-model.trim="formData.username" />
            </ElFormItem>
            <ElFormItem prop="password">
              <ElInput
                :placeholder="$t('login.placeholder[1]')"
                v-model.trim="formData.password"
                type="password"
                radius="8px"
                autocomplete="off"
                show-password
              />
            </ElFormItem>
            <div class="drag-verify">
              <div class="drag-verify-content" :class="{ error: !isPassing && isClickPass }">
                <ArtDragVerify
                  ref="dragVerify"
                  v-model:value="isPassing"
                  :text="$t('login.sliderText')"
                  textColor="var(--art-gray-800)"
                  :successText="$t('login.sliderSuccessText')"
                  :progressBarBg="getCssVar('--el-color-primary')"
                  background="var(--art-gray-200)"
                  handlerBg="var(--art-main-bg-color)"
                />
              </div>
              <p class="error-text" :class="{ 'show-error-text': !isPassing && isClickPass }">{{
                $t('login.placeholder[2]')
              }}</p>
            </div>

            <div class="remember-password">
              <ElCheckbox v-model="formData.rememberPassword">{{
                $t('login.rememberPwd')
              }}</ElCheckbox>
            </div>

            <div style="margin-top: 30px">
              <ElButton
                class="login-btn"
                type="primary"
                @click="handleSubmit"
                :loading="loading"
                v-ripple
              >
                {{ $t('login.btnText') }}
              </ElButton>
            </div>


          </ElForm>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import AppConfig from '@/config'
  import { ElNotification, ElMessage } from 'element-plus'
  import { useUserStore } from '@/store/modules/user'
  import { getCssVar } from '@/utils/ui'
  import { SystemThemeEnum } from '@/enums/appEnum'
  import { useI18n } from 'vue-i18n'
  import { HttpError } from '@/utils/http/error'

  defineOptions({ name: 'Login' })

  const { t } = useI18n()
  import { useSettingStore } from '@/store/modules/setting'
  import type { FormInstance, FormRules } from 'element-plus'

  const settingStore = useSettingStore()
  const { isDark, systemThemeType } = storeToRefs(settingStore)

  const dragVerify = ref()

  const userStore = useUserStore()
  const router = useRouter()
  const isPassing = ref(false)
  const isClickPass = ref(false)

  const systemName = AppConfig.systemInfo.name
  const formRef = ref<FormInstance>()

  const formData = reactive({
    username: '',
    password: '',
    rememberPassword: true
  })

  const rules = computed<FormRules>(() => ({
    username: [{ required: true, message: t('login.placeholder[0]'), trigger: 'blur' }],
    password: [{ required: true, message: t('login.placeholder[1]'), trigger: 'blur' }]
  }))

  const loading = ref(false)

  // 登录
  const handleSubmit = async () => {
    if (!formRef.value) return

    try {
      // 表单验证
      const valid = await formRef.value.validate()
      if (!valid) return

      // 拖拽验证
      if (!isPassing.value) {
        isClickPass.value = true
        return
      }

      loading.value = true

      // 登录请求
      const { username, password } = formData

      const loginResponse = await UserService.login({
        username,
        password,
        captcha: '1234' // 临时验证码
      })

      // 验证token
      if (!loginResponse.token) {
        throw new Error('Login failed - no token received')
      }

      // 存储token和用户信息
      userStore.setToken(loginResponse.token)
      const profileData = await UserService.getUserInfo()
      userStore.setUserInfo(profileData.user_info)
      userStore.setPermissions(profileData.permissions)
      userStore.setLoginStatus(true)

      // 登录成功处理
      showLoginSuccessNotice()
      router.push('/')
    } catch (error) {
      // 处理 HttpError
      if (error instanceof HttpError) {
        ElMessage.error(error.message)
        console.error('[Login] HttpError:', error.toLogData())
      } else {
        // 处理非 HttpError
        ElMessage.error('登录失败，请稍后重试')
        console.error('[Login] Unexpected error:', error)
      }
    } finally {
      loading.value = false
      resetDragVerify()
    }
  }

  // 重置拖拽验证
  const resetDragVerify = () => {
    dragVerify.value.reset()
  }

  // 登录成功提示
  const showLoginSuccessNotice = () => {
    setTimeout(() => {
      const userName = userStore.info.name || userStore.info.username || '用户'
      ElNotification({
        title: t('login.success.title'),
        type: 'success',
        duration: 2500,
        zIndex: 10000,
        message: `${t('login.success.message')}, ${userName}!`
      })
    }, 150)
  }



  // 切换主题
  import { useTheme } from '@/composables/useTheme'
  import { UserService } from '@/api/usersApi'

  const toggleTheme = () => {
    let { LIGHT, DARK } = SystemThemeEnum
    useTheme().switchThemeStyles(systemThemeType.value === LIGHT ? DARK : LIGHT)
  }
</script>

<style lang="scss" scoped>
  @use './index';
</style>
