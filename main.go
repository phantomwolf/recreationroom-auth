package main

import (
	"fmt"
	"github.com/PhantomWolf/recreationroom-auth/user"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	confBind          = "bind"
	confPort          = "port"
	confMysqlServer   = "mysql_server"
	confMysqlPort     = "mysql_port"
	confMysqlUser     = "mysql_user"
	confMysqlPassword = "mysql_password"
	confMysqlDatabase = "mysql_database"
	confMysqlOptions  = "mysql_options"
	confDebug         = "debug"
)

func loadConfFiles() {
	viper.SetConfigName("config")
	paths := [...]string{".", "..", "config", "$HOME/.recreationroom", "/etc/recreationroom"}
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("Unable to load config file: %s\n", err.Error())
	}
}

func parseCliArgs() {
	flag.String(confBind, "127.0.0.1", "Bind to this address. Default: 127.0.0.1")
	flag.Int(confPort, 8080, "Port number to use. Default: 8080")
	flag.String(confMysqlServer, "127.0.0.1", "MySQL database server")
	flag.Int(confMysqlPort, 3306, "MySQL database port. Default: 3306")
	flag.String(confMysqlUser, "", "MySQL database user")
	flag.String(confMysqlPassword, "", "MySQL database password")
	flag.String(confMysqlDatabase, "recreationroom", "MySQL database name")
	flag.String(confMysqlOptions, "charset=utf8mb4&parseTime=true", "MySQL connection options. Default: charset=utf8mb4&parseTime=true")
	flag.Bool(confDebug, false, "Debug mode")
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)
}

func main() {
	loadConfFiles()
	parseCliArgs()

	// Set log level
	if viper.GetBool(confDebug) {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Setup User service
	mysqlDSN := fmt.Sprintf("%s:%s@(%s:%d)/%s",
		viper.GetString(confMysqlUser),
		viper.GetString(confMysqlPassword),
		viper.GetString(confMysqlServer),
		viper.GetInt(confMysqlPort),
		viper.GetString(confMysqlDatabase),
	)
	if len(viper.GetString(confMysqlOptions)) != 0 {
		mysqlDSN += "?" + viper.GetString(confMysqlOptions)
	}
	db, err := gorm.Open("mysql", mysqlDSN)
	if err != nil {
		log.Panicf("Database connection failure: %s\n", err.Error())
	}
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)

	// Setup handlers
	r := mux.NewRouter()
	user.MakeHandler(userService, r)

	// Start HTTP server
	http.Handle("/", r)
	addr := fmt.Sprintf("%s:%d", viper.GetString(confBind), viper.GetInt(confPort))
	errs := make(chan error, 2)
	go func() {
		errs <- http.ListenAndServe(addr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
	}()
	log.Infof("terminated: %s\n", (<-errs).Error())
}
