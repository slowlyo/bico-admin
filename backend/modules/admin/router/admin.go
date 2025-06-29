package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"bico-admin/core/config"
	"bico-admin/core/middleware"
	"bico-admin/core/permission"
	coreRouter "bico-admin/core/router"
	"bico-admin/modules/admin/handler"
)

// SetupRoutes 设置后台管理路由
func SetupRoutes(app fiber.Router, db *gorm.DB) {
	// 获取配置
	cfg := config.New()

	// 设置认证路由（不需要认证的路由）- 直接在admin路由组下
	coreRouter.SetupAuthRoutes(app, db, cfg)

	// 需要认证的路由 - 直接在admin路由组下
	protected := app.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))

	// 创建权限中间件实例
	permissionMiddleware := permission.NewPermissionMiddleware(db)

	// 后台管理处理器
	dashboardHandler := handler.NewDashboardHandler(db)
	userHandler := handler.NewUserHandler(db)
	roleHandler := handler.NewRoleHandler(db)
	rolePermissionHandler := handler.NewRolePermissionHandler(db)

	// 后台管理路由（需要认证）
	{
		// 仪表板（所有登录用户都可以访问）
		protected.Get("/dashboard", dashboardHandler.GetDashboard)
		protected.Get("/dashboard/stats", dashboardHandler.GetStats)

		// 用户管理
		users := protected.Group("/users")
		{
			// 查看用户列表和详情
			users.Get("/", permissionMiddleware.RequirePermission(permission.UserView), userHandler.GetUsers)
			users.Get("/:id", permissionMiddleware.RequirePermission(permission.UserView), userHandler.GetUser)

			// 创建用户
			users.Post("/", permissionMiddleware.RequirePermission(permission.UserCreate), userHandler.CreateUser)

			// 更新用户
			users.Put("/:id", permissionMiddleware.RequirePermission(permission.UserUpdate), userHandler.UpdateUser)

			// 删除用户
			users.Delete("/:id", permissionMiddleware.RequirePermission(permission.UserDelete), userHandler.DeleteUser)
			users.Delete("/batch", permissionMiddleware.RequirePermission(permission.UserDelete), userHandler.BatchDeleteUsers)

			// 用户状态管理
			users.Put("/:id/status", permissionMiddleware.RequirePermission(permission.UserManageStatus), userHandler.UpdateUserStatus)

			// 密码管理
			users.Put("/:id/password", permissionMiddleware.RequirePermission(permission.UserUpdate), userHandler.ChangeUserPassword)
			users.Put("/:id/reset-password", permissionMiddleware.RequirePermission(permission.UserResetPassword), userHandler.ResetUserPassword)
		}

		// 角色管理
		roles := protected.Group("/roles")
		{
			// 查看角色列表和详情
			roles.Get("/", permissionMiddleware.RequirePermission(permission.RoleView), roleHandler.GetRoles)
			roles.Get("/:id", permissionMiddleware.RequirePermission(permission.RoleView), roleHandler.GetRole)

			// 创建角色
			roles.Post("/", permissionMiddleware.RequirePermission(permission.RoleCreate), roleHandler.CreateRole)

			// 更新角色
			roles.Put("/:id", permissionMiddleware.RequirePermission(permission.RoleUpdate), roleHandler.UpdateRole)

			// 删除角色
			roles.Delete("/:id", permissionMiddleware.RequirePermission(permission.RoleDelete), roleHandler.DeleteRole)
			roles.Delete("/batch", permissionMiddleware.RequirePermission(permission.RoleDelete), roleHandler.BatchDeleteRoles)

			// 角色状态管理
			roles.Put("/:id/status", permissionMiddleware.RequirePermission(permission.RoleUpdate), roleHandler.UpdateRoleStatus)

			// 角色权限管理 - 使用新的权限管理处理器
			roles.Get("/:id/permissions", permissionMiddleware.RequirePermission(permission.RoleView), rolePermissionHandler.GetRolePermissions)
			roles.Put("/:id/permissions", permissionMiddleware.RequirePermission(permission.RoleAssignPermissions), rolePermissionHandler.AssignRolePermissions)
			roles.Delete("/:roleId/permissions/:permissionCode", permissionMiddleware.RequirePermission(permission.RoleAssignPermissions), rolePermissionHandler.RemoveRolePermission)
		}

		// 权限管理 - 权限定义在代码中，提供查询接口
		permissions := protected.Group("/permissions")
		{
			// 获取所有权限定义（需要角色查看权限）
			permissions.Get("/", permissionMiddleware.RequirePermission(permission.RoleView), rolePermissionHandler.GetAllPermissions)

			permissions.Get("/tree", permissionMiddleware.RequirePermission(permission.RoleView), func(c *fiber.Ctx) error {
				// 返回权限树结构，按分类组织
				categories := permission.GetPermissionsByCategory()
				var treeData []map[string]interface{}

				for categoryName, perms := range categories {
					categoryNode := map[string]interface{}{
						"id":       categoryName,
						"name":     categoryName,
						"code":     categoryName,
						"children": make([]map[string]interface{}, 0),
					}

					for _, perm := range perms {
						permNode := map[string]interface{}{
							"id":          perm.Code,
							"name":        perm.Name,
							"code":        perm.Code,
							"description": perm.Description,
							"children":    []interface{}{},
						}
						categoryNode["children"] = append(categoryNode["children"].([]map[string]interface{}), permNode)
					}

					treeData = append(treeData, categoryNode)
				}

				return c.JSON(fiber.Map{
					"code": 200,
					"data": treeData,
				})
			})

			// 注意：角色权限查询已移至 /roles/:id/permissions 路由
		}
	}
}
