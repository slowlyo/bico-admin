<!-- 管理员用户管理 -->
<!-- art-full-height 自动计算出页面剩余高度 -->
<!-- art-table-card 一个符合系统样式的 class，同时自动撑满剩余高度 -->
<template>
  <div class="admin-user-page art-full-height">
    <!-- 搜索栏 -->
    <AdminUserSearch
      @reset="resetSearchParams"
      @search="handleSearch"
    />

    <ElCard class="art-table-card" shadow="never">
      <!-- 表格头部 -->
      <ArtTableHeader v-model:columns="columnChecks" @refresh="refresh">
        <template #left>
          <ElButton type="primary" @click="showDialog('add')">新建管理员</ElButton>
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

      <!-- 管理员用户弹窗 -->
      <AdminUserDialog
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
  import { useTable } from '@/composables/useTable'
  import { AdminUserService, RoleService } from '@/api/adminUserApi'
  import AdminUserSearch from './modules/admin-user-search.vue'
  import AdminUserDialog from './modules/admin-user-dialog.vue'
  import { ElMessage, ElMessageBox, ElSwitch, ElTag } from 'element-plus'

  defineOptions({ name: 'AdminUser' })

  type AdminUserInfo = Api.AdminUser.AdminUserInfo
  const { getAdminUserList } = AdminUserService

  // 弹窗相关
  const dialogType = ref<Form.DialogType>('add')
  const dialogVisible = ref(false)
  const currentUserData = ref<AdminUserInfo | undefined>(undefined)

  // 选中行
  const selectedRows = ref<AdminUserInfo[]>([])

  // 角色列表
  const roleList = ref<Api.Role.RoleOption[]>([])





  /**
   * 加载角色列表
   */
  const loadRoles = async () => {
    try {
      const response = await RoleService.getActiveRoles()
      roleList.value = response || []
    } catch (error) {
      console.error('获取角色列表失败:', error)
    }
  }

  // 组件挂载时加载角色列表
  onMounted(() => {
    loadRoles()
  })

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
  } = useTable<AdminUserInfo>({
    // 核心配置
    core: {
      apiFn: getAdminUserList,
      apiParams: {
        page: 1,
        page_size: 20,
        username: '',
        name: '',
        status: undefined,
        role_id: undefined,
        sort_by: '',
        sort_desc: false
      },
      // 配置分页字段映射
      paginationKey: {
        current: 'page',
        size: 'page_size'
      },
      columnsFactory: () => [
        { type: 'selection' }, // 勾选列
        { type: 'index', width: 60, label: '序号' }, // 序号
        {
          prop: 'id',
          label: 'ID',
          width: 80
        },
        {
          prop: 'avatar',
          label: '头像',
          width: 80,
          formatter: (row) => {
            return h('img', {
              class: 'avatar',
              src: row.avatar || '/default-avatar.png',
              style: 'width: 40px; height: 40px; border-radius: 6px;'
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
          prop: 'email',
          label: '邮箱',
          minWidth: 180
        },
        {
          prop: 'phone',
          label: '手机号',
          minWidth: 120
        },
        {
          prop: 'status',
          label: '状态',
          width: 100,
          formatter: (row) => {
            return h(ElSwitch, {
              modelValue: row.status === 1,
              disabled: !row.can_disable || statusUpdateLoading.value.has(row.id),
              loading: statusUpdateLoading.value.has(row.id),
              checkedText: '启用',
              uncheckedText: '禁用',
              onChange: (val: string | number | boolean) => handleStatusChange(row, val ? 1 : 0)
            })
          }
        },
        {
          prop: 'roles',
          label: '角色',
          minWidth: 150,
          formatter: (row) => {
            return h('div', {
              class: 'flex flex-wrap gap-1',
              style: { padding: '4px 0' }
            },
              row.roles?.map((role: Api.AdminUser.AdminUserRole) =>
                h(ElTag, {
                  key: role.id,
                  type: 'primary',
                  size: 'small',
                  effect: 'light',
                  style: { margin: '2px' }
                }, () => role.name)
              ) || []
            )
          }
        },
        {
          prop: 'last_login_at',
          label: '最后登录时间',
          minWidth: 160,
          sortable: true
        },
        {
          prop: 'created_at',
          label: '创建时间',
          minWidth: 160,
          sortable: true
        },
        {
          prop: 'operation',
          label: '操作',
          width: 120,
          fixed: 'right', // 固定列
          formatter: (row) =>
            h('div', [
              h(ArtButtonTable, {
                type: 'edit',
                onClick: () => showDialog('edit', row)
              }),
              h(ArtButtonTable, {
                type: 'delete',
                disabled: !row.can_delete,
                onClick: () => deleteAdminUser(row)
              })
            ])
        }
      ]
    },
    // 数据处理
    transform: {
      // 响应数据适配器 - 适配后端API响应格式
      responseAdapter: (response: any) => {
        // HTTP工具已经提取了data，所以response直接是分页数据
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
      // 数据转换器 - 处理管理员用户数据
      dataTransformer: (records: any) => {
        // 类型守卫检查
        if (!Array.isArray(records)) {
          console.warn('数据转换器: 期望数组类型，实际收到:', typeof records)
          return []
        }

        // 处理管理员用户数据
        return records.map((item: any) => {
          return {
            ...item,
            avatar: item.avatar || '/default-avatar.png'
          }
        })
      }
    },
    // 性能优化
    performance: {
      enableCache: true, // 是否开启缓存
      cacheTime: 10 * 60 * 1000 // 缓存时间 10分钟
    },
    // 生命周期钩子
    hooks: {
      onError: (error) => ElMessage.error(error.message) // 错误处理
      // onSuccess: (data) => console.log('数据加载成功:', data), // 成功处理
      // onCacheHit: (data) => console.log('缓存命中:', data), // 缓存命中处理
      // resetFormCallback: () => console.log('重置表单')
    },
    // 调试配置
    debug: {
      enableLog: true // 是否开启日志
    }
  })

  /**
   * 显示管理员用户弹窗
   */
  const showDialog = (type: Form.DialogType, row?: AdminUserInfo): void => {
    console.log('打开弹窗:', { type, row })
    dialogType.value = type
    currentUserData.value = row
    nextTick(() => {
      dialogVisible.value = true
    })
  }

  // 状态更新防抖
  const statusUpdateLoading = ref<Set<number>>(new Set())

  /**
   * 更新管理员用户状态
   */
  const handleStatusChange = async (row: AdminUserInfo, status: number) => {
    // 防止重复请求
    if (statusUpdateLoading.value.has(row.id)) {
      return
    }

    statusUpdateLoading.value.add(row.id)

    try {
      // 获取完整响应，包括message
      const response = await AdminUserService.updateAdminUserStatus(row.id, status)
      const successMessage = response?.msg || '状态更新成功'
      ElMessage.success(successMessage)
      refresh() // 刷新表格数据
    } catch (error: any) {
      console.error('状态更新失败:', error)
      // 显示后端返回的具体错误消息
      const errorMessage = error?.message || error?.msg || '状态更新失败，请重试！'
      ElMessage.error(errorMessage)
    } finally {
      statusUpdateLoading.value.delete(row.id)
    }
  }

  /**
   * 删除管理员用户
   */
  const deleteAdminUser = (row: AdminUserInfo): void => {
    console.log('删除管理员用户:', row)
    ElMessageBox.confirm(`确定要删除管理员用户 "${row.name}" 吗？`, '删除管理员', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        await AdminUserService.deleteAdminUser(row.id)
        ElMessage.success('删除成功')
        refreshAfterDelete() // 智能删除后刷新
      } catch (error) {
        console.error('删除失败:', error)
        ElMessage.error('删除失败请重试！')
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
      currentUserData.value = undefined
    } catch (error) {
      console.error('提交失败:', error)
    }
  }

  /**
   * 处理表格行选择变化
   */
  const handleSelectionChange = (selection: AdminUserInfo[]): void => {
    selectedRows.value = selection
    console.log('选中行数据:', selectedRows.value)
  }

  /**
   * 处理搜索事件
   */
  const handleSearch = (searchFormData: any): void => {
    console.log('搜索参数:', searchFormData)

    // 更新搜索参数
    Object.assign(searchParams, searchFormData)

    // 重新获取数据
    getDataByPage()
  }

  /**
   * 处理表格排序变化
   */
  const handleSortChange = (sortInfo: any): void => {
    console.log('排序变化:', sortInfo)

    // 更新排序参数
    if (sortInfo.prop && sortInfo.order) {
      // 将 Element Plus 的排序方向转换为后端需要的格式
      const sortDesc = sortInfo.order === 'descending'

      // 更新搜索参数中的排序字段
      Object.assign(searchParams, {
        sort_by: sortInfo.prop,
        sort_desc: sortDesc
      })

      // 重新获取数据
      getDataByPage()
    } else {
      // 清除排序
      Object.assign(searchParams, {
        sort_by: '',
        sort_desc: false
      })

      // 重新获取数据
      getDataByPage()
    }
  }
</script>

<style lang="scss" scoped>
  .admin-user-page {
    :deep(.avatar) {
      width: 40px;
      height: 40px;
      border-radius: 6px;
      object-fit: cover;
    }
  }
</style>
