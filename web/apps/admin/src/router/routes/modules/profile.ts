import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

// 个人中心从头像菜单进入，不出现在侧边栏。
const routes: RouteRecordRaw[] = [
  {
    name: 'Profile',
    path: '/profile',
    component: () => import('#/views/_core/profile/index.vue'),
    meta: {
      hideInMenu: true,
      icon: 'lucide:user',
      title: $t('page.auth.profile'),
    },
  },
];

export default routes;
