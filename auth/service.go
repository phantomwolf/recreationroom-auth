package auth

import (
	"context"
	"errors"
	"github.com/PhantomWolf/recreationroom-auth/session"
	"github.com/PhantomWolf/recreationroom-auth/user"
	"log"
	"regexp"
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
	userRepo user.Repository
	sessRepo session.Repository
}

func New(userRepo user.Repository, sessRepo session.Repository) {
	return &service{userRepo: userRepo, sessRepo: sessRepo}
}

// Create a session for user
func (serv *service) Login(ctx context.Context, nameOrEmail string, password string) (*session.Session, error) {
	const emailRegexp = "[\\w_\\-.]+@[\\w_\\-.]+"
	const nameRegexp = "[\\w_]+"
	if matched, err := regexp.MatchString(emailRegexp, nameOrEmail); err != nil &&
}

func (serv *service) Unregister(ctx context.Context, uid int64, password string, sid string) error {
	sess, err := serv.sessRepo.Find(sid)
	if err != nil {
		log.Printf("[auth/service.go] User %d not logged in\n", uid)
		return errors.New("Not logged in")
	}
	if sess.Expired() {
		log.Printf("[auth/service.go] Session of user %d has expired\n", uid)
		return errors.new("Session expired")
	}

	users := serv.userRepo.Query(&user.User{ID: uid})
	if users == nil {
		log.Printf("[auth/service.go] User %s not found\n", uid)
		return errors.New("User not found")
	}
	if users[0].Password != password {
		log.Printf("[auth/service.go] Incorrect password of user %s\n", users[0].ID)
		return errors.new("Incorrect password")
	}
	// We've verified the user has logged in and
	err = serv.userRepo.Remove(user)
	if err != nil {
		log.Printf("[auth/service.go] Removing user %d failed: %s\n", uid, err.Error())
		return errors.New("Removing user failed")
	}
	return nil
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
