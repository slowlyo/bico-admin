<!-- 角色管理页面 -->
<template>
  <div class="art-full-height">
    <RoleSearch
      v-model="searchForm"
      @search="handleSearch"
      @reset="resetSearchParams"
    ></RoleSearch>

    <ElCard
      class="art-table-card"
      shadow="never"
    >
      <ArtTableHeader
        v-model:columns="columnChecks"
        :loading="loading"
        @refresh="refreshData"
      >
        <template #left>
          <ElSpace wrap>
            <ElButton type="primary" @click="showDialog('add')" v-ripple v-auth="'system:admin_role:create'">新增角色</ElButton>
          </ElSpace>
        </template>
      </ArtTableHeader>

      <!-- 表格 -->
      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @selection-change="handleSelectionChange"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
        @sort-change="handleSortChange"
      >
      </ArtTable>
    </ElCard>

    <!-- 角色编辑弹窗 -->
    <RoleEditDialog
      v-model="dialogVisible"
      :dialog-type="dialogType"
      :role-data="currentRoleData"
      @success="refreshData"
    />

    <!-- 菜单权限抽屉 -->
    <RolePermissionDrawer
      v-model="permissionDrawer"
      :role-data="currentRoleData"
      @success="refreshData"
    />
  </div>
</template>

<script setup lang="ts">
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { useTable } from '@/hooks/core/useTable'
  import { fetchGetRoleList, fetchDeleteRole } from '@/api/system-manage'
  import RoleSearch from './modules/role-search.vue'
  import RoleEditDialog from './modules/role-edit-dialog.vue'
  import RolePermissionDrawer from './modules/role-permission-drawer.vue'
  import { ElTag, ElMessageBox, ElMessage } from 'element-plus'
  import { h } from 'vue'
  import { useAuth } from '@/hooks/core/useAuth'

  defineOptions({ name: 'Role' })

  type RoleListItem = Api.SystemManage.RoleListItem

  const { hasAuth } = useAuth()

  // 搜索表单
  const searchForm = ref({
    name: undefined,
    code: undefined,
    description: undefined,
    enabled: undefined
  })

  const dialogVisible = ref(false)
  const permissionDrawer = ref(false)
  const currentRoleData = ref<RoleListItem | undefined>(undefined)

  const {
    columns,
    columnChecks,
    data,
    loading,
    pagination,
    getData,
    searchParams,
    resetSearchParams,
    handleSizeChange,
    handleCurrentChange,
    handleSortChange,
    refreshData
  } = useTable({
    // 核心配置
    core: {
      apiFn: fetchGetRoleList,
      apiParams: {
        page: 1,
        pageSize: 20
      },
      // 排除 apiParams 中的属性
      excludeParams: ['daterange'],
      columnsFactory: () => [
        {
          prop: 'id',
          label: 'ID',
          width: 80,
          align: 'center'
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
          minWidth: 150,
          showOverflowTooltip: true
        },
        {
          prop: 'enabled',
          label: '状态',
          width: 100,
          align: 'center',
          formatter: (row: any) => {
            const statusConfig = row.enabled
              ? { type: 'success', text: '启用' }
              : { type: 'warning', text: '禁用' }
            return h(
              ElTag,
              { type: statusConfig.type as 'success' | 'warning' },
              () => statusConfig.text
            )
          }
        },
        {
          prop: 'created_at',
          label: '创建时间',
          width: 180,
          align: 'center',
          sortable: 'custom'
        },
        {
          prop: 'operation',
          label: '操作',
          width: 160,
          fixed: 'right',
          align: 'center',
          auth: [
            'system:admin_role:permission',
            'system:admin_role:edit',
            'system:admin_role:delete'
          ],
          formatter: (row: any) => {
            const buttons = []
            if (hasAuth('system:admin_role:permission')) {
              buttons.push(
                h(ArtButtonTable, {
                  type: 'view',
                  icon: 'ri:shield-keyhole-line',
                  onClick: () => showPermissionDrawer(row)
                })
              )
            }
            if (hasAuth('system:admin_role:edit')) {
              buttons.push(
                h(ArtButtonTable, {
                  type: 'edit',
                  onClick: () => showDialog('edit', row)
                })
              )
            }
            if (hasAuth('system:admin_role:delete')) {
              buttons.push(
                h(ArtButtonTable, {
                  type: 'delete',
                  onClick: () => deleteRole(row)
                })
              )
            }
            return h('div', buttons)
          }
        }
      ]
    }
  })

  /**
   * 处理表格行选择变化
   */
  const handleSelectionChange = (selection: RoleListItem[]) => {
    console.log('选中行数据:', selection)
  }

  const dialogType = ref<'add' | 'edit'>('add')

  const showDialog = (type: 'add' | 'edit', row?: RoleListItem) => {
    dialogVisible.value = true
    dialogType.value = type
    currentRoleData.value = row
  }

  /**
   * 搜索处理
   * @param params 搜索参数
   */
  const handleSearch = (params: Record<string, any>) => {
    // 搜索参数赋值
    Object.assign(searchParams, params)
    getData()
  }

  const showPermissionDrawer = (row?: RoleListItem) => {
    permissionDrawer.value = true
    currentRoleData.value = row
  }

  const deleteRole = (row: RoleListItem) => {
    ElMessageBox.confirm(`确定删除角色"${row.name}"吗？此操作不可恢复！`, '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
      .then(async () => {
        try {
          await fetchDeleteRole(row.id)
          ElMessage.success('删除成功')
          refreshData()
        } catch (error) {
          console.error('删除角色失败:', error)
        }
      })
      .catch(() => {
        ElMessage.info('已取消删除')
      })
  }
</script>
