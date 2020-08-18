package Config

import (
	"gopkg.in/redis.v5"
)

func RedisConnect(options *redis.Options) *redis.Client {
	redisClient := redis.NewClient(options)

	return redisClient
}
