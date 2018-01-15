package main

import (
	"flag"
	"fmt"
	"github.com/PhantomWolf/recreationroom-auth/user"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
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
	flag.StringVar(&ret.user, "user", "admin", "Database user")
	flag.StringVar(&ret.password, "password", "test", "Database password")
	flag.StringVar(&ret.database, "database", "recreationroom", "Database name")
	flag.Parse()
	return ret
}

func genDSN(settings *DBSettings) string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4",
		settings.user,
		settings.password,
		settings.server,
		settings.port,
		settings.database,
	)
}

func main() {
	settings := parseArgs()
	dsn := genDSN(&settings)
	// Connect to MySQL
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	// Auto Migration
	log.Println("auto migrating...")
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&user.User{})
}
