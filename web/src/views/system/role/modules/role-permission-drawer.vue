<template>
  <ElDrawer
    v-model="visible"
    title="配置权限"
    size="500px"
    @close="handleClose"
  >
    <div class="flex flex-col h-full overflow-hidden">
      <div class="mb-3 flex items-center gap-2 border-b border-[var(--art-card-border)] pb-3">
        <ElButton
          plain
          @click="handleSelectAll"
          :disabled="allLeafPermissionKeys.length === 0"
        >
          全选
        </ElButton>
        <ElButton
          plain
          @click="handleInvertSelection"
          :disabled="allLeafPermissionKeys.length === 0"
        >
          反选
        </ElButton>
      </div>

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
  import { ElMessage, type TreeInstance } from 'element-plus'

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

  const treeRef = ref<TreeInstance>()
  const permissionTree = ref<Api.SystemManage.Permission[]>([])

  /**
   * 加载权限树
   */
  const loadPermissionTree = async () => {
    try {
      const res = await fetchGetAllPermissions()
      permissionTree.value = res
    } catch (error) {
      console.error('获取权限树失败:', error)
    }
  }

  /**
   * 加载角色已有权限
   */
  const loadRolePermissions = async () => {
    // 未传角色时无需加载角色权限。
    if (!props.roleData?.id) return

    try {
      const res = await fetchGetRolePermissions(props.roleData.id)
      if (res && res.permissions) {
        // 保持叶子节点选中状态，避免父节点冗余回显。
        const leafKeys = filterRedundantKeys(res.permissions, permissionTree.value)
        treeRef.value?.setCheckedKeys(leafKeys)
      }
    } catch (error) {
      console.error('获取角色权限失败:', error)
    }
  }

  /**
   * 收集所有叶子节点权限键
   * @param tree 权限树
   * @returns 叶子节点键列表
   */
  const collectLeafKeys = (tree: Api.SystemManage.Permission[]): string[] => {
    const leafKeys: string[] = []

    tree.forEach((node) => {
      // 有子节点时继续递归收集。
      if (node.children?.length) {
        leafKeys.push(...collectLeafKeys(node.children))
        return
      }

      leafKeys.push(node.key)
    })

    return leafKeys
  }

  /**
   * 获取当前选中的叶子节点权限键
   * @returns 已选中的叶子节点键列表
   */
  const getCheckedLeafKeys = (): string[] => {
    const checkedKeys = (treeRef.value?.getCheckedKeys() || []) as string[]
    const leafKeySet = new Set(allLeafPermissionKeys.value)
    return checkedKeys.filter((key) => leafKeySet.has(String(key))).map((key) => String(key))
  }

  /**
   * 全选权限
   */
  const handleSelectAll = () => {
    treeRef.value?.setCheckedKeys(allLeafPermissionKeys.value)
  }

  /**
   * 反选权限
   */
  const handleInvertSelection = () => {
    const checkedLeafKeySet = new Set(getCheckedLeafKeys())
    const nextCheckedKeys = allLeafPermissionKeys.value.filter((key) => !checkedLeafKeySet.has(key))
    treeRef.value?.setCheckedKeys(nextCheckedKeys)
  }

  /**
   * 找到某个 key 的所有父节点 key
   * @param key 当前节点键
   * @param tree 权限树
   * @param path 当前路径
   * @returns 父节点键列表
   */
  const findParents = (
    key: string,
    tree: Api.SystemManage.Permission[],
    path: string[] = []
  ): string[] | null => {
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
   * 过滤掉冗余的父节点键值
   * @param perms 权限键列表
   * @param tree 权限树
   * @returns 去重后的叶子节点权限键
   */
  const filterRedundantKeys = (perms: string[], tree: Api.SystemManage.Permission[]): string[] => {
    const filtered = new Set<string>()
    perms.forEach((p) => {
      const hasChild = perms.some((other) => {
        // 当前节点自身不参与比较。
        if (other === p) return false
        const parents = findParents(other, tree)
        return parents?.includes(p)
      })
      // 仅保留最终叶子节点。
      if (!hasChild) filtered.add(p)
    })
    return Array.from(filtered)
  }

  /**
   * 展开权限，包含所有父节点
   * @param perms 权限键列表
   * @param tree 权限树
   * @returns 包含父节点的完整权限键
   */
  const expandPerms = (perms: string[], tree: Api.SystemManage.Permission[]): string[] => {
    const expanded = new Set<string>()
    perms.forEach((p) => {
      expanded.add(p)
      findParents(p, tree)?.forEach((parent) => expanded.add(parent))
    })
    return Array.from(expanded)
  }

  /**
   * 控制抽屉显隐
   */
  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
  })

  /**
   * 全部叶子节点权限键
   */
  const allLeafPermissionKeys = computed(() => {
    return collectLeafKeys(permissionTree.value)
  })

  const defaultProps = {
    children: 'children',
    label: 'label'
  }

  /**
   * 监听抽屉打开状态并初始化数据
   */
  watch(
    () => props.modelValue,
    async (newVal) => {
      // 仅在抽屉打开时加载数据。
      if (!newVal) return

      // 首次打开时才加载权限树。
      if (permissionTree.value.length === 0) {
        await loadPermissionTree()
      }

      // 有角色数据时才回填已选权限。
      if (props.roleData) {
        await loadRolePermissions()
      }
    }
  )

  /**
   * 关闭抽屉
   */
  const handleClose = () => {
    visible.value = false
    treeRef.value?.setCheckedKeys([])
  }

  /**
   * 保存权限配置
   */
  const savePermission = async () => {
    // 无角色时不允许保存。
    if (!props.roleData?.id) return

    try {
      const checkedKeys = getCheckedLeafKeys()
      // 保存时补齐所有父节点，保持后端权限结构完整。
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
