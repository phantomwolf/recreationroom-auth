package main

import (
	"flag"
	"fmt"
	"github.com/PhantomWolf/recreationroom-auth/user"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

// Info for connecting to a database
type DBSettings struct {
	server   string
	port     int
	user     string
	password string
	database string
}

// Parse command-line arguments
func parseArgs() DBSettings {
	var ret DBSettings
	flag.StringVar(&ret.server, "host", "localhost", "Database server")
	flag.IntVar(&ret.port, "port", 3306, "Database port")
	flag.StringVar(&ret.user, "user", "oldman", "Database user")
	flag.StringVar(&ret.password, "password", "test", "Database password")
	flag.StringVar(&ret.database, "database", "recreationroom", "Database name")
	flag.Parse()
	return ret
}

func genConnStr(settings *DBSettings) string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb3",
		settings.user,
		settings.password,
		settings.server,
		settings.port,
		settings.database,
	)
}

func main() {
	settings := parseArgs()
	connStr := genConnStr(&settings)
	fmt.Println("DSN:", connStr)
	db, err := gorm.Open("mysql", connStr)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(255)
	}

	db.Set("gorm:table_options", "ENGINE=InnoDB").Set("gorm:table_options", "CHARACTER SET utf8mb3").AutoMigrate(&user.User{})
	// If there's no admin account created, create it
	var u user.User
	if err := db.Where("ID = ?", "admin").First(&u).Error; err != nil {
		u := &user.User{Name: "admin", Password: "test", Email: "foo@bar.com", Nickname: "Oldman"}
		db.Create(u)
	}

	defer db.Close()
}
