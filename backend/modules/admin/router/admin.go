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

	// 后台管理处理器
	dashboardHandler := handler.NewDashboardHandler(db)
	userHandler := handler.NewUserHandler(db)
	roleHandler := handler.NewRoleHandler(db)
	rolePermissionHandler := handler.NewRolePermissionHandler(db)

	// 后台管理路由（需要认证）
	{
		// 仪表板
		protected.Get("/dashboard", dashboardHandler.GetDashboard)
		protected.Get("/dashboard/stats", dashboardHandler.GetStats)

		// 用户管理
		users := protected.Group("/users")
		{
			users.Get("/", userHandler.GetUsers)
			users.Post("/", userHandler.CreateUser)
			users.Get("/:id", userHandler.GetUser)
			users.Put("/:id", userHandler.UpdateUser)
			users.Delete("/:id", userHandler.DeleteUser)

			// 批量操作
			users.Delete("/batch", userHandler.BatchDeleteUsers)

			// 用户状态管理
			users.Put("/:id/status", userHandler.UpdateUserStatus)

			// 密码管理
			users.Put("/:id/password", userHandler.ChangeUserPassword)
			users.Put("/:id/reset-password", userHandler.ResetUserPassword)
		}

		// 角色管理
		roles := protected.Group("/roles")
		{
			roles.Get("/", roleHandler.GetRoles)
			roles.Post("/", roleHandler.CreateRole)
			roles.Get("/:id", roleHandler.GetRole)
			roles.Put("/:id", roleHandler.UpdateRole)
			roles.Delete("/:id", roleHandler.DeleteRole)

			// 批量操作
			roles.Delete("/batch", roleHandler.BatchDeleteRoles)

			// 角色状态管理
			roles.Put("/:id/status", roleHandler.UpdateRoleStatus)

			// 角色权限管理 - 使用新的权限管理处理器
			roles.Get("/:id/permissions", rolePermissionHandler.GetRolePermissions)
			roles.Put("/:id/permissions", rolePermissionHandler.AssignRolePermissions)
			roles.Delete("/:roleId/permissions/:permissionCode", rolePermissionHandler.RemoveRolePermission)
		}

		// 权限管理 - 权限定义在代码中，提供查询接口
		permissions := protected.Group("/permissions")
		{
			// 获取所有权限定义
			permissions.Get("/", rolePermissionHandler.GetAllPermissions)

			permissions.Get("/tree", func(c *fiber.Ctx) error {
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
