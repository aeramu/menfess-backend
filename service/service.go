package service

import (
	"context"
	"github.com/aeramu/menfess-backend/constants"
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
	panic("implement me")
}

func (s *service) UpdateProfile(ctx context.Context, req api.UpdateProfileReq) (*api.UpdateProfileRes, error) {
	panic("implement me")
}

func (s *service) GetUser(ctx context.Context, req api.GetUserReq) (*api.GetUserRes, error) {
	panic("implement me")
}

func (s *service) GetMenfessList(ctx context.Context, req api.GetMenfessListReq) (*api.GetMenfessListRes, error) {
	panic("implement me")
}

func (s *service) GetPost(ctx context.Context, req api.GetPostReq) (*api.GetPostRes, error) {
	panic("implement me")
}

func (s *service) GetPostList(ctx context.Context, req api.GetPostListReq) (*api.GetPostListRes, error) {
	panic("implement me")
}

func (s *service) CreatePost(ctx context.Context, req api.CreatePostReq) (*api.CreatePostRes, error) {
	panic("implement me")
}

func (s *service) LikePost(ctx context.Context, req api.LikePostReq) (*api.LikePostRes, error) {
	panic("implement me")
}
