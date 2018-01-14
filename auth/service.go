package auth

import (
	"context"
	"github.com/PhantomWolf/recreationroom-auth/user"
)

type Service interface {
	Register(ctx context.Context, name string, password string, email string)
	Login(ctx context.Context, nameOrEmail string)
	Logout(ctx context.Context, name string)
}

type service struct {
}
