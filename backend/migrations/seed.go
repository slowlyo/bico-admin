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

	// 1. 初始化权限数据
	initPermissions(db)

	// 2. 初始化角色数据
	initRoles(db)

	// 3. 初始化超级管理员用户
	initSuperAdmin(db)

	// 4. 分配超级管理员权限
	assignSuperAdminPermissions(db)

	log.Println("RBAC权限系统数据初始化完成!")
}

// initPermissions 初始化权限数据
func initPermissions(db *gorm.DB) {
	log.Println("初始化权限数据...")

	// 从代码中获取所有权限定义
	permissionDefs := []struct {
		Code        string
		Name        string
		Description string
		Category    string
	}{
		// 系统管理权限
		{"system:view", "查看系统", "查看系统管理页面", "系统管理"},
		{"system:manage", "管理系统", "管理系统设置", "系统管理"},

		// 用户管理权限
		{"user:view", "查看用户", "查看用户列表和详情", "用户管理"},
		{"user:create", "创建用户", "创建新用户", "用户管理"},
		{"user:update", "编辑用户", "编辑用户信息", "用户管理"},
		{"user:delete", "删除用户", "删除用户", "用户管理"},
		{"user:manage_status", "管理用户状态", "启用/禁用用户", "用户管理"},
		{"user:reset_password", "重置密码", "重置用户密码", "用户管理"},

		// 角色管理权限
		{"role:view", "查看角色", "查看角色列表和详情", "角色管理"},
		{"role:create", "创建角色", "创建新角色", "角色管理"},
		{"role:update", "编辑角色", "编辑角色信息", "角色管理"},
		{"role:delete", "删除角色", "删除角色", "角色管理"},
		{"role:assign_permissions", "分配权限", "为角色分配权限", "角色管理"},
	}

	for _, permDef := range permissionDefs {
		// 检查权限是否已存在
		var existingPerm model.Permission
		if err := db.Where("code = ?", permDef.Code).First(&existingPerm).Error; err == nil {
			log.Printf("权限 '%s' 已存在，跳过", permDef.Code)
			continue
		}

		// 创建权限记录
		permission := model.Permission{
			Name:        permDef.Name,
			Code:        permDef.Code,
			Type:        model.PermissionTypeAPI,
			Description: permDef.Description,
			Status:      model.PermissionStatusActive,
		}

		if err := db.Create(&permission).Error; err != nil {
			log.Printf("创建权限 '%s' 失败: %v", permDef.Code, err)
		} else {
			log.Printf("创建权限: %s - %s", permDef.Code, permDef.Name)
		}
	}
}

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

// assignSuperAdminPermissions 为超级管理员角色分配所有权限
func assignSuperAdminPermissions(db *gorm.DB) {
	log.Println("为超级管理员角色分配权限...")

	// 获取超级管理员角色
	var superAdminRole model.Role
	if err := db.Where("code = ?", "super_admin").First(&superAdminRole).Error; err != nil {
		log.Printf("未找到超级管理员角色: %v", err)
		return
	}

	// 获取所有权限
	var allPermissions []model.Permission
	if err := db.Where("status = ?", model.PermissionStatusActive).Find(&allPermissions).Error; err != nil {
		log.Printf("获取权限列表失败: %v", err)
		return
	}

	// 检查是否已经分配了权限
	var existingCount int64
	db.Model(&model.RolePermission{}).Where("role_id = ?", superAdminRole.ID).Count(&existingCount)

	if existingCount > 0 {
		log.Printf("超级管理员角色已分配 %d 个权限，跳过", existingCount)
		return
	}

	// 为超级管理员角色分配所有权限
	for _, perm := range allPermissions {
		rolePermission := model.RolePermission{
			RoleID:         superAdminRole.ID,
			PermissionCode: perm.Code,
		}

		if err := db.Create(&rolePermission).Error; err != nil {
			log.Printf("分配权限 '%s' 失败: %v", perm.Code, err)
		}
	}

	log.Printf("为超级管理员角色分配了 %d 个权限", len(allPermissions))
}
