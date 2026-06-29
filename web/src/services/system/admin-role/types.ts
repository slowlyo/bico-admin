/**
 * 角色管理 - 特殊接口类型
 */

/** 权限树节点 */
export interface Permission {
  key: string;
  label: string;
  children?: Permission[];
}

/** 角色（用于下拉选择） */
export interface AdminRole {
  id: number;
  name: string;
  code: string;
}

/** 更新权限参数 */
export interface UpdatePermissionsParams {
  permissions: string[];
}
