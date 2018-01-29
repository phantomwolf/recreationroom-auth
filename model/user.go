package model

import (
	"errors"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"regexp"
)

const (
	MaxNameLength = 30
)

type User struct {
	ID       int64  `gorm:"type:BIGINT;PRIMARY KEY;AUTO_INCREMENT" json:"id"`
	Name     string `gorm:"type:VARCHAR(30);NOT NULL;UNIQUE" json:"name"`
	Password string `gorm:"type:VARCHAR(50);NOT NULL" json:"password"`
	Email    string `gorm:"type:VARCHAR(255);NOT NULL;UNIQUE" json:"email"`
}

func (user *User) SetName(name string) error {
	if len(name) > 30 {
		return errors.New("User name too long(" + string(MaxNameLength) + " at most)")
	} else if len(name) == 0 {
		return errors.New("User name can't be empty")
	}
	user.Name = name
	return nil
}

func (user *User) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("Password can't be empty")
	}
	user.Password = password
	return nil
}

func (user *User) SetEmail(email string) error {
	matched, _ := regexp.MatchString("[\\w_\\-.]+@[\\w_\\-.]+", email)
	if matched == false {
		return errors.New("Invalid email address: " + email)
	}
	user.Email = email
	return nil
}

func NewUser(name string, password string, email string) (*User, error) {
	user := &User{}
	return user, nil
}
