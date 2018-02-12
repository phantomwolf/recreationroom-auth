package user

import (
	"errors"
	"log"
	"regexp"
	"time"
)

const (
	maxNameLength     = 50
	maxPasswordLength = 50
)

type User struct {
	ID        uint64 `gorm:"PRIMARY_KEY" json:"id"`
	Name      string `gorm:"type:VARCHAR(50);NOT NULL;UNIQUE" json:"name"`
	Password  string `gorm:"type:VARCHAR(50);NOT NULL" json:"password"`
	Email     string `gorm:"type:VARCHAR(180);NOT NULL;UNIQUE" json:"email"`
	DeletedAt *time.Time
}

func (user *User) SetName(name string) error {
	if length := len(name); length == 0 || length > maxNameLength {
		log.Printf("[model.user] Invalid user name(length %d)\n", length)
		return errors.New("Invalid user name")
	}
	user.Name = name
	return nil
}

func (user *User) SetPassword(password string) error {
	if length := len(password); length == 0 || length > maxPasswordLength {
		log.Printf("[model.user] Invalid password(length %d)\n", length)
		return errors.New("Invalid password")
	}
	user.Password = password
	return nil
}

func (user *User) SetEmail(email string) error {
	matched, _ := regexp.MatchString("[\\w_\\-.]+@[\\w_\\-.]+", email)
	if matched == false {
		log.Printf("[user/user.go] Invalid email address: %s\n", email)
		return errors.New("Invalid email address")
	}
	user.Email = email
	return nil
}

func New(name string, password string, email string) (*User, error) {
	user := &User{}
	if err := user.SetName(name); err != nil {
		return nil, err
	}
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	if err := user.SetEmail(email); err != nil {
		return nil, err
	}
	return user, nil
}
