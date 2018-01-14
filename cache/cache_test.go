package cache

import (
	"testing"
)

func TestRedisClient(t *testing.T) {
	// Connect to redis
	client := RedisClient()
	// Store a hash
	client.HSet("128", "name", "baka")
	client.HSet("128", "gender", "male")
	// Retrieve data
	data, err := client.HGetAll("128").Result()
	if err != nil {
		t.Errorf("Failed to retrieve data: %v\n", err)
	}
	t.Logf("data retrieved: %v\n", data)
	// Delete data
	client.Del("128")
	client.Close()
}
