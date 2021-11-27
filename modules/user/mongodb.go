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
		Name:     user.Profile.Name,
		Avatar:   user.Profile.Avatar,
		Bio:      user.Profile.Bio,
		Type:     "user",
		Follow:   &[]primitive.ObjectID{},
	}
	if err := u.user.Save(ctx, id, model); err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (u *userModule) SaveProfile(ctx context.Context, user entity.User) error {
	model := User{
		ID:       mongolib.ObjectID(user.ID),
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
	var model Users
	if err := u.user.Aggregate().
		Match(mongolib.Filter().Equal("type", "menfess")).
		Lookup("menfess", "_id", "user_id", "menfess").
		Unwind("$menfess").
		AddField("ordering", "$menfess.ordering").
		Sort("ordering", mongolib.Ascending).
		Exec(ctx).Consume(&model); err != nil {
		return nil, err
	}
	return model.Entity(), nil
}

func (u *userModule) GetFollowedUserID(ctx context.Context, userID string) ([]string, error) {
	var model User
	if err := u.user.FindByID(ctx, mongolib.ObjectID(userID)).Consume(&model); err != nil {
		return nil, err
	}
	if model.Follow == nil{
		model.Follow = &[]primitive.ObjectID{}
	}
	result := make([]string, len(*model.Follow))
	for i, v := range *model.Follow {
		result[i] = v.Hex()
	}
	return result, nil
}

func (u *userModule) UpdateFollowStatus(ctx context.Context, follower, followed, status string) error {
	if status == constants.FollowStatusActive {
		if err := u.user.Query().
			Equal("_id", mongolib.ObjectID(follower)).
			Push("follow", mongolib.ObjectID(followed)).
			Update(ctx); err != nil {
				return err
		}
	} else {
		if err := u.user.Query().
			Equal("_id", mongolib.ObjectID(follower)).
			Pull("follow", mongolib.ObjectID(followed)).
			Update(ctx); err != nil {
			return err
		}
	}
	return nil
}

type User struct {
	ID     primitive.ObjectID    `bson:"_id"`
	Name   string                `bson:"name"`
	Avatar string                `bson:"avatar"`
	Bio    string                `bson:"bio"`
	Type   string                `bson:"type"`
	Follow *[]primitive.ObjectID `bson:"follow,omitempty"`
}

func (u User) Entity() *entity.User {
	return &entity.User{
		ID:      u.ID.Hex(),
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