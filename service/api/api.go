package api

import (
	"github.com/aeramu/menfess-backend/constants"
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
	Email     string
	Password  string
	PushToken string
}

type RegisterRes struct {
	Token string
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
		return constants.ErrInvalidEmail
	}
	if req.Password == "" {
		return constants.ErrInvalidPassword
	}
	if req.PushToken == "" {
		return constants.ErrInvalidPushToken
	}
	return nil
}

func (req RegisterReq) Validate() error {
	if err := utils.ValidateEmail(req.Email); err != nil {
		return constants.ErrInvalidEmail
	}
	if req.Password == "" {
		return constants.ErrInvalidPassword
	}
	if req.PushToken == "" {
		return constants.ErrInvalidPushToken
	}
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
