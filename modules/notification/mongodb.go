package notification

import (
	"context"
	"github.com/aeramu/mongolib"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *notificationModule) findPushToken(ctx context.Context, userID string) ([]string, error){
	var model PushToken
	err := m.pushToken.Query().Equal("_id", mongolib.ObjectID(userID)).FindOne(ctx).Consume(&model)
	if err != nil {
		if err == mongolib.ErrNotFound {
			return []string{}, nil
		}
		return nil, err
	}

	return convertMapStringToArray(model.Token), nil
}

func (m *notificationModule) findAllPushToken(ctx context.Context) ([]string, error) {
	var model []PushToken
	err := m.pushToken.Query().Find(ctx).Consume(&model)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, v := range model {
		result = append(result, convertMapStringToArray(v.Token)...)
	}
	return result, err
}

func (m *notificationModule) insertPushToken(ctx context.Context, userID string, token string) error {
	var model PushToken
	err := m.pushToken.Query().Equal("_id", mongolib.ObjectID(userID)).FindOne(ctx).Consume(&model)
	if err != nil {
		if err == mongolib.ErrNotFound {
			model.ID = mongolib.ObjectID(userID)
			model.Token = map[string]bool{}
		} else {
			return err
		}
	}

	_, ok := model.Token[token]
	if !ok {
		model.Token[token] = true
		if err := m.pushToken.Save(ctx, model.ID, model); err != nil {
			return err
		}
	}

	return nil
}

func (m *notificationModule) removePushToken(ctx context.Context, userID string, token string) error {
	var model PushToken
	err := m.pushToken.Query().Equal("_id", mongolib.ObjectID(userID)).FindOne(ctx).Consume(&model)
	if err != nil {
		return err
	}

	_, ok := model.Token[token]
	if ok {
		delete(model.Token, token)
		if err := m.pushToken.Save(ctx, model.ID, model); err != nil {
			return err
		}
	}

	return nil
}

type PushToken struct {
	ID    primitive.ObjectID `bson:"_id"`
	Token map[string]bool    `bson:"token"`
}

func convertMapStringToArray(m map[string]bool) (arr []string) {
	for key, _ := range m {
		if key != "" {
			arr = append(arr, key)
		}
	}
	return
}