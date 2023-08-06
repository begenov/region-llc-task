package redis

import (
	"context"
	"fmt"

	"github.com/begenov/region-llc-task/internal/config"
	"github.com/go-redis/redis/v8"
)

func CreateClient(config config.ConfigRedis) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redisClient.Ping(): %v", err)
	}
	return redisClient, err
}
