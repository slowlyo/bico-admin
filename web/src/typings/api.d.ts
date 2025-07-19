/**
 * namespace: Api
 *
 * 通用接口类型定义
 * 在.vue文件使用会报错，需要在 eslint.config.mjs 中配置 globals: { Api: 'readonly' }
 *
 * 注意：具体业务相关的类型定义已移动到对应的API文件中：
 * - 角色相关类型：见 adminRoleApi.ts 中的 RoleTypes
 * - 管理员用户相关类型：见 adminUserApi.ts 中的 AdminUserTypes
 * - 用户和认证相关类型：见 usersApi.ts 中的 AuthTypes 和 UserTypes
 */
declare namespace Api {
  /** 基础类型 */
  namespace Http {
    /** 基础响应 */
    interface BaseResponse<T = any> {
      // 状态码
      code: number
      // 消息 (支持两种字段名)
      msg?: string
      message?: string
      // 数据
      data: T
    }
  }

  /** 通用类型 */
  namespace Common {
    /** 分页参数 */
    interface PaginatingParams {
      /** 当前页码 */
      page: number
      /** 每页条数 */
      page_size: number
      /** 总条数 */
      total: number
    }

    /** 通用搜索参数 */
    type PaginatingSearchParams = Pick<PaginatingParams, 'page' | 'page_size'>

    /** 启用状态 */
    type EnableStatus = '1' | '2'
  }
}
