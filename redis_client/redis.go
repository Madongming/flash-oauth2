package redis_client

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func Init(redisURL string) *redis.Client {
	// 解析Redis URL
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		// 如果解析失败，使用默认配置
		opts = &redis.Options{
			Addr: "localhost:6379",
		}
	}

	client := redis.NewClient(opts)

	// 测试连接
	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	return client
}
