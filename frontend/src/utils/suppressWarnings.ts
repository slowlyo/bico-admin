/**
 * 屏蔽第三方库的已知警告
 * 这些警告通常是由于第三方库（如Ant Design）使用了已废弃的React API导致的
 * 在库更新之前，我们可以安全地屏蔽这些警告
 */

// 需要屏蔽的警告消息列表
const SUPPRESSED_WARNINGS = [
  'findDOMNode is deprecated',
  'componentWillReceiveProps has been renamed',
  'componentWillMount has been renamed',
  'componentWillUpdate has been renamed',
];

/**
 * 初始化警告屏蔽
 * 只在开发环境中生效
 */
export function initWarningSuppress() {
  if (process.env.NODE_ENV !== 'development') {
    return;
  }

  // 屏蔽 console.warn
  const originalConsoleWarn = console.warn;
  console.warn = (...args: any[]) => {
    const message = args[0];
    if (typeof message === 'string') {
      const shouldSuppress = SUPPRESSED_WARNINGS.some(warning => 
        message.includes(warning)
      );
      if (shouldSuppress) {
        return;
      }
    }
    originalConsoleWarn.apply(console, args);
  };

  // 屏蔽 console.error 中的特定React警告
  const originalConsoleError = console.error;
  console.error = (...args: any[]) => {
    const message = args[0];
    if (typeof message === 'string') {
      const shouldSuppress = SUPPRESSED_WARNINGS.some(warning => 
        message.includes(warning)
      );
      if (shouldSuppress) {
        return;
      }
    }
    originalConsoleError.apply(console, args);
  };

  console.log('🔇 开发环境警告屏蔽已启用');
}
