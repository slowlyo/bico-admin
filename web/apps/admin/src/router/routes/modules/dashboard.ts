import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    name: 'Dashboard',
    path: '/dashboard',
    component: () => import('#/views/dashboard/index.vue'),
    meta: {
      affixTab: true,
      authority: ['dashboard:menu'],
      icon: 'lucide:layout-dashboard',
      order: -1,
      title: '工作台',
    },
  },
];

export default routes;
