<template>
  <ElDialog
    v-model="dialogVisible"
    :title="dialogType === 'add' ? '新建角色' : '编辑角色'"
    width="600px"
    align-center
  >
    <ElForm ref="formRef" :model="formData" :rules="rules" label-width="100px">
      <ElRow :gutter="20">
        <ElCol :span="12">
          <ElFormItem label="角色名称" prop="name">
            <ElInput v-model="formData.name" placeholder="请输入角色名称" />
          </ElFormItem>
        </ElCol>
        <ElCol :span="12">
          <ElFormItem label="角色代码" prop="code">
            <ElInput 
              v-model="formData.code" 
              placeholder="请输入角色代码"
              :disabled="dialogType === 'edit'"
            />
          </ElFormItem>
        </ElCol>
      </ElRow>

      <ElFormItem label="状态" prop="status">
        <ElSwitch
          v-model="enabled"
          checked-text="启用"
          unchecked-text="禁用"
        />
      </ElFormItem>

      <ElFormItem label="描述" prop="description">
        <ElInput 
          v-model="formData.description" 
          type="textarea" 
          :rows="3"
          placeholder="请输入角色描述"
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
  import { AdminRoleService, type RoleTypes } from '@/api/adminRoleApi'

  interface Props {
    visible: boolean
    type: string
    roleData?: RoleTypes.RoleInfo
  }

  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'submit'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

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
  const formData = reactive<RoleTypes.RoleCreateRequest>({
    name: '',
    code: '',
    description: '',
    status: 1,
    permissions: []
  })

  // 启用状态的计算属性
  const enabled = computed({
    get: () => formData.status === 1,
    set: (value) => {
      formData.status = value ? 1 : 0
    }
  })

  // 表单验证规则
  const rules: FormRules = {
    name: [
      { required: true, message: '请输入角色名称', trigger: 'blur' },
      { min: 1, max: 100, message: '角色名称长度为1-100个字符', trigger: 'blur' }
    ],
    code: [
      { required: true, message: '请输入角色代码', trigger: 'blur' },
      { min: 1, max: 50, message: '角色代码长度为1-50个字符', trigger: 'blur' },
      { pattern: /^[a-zA-Z0-9_]+$/, message: '角色代码只能包含字母、数字和下划线', trigger: 'blur' }
    ],
    description: [
      { max: 500, message: '描述长度不能超过500个字符', trigger: 'blur' }
    ]
  }

  // 初始化表单数据
  const initFormData = () => {
    const isEdit = props.type === 'edit' && props.roleData
    const row = props.roleData

    if (isEdit && row) {
      Object.assign(formData, {
        name: row.name || '',
        code: row.code || '',
        description: row.description || '',
        status: row.status,
        permissions: row.permissions?.map(p => p.permission_code) || []
      })
    } else {
      // 重置为初始值
      Object.assign(formData, {
        name: '',
        code: '',
        description: '',
        status: 1,
        permissions: []
      })
    }
  }

  // 统一监听对话框状态变化
  watch(
    () => [props.visible, props.type, props.roleData],
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
            await AdminRoleService.createRole(formData)
          } else {
            const updateData = { ...formData } as RoleTypes.RoleUpdateRequest
            await AdminRoleService.updateRole(props.roleData!.id, updateData)
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
