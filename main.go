package main

import (
	"github.com/PhantomWolf/recreationroom-auth/config"
	"github.com/PhantomWolf/recreationroom-auth/user"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	config.Load()
	log.SetLevel(log.DebugLevel)
	// Init User service
	db, err := gorm.Open(config.DatabaseBackend(), config.DSN())
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
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
