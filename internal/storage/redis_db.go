package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
)

var Ctx = context.Background()

func CreateRedisClient(dbNum int) *redis.Client {

	// Initialize the Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       dbNum,
	})

	return rdb
}
