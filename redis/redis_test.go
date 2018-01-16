package cache

import (
	"testing"
)

func TestRedisClient(t *testing.T) {
	// Connect to redis
	client := RedisClient()
	// Store a hash
	ret, err := client.HSet("128", "name", "baka").Result()
	ret, err = client.HSetNX("128", "name", "fool").Result()
	t.Log(ret, err)
	ret, err = client.HSet("128", "gender", "male").Result()
	t.Log(ret, err)
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
