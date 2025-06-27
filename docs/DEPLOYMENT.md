# 部署指南

## 🎯 部署概述

Bico Admin 支持多种部署方式，从简单的单机部署到复杂的容器化集群部署。

## 🖥️ 单机部署

### 环境要求
- **操作系统**: Linux (Ubuntu 20.04+ / CentOS 7+)
- **Go**: 1.21+
- **Node.js**: 18+
- **MySQL**: 8.0+
- **Redis**: 6.0+ (可选)
- **Nginx**: 1.18+ (推荐)

### 1. 服务器准备
```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装必要软件
sudo apt install -y git curl wget unzip

# 安装Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 安装pnpm
npm install -g pnpm
```

### 2. 数据库安装
```bash
# 安装MySQL
sudo apt install -y mysql-server

# 配置MySQL
sudo mysql_secure_installation

# 创建数据库和用户
sudo mysql -u root -p
```

```sql
CREATE DATABASE bico_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'bico_admin'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON bico_admin.* TO 'bico_admin'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### 3. 项目部署
```bash
# 创建部署目录
sudo mkdir -p /opt/bico-admin
sudo chown $USER:$USER /opt/bico-admin
cd /opt/bico-admin

# 克隆项目
git clone https://github.com/your-username/bico-admin.git .

# 构建项目
make build

# 配置环境变量
cp backend/.env.example backend/.env
# 编辑 backend/.env 文件配置数据库连接等信息
```

### 4. 配置系统服务
```bash
# 创建systemd服务文件
sudo tee /etc/systemd/system/bico-admin.service > /dev/null <<EOF
[Unit]
Description=Bico Admin Backend Service
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/bico-admin/backend
ExecStart=/opt/bico-admin/backend/bin/server
Restart=always
RestartSec=5
Environment=PATH=/usr/local/go/bin:/usr/bin:/bin

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable bico-admin
sudo systemctl start bico-admin
sudo systemctl status bico-admin
```

### 5. 配置Nginx
```bash
# 安装Nginx
sudo apt install -y nginx

# 创建配置文件
sudo tee /etc/nginx/sites-available/bico-admin > /dev/null <<EOF
server {
    listen 80;
    server_name your-domain.com;
    
    # 前端静态文件
    location / {
        root /opt/bico-admin/frontend/dist;
        try_files \$uri \$uri/ /index.html;
    }
    
    # 后端API代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
    
    location /admin/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
    
    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
EOF

# 启用站点
sudo ln -s /etc/nginx/sites-available/bico-admin /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

## 🐳 Docker 部署

### 1. 创建 Dockerfile
```dockerfile
# backend/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/bin/server .
COPY --from=builder /app/storage ./storage

EXPOSE 8080
CMD ["./server"]
```

```dockerfile
# frontend/Dockerfile
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install

COPY . .
RUN pnpm build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### 2. Docker Compose 配置
```yaml
# docker-compose.yml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: bico-admin-mysql
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: bico_admin
      MYSQL_USER: bico_admin
      MYSQL_PASSWORD: password
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    container_name: bico-admin-redis
    ports:
      - "6379:6379"
    restart: unless-stopped

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: bico-admin-backend
    environment:
      - DB_HOST=mysql
      - DB_USER=bico_admin
      - DB_PASSWORD=password
      - DB_NAME=bico_admin
      - REDIS_HOST=redis
    ports:
      - "8080:8080"
    depends_on:
      - mysql
      - redis
    restart: unless-stopped

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: bico-admin-frontend
    ports:
      - "80:80"
    depends_on:
      - backend
    restart: unless-stopped

volumes:
  mysql_data:
```

### 3. 部署命令
```bash
# 构建和启动
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 更新服务
docker-compose pull
docker-compose up -d --force-recreate
```

## ☁️ 云服务部署

### AWS 部署
```bash
# 使用 AWS ECS 或 EKS
# 配置 RDS MySQL 实例
# 配置 ElastiCache Redis 实例
# 配置 ALB 负载均衡器
```

### 阿里云部署
```bash
# 使用容器服务 ACK
# 配置 RDS MySQL 实例
# 配置 Redis 实例
# 配置 SLB 负载均衡器
```

## 🔧 环境配置

### 生产环境配置
```env
# backend/.env
ENV=production
PORT=8080

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=bico_admin
DB_PASSWORD=your_secure_password
DB_NAME=bico_admin

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password

# JWT配置
JWT_SECRET=your_very_secure_jwt_secret_key
JWT_EXPIRE=24h

# 日志配置
LOG_LEVEL=info
LOG_PATH=/var/log/bico-admin
```

### SSL/HTTPS 配置
```bash
# 使用 Let's Encrypt 获取免费SSL证书
sudo apt install -y certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo crontab -e
# 添加: 0 12 * * * /usr/bin/certbot renew --quiet
```

## 📊 监控和日志

### 日志配置
```bash
# 配置日志轮转
sudo tee /etc/logrotate.d/bico-admin > /dev/null <<EOF
/var/log/bico-admin/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 www-data www-data
    postrotate
        systemctl reload bico-admin
    endscript
}
EOF
```

### 监控配置
```bash
# 安装监控工具
# Prometheus + Grafana
# 或使用云服务监控
```

## 🔄 更新和维护

### 应用更新
```bash
# 备份数据库
mysqldump -u bico_admin -p bico_admin > backup_$(date +%Y%m%d_%H%M%S).sql

# 更新代码
cd /opt/bico-admin
git pull origin main

# 重新构建
make build

# 重启服务
sudo systemctl restart bico-admin
sudo systemctl restart nginx
```

### 数据库迁移
```bash
# 运行数据库迁移
cd /opt/bico-admin/backend
./bin/server migrate
```

## 🛡️ 安全配置

### 防火墙配置
```bash
# 配置UFW防火墙
sudo ufw enable
sudo ufw allow ssh
sudo ufw allow 80
sudo ufw allow 443
sudo ufw deny 8080  # 只允许内部访问
```

### 安全加固
```bash
# 禁用root登录
sudo sed -i 's/PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config

# 配置fail2ban
sudo apt install -y fail2ban

# 定期安全更新
sudo apt install -y unattended-upgrades
```

## 📋 部署检查清单

### 部署前检查
- [ ] 服务器资源充足 (CPU, 内存, 磁盘)
- [ ] 数据库连接正常
- [ ] 环境变量配置正确
- [ ] SSL证书配置 (生产环境)
- [ ] 防火墙规则配置
- [ ] 备份策略制定

### 部署后验证
- [ ] 应用服务正常启动
- [ ] 数据库连接正常
- [ ] API接口响应正常
- [ ] 前端页面加载正常
- [ ] 用户登录功能正常
- [ ] 日志记录正常
- [ ] 监控指标正常

## 🆘 故障排除

### 常见问题
1. **服务启动失败**
   - 检查端口占用: `sudo netstat -tlnp | grep :8080`
   - 检查日志: `sudo journalctl -u bico-admin -f`

2. **数据库连接失败**
   - 检查数据库服务: `sudo systemctl status mysql`
   - 检查连接配置: 验证.env文件中的数据库配置

3. **前端页面无法访问**
   - 检查Nginx配置: `sudo nginx -t`
   - 检查Nginx状态: `sudo systemctl status nginx`

### 日志查看
```bash
# 应用日志
sudo journalctl -u bico-admin -f

# Nginx日志
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log

# 系统日志
sudo tail -f /var/log/syslog
```
