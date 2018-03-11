package user

import (
	"errors"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	maxNameLength     = 50
	maxPasswordLength = 50
)

var (
	ErrInvalidUsername = errors.New("Invalid user name")
	ErrInvalidPassword = errors.New("Invalid password")
	ErrInvalidEmail    = errors.New("Invalid email")
)

type User struct {
	ID        int64      `gorm:"PRIMARY_KEY" json:"id,string"`
	Name      string     `gorm:"type:VARCHAR(50);NOT NULL;UNIQUE" json:"name"`
	Password  string     `gorm:"type:VARCHAR(50);NOT NULL" json:"password"`
	Email     string     `gorm:"type:VARCHAR(180);NOT NULL;UNIQUE" json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

func New(name string, password string, email string) (*User, error) {
	user := &User{CreatedAt: time.Now()}
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

func (user *User) SetName(name string) error {
	if length := len(name); length == 0 || length > maxNameLength {
		log.Debugf("[user/user.go:SetName] Invalid user name length %d\n", length)
		return ErrInvalidUsername
	}
	user.Name = name
	return nil
}

func (user *User) SetPassword(password string) error {
	if length := len(password); length == 0 || length > maxPasswordLength {
		log.Debugf("[user/user.go:SetPassword] Invalid password length %d\n", length)
		return ErrInvalidPassword
	}
	user.Password = password
	return nil
}

func (user *User) SetEmail(email string) error {
	matched, _ := regexp.MatchString("[\\w_\\-.]+@[\\w_\\-.]+", email)
	if matched == false {
		log.Debugf("[user/user.go:SetEmail] Invalid email address: %s\n", email)
		return ErrInvalidEmail
	}
	user.Email = email
	return nil
}
