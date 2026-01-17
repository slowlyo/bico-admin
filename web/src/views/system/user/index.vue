<!-- 用户管理页面 -->
<!-- art-full-height 自动计算出页面剩余高度 -->
<!-- art-table-card 一个符合系统样式的 class，同时自动撑满剩余高度 -->
<!-- 更多 useTable 使用示例请移步至 功能示例 下面的高级表格示例或者查看官方文档 -->
<!-- useTable 文档：https://www.artd.pro/docs/zh/guide/hooks/use-table.html -->
<template>
  <div class="user-page art-full-height">
    <!-- 搜索栏 -->
    <UserSearch v-model="searchForm" @search="handleSearch" @reset="resetSearchParams"></UserSearch>

    <ElCard class="art-table-card" shadow="never">
      <!-- 表格头部 -->
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElSpace wrap>
            <ElButton type="primary" @click="showDialog('add')" v-ripple v-auth="'system:admin_user:create'">新增用户</ElButton>
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

      <!-- 用户弹窗 -->
      <UserDialog
        v-model:visible="dialogVisible"
        :type="dialogType"
        :user-data="currentUserData"
        @submit="handleDialogSubmit"
      />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { ACCOUNT_TABLE_DATA } from '@/mock/temp/formData'
  import { useTable } from '@/hooks/core/useTable'
  import { fetchGetUserList } from '@/api/system-manage'
  import UserSearch from './modules/user-search.vue'
  import UserDialog from './modules/user-dialog.vue'
  import { ElTag, ElMessageBox, ElImage } from 'element-plus'
  import { DialogType } from '@/types'
  import { useAuth } from '@/hooks/core/useAuth'

  defineOptions({ name: 'User' })

  type UserListItem = Api.SystemManage.UserListItem

  const { hasAuth } = useAuth()

  // 弹窗相关
  const dialogType = ref<DialogType>('add')
  const dialogVisible = ref(false)
  const currentUserData = ref<Partial<UserListItem>>({})

  // 选中行
  const selectedRows = ref<UserListItem[]>([])

  // 搜索表单
  const searchForm = ref({
    username: undefined,
    name: undefined,
    enabled: undefined
  })

  // 用户状态配置
  const USER_STATUS_CONFIG = {
    true: { type: 'success' as const, text: '启用' },
    false: { type: 'danger' as const, text: '禁用' }
  } as const

  /**
   * 获取用户状态配置
   */
  const getUserStatusConfig = (enabled: boolean) => {
    return (
      USER_STATUS_CONFIG[String(enabled) as keyof typeof USER_STATUS_CONFIG] || {
        type: 'info' as const,
        text: '未知'
      }
    )
  }

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
      apiFn: fetchGetUserList,
      apiParams: {
        page: 1,
        pageSize: 20,
        ...searchForm.value
      },
      // 自定义分页字段映射，未设置时将使用全局配置 tableConfig.ts 中的 paginationKey
      // paginationKey: {
      //   current: 'pageNum',
      //   size: 'pageSize'
      // },
      columnsFactory: () => [
        { type: 'selection' }, // 勾选列
        { prop: 'id', label: 'ID', width: 80, align: 'center' },
        {
          prop: 'avatar',
          label: '头像',
          width: 80,
          align: 'center',
          formatter: (row: any) => {
            return h(ElImage, {
              class: 'size-9.5 rounded-md vertical-middle',
              src: row.avatar || '/src/assets/imgs/user/avatar.webp',
              previewSrcList: [row.avatar || '/src/assets/imgs/user/avatar.webp'],
              previewTeleported: true
            })
          }
        },
        {
          prop: 'username',
          label: '用户名',
          minWidth: 120
        },
        {
          prop: 'name',
          label: '姓名',
          minWidth: 120
        },
        {
          prop: 'roles',
          label: '角色',
          minWidth: 180,
          formatter: (row: any) => {
            return h(
              ElSpace,
              { size: 4, wrap: true },
              () =>
                row.roles?.map((role: any) =>
                  h(ElTag, { size: 'small', effect: 'plain' }, () => role.name)
                ) || []
            )
          }
        },
        {
          prop: 'enabled',
          label: '状态',
          width: 100,
          align: 'center',
          formatter: (row: any) => {
            const statusConfig = getUserStatusConfig(row.enabled)
            return h(ElTag, { type: statusConfig.type }, () => statusConfig.text)
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
          width: 120,
          fixed: 'right',
          align: 'center',
          auth: ['system:admin_user:edit', 'system:admin_user:delete'],
          formatter: (row: any) => {
            const buttons = []
            if (hasAuth('system:admin_user:edit')) {
              buttons.push(
                h(ArtButtonTable, {
                  type: 'edit',
                  onClick: () => showDialog('edit', row)
                })
              )
            }
            if (hasAuth('system:admin_user:delete')) {
              buttons.push(
                h(ArtButtonTable, {
                  type: 'delete',
                  onClick: () => deleteUser(row)
                })
              )
            }
            return h('div', buttons)
          }
        }
      ]
    },
    // 数据处理
    transform: {
      // 数据转换器
      dataTransformer: (list) => {
        if (!Array.isArray(list)) {
          return []
        }
        return list
      }
    }
  })

  /**
   * 搜索处理
   * @param params 参数
   */
  const handleSearch = (params: Record<string, any>) => {
    console.log(params)
    // 搜索参数赋值
    Object.assign(searchParams, params)
    getData()
  }

  /**
   * 显示用户弹窗
   */
  const showDialog = (type: DialogType, row?: UserListItem): void => {
    console.log('打开弹窗:', { type, row })
    dialogType.value = type
    currentUserData.value = row || {}
    nextTick(() => {
      dialogVisible.value = true
    })
  }

  /**
   * 删除用户
   */
  const deleteUser = (row: UserListItem): void => {
    ElMessageBox.confirm(`确定要删除用户 "${row.username}" 吗？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        const { fetchDeleteUser } = await import('@/api/system-manage')
        await fetchDeleteUser(row.id)
        ElMessage.success('删除成功')
        refreshData()
      } catch (error) {
        console.error('删除失败:', error)
      }
    })
  }

  /**
   * 处理弹窗提交事件
   */
  const handleDialogSubmit = async () => {
    try {
      dialogVisible.value = false
      currentUserData.value = {}
      refreshData()
    } catch (error) {
      console.error('提交失败:', error)
    }
  }

  /**
   * 处理表格行选择变化
   */
  const handleSelectionChange = (selection: UserListItem[]): void => {
    selectedRows.value = selection
    console.log('选中行数据:', selectedRows.value)
  }
</script>
