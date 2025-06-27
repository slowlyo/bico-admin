package config

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

// InitRedis 初始化Redis连接
func InitRedis(cfg *Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试连接
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("Redis connected successfully")
	return rdb, nil
}

// GetRedis 获取Redis实例（单例模式）
var redisInstance *redis.Client

func GetRedis() *redis.Client {
	return redisInstance
}

func SetRedis(rdb *redis.Client) {
	redisInstance = rdb
}
