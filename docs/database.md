# 数据库配置

本项目支持多种数据库，可根据不同场景选择合适的数据库类型。

## 支持的数据库

- **MySQL** - 生产环境推荐，高并发支持
- **PostgreSQL** - 功能丰富，适合复杂查询
- **SQLite** - 开发环境友好，零配置

## 配置方式

### 配置文件

在 `config/config.yaml` 中配置数据库连接：

```yaml
database:
  driver: "mysql"                    # 数据库类型: mysql/postgresql/sqlite
  host: "localhost"                  # 数据库主机
  port: 3306                         # 数据库端口
  username: "root"                   # 用户名
  password: "password"               # 密码
  database: "bico_admin"             # 数据库名
  charset: "utf8mb4"                 # 字符集
  max_idle_conns: 10                 # 最大空闲连接数
  max_open_conns: 100                # 最大打开连接数
  conn_max_lifetime: "1h"            # 连接最大生存时间
```

### 环境变量

```bash
export DATABASE_DRIVER=mysql
export DATABASE_HOST=localhost
export DATABASE_PORT=3306
export DATABASE_USERNAME=root
export DATABASE_PASSWORD=password
export DATABASE_DATABASE=bico_admin
```

## 不同数据库配置示例

### MySQL 配置

```yaml
database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "bico_admin"
  charset: "utf8mb4"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: "1h"
```

### PostgreSQL 配置

```yaml
database:
  driver: "postgresql"
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "password"
  database: "bico_admin"
  sslmode: "disable"                 # SSL 模式: disable/require/verify-full
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: "1h"
```

### SQLite 配置

```yaml
database:
  driver: "sqlite"
  database: "data/bico_admin.db"     # 数据库文件路径
  max_idle_conns: 5
  max_open_conns: 10
  conn_max_lifetime: "1h"
```

## 使用场景建议

### 开发环境
- **推荐**: SQLite
- **优势**: 零配置、快速启动、便于调试
```yaml
database:
  driver: "sqlite"
  database: "dev.db"
```

### 测试环境
- **推荐**: SQLite (内存模式)
- **优势**: 测试隔离、自动清理
```yaml
database:
  driver: "sqlite"
  database: ":memory:"
```

### 生产环境
- **推荐**: MySQL 或 PostgreSQL
- **优势**: 高并发、高可用、丰富特性

#### 小型项目
```yaml
database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "bico_admin"
  password: "secure_password"
  database: "bico_admin"
```

#### 大型项目
```yaml
database:
  driver: "postgresql"
  host: "db.example.com"
  port: 5432
  username: "bico_admin"
  password: "secure_password"
  database: "bico_admin"
  sslmode: "require"
  max_idle_conns: 20
  max_open_conns: 200
  conn_max_lifetime: "30m"
```

## 数据库初始化

应用启动时会自动：
1. 创建数据库表结构
2. 初始化基础数据（超级管理员、权限等）
3. 执行数据库迁移

## 性能优化

### 连接池配置
- `max_idle_conns`: 根据并发量调整，通常为 CPU 核数的 2-4 倍
- `max_open_conns`: 不超过数据库最大连接数限制
- `conn_max_lifetime`: 避免长连接问题，建议 30m-1h

### 索引优化
项目已为常用查询字段添加索引：
- 用户表：username、email
- 角色表：name
- 权限表：code

## 数据库迁移

### 切换数据库类型

1. **备份现有数据**
2. **修改配置文件**
3. **重启应用**（自动创建新表结构）
4. **导入数据**（如需要）

### 数据导出导入

```bash
# MySQL 导出
mysqldump -u root -p bico_admin > backup.sql

# PostgreSQL 导出
pg_dump -U postgres bico_admin > backup.sql

# SQLite 导出
sqlite3 bico_admin.db .dump > backup.sql
```

## 故障排除

### 连接问题
1. 检查数据库服务是否启动
2. 验证连接参数（主机、端口、用户名、密码）
3. 确认网络连通性和防火墙设置

### 权限问题
1. 确保数据库用户有足够权限
2. 检查数据库和表的访问权限
3. 验证 SSL/TLS 配置（如启用）

### 性能问题
1. 监控连接池使用情况
2. 检查慢查询日志
3. 优化数据库配置参数

### 日志调试

启用数据库调试日志：
```yaml
app:
  debug: true
database:
  log_level: "info"  # 可选: silent/error/warn/info
```

## 最佳实践

1. **生产环境使用专用数据库用户**，避免使用 root 用户
2. **定期备份数据库**，制定备份策略
3. **监控数据库性能**，及时发现问题
4. **使用连接池**，避免频繁建立连接
5. **合理设置超时时间**，防止长时间阻塞
