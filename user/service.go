package user

import (
	"context"
	"regexp"

	log "github.com/sirupsen/logrus"
)

type Service interface {
	Create(ctx context.Context, name string, password string, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Patch(ctx context.Context, data map[string]interface{}) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*User, error)
	ResetPassword(ctx context.Context, nameOrEmail string) error
	CreatePassword(ctx context.Context, id int64, token string, newPassword string) error
	UpdatePassword(ctx context.Context, id int64, password string, newPassword string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ResetPassword finds user by name or email
// and generate a password reset token
func (serv *service) ResetPassword(ctx context.Context, nameOrEmail string) error {
	matched, _ := regexp.MatchString("[\\w_\\-.]+@[\\w_\\-.]+", nameOrEmail)
	query := &User{}
	if matched {
		query.Email = nameOrEmail
	} else {
		query.Name = nameOrEmail
	}

	users, err := serv.repo.Query(query)
	if err != nil {
		log.Debugf("[user/service] User '%s' not found: %s\n", nameOrEmail, err.Error())
		return err
	}
	user := &users[0]

	token, err := user.SetToken()
	if err != nil {
		log.Debugf("[user/service] Token generation failed: %s\n", err.Error())
		return err
	}
	//TODO: Remove this after finishing debugging
	log.Debugf("Password reset token: %s\n", token)

	data := map[string]interface{}{
		"token":        user.Token,
		"token_expire": user.TokenExpire,
	}
	if err := serv.repo.Patch(data); err != nil {
		log.Debugf("[user/service] Generating password reset token failed: %s\n", err.Error())
		return err
	}

	//TODO: Send password reset link through email

	return nil
}

// CreatePassword creates a password for a user(requires a password reset token).
// Returns nil if success; otherwise, return the error
func (serv *service) CreatePassword(ctx context.Context, id int64, token string, newPassword string) error {
	users, err := serv.repo.Query(&User{ID: id})
	if err != nil {
		log.Debugf("[user/service] User %d not found: %s\n", id, err.Error())
		return err
	}
	user := &users[0]

	if err := user.VerifyToken(token); err != nil {
		return err
	}

	if err := user.SetPassword(newPassword); err != nil {
		return err
	}
	user.ClearToken()

	if err := serv.repo.Update(user); err != nil {
		return err
	}

	return nil
}

// UpdatePassword sets password to <newPassword> for user <id>(requires current password <password>)
func (serv *service) UpdatePassword(ctx context.Context, id int64, password string, newPassword string) error {
	users, err := serv.repo.Query(&User{ID: id})
	if err != nil {
		log.Debugf("[user/service] User %d not found: %s\n", id, err.Error())
		return err
	}
	user := &users[0]

	if err := user.VerifyPassword(password); err != nil {
		return err
	}

	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	if err := serv.repo.Update(user); err != nil {
		return err
	}

	return nil
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
		log.Debugf("[user/service.go] Deleting user %d failed: %s\n", id, err.Error())
		return err
	}
	return nil
}

func (serv *service) Get(ctx context.Context, id int64) (*User, error) {
	users, err := serv.repo.Query(&User{ID: id})
	if err != nil {
		log.Debugf("[user/service.go] Getting user %d failed: %s\n", id, err.Error())
		return nil, err
	}
	return &users[0], nil
}
