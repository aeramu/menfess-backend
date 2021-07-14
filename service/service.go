package service

import (
	"context"
	"github.com/aeramu/menfess-backend/constants"
	"github.com/aeramu/menfess-backend/entity"
	"github.com/aeramu/menfess-backend/service/api"
)

type Service interface {
	Login(ctx context.Context, req api.LoginReq) (*api.LoginRes, error)
	Register(ctx context.Context, req api.RegisterReq) (*api.RegisterRes, error)
	UpdateProfile(ctx context.Context, req api.UpdateProfileReq) (*api.UpdateProfileRes, error)
	GetUser(ctx context.Context, req api.GetUserReq) (*api.GetUserRes, error)
	GetMenfessList(ctx context.Context, req api.GetMenfessListReq) (*api.GetMenfessListRes, error)
	GetPost(ctx context.Context, req api.GetPostReq) (*api.GetPostRes, error)
	GetPostList(ctx context.Context, req api.GetPostListReq) (*api.GetPostListRes, error)
	CreatePost(ctx context.Context, req api.CreatePostReq) (*api.CreatePostRes, error)
	LikePost(ctx context.Context, req api.LikePostReq) (*api.LikePostRes, error)
}

func NewService(adapter Adapter) Service {
	return &service {
		adapter: adapter,
	}
}

type service struct {
	adapter Adapter
}

func (s *service) Login(ctx context.Context, req api.LoginReq) (*api.LoginRes, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	user, err := s.adapter.UserModule.FindUserByEmail(ctx, req.Email)
	if err != nil {
		if err == constants.ErrUserNotFound {
			return nil, constants.ErrUserNotFound
		}
		s.adapter.LogModule.Log(err, req, "[Login] failed find user from repo")
		return nil, constants.ErrInternalServerError
	}

	if err := s.adapter.AuthModule.ComparePassword(ctx, user.Account.Password, req.Password); err != nil {
		return nil, constants.ErrWrongPassword
	}

	if err := s.adapter.NotificationModule.AddPushToken(ctx, user.ID, req.PushToken); err != nil {
		s.adapter.LogModule.Log(err, req, "[Login] failed add notification push token")
		return nil, constants.ErrInternalServerError
	}

	token, err := s.adapter.AuthModule.GenerateToken(ctx, *user)
	if err != nil {
		s.adapter.LogModule.Log(err, req, "[Login] failed generate token")
		return nil, constants.ErrInternalServerError
	}

	return &api.LoginRes{Token: token}, nil
}

func (s *service) Register(ctx context.Context, req api.RegisterReq) (*api.RegisterRes, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.adapter.UserModule.FindUserByEmail(ctx, req.Email); err != nil {
		if err != constants.ErrUserNotFound {
			s.adapter.LogModule.Log(err, req, "[Register] failed find user from repo")
			return nil, constants.ErrInternalServerError
		}
	} else {
		return nil, constants.ErrEmailAlreadyRegistered
	}

	hash, err := s.adapter.AuthModule.HashPassword(ctx, req.Password)
	if err != nil {
		s.adapter.LogModule.Log(err, req, "[Register] failed hash password")
		return nil, constants.ErrInternalServerError
	}

	user := entity.User{
		Account: entity.Account{
			Email:    req.Email,
			Password: hash,
		},
	}

	id, err := s.adapter.UserModule.InsertUser(ctx, user)
	if err != nil {
		s.adapter.LogModule.Log(err, req, "[Register] failed insert user")
		return nil, constants.ErrInternalServerError
	}
	user.ID = id

	token, err := s.adapter.AuthModule.GenerateToken(ctx, user)
	if err != nil {
		s.adapter.LogModule.Log(err, req, "[Register] failed generate token")
		return nil, constants.ErrInternalServerError
	}

	return &api.RegisterRes{Token: token}, nil
}

// TODO: Refactor case when user not found, expected to insert new profile
func (s *service) UpdateProfile(ctx context.Context, req api.UpdateProfileReq) (*api.UpdateProfileRes, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	profile, err := s.adapter.UserModule.FindUserByID(ctx, req.ID)
	if err != nil {
		s.adapter.LogModule.Log(err, req, "[UpdateProfile] failed get user")
		return nil, constants.ErrInternalServerError
	}

	profile.Profile.Name = req.Name
	profile.Profile.Avatar = req.Avatar
	profile.Profile.Bio = req.Bio

	if err := s.adapter.UserModule.SaveProfile(ctx, *profile); err != nil {
		s.adapter.LogModule.Log(err, req, "[UpdateProfile] failed save profile")
		return nil, constants.ErrInternalServerError
	}

	return &api.UpdateProfileRes{Message: "success"}, nil
}

func (s *service) GetUser(ctx context.Context, req api.GetUserReq) (*api.GetUserRes, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	user, err := s.adapter.UserModule.FindUserByID(ctx, req.ID)
	if err != nil {
		s.adapter.LogModule.Log(err, req, "[GetUser] failed get user")
		return nil, constants.ErrInternalServerError
	}

	return &api.GetUserRes{User: *user}, nil
}

func (s *service) GetMenfessList(ctx context.Context, req api.GetMenfessListReq) (*api.GetMenfessListRes, error) {
	menfessList, err := s.adapter.UserModule.FindMenfessList(ctx)
	if err != nil {
		s.adapter.LogModule.Log(err, req, "[GetMenfessList] failed get menfess list")
		return nil, constants.ErrInternalServerError
	}

	return &api.GetMenfessListRes{MenfessList: menfessList}, nil
}

func (s *service) GetPost(ctx context.Context, req api.GetPostReq) (*api.GetPostRes, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	post, err := s.adapter.PostModule.FindPostByID(ctx, req.ID, req.UserID)
	if err != nil {
		if err == constants.ErrPostNotFound {
			return nil, constants.ErrPostNotFound
		}
		s.adapter.LogModule.Log(err, req, "[GetPost] failed get post")
		return nil, constants.ErrInternalServerError
	}

	return &api.GetPostRes{Post: *post}, nil
}

func (s *service) GetPostList(ctx context.Context, req api.GetPostListReq) (*api.GetPostListRes, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	postList, pagination, err := s.adapter.PostModule.FindPostListByParentIDAndAuthorIDs(ctx,
		req.ParentID,
		req.AuthorIDs,
		req.UserID,
		req.Pagination)
	if err != nil {
		s.adapter.LogModule.Log(err, req, "[GetPostList] failed get post list")
		return nil, constants.ErrInternalServerError
	}

	return &api.GetPostListRes{
		PostList:   postList,
		Pagination: *pagination,
	}, nil
}

func (s *service) CreatePost(ctx context.Context, req api.CreatePostReq) (*api.CreatePostRes, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.adapter.UserModule.FindUserByID(ctx, req.UserID); err != nil {
		if err == constants.ErrUserNotFound {
			return nil, constants.ErrUserNotFound
		}
		s.adapter.LogModule.Log(err, req, "[CreatePost] failed get user")
		return nil, constants.ErrInternalServerError
	}

	if err := s.adapter.PostModule.SavePost(ctx, entity.Post{
		Body:         req.Body,
		RepliesCount: 0,
		LikesCount:   0,
		Parent:       &entity.Post{ID: req.ParentID},
		Author:       &entity.User{ID: req.AuthorID},
		User:         entity.User{ID: req.UserID},
	}); err != nil {
		s.adapter.LogModule.Log(err, req, "[CreatePost] failed save post")
		return nil, constants.ErrInternalServerError
	}

	return &api.CreatePostRes{Message: "success"}, nil
}

func (s *service) LikePost(ctx context.Context, req api.LikePostReq) (*api.LikePostRes, error) {
	panic("implement me")
}
