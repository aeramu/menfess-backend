package service

import (
	"context"
	"github.com/aeramu/menfess-backend/constants"
	"github.com/aeramu/menfess-backend/entity"
	"github.com/aeramu/menfess-backend/service/api"
)

type Service interface {
	// User
	Login(ctx context.Context, req api.LoginReq) (*api.LoginRes, error)
	Register(ctx context.Context, req api.RegisterReq) (*api.RegisterRes, error)
	Logout(ctx context.Context, req api.LogoutReq) (*api.LogoutRes, error)
	UpdateProfile(ctx context.Context, req api.UpdateProfileReq) (*api.UpdateProfileRes, error)
	GetUser(ctx context.Context, req api.GetUserReq) (*api.GetUserRes, error)
	FollowUser(ctx context.Context, req api.FollowUserReq) (*api.FollowUserRes, error)
	GetMenfessList(ctx context.Context, req api.GetMenfessListReq) (*api.GetMenfessListRes, error)

	// Post
	Feed(ctx context.Context, req api.FeedReq) (*api.FeedRes, error)
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
	user := entity.User{}
	id, err := s.adapter.UserModule.InsertUser(ctx, user)
	if err != nil {
		s.adapter.LogModule.Log(err, req, "[Register] failed insert user")
		return nil, constants.ErrInternalServerError
	}
	user.ID = id

	if err := s.adapter.NotificationModule.AddPushToken(ctx, user.ID, req.PushToken); err != nil {
		s.adapter.LogModule.Log(err, req, "[Login] failed add notification push token")
		return nil, constants.ErrInternalServerError
	}

	token, err := s.adapter.AuthModule.GenerateToken(ctx, user)
	if err != nil {
		s.adapter.LogModule.Log(err, req, "[Register] failed generate token")
		return nil, constants.ErrInternalServerError
	}

	return &api.RegisterRes{Token: token}, nil
}

func (s *service) Logout(ctx context.Context, req api.LogoutReq) (*api.LogoutRes, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if err := s.adapter.NotificationModule.RemovePushToken(ctx, req.UserID, req.PushToken); err != nil {
		s.adapter.LogModule.Log(err, req, "[Logout] Failed remove push token")
		return nil, constants.ErrInternalServerError
	}

	return &api.LogoutRes{Message: "Success"}, nil
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

	followed, err := s.adapter.UserModule.GetFollowedUserID(ctx, req.UserID)
	if err != nil {
		s.adapter.LogModule.Log(err, req, "[FollowUser] failed get followed user")
		return nil, constants.ErrInternalServerError
	}

	return &api.GetMenfessListRes{MenfessList: menfessList, FollowedIDs: followed}, nil
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
		nil,
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

	user, err := s.adapter.UserModule.FindUserByID(ctx, req.UserID)
	if err != nil {
		if err == constants.ErrUserNotFound {
			return nil, constants.ErrUserNotFound
		}
		s.adapter.LogModule.Log(err, req, "[CreatePost] failed get user")
		return nil, constants.ErrInternalServerError
	}

	post := entity.Post{
		Body:         req.Body,
		Parent:       &entity.Post{ID: req.ParentID},
		Author:       entity.User{ID: req.AuthorID},
		User:         *user,
	}
	id, err := s.adapter.PostModule.InsertPost(ctx, post)
	if err != nil {
		s.adapter.LogModule.Log(err, req, "[CreatePost] failed save post")
		return nil, constants.ErrInternalServerError
	}
	post.ID = id

	if req.ParentID != "" {
		parent, err := s.adapter.PostModule.FindPostByID(ctx, req.ParentID, "")
		if err != nil {
			s.adapter.LogModule.Log(err, req, "[CreatePost] failed get post")
			return nil, constants.ErrInternalServerError
		}
		if err := s.adapter.NotificationModule.SendCommentNotification(ctx, post, *parent); err != nil {
			s.adapter.LogModule.Log(err, req, "[CreatePost] failed send notification")
		}
	} else {
		if err := s.adapter.NotificationModule.BroadcastNewPostNotification(ctx, post); err != nil {
			s.adapter.LogModule.Log(err, req, "[CreatePost] failed send notification")
		}
	}

	return &api.CreatePostRes{Message: "success"}, nil
}

func (s *service) LikePost(ctx context.Context, req api.LikePostReq) (*api.LikePostRes, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	user, err := s.adapter.UserModule.FindUserByID(ctx, req.UserID)
	if err != nil {
		if err == constants.ErrUserNotFound {
			return nil, constants.ErrUserNotFound
		}
		s.adapter.LogModule.Log(err, req, "[LikePost] failed get user")
		return nil, constants.ErrInternalServerError
	}

	post, err := s.adapter.PostModule.FindPostByID(ctx, req.PostID, req.UserID)
	if err != nil {
		if err == constants.ErrPostNotFound {
			return nil, constants.ErrPostNotFound
		}
		s.adapter.LogModule.Log(err, req, "[LikePost] failed get post")
		return nil, constants.ErrInternalServerError
	}

	if post.IsLiked {
		if err := s.adapter.PostModule.UnlikePost(ctx, req.PostID, req.UserID); err != nil {
			s.adapter.LogModule.Log(err, req, "[LikePost] failed unlike post")
			return nil, constants.ErrInternalServerError
		}
	} else {
		if err := s.adapter.PostModule.LikePost(ctx, req.PostID, req.UserID); err != nil {
			s.adapter.LogModule.Log(err, req, "[LikePost] failed like post")
			return nil, constants.ErrInternalServerError
		}
		if err := s.adapter.NotificationModule.SendLikeNotification(ctx, *user, *post); err != nil {
			s.adapter.LogModule.Log(err, req, "[LikePost] failed send notification")
		}
	}

	return &api.LikePostRes{Message: "success"}, nil
}

func (s *service) FollowUser(ctx context.Context, req api.FollowUserReq) (*api.FollowUserRes, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	followed, err := s.adapter.UserModule.GetFollowedUserID(ctx, req.UserID)
	if err != nil {
		s.adapter.LogModule.Log(err, req, "[FollowUser] failed get followed user")
		return nil, constants.ErrInternalServerError
	}

	status := constants.FollowStatusActive
	for _, v := range followed {
		if req.FollowedID == v {
			status = constants.FollowStatusInactive
		}
	}

	if err := s.adapter.UserModule.UpdateFollowStatus(ctx, req.UserID, req.FollowedID, status); err != nil {
		s.adapter.LogModule.Log(err, req, "[FollowUser] failed update follow status")
		return nil, constants.ErrInternalServerError
	}
	return &api.FollowUserRes{Message: "success"}, nil
}

func (s *service) Feed(ctx context.Context, req api.FeedReq) (*api.FeedRes, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var followed []string
	if req.Type == constants.FeedTypeFollow {
		userIDs, err := s.adapter.UserModule.GetFollowedUserID(ctx, req.UserID)
		if err != nil {
			s.adapter.LogModule.Log(err, req, "[Feed] failed get followed user")
			return nil, constants.ErrInternalServerError
		}
		if len(userIDs) == 0 {
			return nil, constants.ErrUserNotFollowAnyone
		}
		followed = userIDs
	}

	postList, pagination, err := s.adapter.PostModule.FindPostListByParentIDAndAuthorIDs(ctx,
		"",
		followed,
		req.UserID,
		req.Pagination)
	if err != nil {
		s.adapter.LogModule.Log(err, req, "[GetPostList] failed get post list")
		return nil, constants.ErrInternalServerError
	}

	return &api.FeedRes{
		PostList:   postList,
		Pagination: *pagination,
	}, nil
}
