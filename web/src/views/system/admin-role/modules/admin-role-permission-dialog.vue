<template>
  <ElDialog
    v-model="dialogVisible"
    :title="`配置权限 - ${roleData?.name || ''}`"
    width="600px"
    align-center
  >
    <div v-loading="loading" class="permission-container">
      <div class="permission-actions">
        <ElButton @click="toggleExpandAll">
          {{ isExpandAll ? '全部收起' : '全部展开' }}
        </ElButton>
        <ElButton @click="toggleSelectAll">
          {{ isSelectAll ? '取消全选' : '全部选择' }}
        </ElButton>
      </div>

      <ElScrollbar height="400px" class="permission-tree-container">
        <ElTree
          ref="treeRef"
          :data="permissionTree"
          show-checkbox
          node-key="key"
          :default-expand-all="isExpandAll"
          :default-checked-keys="selectedPermissions"
          :props="treeProps"
          @check="handleTreeCheck"
        >
        </ElTree>
      </ElScrollbar>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <ElButton @click="dialogVisible = false">取消</ElButton>
        <ElButton type="primary" :loading="submitLoading" @click="handleSubmit">
          保存
        </ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { ElMessage } from 'element-plus'
  import { AdminRoleService } from '@/api/adminRoleApi'

  interface Props {
    visible: boolean
    roleData?: Api.Role.RoleInfo
  }

  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'submit'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  const loading = ref(false)
  const submitLoading = ref(false)
  const treeRef = ref()
  const isExpandAll = ref(true)
  const isSelectAll = ref(false)

  // 权限树数据
  const permissionTree = ref<Api.Role.PermissionTreeNode[]>([])
  const selectedPermissions = ref<string[]>([])

  // 对话框显示控制
  const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
  })

  // 树组件属性配置
  const treeProps = {
    children: 'children',
    label: 'title'
  }

  // 加载权限树数据
  const loadPermissionTree = async () => {
    if (!props.roleData) return

    loading.value = true
    try {
      const response = await AdminRoleService.getPermissionTree(props.roleData.id)
      permissionTree.value = response || []
      
      // 提取已选中的权限
      selectedPermissions.value = extractSelectedPermissions(permissionTree.value)
      
      // 设置树的选中状态
      nextTick(() => {
        if (treeRef.value) {
          treeRef.value.setCheckedKeys(selectedPermissions.value)
        }
      })
    } catch (error) {
      console.error('获取权限树失败:', error)
      ElMessage.error('获取权限数据失败')
    } finally {
      loading.value = false
    }
  }

  // 提取已选中的权限
  const extractSelectedPermissions = (nodes: Api.Role.PermissionTreeNode[]): string[] => {
    const selected: string[] = []
    
    const traverse = (nodeList: Api.Role.PermissionTreeNode[]) => {
      nodeList.forEach(node => {
        if (node.selected) {
          selected.push(node.key)
        }
        if (node.children && node.children.length > 0) {
          traverse(node.children)
        }
      })
    }
    
    traverse(nodes)
    return selected
  }

  // 获取所有节点的key
  const getAllNodeKeys = (nodes: Api.Role.PermissionTreeNode[]): string[] => {
    const keys: string[] = []
    
    const traverse = (nodeList: Api.Role.PermissionTreeNode[]) => {
      nodeList.forEach(node => {
        keys.push(node.key)
        if (node.children && node.children.length > 0) {
          traverse(node.children)
        }
      })
    }
    
    traverse(nodes)
    return keys
  }

  // 切换展开/收起
  const toggleExpandAll = () => {
    const tree = treeRef.value
    if (!tree) return

    const nodes = tree.store.nodesMap
    for (const node in nodes) {
      nodes[node].expanded = !isExpandAll.value
    }

    isExpandAll.value = !isExpandAll.value
  }

  // 切换全选/取消全选
  const toggleSelectAll = () => {
    const tree = treeRef.value
    if (!tree) return

    if (!isSelectAll.value) {
      const allKeys = getAllNodeKeys(permissionTree.value)
      tree.setCheckedKeys(allKeys)
    } else {
      tree.setCheckedKeys([])
    }

    isSelectAll.value = !isSelectAll.value
  }

  // 处理树节点选中变化
  const handleTreeCheck = () => {
    const tree = treeRef.value
    if (!tree) return

    const checkedKeys = tree.getCheckedKeys()
    const allKeys = getAllNodeKeys(permissionTree.value)

    isSelectAll.value = checkedKeys.length === allKeys.length && allKeys.length > 0
  }

  // 提交权限配置
  const handleSubmit = async () => {
    if (!props.roleData || !treeRef.value) return

    submitLoading.value = true
    try {
      const checkedKeys = treeRef.value.getCheckedKeys()
      await AdminRoleService.updateRolePermissions(props.roleData.id, checkedKeys)
      
      ElMessage.success('权限配置成功')
      dialogVisible.value = false
      emit('submit')
    } catch (error) {
      console.error('权限配置失败:', error)
      ElMessage.error('权限配置失败请重试！')
    } finally {
      submitLoading.value = false
    }
  }

  // 监听对话框显示状态
  watch(
    () => props.visible,
    (visible) => {
      if (visible && props.roleData) {
        loadPermissionTree()
      }
    },
    { immediate: true }
  )
</script>

<style lang="scss" scoped>
  .permission-container {
    .permission-actions {
      margin-bottom: 16px;
      display: flex;
      gap: 8px;
    }

    .permission-tree-container {
      border: 1px solid var(--el-border-color);
      border-radius: 4px;
      padding: 8px;
    }

    .tree-node {
      display: flex;
      align-items: center;
      justify-content: space-between;
      width: 100%;
      
      .node-label {
        flex: 1;
      }
    }
  }

  .dialog-footer {
    text-align: right;
  }
</style>
