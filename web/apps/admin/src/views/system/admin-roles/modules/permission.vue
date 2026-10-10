<script lang="ts" setup>
import type { AdminRole, PermissionNode } from '#/api/system/role';

import { ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';

import { Tree, message } from 'antdv-next';

import {
  listPermissionsApi,
  updateRolePermissionsApi,
} from '#/api/system/role';

const emit = defineEmits<{ success: [] }>();

const treeData = ref<Array<Record<string, unknown>>>([]);
const permissions = ref<PermissionNode[]>([]);
const selectedKeys = ref<string[]>([]);

// 权限树换成 Tree 需要的 title/key。
function toTree(nodes: PermissionNode[]): Array<Record<string, unknown>> {
  return nodes.map((node) => ({
    children: node.children ? toTree(node.children) : undefined,
    key: node.key,
    title: node.label,
  }));
}

// 从某个权限往上找父节点。找不到说明它不在当前树里。
function findParents(
  key: string,
  tree: PermissionNode[],
  path: string[] = [],
): null | string[] {
  for (const node of tree) {
    if (node.key === key) {
      return path;
    }
    if (node.children) {
      const result = findParents(key, node.children, [...path, node.key]);
      if (result) {
        return result;
      }
    }
  }
  return null;
}

// 勾选展示只保留叶子，父节点交给 Tree 计算半选。
function filterRedundant(perms: string[], tree: PermissionNode[]) {
  return perms.filter((perm) => {
    const hasChild = perms.some((other) => {
      if (other === perm) {
        return false;
      }
      return findParents(other, tree)?.includes(perm);
    });
    return !hasChild;
  });
}

// 保存时把叶子的父节点补回去，菜单权限依赖父级 key。
function expandPerms(perms: string[], tree: PermissionNode[]) {
  const expanded = new Set<string>();
  perms.forEach((perm) => {
    expanded.add(perm);
    findParents(perm, tree)?.forEach((parent) => expanded.add(parent));
  });
  return [...expanded];
}

// Tree 父子关联时可能返回对象，只取已勾选的 key。
function onCheck(keys: Array<number | string> | { checked: Array<number | string> }) {
  const checked = Array.isArray(keys) ? keys : keys.checked;
  selectedKeys.value = checked.map(String);
}

const [Drawer, drawerApi] = useVbenDrawer({
  async onConfirm() {
    const role = drawerApi.getData<AdminRole>();
    if (!role?.id) {
      return;
    }
    drawerApi.lock();
    try {
      await updateRolePermissionsApi(
        role.id,
        expandPerms(selectedKeys.value, permissions.value),
      );
      message.success('权限配置成功');
      emit('success');
      drawerApi.close();
    } finally {
      drawerApi.unlock();
    }
  },
  async onOpenChange(isOpen) {
    if (!isOpen) {
      return;
    }
    const role = drawerApi.getData<AdminRole>();
    if (permissions.value.length === 0) {
      permissions.value = (await listPermissionsApi()) ?? [];
      treeData.value = toTree(permissions.value);
    }
    selectedKeys.value = filterRedundant(
      role?.permissions ?? [],
      permissions.value,
    );
  },
  title: '配置权限',
});
</script>

<template>
  <Drawer class="w-[480px]">
    <Tree
      v-if="treeData.length > 0"
      checkable
      default-expand-all
      :checked-keys="selectedKeys"
      :tree-data="treeData"
      @check="onCheck"
    />
  </Drawer>
</template>
