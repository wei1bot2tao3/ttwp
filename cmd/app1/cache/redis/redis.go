package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

// 定义一个链接
var rdb *redis.Client
var ctx = context.Background()

// 初始化链接
func Redis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 密码
		DB:       0,  // 选择Redis数据库
	})

	// 检查Redis连接是否正常
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
