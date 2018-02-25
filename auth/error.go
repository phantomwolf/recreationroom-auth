package auth

import (
	"github.com/PhantomWolf/recreationroom-auth/error"
)

type ErrUserNotLoggedIn error.ErrorString
type ErrUserNotFound error.ErrorString
type ErrIncorrectPassword error.ErrorString
type ErrInvalidUserInfo error.ErrorString
type ErrRegisterFailure error.ErrorString
type ErrLogoutFailure error.ErrorString
