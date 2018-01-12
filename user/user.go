package user

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	ID       uint64
	Name     string
	Password string
}

func New(name string, password string) *User {
	return &User{
		ID:       0,
		Name:     name,
		Password: password,
	}
}
