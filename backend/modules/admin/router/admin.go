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

			// 权限管理
			roles.Get("/:id/permissions", roleHandler.GetRolePermissions)
			roles.Put("/:id/permissions", roleHandler.AssignPermissions)
		}

		// 权限管理 - 改为代码配置，提供权限列表查询接口
		permissions := protected.Group("/permissions")
		{
			permissions.Get("/", func(c *fiber.Ctx) error {
				// 返回所有权限配置，兼容前端分页格式
				allPermissions := permission.AllPermissions
				return c.JSON(fiber.Map{
					"data":     allPermissions,
					"total":    len(allPermissions),
					"success":  true,
					"current":  1,
					"pageSize": len(allPermissions),
				})
			})

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

			permissions.Get("/roles/:role", func(c *fiber.Ctx) error {
				role := c.Params("role")
				userPermissions := permission.GetUserPermissions(role)
				return c.JSON(fiber.Map{
					"code": 200,
					"data": userPermissions,
				})
			})
		}
	}
}
