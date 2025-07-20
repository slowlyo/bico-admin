# 🐳 Docker 部署指南

## 📋 概述

本文档介绍如何使用 Docker 部署 bico-admin 项目，包括开发环境和生产环境的部署方式。

## 🚀 快速开始

### 1. 构建镜像

```bash
# 使用 Makefile (推荐)
make docker-build

# 或直接使用脚本
./scripts/docker-build.sh
```

### 2. 开发环境部署

```bash
# 启动开发环境
make docker-dev

# 访问应用
# 前端: http://localhost:8899
# 后端API: http://localhost:8899/admin-api
```

### 3. 停止服务

```bash
make docker-stop
```

## 🏗️ 镜像架构

### 多阶段构建

1. **前端构建阶段**: 使用 Node.js 18 Alpine 构建前端资源
2. **后端构建阶段**: 使用 Go 1.23 Alpine 编译后端应用
3. **运行阶段**: 使用轻量级 Alpine 镜像运行应用

### 安全特性

- ✅ 非 root 用户运行 (appuser:1001)
- ✅ 最小化镜像体积
- ✅ 健康检查机制
- ✅ 安全的文件权限设置

## 🔧 配置管理

### 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `BICO_APP_ENVIRONMENT` | production | 运行环境 |
| `BICO_APP_DEBUG` | false | 调试模式 |
| `BICO_LOG_LEVEL` | info | 日志级别 |
| `BICO_SERVER_PORT` | 8899 | 服务端口 |
| `DB_HOST` | localhost | 数据库主机 |
| `DB_PASSWORD` | - | 数据库密码 |
| `REDIS_HOST` | localhost | Redis 主机 |
| `JWT_SECRET` | - | JWT 密钥 (必须设置) |

### 配置文件

- 开发环境: `config/app.yml`
- 生产环境: `config/app.prod.yml`

## 🌍 生产环境部署

### 使用 Docker Compose

```bash
# 复制并编辑环境变量
cp .env.example .env
vim .env

# 启动生产环境
docker-compose -f docker-compose.prod.yml up -d
```

### 环境变量文件 (.env)

```bash
# 数据库配置
DB_PASSWORD=your_secure_password
MYSQL_ROOT_PASSWORD=your_root_password

# Redis 配置
REDIS_PASSWORD=your_redis_password

# JWT 配置 (必须更改)
JWT_SECRET=your-very-secure-jwt-secret-key

# 其他配置
UPLOAD_BASE_URL=https://your-domain.com
```

## 📊 健康检查

容器包含内置健康检查：

```bash
# 检查容器健康状态
docker ps

# 查看健康检查日志
docker inspect --format='{{json .State.Health}}' bico-admin-app
```

## 🔍 故障排查

### 查看日志

```bash
# 查看应用日志
docker logs bico-admin-app

# 实时查看日志
docker logs -f bico-admin-app
```

### 常见问题

1. **端口冲突**: 确保 8899 端口未被占用
2. **权限问题**: 检查数据目录权限设置
3. **配置错误**: 验证环境变量和配置文件

### 进入容器调试

```bash
# 进入运行中的容器
docker exec -it bico-admin-app sh

# 检查配置
cat /app/config/app.yml
```

## 📁 数据持久化

### 重要目录

- `/app/data`: 数据文件 (数据库、上传文件)
- `/app/logs`: 日志文件
- `/app/config`: 配置文件

### 备份建议

```bash
# 备份数据
docker cp bico-admin-app:/app/data ./backup/data-$(date +%Y%m%d)

# 备份数据库 (MySQL)
docker exec bico-admin-mysql mysqldump -u root -p bico_admin > backup/db-$(date +%Y%m%d).sql
```

## 🚀 性能优化

### 镜像优化

- 使用多阶段构建减少镜像大小
- 利用 Docker 缓存层优化构建速度
- 最小化运行时依赖

### 运行时优化

- 合理设置资源限制
- 配置适当的健康检查间隔
- 使用 Redis 作为缓存后端

## 📝 注意事项

1. **生产环境必须更改 JWT_SECRET**
2. **建议使用外部数据库 (MySQL/PostgreSQL)**
3. **定期备份数据和配置**
4. **监控容器资源使用情况**
5. **及时更新镜像版本**

## 🆘 获取帮助

如果遇到问题，请：

1. 查看应用日志
2. 检查配置文件
3. 验证环境变量
4. 参考项目文档
