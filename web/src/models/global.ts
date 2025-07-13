// 全局共享数据
import { useState, useCallback } from 'react';
import { UserInfo } from '@/services/auth';

export interface GlobalState {
  currentUser?: UserInfo;
  permissions?: string[]; // 权限代码字符串数组
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
  const setPermissions = useCallback((permissions: string[]) => {
    setGlobalState(prev => ({
      ...prev,
      permissions,
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
  }, []);

  // 检查权限
  const hasPermission = useCallback((permission: string) => {
    if (!globalState.permissions) return false;
    return globalState.permissions.includes(permission);
  }, [globalState.permissions]);

  return {
    ...globalState,
    setCurrentUser,
    setPermissions,
    logout,
    hasPermission,
  };
};

export default useGlobal;
