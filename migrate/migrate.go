package main

import (
	"github.com/PhantomWolf/recreationroom-auth/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
)

func main() {
	db, err := gorm.Open("mysql", "admin:password@(localhost:3306)/recreationroom?charset=utf8mb4")
	if err != nil {
		log.Println(err)
		os.Exit(255)
	}
	defer db.Close()
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.User{})
}
