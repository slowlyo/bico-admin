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

	// 创建测试用户
	createTestUser(db)

	log.Println("Seed data created successfully!")
}

func createTestUser(db *gorm.DB) {
	// 检查是否已存在测试用户
	var existingUser model.User
	if err := db.Where("username = ?", "admin").First(&existingUser).Error; err == nil {
		log.Println("Test user 'admin' already exists")
		return
	}

	// 创建测试用户
	user := model.User{
		Username: "admin",
		Email:    "admin@bico-admin.com",
		Password: "123456", // 这个密码会被自动加密
		Nickname: "管理员",
		Status:   model.UserStatusActive,
	}

	// 加密密码
	if err := user.HashPassword(); err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// 保存用户
	if err := db.Create(&user).Error; err != nil {
		log.Fatal("Failed to create test user:", err)
	}

	log.Printf("Test user created: username=%s, password=123456", user.Username)
}
