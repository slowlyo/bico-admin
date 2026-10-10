import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      authority: ['developer'],
      icon: 'lucide:wrench',
      order: 20,
      title: '开发者工具',
    },
    name: 'Developer',
    path: '/developer',
    redirect: '/developer/api-docs',
    children: [
      {
        name: 'ApiDocs',
        path: '/developer/api-docs',
        component: () => import('#/views/developer/api-docs/index.vue'),
        meta: {
          icon: 'lucide:book-open',
          title: '接口文档',
        },
      },
      {
        meta: {
          icon: 'lucide:flask-conical',
          title: '示例',
        },
        name: 'DeveloperDemo',
        path: '/developer/demo',
        redirect: '/developer/demo/excel',
        children: [
          {
            name: 'DemoExcel',
            path: '/developer/demo/excel',
            component: () => import('#/views/developer/excel/index.vue'),
            meta: {
              title: 'Excel 导入导出',
            },
          },
          {
            name: 'DemoEditor',
            path: '/developer/demo/editor',
            component: () => import('#/views/developer/editor/index.vue'),
            meta: {
              title: '编辑器示例',
            },
          },
        ],
      },
    ],
  },
];

export default routes;
