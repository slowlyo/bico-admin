import type { AuthDirective, HighlightDirective, RippleDirective } from '@/directives'

declare module 'vue' {
  export interface GlobalDirectives {
    vAuth: AuthDirective
    vHighlight: HighlightDirective
    vRipple: RippleDirective
  }
}
