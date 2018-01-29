package service

import (
	"context"
	"github.com/PhantomWolf/recreationroom-auth/model"
)

type Auth interface {
	Register(ctx context.Context, name string, password string, email string) (*model.User, error)
	Unregister(ctx context.Context)
	Login(ctx context.Context, nameOrEmail string, password string) (*model.User, error)
	Logout(ctx context.Context, user *model.User) error
}
