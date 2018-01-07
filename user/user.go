package user

import (
	_ "errors"
        _ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	ID       int64
	Name     string
	Password string
}

func (u *User) Delete() {
}

func (u *User) ChangeName(name string) {
	u.Name = name
}

func (u *User) ChangePassword(password string) {
	u.Password = password
}

func New(name string, password string) *User {
	return &User{
		ID:       8,
		Name:     name,
		Password: password,
	}
}

