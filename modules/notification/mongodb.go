package notification

import (
	"context"
	"github.com/aeramu/mongolib"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *notificationModule) findPushToken(ctx context.Context, userID string) ([]string, error){
	var models []PushToken
	err := m.pushToken.Query().Equal("user_id", mongolib.ObjectID(userID)).Find(ctx).Consume(&models)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(models))
	for i, v := range models {
		result[i] = v.Token
	}

	return result, nil
}

func (m *notificationModule) findAllPushToken(ctx context.Context) ([]PushToken, error) {
	var model []PushToken
	err := m.pushToken.Query().Find(ctx).Consume(&model)
	if err != nil {
		return nil, err
	}

	return model, err
}

func (m *notificationModule) insertPushToken(ctx context.Context, userID string, token string) error {
	var model PushToken
	err := m.pushToken.Query().
		Equal("user_id", mongolib.ObjectID(userID)).
		Equal("token", token).
		FindOne(ctx).Consume(&model)
	if err == nil {
		return nil
	}
	if err != mongolib.ErrNotFound{
		return err
	}

	model.ID = mongolib.NewObjectID()
	model.Token = token
	model.UserID = mongolib.ObjectID(userID)
	if err := m.pushToken.Save(ctx, model.ID, model); err != nil {
		return err
	}

	return nil
}

func (m *notificationModule) removePushToken(ctx context.Context, userID string, token string) error {
	err := m.oldPushToken.Query().
		Equal("user_id", mongolib.ObjectID(userID)).
		Equal("token", token).
		DeleteOne(ctx)
	if err != nil {
		return err
	}

	return nil
}

type PushToken struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserID primitive.ObjectID `bson:"user_id"`
	Token  string             `bson:"token"`
}