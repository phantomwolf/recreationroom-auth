package session

import (
	"github.com/spf13/viper"
)

const (
	backendField = "session.backend"
)

type Backend interface {
	Load(key string) (map[string]string, error)
	Save(key string, data map[string]string) error
	Delete(key string) error
}

func getBackend() Backend {
	backend := viper.GetString(backendField)
	switch backend {
	case "redis":
		return getRedisBackend()
	}
	if len(backend) == 0 {
		panic("session.backend not set")
	} else {
		panic("Invalid session.backend: " + backend)
	}
}
