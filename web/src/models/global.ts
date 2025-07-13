// 全局共享数据
import { useState, useCallback } from 'react';
import { UserInfo, Permission, Menu } from '@/services/auth';

export interface GlobalState {
  currentUser?: UserInfo;
  permissions?: Permission[];
  menus?: Menu[];
  isLogin: boolean;
}

const useGlobal = () => {
  const [globalState, setGlobalState] = useState<GlobalState>({
    isLogin: false,
  });

  // 设置用户信息
  const setCurrentUser = useCallback((user: UserInfo | undefined) => {
    setGlobalState(prev => ({
      ...prev,
      currentUser: user,
      isLogin: !!user,
    }));
  }, []);

  // 设置权限
  const setPermissions = useCallback((permissions: Permission[]) => {
    setGlobalState(prev => ({
      ...prev,
      permissions,
    }));
  }, []);

  // 设置菜单
  const setMenus = useCallback((menus: Menu[]) => {
    setGlobalState(prev => ({
      ...prev,
      menus,
    }));
  }, []);

  // 登出
  const logout = useCallback(() => {
    setGlobalState({
      isLogin: false,
    });
    // 清除本地存储
    localStorage.removeItem('token');
    localStorage.removeItem('userInfo');
    localStorage.removeItem('permissions');
    localStorage.removeItem('menus');
  }, []);

  // 检查权限
  const hasPermission = useCallback((permission: string) => {
    if (!globalState.permissions) return false;
    return globalState.permissions.some(p => p.sign === permission);
  }, [globalState.permissions]);

  return {
    ...globalState,
    setCurrentUser,
    setPermissions,
    setMenus,
    logout,
    hasPermission,
  };
};

export default useGlobal;
