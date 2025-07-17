<!-- 管理员角色管理 -->
<!-- art-full-height 自动计算出页面剩余高度 -->
<!-- art-table-card 一个符合系统样式的 class，同时自动撑满剩余高度 -->
<template>
  <div class="admin-role-page art-full-height">
    <!-- 搜索栏 -->
    <AdminRoleSearch
      @reset="resetSearchParams"
      @search="handleSearch"
    />

    <ElCard class="art-table-card" shadow="never">
      <!-- 表格头部 -->
      <ArtTableHeader v-model:columns="columnChecks" @refresh="refresh">
        <template #left>
          <ElButton
            v-if="hasAuth('create')"
            type="primary"
            @click="showDialog('add')"
          >
            新建角色
          </ElButton>
        </template>
      </ArtTableHeader>

      <!-- 表格 -->
      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        :table-config="{ rowKey: 'id' }"
        :layout="{ marginTop: 10 }"
        @row:selection-change="handleSelectionChange"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
        @sort:change="handleSortChange"
      >
      </ArtTable>

      <!-- 管理员角色弹窗 -->
      <AdminRoleDialog
        v-model:visible="dialogVisible"
        :type="dialogType"
        :role-data="currentRoleData"
        @submit="handleDialogSubmit"
      />

      <!-- 权限配置弹窗 -->
      <AdminRolePermissionDialog
        v-model:visible="permissionDialogVisible"
        :role-data="currentRoleData"
        @submit="handlePermissionSubmit"
      />
    </ElCard>
  </div>
</template>



<script setup lang="ts">
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { useTable } from '@/composables/useTable'
  import { useAuth } from '@/composables/useAuth'
  import { AdminRoleService } from '@/api/adminRoleApi'
  import AdminRoleSearch from './modules/admin-role-search.vue'
  import AdminRoleDialog from './modules/admin-role-dialog.vue'
  import AdminRolePermissionDialog from './modules/admin-role-permission-dialog.vue'
  import { ElMessage, ElMessageBox, ElSwitch } from 'element-plus'
  import type { ColumnOption } from '@/types/component'

  defineOptions({ name: 'AdminRole' })

  // 权限检查
  const { hasAuth, hasAnyAuth } = useAuth('system.role')

  type RoleInfo = Api.Role.RoleInfo
  const { getRoleList } = AdminRoleService

  // 弹窗相关
  const dialogType = ref<Form.DialogType>('add')
  const dialogVisible = ref(false)
  const permissionDialogVisible = ref(false)
  const currentRoleData = ref<RoleInfo | undefined>(undefined)

  // 选中行
  const selectedRows = ref<RoleInfo[]>([])

  const {
    columns,
    columnChecks,
    tableData: data,
    isLoading: loading,
    paginationState: pagination,
    searchState: searchParams,
    searchData: getDataByPage,
    resetSearch: resetSearchParams,
    onPageSizeChange: handleSizeChange,
    onCurrentPageChange: handleCurrentChange,
    refreshAll: refresh,
    refreshAfterCreate: refreshAfterAdd,
    refreshAfterUpdate: refreshAfterEdit,
    refreshAfterRemove: refreshAfterDelete
  } = useTable<RoleInfo>({
    // 核心配置
    core: {
      apiFn: getRoleList,
      apiParams: {
        page: 1,
        page_size: 20,
        name: '',
        code: '',
        status: undefined,
        sort_by: '',
        sort_desc: false
      },
      // 配置分页字段映射
      paginationKey: {
        current: 'page',
        size: 'page_size'
      },
      columnsFactory: () => {
        const columns: ColumnOption<RoleInfo>[] = [
          {
            prop: 'id',
            label: 'ID',
            width: 80
          },
          {
            prop: 'name',
            label: '角色名称',
            minWidth: 120
          },
          {
            prop: 'code',
            label: '角色代码',
            minWidth: 120
          },
          {
            prop: 'description',
            label: '描述',
            minWidth: 200
          },
          {
            prop: 'status',
            label: '状态',
            width: 100,
            formatter: (row: RoleInfo) => {
              return h(ElSwitch, {
                modelValue: row.status === 1,
                disabled: !hasAuth('update') || !row.can_edit || statusUpdateLoading.value.has(row.id),
                loading: statusUpdateLoading.value.has(row.id),
                checkedText: '启用',
                uncheckedText: '禁用',
                onChange: (val: string | number | boolean) => handleStatusChange(row, val ? 1 : 0)
              })
            }
          },
          {
            prop: 'user_count',
            label: '用户数量',
            width: 100
          },
          {
            prop: 'created_at',
            label: '创建时间',
            minWidth: 160,
            sortable: true
          },
          {
            prop: 'updated_at',
            label: '更新时间',
            minWidth: 160,
            sortable: true
          }
        ]

        // 只有在有编辑或删除权限时才添加操作列
        if (hasAnyAuth(['update', 'delete'])) {
          columns.push({
            prop: 'operation',
            label: '操作',
            width: 180,
            fixed: 'right',
            formatter: (row: RoleInfo) => {
              const buttons = []

              // 编辑按钮
              if (hasAuth('update')) {
                buttons.push(h(ArtButtonTable, {
                  type: 'edit',
                  disabled: !row.can_edit,
                  onClick: () => showDialog('edit', row)
                }))
              }

              // 权限配置按钮
              if (hasAuth('update')) {
                buttons.push(h(ArtButtonTable, {
                  type: 'view',
                  icon: '&#xe755;',
                  text: '权限',
                  disabled: !row.can_edit,
                  onClick: () => showPermissionDialog(row)
                }))
              }

              // 删除按钮
              if (hasAuth('delete')) {
                buttons.push(h(ArtButtonTable, {
                  type: 'delete',
                  disabled: !row.can_delete,
                  onClick: () => deleteRole(row)
                }))
              }

              return h('div', { class: 'flex gap-2' }, buttons)
            }
          })
        }

        return columns
      }
    },
    // 数据处理
    transform: {
      // 响应数据适配器
      responseAdapter: (response: any) => {
        if (response && response.list) {
          return {
            records: response.list || [],
            total: response.total || 0,
            current: response.page || 1,
            size: response.page_size || 20
          }
        }
        return { records: [], total: 0, current: 1, size: 20 }
      },
      // 数据转换器
      dataTransformer: (records: any) => {
        if (!Array.isArray(records)) {
          return []
        }
        return records
      }
    },
    // 性能优化
    performance: {
      enableCache: true,
      cacheTime: 10 * 60 * 1000
    },
  })

  /**
   * 显示角色弹窗
   */
  const showDialog = (type: Form.DialogType, row?: RoleInfo): void => {
    dialogType.value = type
    currentRoleData.value = row
    nextTick(() => {
      dialogVisible.value = true
    })
  }

  /**
   * 显示权限配置弹窗
   */
  const showPermissionDialog = (row: RoleInfo): void => {
    currentRoleData.value = row
    nextTick(() => {
      permissionDialogVisible.value = true
    })
  }

  // 状态更新防抖
  const statusUpdateLoading = ref<Set<number>>(new Set())

  /**
   * 更新角色状态
   */
  const handleStatusChange = async (row: RoleInfo, status: number) => {
    // 权限检查
    if (!hasAuth('update')) {
      ElMessage.error('权限不足，无法修改角色状态')
      return
    }

    if (statusUpdateLoading.value.has(row.id)) {
      return
    }

    statusUpdateLoading.value.add(row.id)

    try {
      const response = await AdminRoleService.updateRoleStatus(row.id, status)
      const successMessage = response?.msg || '状态更新成功'
      ElMessage.success(successMessage)
      refresh()
    } catch (error: any) {
      // HTTP 拦截器已经显示了具体的错误消息，这里不再重复显示
    } finally {
      statusUpdateLoading.value.delete(row.id)
    }
  }

  /**
   * 删除角色
   */
  const deleteRole = (row: RoleInfo): void => {
    ElMessageBox.confirm(`确定要删除角色 "${row.name}" 吗？`, '删除角色', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        await AdminRoleService.deleteRole(row.id)
        ElMessage.success('删除成功')
        refreshAfterDelete()
      } catch (error) {
        // HTTP 拦截器已经显示了具体的错误消息，这里不再重复显示
      }
    }).catch(() => {
      // 用户取消删除
    })
  }

  /**
   * 处理弹窗提交事件
   */
  const handleDialogSubmit = async () => {
    try {
      dialogVisible.value = false
      await (dialogType.value === 'add' ? refreshAfterAdd() : refreshAfterEdit())
      currentRoleData.value = undefined
    } catch (error) {
    }
  }

  /**
   * 处理权限配置提交事件
   */
  const handlePermissionSubmit = async () => {
    try {
      permissionDialogVisible.value = false
      refresh()
      currentRoleData.value = undefined
    } catch (error) {
    }
  }

  /**
   * 处理表格行选择变化
   */
  const handleSelectionChange = (selection: RoleInfo[]): void => {
    selectedRows.value = selection
  }

  /**
   * 处理搜索事件
   */
  const handleSearch = (searchFormData: any): void => {
    Object.assign(searchParams, searchFormData)
    getDataByPage()
  }

  /**
   * 处理表格排序变化
   */
  const handleSortChange = (sortInfo: any): void => {
    if (sortInfo.prop && sortInfo.order) {
      const sortDesc = sortInfo.order === 'descending'
      Object.assign(searchParams, {
        sort_by: sortInfo.prop,
        sort_desc: sortDesc
      })
      getDataByPage()
    } else {
      Object.assign(searchParams, {
        sort_by: '',
        sort_desc: false
      })
      getDataByPage()
    }
  }
</script>

<style lang="scss" scoped>
  .admin-role-page {
    .flex {
      display: flex;
    }

    .gap-2 {
      gap: 8px;
    }
  }
</style>
