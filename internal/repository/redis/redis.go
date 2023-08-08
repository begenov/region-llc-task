package redis

import (
	"time"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/begenov/region-llc-task/pkg/logger"
	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(rdb *redis.Client) *Redis {
	return &Redis{
		client: rdb,
	}
}

func (r *Redis) Set(key string, value string, expiration time.Duration) error {

	err := r.client.Set(key, value, expiration).Err()
	if err != nil {
		logger.Errorf("r.client.Set(): %v\t%s", err, key)
		return err
	}

	return nil
}

func (r *Redis) Get(key string) (string, error) {

	val := r.client.Get(key).String()

	if val == "" {
		logger.Errorf("Not Found: %v", key)
		return "", domain.ErrNotFound
	}

	return val, nil
}

func (r *Redis) Delete(key string) error {
	err := r.client.Del(key).Err()
	if err != nil {
		logger.Errorf("r.client.Del(): %v\t%s", err, key)
		return err
	}

	return nil
}
