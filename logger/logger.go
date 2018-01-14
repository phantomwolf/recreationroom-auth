package logger

import (
	"log"
	"os"
	"sync"
)

var logger *log.Logger
var once sync.Once

func Get() *log.Logger {
	once.Do(func() {
		logger = log.New(os.Stderr, "auth: ", log.Ldate&log.Ltime&log.Lshortfile)
	})
	return logger
}
