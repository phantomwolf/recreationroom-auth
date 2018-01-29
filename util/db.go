package util

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
)

var db *gorm.DB
var once sync.Once

func DB() *gorm.DB {
	once.Do(func() {
		db, err := gorm.Open("mysql", "admin:password@(localhost:3306)/recreationroom?charset=utf8mb4")
		if err != nil {
			panic(err)
		}
		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(100)
	})
	return db
}
