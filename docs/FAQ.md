# 常见问题解答 (FAQ)

## 🚀 快速开始相关

### Q: 如何快速搭建开发环境？
A: 使用以下命令快速初始化项目：
```bash
make init    # 安装依赖并创建配置文件
make dev     # 启动开发环境
```

### Q: 首次启动时遇到数据库连接错误怎么办？
A: 请检查以下几点：
1. 确保MySQL服务已启动
2. 检查 `backend/.env` 文件中的数据库配置
3. 确保数据库用户有足够的权限
4. 运行 `make db-migrate` 初始化数据库表

### Q: 前端页面无法访问后端API怎么办？
A: 检查以下配置：
1. 后端服务是否正常启动 (http://localhost:8080/health)
2. 前端 `.env` 文件中的API地址配置是否正确
3. 检查CORS配置是否允许前端域名访问

## 🔧 开发相关

### Q: 如何添加新的业务模块？
A: 按照以下步骤添加新模块：
1. 在对应目录下创建 model、repository、service、handler 文件
2. 在 router 中注册新的路由
3. 在数据库迁移中添加新表结构
4. 前端创建对应的页面组件

### Q: 如何使用业务封装的CRUD方法？
A: 示例代码：
```go
// 创建CRUD服务
crudService := business.NewCRUDService[model.Content](db)

// 使用封装的方法
result, err := crudService.List(business.ListParams{
    Page:     1,
    PageSize: 10,
    Sort:     "id",
    Order:    "desc",
})
```

### Q: 如何自定义验证规则？
A: 在 `pkg/validator` 包中注册自定义验证器：
```go
func init() {
    validator.Validator.RegisterValidation("custom", validateCustom)
}

func validateCustom(fl validator.FieldLevel) bool {
    // 自定义验证逻辑
    return true
}
```

## 🗄️ 数据库相关

### Q: 如何添加新的数据库表？
A: 
1. 在对应模块的 `model` 目录下定义结构体
2. 在 `core/config/database.go` 的 `autoMigrate` 函数中添加新模型
3. 重启服务，GORM会自动创建表结构

### Q: 如何执行数据库迁移？
A: 使用以下命令：
```bash
make db-migrate    # 运行迁移
make db-seed       # 填充测试数据
```

### Q: 如何处理数据库关联查询？
A: 使用GORM的预加载功能：
```go
// 预加载用户的角色信息
var user model.User
db.Preload("Roles").First(&user, id)

// 注意：权限信息不再通过关联查询获取
// 权限数据通过单独的API接口获取，基于role_permissions表查询
```

## 🔐 认证授权相关

### Q: JWT Token过期后如何处理？
A: 
1. 前端检测到401状态码时自动跳转到登录页
2. 可以实现Token自动刷新机制
3. 后端提供Token刷新接口

### Q: 如何实现权限控制？
A: 
1. 后端使用中间件验证用户权限
2. 前端使用 `useCan` Hook控制组件显示
3. 路由级别的权限控制

### Q: 如何添加新的权限？
A:
1. 在 `backend/core/permission/config.go` 中添加权限常量和定义
2. 在前端 `frontend/src/constants/permissions.ts` 中添加对应的权限常量
3. 通过系统界面为角色分配新权限

## 🎨 前端相关

### Q: 如何自定义Ant Design主题？
A: 在 `App.tsx` 中配置主题：
```typescript
const theme = {
  token: {
    colorPrimary: '#1890ff',
    borderRadius: 6,
  },
};

<ConfigProvider theme={theme}>
  {/* 应用内容 */}
</ConfigProvider>
```

### Q: 如何处理表单验证？
A: 使用Ant Design的Form组件：
```typescript
<Form.Item
  name="email"
  rules={[
    { required: true, message: '请输入邮箱' },
    { type: 'email', message: '请输入有效的邮箱地址' },
  ]}
>
  <Input />
</Form.Item>
```

### Q: 如何实现国际化？
A: 
1. 安装 `react-i18next`
2. 配置语言文件
3. 使用 `useTranslation` Hook

## 🚀 部署相关

### Q: 如何部署到生产环境？
A: 
1. 使用 `make build` 构建项目
2. 配置Nginx反向代理
3. 使用systemd管理后端服务
4. 配置SSL证书

### Q: Docker部署时遇到问题怎么办？
A: 
1. 检查Docker镜像是否构建成功
2. 查看容器日志：`docker-compose logs -f`
3. 确保环境变量配置正确
4. 检查网络连接和端口映射

### Q: 如何配置HTTPS？
A: 
1. 获取SSL证书（Let's Encrypt免费证书）
2. 配置Nginx SSL设置
3. 更新前端API地址为HTTPS

## 🔍 故障排除

### Q: 服务启动失败怎么办？
A: 
1. 检查端口是否被占用：`netstat -tlnp | grep :8080`
2. 查看服务日志：`journalctl -u bico-admin -f`
3. 检查配置文件是否正确
4. 确保依赖服务（MySQL、Redis）正常运行

### Q: 前端构建失败怎么办？
A: 
1. 清理缓存：`pnpm clean` 或删除 `node_modules`
2. 重新安装依赖：`pnpm install`
3. 检查TypeScript类型错误
4. 查看构建日志中的具体错误信息

### Q: API请求失败怎么办？
A: 
1. 检查网络连接
2. 验证API地址是否正确
3. 检查请求头和参数格式
4. 查看浏览器开发者工具的网络面板

## 📊 性能优化

### Q: 如何优化数据库查询性能？
A: 
1. 添加适当的数据库索引
2. 使用分页查询避免大量数据加载
3. 实现查询结果缓存
4. 优化SQL查询语句

### Q: 如何优化前端性能？
A: 
1. 使用React.memo优化组件渲染
2. 实现代码分割和懒加载
3. 优化图片和静态资源
4. 使用CDN加速资源加载

### Q: 如何处理大量数据的列表？
A: 
1. 实现虚拟滚动
2. 使用分页加载
3. 添加搜索和筛选功能
4. 考虑服务端分页

## 🔧 开发工具

### Q: 推荐使用哪些开发工具？
A: 
- **编辑器**: VSCode + Go扩展 + React扩展
- **API测试**: Postman 或 Insomnia
- **数据库管理**: MySQL Workbench 或 DBeaver
- **版本控制**: Git + GitHub/GitLab

### Q: 如何配置代码格式化？
A: 
1. 后端使用 `go fmt` 和 `go vet`
2. 前端配置ESLint和Prettier
3. 在VSCode中启用保存时自动格式化
4. 配置Git pre-commit钩子

## 📝 其他问题

### Q: 如何贡献代码？
A: 
1. Fork项目到自己的GitHub
2. 创建功能分支进行开发
3. 提交Pull Request
4. 等待代码审查和合并

### Q: 如何报告Bug？
A: 
1. 在GitHub Issues中创建新问题
2. 提供详细的错误描述和复现步骤
3. 包含相关的日志和截图
4. 说明运行环境信息

### Q: 如何获取技术支持？
A: 
1. 查看项目文档和FAQ
2. 搜索GitHub Issues中的相关问题
3. 在社区论坛提问
4. 联系项目维护者

---

💡 **提示**: 如果您的问题不在此列表中，请在GitHub Issues中提出，我们会及时回复并更新FAQ。
