package main

import (
	"fmt"
	"github.com/PhantomWolf/recreationroom-auth/user"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Info for connecting to a database
type Args struct {
	host     string
	user     string
	password string
        port     int
	database string
}

func main() {
	// Parse command-line arguments
	var args = Args{
            host:     "127.0.0.1",
            user:     "manager",
            port:     3306,
            password: "test",
            database: "recreationroom",
        }
	// Connect to database
	connString := fmt.Sprintf("%s:%s@(%s:%d)/%s",
		args.user,
		args.password,
                args.host,
                args.port,
		args.database,
        )
        fmt.Println("DSN:", connString)
	db, err := gorm.Open("mysql", connString)
        if err != nil {
            fmt.Println("error: ", err)
        }

        db.AutoMigrate(&user.User{})

        defer db.Close()
}
