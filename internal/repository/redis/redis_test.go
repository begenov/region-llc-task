package redis

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/begenov/region-llc-task/pkg/logger"
	"github.com/begenov/region-llc-task/pkg/utils"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/require"
)

var redisRepo *Redis

func init() {
	s, err := miniredis.Run()
	if err != nil {
		logger.Fatalf("miniredis.Run(): %v", err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	redisRepo = NewRedis(client)

}

func setRedis(t *testing.T) string {
	key := utils.RandomString(10)
	value := utils.RandomString(15)
	exp := time.Minute
	err := redisRepo.Set(key, value, exp)
	require.NoError(t, err)
	return key
}

func TestRedis_Set(t *testing.T) {
	setRedis(t)
}

func TestRedis_Get(t *testing.T) {
	key := setRedis(t)

	res, err := redisRepo.Get(key)
	require.NoError(t, err)
	require.NotEmpty(t, res)
}

func TestRedis_Delete(t *testing.T) {
	key := setRedis(t)
	err := redisRepo.Delete(key)
	require.NoError(t, err)
}
