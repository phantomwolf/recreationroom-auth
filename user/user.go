package user

import (
	"crypto/rand"
	"github.com/lytics/base62"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

const (
	maxNameLength = 50

	//TODO: let user set this value in config file
	bcryptCost = 11
)

type User struct {
	ID          int64      `gorm:"PRIMARY_KEY" json:"id,string"`
	Name        string     `gorm:"type:VARCHAR(50);NOT NULL;UNIQUE" json:"name"`
	Password    string     `gorm:"type:CHAR(60);NOT NULL" json:"-"`
	Email       string     `gorm:"type:VARCHAR(180);NOT NULL;UNIQUE" json:"email"`
	Token       string     `gorm:"type:VARCHAR(200)" json:"-"`
	TokenExpire *time.Time `json:"-"`

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
		return ErrUserInvalidName
	}
	user.Name = name
	return nil
}

func (user *User) SetPassword(password string) error {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	user.Password = string(encrypted)
	return nil
}

func (user *User) VerifyPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return ErrUserWrongLoginOrPassword
	}
	return nil
}

func (user *User) SetEmail(email string) error {
	matched, _ := regexp.MatchString("[\\w_\\-.]+@[\\w_\\-.]+", email)
	if matched == false {
		return ErrUserInvalidEmail
	}
	user.Email = email
	return nil
}

func (user *User) SetToken() (string, error) {
	bytes := make([]byte, 100)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := base62.StdEncoding.EncodeToString(bytes)
	// Only a hash of token is stored
	encrypted, err := bcrypt.GenerateFromPassword([]byte(token), bcryptCost)
	if err != nil {
		return "", err
	}
	user.Token = string(encrypted)
	// Token will expire after 20 mins
	expire := time.Now().Add(time.Minute * 20)
	user.TokenExpire = &expire
	return token, nil
}

func (user *User) VerifyToken(token string) error {
	if user.Token == "" || user.TokenExpire == nil || user.TokenExpire.Before(time.Now()) {
		return ErrUserTokenExpired
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Token), []byte(token)); err != nil {
		return err
	}
	return nil
}

func (user *User) ClearToken() {
	user.Token = ""
	user.TokenExpire = nil
}
