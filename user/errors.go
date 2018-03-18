package user

import (
	"errors"
)

const (
	StatusError = "error"
	StatusOK    = "ok"
)

var (
	ErrInvalidRequest = errors.New("Invalid request")
	ErrUnknownError   = errors.New("Unknown error")

	ErrUserInvalidName          = errors.New("Invalid user name")
	ErrUserInvalidPassword      = errors.New("Invalid password")
	ErrUserInvalidEmail         = errors.New("Invalid email")
	ErrUserTokenExpired         = errors.New("Password reset token expired")
	ErrUserWrongLoginOrPassword = errors.New("Incorrect login or password")
)
