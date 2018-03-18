package main

import (
	"github.com/PhantomWolf/recreationroom-auth/config"
	"github.com/PhantomWolf/recreationroom-auth/user"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	config.Load()
	db, err := gorm.Open(config.DatabaseBackend(), config.DSN())
	if err != nil {
		log.Printf("Database connection error: %s\n", err.Error())
		os.Exit(255)
	}
	defer db.Close()

	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&user.User{})
}
