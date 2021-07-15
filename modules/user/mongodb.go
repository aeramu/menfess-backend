package user

import (
	"context"
	"github.com/aeramu/menfess-backend/constants"
	"github.com/aeramu/menfess-backend/entity"
	"github.com/aeramu/menfess-backend/service"
	"github.com/aeramu/mongolib"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewUserModule(db *mongolib.Database) service.UserModule {
	return &userModule{user: db.Coll("user")}
}

type userModule struct {
	user *mongolib.Collection
}

func (u *userModule) FindUserByID(ctx context.Context, id string) (*entity.User, error) {
	var model User
	if err := u.user.FindByID(ctx, mongolib.ObjectID(id)).Consume(&model); err != nil {
		if err == mongolib.ErrNotFound {
			return nil, constants.ErrUserNotFound
		}
		return nil, err
	}
	return model.Entity(), nil
}

func (u *userModule) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var model User
	if err := u.user.Query().Equal("email", email).FindOne(ctx).Consume(&model); err != nil {
		if err == mongolib.ErrNotFound {
			return nil, constants.ErrUserNotFound
		}
		return nil, err
	}
	return model.Entity(), nil
}

func (u *userModule) InsertUser(ctx context.Context, user entity.User) (string, error) {
	id := mongolib.NewObjectID()
	model := User{
		ID:       id,
		Email:    user.Account.Email,
		Password: user.Account.Password,
		Name:     user.Profile.Name,
		Avatar:   user.Profile.Avatar,
		Bio:      user.Profile.Bio,
		Type:     "user",
	}
	if err := u.user.Save(ctx, id, model); err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (u *userModule) SaveProfile(ctx context.Context, user entity.User) error {
	model := User{
		ID:       mongolib.ObjectID(user.ID),
		Email:    user.Account.Email,
		Password: user.Account.Password,
		Name:     user.Profile.Name,
		Avatar:   user.Profile.Avatar,
		Bio:      user.Profile.Bio,
		Type:     "user",
	}
	if err := u.user.Save(ctx, mongolib.ObjectID(user.ID), model); err != nil {
		return err
	}
	return nil
}

func (u *userModule) FindMenfessList(ctx context.Context) ([]entity.User, error) {
	model := Users{}
	if err := u.user.Query().Equal("type", "menfess").Find(ctx).Consume(model); err != nil {
		return nil, err
	}
	return model.Entity(), nil
}

type User struct {
	ID primitive.ObjectID `bson:"_id"`
	Email string `bson:"email"`
	Password string `bson:"password"`
	Name string `bson:"name"`
	Avatar string `bson:"avatar"`
	Bio string `bson:"bio"`
	Type string `bson:"type"`
}

func (u User) Entity() *entity.User {
	return &entity.User{
		ID:      u.ID.Hex(),
		Account: entity.Account{
			Email:    u.Email,
			Password: u.Password,
		},
		Profile: entity.Profile{
			Name:   u.Name,
			Avatar: u.Avatar,
			Bio:    u.Bio,
		},
	}
}

type Users []User

func (u Users) Entity() []entity.User {
	var entities []entity.User
	for _, v := range u {
		entities = append(entities, *v.Entity())
	}
	return entities
}