package model

import (
	"context"
	"errors"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"regexp"
)

const (
	MaxNameLength     = 50
	MaxPasswordLength = 50
)

type User struct {
	ID       int64  `gorm:"AUTO_INCREMENT" json:"id"`
	Name     string `gorm:"type:VARCHAR(50);NOT NULL;UNIQUE" json:"name"`
	Password string `gorm:"type:VARCHAR(50);NOT NULL" json:"password"`
	Email    string `gorm:"type:VARCHAR(180);NOT NULL;UNIQUE" json:"email"`
}

func (user *User) SetName(ctx context.Context, name string) error {
	if length := len(name); length == 0 || length > MaxNameLength {
		log.Printf("[model.user] Invalid user name(length %d)\n", length)
		return errors.New("Invalid user name")
	}
	user.Name = name
	return nil
}

func (user *User) SetPassword(ctx context.Context, password string) error {
	if length := len(password); length == 0 || length > MaxPasswordLength {
		log.Printf("[model.user] Invalid password(length %d)\n", length)
		return errors.New("Invalid password")
	}
	user.Password = password
	return nil
}

func (user *User) SetEmail(ctx context.Context, email string) error {
	matched, _ := regexp.MatchString("[\\w_\\-.]+@[\\w_\\-.]+", email)
	if matched == false {
		log.Printf("[model.user] Invalid email address: %s\n", email)
		return errors.New("Invalid email address")
	}
	user.Email = email
	return nil
}

func NewUser(ctx context.Context, name string, password string, email string) (*User, error) {
	user := &User{Name: name, Password: password, Email: email}
	return user, nil
}
