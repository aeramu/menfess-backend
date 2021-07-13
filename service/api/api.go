package api

import (
	"errors"
	"github.com/aeramu/menfess-backend/utils"
)

type LoginReq struct {
	Email     string
	Password  string
	PushToken string
}

type LoginRes struct {
	Token string
}

type RegisterReq struct {
}

type RegisterRes struct {
}

type UpdateProfileReq struct {
}

type UpdateProfileRes struct {
}

type GetUserReq struct {
}

type GetUserRes struct {
}

type GetMenfessListReq struct {
}

type GetMenfessListRes struct {
}

type GetPostReq struct {
}

type GetPostRes struct {
}

type GetPostListReq struct {
}

type GetPostListRes struct {
}

type CreatePostReq struct {
}

type CreatePostRes struct {
}

type LikePostReq struct {
}

type LikePostRes struct {
}

func (req LoginReq) Validate() error {
	if err := utils.ValidateEmail(req.Email); err != nil {
		return errors.New("email is not valid")
	}
	if req.Password == "" {
		return errors.New("password is not valid")
	}
	if req.PushToken == "" {
		return errors.New("push token is not valid")
	}
	return nil
}

func (req RegisterReq) Validate() error {
	return nil
}

func (req UpdateProfileReq) Validate() error {
	return nil
}

func (req GetUserReq) Validate() error {
	return nil
}

func (req GetMenfessListReq) Validate() error {
	return nil
}

func (req GetPostReq) Validate() error {
	return nil
}

func (req GetPostListReq) Validate() error {
	return nil
}

func (req CreatePostReq) Validate() error {
	return nil
}

func (req LikePostReq) Validate() error {
	return nil
}
