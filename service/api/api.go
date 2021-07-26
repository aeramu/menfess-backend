package api

import (
	"github.com/aeramu/menfess-backend/constants"
	"github.com/aeramu/menfess-backend/entity"
	"github.com/aeramu/menfess-backend/utils"
	"strings"
)

type PaginationReq struct {
	First int
	After string
}

type PaginationRes struct {
	EndCursor   string
	HasNextPage bool
}

type LoginReq struct {
	Email     string
	Password  string
	PushToken string
}

type LoginRes struct {
	Token string
}

func (req LoginReq) Validate() error {
	req.Email = strings.ToLower(req.Email)
	if err := utils.ValidateEmail(req.Email); err != nil {
		return constants.ErrInvalidEmail
	}
	if req.Password == "" {
		return constants.ErrInvalidPassword
	}
	//if req.PushToken == "" {
	//	return constants.ErrInvalidPushToken
	//}
	return nil
}

type RegisterReq struct {
	Email     string
	Password  string
	PushToken string
}

type RegisterRes struct {
	Token string
}

func (req RegisterReq) Validate() error {
	req.Email = strings.ToLower(req.Email)
	if err := utils.ValidateEmail(req.Email); err != nil {
		return constants.ErrInvalidEmail
	}
	if req.Password == "" {
		return constants.ErrInvalidPassword
	}
	//if req.PushToken == "" {
	//	return constants.ErrInvalidPushToken
	//}
	return nil
}

type UpdateProfileReq struct {
	ID     string
	Name   string
	Avatar string
	Bio    string
}

type UpdateProfileRes struct {
	Message string
}

func (req *UpdateProfileReq) Validate() error {
	if req.ID == "" {
		return constants.ErrInvalidID
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return constants.ErrInvalidName
	}
	req.Avatar = strings.TrimSpace(req.Avatar)
	if req.Avatar == "" {
		return constants.ErrInvalidAvatar
	}
	return nil
}

type GetUserReq struct {
	ID string
}

type GetUserRes struct {
	User entity.User
}

func (req GetUserReq) Validate() error {
	if req.ID == "" {
		return constants.ErrInvalidID
	}
	return nil
}

type GetMenfessListReq struct {
}

type GetMenfessListRes struct {
	MenfessList []entity.User
}

func (req GetMenfessListReq) Validate() error {
	return nil
}

type GetPostReq struct {
	ID     string
	UserID string
}

type GetPostRes struct {
	Post entity.Post
}

func (req GetPostReq) Validate() error {
	if req.ID == "" {
		return constants.ErrInvalidID
	}
	if req.UserID == "" {
		return constants.ErrInvalidUserID
	}
	return nil
}

type GetPostListReq struct {
	ParentID   string
	AuthorIDs  []string
	UserID     string
	Pagination PaginationReq
}

type GetPostListRes struct {
	PostList   []entity.Post
	Pagination PaginationRes
}

func (req *GetPostListReq) Validate() error {
	if req.UserID == "" {
		return constants.ErrInvalidUserID
	}
	if req.Pagination.First < 1 {
		req.Pagination.First = 20
	}
	return nil
}

type CreatePostReq struct {
	Body     string
	UserID   string
	AuthorID string
	ParentID string
}

type CreatePostRes struct {
	Message string
}

func (req CreatePostReq) Validate() error {
	if req.Body == "" {
		return constants.ErrInvalidBody
	}
	if req.UserID == "" {
		return constants.ErrInvalidUserID
	}
	return nil
}

type LikePostReq struct {
	PostID string
	UserID string
}

type LikePostRes struct {
	Message string
}

func (req LikePostReq) Validate() error {
	if req.PostID == "" {
		return constants.ErrInvalidPostID
	}
	if req.UserID == "" {
		return constants.ErrInvalidUserID
	}
	return nil
}

type LogoutReq struct {
	UserID    string
	PushToken string
}

type LogoutRes struct {
	Message string
}

func (req LogoutReq) Validate() error {
	if req.UserID == "" {
		return constants.ErrInvalidUserID
	}
	if req.PushToken == "" {
		return constants.ErrInvalidPushToken
	}
	return nil
}
