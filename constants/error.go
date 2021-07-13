package constants

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")

	ErrUserNotFound = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
)
