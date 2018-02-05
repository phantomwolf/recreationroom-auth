package session

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"log"
	"sync"
)

const (
	serverField   = "session.redis.server"
	portField     = "session.redis.port"
	passwordField = "session.redis.password"
	dbField       = "session.redis.db"
)

var redisBackend *RedisBackend
var redisOnce sync.Once

type RedisBackend struct {
	client *redis.Client
}

func getRedisBackend() Backend {
	redisOnce.Do(func() {
		addr := fmt.Sprintf("%s:%d", viper.GetString(serverField), viper.GetInt(portField))
		password := viper.GetString(passwordField)
		db := viper.GetInt(dbField)
		options := &redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}
		client := redis.NewClient(options)
		if client == nil {
			panic("Redis connection failed")
		}
		redisBackend = &RedisBackend{client: client}
	})
	return redisBackend
}

func (backend *RedisBackend) Load(key string) (map[string]string, error) {
	data, err := backend.client.HGetAll(key).Result()
	if err != nil {
		log.Printf("[session/storage.go] Loading key %s failed: %s\n", key, err.Error())
		return nil, err
	}
	if data == nil || len(data) == 0 {
		log.Printf("[session/storage.go] No such key %s in redis\n", key)
		return nil, errors.New("No such key")
	}
	return data, nil
}

func (backend *RedisBackend) Save(key string, data map[string]string) error {
	// map[string]string => map[string]interface{}
	tmp := map[string]interface{}{}
	for k, v := range data {
		tmp[k] = v
	}
	err := backend.client.HMSet(key, tmp).Err()
	if err != nil {
		log.Printf("[session/storage.go] Saving key %s failed: %s\n", key, err.Error())
		return err
	}
	return nil
}

func (backend *RedisBackend) Delete(key string) error {
	err := backend.client.Del(key).Err()
	if err != nil {
		log.Printf("[session/storage.go] Deleting key %s failed: %s\n", key, err.Error())
		return err
	}
	return nil
}
