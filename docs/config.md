# 配置文件说明

## 配置文件查找顺序

程序启动时按以下优先级查找配置：

1. 命令行参数 `-c/--config` 指定路径
2. `./config.yaml`（项目根目录）
3. `./config/config.yaml`

对应实现：`internal/core/config/config.go` 的 `findConfigFile`。

## 常用启动方式

```bash
# 自动查找配置
./bico-admin serve

# 指定配置
./bico-admin serve -c config/config.yaml

# 迁移
./bico-admin migrate -c config/config.yaml
```

## 配置项总览

完整模板见 [config/config.yaml](../config/config.yaml)。

### server

```yaml
server:
  port: 8080
  mode: debug           # debug/release
  embed_static: true    # 是否启用嵌入前端静态资源
  admin_path: /admin    # 前端访问前缀
```

### app

```yaml
app:
  name: Bico Admin
  logo: /admin/logo.png
```

### database

支持：`sqlite` / `mysql` / `postgres`

```yaml
database:
  driver: sqlite
  max_idle_conns: 10
  max_open_conns: 100
  sqlite:
    path: storage/data.db
```

### log

```yaml
log:
  level: debug
  format: console       # console/json
  output: stdout        # stdout/stderr/文件路径
```

### cache

```yaml
cache:
  driver: memory        # memory/redis
  redis:
    host: localhost
    port: 6379
    password: ""
    db: 0
```

### jwt

```yaml
jwt:
  secret: "bico-admin-secret-key-change-in-production"
  expire_hours: 168
```

### rate_limit

```yaml
rate_limit:
  enabled: true
  rps: 100
  burst: 200
```

### upload

支持：`local` / `qiniu` / `aliyun`

```yaml
upload:
  driver: local
  max_size: 10485760
  local:
    base_path: storage/uploads
    serve_path: /uploads
    url_prefix: http://localhost:8080/uploads
```

## Docker 使用建议

`docker-compose.yml` 默认挂载：

- `./config.yaml -> /app/config.yaml`

建议部署前执行：

```bash
cp config/config.prod.yaml config.yaml
# 然后修改 config.yaml 的数据库、JWT 等敏感配置
```

## 注意事项

1. `jwt.secret` 生产环境必须替换。
2. `database.driver` 切换后需确保对应依赖服务可连通。
3. `server.admin_path` 变更后，访问前端地址也会变化。
4. 配置文件会被监听并更新到内存对象，但并不等于所有配置都会立即生效（见 [配置热更新](./config-hot-reload.md)）。
