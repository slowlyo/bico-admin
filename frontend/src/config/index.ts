export const config = {
  // API配置
  adminApiUrl: import.meta.env.VITE_ADMIN_API_URL || 'http://localhost:8080/admin/api',
  apiUrl: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',
  
  // 应用配置
  appName: 'Bico Admin',
  version: '1.0.0',
  
  // 开发配置
  isDev: import.meta.env.DEV,
  isProd: import.meta.env.PROD,
};

export default config;
