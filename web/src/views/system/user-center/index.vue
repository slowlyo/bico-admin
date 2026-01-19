<!-- 个人中心页面 -->
<template>
  <div class="w-full h-full p-0 bg-transparent border-none shadow-none text-g-700">
    <div class="flex flex-col lg:flex-row items-start mt-2.5 gap-5">
      <!-- 左侧：头像设置 -->
      <div class="w-full lg:w-96 shrink-0">
        <div class="art-card-sm">
          <div class="p-4 border-b border-g-300">
            <h1 class="text-base font-medium">头像设置</h1>
          </div>
          <div class="p-6 text-center">
            <div class="flex flex-col items-center justify-center py-8">
              <div
                class="relative w-32 h-32 rounded-full overflow-hidden cursor-pointer group"
                @click="triggerUpload"
              >
                <img
                  class="w-full h-full object-cover border border-g-200 dark:border-g-300 rounded-full shadow-sm"
                  :src="userInfo.avatar || '/src/assets/imgs/user/avatar.webp'"
                />
                <div
                  class="absolute inset-0 bg-black/50 flex flex-col items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300 text-white"
                >
                  <ArtSvgIcon icon="ri:camera-line" class="text-2xl mb-1" />
                  <span class="text-xs">{{ uploading ? '上传中...' : '点击上传' }}</span>
                </div>
                <input
                  ref="fileInput"
                  type="file"
                  class="hidden"
                  accept="image/*"
                  @change="handleAvatarChange"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧：表单内容 -->
      <div class="flex-1 flex flex-col gap-5">
        <!-- 基本信息 -->
        <div class="art-card-sm">
          <div class="p-4 border-b border-g-300">
            <h1 class="text-base font-medium">基本信息</h1>
          </div>

          <ElForm
            :model="form"
            class="p-6"
            ref="ruleFormRef"
            :rules="rules"
            label-width="80px"
          >
            <ElFormItem label="用户名" prop="username">
              <ElInput v-model="form.username" disabled />
            </ElFormItem>
            <ElFormItem label="姓名" prop="name">
              <ElInput v-model="form.name" placeholder="请输入姓名" />
            </ElFormItem>

            <div class="flex justify-start ml-20 mt-4">
              <ElButton type="primary" class="w-24" :loading="updating" @click="saveProfile">
                保存
              </ElButton>
            </div>
          </ElForm>
        </div>

        <!-- 修改密码 -->
        <div class="art-card-sm">
          <div class="p-4 border-b border-g-300">
            <h1 class="text-base font-medium">修改密码</h1>
          </div>

          <ElForm
            :model="pwdForm"
            class="p-6"
            ref="pwdFormRef"
            :rules="pwdRules"
            label-width="80px"
          >
            <ElFormItem label="原密码" prop="old_password">
              <ElInput
                v-model="pwdForm.old_password"
                type="password"
                placeholder="请输入原密码"
                show-password
              />
            </ElFormItem>

            <ElFormItem label="新密码" prop="new_password">
              <ElInput
                v-model="pwdForm.new_password"
                type="password"
                placeholder="请输入新密码"
                show-password
              />
            </ElFormItem>

            <ElFormItem label="确认密码" prop="confirm_password">
              <ElInput
                v-model="pwdForm.confirm_password"
                type="password"
                placeholder="请再次输入新密码"
                show-password
              />
            </ElFormItem>

            <div class="flex justify-start ml-20 mt-4">
              <ElButton type="primary" class="w-24" :loading="pwdUpdating" @click="changePwd">
                修改密码
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
  import { fetchUpdateProfile, fetchChangePassword, fetchGetUserInfo, fetchUploadAvatar } from '@/api/auth'
  import type { FormInstance, FormRules } from 'element-plus'

  defineOptions({ name: 'UserCenter' })

  const userStore = useUserStore()
  const userInfo = computed(() => userStore.getUserInfo)

  const ruleFormRef = ref<FormInstance>()
  const pwdFormRef = ref<FormInstance>()
  const fileInput = ref<HTMLInputElement>()

  const uploading = ref(false)
  const updating = ref(false)
  const pwdUpdating = ref(false)

  /**
   * 用户信息表单
   */
  const form = reactive({
    name: '',
    username: '',
    avatar: ''
  })

  watch(
    () => userInfo.value,
    (val) => {
      if (val) {
        form.name = val.name || ''
        form.username = val.username || ''
        form.avatar = val.avatar || ''
      }
    },
    { immediate: true, deep: true }
  )

  /**
   * 密码修改表单
   */
  const pwdForm = reactive({
    old_password: '',
    new_password: '',
    confirm_password: ''
  })

  /**
   * 基本信息验证规则
   */
  const rules = reactive<FormRules>({
    name: [{ required: true, message: '请输入姓名', trigger: 'blur' }]
  })

  /**
   * 密码修改验证规则
   */
  const pwdRules = reactive<FormRules>({
    old_password: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
    new_password: [
      { required: true, message: '请输入新密码', trigger: 'blur' },
      { min: 6, message: '密码长度至少6位', trigger: 'blur' }
    ],
    confirm_password: [
      { required: true, message: '请确认新密码', trigger: 'blur' },
      {
        validator: (_rule, value, callback) => {
          if (value !== pwdForm.new_password) {
            callback(new Error('两次输入密码不一致'))
          } else {
            callback()
          }
        },
        trigger: 'blur'
      }
    ]
  })

  /**
   * 触发文件上传
   */
  const triggerUpload = () => {
    if (uploading.value) return
    fileInput.value?.click()
  }

  /**
   * 处理头像上传
   */
  const handleAvatarChange = async (event: Event) => {
    const target = event.target as HTMLInputElement
    if (!target.files?.length) return

    const file = target.files[0]
    uploading.value = true

    try {
      const res: any = await fetchUploadAvatar(file)
      const avatarUrl = res.url || (res.data && res.data.url)
      
      if (avatarUrl) {
        // 更新个人资料中的头像
        await fetchUpdateProfile({ avatar: avatarUrl })
        ElMessage.success('头像上传成功')
        // 刷新用户信息
        refreshUserInfo()
      }
    } catch (error) {
      console.error('上传头像失败:', error)
    } finally {
      uploading.value = false
      if (fileInput.value) fileInput.value.value = ''
    }
  }

  /**
   * 刷新用户信息
   */
  const refreshUserInfo = async () => {
    try {
      const newUserInfo = await fetchGetUserInfo()
      userStore.setUserInfo(newUserInfo)
    } catch (error) {
      console.error('刷新用户信息失败:', error)
    }
  }

  /**
   * 保存个人信息
   */
  const saveProfile = async () => {
    if (!ruleFormRef.value) return
    await ruleFormRef.value.validate(async (valid) => {
      if (valid) {
        updating.value = true
        try {
          await fetchUpdateProfile({
            name: form.name,
            avatar: form.avatar
          })
          ElMessage.success('保存成功')
          refreshUserInfo()
        } catch (error) {
          console.error('保存失败:', error)
        } finally {
          updating.value = false
        }
      }
    })
  }

  /**
   * 修改密码
   */
  const changePwd = async () => {
    if (!pwdFormRef.value) return
    await pwdFormRef.value.validate(async (valid) => {
      if (valid) {
        pwdUpdating.value = true
        try {
          await fetchChangePassword({
            old_password: pwdForm.old_password,
            new_password: pwdForm.new_password
          })
          ElMessage.success('修改成功')
          // 清空表单
          pwdForm.old_password = ''
          pwdForm.new_password = ''
          pwdForm.confirm_password = ''
          const formRef = pwdFormRef.value
          if (formRef) {
            formRef.resetFields()
          }
        } catch (error) {
          console.error('修改密码失败:', error)
        } finally {
          pwdUpdating.value = false
        }
      }
    })
  }
</script>
