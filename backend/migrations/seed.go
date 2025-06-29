package main

import (
	"log"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"bico-admin/core/config"
	"bico-admin/core/model"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 初始化配置
	cfg := config.New()

	// 初始化数据库
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("开始初始化RBAC权限系统数据...")

	// 1. 初始化角色数据
	initRoles(db)

	// 2. 初始化超级管理员用户
	initSuperAdmin(db)

	// 3. 分配超级管理员权限（基于代码定义的权限）
	assignSuperAdminPermissions(db)

	log.Println("RBAC权限系统数据初始化完成!")
}

// 注意：权限数据不再存储在数据库中
// 权限定义完全基于代码，位于 backend/core/permission/config.go

// initRoles 初始化角色数据
func initRoles(db *gorm.DB) {
	log.Println("初始化角色数据...")

	roles := []model.Role{
		{
			Name:        "超级管理员",
			Code:        "super_admin",
			Description: "系统超级管理员，拥有所有权限，不可删除或编辑",
			Status:      model.RoleStatusActive,
		},
		{
			Name:        "管理员",
			Code:        "manager",
			Description: "系统管理员，拥有大部分管理权限",
			Status:      model.RoleStatusActive,
		},
		{
			Name:        "普通用户",
			Code:        "user",
			Description: "普通用户，拥有基础权限",
			Status:      model.RoleStatusActive,
		},
	}

	for _, role := range roles {
		// 检查角色是否已存在
		var existingRole model.Role
		if err := db.Where("code = ?", role.Code).First(&existingRole).Error; err == nil {
			log.Printf("角色 '%s' 已存在，跳过", role.Code)
			continue
		}

		if err := db.Create(&role).Error; err != nil {
			log.Printf("创建角色 '%s' 失败: %v", role.Code, err)
		} else {
			log.Printf("创建角色: %s - %s", role.Code, role.Name)
		}
	}
}

// initSuperAdmin 初始化超级管理员用户
func initSuperAdmin(db *gorm.DB) {
	log.Println("初始化超级管理员用户...")

	// 检查是否已存在超级管理员用户
	var existingUser model.User
	if err := db.Where("username = ?", "admin").First(&existingUser).Error; err == nil {
		log.Println("超级管理员用户 'admin' 已存在")

		// 更新用户角色为 super_admin（如果不是的话）
		if existingUser.Role != "super_admin" {
			existingUser.Role = "super_admin"
			if err := db.Save(&existingUser).Error; err != nil {
				log.Printf("更新用户角色失败: %v", err)
			} else {
				log.Println("已将现有admin用户角色更新为super_admin")
			}
		}
		return
	}

	// 创建超级管理员用户
	user := model.User{
		Username: "admin",
		Email:    "admin@bico-admin.com",
		Password: "123456", // 这个密码会被自动加密
		Nickname: "超级管理员",
		Role:     "super_admin",
		Status:   model.UserStatusActive,
	}

	// 加密密码
	if err := user.HashPassword(); err != nil {
		log.Fatal("密码加密失败:", err)
	}

	// 保存用户
	if err := db.Create(&user).Error; err != nil {
		log.Fatal("创建超级管理员用户失败:", err)
	}

	log.Printf("超级管理员用户创建成功: username=%s, password=123456", user.Username)
}

// assignSuperAdminPermissions 为超级管理员角色分配所有权限（基于代码定义）
func assignSuperAdminPermissions(db *gorm.DB) {
	log.Println("为超级管理员角色分配权限...")

	// 获取超级管理员角色
	var superAdminRole model.Role
	if err := db.Where("code = ?", "super_admin").First(&superAdminRole).Error; err != nil {
		log.Printf("未找到超级管理员角色: %v", err)
		return
	}

	// 从代码中获取所有权限代码
	allPermissionCodes := []string{
		// 用户管理权限
		"user:view",
		"user:create",
		"user:update",
		"user:delete",
		"user:manage_status",
		"user:reset_password",

		// 角色管理权限
		"role:view",
		"role:create",
		"role:update",
		"role:delete",
		"role:assign_permissions",
	}

	// 检查是否已经分配了权限
	var existingCount int64
	db.Model(&model.RolePermission{}).Where("role_id = ?", superAdminRole.ID).Count(&existingCount)

	if existingCount > 0 {
		log.Printf("超级管理员角色已分配 %d 个权限，跳过", existingCount)
		return
	}

	// 为超级管理员角色分配所有权限
	for _, permissionCode := range allPermissionCodes {
		rolePermission := model.RolePermission{
			RoleID:         superAdminRole.ID,
			PermissionCode: permissionCode,
		}

		if err := db.Create(&rolePermission).Error; err != nil {
			log.Printf("分配权限 '%s' 失败: %v", permissionCode, err)
		}
	}

	log.Printf("为超级管理员角色分配了 %d 个权限", len(allPermissionCodes))
}
