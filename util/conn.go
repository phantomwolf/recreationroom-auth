package util

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

var db *sql.DB
var once sync.Once
