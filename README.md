# Bico Admin

基于Go语言的多端管理系统，包含master、admin、api三个端点，以admin端为主要开发重点。

## 项目特点

- 🎯 **多端分离架构**: master端(主控)、admin端(管理)、api端(接口)
- 🏗️ **清晰的分层架构**: Handler → Service → Repository → Database
- 🔧 **依赖注入**: 使用Google Wire进行编译时依赖注入
- 📝 **统一响应格式**: 标准化API响应结构
- 🔐 **权限控制**: 不同端点的权限隔离
- 📊 **完善的日志**: 基于Zap的结构化日志
- 🗄️ **数据库支持**: MySQL/PostgreSQL/SQLite + Redis缓存

## 技术栈

- **Web框架**: Gin
- **ORM框架**: GORM
- **数据库**: MySQL/PostgreSQL/SQLite
- **缓存**: Redis
- **依赖注入**: Google Wire
- **配置管理**: Viper
- **日志**: Zap
- **文档**: Swagger