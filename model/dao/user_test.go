package dao

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestRead(t *testing.T) {
	db, err := sql.Open("mysql", "admin:password@(localhost)/recreationroom")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}

	userDao, err := NewUser("mysql")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", db)
	userEntity, err := userDao.Read(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("user entity: %v", *userEntity)
}
