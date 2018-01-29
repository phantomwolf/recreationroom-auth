package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/PhantomWolf/recreationroom-auth/model"
	"github.com/PhantomWolf/recreationroom-auth/util"
)

type Auth interface {
	// Create a new user, returning its uid and error object
	Register(ctx context.Context, name string, password string, email string) (int64, error)
	// Permanently delete a user from database. Need to verify password
	Unregister(ctx context.Context, uid int64, password string) error
	// User login. Create a session
	Login(ctx context.Context, nameOrEmail string, password string) error
	// User logout. Remove session
	Logout(ctx context.Context, uid int64, sessionID string) error
	// Check if a user has logged in by querying redis
	//IsLogined(ctx context.Context, uid int64) bool
}

type auth struct {
}

func (s *auth) Register(ctx context.Context, name string, password string, email string) (int64, error) {
	db := util.DB()
	user := &model.User{}
	db.First(user, &model.User{Name: name})
	if user.ID != 0 {
		// User already registered
		return -1, errors.New(fmt.Sprintf("User %s already registered", name))
	}
	// Create new user
	if err := user.SetName(name); err != nil {
		return -1, err
	}
	if err := user.SetPassword(password); err != nil {
		return -1, err
	}
	if err := user.SetEmail(email); err != nil {
		return -1, err
	}
	db.Create(user)
	if db.NewRecord(user) {
		// INSERT failed
		return -1, errors.New("Failed to create user")
	}
	return user.ID, nil
}
