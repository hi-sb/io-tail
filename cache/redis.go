package cache

import "github.com/go-redis/redis"

var RedisClient *redis.Client

func RedisServiceClientInit(url string, password string) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
	})
}

