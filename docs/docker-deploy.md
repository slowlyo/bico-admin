# Docker 部署指南

## 前置说明

仓库自带：

- [Dockerfile](../Dockerfile)
- [docker-compose.yml](../docker-compose.yml)

`docker-compose.yml` 会启动 3 个服务：`app`、`mysql`、`redis`。

## 快速部署

```bash
# 1) 准备配置
cp config/config.prod.yaml config.yaml

# 2) 按需修改 config.yaml
# - database.driver（建议 mysql）
# - database.mysql.*
# - cache.driver（建议 redis）
# - cache.redis.*
# - jwt.secret

# 3) 启动
docker-compose up -d --build

# 4) 执行迁移
docker exec bico-admin ./bico-admin migrate
```

访问地址：

- 后台：`http://localhost:8080/admin/`
- Swagger：`http://localhost:8080/swagger/index.html`
- 健康检查：`http://localhost:8080/health`

## 挂载目录

`docker-compose.yml` 当前挂载：

- `./config.yaml -> /app/config.yaml`
- `./data -> /app/data`
- `./logs -> /app/logs`

## 常用命令

```bash
docker-compose up -d
docker-compose down
docker-compose restart app
docker-compose logs -f app
docker exec -it bico-admin sh
docker exec bico-admin ./bico-admin migrate
docker-compose up -d --build
```

## 排障

### 1) 容器启动失败

```bash
docker-compose logs app
docker exec bico-admin cat /app/config.yaml
```

### 2) 数据库连接失败

```bash
docker-compose logs mysql
docker exec bico-mysql mysql -uroot -proot123456 -e "SELECT 1"
```

### 3) 端口冲突

修改 `docker-compose.yml`：

```yaml
ports:
  - "8081:8080"
```

访问地址改为 `http://localhost:8081/admin/`。

## 生产建议

1. 必须替换 `jwt.secret`。
2. 不使用默认数据库密码。
3. `server.mode` 使用 `release`。
4. 日志建议 `json + 文件输出`。
5. 建议在反向代理层启用 HTTPS。
