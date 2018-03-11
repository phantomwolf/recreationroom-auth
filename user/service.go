package user

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
)

var (
	ErrUserNotFound = errors.New("User not found")
)

type Service interface {
	Create(ctx context.Context, name string, password string, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Patch(ctx context.Context, data map[string]interface{}) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (serv *service) Create(ctx context.Context, name string, password string, email string) (*User, error) {
	user, err := New(name, password, email)
	if err != nil {
		log.Debugf("[user/service.go:Create] User creation failed: %s\n", err.Error())
		return nil, err
	}

	user, err = serv.repo.Add(user)
	if err != nil {
		log.Debugf("[user/service.go:Create] Adding user failed: %s\n", err.Error())
		return nil, err
	}
	return user, nil
}

func (serv *service) Update(ctx context.Context, user *User) error {
	if err := serv.repo.Update(user); err != nil {
		log.Debugf("[user/service.go:Update] Updating user %v failed: %s\n", user, err.Error())
		return err
	}
	return nil
}

func (serv *service) Patch(ctx context.Context, data map[string]interface{}) error {
	if err := serv.repo.Patch(data); err != nil {
		log.Debugf("[user/service.go:Patch] Patching user %v failed: %s\n", data, err.Error())
		return err
	}
	return nil
}

func (serv *service) Delete(ctx context.Context, id int64) error {
	if err := serv.repo.Remove(&User{ID: id}); err != nil {
		log.Debugf("[user/service.go:Delete] Unregistering user %d failed: %s\n", id, err.Error())
		return err
	}
	return nil
}

func (serv *service) Get(ctx context.Context, id int64) (*User, error) {
	users, err := serv.repo.Query(&User{ID: id})
	if err != nil {
		log.Debugf("[user/service.go:Show] Query uid %d failed: %s\n", id, err.Error())
		return nil, err
	}
	if len(users) == 0 {
		log.Debugf("[user/service.go:Show] User %d not found\n", id)
		return nil, ErrUserNotFound
	}
	return users[0], nil
}
