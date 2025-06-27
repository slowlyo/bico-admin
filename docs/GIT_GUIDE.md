# Git 使用指南

## 📋 .gitignore 配置说明

本项目的 `.gitignore` 文件已经配置了以下忽略规则：

### 🔧 **后端 Go 项目忽略**
- ✅ 二进制文件：`backend/bin/`, `*.exe`, `*.dll` 等
- ✅ 测试文件：`*.test`, `coverage.html` 等
- ✅ 依赖目录：`backend/vendor/`
- ✅ 调试文件：`debug`, `*.prof` 等
- ✅ 热重载工具：`.air.toml`, `tmp/`

### 🎨 **前端 Node.js 项目忽略**
- ✅ 依赖目录：`frontend/node_modules/`
- ✅ 构建输出：`frontend/dist/`, `frontend/build/`
- ✅ 缓存目录：`.cache/`, `.vite/`
- ✅ 测试覆盖率：`coverage/`
- ✅ 日志文件：`npm-debug.log*`, `yarn-error.log*`

### 🔒 **敏感信息忽略**
- ✅ 环境变量：`.env`, `backend/.env`, `frontend/.env`
- ✅ 证书密钥：`*.pem`, `*.key`, `*.crt`
- ✅ SSH密钥：`id_rsa`, `*.ppk`

### 📁 **存储文件忽略**
- ✅ 上传文件：`backend/storage/uploads/*`
- ✅ 日志文件：`backend/storage/logs/*.log`
- ✅ 缓存文件：`backend/storage/cache/*`
- ✅ 临时文件：`backend/storage/temp/*`

### 💻 **开发环境忽略**
- ✅ 编辑器配置：`.vscode/`, `.idea/`
- ✅ 操作系统文件：`.DS_Store`, `Thumbs.db`
- ✅ 备份文件：`*.backup`, `*.bak`

## 🚀 Git 工作流建议

### 1. 初始化仓库
```bash
git init
git add .
git commit -m "feat: 初始化 Bico Admin 项目"
```

### 2. 添加远程仓库
```bash
git remote add origin <your-repository-url>
git push -u origin main
```

### 3. 日常开发流程
```bash
# 创建功能分支
git checkout -b feature/your-feature-name

# 开发完成后提交
git add .
git commit -m "feat: 添加新功能"

# 推送到远程
git push origin feature/your-feature-name

# 合并到主分支
git checkout main
git merge feature/your-feature-name
git push origin main
```

## 📝 提交信息规范

建议使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```bash
# 功能开发
git commit -m "feat: 添加用户认证功能"

# 问题修复
git commit -m "fix: 修复登录页面样式问题"

# 文档更新
git commit -m "docs: 更新API文档"

# 代码重构
git commit -m "refactor: 重构用户服务层代码"

# 性能优化
git commit -m "perf: 优化数据库查询性能"

# 测试相关
git commit -m "test: 添加用户服务单元测试"

# 构建相关
git commit -m "build: 更新依赖版本"

# 配置相关
git commit -m "chore: 更新开发环境配置"
```

## 🔍 检查忽略状态

### 检查文件是否被忽略
```bash
# 检查单个文件
git check-ignore backend/bin/server

# 检查多个文件
git check-ignore backend/.env frontend/node_modules

# 显示详细信息
git check-ignore -v backend/bin/server
```

### 查看被忽略的文件
```bash
# 查看所有被忽略的文件
git status --ignored

# 只查看被忽略的文件
git status --ignored --porcelain | grep '^!!'
```

## 🛠️ 常用 Git 命令

### 状态查看
```bash
git status              # 查看工作区状态
git log --oneline       # 查看提交历史
git diff                # 查看工作区变更
git diff --cached       # 查看暂存区变更
```

### 分支管理
```bash
git branch              # 查看本地分支
git branch -r           # 查看远程分支
git branch -a           # 查看所有分支
git checkout -b <name>  # 创建并切换分支
git branch -d <name>    # 删除分支
```

### 撤销操作
```bash
git checkout -- <file> # 撤销工作区修改
git reset HEAD <file>   # 撤销暂存区修改
git reset --hard HEAD^  # 撤销最后一次提交
```

## ⚠️ 注意事项

1. **环境变量文件**：
   - `.env` 文件已被忽略，请复制 `.env.example` 并修改配置
   - 不要将包含敏感信息的 `.env` 文件提交到仓库

2. **构建文件**：
   - `backend/bin/` 和 `frontend/dist/` 已被忽略
   - 这些文件会在构建时自动生成

3. **依赖文件**：
   - `node_modules/` 已被忽略
   - 使用 `package.json` 和锁定文件管理依赖

4. **存储文件**：
   - 上传的文件和日志文件已被忽略
   - 使用 `.gitkeep` 文件保持目录结构

5. **锁定文件**：
   - 当前未忽略 `pnpm-lock.yaml` 等锁定文件
   - 建议团队统一是否提交锁定文件

## 🔧 自定义配置

如需修改忽略规则，请编辑根目录的 `.gitignore` 文件：

```bash
# 添加自定义忽略规则
echo "your-custom-file" >> .gitignore

# 忽略特定目录
echo "your-directory/" >> .gitignore

# 不忽略特定文件（使用 ! 前缀）
echo "!important-file.txt" >> .gitignore
```
