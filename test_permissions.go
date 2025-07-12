package main

import (
	"bico-admin/internal/admin/definitions"
	"fmt"
)

func main() {
	fmt.Println("=== 所有权限 ===")
	allGroups := definitions.GetAllPermissions()
	for _, group := range allGroups {
		fmt.Printf("\n模块: %s (%s)\n", group.Name, group.Module)
		for _, perm := range group.Permissions {
			fmt.Printf("  %s - %s\n", perm.Code, perm.Name)
		}
	}

	fmt.Println("\n=== 所有权限代码 ===")
	allCodes := definitions.GetPermissionCodes()
	for _, code := range allCodes {
		fmt.Println(code)
	}

	fmt.Println("\n=== 用户模块权限 ===")
	userPerms := definitions.GetPermissionsByModule("user")
	for _, perm := range userPerms {
		fmt.Printf("%s - %s\n", perm.Code, perm.Name)
	}
}
