package main

import (
	"encoding/json"
	"fmt"
	"time"

	"bico-admin/internal/admin/types"
	"bico-admin/pkg/utils"
)

func main() {
	// 测试时间格式化
	now := time.Now()
	
	// 测试管理员用户响应
	userResponse := types.AdminUserResponse{
		ID:          1,
		Username:    "admin",
		Name:        "管理员",
		Status:      1,
		StatusText:  "启用",
		LastLoginAt: &utils.FormattedTime{Time: now},
		CreatedAt:   utils.NewFormattedTime(now),
		UpdatedAt:   utils.NewFormattedTime(now),
	}
	
	userJSON, err := json.MarshalIndent(userResponse, "", "  ")
	if err != nil {
		fmt.Printf("序列化用户响应失败: %v\n", err)
		return
	}
	
	fmt.Println("管理员用户响应JSON:")
	fmt.Println(string(userJSON))
	fmt.Println()
	
	// 测试角色响应
	roleResponse := types.RoleResponse{
		ID:          1,
		Name:        "超级管理员",
		Code:        "super_admin",
		Description: "系统超级管理员",
		Status:      1,
		StatusText:  "启用",
		UserCount:   1,
		CanEdit:     true,
		CanDelete:   false,
		CreatedAt:   utils.NewFormattedTime(now),
		UpdatedAt:   utils.NewFormattedTime(now),
	}
	
	roleJSON, err := json.MarshalIndent(roleResponse, "", "  ")
	if err != nil {
		fmt.Printf("序列化角色响应失败: %v\n", err)
		return
	}
	
	fmt.Println("角色响应JSON:")
	fmt.Println(string(roleJSON))
}
