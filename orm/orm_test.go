package orm

import (
	"testing"
)

func TestConn(t *testing.T) {
	db := Conn()
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(50)
}
