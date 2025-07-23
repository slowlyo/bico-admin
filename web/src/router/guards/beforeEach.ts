import type { Router, RouteLocationNormalized, NavigationGuardNext } from 'vue-router'
import { ref, nextTick } from 'vue'
import NProgress from 'nprogress'
import { useSettingStore } from '@/store/modules/setting'
import { useUserStore } from '@/store/modules/user'
import { useMenuStore } from '@/store/modules/menu'
import { setWorktab } from '@/utils/navigation'
import { setPageTitle, setSystemTheme } from '../utils/utils'
import { menuService } from '@/api/menuApi'
import { registerDynamicRoutes } from '../utils/registerRoutes'
import { AppRouteRecord } from '@/types/router'
import { RoutesAlias } from '../routesAlias'
import { menuDataToRouter } from '../utils/menuToRouter'
import { asyncRoutes } from '../routes/asyncRoutes'
import { loadingService } from '@/utils/ui'
import { useCommon } from '@/composables/useCommon'
import { useWorktabStore } from '@/store/modules/worktab'
import { filterMenuItems, hasRoutePermission, debugMenuPermissions } from '@/utils/menuPermission'

// 前端权限模式 loading 关闭延时，提升用户体验
const LOADING_DELAY = 300

// 是否已注册动态路由
const isRouteRegistered = ref(false)

// 跟踪是否需要关闭 loading
const pendingLoading = ref(false)

/**
 * 设置路由全局前置守卫
 */
export function setupBeforeEachGuard(router: Router): void {
  router.beforeEach(
    async (
      to: RouteLocationNormalized,
      from: RouteLocationNormalized,
      next: NavigationGuardNext
    ) => {
      try {
        await handleRouteGuard(to, from, next, router)
      } catch (error) {
        console.error('路由守卫处理失败:', error)
        next('/exception/500')
      }
    }
  )

  // 设置后置守卫以关闭 loading 和进度条
  setupAfterEachGuard(router)
}

/**
 * 设置路由全局后置守卫
 */
function setupAfterEachGuard(router: Router): void {
  router.afterEach(() => {
    // 关闭进度条
    const settingStore = useSettingStore()
    if (settingStore.showNprogress) {
      NProgress.done()
    }

    // 关闭 loading 效果
    if (pendingLoading.value) {
      nextTick(() => {
        loadingService.hideLoading()
        pendingLoading.value = false
      })
    }
  })
}

/**
 * 处理路由守卫逻辑
 */
async function handleRouteGuard(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext,
  router: Router
): Promise<void> {
  const settingStore = useSettingStore()
  const userStore = useUserStore()

  // 处理进度条
  if (settingStore.showNprogress) {
    NProgress.start()
  }

  // 设置系统主题
  setSystemTheme(to)

  // 处理登录状态
  if (!(await handleLoginStatus(to, userStore, next))) {
    return
  }

  // 处理动态路由注册
  if (!isRouteRegistered.value && userStore.isLogin) {
    await handleDynamicRoutes(to, router, next)
    return
  }

  // 处理根路径跳转到首页
  if (userStore.isLogin && isRouteRegistered.value && handleRootPathRedirect(to, next)) {
    return
  }

  // 处理已知的匹配路由
  if (to.matched.length > 0) {
    // 检查路由权限
    if (!checkRoutePermission(to, userStore)) {
      next(RoutesAlias.Exception403)
      return
    }

    setWorktab(to)
    setPageTitle(to)
    next()
    return
  }

  // 尝试刷新路由重新注册
  if (userStore.isLogin && !isRouteRegistered.value) {
    await handleDynamicRoutes(to, router, next)
    return
  }

  // 未匹配到路由，跳转到 404
  next(RoutesAlias.Exception404)
}

/**
 * 处理登录状态
 */
async function handleLoginStatus(
  to: RouteLocationNormalized,
  userStore: ReturnType<typeof useUserStore>,
  next: NavigationGuardNext
): Promise<boolean> {
  // 如果用户未登录且访问需要登录的页面，跳转到登录页
  if (!userStore.isLogin && to.path !== RoutesAlias.Login && !to.meta.noLogin) {
    await userStore.logOut()
    next(RoutesAlias.Login)
    return false
  }

  // 如果用户已登录且访问登录页，重定向到首页
  if (userStore.isLogin && to.path === RoutesAlias.Login) {
    const { homePath } = useCommon()
    const redirectPath = homePath.value || '/'
    next({ path: redirectPath, replace: true })
    return false
  }

  return true
}

/**
 * 处理动态路由注册
 */
async function handleDynamicRoutes(
  to: RouteLocationNormalized,
  router: Router,
  next: NavigationGuardNext
): Promise<void> {
  try {
    // 显示 loading 并标记 pending
    pendingLoading.value = true
    loadingService.showLoading()

    await getMenuData(router)

    // 处理根路径跳转
    if (handleRootPathRedirect(to, next)) {
      return
    }

    next({
      path: to.path,
      query: to.query,
      hash: to.hash,
      replace: true
    })
  } catch (error) {
    console.error('动态路由注册失败:', error)
    next('/exception/500')
  }
}

/**
 * 获取菜单数据
 */
async function getMenuData(router: Router): Promise<void> {
  try {
    if (useCommon().isFrontendMode.value) {
      await processFrontendMenu(router)
    } else {
      await processBackendMenu(router)
    }
  } catch (error) {
    await handleMenuError(error)
    throw error
  }
}

/**
 * 处理前端控制模式的菜单逻辑
 */
async function processFrontendMenu(router: Router): Promise<void> {
  const userStore = useUserStore()
  const userPermissions = userStore.permissions

  // 将路由配置转换为菜单数据
  const allMenus = asyncRoutes.map((route) => menuDataToRouter(route))

  // 根据用户权限过滤菜单
  const filteredMenus = filterMenuItems(allMenus, userPermissions)

  // 开发环境下输出调试信息
  if (import.meta.env.DEV) {
    debugMenuPermissions(allMenus, userPermissions)
  }

  // 检查过滤后是否还有可访问的菜单
  if (filteredMenus.length === 0) {
    throw new Error('您没有权限访问任何菜单，请联系管理员')
  }

  // 添加延时以提升用户体验
  await new Promise((resolve) => setTimeout(resolve, LOADING_DELAY))

  await registerAndStoreMenu(router, filteredMenus)
}

/**
 * 处理后端控制模式的菜单逻辑
 */
async function processBackendMenu(router: Router): Promise<void> {
  const { menuList } = await menuService.getMenuList()
  await registerAndStoreMenu(router, menuList)
}

/**
 * 注册路由并存储菜单数据
 */
async function registerAndStoreMenu(router: Router, menuList: AppRouteRecord[]): Promise<void> {
  if (!isValidMenuList(menuList)) {
    throw new Error('获取菜单列表失败，请重新登录')
  }

  const menuStore = useMenuStore()
  menuStore.setMenuList(menuList)
  registerDynamicRoutes(router, menuList)
  isRouteRegistered.value = true
  useWorktabStore().validateWorktabs(router)
}

/**
 * 处理菜单相关错误
 */
async function handleMenuError(error: unknown): Promise<void> {
  console.error('菜单处理失败:', error)
  await useUserStore().logOut()
  throw error instanceof Error ? error : new Error('获取菜单列表失败，请重新登录')
}



/**
 * 验证菜单列表是否有效
 */
function isValidMenuList(menuList: AppRouteRecord[]): boolean {
  return Array.isArray(menuList) && menuList.length > 0
}

/**
 * 重置路由相关状态
 */
export function resetRouterState(): void {
  isRouteRegistered.value = false
  const menuStore = useMenuStore()
  menuStore.removeAllDynamicRoutes()
  menuStore.setMenuList([])
}

/**
 * 检查路由权限
 */
function checkRoutePermission(
  to: RouteLocationNormalized,
  userStore: ReturnType<typeof useUserStore>
): boolean {
  // 如果路由没有权限要求，允许访问
  const permissions = to.meta?.permissions as string[] | undefined
  if (!permissions || permissions.length === 0) {
    return true
  }

  // 检查用户是否有访问权限
  const userPermissions = userStore.permissions
  return permissions.some((permission: string) =>
    userPermissions.includes(permission)
  )
}

/**
 * 处理根路径跳转到首页
 */
function handleRootPathRedirect(to: RouteLocationNormalized, next: NavigationGuardNext): boolean {
  if (to.path === '/') {
    const { homePath } = useCommon()
    if (homePath.value) {
      next({ path: homePath.value, replace: true })
      return true
    }
  }
  return false
}
