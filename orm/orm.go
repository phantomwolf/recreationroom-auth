package orm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
)

var db *gorm.DB
var once sync.Once

func Conn() *gorm.DB {
	once.Do(func() {
		if db == nil {
			db = gorm.Open(
				"mysql",
				"oldman:test@(localhost:3306)/recreationroom",
			)
		}
	})
	return db
}

func Close() {
	if db != nil {
		db.Close()
	}
}
