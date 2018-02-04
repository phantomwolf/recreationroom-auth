package session

import (
	"github.com/spf13/viper"
)

type Backend interface {
	Load(key string) (map[string]string, error)
	Save(key string, data map[string]string) error
	Delete(key string) error
}

func getBackend() Backend {
	switch viper.GetString("session.backend") {
	case "redis":
		return
	}
}
