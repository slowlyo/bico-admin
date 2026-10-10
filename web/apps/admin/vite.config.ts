import { defineConfig } from '@vben/vite-config';

export default defineConfig(async () => {
  return {
    application: {},
    vite: {
      build: {
        // 产物落到 web/dist。Go 嵌入时按 dist 这个目录名读取，不能放在应用目录里。
        emptyOutDir: true,
        outDir: '../../dist',
        // 依赖 @v-c/mini-decimal 使用了 BigInt 字面量，es2015 无法表达会产生转换警告；
        // es2020 是所有现代浏览器的共同基线，对后台管理系统足够。
        target: 'es2020',
      },
      server: {
        proxy: {
          // 开发时把后台接口转到本地 Go 服务，不再依赖已删除的 Nitro Mock。
          '/admin-api': {
            changeOrigin: true,
            target: 'http://localhost:8080',
            ws: true,
          },
          // 接口文档页直接嵌后端 Swagger，开发服务器需要把这个路径转过去。
          '/swagger': {
            changeOrigin: true,
            target: 'http://localhost:8080',
          },
        },
      },
    },
  };
});
