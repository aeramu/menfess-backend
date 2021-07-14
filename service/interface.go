package service

import (
	"context"
	"github.com/aeramu/menfess-backend/entity"
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
}

type PostModule interface {
}

type NotificationModule interface {
	AddPushToken(ctx context.Context, userID string, pushToken string) error
}

type LogModule interface {
	Log(err error, payload interface{}, message string)
}
