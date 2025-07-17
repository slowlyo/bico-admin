<template>
  <div class="profile-page art-full-height">
    <div class="profile-container">
      <!-- 个人信息卡片 -->
      <ElCard class="profile-card art-table-card" shadow="never" :style="{ borderRadius: settingStore.getCustomRadius }">
        <template #header>
          <div class="card-header">
            <span>个人信息</span>
          </div>
        </template>

        <ElForm
          ref="ruleFormRef"
          :model="form"
          :rules="rules"
          class="profile-form"
          label-width="100px"
          label-position="left"
        >
          <!-- 头像上传 -->
          <ElFormItem label="头像" prop="avatar">
            <AvatarUpload
              v-model="form.avatar"
              :disabled="!isEdit"
              @change="handleAvatarChange"
            />
          </ElFormItem>

          <ElRow :gutter="20">
            <ElCol :span="12">
              <ElFormItem label="姓名" prop="name">
                <ElInput v-model="form.name" :disabled="!isEdit" placeholder="请输入姓名" />
              </ElFormItem>
            </ElCol>
            <ElCol :span="12">
              <ElFormItem label="邮箱" prop="email">
                <ElInput v-model="form.email" :disabled="!isEdit" placeholder="请输入邮箱" />
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElRow :gutter="20">
            <ElCol :span="12">
              <ElFormItem label="手机号" prop="phone">
                <ElInput v-model="form.phone" :disabled="!isEdit" placeholder="请输入手机号" />
              </ElFormItem>
            </ElCol>
            <ElCol :span="12">
              <ElFormItem label="用户名">
                <ElInput :value="userInfo.username" disabled placeholder="用户名不可修改" />
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElFormItem>
            <ElButton type="primary" @click="handleProfileSubmit" :loading="profileLoading">
              {{ isEdit ? '保存' : '编辑' }}
            </ElButton>
            <ElButton v-if="isEdit" @click="cancelEdit">取消</ElButton>
          </ElFormItem>
        </ElForm>
      </ElCard>
      <!-- 修改密码卡片 -->
      <ElCard class="password-card art-table-card" shadow="never" :style="{ borderRadius: settingStore.getCustomRadius }">
        <template #header>
          <div class="card-header">
            <span>修改密码</span>
          </div>
        </template>

        <ElForm
          ref="pwdFormRef"
          :model="pwdForm"
          :rules="pwdRules"
          class="password-form"
          label-width="100px"
          label-position="left"
        >
          <ElRow :gutter="20">
            <ElCol :span="12">
              <ElFormItem label="当前密码" prop="old_password">
                <ElInput
                  v-model="pwdForm.old_password"
                  type="password"
                  :disabled="!isEditPwd"
                  show-password
                  placeholder="请输入当前密码"
                />
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElRow :gutter="20">
            <ElCol :span="12">
              <ElFormItem label="新密码" prop="new_password">
                <ElInput
                  v-model="pwdForm.new_password"
                  type="password"
                  :disabled="!isEditPwd"
                  show-password
                  placeholder="请输入新密码"
                />
              </ElFormItem>
            </ElCol>
            <ElCol :span="12">
              <ElFormItem label="确认新密码" prop="confirm_password">
                <ElInput
                  v-model="pwdForm.confirm_password"
                  type="password"
                  :disabled="!isEditPwd"
                  show-password
                  placeholder="请再次输入新密码"
                />
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElFormItem>
            <ElButton type="primary" @click="handlePasswordSubmit" :loading="passwordLoading">
              {{ isEditPwd ? '保存' : '编辑' }}
            </ElButton>
            <ElButton v-if="isEditPwd" @click="cancelPasswordEdit">取消</ElButton>
          </ElFormItem>
        </ElForm>
      </ElCard>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, computed, onMounted, watch } from 'vue'
  import { useUserStore } from '@/store/modules/user'
  import { useSettingStore } from '@/store/modules/setting'
  import { UserService } from '@/api/usersApi'
  import { ElForm, ElMessage, FormInstance, FormRules } from 'element-plus'
  import AvatarUpload from '@/components/custom/avatar-upload/index.vue'

  defineOptions({ name: 'UserCenter' })

  const userStore = useUserStore()
  const settingStore = useSettingStore()
  const userInfo = computed(() => userStore.getUserInfo)

  // 编辑状态
  const isEdit = ref(false)
  const isEditPwd = ref(false)
  const profileLoading = ref(false)
  const passwordLoading = ref(false)

  // 表单引用
  const ruleFormRef = ref<FormInstance>()
  const pwdFormRef = ref<FormInstance>()

  // 个人信息表单
  const form = reactive<Api.User.ProfileUpdateRequest>({
    name: '',
    avatar: '',
    email: '',
    phone: ''
  })

  // 密码表单
  const pwdForm = reactive<Api.User.ChangePasswordRequest & { confirm_password: string }>({
    old_password: '',
    new_password: '',
    confirm_password: ''
  })

  // 个人信息验证规则
  const rules = reactive<FormRules>({
    name: [
      { required: true, message: '请输入姓名', trigger: 'blur' },
      { min: 1, max: 100, message: '姓名长度为1-100个字符', trigger: 'blur' }
    ],
    email: [
      { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' },
      { max: 100, message: '邮箱长度不能超过100个字符', trigger: 'blur' }
    ],
    phone: [
      { max: 20, message: '手机号长度不能超过20个字符', trigger: 'blur' }
    ],
    avatar: [
      { max: 255, message: '头像URL长度不能超过255个字符', trigger: 'blur' }
    ]
  })

  // 密码验证规则
  const pwdRules = reactive<FormRules>({
    old_password: [
      { required: true, message: '请输入当前密码', trigger: 'blur' },
      { min: 6, max: 100, message: '密码长度为6-100个字符', trigger: 'blur' }
    ],
    new_password: [
      { required: true, message: '请输入新密码', trigger: 'blur' },
      { min: 6, max: 100, message: '密码长度为6-100个字符', trigger: 'blur' }
    ],
    confirm_password: [
      { required: true, message: '请确认新密码', trigger: 'blur' },
      {
        validator: (_rule, value, callback) => {
          if (value !== pwdForm.new_password) {
            callback(new Error('两次输入的密码不一致'))
          } else {
            callback()
          }
        },
        trigger: 'blur'
      }
    ]
  })

  // 头像变更处理
  const handleAvatarChange = (value: string) => {
    form.avatar = value
  }

  // 初始化表单数据
  const initFormData = () => {
    if (userInfo.value) {
      form.name = userInfo.value.name || ''
      form.avatar = userInfo.value.avatar || ''
      form.email = userInfo.value.email || ''
      form.phone = userInfo.value.phone || ''
      console.log('表单数据已初始化:', form)
    } else {
      console.log('用户信息尚未加载')
    }
  }



  // 取消编辑
  const cancelEdit = () => {
    isEdit.value = false
    initFormData()
    ruleFormRef.value?.clearValidate()
  }

  // 取消密码编辑
  const cancelPasswordEdit = () => {
    isEditPwd.value = false
    pwdForm.old_password = ''
    pwdForm.new_password = ''
    pwdForm.confirm_password = ''
    pwdFormRef.value?.clearValidate()
  }



  // 处理个人信息提交
  const handleProfileSubmit = async () => {
    if (!isEdit.value) {
      isEdit.value = true
      return
    }

    if (!ruleFormRef.value) return

    await ruleFormRef.value.validate(async (valid) => {
      if (valid) {
        profileLoading.value = true
        try {
          await UserService.updateProfile(form)
          ElMessage.success('个人信息更新成功')
          isEdit.value = false
          // 刷新用户信息
          const profileData = await UserService.getUserInfo()
          userStore.setUserInfo(profileData.user_info)
          userStore.setPermissions(profileData.permissions)
          initFormData()
        } catch (error) {
          console.error('更新个人信息失败:', error)
          ElMessage.error('更新失败，请重试！')
        } finally {
          profileLoading.value = false
        }
      }
    })
  }

  // 处理密码提交
  const handlePasswordSubmit = async () => {
    if (!isEditPwd.value) {
      isEditPwd.value = true
      return
    }

    if (!pwdFormRef.value) return

    await pwdFormRef.value.validate(async (valid) => {
      if (valid) {
        passwordLoading.value = true
        try {
          const { confirm_password, ...passwordData } = pwdForm
          await UserService.changePassword(passwordData)
          ElMessage.success('密码修改成功')
          isEditPwd.value = false
          // 清空密码表单
          pwdForm.old_password = ''
          pwdForm.new_password = ''
          pwdForm.confirm_password = ''
        } catch (error) {
          console.error('修改密码失败:', error)
          ElMessage.error('修改密码失败，请重试！')
        } finally {
          passwordLoading.value = false
        }
      }
    })
  }

  // 组件挂载时初始化数据
  onMounted(() => {
    initFormData()
  })

  // 监听用户信息变化
  watch(
    () => userInfo.value,
    () => {
      initFormData()
    },
    { deep: true }
  )
</script>

<style lang="scss">
  .user {
    .icon {
      width: 1.4em;
      height: 1.4em;
      overflow: hidden;
      vertical-align: -0.15em;
      fill: currentcolor;
    }
  }
</style>

<style lang="scss" scoped>
  .profile-page {
    padding: 20px;
    height: 100%;

    .profile-container {
      height: 100%;
      display: flex;
      flex-direction: column;
      gap: 20px;

      .profile-card {
        flex: none;
        min-height: auto;

        .card-header {
          font-size: 18px;
          font-weight: 600;
        }
      }

      .password-card {
        flex: none;
        min-height: auto;

        .card-header {
          font-size: 18px;
          font-weight: 600;
        }
      }

      .profile-form,
      .password-form {


          .el-form-item {
            margin-bottom: 20px;
          }
        }
      }
    }
  

  /* 响应式设计 */
  @media (max-width: 768px) {
    .profile-page {
      .profile-container {
        .profile-form,
        .password-form {
          // 响应式样式可以在这里添加
        }
      }
    }
  }
</style>
