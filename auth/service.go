package auth

import (
	"context"
	"fmt"
	"github.com/PhantomWolf/recreationroom-auth/user"
	session "github.com/PhantomWolf/recreationroom-session"
	"log"
	"regexp"
)

type Service interface {
	// Create a new user, returning its uid and error object
	Register(ctx context.Context, name string, password string, email string) error
	// Permanently delete a user from database. Need to verify password
	Unregister(ctx context.Context, uid uint64, password string, sid string) error
	// User login. Create a session
	Login(ctx context.Context, nameOrEmail string, password string) (*session.Session, error)
	// User logout. Remove session
	Logout(ctx context.Context, uid uint64, sid string) error
}

type service struct {
	userRepo user.Repository
	sessRepo session.Repository
}

func New(userRepo user.Repository, sessRepo session.Repository) Service {
	return &service{userRepo: userRepo, sessRepo: sessRepo}
}

func (serv *service) Logout(ctx context.Context, uid uint64, sid string) error {
	if !serv.online(ctx, uid, sid) {
		return ErrUserNotLoggedIn
	}

	err := serv.sessRepo.Remove(uid)
	if err != nil {
		return ErrLogoutFailure
	}
	return nil
}

func (serv *service) Login(ctx context.Context, nameOrEmail string, password string) (*session.Session, error) {
	spec := &user.User{}
	if matched, _ := regexp.MatchString("[\\w\\-_.]@[\\w\\-_.]", nameOrEmail); matched {
		spec.Email = nameOrEmail
	} else if matched, _ := regexp.MatchString("[\\w\\-_.]", nameOrEmail); matched {
		spec.Name = nameOrEmail
	} else {
		return nil, ErrInvalidUser
	}

	users := serv.userRepo.Query(spec)
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}

	// Verify password
	// TODO: Hash password with salt
	if password != users[0].Password {
		return nil, ErrIncorrectPassword
	}

	// Create session for user
	sess, err := session.New(users[0].ID, 60*48)
	if err != nil {
		return nil, err
	}
	serv.sessRepo.Remove(users[0].ID)
	serv.sessRepo.Add(sess)
	return sess, nil
}

func (serv *service) online(ctx context.Context, uid uint64, sid string) bool {
	sess, err := serv.sessRepo.Find(uid)
	if err != nil || sess.ID() != sid {
		return false
	}
	return true
}

func (serv *service) Unregister(ctx context.Context, uid uint64, password string, sid string) error {
	if serv.online(ctx, uid, sid) {
		log.Printf("[auth/service.go] User %d not logged in\n", uid)
		return ErrUserNotLoggedIn
	}

	users := serv.userRepo.Query(&user.User{ID: uid})
	if len(users) == 0 {
		log.Printf("[auth/service.go] User %d not found", uid)
		return ErrUserNotFound
	}
	if users[0].Password != password {
		log.Printf("[auth/service.go] Incorrect password of user %d\n", uid)
		return ErrIncorrectPassword
	}

	if err := serv.userRepo.Remove(&user.User{ID: uid}); err != nil {
		log.Printf("[auth/service.go] Removing user %d failed\n", uid)
		return ErrUnregisterFailure
	}
	serv.sessRepo.Remove(uid)
	return nil
}

func (serv *service) Register(ctx context.Context, name string, password string, email string) error {
	u, err := user.New(name, password, email)
	if err != nil {
		log.Printf("[auth/service.go] Invalid user: %s\n", err.Error())
		return ErrInvalidUser
	}

	_, err := serv.userRepo.Add(u)
	if err != nil {
		log.Printf("[auth/service.go] Registering failed: %s\n", err.Error())
		return ErrRegisterFailure
	}
	return nil
}
