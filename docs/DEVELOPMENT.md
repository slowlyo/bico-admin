# 开发指南

## 🚀 快速开始

### 1. 项目初始化
```bash
# 初始化项目（安装依赖 + 创建配置文件）
make init
```

### 2. 配置环境变量
```bash
# 编辑后端配置
vim backend/.env

# 编辑前端配置  
vim frontend/.env
```

### 3. 启动开发环境
```bash
# 同时启动前后端开发服务
make dev

# 或者分别启动
make dev-backend    # 启动后端服务
make dev-frontend   # 启动前端服务
```

## 📋 常用命令

### 开发相关
```bash
make dev            # 启动开发环境
make up             # 启动开发环境（别名）
make stop           # 停止所有服务
make down           # 停止所有服务（别名）
make restart        # 重启开发服务
make status         # 检查服务状态
```

### 构建相关
```bash
make build          # 构建前后端项目
make build-backend  # 只构建后端
make build-frontend # 只构建前端
make clean          # 清理构建文件
```

### 测试和检查
```bash
make test           # 运行测试
make lint           # 代码检查和格式化
```

### 数据库相关
```bash
make db-migrate     # 运行数据库迁移
make db-seed        # 填充测试数据
```

### Docker相关
```bash
make docker-build   # 构建Docker镜像
make docker-up      # 启动Docker服务
make docker-down    # 停止Docker服务
make docker-logs    # 查看Docker日志
```

## 🌐 访问地址

### 开发环境
- 前端: http://localhost:5174
- 后端API: http://localhost:8080/api
- 后台管理API: http://localhost:8080/admin/api
- 健康检查: http://localhost:8080/health

### 生产环境
- 前端: http://localhost:4173
- 后端: http://localhost:8080

## 📁 项目结构

```
bico-admin/
├── backend/          # 后端Go服务
├── frontend/         # 前端React应用
├── Makefile         # 项目管理命令
├── README.md        # 项目说明
└── DEVELOPMENT.md   # 开发指南
```

## 🔧 开发工作流

1. **启动开发环境**
   ```bash
   make dev
   ```

2. **进行开发**
   - 后端代码修改会自动重启服务
   - 前端代码修改会自动热更新

3. **代码检查**
   ```bash
   make lint
   ```

4. **运行测试**
   ```bash
   make test
   ```

5. **构建项目**
   ```bash
   make build
   ```

## 🐛 故障排除

### 端口冲突
如果遇到端口冲突，可以：
1. 停止占用端口的服务：`make stop`
2. 或修改配置文件中的端口号

### 依赖问题
```bash
# 重新安装依赖
make install
```

### 构建失败
```bash
# 清理后重新构建
make clean
make build
```

### 查看日志
```bash
# 查看服务日志
make logs

# 查看Docker日志
make docker-logs
```

## 📝 注意事项

1. **首次运行**：确保先运行 `make init` 初始化项目
2. **数据库配置**：修改 `backend/.env` 中的数据库连接信息
3. **API地址**：修改 `frontend/.env` 中的API地址配置
4. **并发启动**：`make dev` 会同时启动前后端，如需单独启动请使用对应命令

## 🤝 贡献指南

1. 遵循项目代码规范
2. 提交前运行 `make lint` 检查代码
3. 确保 `make test` 通过
4. 提交信息遵循 Conventional Commits 规范
