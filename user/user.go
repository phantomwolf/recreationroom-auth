package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	maxNameLength = 50
	//TODO: let user set this value in config file
	bcryptCost = 11
)

var (
	ErrInvalidUsername = errors.New("Invalid user name")
	ErrInvalidPassword = errors.New("Invalid password")
	ErrInvalidEmail    = errors.New("Invalid email")
)

type User struct {
	ID        int64      `gorm:"PRIMARY_KEY" json:"id,string"`
	Name      string     `gorm:"type:VARCHAR(50);NOT NULL;UNIQUE" json:"name"`
	Password  string     `gorm:"type:CHAR(60);NOT NULL" json:"-"`
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
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		log.Debugf("[user/user.go] Failed to set password: %s\n", err.Error())
		return err
	}
	user.Password = string(encrypted)
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

func (user *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}
