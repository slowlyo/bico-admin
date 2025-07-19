<template>
  <ElDialog
    v-model="dialogVisible"
    :title="dialogType === 'add' ? '新建管理员' : '编辑管理员'"
    width="600px"
    align-center
  >
    <ElForm ref="formRef" :model="formData" :rules="rules" label-width="100px">
      <ElRow :gutter="20">
        <ElCol :span="12">
          <ElFormItem label="用户名" prop="username">
            <ElInput 
              v-model="formData.username" 
              placeholder="请输入用户名"
              :disabled="dialogType === 'edit'"
            />
          </ElFormItem>
        </ElCol>
        <ElCol :span="12">
          <ElFormItem label="姓名" prop="name">
            <ElInput v-model="formData.name" placeholder="请输入姓名" />
          </ElFormItem>
        </ElCol>
      </ElRow>

      <ElRow :gutter="20">
        <ElCol :span="12">
          <ElFormItem label="邮箱" prop="email">
            <ElInput v-model="formData.email" placeholder="请输入邮箱" />
          </ElFormItem>
        </ElCol>
        <ElCol :span="12">
          <ElFormItem label="手机号" prop="phone">
            <ElInput v-model="formData.phone" placeholder="请输入手机号" />
          </ElFormItem>
        </ElCol>
      </ElRow>

      <ElRow :gutter="20">
        <ElCol :span="12">
          <ElFormItem label="密码" prop="password">
            <ElInput
              v-model="formData.password"
              type="password"
              :placeholder="dialogType === 'edit' ? '留空表示不修改密码' : '请输入密码'"
              show-password
            />
          </ElFormItem>
        </ElCol>
        <ElCol :span="12">
          <ElFormItem label="确认密码" prop="confirmPassword">
            <ElInput
              v-model="formData.confirmPassword"
              type="password"
              :placeholder="dialogType === 'edit' ? '留空表示不修改密码' : '请再次输入密码'"
              show-password
            />
          </ElFormItem>
        </ElCol>
      </ElRow>

      <ElRow :gutter="20">
        <ElCol :span="12">
          <ElFormItem label="状态" prop="enabled">
            <ElSwitch
              v-model="formData.enabled"
              checked-text="启用"
              unchecked-text="禁用"
            />
          </ElFormItem>
        </ElCol>
        <ElCol :span="12">
          <!-- 占位列 -->
        </ElCol>
      </ElRow>

      <ElFormItem label="角色" prop="role_ids">
        <ElSelect 
          v-model="formData.role_ids" 
          multiple 
          filterable
          placeholder="请选择角色"
          style="width: 100%"
        >
          <ElOption
            v-for="role in roleList"
            :key="role.id"
            :value="role.id"
            :label="role.name"
          />
        </ElSelect>
      </ElFormItem>

      <ElFormItem label="头像" prop="avatar">
        <AvatarUpload
          v-model="formData.avatar"
          :size="80"
          @change="handleAvatarChange"
        />
      </ElFormItem>

      <ElFormItem label="备注" prop="remark">
        <ElInput 
          v-model="formData.remark" 
          type="textarea" 
          :rows="3"
          placeholder="请输入备注"
        />
      </ElFormItem>
    </ElForm>

    <template #footer>
      <div class="dialog-footer">
        <ElButton @click="dialogVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="submitLoading" @click="handleSubmit">
          {{ dialogType === 'add' ? '创建' : '更新' }}
        </ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import type { FormInstance, FormRules } from 'element-plus'
  import { ElMessage } from 'element-plus'
  import { AdminUserService, type AdminUserTypes, type RoleTypes } from '@/api/adminUserApi'
  import AvatarUpload from '@/components/custom/avatar-upload/index.vue'

  interface Props {
    visible: boolean
    type: string
    userData?: AdminUserTypes.AdminUserInfo
    roleOptions: RoleTypes.RoleOption[]
  }

  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'submit'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  // 角色列表数据 - 使用 props 传入的数据
  const roleList = computed(() => props.roleOptions)
  const submitLoading = ref(false)

  // 对话框显示控制
  const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
  })

  const dialogType = computed(() => props.type)

  // 表单实例
  const formRef = ref<FormInstance>()

  // 表单数据
  const formData = reactive<AdminUserTypes.AdminUserCreateRequest & { confirmPassword: string }>({
    username: '',
    password: '',
    confirmPassword: '',
    name: '',
    avatar: '',
    email: '',
    phone: '',
    remark: '',
    enabled: true,
    role_ids: []
  })

  // 密码确认验证器
  const validateConfirmPassword = (_rule: any, value: string, callback: any) => {
    const isAdd = dialogType.value === 'add'
    const hasPassword = !!formData.password

    // 新增时确认密码必填，编辑时如果有密码则确认密码必填
    if ((isAdd || hasPassword) && !value) {
      callback(new Error('请再次输入密码'))
      return
    }

    // 如果有确认密码值，检查是否与密码一致
    if (value && value !== formData.password) {
      callback(new Error('两次输入的密码不一致'))
      return
    }

    callback()
  }

  // 表单验证规则
  const rules: FormRules = {
    username: [
      { required: true, message: '请输入用户名', trigger: 'blur' },
      { min: 3, max: 50, message: '用户名长度为3-50个字符', trigger: 'blur' }
    ],
    password: [
      {
        validator: (_rule: any, value: string, callback: any) => {
          // 新增时密码必填
          if (dialogType.value === 'add' && !value) {
            callback(new Error('请输入密码'))
            return
          }
          // 如果有值，检查长度
          if (value && (value.length < 6 || value.length > 100)) {
            callback(new Error('密码长度为6-100个字符'))
            return
          }
          // 如果密码有变化，需要重新验证确认密码
          if (value && formData.confirmPassword) {
            formRef.value?.validateField('confirmPassword')
          }
          callback()
        },
        trigger: 'blur'
      }
    ],
    confirmPassword: [
      { validator: validateConfirmPassword, trigger: 'blur' }
    ],
    name: [
      { required: true, message: '请输入姓名', trigger: 'blur' },
      { max: 100, message: '姓名长度不能超过100个字符', trigger: 'blur' }
    ],
    email: [
      { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
    ]
  }



  // 初始化表单数据
  const initFormData = () => {
    const isEdit = props.type === 'edit' && props.userData
    const row = props.userData

    if (isEdit && row) {
      Object.assign(formData, {
        username: row.username || '',
        password: '', // 编辑时密码为空，表示不修改
        confirmPassword: '', // 编辑时确认密码也为空
        name: row.name || '',
        avatar: row.avatar || '',
        email: row.email || '',
        phone: row.phone || '',
        remark: row.remark || '',
        enabled: row.status === 1,
        role_ids: row.roles?.map(role => role.id) || []
      })
    } else {
      // 重置为初始值
      Object.assign(formData, {
        username: '',
        password: '',
        confirmPassword: '',
        name: '',
        avatar: '',
        email: '',
        phone: '',
        remark: '',
        enabled: true,
        role_ids: []
      })
    }
  }

  // 头像变更处理
  const handleAvatarChange = (value: string) => {
    formData.avatar = value
  }

  // 统一监听对话框状态变化
  watch(
    () => [props.visible, props.type, props.userData],
    ([visible]) => {
      if (visible) {
        initFormData()
        nextTick(() => {
          formRef.value?.clearValidate()
        })
      }
    },
    { immediate: true }
  )

  // 提交表单
  const handleSubmit = async () => {
    if (!formRef.value) return

    await formRef.value.validate(async (valid) => {
      if (valid) {
        submitLoading.value = true
        try {
          // 准备提交数据，移除确认密码字段
          const { confirmPassword, ...submitData } = formData

          if (dialogType.value === 'add') {
            await AdminUserService.createAdminUser(submitData)
          } else {
            const updateData = { ...submitData } as AdminUserTypes.AdminUserUpdateRequest
            // 如果密码为空，则不更新密码
            if (!updateData.password) {
              delete updateData.password
            }
            await AdminUserService.updateAdminUser(props.userData!.id, updateData)
          }

          ElMessage.success(dialogType.value === 'add' ? '创建成功' : '更新成功')
          dialogVisible.value = false
          emit('submit')
        } catch (error) {
          console.error('提交失败:', error)
          // HTTP 拦截器已经显示了具体的错误消息，这里不再重复显示
        } finally {
          submitLoading.value = false
        }
      }
    })
  }


</script>

<style lang="scss" scoped>
  .dialog-footer {
    text-align: right;
  }
</style>
