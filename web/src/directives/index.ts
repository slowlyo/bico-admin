import type { App } from 'vue'
import { setupAuthDirective } from './core/auth'
import { setupHighlightDirective } from './business/highlight'
import { setupRippleDirective } from './business/ripple'
export function setupGlobDirectives(app: App) {
  setupAuthDirective(app) // 权限指令
  setupHighlightDirective(app) // 高亮指令
  setupRippleDirective(app) // 水波纹指令
}
