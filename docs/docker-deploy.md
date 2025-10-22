# Docker 部署指南

## 快速开始

### 1. 使用 Docker Compose（推荐）

```bash
# 克隆项目
git clone <your-repo-url>
cd bico-admin

# 修改配置（复制到项目根目录）
cp config/config.prod.yaml config.yaml
# 编辑 config.yaml，修改数据库密码、JWT密钥等

# 一键启动
docker-compose up -d

# 执行数据库迁移
docker exec bico-admin ./bico-admin migrate

# 查看日志
docker-compose logs -f app
```

访问 http://localhost:8080

### 2. 单独构建镜像

```bash
# 构建镜像
docker build -t bico-admin:latest .

# 运行容器
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  --name bico-admin \
  bico-admin:latest
```

## 配置说明

### docker-compose.yml 服务

- **app**: 主应用
  - 端口: 8080
  - 依赖: MySQL、Redis
  
- **mysql**: 数据库
  - 端口: 3306
  - 默认密码: root123456（生产环境请修改）
  
- **redis**: 缓存
  - 端口: 6379

### 目录挂载

```
./config.yaml  → /app/config.yaml  # 配置文件
./data         → /app/data         # 数据文件（上传等）
./logs         → /app/logs         # 日志文件
```

## 生产环境配置

### 1. 修改敏感信息

编辑 `config.yaml`:

```yaml
jwt:
  secret: "使用强随机字符串"

database:
  mysql:
    password: "使用强密码"
```

同步修改 `docker-compose.yml`:

```yaml
mysql:
  environment:
    MYSQL_ROOT_PASSWORD: "与配置文件中的密码一致"
```

### 2. 反向代理（Nginx）

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 3. 启用 HTTPS

```bash
# 使用 Let's Encrypt
certbot --nginx -d your-domain.com
```

## 常用命令

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 重启服务
docker-compose restart app

# 查看日志
docker-compose logs -f app

# 进入容器
docker exec -it bico-admin sh

# 执行数据库迁移
docker exec bico-admin ./bico-admin migrate

# 重新构建镜像
docker-compose up -d --build

# 清理数据（危险操作）
docker-compose down -v
```

## 备份与恢复

### 数据库备份

```bash
# 备份
docker exec bico-mysql mysqldump -uroot -proot123456 bico_admin > backup.sql

# 恢复
docker exec -i bico-mysql mysql -uroot -proot123456 bico_admin < backup.sql
```

### 数据卷备份

```bash
# 备份 MySQL 数据
docker run --rm -v bico-admin_mysql-data:/data -v $(pwd):/backup alpine \
  tar czf /backup/mysql-data.tar.gz /data

# 恢复
docker run --rm -v bico-admin_mysql-data:/data -v $(pwd):/backup alpine \
  tar xzf /backup/mysql-data.tar.gz -C /
```

## 监控与日志

### 日志收集

生产环境建议使用 `json` 格式日志，配合 ELK/Loki 等日志系统：

```yaml
log:
  level: info
  format: json
  output: /app/logs/app.log
```

### 健康检查

在 `docker-compose.yml` 中添加：

```yaml
app:
  healthcheck:
    test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/admin-api/app-config"]
    interval: 30s
    timeout: 3s
    retries: 3
```

## 性能优化

### 1. 限制资源使用

```yaml
app:
  deploy:
    resources:
      limits:
        cpus: '2'
        memory: 1G
      reservations:
        cpus: '0.5'
        memory: 512M
```

### 2. 启用 Gzip 压缩

在 Nginx 反向代理中启用 gzip。

### 3. 数据库连接池

```yaml
database:
  max_idle_conns: 10
  max_open_conns: 100
```

## 故障排查

### 容器无法启动

```bash
# 查看日志
docker-compose logs app

# 检查配置文件
docker exec bico-admin cat /app/config/config.yaml
```

### 无法连接数据库

```bash
# 检查网络
docker network inspect bico-admin_bico-network

# 测试连接
docker exec bico-mysql mysql -uroot -proot123456 -e "SELECT 1"
```

### 端口被占用

```bash
# 修改 docker-compose.yml 中的端口映射
ports:
  - "8081:8080"  # 将宿主机端口改为 8081
```

## 安全建议

1. ✅ 修改默认密码（数据库、JWT密钥）
2. ✅ 使用环境变量存储敏感信息
3. ✅ 定期更新镜像
4. ✅ 限制容器权限（非 root 运行）
5. ✅ 启用防火墙，仅开放必要端口
6. ✅ 使用 HTTPS
7. ✅ 定期备份数据
