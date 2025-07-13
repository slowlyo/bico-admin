package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bico-admin/pkg/cache"
)

func TestMemoryCache(t *testing.T) {
	ctx := context.Background()

	// 创建内存缓存
	config := cache.Config{
		Driver: "memory",
		Memory: cache.MemoryConfig{
			MaxSize:           10,
			DefaultExpiration: 1 * time.Second,
			CleanupInterval:   500 * time.Millisecond,
		},
	}

	cacheInstance, err := cache.NewCache(config)
	require.NoError(t, err)
	defer cacheInstance.Close()

	// 测试基本操作
	t.Run("基本操作", func(t *testing.T) {
		// 设置缓存
		err := cacheInstance.Set(ctx, "test_key", "test_value", 5*time.Second)
		assert.NoError(t, err)

		// 获取缓存
		value, err := cacheInstance.Get(ctx, "test_key")
		assert.NoError(t, err)
		assert.Equal(t, "test_value", value)

		// 检查存在
		exists, err := cacheInstance.Exists(ctx, "test_key")
		assert.NoError(t, err)
		assert.True(t, exists)

		// 删除缓存
		err = cacheInstance.Delete(ctx, "test_key")
		assert.NoError(t, err)

		// 检查不存在
		exists, err = cacheInstance.Exists(ctx, "test_key")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	// 测试过期
	t.Run("过期测试", func(t *testing.T) {
		err := cacheInstance.Set(ctx, "expire_key", "expire_value", 100*time.Millisecond)
		assert.NoError(t, err)

		// 立即获取应该成功
		value, err := cacheInstance.Get(ctx, "expire_key")
		assert.NoError(t, err)
		assert.Equal(t, "expire_value", value)

		// 等待过期
		time.Sleep(200 * time.Millisecond)

		// 获取应该失败
		_, err = cacheInstance.Get(ctx, "expire_key")
		assert.Equal(t, cache.ErrCacheNotFound, err)
	})

	// 测试容量限制
	t.Run("容量限制", func(t *testing.T) {
		// 清空缓存
		cacheInstance.Clear(ctx)

		// 填满缓存
		for i := 0; i < 10; i++ {
			err := cacheInstance.Set(ctx, fmt.Sprintf("key_%d", i), fmt.Sprintf("value_%d", i), 1*time.Hour)
			assert.NoError(t, err)
		}

		// 尝试添加第11个应该失败
		err := cacheInstance.Set(ctx, "key_11", "value_11", 1*time.Hour)
		assert.Error(t, err)
	})
}

func TestCacheManager(t *testing.T) {
	ctx := context.Background()

	// 创建缓存
	config := cache.Config{Driver: "memory"}
	cacheInstance, err := cache.NewCache(config)
	require.NoError(t, err)
	defer cacheInstance.Close()

	manager := cache.NewManager(cacheInstance)

	// 测试JSON操作
	t.Run("JSON操作", func(t *testing.T) {
		user := map[string]interface{}{
			"id":   123,
			"name": "John Doe",
			"age":  30,
		}

		// 设置JSON
		err := manager.SetJSON(ctx, "user:123", user, 5*time.Minute)
		assert.NoError(t, err)

		// 获取JSON
		var retrievedUser map[string]interface{}
		err = manager.GetJSON(ctx, "user:123", &retrievedUser)
		assert.NoError(t, err)
		assert.Equal(t, user, retrievedUser)
	})

	// 测试GetOrSet
	t.Run("GetOrSet", func(t *testing.T) {
		callCount := 0
		setter := func() (string, error) {
			callCount++
			return "generated_value", nil
		}

		// 第一次调用应该执行setter
		value, err := manager.GetOrSet(ctx, "generated_key", setter, 5*time.Minute)
		assert.NoError(t, err)
		assert.Equal(t, "generated_value", value)
		assert.Equal(t, 1, callCount)

		// 第二次调用应该从缓存获取
		value, err = manager.GetOrSet(ctx, "generated_key", setter, 5*time.Minute)
		assert.NoError(t, err)
		assert.Equal(t, "generated_value", value)
		assert.Equal(t, 1, callCount) // setter不应该被再次调用
	})

	// 测试Remember
	t.Run("Remember", func(t *testing.T) {
		callCount := 0
		fn := func() (interface{}, error) {
			callCount++
			return map[string]string{"result": "computed"}, nil
		}

		// 第一次调用应该执行函数
		result, err := manager.Remember(ctx, "computed_key", 5*time.Minute, fn)
		assert.NoError(t, err)
		expected := map[string]string{"result": "computed"}
		assert.Equal(t, expected, result)
		assert.Equal(t, 1, callCount)

		// 第二次调用应该从缓存获取
		result, err = manager.Remember(ctx, "computed_key", 5*time.Minute, fn)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.Equal(t, 1, callCount) // 函数不应该被再次调用
	})
}

func TestCacheError(t *testing.T) {
	// 测试无效驱动
	config := cache.Config{Driver: "invalid"}
	_, err := cache.NewCache(config)
	assert.Error(t, err)

	var cacheErr *cache.CacheError
	assert.ErrorAs(t, err, &cacheErr)
	assert.Equal(t, "new", cacheErr.Op)
}

// 基准测试
func BenchmarkMemoryCache(b *testing.B) {
	ctx := context.Background()
	config := cache.Config{Driver: "memory"}
	cacheInstance, _ := cache.NewCache(config)
	defer cacheInstance.Close()

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cacheInstance.Set(ctx, fmt.Sprintf("key_%d", i), fmt.Sprintf("value_%d", i), 1*time.Hour)
		}
	})

	b.Run("Get", func(b *testing.B) {
		// 预设一些数据
		for i := 0; i < 1000; i++ {
			cacheInstance.Set(ctx, fmt.Sprintf("key_%d", i), fmt.Sprintf("value_%d", i), 1*time.Hour)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cacheInstance.Get(ctx, fmt.Sprintf("key_%d", i%1000))
		}
	})
}
