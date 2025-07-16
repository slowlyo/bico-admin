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
              placeholder="请输入密码"
              show-password
            />
          </ElFormItem>
        </ElCol>
        <ElCol :span="12">
          <ElFormItem label="状态" prop="enabled">
            <ElSwitch 
              v-model="formData.enabled"
              checked-text="启用"
              unchecked-text="禁用"
            />
          </ElFormItem>
        </ElCol>
      </ElRow>

      <ElFormItem label="角色" prop="role_ids">
        <ElSelect 
          v-model="formData.role_ids" 
          multiple 
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
        <ElInput v-model="formData.avatar" placeholder="请输入头像URL" />
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
  import { AdminUserService, RoleService } from '@/api/adminUserApi'

  interface Props {
    visible: boolean
    type: string
    userData?: Api.AdminUser.AdminUserInfo
  }

  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'submit'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  // 角色列表数据
  const roleList = ref<Api.Role.RoleOption[]>([])
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
  const formData = reactive<Api.AdminUser.AdminUserCreateRequest>({
    username: '',
    password: '',
    name: '',
    avatar: '',
    email: '',
    phone: '',
    remark: '',
    enabled: true,
    role_ids: []
  })

  // 表单验证规则
  const rules: FormRules = {
    username: [
      { required: true, message: '请输入用户名', trigger: 'blur' },
      { min: 3, max: 50, message: '用户名长度为3-50个字符', trigger: 'blur' }
    ],
    password: [
      { 
        required: dialogType.value === 'add', 
        message: '请输入密码', 
        trigger: 'blur' 
      },
      { min: 6, max: 100, message: '密码长度为6-100个字符', trigger: 'blur' }
    ],
    name: [
      { required: true, message: '请输入姓名', trigger: 'blur' },
      { max: 100, message: '姓名长度不能超过100个字符', trigger: 'blur' }
    ],
    email: [
      { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
    ]
  }

  // 加载角色列表
  const loadRoles = async () => {
    try {
      const response = await RoleService.getActiveRoles()
      roleList.value = response || []
    } catch (error) {
      console.error('获取角色列表失败:', error)
    }
  }

  // 初始化表单数据
  const initFormData = () => {
    const isEdit = props.type === 'edit' && props.userData
    const row = props.userData

    if (isEdit && row) {
      Object.assign(formData, {
        username: row.username || '',
        password: '', // 编辑时密码为空，表示不修改
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
          if (dialogType.value === 'add') {
            await AdminUserService.createAdminUser(formData)
          } else {
            const updateData = { ...formData } as Api.AdminUser.AdminUserUpdateRequest
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
          ElMessage.error('操作失败请重试！')
        } finally {
          submitLoading.value = false
        }
      }
    })
  }

  // 组件挂载时加载角色列表
  onMounted(() => {
    loadRoles()
  })
</script>

<style lang="scss" scoped>
  .dialog-footer {
    text-align: right;
  }
</style>
