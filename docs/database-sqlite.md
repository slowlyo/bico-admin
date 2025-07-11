# SQLite 数据库支持

本项目现已支持 SQLite 数据库，适合开发环境和小型部署场景。

## 特性

- 🚀 **零配置**: 无需安装额外数据库服务
- 📁 **文件存储**: 数据存储在单个文件中，便于备份和迁移
- 🔧 **开发友好**: 适合本地开发和测试
- ⚡ **高性能**: 启用 WAL 模式提升并发性能

## 配置方式

### 1. 使用配置文件

项目默认已配置为使用 SQLite，如需修改数据库配置：
```yaml
database:
  driver: "sqlite"
  database: "data/bico_admin.db"  # 数据库文件路径
  max_idle_conns: 5
  max_open_conns: 10
  conn_max_lifetime: "1h"
```

### 2. 环境变量

```bash
export DATABASE_DRIVER=sqlite
export DATABASE_DATABASE=data/bico_admin.db
```

## 数据库文件路径

- **相对路径**: `data/bico_admin.db` (推荐)
- **绝对路径**: `/path/to/database.db`
- **内存数据库**: `:memory:` (仅用于测试)

## 使用场景

### 开发环境
```yaml
database:
  driver: "sqlite"
  database: "dev.db"
```

### 测试环境
```yaml
database:
  driver: "sqlite"
  database: ":memory:"  # 内存数据库，测试后自动清理
```

### 生产环境（小型项目）
```yaml
database:
  driver: "sqlite"
  database: "/var/lib/bico-admin/production.db"
```

## 性能优化

SQLite 连接已自动配置以下优化：

1. **WAL 模式**: 提升并发读写性能
2. **连接池**: 合理配置连接数量

## 注意事项

1. **并发限制**: SQLite 适合读多写少的场景
2. **文件权限**: 确保应用有数据库文件的读写权限
3. **备份策略**: 定期备份数据库文件
4. **目录创建**: 系统会自动创建数据库文件所在目录

## 迁移指南

### 从 MySQL/PostgreSQL 迁移到 SQLite

1. 导出现有数据
2. 修改配置文件
3. 重启应用（会自动创建表结构）
4. 导入数据

### 从 SQLite 迁移到 MySQL/PostgreSQL

1. 备份 SQLite 数据
2. 配置目标数据库
3. 修改配置文件
4. 重启应用并导入数据

## 故障排除

### 常见问题

1. **权限错误**: 检查数据库文件和目录权限
2. **磁盘空间**: 确保有足够的磁盘空间
3. **文件锁定**: 确保没有其他进程占用数据库文件

### 日志查看

启用调试模式查看详细的数据库操作日志：
```yaml
app:
  debug: true
```
