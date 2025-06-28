export default (initialState: any) => {
  // 在这里按照初始化数据定义项目中的权限，统一管理
  // 参考文档 https://umijs.org/docs/max/access
  const { currentUser } = initialState ?? {};

  return {
    // 是否可以访问管理员功能
    canSeeAdmin: currentUser && currentUser.role === 'admin',
    // 是否可以管理用户
    canManageUsers: currentUser && ['admin', 'manager'].includes(currentUser.role),
    // 是否已登录
    isLogin: !!currentUser,
  };
};
