package util

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

var database struct {
	db   *sql.DB
	once sync.Once
}

func Db() *sql.DB {
	database.once.Do(func() {
		db, err := sql.Open("mysql", "admin:password@(localhost)/recreationroom")
		if err != nil {
			panic(err)
		}
		db.SetMaxIdleConns(5)
		db.SetMaxOpenConns(1000)
		database.db = db
	})
	return database.db
}

func Conn(ctx context.Context) (*sql.Conn, error) {
	conn, err := Db().Conn(ctx)
	return conn, err
}
