package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	// Database config
	databaseBackend       = "database.backend"
	databaseMysqlServer   = "database.mysql.server"
	databaseMysqlPort     = "database.mysql.port"
	databaseMysqlUser     = "database.mysql.user"
	databaseMysqlPassword = "database.mysql.password"
	databaseMysqlDatabase = "database.mysql.database"
	databaseMysqlOptions  = "database.mysql.options"
)

func Load() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("$HOME/.recreationroom")
	viper.AddConfigPath("/etc/recreationroom")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error loading config file: %s\n", err.Error()))
	}
}

func DatabaseBackend() string {
	return viper.GetString(databaseBackend)
}

func DSN() string {
	var dsn string
	backend := viper.GetString(databaseBackend)
	switch backend {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@(%s:%d)/%s",
			viper.GetString(databaseMysqlUser),
			viper.GetString(databaseMysqlPassword),
			viper.GetString(databaseMysqlServer),
			viper.GetInt(databaseMysqlPort),
			viper.GetString(databaseMysqlDatabase),
		)
		options := viper.GetString(databaseMysqlOptions)
		if len(options) != 0 {
			dsn += fmt.Sprintf("?%s", options)
		}

	default:
		panic(fmt.Errorf("Unsupported database: %s\n", backend))
	}
	return dsn
}
