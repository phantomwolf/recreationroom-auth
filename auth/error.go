package auth

import (
	"errors"
)

var (
	ErrInvalidArgument   = errors.New("Invalid argument")
	ErrUserNotLoggedIn   = errors.New("User not logged in")
	ErrUserNotFound      = errors.New("User not found")
	ErrIncorrectPassword = errors.New("Incorrect password")
	ErrInvalidUser       = errors.New("Invalid user")
	ErrRegisterFailure   = errors.New("Register failed")
	ErrLogoutFailure     = errors.New("Logout failed")
	ErrUnregisterFailure = errors.New("Unregister failure")
)
