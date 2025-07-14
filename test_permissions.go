package main

import (
	"fmt"

	"bico-admin/internal/admin/definitions"
)

func main() {
	fmt.Println("=== 权限初始化测试 ===")

	// 1. 测试权限定义加载
	fmt.Println("\n1. 权限定义加载测试:")
	allPermissions := definitions.GetAllPermissions()
	fmt.Printf("根权限数量: %d\n", len(allPermissions))

	for _, rootPerm := range allPermissions {
		fmt.Printf("- %s (%s) - 类型: %s\n", rootPerm.Name, rootPerm.Code, rootPerm.Type)
		printChildren(rootPerm.Children, "  ")
	}

	// 2. 测试扁平化权限列表
	fmt.Println("\n2. 扁平化权限列表测试:")
	flatPermissions := definitions.GetAllPermissionsFlat()
	fmt.Printf("总权限数量: %d\n", len(flatPermissions))

	actionCount := 0
	moduleCount := 0
	for _, perm := range flatPermissions {
		if perm.Type == definitions.PermissionTypeAction {
			actionCount++
		} else if perm.Type == definitions.PermissionTypeModule {
			moduleCount++
		}
	}
	fmt.Printf("模块权限: %d, 操作权限: %d\n", moduleCount, actionCount)

	// 3. 测试权限代码列表
	fmt.Println("\n3. 权限代码列表测试:")
	codes := definitions.GetPermissionCodes()
	fmt.Printf("权限代码数量: %d\n", len(codes))

	// 4. 测试权限查找
	fmt.Println("\n4. 权限查找测试:")
	testCodes := []string{
		"system",
		"system.admin_user",
		"system.admin_user:list",
		"system.role:create",
		"nonexistent",
	}

	for _, code := range testCodes {
		perm := definitions.GetPermissionByCode(code)
		if perm != nil {
			fmt.Printf("✅ %s -> %s (类型: %s, 级别: %d)\n", code, perm.Name, perm.Type, perm.Level)
		} else {
			fmt.Printf("❌ %s -> 未找到\n", code)
		}
	}

	// 5. 测试API权限检查
	fmt.Println("\n5. API权限检查测试:")
	userPermissions := []string{
		"system.admin_user:list",
		"system.admin_user:create",
		"system.role:list",
	}

	testAPIs := []string{
		"/api/admin/admin-users",
		"/api/admin/admin-users/:id",
		"/api/admin/roles",
		"/api/admin/nonexistent",
	}

	for _, api := range testAPIs {
		hasPermission := definitions.HasPermissionForAPI(userPermissions, api)
		if hasPermission {
			fmt.Printf("✅ %s -> 有权限\n", api)
		} else {
			fmt.Printf("❌ %s -> 无权限\n", api)
		}
	}

	// 6. 测试按钮权限检查
	fmt.Println("\n6. 按钮权限检查测试:")
	testButtons := []string{
		"search",
		"create",
		"edit",
		"delete",
		"nonexistent",
	}

	for _, button := range testButtons {
		hasPermission := definitions.HasPermissionForButton(userPermissions, button)
		if hasPermission {
			fmt.Printf("✅ %s -> 有权限\n", button)
		} else {
			fmt.Printf("❌ %s -> 无权限\n", button)
		}
	}

	fmt.Println("\n=== 测试完成 ===")
}

func printChildren(children []definitions.Permission, indent string) {
	for _, child := range children {
		fmt.Printf("%s- %s (%s) - 类型: %s, 级别: %d\n",
			indent, child.Name, child.Code, child.Type, child.Level)
		if len(child.Children) > 0 {
			printChildren(child.Children, indent+"  ")
		}
	}
}
