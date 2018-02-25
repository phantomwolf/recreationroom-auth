package auth

import (
	"context"
	"fmt"
	"github.com/PhantomWolf/recreationroom-auth/session"
	"github.com/PhantomWolf/recreationroom-auth/user"
	"log"
	"regexp"
)

type Service interface {
	// Create a new user, returning its uid and error object
	Register(ctx context.Context, name string, password string, email string) (uint64, error)
	// Permanently delete a user from database. Need to verify password
	Unregister(ctx context.Context, uid uint64, password string, sid string) error
	// User login. Create a session
	Login(ctx context.Context, nameOrEmail string, password string) (*session.Session, error)
	// User logout. Remove session
	Logout(ctx context.Context, uid uint64, sid string) error
	// Check if a user has logged in by querying redis
	Online(ctx context.Context, uid uint64, sid string) bool
}

type service struct {
	userRepo user.Repository
	sessRepo session.Repository
}

func New(userRepo user.Repository, sessRepo session.Repository) Service {
	return &service{userRepo: userRepo, sessRepo: sessRepo}
}

func (serv *service) Logout(ctx context.Context, uid uint64, sid string) error {
	if !serv.Online(ctx, uid, sid) {
		return &ErrUserNotLoggedIn{"User not logged in. No need to logout"}
	}

	err := serv.sessRepo.Remove(uid)
	if err != nil {
		return &ErrLogoutFailure{err.Error()}
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
		return &ErrInvalidUserInfo{"Invalid name or email"}
	}

	users := serv.userRepo.Query(spec)
	if len(users) == 0 {
		return &ErrUserNotFound{fmt.Sprintf("Usere %s not found", nameOrEmail)}
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

func (serv *service) Online(ctx context.Context, uid uint64, sid string) bool {
	sess, err := serv.sessRepo.Find(uid)
	if err != nil || sess.ID() != sid {
		return false
	}
	return true
}

func (serv *service) Unregister(ctx context.Context, uid uint64, password string, sid string) error {
	if serv.Online(ctx, uid, sid) {
		log.Printf("[auth/service.go] User %d not logged in\n", uid)
		return &ErrUserNotLoggedIn{"User must login in before unregistering"}
	}

	users := serv.userRepo.Query(&user.User{ID: uid})
	if len(users) == 0 {
		log.Printf("[auth/service.go] User %d not found", uid)
		return &ErrUserNotFound{"User doesn't exist"}
	}
	if users[0].Password != password {
		log.Printf("[auth/service.go] Incorrect password of user %d\n", uid)
		return &ErrIncorrectPassword{"Incorrect password"}
	}

	if err := serv.userRepo.Remove(uid); err != nil {
		log.Printf("[auth/service.go] Removing user %d failed\n", uid)
		return err
	}
	serv.sessRepo.Remove(uid)
	return nil
}

func (serv *service) Register(ctx context.Context, name string, password string, email string) (uint64, error) {
	u, err := user.New(name, password, email)
	if err != nil {
		log.Printf("[auth/service.go] Invalid user: %s\n", err.Error())
		return 0, &ErrInvalidUserInfo{err.Error()}
	}

	uid, err := serv.userRepo.Add(u)
	if err != nil {
		log.Printf("[auth/service.go] Registering failed: %s\n", err.Error())
		return 0, &ErrRegisterFailure{err.Error()}
	}
	return uid, nil
}
