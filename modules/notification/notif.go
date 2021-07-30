package notification

import (
	"context"
	"fmt"
	"github.com/aeramu/menfess-backend/constants"
	"github.com/aeramu/menfess-backend/entity"
	"github.com/aeramu/menfess-backend/service"
	"github.com/aeramu/mongolib"
)

func NewNotificationModule(db *mongolib.Database) service.NotificationModule {
	return &notificationModule{
		pushToken: db.Coll("push_token"),
	}
}

type notificationModule struct {
	pushToken *mongolib.Collection
}

const(
	likeNotificationTitle = "%s like your post"
	commentNotificationTitle = "%s comment on your post"
	newPostNotificationTitle = "%s post a menfess"
)

func (m *notificationModule) AddPushToken(ctx context.Context, userID string, pushToken string) error {
	if err := m.insertPushToken(ctx, userID, pushToken); err != nil {
		return err
	}
	return nil
}

func (m *notificationModule) RemovePushToken(ctx context.Context, userID string, pushToken string) error {
	if err := m.removePushToken(ctx, userID, pushToken); err != nil {
		if err == mongolib.ErrNotFound {
			return constants.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (m *notificationModule) SendLikeNotification(ctx context.Context, user entity.User, post entity.Post) error {
	tokens, err := m.findPushToken(ctx, post.User.ID)
	if err != nil {
		return err
	}

	if err := m.sendNotification(
		tokens,
		fmt.Sprintf(likeNotificationTitle, user.Profile.Name),
		post.Body,
		Data{PostID: post.ID},
	); err != nil {
		return err
	}

	return nil
}

func (m *notificationModule) SendCommentNotification(ctx context.Context, comment entity.Post, parent entity.Post) error {
	tokens, err := m.findPushToken(ctx, parent.User.ID)
	if err != nil {
		return err
	}

	if err := m.sendNotification(
		tokens,
		fmt.Sprintf(commentNotificationTitle, comment.User.Profile.Name),
		comment.Body,
		Data{PostID: parent.ID},
	); err != nil {
		return err
	}

	return nil
}

func (m *notificationModule) BroadcastNewPostNotification(ctx context.Context, post entity.Post) error {
	tokens, err := m.findAllPushToken(ctx)
	if err != nil {
		return err
	}

	//if utils.RandomChance(50, 100) {
	//	return nil
	//}

	if err := m.sendNotification(
		tokens,
		fmt.Sprintf(newPostNotificationTitle, post.User.Profile.Name),
		post.Body,
		Data{PostID: post.ID},
	); err != nil {
		return err
	}

	return nil
}

type Data struct {
	PostID string `json:"postID"`
}