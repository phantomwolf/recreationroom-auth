package user

import (
	"errors"
)

const (
	CodeSuccess = iota // 0
	CodeInvalidRequest
	CodeUnknownError

	CodeUserTokenExpired
	CodeUserWrongLoginOrPassword
	CodeUserCreateFailure
	CodeUserUpdateFailure
	CodeUserDeleteFailure
	CodeUserGetFailure

	CodePasswordResetFailure
	CodePasswordCreateFailure
	CodePasswordUpdateFailure

	StatusOK    = "ok"
	StatusError = "error"
)

var (
	ErrInvalidRequest = errors.New("Invalid request")
	ErrUnknownError   = errors.New("Unknown error")

	ErrUserInvalidName          = errors.New("Invalid user name")
	ErrUserInvalidPassword      = errors.New("Invalid password")
	ErrUserInvalidEmail         = errors.New("Invalid email")
	ErrUserTokenExpired         = errors.New("Password reset token expired")
	ErrUserWrongLoginOrPassword = errors.New("Incorrect login or password")
	ErrUserAlreadyExists        = errors.New("User already exists")
	ErrUserNotFound             = errors.New("User not found")
)
