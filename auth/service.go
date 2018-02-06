package service

import (
	"context"
	"errors"
	"github.com/PhantomWolf/recreationroom-auth/model"
	"github.com/PhantomWolf/recreationroom-auth/session"
	"github.com/PhantomWolf/recreationroom-auth/util"
	"log"
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

// uid: user id
// password: user password
// sid: session id
func Unregister(ctx context.Context, uid int64, password string, sid string) error {
	orm := util.ORM()
	sess, err := session.Load(sid)
	if err != nil {
		log.Printf("[service/auth.go] User %d must login before unregistering", uid)
		return errors.new("User not logged in")
	}

}

func Register(ctx context.Context, name string, password string, email string) (int64, error) {
	orm := util.ORM()
	user := &model.User{}
	// Query user with the same name
	orm.First(user, &model.User{Name: name})
	if user.ID != 0 {
		log.Printf("User %s already registered\n", name)
		return -1, errors.New("User already registered")
	}
	// Create new user
	if err := user.SetName(ctx, name); err != nil {
		return -1, err
	}
	if err := user.SetPassword(ctx, password); err != nil {
		return -1, err
	}
	if err := user.SetEmail(ctx, email); err != nil {
		return -1, err
	}
	orm.Create(user)
	if orm.NewRecord(user) {
		log.Printf("[service.auth] Failed to create user %s\n", user.Name)
		return -1, errors.New("User creation failed")
	}
	return user.ID, nil
}
