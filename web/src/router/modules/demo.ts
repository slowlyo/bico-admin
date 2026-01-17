import { AppRouteRecord } from '@/types/router'
import { useSettingStore } from '@/store/modules/setting'

const demoRoutes: AppRouteRecord = {
  path: '/demo',
  name: 'Demo',
  component: '/index/index',
  redirect: '/demo/excel',
  meta: {
    title: '示例页面',
    icon: 'ri:flask-line',
    order: 100,
    // 根据 debug 开启状态显示菜单
    get hideInMenu() {
      const settingStore = useSettingStore()
      return !settingStore.appConfig.debug
    }
  },
  children: [
    {
      path: 'excel',
      name: 'DemoExcel',
      component: '/demo/excel/index',
      meta: {
        title: 'Excel 导入导出',
        icon: 'ri:file-excel-2-line'
      }
    },
    {
      path: 'editor',
      name: 'DemoEditor',
      component: '/demo/editor/index',
      meta: {
        title: '富文本编辑器',
        icon: 'ri:edit-2-line'
      }
    }
  ]
}

export default demoRoutes
