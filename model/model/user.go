package model

import (
	"context"
	"errors"
	"github.com/PhantomWolf/recreationroom-auth/model/entity"
	"regexp"
)

const (
	MaxNameLength     = 30
	MaxPasswordLength = 30
)

type User struct {
	u entity.User
}

func (user *User) Id(ctx context.Context) int64 {
	return user.u.Id
}

func (user *User) Name(ctx context.Context) string {
	return user.u.Name
}

// User name must be non-empty string
func (user *User) SetName(ctx context.Context, name string) error {
	if len(name) == 0 {
		return errors.New("User name can't be empty string")
	}
	if len(name) > MaxNameLength {
		return errors.New("User name length must be less than " + string(MaxNameLength))
	}
	user.u.Name = name
	return nil
}

func (user *User) Password(ctx context.Context) string {
	return user.u.Password
}

func (user *User) SetPassword(ctx context.Context, password string) error {
	if len(password) == 0 {
		return errors.New("Password can't be empty")
	}
	if len(password) > MaxPasswordLength {
		return errors.New("Password can't be longer than " + string(MaxPasswordLength))
	}
	user.u.Password = password
	return nil
}

func (user *User) Email(ctx context.Context) string {
	return user.u.Email
}

func (user *User) SetEmail(ctx context.Context, email string) error {
	pattern := "[\\w\\-._]+@[\\w\\-.]+"
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		panic(err)
	}
	if !matched {
		return errors.New("Invalid email address: " + email + "\n")
	}
	user.u.Email = email
	return nil
}
