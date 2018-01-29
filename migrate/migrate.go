package main

import (
	"github.com/PhantomWolf/recreationroom-auth/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

func main() {
	db, err := gorm.Open("mysql", "admin:password@(localhost:3306)/recreationroom?charset=utf8mb4")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&model.User{})
}
