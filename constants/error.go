package constants

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")

	ErrUserNotFound = errors.New("user not found")
	ErrInvalidID = errors.New("id is not valid")
	ErrInvalidUserID = errors.New("user id is not valid")
	ErrInvalidFollowedID = errors.New("followed id is not valid")
	ErrInvalidEmail = errors.New("email is not valid")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidPassword = errors.New("password is not valid")
	ErrWrongPassword = errors.New("wrong password")
	ErrInvalidPushToken = errors.New("push token is not valid")
	ErrInvalidName = errors.New("name is not valid")
	ErrInvalidAvatar = errors.New("avatar is not valid")

	ErrPostNotFound = errors.New("post not found")
	ErrInvalidPostID = errors.New("post id is not valid")
	ErrInvalidBody = errors.New("body is not valid")

	ErrInvalidFeedType = errors.New("invalid feed type")
)
