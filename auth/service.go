package auth

import (
	"context"
	"errors"
	"github.com/PhantomWolf/recreationroom-auth/session"
	"github.com/PhantomWolf/recreationroom-auth/user"
	"log"
)

type Service interface {
	// Create a new user, returning its uid and error object
	Register(ctx context.Context, name string, password string, email string) (uint64, error)
	// Permanently delete a user from database. Need to verify password
	Unregister(ctx context.Context, uid int64, password string) error
	// User login. Create a session
	Login(ctx context.Context, nameOrEmail string, password string) error
	// User logout. Remove session
	Logout(ctx context.Context, uid int64, sessionID string) error
	// Check if a user has logged in by querying redis
	//IsLogined(ctx context.Context, uid int64) bool
}

type service struct {
	repo user.Repository
}

func New(repo user.Repository) {
	return &service{repo: repo}
}

func (serv *service) Unregister(ctx context.Context, uid int64, password string) error {

}

func (serv *service) Register(ctx context.Context, name string, password string, email string) (uint64, error) {
	u, err := user.New(name, password, email)
	if err != nil {
		log.Printf("[auth/service.go] Invalid user")
		return 0, err
	}

	uid, err := serv.repo.Add(u)
	if err != nil {
		log.Printf("[auth/service.go] Registering failed")
		return 0, err
	}
	return uid, nil
}
