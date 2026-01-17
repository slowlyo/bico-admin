<!-- 登录页面 -->
<template>
  <div class="flex w-full h-screen">
    <LoginLeftView />

    <div class="relative flex-1">
      <AuthTopBar />

      <div class="auth-right-wrap">
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
            <ElFormItem prop="username" class="mb-8">
              <ElInput
                class="custom-height"
                :placeholder="$t('login.placeholder.username')"
                v-model.trim="formData.username"
              />
            </ElFormItem>
            <ElFormItem prop="password" class="mb-8">
              <ElInput
                class="custom-height"
                :placeholder="$t('login.placeholder.password')"
                v-model.trim="formData.password"
                type="password"
                autocomplete="off"
                show-password
              />
            </ElFormItem>
            <ElFormItem prop="captchaCode">
              <div class="flex w-full gap-4">
                <ElInput
                  class="custom-height flex-1"
                  :placeholder="$t('login.placeholder.captcha')"
                  v-model.trim="formData.captchaCode"
                />
                <div
                  class="captcha-img cursor-pointer bg-gray-100 dark:bg-zinc-800 rounded overflow-hidden flex items-center justify-center"
                  style="width: 120px; height: 40px"
                  @click="getCaptcha"
                >
                  <img v-if="captchaImg" :src="captchaImg" class="w-full h-full" alt="captcha" />
                  <span v-else class="text-xs text-gray-400">{{ $t('common.loading') }}</span>
                </div>
              </div>
            </ElFormItem>

            <div class="flex-cb mt-2 text-sm">
              <ElCheckbox v-model="formData.rememberPassword">{{
                $t('login.rememberPwd')
              }}</ElCheckbox>
            </div>

            <div style="margin-top: 30px">
              <ElButton
                class="w-full custom-height"
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
  import { useUserStore } from '@/store/modules/user'
  import { useI18n } from 'vue-i18n'
  import { HttpError } from '@/utils/http/error'
  import { fetchLogin, fetchCaptcha, fetchAppConfig, fetchGetUserInfo } from '@/api/auth'
  import { ElNotification, type FormInstance, type FormRules } from 'element-plus'
  import { useSettingStore } from '@/store/modules/setting'
  import { StorageConfig } from '@/utils/storage/storage-config'

  defineOptions({ name: 'Login' })

  const settingStore = useSettingStore()
  const { t } = useI18n()

  const userStore = useUserStore()
  const router = useRouter()
  const route = useRoute()

  const formRef = ref<FormInstance>()

  const formData = reactive({
    username: '',
    password: '',
    captchaId: '',
    captchaCode: '',
    rememberPassword: false
  })

  const captchaImg = ref('')

  const rules = computed<FormRules>(() => ({
    username: [{ required: true, message: t('login.placeholder.username'), trigger: 'blur' }],
    password: [{ required: true, message: t('login.placeholder.password'), trigger: 'blur' }],
    captchaCode: [{ required: true, message: t('login.placeholder.captcha'), trigger: 'blur' }]
  }))

  const loading = ref(false)

  // 获取图形验证码
  const getCaptcha = async () => {
    try {
      const { id, image } = await fetchCaptcha()
      formData.captchaId = id
      captchaImg.value = image
    } catch (error) {
      console.error('获取验证码失败:', error)
    }
  }

  // 获取应用配置
  const getAppConfig = async () => {
    try {
      const config = await fetchAppConfig()
      if (config) {
        settingStore.setAppConfig(config)
      }
    } catch (error) {
      console.error('获取应用配置失败:', error)
    }
  }

  // 加载记住的密码
  const loadRememberedPwd = () => {
    const saved = localStorage.getItem(StorageConfig.REMEMBER_PWD_KEY)
    if (saved) {
      try {
        const { username, password } = JSON.parse(window.atob(saved))
        formData.username = username
        formData.password = password
        formData.rememberPassword = true
      } catch (e) {
        localStorage.removeItem(StorageConfig.REMEMBER_PWD_KEY)
      }
    }
  }

  onMounted(() => {
    getCaptcha()
    getAppConfig()
    loadRememberedPwd()
  })

  // 登录
  const handleSubmit = async () => {
    if (!formRef.value) return

    try {
      // 表单验证
      const valid = await formRef.value.validate()
      if (!valid) return

      loading.value = true

      // 登录请求
      const { username, password, captchaId, captchaCode } = formData

      const { token } = await fetchLogin({
        username,
        password,
        captchaId,
        captchaCode
      })

      // 验证token
      if (!token) {
        throw new Error('Login failed - no token received')
      }

      // 存储 token 和登录状态
      userStore.setToken(token)
      userStore.setLoginStatus(true)

      // 获取用户信息
      const userInfo = await fetchGetUserInfo()
      userStore.setUserInfo(userInfo)

      // 登录成功处理
      showLoginSuccessNotice()

      // 处理记住密码
      if (formData.rememberPassword) {
        const loginInfo = {
          username: formData.username,
          password: formData.password
        }
        localStorage.setItem(StorageConfig.REMEMBER_PWD_KEY, window.btoa(JSON.stringify(loginInfo)))
      } else {
        localStorage.removeItem(StorageConfig.REMEMBER_PWD_KEY)
      }

      // 获取 redirect 参数，如果存在则跳转到指定页面，否则跳转到首页
      const redirect = route.query.redirect as string
      router.push(redirect || '/')
    } catch (error) {
      // 登录失败刷新验证码
      getCaptcha()
      // 处理 HttpError
      if (error instanceof HttpError) {
        // console.log(error.code)
      } else {
        // 处理非 HttpError
        // ElMessage.error('登录失败，请稍后重试')
        console.error('[Login] Unexpected error:', error)
      }
    } finally {
      loading.value = false
    }
  }

  // 登录成功提示
  const showLoginSuccessNotice = () => {
    const userInfo = userStore.getUserInfo
    setTimeout(() => {
      ElNotification({
        title: t('login.success.title'),
        type: 'success',
        duration: 2500,
        zIndex: 10000,
        message: `${t('login.success.message')}, ${userInfo.name || userInfo.username}!`
      })
    }, 1000)
  }
</script>

<style scoped>
  @import './style.css';
</style>

<style lang="scss" scoped></style>
