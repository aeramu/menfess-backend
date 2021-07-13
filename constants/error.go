package constants

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")

	ErrUserNotFound = errors.New("user not found")
	ErrInvalidEmail = errors.New("email is not valid")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidPassword = errors.New("password is not valid")
	ErrWrongPassword = errors.New("wrong password")
	ErrInvalidPushToken = errors.New("push token is not valid")
)
