package main

import (
	"flag"
	"fmt"
	"github.com/PhantomWolf/recreationroom-auth/user"
	"github.com/go-kit/kit"
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
	return fmt.Sprintf("%s:%s@(%s:%d)/%s",
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

	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&user.User{})
	defer db.Close()
}
