# 配置文件管理

## 配置文件查找

应用支持**自动查找**配置文件，无需手动指定路径：

### 查找顺序

1. **命令行指定**：`-c` 或 `--config` 参数指定的路径（如果存在）
2. **项目根目录**：`./config.yaml`（Docker 友好）
3. **config 目录**：`./config/config.yaml`（传统位置）

### 使用示例

```bash
# 自动查找（推荐）
./bico-admin serve

# 指定配置文件
./bico-admin serve -c /path/to/config.yaml
./bico-admin serve --config custom.yaml

# 查看帮助
./bico-admin --help
```

### Docker 部署

Docker 环境推荐将配置文件放在项目根目录：

```bash
# 目录结构
bico-admin/
├── config.yaml          # 生产配置（根目录）
├── config/
│   ├── config.yaml      # 开发配置
│   └── config.prod.yaml # 生产配置模板
└── docker-compose.yml

# docker-compose.yml 挂载
volumes:
  - ./config.yaml:/app/config.yaml
```

### 本地开发

本地开发可以使用 `config/config.yaml`：

```bash
bico-admin/
├── config/
│   └── config.yaml  # 开发配置
└── cmd/main.go

# 直接运行（自动找到 config/config.yaml）
go run cmd/main.go serve
make serve
```

## 配置文件结构

完整配置示例见 `config/config.yaml`：

```yaml
server:
  port: 8080
  mode: debug  # debug/release
  embed_static: false

app:
  name: Bico Admin
  logo: /logo.png

database:
  driver: sqlite  # sqlite/mysql
  max_idle_conns: 10
  max_open_conns: 100
  
  sqlite:
    path: storage/data.db
  
  mysql:
    host: localhost
    port: 3306
    username: root
    password: root
    database: bico_admin
    charset: utf8mb4

log:
  level: debug       # debug/info/warn/error/fatal
  format: console    # console/json
  output: stdout     # stdout/stderr/文件路径

cache:
  driver: memory  # memory/redis
  
  redis:
    host: localhost
    port: 6379
    password: ""
    db: 0

jwt:
  secret: "bico-admin-secret-key-change-in-production"
  expire_hours: 168  # 7天

upload:
  driver: local  # local/qiniu/aliyun/tencent
  max_size: 10485760  # 10MB
  allowed_types:
    - image/jpeg
    - image/png
    - image/gif
    - image/webp
  
  local:
    base_path: storage/uploads
    serve_path: /uploads
    url_prefix: http://localhost:8080/uploads
```

## 环境配置

### 开发环境

使用 `config/config.yaml`：

```yaml
server:
  mode: debug
  embed_static: false

database:
  driver: sqlite
  sqlite:
    path: storage/data.db

cache:
  driver: memory

log:
  level: debug
  format: console
```

### 生产环境

使用 `config.yaml`（项目根目录）：

```yaml
server:
  mode: release
  embed_static: true

database:
  driver: mysql
  mysql:
    host: mysql
    port: 3306
    username: root
    password: strong-password
    database: bico_admin

cache:
  driver: redis
  redis:
    host: redis
    port: 6379

log:
  level: info
  format: json
  output: /app/logs/app.log

jwt:
  secret: "change-this-to-random-string"
```

## 配置管理最佳实践

### 1. 敏感信息管理

**不要**将生产环境的敏感信息提交到 Git：

```bash
# .gitignore
config.yaml
config.prod.yaml
*.local.yaml
```

**使用**配置模板：

```bash
# 提供模板
config/config.prod.yaml.example

# 部署时复制并修改
cp config/config.prod.yaml.example config.yaml
vim config.yaml  # 修改敏感信息
```

### 2. 多环境配置

```bash
config/
├── config.yaml          # 开发环境（默认）
├── config.prod.yaml     # 生产环境模板
├── config.test.yaml     # 测试环境
└── config.local.yaml    # 本地个性化配置（不提交）

# 使用不同配置
./bico-admin serve -c config/config.test.yaml
```

### 3. 环境变量覆盖

对于 Docker 部署，可以通过环境变量覆盖部分配置：

```yaml
# docker-compose.yml
services:
  app:
    environment:
      - DB_HOST=mysql
      - DB_PASSWORD=${DB_PASSWORD}
      - JWT_SECRET=${JWT_SECRET}
```

配置文件支持环境变量（需要额外实现）：

```yaml
database:
  mysql:
    host: ${DB_HOST:-localhost}
    password: ${DB_PASSWORD}
```

### 4. 配置验证

启动时会自动验证配置：

```bash
./bico-admin serve
# ✅ 配置加载成功: config.yaml
# ❌ 配置文件未找到，已尝试的路径: , [config.yaml config/config.yaml]
```

## 故障排查

### 配置文件找不到

**错误信息**：
```
配置文件未找到，已尝试的路径: , [config.yaml config/config.yaml]
```

**解决方法**：
1. 检查当前目录是否正确（应该在项目根目录）
2. 确认配置文件存在：`ls -la config.yaml config/config.yaml`
3. 手动指定路径：`./bico-admin serve -c /path/to/config.yaml`

### 配置解析失败

**错误信息**：
```
解析配置文件失败: yaml: unmarshal errors
```

**解决方法**：
1. 检查 YAML 语法是否正确（缩进、冒号等）
2. 使用在线工具验证：https://www.yamllint.com/
3. 参考 `config/config.yaml` 示例

### Docker 容器配置问题

**错误信息**：
```
读取配置文件失败: open /app/config/config.yaml: no such file or directory
```

**解决方法**：
```bash
# 检查挂载是否正确
docker exec bico-admin ls -la /app/config.yaml

# 更新 docker-compose.yml
volumes:
  - ./config.yaml:/app/config.yaml  # 确保路径正确
```

## 配置热加载

当前版本不支持配置热加载，修改配置后需要重启：

```bash
# 本地开发
make serve  # Ctrl+C 停止，再次运行

# Docker 部署
docker-compose restart app
```

未来版本计划支持部分配置的热加载（如日志级别）。
