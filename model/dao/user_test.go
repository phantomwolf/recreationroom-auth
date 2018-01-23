package dao

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestRead(t *testing.T) {
	db, _ := sql.Open("mysql", "admin:password@(localhost)/recreationroom")
	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", db)
	conn, _ := db.Conn(ctx)
	defer db.Close()

	userDao, err := NewUser("mysql")
	if err != nil {
		t.Fatal(err)
	}

	userEntity, err := userDao.Read(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("user entity: %v", *userEntity)
}
