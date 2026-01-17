<template>
  <ElDrawer
    v-model="visible"
    title="配置权限"
    size="500px"
    @close="handleClose"
  >
    <div class="flex flex-col h-full overflow-hidden">
      <div class="flex-1 overflow-auto">
        <ElTree
          ref="treeRef"
          :data="permissionTree"
          show-checkbox
          node-key="key"
          :default-expand-all="true"
          :props="defaultProps"
        >
          <template #default="{ data }">
            <span>{{ data.label }}</span>
          </template>
        </ElTree>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2">
        <ElButton @click="handleClose">取消</ElButton>
        <ElButton type="primary" @click="savePermission">保存</ElButton>
      </div>
    </template>
  </ElDrawer>
</template>

<script setup lang="ts">
  import {
    fetchGetAllPermissions,
    fetchGetRolePermissions,
    fetchUpdateRolePermissions
  } from '@/api/system-manage'
  import { ElMessage } from 'element-plus'

  type RoleListItem = Api.SystemManage.RoleListItem

  interface Props {
    modelValue: boolean
    roleData?: RoleListItem
  }

  interface Emits {
    (e: 'update:modelValue', value: boolean): void
    (e: 'success'): void
  }

  const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
    roleData: undefined
  })

  const emit = defineEmits<Emits>()

  const treeRef = ref()
  const permissionTree = ref<Api.SystemManage.Permission[]>([])

  const loadPermissionTree = async () => {
    try {
      const res = await fetchGetAllPermissions()
      permissionTree.value = res
    } catch (error) {
      console.error('获取权限树失败:', error)
    }
  }

  const loadRolePermissions = async () => {
    if (!props.roleData?.id) return
    try {
      const res = await fetchGetRolePermissions(props.roleData.id)
      if (res && res.permissions) {
        // web_old 逻辑: 过滤掉冗余的父节点，只给树组件设置叶子节点
        const leafKeys = filterRedundantKeys(res.permissions, permissionTree.value)
        treeRef.value?.setCheckedKeys(leafKeys)
      }
    } catch (error) {
      console.error('获取角色权限失败:', error)
    }
  }

  /**
   * 找到某个 key 的所有父节点 key
   */
  const findParents = (key: string, tree: Api.SystemManage.Permission[], path: string[] = []): string[] | null => {
    for (const node of tree) {
      if (node.key === key) return path
      if (node.children) {
        const result = findParents(key, node.children, [...path, node.key])
        if (result) return result
      }
    }
    return null
  }

  /**
   * 过滤掉冗余的父节点键值（web_old 逻辑）
   */
  const filterRedundantKeys = (perms: string[], tree: Api.SystemManage.Permission[]): string[] => {
    const filtered = new Set<string>()
    perms.forEach((p) => {
      const hasChild = perms.some((other) => {
        if (other === p) return false
        const parents = findParents(other, tree)
        return parents?.includes(p)
      })
      if (!hasChild) filtered.add(p)
    })
    return Array.from(filtered)
  }

  /**
   * 展开权限，包含所有父节点（web_old 逻辑）
   */
  const expandPerms = (perms: string[], tree: Api.SystemManage.Permission[]): string[] => {
    const expanded = new Set<string>()
    perms.forEach((p) => {
      expanded.add(p)
      findParents(p, tree)?.forEach((parent) => expanded.add(parent))
    })
    return Array.from(expanded)
  }

  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
  })

  const defaultProps = {
    children: 'children',
    label: 'label'
  }

  watch(
    () => props.modelValue,
    async (newVal) => {
      if (newVal) {
        if (permissionTree.value.length === 0) {
          await loadPermissionTree()
        }
        if (props.roleData) {
          await loadRolePermissions()
        }
      }
    }
  )

  const handleClose = () => {
    visible.value = false
    treeRef.value?.setCheckedKeys([])
  }

  const savePermission = async () => {
    if (!props.roleData?.id) return

    try {
      const checkedKeys = treeRef.value?.getCheckedKeys() || []
      // 保存时，复刻 web_old 逻辑：带上所有父节点
      const fullPermissions = expandPerms(checkedKeys, permissionTree.value)

      await fetchUpdateRolePermissions(props.roleData.id, fullPermissions)

      ElMessage.success('权限保存成功')
      emit('success')
      handleClose()
    } catch (error) {
      console.error('保存权限失败:', error)
    }
  }
</script>
