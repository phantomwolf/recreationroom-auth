package cache

import (
	"github.com/go-redis/redis"
	"sync"
)

var redisClient *redis.Client
var once sync.Once

func RedisClient() *redis.Client {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		})
		_, err := redisClient.Ping().Result()
		if err != nil {
			panic("redis connection failed")
		}
	})
	return redisClient
}

func RedisClose() {
	if redisClient != nil {
		redisClient.Close()
	}
}
