# AGENTS.md

Go (Gin + GORM) 后端 + Vue 3 (Vben Admin + antdv-next + vxe-table) 前端的后台管理系统。详细文档见 `docs/`。

## 验证

改完代码执行 `make check`（gofmt、后端 build/vet/test、前端 `vue-tsc` 与 vitest），失败即未完成。前端依赖用 `pnpm`，目录是 `web/`，管理后台应用在 `web/apps/admin/`。

## 新增标准 CRUD 模块

优先用生成器，不要手写样板：

```bash
make gen NAME=article TITLE=文章 FIELDS="title:string:标题,content:text:内容,sort:int:排序"
```

字段格式 `列名:类型:标签`，类型可选 `string/text/int/bool`。命令会生成 model、handler、前端 api/页面/表单，注册到 `internal/admin/module.go`、`internal/migrate/migrate.go` 和前端路由，并刷新 Swagger 与前端类型。生成后只改业务逻辑（校验、关联、缓存失效）。

## 目录

- `internal/admin/{handler,model,service}`：后台业务。
- `internal/pkg/crud`：声明式 CRUD 框架（`CRUDHandler` + hook），见 `docs/crud-pkg.md`。
- `web/apps/admin/src/api/crud.ts`：`createCrudApi`；`hooks/use-crud-grid.ts`、`hooks/use-crud-form.ts`，见 `docs/frontend-crud.md`。
- `web/apps/admin/src/api/generated/schema.d.ts`：由 Swagger 生成，禁止手改。

## 约定

- 标准增删改查走 `CRUDHandler` 的 hook，不要手写 Gin 路由方法。
- 返回给前端的业务错误用 `crud.NewBizError("提示")`；其他错误只记日志，前端看到统一的"服务器内部错误"。处理器里不要再用 `h.Error(c, err.Error())`，统一 `h.Fail(c, err)`；服务层要展示给用户的错误定义成 `BizError`。
- 缓存失效等外部副作用放在 `AfterCreateCommit` / `AfterUpdateCommit` / `AfterDeleteCommit`，不要放在事务内。
- 查重用 `UniqueFields`，同时给列建唯一索引。排序字段必须加入 `SortableFields` 白名单。
- 修改后端请求或响应结构后执行 `make swagger web-types`，让前端类型同步。
- 前端列表页用 `useCrudGrid`，弹窗表单用 `useCrudForm`，接口用 `createCrudApi`，按钮权限用 `v-access:code`，权限码格式 `{命名空间}:{模块}:{menu|list|create|edit|delete}`。
- 注释使用中文，说明意图而不是复述代码；未被要求时不要执行 git 提交或分支操作。
