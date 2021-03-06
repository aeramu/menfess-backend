package service

import (
	"context"
	"github.com/aeramu/menfess-backend/entity"
	"github.com/aeramu/menfess-backend/service/api"
)

type Adapter struct {
	UserModule         UserModule
	PostModule         PostModule
	AuthModule         AuthModule
	NotificationModule NotificationModule
	LogModule          LogModule
}

type AuthModule interface {
	GenerateToken(ctx context.Context, user entity.User) (string, error)
	ComparePassword(ctx context.Context, hash string, password string) error
	HashPassword(ctx context.Context, password string) (string, error)
}

type UserModule interface {
	FindUserByID(ctx context.Context, id string) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	InsertUser(ctx context.Context, user entity.User) (string, error)
	SaveProfile(ctx context.Context, user entity.User) error
	FindMenfessList(ctx context.Context) ([]entity.User, error)
	GetFollowedUserID(ctx context.Context, userID string) ([]string, error)
	UpdateFollowStatus(ctx context.Context, follower, followed, status string) error
}

type PostModule interface {
	FindPostByID(ctx context.Context, id string, userID string) (*entity.Post, error)
	FindPostListByParentIDAndAuthorIDs(ctx context.Context,
		parentID string,
		authorIDs []string,
		userID string,
		pagination api.PaginationReq,
	) ([]entity.Post, *api.PaginationRes, error)
	InsertPost(ctx context.Context, post entity.Post) (string, error)
	LikePost(ctx context.Context, postID string, userID string) error
	UnlikePost(ctx context.Context, postID string, userID string) error
}

type NotificationModule interface {
	AddPushToken(ctx context.Context, userID string, pushToken string) error
	RemovePushToken(ctx context.Context, userID string, pushToken string) error
	SendLikeNotification(ctx context.Context, user entity.User, post entity.Post) error
	SendCommentNotification(ctx context.Context, comment entity.Post, parent entity.Post) error
	BroadcastNewPostNotification(ctx context.Context, post entity.Post) error
}

type LogModule interface {
	Log(err error, payload interface{}, message string)
}
