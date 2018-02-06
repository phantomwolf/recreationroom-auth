package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func Load() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.recreationroom")
	viper.AddConfigPath("/etc/recreationroom")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error loading config file: %s\n", err.Error()))
	}
}
