<template>
  <ElDialog
    v-model="dialogVisible"
    :title="dialogType === 'add' ? '添加用户' : '编辑用户'"
    width="30%"
    align-center
    @close="handleClose"
  >
    <ElForm
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="80px"
      :validate-on-rule-change="false"
    >
      <ElFormItem label="头像" prop="avatar">
        <ElUpload
          class="avatar-uploader"
          action=""
          accept="image/*"
          :show-file-list="false"
          :http-request="handleAvatarUpload"
          :before-upload="beforeAvatarUpload"
        >
          <img v-if="formData.avatar" :src="formData.avatar" class="avatar" />
          <ElIcon v-else class="avatar-uploader-icon"><Plus /></ElIcon>
        </ElUpload>
      </ElFormItem>
      <ElFormItem label="用户名" prop="username">
        <ElInput v-model="formData.username" placeholder="请输入用户名" :disabled="dialogType === 'edit'" />
      </ElFormItem>
      <ElFormItem label="姓名" prop="name">
        <ElInput v-model="formData.name" placeholder="请输入姓名" />
      </ElFormItem>
      <ElFormItem label="密码" prop="password" v-if="dialogType === 'add'">
        <ElInput v-model="formData.password" type="password" placeholder="请输入密码" show-password />
      </ElFormItem>
      <ElFormItem label="新密码" prop="password" v-else>
        <ElInput v-model="formData.password" type="password" placeholder="不修改请留空" show-password />
      </ElFormItem>
      <ElFormItem label="角色" prop="role_ids">
        <ElSelect v-model="formData.role_ids" multiple filterable placeholder="请选择角色">
          <ElOption
            v-for="role in roleList"
            :key="role.id"
            :value="role.id"
            :label="role.name"
          />
        </ElSelect>
      </ElFormItem>
      <ElFormItem label="状态" prop="enabled">
        <ElSwitch v-model="formData.enabled" />
      </ElFormItem>
    </ElForm>
    <template #footer>
      <div class="dialog-footer">
        <ElButton @click="dialogVisible = false">取消</ElButton>
        <ElButton type="primary" @click="handleSubmit">提交</ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { fetchGetAllRoles, fetchCreateUser, fetchUpdateUser, fetchUploadFile } from '@/api/system-manage'
  import type { FormInstance, FormRules, UploadProps } from 'element-plus'
  import { Plus } from '@element-plus/icons-vue'

  interface Props {
    visible: boolean
    type: string
    userData?: Partial<Api.SystemManage.UserListItem>
  }

  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'submit'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  // 角色列表数据
  const roleList = ref<Api.SystemManage.RoleListItem[]>([])

  const loadRoleList = async () => {
    try {
      const res = await fetchGetAllRoles()
      roleList.value = res
    } catch (error) {
      console.error('获取角色列表失败:', error)
    }
  }

  // 对话框显示控制
  const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
  })

  const dialogType = computed(() => props.type)

  // 表单实例
  const formRef = ref<FormInstance>()

  // 表单数据
  const formData = reactive({
    username: '',
    password: '' as string | undefined,
    name: '',
    avatar: '',
    enabled: true,
    role_ids: [] as number[]
  })

  // 表单验证规则
  const rules = computed<FormRules>(() => ({
    username: [
      { required: true, message: '请输入用户名', trigger: 'blur' },
      { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
    ],
    name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
    password: props.type === 'add' ? [{ required: true, message: '请输入密码', trigger: 'blur' }] : []
  }))

  /**
   * 初始化表单数据
   * 根据对话框类型（新增/编辑）填充表单
   */
  const initFormData = () => {
    const isEdit = props.type === 'edit' && props.userData
    const row = props.userData

    Object.assign(formData, {
      username: isEdit && row ? row.username || '' : '',
      password: '',
      name: isEdit && row ? row.name || '' : '',
      avatar: isEdit && row ? row.avatar || '' : '',
      enabled: isEdit && row ? !!row.enabled : true,
      role_ids: isEdit && row ? row.roles?.map((r) => r.id) || [] : []
    })
  }

  onMounted(() => {
    loadRoleList()
  })

  watch(
    () => props.visible,
    (visible) => {
      if (visible) {
        initFormData()
        nextTick(() => {
          formRef.value?.clearValidate()
        })
      }
    }
  )

  // 监听类型和数据的变化，仅在弹窗打开时初始化
  watch(
    () => [props.type, props.userData],
    () => {
      if (props.visible) {
        initFormData()
        nextTick(() => {
          formRef.value?.clearValidate()
        })
      }
    },
    { deep: true }
  )

  /**
   * 上传头像前校验
   */
  const beforeAvatarUpload: UploadProps['beforeUpload'] = (rawFile) => {
    const isImage = rawFile.type.startsWith('image/')
    if (!isImage) {
      ElMessage.error('头像图片必须是图片格式!')
      return false
    } else if (rawFile.size / 1024 / 1024 > 2) {
      ElMessage.error('头像图片大小不能超过 2MB!')
      return false
    }
    return true
  }

  /**
   * 自定义头像上传
   */
  const handleAvatarUpload = async (options: any) => {
    try {
      const res = await fetchUploadFile(options.file)
      formData.avatar = res.url
    } catch (error) {
      console.error('上传头像失败:', error)
    }
  }

  /**
   * 关闭弹窗并重置表单
   */
  const handleClose = () => {
    dialogVisible.value = false
    setTimeout(() => {
      initFormData()
      formRef.value?.clearValidate()
    }, 300)
  }

  /**
   * 提交表单
   * 验证通过后触发提交事件
   */
  const handleSubmit = async () => {
    if (!formRef.value) return

    await formRef.value.validate(async (valid) => {
      if (valid) {
        try {
          if (props.type === 'add') {
            await fetchCreateUser(formData as any)
            ElMessage.success('添加成功')
          } else {
            const id = props.userData?.id
            if (!id) return
            const updateData = { ...formData }
            if (!updateData.password) {
              delete updateData.password
            }
            await fetchUpdateUser(id, updateData as any)
            ElMessage.success('更新成功')
          }
          dialogVisible.value = false
          emit('submit')
        } catch (error) {
          console.error('提交用户失败:', error)
        }
      }
    })
  }
</script>

<style scoped>
  .avatar-uploader :deep(.el-upload) {
    position: relative;
    overflow: hidden;
    cursor: pointer;
    border: 1px dashed var(--el-border-color);
    border-radius: 6px;
    transition: var(--el-transition-duration-fast);
  }

  .avatar-uploader :deep(.el-upload:hover) {
    border-color: var(--el-color-primary);
  }

  .el-icon.avatar-uploader-icon {
    width: 100px;
    height: 100px;
    font-size: 28px;
    color: #8c939d;
    text-align: center;
  }

  .avatar {
    display: block;
    width: 100px;
    height: 100px;
    object-fit: cover;
  }
</style>
