import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      authority: ['system:manage'],
      icon: 'lucide:settings',
      order: 10,
      title: '系统管理',
    },
    name: 'System',
    path: '/system',
    redirect: '/system/admin-users',
    children: [
      {
        name: 'AdminUsers',
        path: '/system/admin-users',
        component: () => import('#/views/system/admin-users/index.vue'),
        meta: {
          authority: ['system:admin_user:menu'],
          icon: 'lucide:users',
          title: '用户管理',
        },
      },
      {
        name: 'AdminRoles',
        path: '/system/admin-roles',
        component: () => import('#/views/system/admin-roles/index.vue'),
        meta: {
          authority: ['system:admin_role:menu'],
          icon: 'lucide:shield',
          title: '角色管理',
        },
      },
      // gen:routes
    ],
  },
];

export default routes;
