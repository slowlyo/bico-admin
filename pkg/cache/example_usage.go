package cache

import (
	"context"
	"fmt"
	"time"
)

// ExampleUsage 展示缓存包的使用方法
func ExampleUsage() {
	ctx := context.Background()

	// 1. 创建内存缓存（默认）
	memoryConfig := Config{
		Driver: "memory", // 可以省略，默认就是memory
		Memory: MemoryConfig{
			MaxSize:           1000,
			DefaultExpiration: 10 * time.Minute,
			CleanupInterval:   5 * time.Minute,
		},
	}

	cache, err := NewCache(memoryConfig)
	if err != nil {
		panic(err)
	}
	defer cache.Close()

	// 2. 基本操作
	// 设置缓存
	err = cache.Set(ctx, "user:123", "John Doe", 5*time.Minute)
	if err != nil {
		fmt.Printf("设置缓存失败: %v\n", err)
	}

	// 获取缓存
	value, err := cache.Get(ctx, "user:123")
	if err != nil {
		fmt.Printf("获取缓存失败: %v\n", err)
	} else {
		fmt.Printf("用户信息: %s\n", value)
	}

	// 检查键是否存在
	exists, err := cache.Exists(ctx, "user:123")
	if err != nil {
		fmt.Printf("检查键存在失败: %v\n", err)
	} else {
		fmt.Printf("键是否存在: %v\n", exists)
	}

	// 3. 使用缓存管理器
	manager := NewManager(cache)

	// 设置JSON对象
	user := map[string]interface{}{
		"id":   123,
		"name": "John Doe",
		"age":  30,
	}
	err = manager.SetJSON(ctx, "user:123:profile", user, 10*time.Minute)
	if err != nil {
		fmt.Printf("设置JSON缓存失败: %v\n", err)
	}

	// 获取JSON对象
	var retrievedUser map[string]interface{}
	err = manager.GetJSON(ctx, "user:123:profile", &retrievedUser)
	if err != nil {
		fmt.Printf("获取JSON缓存失败: %v\n", err)
	} else {
		fmt.Printf("用户资料: %+v\n", retrievedUser)
	}

	// 4. GetOrSet 模式
	value, err = manager.GetOrSet(ctx, "expensive:calculation", func() (string, error) {
		// 模拟耗时计算
		time.Sleep(100 * time.Millisecond)
		return "计算结果", nil
	}, 1*time.Hour)
	if err != nil {
		fmt.Printf("GetOrSet失败: %v\n", err)
	} else {
		fmt.Printf("计算结果: %s\n", value)
	}

	// 5. Remember 模式（记忆化）
	result, err := manager.Remember(ctx, "api:data", 30*time.Minute, func() (interface{}, error) {
		// 模拟API调用
		return map[string]string{
			"status": "success",
			"data":   "API响应数据",
		}, nil
	})
	if err != nil {
		fmt.Printf("Remember失败: %v\n", err)
	} else {
		fmt.Printf("API数据: %+v\n", result)
	}

	// 6. 获取统计信息
	stats, err := manager.GetStats(ctx)
	if err != nil {
		fmt.Printf("获取统计信息失败: %v\n", err)
	} else {
		fmt.Printf("缓存统计: 驱动=%s, 键数量=%d\n", stats.Driver, stats.KeyCount)
	}
}

// ExampleRedisUsage 展示Redis缓存的使用
func ExampleRedisUsage() {
	ctx := context.Background()

	// 创建Redis缓存
	redisConfig := Config{
		Driver: "redis",
		Redis: RedisConfig{
			Host:      "localhost",
			Port:      6379,
			Database:  0,
			KeyPrefix: "myapp:",
		},
	}

	cache, err := NewCache(redisConfig)
	if err != nil {
		fmt.Printf("创建Redis缓存失败: %v\n", err)
		return
	}
	defer cache.Close()

	manager := NewManager(cache)

	// 使用方法与内存缓存相同
	err = manager.SetString(ctx, "session:abc123", "user_data", 1*time.Hour)
	if err != nil {
		fmt.Printf("设置会话失败: %v\n", err)
	}

	value, err := manager.GetString(ctx, "session:abc123")
	if err != nil {
		fmt.Printf("获取会话失败: %v\n", err)
	} else {
		fmt.Printf("会话数据: %s\n", value)
	}
}

// ExampleCachePatterns 展示常见的缓存模式
func ExampleCachePatterns() {
	ctx := context.Background()

	// 创建缓存
	cache, _ := NewCache(Config{Driver: "memory"})
	defer cache.Close()

	manager := NewManager(cache)

	// 1. 缓存穿透保护
	getUserFromCache := func(userID string) (map[string]interface{}, error) {
		var user map[string]interface{}
		err := manager.GetOrSetJSON(ctx, "user:"+userID, &user, func() (interface{}, error) {
			// 模拟数据库查询
			if userID == "999" {
				// 用户不存在，缓存空值防止缓存穿透
				return map[string]interface{}{"exists": false}, nil
			}
			return map[string]interface{}{
				"id":     userID,
				"name":   "User " + userID,
				"exists": true,
			}, nil
		}, 5*time.Minute)

		return user, err
	}

	user, _ := getUserFromCache("123")
	fmt.Printf("用户: %+v\n", user)

	// 2. 多级缓存
	getDataWithFallback := func(key string) (string, error) {
		// 先从L1缓存（内存）获取
		if value, err := manager.GetString(ctx, "l1:"+key); err == nil {
			return value, nil
		}

		// 再从L2缓存（Redis）获取
		// 这里简化为同一个缓存实例
		if value, err := manager.GetString(ctx, "l2:"+key); err == nil {
			// 回填L1缓存
			manager.SetString(ctx, "l1:"+key, value, 1*time.Minute)
			return value, nil
		}

		// 从数据源获取
		value := "从数据源获取的数据"

		// 设置到两级缓存
		manager.SetString(ctx, "l1:"+key, value, 1*time.Minute)
		manager.SetString(ctx, "l2:"+key, value, 10*time.Minute)

		return value, nil
	}

	data, _ := getDataWithFallback("important_data")
	fmt.Printf("数据: %s\n", data)
}
