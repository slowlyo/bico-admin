import type { App } from 'vue'
import { setupAuthDirective, type AuthDirective } from './core/auth'
import { setupHighlightDirective, type HighlightDirective } from './business/highlight'
import { setupRippleDirective, type RippleDirective } from './business/ripple'

export function setupGlobDirectives(app: App) {
  setupAuthDirective(app) // 权限指令
  setupHighlightDirective(app) // 高亮指令
  setupRippleDirective(app) // 水波纹指令
}

export type { AuthDirective, HighlightDirective, RippleDirective }
