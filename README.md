# Bico Admin

基于 Go 构建的现代化后台管理系统。

## 技术栈

- **[Gin](https://github.com/gin-gonic/gin)** - HTTP Web 框架
- **[GORM](https://gorm.io/)** - ORM 数据库操作
- **[Viper](https://github.com/spf13/viper)** - 配置管理
- **[Dig](https://github.com/uber-go/dig)** - 依赖注入容器
- **[Cobra](https://github.com/spf13/cobra)** - 命令行框架

## 快速开始

```bash
# 修改配置文件
vim config/config.yaml

# 执行数据库迁移
make migrate

# 启动服务
make serve
```

## 常用命令

```bash
make help      # 查看所有可用命令
make build     # 编译应用
make serve     # 启动服务
make migrate   # 数据库迁移
make clean     # 清理构建产物
make tidy      # 整理依赖
```

## 手动命令

```bash
# 启动服务
go run cmd/main.go serve

# 指定配置文件
go run cmd/main.go serve -c config/prod.yaml

# 数据库迁移
go run cmd/main.go migrate

# 查看帮助
go run cmd/main.go --help
```

## 项目文档

详细的项目文档请查看 [docs](./docs) 目录：

- [项目结构说明](./docs/structure.md)

## 开发环境

- Go 1.21+
- MySQL 5.7+

## License

MIT
