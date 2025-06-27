# Bico Admin Frontend

基于 Refine 的后台管理前端应用，已清理无用依赖和推广内容。

## 已完成的清理工作

### 🗑️ 移除的组件和依赖
- ✅ GitHubBanner - 移除GitHub推广横幅
- ✅ DevtoolsPanel & DevtoolsProvider - 移除开发调试工具
- ✅ RefineKbar & RefineKbarProvider - 移除命令面板
- ✅ @refinedev/devtools - 移除开发工具依赖
- ✅ @refinedev/kbar - 移除命令面板依赖
- ✅ @refinedev/cli - 移除CLI工具依赖
- ✅ @uiw/react-md-editor - 移除Markdown编辑器，替换为普通TextArea
- ✅ projectId配置 - 移除Refine项目追踪ID

### 📦 更新的依赖
- ✅ antd: 5.17.0 → 5.22.0 (最新版本)
- ✅ 项目名称: refine-project → bico-admin-frontend
- ✅ 脚本命令: 使用原生vite命令替代refine命令

### 🔧 新增配置
- ✅ 配置文件: `src/config/index.ts`
- ✅ 环境变量: `.env.example`
- ✅ API地址配置化

## 🚀 启动项目

```bash
# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 构建生产版本
pnpm build

# 预览生产版本
pnpm preview
```

## 📁 目录结构

```
frontend/
├── src/
│   ├── components/     # 通用组件
│   ├── pages/         # 页面组件
│   ├── contexts/      # React上下文
│   ├── config/        # 配置文件
│   ├── App.tsx        # 主应用组件
│   └── index.tsx      # 入口文件
├── public/            # 静态资源
├── .env.example       # 环境变量示例
└── package.json       # 项目配置
```

## 🔗 API配置

默认API地址配置：
- 后台管理API: `http://localhost:8080/admin/api`
- 对外API: `http://localhost:8080/api`

可通过环境变量自定义：
```bash
cp .env.example .env
# 编辑 .env 文件修改API地址
```

## 📝 注意事项

1. 已移除所有Refine的调试和推广组件
2. 使用原生Vite命令，不再依赖Refine CLI
3. Markdown编辑器已替换为普通文本框
4. 保留了Refine的核心功能和UI组件
5. 项目结构保持简洁，便于AI理解和维护
